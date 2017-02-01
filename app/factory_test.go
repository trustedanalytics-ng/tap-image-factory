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
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics/tap-catalog/builder"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
)

type GomockCalls []*gomock.Call

type SimpleTestCaseDefinition struct {
	mockedCalls       func() GomockCalls
	testDescription   string
	shouldReturnError bool
}

func makePatchUserFriendly(oldState interface{}, newState interface{}, t *testing.T) []catalogModels.Patch {
	patch, err := builder.MakePatchWithPreviousValue("state", newState, oldState, catalogModels.OperationUpdate)
	if err != nil {
		t.Fatal(err)
	}
	return []catalogModels.Patch{patch}
}

func TestBuildAndPushImage(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)
	factory := Factory{}

	fakeImage := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeExec,
	}

	Convey("Test BuildAndPushImage", t, func() {
		Convey("Given that all underneath functions invokes successfully", func() {
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImage.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImage, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().GetImage(fakeImage.Id).Return(fakeImage, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImage.Id, gomock.Any()).Return(nil),
				mocks.MockDockerConnector.EXPECT().
					CreateImage(gomock.Any(), fakeImage.Type, fakeImage.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImage.Id)).
					Return(nil),
				mocks.MockDockerConnector.EXPECT().PushImage(GetImageWithHubAddressWithoutProtocol(fakeImage.Id)).Return(nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImage.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateReady, t)).
					Return(catalogModels.Image{}, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(fakeImage.Id).Return(http.StatusNoContent, nil),
			)

			err := factory.BuildAndPushImage(fakeImage.Id)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given that catalog's update image return error", func() {
			mocks.MockTapCatalogApiConnector.EXPECT().
				UpdateImage(fakeImage.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
				Return(catalogModels.Image{}, http.StatusInternalServerError, errors.New("unknown error"))

			err := factory.BuildAndPushImage(fakeImage.Id)

			Convey("Error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given that catalog's get image return error", func() {
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImage.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImage, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					GetImage(fakeImage.Id).
					Return(catalogModels.Image{}, http.StatusInternalServerError, errors.New("unknown error")),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImage.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateError, t)).
					Return(fakeImage, http.StatusOK, nil),
			)

			err := factory.BuildAndPushImage(fakeImage.Id)

			Convey("Error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestMarkImageAs(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	statesToTest := []catalogModels.ImageState{
		catalogModels.ImageStateReady,
		catalogModels.ImageStateBuilding,
		catalogModels.ImageStateError,
		catalogModels.ImageStatePending,
	}

	Convey("Test markImageAs", t, func() {
		imageID := "fake-id"

		for _, testedState := range statesToTest {
			Convey(fmt.Sprintf("For state=%v", testedState), func() {

				if testedState != catalogModels.ImageStatePending {
					mocks.MockTapCatalogApiConnector.EXPECT().UpdateImage(imageID, gomock.Any()).Return(catalogModels.Image{}, http.StatusOK, nil)
				}

				err := markImageAs(testedState, imageID)

				if testedState == catalogModels.ImageStatePending {
					Convey("error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					Convey("No error should be returned", func() {
						So(err, ShouldBeNil)
					})
				}
			})
		}
		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestGetImageInfoFromCatalog(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	Convey("Test getImageInfoFromCatalog", t, func() {
		imageID := "fake-id"

		Convey("Given that catalog client returns proper response", func() {
			mocks.MockTapCatalogApiConnector.EXPECT().GetImage(imageID).Return(catalogModels.Image{Id: imageID}, http.StatusOK, nil)

			image, err := getImageInfoFromCatalog(imageID)

			Convey("No error should be returned, image returned should have given image ID", func() {
				So(err, ShouldBeNil)
				So(image.Id, ShouldEqual, imageID)
			})

		})
		Convey("Given that catalog client returns error", func() {
			mocks.MockTapCatalogApiConnector.EXPECT().GetImage(imageID).Return(catalogModels.Image{}, http.StatusInternalServerError, errors.New("some error"))

			_, err := getImageInfoFromCatalog(imageID)

			Convey("Error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestApp_BuildImage(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	fakeImage := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeExec,
	}

	testCases := []SimpleTestCaseDefinition{
		{
			mockedCalls: func() GomockCalls {
				return GomockCalls{
					mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImage.Id,
						gomock.Any()).Return(nil),
					mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImage.Type,
						fakeImage.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImage.Id)).Return(nil),
				}
			},
			testDescription: "Given that blob-store and docker respond properly",
		},
		{
			mockedCalls: func() GomockCalls {
				return GomockCalls{
					mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImage.Id,
						gomock.Any()).Return(errors.New("get error")),
				}
			},
			testDescription:   "Given that blob-store respond badly",
			shouldReturnError: true,
		},
		{
			mockedCalls: func() GomockCalls {
				return GomockCalls{
					mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImage.Id, gomock.Any()).Return(nil),
					mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImage.Type,
						fakeImage.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImage.Id)).Return(errors.New("oops")),
				}
			},
			testDescription:   "Given that docker respond badly",
			shouldReturnError: true,
		},
	}

	Convey("Test buildImage", t, func() {
		for _, testCase := range testCases {
			Convey(testCase.testDescription, func() {
				gomock.InOrder(testCase.mockedCalls()...)

				imageTag, err := buildImage(fakeImage)

				if testCase.shouldReturnError {
					Convey("Error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					Convey("No error should be returned, image returned should have given image ID", func() {
						So(err, ShouldBeNil)
						So(imageTag, ShouldEqual, GetImageWithHubAddressWithoutProtocol(fakeImage.Id))
					})
				}
			})
		}
		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestPushImage(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	Convey("Test pushImage", t, func() {
		imageID := "image-id"
		imageTag := "image-tag"

		Convey("Given that docker connector return proper response", func() {
			mocks.MockDockerConnector.EXPECT().PushImage(imageTag).Return(nil)

			err := pushImage(imageID, imageTag)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})

		})

		Convey("Given that docker connector return error", func() {
			mocks.MockDockerConnector.EXPECT().PushImage(imageTag).Return(errors.New("fake error"))

			err := pushImage(imageID, imageTag)

			Convey("Error should be returned", func() {
				So(err, ShouldNotBeNil)
			})

		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func TestCleanupBlob(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	imageID := "image-id"

	testCases := []SimpleTestCaseDefinition{
		{
			mockedCalls: func() GomockCalls {
				return GomockCalls{
					mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(imageID).Return(http.StatusNoContent, nil),
				}
			},
			testDescription: "Expect blob-store client to be invoked (proper response)",
		},
		{
			mockedCalls: func() GomockCalls {
				return GomockCalls{
					mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(imageID).Return(http.StatusInternalServerError, errors.New("oops")),
				}
			},
			testDescription:   "Expect blob-store client to be invoked (bad response)",
			shouldReturnError: true,
		},
	}

	Convey("Test CleanupBlob", t, func() {
		for _, testCase := range testCases {
			Convey(testCase.testDescription, func() {
				gomock.InOrder(testCase.mockedCalls()...)

				cleanupBlob(imageID)
			})
		}

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}
