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
// Source: app/archive_reader.go

package app

import (
	tar "archive/tar"
	gzip "compress/gzip"
	gomock "github.com/golang/mock/gomock"
	io "io"
)

// Mock of ArchiveReader interface
type MockArchiveReader struct {
	ctrl     *gomock.Controller
	recorder *_MockArchiveReaderRecorder
}

// Recorder for MockArchiveReader (not exported)
type _MockArchiveReaderRecorder struct {
	mock *MockArchiveReader
}

func NewMockArchiveReader(ctrl *gomock.Controller) *MockArchiveReader {
	mock := &MockArchiveReader{ctrl: ctrl}
	mock.recorder = &_MockArchiveReaderRecorder{mock}
	return mock
}

func (_m *MockArchiveReader) EXPECT() *_MockArchiveReaderRecorder {
	return _m.recorder
}

func (_m *MockArchiveReader) NewGzipReader(reader io.Reader) (*gzip.Reader, error) {
	ret := _m.ctrl.Call(_m, "NewGzipReader", reader)
	ret0, _ := ret[0].(*gzip.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockArchiveReaderRecorder) NewGzipReader(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewGzipReader", arg0)
}

func (_m *MockArchiveReader) NewTarReader(reader io.Reader) *tar.Reader {
	ret := _m.ctrl.Call(_m, "NewTarReader", reader)
	ret0, _ := ret[0].(*tar.Reader)
	return ret0
}

func (_mr *_MockArchiveReaderRecorder) NewTarReader(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewTarReader", arg0)
}
