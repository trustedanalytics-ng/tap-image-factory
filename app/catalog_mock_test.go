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
// Source: vendor/github.com/trustedanalytics-ng/tap-catalog/client/client.go

package app

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/trustedanalytics-ng/tap-catalog/models"
)

// Mock of TapCatalogApi interface
type MockTapCatalogApi struct {
	ctrl     *gomock.Controller
	recorder *_MockTapCatalogApiRecorder
}

// Recorder for MockTapCatalogApi (not exported)
type _MockTapCatalogApiRecorder struct {
	mock *MockTapCatalogApi
}

func NewMockTapCatalogApi(ctrl *gomock.Controller) *MockTapCatalogApi {
	mock := &MockTapCatalogApi{ctrl: ctrl}
	mock.recorder = &_MockTapCatalogApiRecorder{mock}
	return mock
}

func (_m *MockTapCatalogApi) EXPECT() *_MockTapCatalogApiRecorder {
	return _m.recorder
}

func (_m *MockTapCatalogApi) AddApplication(application models.Application) (models.Application, int, error) {
	ret := _m.ctrl.Call(_m, "AddApplication", application)
	ret0, _ := ret[0].(models.Application)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddApplication(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddApplication", arg0)
}

func (_m *MockTapCatalogApi) AddImage(image models.Image) (models.Image, int, error) {
	ret := _m.ctrl.Call(_m, "AddImage", image)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddImage", arg0)
}

func (_m *MockTapCatalogApi) AddService(service models.Service) (models.Service, int, error) {
	ret := _m.ctrl.Call(_m, "AddService", service)
	ret0, _ := ret[0].(models.Service)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddService", arg0)
}

func (_m *MockTapCatalogApi) AddServiceInstance(serviceId string, instance models.Instance) (models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "AddServiceInstance", serviceId, instance)
	ret0, _ := ret[0].(models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddServiceInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddServiceInstance", arg0, arg1)
}

func (_m *MockTapCatalogApi) AddServiceBrokerInstance(serviceId string, instance models.Instance) (models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "AddServiceBrokerInstance", serviceId, instance)
	ret0, _ := ret[0].(models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddServiceBrokerInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddServiceBrokerInstance", arg0, arg1)
}

func (_m *MockTapCatalogApi) AddApplicationInstance(applicationId string, instance models.Instance) (models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "AddApplicationInstance", applicationId, instance)
	ret0, _ := ret[0].(models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddApplicationInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddApplicationInstance", arg0, arg1)
}

func (_m *MockTapCatalogApi) AddTemplate(template models.Template) (models.Template, int, error) {
	ret := _m.ctrl.Call(_m, "AddTemplate", template)
	ret0, _ := ret[0].(models.Template)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) AddTemplate(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddTemplate", arg0)
}

func (_m *MockTapCatalogApi) GetApplication(applicationId string) (models.Application, int, error) {
	ret := _m.ctrl.Call(_m, "GetApplication", applicationId)
	ret0, _ := ret[0].(models.Application)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetApplication(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetApplication", arg0)
}

func (_m *MockTapCatalogApi) GetCatalogHealth() (int, error) {
	ret := _m.ctrl.Call(_m, "GetCatalogHealth")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) GetCatalogHealth() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCatalogHealth")
}

func (_m *MockTapCatalogApi) GetImage(imageId string) (models.Image, int, error) {
	ret := _m.ctrl.Call(_m, "GetImage", imageId)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetImage", arg0)
}

func (_m *MockTapCatalogApi) GetImageRefs(imageId string) (models.ImageRefsResponse, int, error) {
	ret := _m.ctrl.Call(_m, "GetImageRefs", imageId)
	ret0, _ := ret[0].(models.ImageRefsResponse)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetImageRefs(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetImageRefs", arg0)
}

func (_m *MockTapCatalogApi) GetInstance(instanceId string) (models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "GetInstance", instanceId)
	ret0, _ := ret[0].(models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetInstance(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetInstance", arg0)
}

func (_m *MockTapCatalogApi) GetInstanceBindings(instanceId string) ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "GetInstanceBindings", instanceId)
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetInstanceBindings(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetInstanceBindings", arg0)
}

func (_m *MockTapCatalogApi) GetServicePlan(serviceId string, planId string) (models.ServicePlan, int, error) {
	ret := _m.ctrl.Call(_m, "GetServicePlan", serviceId, planId)
	ret0, _ := ret[0].(models.ServicePlan)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetServicePlan(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServicePlan", arg0, arg1)
}

func (_m *MockTapCatalogApi) GetService(serviceId string) (models.Service, int, error) {
	ret := _m.ctrl.Call(_m, "GetService", serviceId)
	ret0, _ := ret[0].(models.Service)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetService", arg0)
}

func (_m *MockTapCatalogApi) GetServices() ([]models.Service, int, error) {
	ret := _m.ctrl.Call(_m, "GetServices")
	ret0, _ := ret[0].([]models.Service)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetServices() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServices")
}

func (_m *MockTapCatalogApi) GetLatestIndex() (models.Index, int, error) {
	ret := _m.ctrl.Call(_m, "GetLatestIndex")
	ret0, _ := ret[0].(models.Index)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) GetLatestIndex() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLatestIndex")
}

func (_m *MockTapCatalogApi) ListApplications() ([]models.Application, int, error) {
	ret := _m.ctrl.Call(_m, "ListApplications")
	ret0, _ := ret[0].([]models.Application)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListApplications() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListApplications")
}

func (_m *MockTapCatalogApi) ListApplicationsInstances() ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "ListApplicationsInstances")
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListApplicationsInstances() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListApplicationsInstances")
}

func (_m *MockTapCatalogApi) ListInstances() ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "ListInstances")
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListInstances() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListInstances")
}

func (_m *MockTapCatalogApi) ListImages() ([]models.Image, int, error) {
	ret := _m.ctrl.Call(_m, "ListImages")
	ret0, _ := ret[0].([]models.Image)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListImages() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListImages")
}

func (_m *MockTapCatalogApi) ListServicesInstances() ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "ListServicesInstances")
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListServicesInstances() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListServicesInstances")
}

func (_m *MockTapCatalogApi) ListApplicationInstances(applicationId string) ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "ListApplicationInstances", applicationId)
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListApplicationInstances(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListApplicationInstances", arg0)
}

func (_m *MockTapCatalogApi) ListServiceInstances(serviceId string) ([]models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "ListServiceInstances", serviceId)
	ret0, _ := ret[0].([]models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) ListServiceInstances(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListServiceInstances", arg0)
}

func (_m *MockTapCatalogApi) UpdateApplication(applicationId string, patches []models.Patch) (models.Application, int, error) {
	ret := _m.ctrl.Call(_m, "UpdateApplication", applicationId, patches)
	ret0, _ := ret[0].(models.Application)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdateApplication(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateApplication", arg0, arg1)
}

func (_m *MockTapCatalogApi) UpdateImage(imageId string, patches []models.Patch) (models.Image, int, error) {
	ret := _m.ctrl.Call(_m, "UpdateImage", imageId, patches)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdateImage(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateImage", arg0, arg1)
}

func (_m *MockTapCatalogApi) UpdateInstance(instanceId string, patches []models.Patch) (models.Instance, int, error) {
	ret := _m.ctrl.Call(_m, "UpdateInstance", instanceId, patches)
	ret0, _ := ret[0].(models.Instance)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdateInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateInstance", arg0, arg1)
}

func (_m *MockTapCatalogApi) UpdatePlan(serviceId string, planId string, patches []models.Patch) (models.ServicePlan, int, error) {
	ret := _m.ctrl.Call(_m, "UpdatePlan", serviceId, planId, patches)
	ret0, _ := ret[0].(models.ServicePlan)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdatePlan(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdatePlan", arg0, arg1, arg2)
}

func (_m *MockTapCatalogApi) UpdateService(serviceId string, patches []models.Patch) (models.Service, int, error) {
	ret := _m.ctrl.Call(_m, "UpdateService", serviceId, patches)
	ret0, _ := ret[0].(models.Service)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdateService(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateService", arg0, arg1)
}

func (_m *MockTapCatalogApi) UpdateTemplate(templateId string, patches []models.Patch) (models.Template, int, error) {
	ret := _m.ctrl.Call(_m, "UpdateTemplate", templateId, patches)
	ret0, _ := ret[0].(models.Template)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) UpdateTemplate(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateTemplate", arg0, arg1)
}

func (_m *MockTapCatalogApi) DeleteApplication(applicationId string) (int, error) {
	ret := _m.ctrl.Call(_m, "DeleteApplication", applicationId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) DeleteApplication(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteApplication", arg0)
}

func (_m *MockTapCatalogApi) DeleteService(serviceId string) (int, error) {
	ret := _m.ctrl.Call(_m, "DeleteService", serviceId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) DeleteService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteService", arg0)
}

func (_m *MockTapCatalogApi) DeleteServicePlan(serviceId string, planId string) (int, error) {
	ret := _m.ctrl.Call(_m, "DeleteServicePlan", serviceId, planId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) DeleteServicePlan(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteServicePlan", arg0, arg1)
}

func (_m *MockTapCatalogApi) DeleteImage(imageId string) (int, error) {
	ret := _m.ctrl.Call(_m, "DeleteImage", imageId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) DeleteImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteImage", arg0)
}

func (_m *MockTapCatalogApi) DeleteInstance(instanceId string) (int, error) {
	ret := _m.ctrl.Call(_m, "DeleteInstance", instanceId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockTapCatalogApiRecorder) DeleteInstance(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteInstance", arg0)
}

func (_m *MockTapCatalogApi) WatchInstances(afterIndex uint64) (models.StateChange, int, error) {
	ret := _m.ctrl.Call(_m, "WatchInstances", afterIndex)
	ret0, _ := ret[0].(models.StateChange)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) WatchInstances(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WatchInstances", arg0)
}

func (_m *MockTapCatalogApi) WatchInstance(instanceId string, afterIndex uint64) (models.StateChange, int, error) {
	ret := _m.ctrl.Call(_m, "WatchInstance", instanceId, afterIndex)
	ret0, _ := ret[0].(models.StateChange)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) WatchInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WatchInstance", arg0, arg1)
}

func (_m *MockTapCatalogApi) WatchImages(afterIndex uint64) (models.StateChange, int, error) {
	ret := _m.ctrl.Call(_m, "WatchImages", afterIndex)
	ret0, _ := ret[0].(models.StateChange)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) WatchImages(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WatchImages", arg0)
}

func (_m *MockTapCatalogApi) WatchImage(imageId string, afterIndex uint64) (models.StateChange, int, error) {
	ret := _m.ctrl.Call(_m, "WatchImage", imageId, afterIndex)
	ret0, _ := ret[0].(models.StateChange)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) WatchImage(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WatchImage", arg0, arg1)
}

func (_m *MockTapCatalogApi) CheckStateStability() (models.StateStability, int, error) {
	ret := _m.ctrl.Call(_m, "CheckStateStability")
	ret0, _ := ret[0].(models.StateStability)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockTapCatalogApiRecorder) CheckStateStability() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CheckStateStability")
}
