// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	filestorage "file-uploader/src/services/filestorage"

	mock "github.com/stretchr/testify/mock"

	models "file-uploader/src/models"

	multipart "mime/multipart"
)

// StorageService is an autogenerated mock type for the StorageService type
type StorageService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: data
func (_m *StorageService) Delete(data *models.File) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.File) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Provider provides a mock function with given fields:
func (_m *StorageService) Provider() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Read provides a mock function with given fields: data
func (_m *StorageService) Read(data *models.File) (*filestorage.FileReader, error) {
	ret := _m.Called(data)

	var r0 *filestorage.FileReader
	if rf, ok := ret.Get(0).(func(*models.File) *filestorage.FileReader); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*filestorage.FileReader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.File) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upload provides a mock function with given fields: file, data
func (_m *StorageService) Upload(file *multipart.FileHeader, data *models.File) error {
	ret := _m.Called(file, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(*multipart.FileHeader, *models.File) error); ok {
		r0 = rf(file, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
