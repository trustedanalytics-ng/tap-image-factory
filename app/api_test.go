/**
 * Copyright (c) 2016 Intel Corporation
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
package main

import (
	"testing"
	"net/http/httptest"
	"strings"

	"github.com/gocraft/web"
	. "github.com/smartystreets/goconvey/convey"
)

const URLbuildImage = "/rest/v1/app"

func prepareMocksAndRouter(t *testing.T) (router *web.Router) {
	router = web.New(Context{})
	return router
}

func TestBuildImage(t *testing.T) {
	request := BuildImagePostRequest{ApplicationId: "test-app-id"}
	router := prepareMocksAndRouter(t)
	router.Post(URLbuildImage, (*Context).BuildImage)

	Convey("Test BuildImage", t, func() {
		Convey("Should returns proper response", func() {
			response := sendRequest("POST", URLbuildImage, marshallToJson(t, request), router)
			assertResponse(response, "", 201)
		})
	})
}

func assertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		So(strings.TrimSpace(string(rr.Body.Bytes())), ShouldContainSubstring, body)
	}
	So(rr.Code, ShouldEqual, code)
}
