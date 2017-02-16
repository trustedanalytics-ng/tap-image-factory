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
	"testing"
	"time"

	"github.com/gocraft/web"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPing(t *testing.T) {
	Convey("Given that docker-registry server returns successful ping immediately", t, func() {
		registry, err := NewDockerRegistry(getFakeServerURL())

		Convey("No error should be returned, registry client should be initialized", func() {
			So(err, ShouldBeNil)
			So(registry, ShouldNotBeNil)
		})
	})

	Convey("Given that docker-registry server returns bad response on ping immediately", t, func() {
		mockedServerCtx := RegistryServerContext{
			getPingMockedFunc: func(rw web.ResponseWriter) {
				rw.WriteHeader(http.StatusInternalServerError)
			},
		}

		_, err := NewDockerRegistry(getFakeServerURLWithCustomContext(mockedServerCtx))

		Convey("Error should be returned", func() {
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given that server would answer too long on ping", t, func() {
		mockedServerCtx := RegistryServerContext{
			getPingMockedFunc: func(rw web.ResponseWriter) {
				time.Sleep(3 * time.Second)
			},
		}

		_, err := NewDockerRegistryWithCustomConnectionTimeout(getFakeServerURLWithCustomContext(mockedServerCtx), time.Millisecond)

		Convey("Error should be returned", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestIsImageReady(t *testing.T) {
	fakeImageID := "image-id"
	fakeTag := "image-tag"

	Convey("Given that docker-registry server responds OK for getting manifest", t, func() {
		mockedServerCtx := RegistryServerContext{
			getManifestMockedFunc: func(rw web.ResponseWriter, id, tag string) {
				if id != fakeImageID || tag != fakeTag {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				rw.WriteHeader(http.StatusOK)
			},
		}

		registry, err := NewDockerRegistry(getFakeServerURLWithCustomContext(mockedServerCtx))
		So(err, ShouldBeNil)

		readiness, err := registry.IsImageReady(fakeImageID, fakeTag)
		Convey("IsImageReady should return true and no error", func() {
			So(readiness, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given that docker-registry server responds Not Found for getting manifest", t, func() {
		mockedServerCtx := RegistryServerContext{
			getManifestMockedFunc: func(rw web.ResponseWriter, id, tag string) {
				if id != fakeImageID || tag != fakeTag {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				rw.WriteHeader(http.StatusNotFound)
			},
		}

		registry, err := NewDockerRegistry(getFakeServerURLWithCustomContext(mockedServerCtx))
		So(err, ShouldBeNil)

		readiness, err := registry.IsImageReady(fakeImageID, fakeTag)
		Convey("IsImageReady should return false and no error", func() {
			So(readiness, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given that docker-registry server responds Internal Server Error for getting manifest", t, func() {
		mockedServerCtx := RegistryServerContext{
			getManifestMockedFunc: func(rw web.ResponseWriter, id, tag string) {
				if id != fakeImageID || tag != fakeTag {
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				rw.WriteHeader(http.StatusInternalServerError)
			},
		}

		registry, err := NewDockerRegistry(getFakeServerURLWithCustomContext(mockedServerCtx))
		So(err, ShouldBeNil)

		readiness, err := registry.IsImageReady(fakeImageID, fakeTag)
		Convey("IsImageReady should return false and some error", func() {
			So(readiness, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
	})
}
