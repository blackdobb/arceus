/*
Copyright © 2021 zc2638 <zc2638@qq.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package resource

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zc2638/arceus/pkg/types"

	"github.com/pkgms/go/ctr"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/zc2638/arceus/global"
	"github.com/zc2638/arceus/static"
)

func list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typ := r.URL.Query().Get("type")
		var (
			result []types.Resource
			err    error
		)
		if typ == "custom" {
			result, err = types.BuildResourceList(os.DirFS(global.ResourcePath), "custom")
		} else {
			result, err = types.BuildResourceList(static.Kubernetes, static.KubernetesDir)
		}
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.OK(w, result)
	}
}

func info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		kind := r.URL.Query().Get("kind")
		version := r.URL.Query().Get("version")

		filePath := filepath.Join(group, kind, version) + ".yaml"
		baseFilePath := filepath.Join(static.KubernetesDir, filePath)
		fileData, err := fs.ReadFile(static.Kubernetes, baseFilePath)
		if os.IsNotExist(err) {
			fileData, err = fs.ReadFile(os.DirFS(global.CustomResourcePath), filePath)
		}
		if err != nil {
			ctr.BadRequest(w, errors.New("resource not exist"))
			return
		}
		ctr.Bytes(w, fileData)
	}
}

func tree() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		kind := r.URL.Query().Get("kind")
		version := r.URL.Query().Get("version")

		filePath := filepath.Join(group, kind, version) + ".yaml"
		baseFilePath := filepath.Join(static.KubernetesDir, filePath)
		fileData, err := fs.ReadFile(static.Kubernetes, baseFilePath)
		if os.IsNotExist(err) {
			fileData, err = fs.ReadFile(os.DirFS(global.CustomResourcePath), filePath)
		}
		if err != nil {
			ctr.BadRequest(w, errors.New("resource not exist"))
			return
		}
		// 解析到结构体
		var data apiextensionsV1.CustomResourceDefinition
		if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), fileData, &data); err != nil {
			ctr.InternalError(w, fmt.Errorf("resource parse failed: %s", err))
			return
		}
		// 转换为自定义的树结构
		if len(data.Spec.Versions) == 0 {
			ctr.OK(w, struct{}{})
			return
		}
		apiVersion := data.Spec.Group + "/" + data.Spec.Versions[0].Name
		node := types.BuildNode(data.Spec.Versions[0].Schema.OpenAPIV3Schema, nil, apiVersion, data.Spec.Names.Kind)
		ctr.OK(w, node)
	}
}
