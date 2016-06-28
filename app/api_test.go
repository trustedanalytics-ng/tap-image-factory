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
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

const URLbuildImage = "/api/v1/app"

func TestBuildImage(t *testing.T) {
	os.Setenv("CATALOG_ADDRESS", catalogAddress)
	httpmock.Activate()
	c := NewCatalogConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	request := BuildImagePostRequest{ApplicationId: applicationId}
	router := prepareMocksAndRouter(t)
	context := Context{}
	catalogConnectorMock := CatalogConnectorMock{}
	blobStoreConnectorMock := BlobStoreConnectorMock{}
	dockerClientMock := DockerClientMock{}

	context.BlobStoreConnector = &BlobStoreConnector{Api: &blobStoreConnectorMock}
	context.CatalogConnector = &CatalogConnector{Api: &catalogConnectorMock}
	context.DockerConnector = &DockerHandler{Api: &dockerClientMock}

	router.Post(URLbuildImage, context.BuildImage)

	Convey("Test BuildImage", t, func() {
		Convey("Should return proper response: 201", func() {
			response := sendRequest("POST", URLbuildImage, marshalToJson(t, request), router)
			assertResponse(response, "", 201)
		})
		Convey("Should return invalid body response: 400", func() {
			response := sendRequest("POST", URLbuildImage, []byte{}, router)
			assertResponse(response, "", 400)
		})
	})

}
