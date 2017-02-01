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
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildImage(t *testing.T) {
	mockCtrl, mocks, client, _ := prepareMocksAndClient(t)
	Convey("Test BuildImage", t, func() {
		Convey("Given that BuildAndPushImage invokes successfully", func() {
			sampleImageId := "sampleID"

			mocks.MockFactory.EXPECT().BuildAndPushImage(sampleImageId).Return(nil)

			err := client.BuildImage(sampleImageId)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})

		})

		//BuildImage is asynchronous, so client doesn't return error even dependent function fails
		Convey("Given that BuildAndPushImage returns error", func() {
			sampleImageId := "sampleID"

			mocks.MockFactory.EXPECT().BuildAndPushImage(sampleImageId).Return(errors.New("some error"))

			err := client.BuildImage(sampleImageId)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}
