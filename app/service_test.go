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

func TestGetApplicationDetails(t *testing.T) {
	os.Setenv("CATALOG_ADDRESS", catalogAddress)
	httpmock.Activate()
	c := NewCatalogConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test GetApplicationDetails", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("GET", GetCatalogAddress()+applicationsPath+applicationId,
				httpmock.NewStringResponder(200, `{"id":"`+applicationId+`"}`))
			res, err := c.GetApplicationDetails(applicationId)
			So(err, ShouldBeNil)
			So(res.ApplicationId, ShouldEqual, applicationId)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("GET", GetCatalogAddress()+applicationsPath+applicationId,
				httpmock.NewStringResponder(404, ``))
			res, err := c.GetApplicationDetails(applicationId)
			So(err.Error(), ShouldEqual, "Invalid status: 404")
			So(res.ApplicationId, ShouldBeEmpty)
		})
	})
}

func TestUpdateApplicationState(t *testing.T) {
	os.Setenv("CATALOG_ADDRESS", catalogAddress)
	httpmock.Activate()
	c := NewCatalogConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test UpdateApplicationState", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("PATCH", GetCatalogAddress()+applicationsPath+applicationId,
				httpmock.NewStringResponder(200, ``))
			err := c.UpdateApplicationState(applicationId, "created")
			So(err, ShouldBeNil)

		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("PATCH", GetCatalogAddress()+applicationsPath+applicationId,
				httpmock.NewStringResponder(404, ``))
			err := c.UpdateApplicationState(applicationId, "created")
			So(err.Error(), ShouldEqual, "Invalid status: 404")
		})
	})
}

func TestGetBlob(t *testing.T) {
	os.Setenv("BLOB_STORE_ADDRESS", blobStoreAddress)
	httpmock.Activate()
	c := NewBlobStoreConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test GetBlob", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("GET", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(200, blob))
			res, err := c.GetApplicationBlob(applicationId)
			So(err, ShouldBeNil)
			So(string(res), ShouldEqual, blob)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("GET", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(404, ""))
			res, err := c.GetApplicationBlob(applicationId)
			So(err.Error(), ShouldEqual, "Invalid status: 404")
			So(string(res), ShouldEqual, "")
		})
	})
}

func TestDeleteBlob(t *testing.T) {
	os.Setenv("BLOB_STORE_ADDRESS", blobStoreAddress)
	httpmock.Activate()
	c := NewBlobStoreConnector()
	c.Client.Transport = httpmock.DefaultTransport
	defer httpmock.DeactivateAndReset()

	Convey("Test DeleteBlob", t, func() {
		Convey("Should return proper response", func() {
			httpmock.RegisterResponder("DELETE", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(204, ""))
			err := c.DeleteApplicationBlob(applicationId)
			So(err, ShouldBeNil)
		})

		Convey("Should return not found", func() {
			httpmock.RegisterResponder("DELETE", GetBlobStoreAddress()+blobsPath+blobId,
				httpmock.NewStringResponder(404, ""))
			err := c.DeleteApplicationBlob(applicationId)
			So(err.Error(), ShouldEqual, "Invalid status: 404")
		})
	})
}
