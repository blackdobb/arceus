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
package home

import (
	"net/http"

	"github.com/zc2638/arceus/pkg/version"

	"github.com/pkgms/go/ctr"

	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/swagger"
)

const tag = "home"

// Route handle home routing related
func Route(doc *swagger.API) {
	doc.Tags = append(doc.Tags, swagger.Tag{
		Name:        tag,
		Description: "主页",
	})
	doc.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/",
			endpoint.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctr.Found(w, r, "/web")
				//ctr.OK(w, "Hello Arceus!")
			})),
			endpoint.ResponseSuccess(),
			endpoint.Tags(tag),
		),
		endpoint.New(
			http.MethodGet, "/version",
			endpoint.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctr.OK(w, version.Get())
			})),
			endpoint.ResponseSuccess(endpoint.Schema(version.Version{})),
			endpoint.Tags(tag),
		),
	)
}
