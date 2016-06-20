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
	"github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

const (
	catalogAddress   = "https://test-catalog.org"
	applicationId    = "25de06b4-ac21-4454-bb82-e72bc05f3a5c"
	applicationsPath = "/applications/"
)

func TestGetApplicationDetails(t *testing.T) {

	os.Setenv("CATALOG_ADDRESS", catalogAddress)
	c := NewCatalogConnector()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	print(GetCatalogAddress() + applicationsPath + applicationId)
	httpmock.RegisterResponder("GET", GetCatalogAddress()+applicationsPath+applicationId,
		httpmock.NewStringResponder(200, `{"applicationId":"`+applicationId+`"}`))
	c.Client.Transport = httpmock.DefaultTransport
	res, err := c.GetApplicationDetails(applicationId)
	convey.ShouldBeNil(err)
	convey.ShouldEqual(res.ApplicationId, applicationId)
}

func TestUpdateApplicationState(t *testing.T) {

}

func TestGetBlob(t *testing.T) {
}
