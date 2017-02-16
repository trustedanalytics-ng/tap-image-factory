/**
 * Copyright (c) 2017 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"net/http"
	"net/http/httptest"

	"github.com/gocraft/web"
)

type RegistryServerContext struct {
	getPingMockedFunc     func(rw web.ResponseWriter)
	getManifestMockedFunc func(web.ResponseWriter, string, string)
}

func getFakeServerURL() string {
	context := RegistryServerContext{}
	return getFakeServerURLWithCustomContext(context)
}

func getFakeServerURLWithCustomContext(serverContext RegistryServerContext) string {
	router := web.New(serverContext)
	router.Get("/v2/", serverContext.GetPing)
	router.Get("/v2/:imageid/manifests/:tag", serverContext.GetManifest)

	testServer := httptest.NewServer(router)
	return testServer.URL
}

func (c *RegistryServerContext) GetPing(rw web.ResponseWriter, req *web.Request) {
	if c.getPingMockedFunc != nil {
		c.getPingMockedFunc(rw)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func (c *RegistryServerContext) GetManifest(rw web.ResponseWriter, req *web.Request) {
	if c.getManifestMockedFunc != nil {
		imageID := req.PathParams["imageid"]
		tag := req.PathParams["tag"]
		c.getManifestMockedFunc(rw, imageID, tag)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}
