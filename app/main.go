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
	"net/http"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/image-factory/logger"
)

type Context struct{}

var (
	logger = logger_wrapper.InitLogger("main")
	port = "8080"
)

func main() {
	r := web.New(Context{})
	r.Post("/rest/v1/app", (*Context).BuildImage)

	err := http.ListenAndServe("localhost"+":"+port, r)
	if err != nil {
		logger.Critical("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
