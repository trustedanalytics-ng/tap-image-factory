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
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gocraft/web"
	"github.com/golang/mock/gomock"

	"github.com/trustedanalytics-ng/tap-image-factory/client"
)

type MockPack struct {
	MockBlobStoreConnector     *MockTapBlobStoreApi
	MockTapCatalogApiConnector *MockTapCatalogApi
	MockDockerConnector        *MockImageBuilder
	MockFactory                *MockFactoryAPI
	MockReader                 *MockArchiveReader
}

func prepareMocksAndClient(t *testing.T) (mockCtrl *gomock.Controller, mocks MockPack, client client.TapApiImageFactoryApi, router *web.Router) {
	mockCtrl = gomock.NewController(t)

	mocks = MockPack{
		MockBlobStoreConnector:     NewMockTapBlobStoreApi(mockCtrl),
		MockTapCatalogApiConnector: NewMockTapCatalogApi(mockCtrl),
		MockDockerConnector:        NewMockImageBuilder(mockCtrl),
		MockFactory:                NewMockFactoryAPI(mockCtrl),
		MockReader:                 NewMockArchiveReader(mockCtrl),
	}

	c := Context{
		BlobStoreConnector:     mocks.MockBlobStoreConnector,
		TapCatalogApiConnector: mocks.MockTapCatalogApiConnector,
		DockerConnector:        mocks.MockDockerConnector,
		Factory:                mocks.MockFactory,
		Reader:                 mocks.MockReader,
	}
	ctx = &c

	router = SetupRouter(&c)
	client = getImageFactoryClient(router, t)
	return
}

func getImageFactoryClient(router *web.Router, t *testing.T) client.TapApiImageFactoryApi {
	const user = "user"
	const password = "password"

	os.Setenv("IMAGE_FACTORY_USER", user)
	os.Setenv("IMAGE_FACTORY_PASS", password)

	testServer := httptest.NewServer(router)
	catalogClient, err := client.NewTapImageFactoryApiWithBasicAuth(testServer.URL, user, password)
	if err != nil {
		t.Fatal("ImageFactory client error: ", err)
	}
	return catalogClient
}
