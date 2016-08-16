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
	"github.com/gocraft/web"

	httpGoCommon "github.com/trustedanalytics/tapng-go-common/http"
	"github.com/trustedanalytics/tapng-image-factory/app"
	"os"
)

func main() {
	context := app.Context{}
	context.SetupContext()

	/* Queue Handling */
	go app.StartConsumer(context)
	/* Queue Handling */

	/* REST API handling */
	router := web.New(context)
	router.Middleware(web.LoggerMiddleware)

	router.Get("/healthz", context.GetImageFactoryHealth)

	apiRouter := router.Subrouter(context, "/api/v1")
	route(apiRouter, &context)
	v1AliasRouter := router.Subrouter(context, "/api/v1.0")
	route(v1AliasRouter, &context)

	if os.Getenv("IMAGE_FACTORY_SSL_CERT_FILE_LOCATION") != "" {
		httpGoCommon.StartServerTLS(os.Getenv("IMAGE_FACTORY_SSL_CERT_FILE_LOCATION"),
			os.Getenv("IMAGE_FACTORY_SSL_CERT_FILE_LOCATION"), router)
	} else {
		httpGoCommon.StartServer(router)
	}
	/* REST API handling */
}

func route(router *web.Router, context *app.Context) {
	router.Post("/image", (*context).BuildImage)
}