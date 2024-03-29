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
// Automatically generated by MockGen. DO NOT EDIT!
// Source: app/image_builder.go

package app

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/trustedanalytics-ng/tap-catalog/models"
	io "io"
	os "os"
)

// Mock of ImageBuilder interface
type MockImageBuilder struct {
	ctrl     *gomock.Controller
	recorder *_MockImageBuilderRecorder
}

// Recorder for MockImageBuilder (not exported)
type _MockImageBuilderRecorder struct {
	mock *MockImageBuilder
}

func NewMockImageBuilder(ctrl *gomock.Controller) *MockImageBuilder {
	mock := &MockImageBuilder{ctrl: ctrl}
	mock.recorder = &_MockImageBuilderRecorder{mock}
	return mock
}

func (_m *MockImageBuilder) EXPECT() *_MockImageBuilderRecorder {
	return _m.recorder
}

func (_m *MockImageBuilder) CreateImage(artifact *os.File, imageType models.ImageType, blobType models.BlobType, tag string) error {
	ret := _m.ctrl.Call(_m, "CreateImage", artifact, imageType, blobType, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockImageBuilderRecorder) CreateImage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateImage", arg0, arg1, arg2, arg3)
}

func (_m *MockImageBuilder) buildImage(buildContext io.Reader, imageId string) error {
	ret := _m.ctrl.Call(_m, "buildImage", buildContext, imageId)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockImageBuilderRecorder) buildImage(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "buildImage", arg0, arg1)
}

func (_m *MockImageBuilder) TagImage(imageId string, tag string) error {
	ret := _m.ctrl.Call(_m, "TagImage", imageId, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockImageBuilderRecorder) TagImage(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TagImage", arg0, arg1)
}

func (_m *MockImageBuilder) PushImage(tag string) error {
	ret := _m.ctrl.Call(_m, "PushImage", tag)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockImageBuilderRecorder) PushImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PushImage", arg0)
}
