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
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/trustedanalytics-ng/tap-catalog/builder"
	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
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

	factory := NewFactoryWithCustomProbeTimes(mocks.MockImageChecker, defaultImageReadinessProbeRetries, defaultImageReadinessProbeDelaySeconds*time.Second)

	fakeImageTarGz := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeTarGz,
	}

	fakeImageJar := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeJava,
		BlobType: catalogModels.BlobTypeJar,
	}

	fakeImageExec := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeExec,
	}

	Convey("Test BuildAndPushImage", t, func() {
		Convey("Given TarGz image and that all underneath functions invokes successfully", func() {
			buffer := getTarGzBufferWithFile(t, runFileName)
			gzipReader, tarReader := getReaders(buffer, t)
			defer gzipReader.Close()
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageTarGz.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImageTarGz, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().GetImage(fakeImageTarGz.Id).Return(fakeImageTarGz, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageTarGz.Id, gomock.Any()).Return(nil),
				mocks.MockReader.EXPECT().NewGzipReader(gomock.Any()).Return(gzipReader, nil),
				mocks.MockReader.EXPECT().NewTarReader(gomock.Any()).Return(tarReader),
				mocks.MockDockerConnector.EXPECT().
					CreateImage(gomock.Any(), fakeImageTarGz.Type, fakeImageTarGz.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageTarGz.Id)).
					Return(nil),
				mocks.MockDockerConnector.EXPECT().PushImage(GetImageWithHubAddressWithoutProtocol(fakeImageTarGz.Id)).Return(nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(fakeImageTarGz.Id, defaultImageTag).Return(true, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageTarGz.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateReady, t)).
					Return(catalogModels.Image{}, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(fakeImageTarGz.Id).Return(http.StatusNoContent, nil),
			)

			err := factory.BuildAndPushImage(fakeImageTarGz.Id)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given Jar image and that all underneath functions invokes successfully", func() {
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageJar.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImageJar, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().GetImage(fakeImageJar.Id).Return(fakeImageJar, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageJar.Id, gomock.Any()).Return(nil),
				mocks.MockDockerConnector.EXPECT().
					CreateImage(gomock.Any(), fakeImageJar.Type, fakeImageJar.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageJar.Id)).
					Return(nil),
				mocks.MockDockerConnector.EXPECT().PushImage(GetImageWithHubAddressWithoutProtocol(fakeImageJar.Id)).Return(nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(fakeImageJar.Id, defaultImageTag).Return(true, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageJar.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateReady, t)).
					Return(catalogModels.Image{}, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(fakeImageJar.Id).Return(http.StatusNoContent, nil),
			)

			err := factory.BuildAndPushImage(fakeImageJar.Id)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given Exec image and that all underneath functions invokes successfully", func() {
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageExec.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImageExec, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().GetImage(fakeImageExec.Id).Return(fakeImageExec, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageExec.Id, gomock.Any()).Return(nil),
				mocks.MockDockerConnector.EXPECT().
					CreateImage(gomock.Any(), fakeImageExec.Type, fakeImageExec.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageExec.Id)).
					Return(nil),
				mocks.MockDockerConnector.EXPECT().PushImage(GetImageWithHubAddressWithoutProtocol(fakeImageExec.Id)).Return(nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(fakeImageExec.Id, defaultImageTag).Return(true, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageExec.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateReady, t)).
					Return(catalogModels.Image{}, http.StatusOK, nil),
				mocks.MockBlobStoreConnector.EXPECT().DeleteBlob(fakeImageExec.Id).Return(http.StatusNoContent, nil),
			)

			err := factory.BuildAndPushImage(fakeImageExec.Id)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given that catalog's update image return error", func() {
			mocks.MockTapCatalogApiConnector.EXPECT().
				UpdateImage(fakeImageExec.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
				Return(catalogModels.Image{}, http.StatusInternalServerError, errors.New("unknown error"))

			err := factory.BuildAndPushImage(fakeImageExec.Id)

			Convey("Error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given that catalog's get image return error", func() {
			gomock.InOrder(
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageTarGz.Id, makePatchUserFriendly(catalogModels.ImageStatePending, catalogModels.ImageStateBuilding, t)).
					Return(fakeImageTarGz, http.StatusOK, nil),
				mocks.MockTapCatalogApiConnector.EXPECT().
					GetImage(fakeImageTarGz.Id).
					Return(catalogModels.Image{}, http.StatusInternalServerError, errors.New("unknown error")),
				mocks.MockTapCatalogApiConnector.EXPECT().
					UpdateImage(fakeImageTarGz.Id, makePatchUserFriendly(catalogModels.ImageStateBuilding, catalogModels.ImageStateError, t)).
					Return(fakeImageTarGz, http.StatusOK, nil),
			)

			err := factory.BuildAndPushImage(fakeImageTarGz.Id)

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

	fakeImageTarGz := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeTarGz,
	}

	fakeImageExec := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeGo,
		BlobType: catalogModels.BlobTypeExec,
	}

	fakeImageJar := catalogModels.Image{
		Id:       "fake-id",
		Type:     catalogModels.ImageTypeJava,
		BlobType: catalogModels.BlobTypeJar,
	}

	testCases := []struct {
		SimpleTestCaseDefinition
		image catalogModels.Image
	}{
		{
			image: fakeImageTarGz,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					buffer := getTarGzBufferWithFile(t, runFileName)
					gzipReader, tarReader := getReaders(buffer, t)
					defer gzipReader.Close()
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageTarGz.Id,
							gomock.Any()).Return(nil),
						mocks.MockReader.EXPECT().NewGzipReader(gomock.Any()).Return(gzipReader, nil),
						mocks.MockReader.EXPECT().NewTarReader(gomock.Any()).Return(tarReader),
						mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImageTarGz.Type,
							fakeImageTarGz.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageTarGz.Id)).Return(nil),
					}
				},
				testDescription: "Given that blob-store and docker respond properly",
			},
		},
		{
			image: fakeImageTarGz,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageTarGz.Id,
							gomock.Any()).Return(errors.New("get error")),
					}
				},
				testDescription:   "Given that blob-store API responds badly",
				shouldReturnError: true,
			},
		},
		{
			image: fakeImageTarGz,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					buffer := getTarGzBufferWithFile(t, runFileName)
					gzipReader, tarReader := getReaders(buffer, t)
					defer gzipReader.Close()
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageTarGz.Id, gomock.Any()).Return(nil),
						mocks.MockReader.EXPECT().NewGzipReader(gomock.Any()).Return(gzipReader, nil),
						mocks.MockReader.EXPECT().NewTarReader(gomock.Any()).Return(tarReader),
						mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImageTarGz.Type,
							fakeImageTarGz.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageTarGz.Id)).Return(errors.New("oops")),
					}
				},
				testDescription:   "Given that docker respond badly",
				shouldReturnError: true,
			},
		},
		{
			image: fakeImageTarGz,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					buffer := getTarGzBufferWithFile(t, "some_file.txt")
					gzipReader, tarReader := getReaders(buffer, t)
					defer gzipReader.Close()
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageTarGz.Id, gomock.Any()).Return(nil),
						mocks.MockReader.EXPECT().NewGzipReader(gomock.Any()).Return(gzipReader, nil),
						mocks.MockReader.EXPECT().NewTarReader(gomock.Any()).Return(tarReader),
					}
				},
				testDescription:   fmt.Sprintf("Given that tar.gz file does not contain %s file", runFileName),
				shouldReturnError: true,
			},
		},
		{
			image: fakeImageJar,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageJar.Id,
							gomock.Any()).Return(nil),
						mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImageJar.Type,
							fakeImageJar.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageJar.Id)).Return(nil),
					}
				},
				testDescription: "Given Jar image and blob-store and docker respond properly",
			},
		},
		{
			image: fakeImageExec,
			SimpleTestCaseDefinition: SimpleTestCaseDefinition{
				mockedCalls: func() GomockCalls {
					return GomockCalls{
						mocks.MockBlobStoreConnector.EXPECT().GetBlob(fakeImageExec.Id,
							gomock.Any()).Return(nil),
						mocks.MockDockerConnector.EXPECT().CreateImage(gomock.Any(), fakeImageExec.Type,
							fakeImageExec.BlobType, GetImageWithHubAddressWithoutProtocol(fakeImageExec.Id)).Return(nil),
					}
				},
				testDescription: "Given Exec image and blob-store and docker respond properly",
			},
		},
	}

	Convey("Test buildImage", t, func() {
		for _, testCase := range testCases {
			Convey(testCase.testDescription, func() {
				gomock.InOrder(testCase.mockedCalls()...)

				imageTag, err := buildImage(testCase.image)

				if testCase.shouldReturnError {
					Convey("Error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				} else {
					Convey("No error should be returned, image returned should have given image ID", func() {
						So(err, ShouldBeNil)
						So(imageTag, ShouldEqual, GetImageWithHubAddressWithoutProtocol(testCase.image.Id))
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

func TestEnsureImageReady(t *testing.T) {
	mockCtrl, mocks, _, _ := prepareMocksAndClient(t)

	Convey("Test ensureImageReady", t, func() {
		imageID := "fake-id"
		maxRetries := uint32(3)

		factory := NewFactoryWithCustomProbeTimes(mocks.MockImageChecker, maxRetries, time.Nanosecond).(*Factory)

		Convey("Given that imageProber returns immediately readiness confirmation ", func() {
			mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(true, nil)

			err := factory.ensureImageReady(imageID, defaultImageTag)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given that imageProber returns readiness confirmation after some time", func() {
			gomock.InOrder(
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(true, nil),
			)

			err := factory.ensureImageReady(imageID, defaultImageTag)

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("Given that imageProber exceeds retries", func() {
			gomock.InOrder(
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, nil),
				mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, nil),
			)

			err := factory.ensureImageReady(imageID, defaultImageTag)

			Convey("error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Given that imageProber returns error", func() {
			mocks.MockImageChecker.EXPECT().IsImageReady(imageID, defaultImageTag).Return(false, errors.New("fake error"))

			err := factory.ensureImageReady(imageID, defaultImageTag)

			Convey("error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Reset(func() {
			mockCtrl.Finish()
		})
	})
}

func getTarGzBufferWithFile(t *testing.T, path string) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	bufferWriter := bufio.NewWriter(buffer)
	defer bufferWriter.Flush()

	gzipWriter := gzip.NewWriter(bufferWriter)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	if err := addFile(tarWriter, path); err != nil {
		t.Fatalf("cannot add file %q: %v", path, err)
	}

	return buffer
}

func addFile(writer *tar.Writer, path string) error {
	header := &tar.Header{}
	header.Name = path
	header.Size = 1
	header.Typeflag = tar.TypeReg

	if err := writer.WriteHeader(header); err != nil {
		return err
	}

	if _, err := writer.Write([]byte{'a'}); err != nil {
		return err
	}

	return nil
}

func getReaders(reader io.Reader, t *testing.T) (*gzip.Reader, *tar.Reader) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		t.Fatalf("gzip.newReader failed: %v", err)
	}
	tarReader := tar.NewReader(gzipReader)
	return gzipReader, tarReader
}
