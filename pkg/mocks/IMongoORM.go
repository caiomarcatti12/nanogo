// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	rsql "github.com/caiomarcatti12/nanogo/pkg/rsql"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IMongoORM is an autogenerated mock type for the IMongoORM type
type IMongoORM[T interface{}] struct {
	mock.Mock
}

type IMongoORM_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *IMongoORM[T]) EXPECT() *IMongoORM_Expecter[T] {
	return &IMongoORM_Expecter[T]{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: document
func (_m *IMongoORM[T]) Delete(document T) (bool, error) {
	ret := _m.Called(document)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(T) (bool, error)); ok {
		return rf(document)
	}
	if rf, ok := ret.Get(0).(func(T) bool); ok {
		r0 = rf(document)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(T) error); ok {
		r1 = rf(document)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type IMongoORM_Delete_Call[T interface{}] struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - document T
func (_e *IMongoORM_Expecter[T]) Delete(document interface{}) *IMongoORM_Delete_Call[T] {
	return &IMongoORM_Delete_Call[T]{Call: _e.mock.On("Delete", document)}
}

func (_c *IMongoORM_Delete_Call[T]) Run(run func(document T)) *IMongoORM_Delete_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(T))
	})
	return _c
}

func (_c *IMongoORM_Delete_Call[T]) Return(_a0 bool, _a1 error) *IMongoORM_Delete_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_Delete_Call[T]) RunAndReturn(run func(T) (bool, error)) *IMongoORM_Delete_Call[T] {
	_c.Call.Return(run)
	return _c
}

// DeleteById provides a mock function with given fields: id
func (_m *IMongoORM[T]) DeleteById(id uuid.UUID) (bool, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteById")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (bool, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_DeleteById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteById'
type IMongoORM_DeleteById_Call[T interface{}] struct {
	*mock.Call
}

// DeleteById is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *IMongoORM_Expecter[T]) DeleteById(id interface{}) *IMongoORM_DeleteById_Call[T] {
	return &IMongoORM_DeleteById_Call[T]{Call: _e.mock.On("DeleteById", id)}
}

func (_c *IMongoORM_DeleteById_Call[T]) Run(run func(id uuid.UUID)) *IMongoORM_DeleteById_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *IMongoORM_DeleteById_Call[T]) Return(_a0 bool, _a1 error) *IMongoORM_DeleteById_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_DeleteById_Call[T]) RunAndReturn(run func(uuid.UUID) (bool, error)) *IMongoORM_DeleteById_Call[T] {
	_c.Call.Return(run)
	return _c
}

// FindAll provides a mock function with given fields:
func (_m *IMongoORM[T]) FindAll() ([]T, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
	}

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]T, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []T); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_FindAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAll'
type IMongoORM_FindAll_Call[T interface{}] struct {
	*mock.Call
}

// FindAll is a helper method to define mock.On call
func (_e *IMongoORM_Expecter[T]) FindAll() *IMongoORM_FindAll_Call[T] {
	return &IMongoORM_FindAll_Call[T]{Call: _e.mock.On("FindAll")}
}

func (_c *IMongoORM_FindAll_Call[T]) Run(run func()) *IMongoORM_FindAll_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IMongoORM_FindAll_Call[T]) Return(_a0 []T, _a1 error) *IMongoORM_FindAll_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_FindAll_Call[T]) RunAndReturn(run func() ([]T, error)) *IMongoORM_FindAll_Call[T] {
	_c.Call.Return(run)
	return _c
}

// FindById provides a mock function with given fields: id
func (_m *IMongoORM[T]) FindById(id uuid.UUID) (*T, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FindById")
	}

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*T, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *T); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_FindById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindById'
type IMongoORM_FindById_Call[T interface{}] struct {
	*mock.Call
}

// FindById is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *IMongoORM_Expecter[T]) FindById(id interface{}) *IMongoORM_FindById_Call[T] {
	return &IMongoORM_FindById_Call[T]{Call: _e.mock.On("FindById", id)}
}

func (_c *IMongoORM_FindById_Call[T]) Run(run func(id uuid.UUID)) *IMongoORM_FindById_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *IMongoORM_FindById_Call[T]) Return(_a0 *T, _a1 error) *IMongoORM_FindById_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_FindById_Call[T]) RunAndReturn(run func(uuid.UUID) (*T, error)) *IMongoORM_FindById_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: document
func (_m *IMongoORM[T]) Insert(document T) (uuid.UUID, error) {
	ret := _m.Called(document)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(T) (uuid.UUID, error)); ok {
		return rf(document)
	}
	if rf, ok := ret.Get(0).(func(T) uuid.UUID); ok {
		r0 = rf(document)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(T) error); ok {
		r1 = rf(document)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type IMongoORM_Insert_Call[T interface{}] struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - document T
func (_e *IMongoORM_Expecter[T]) Insert(document interface{}) *IMongoORM_Insert_Call[T] {
	return &IMongoORM_Insert_Call[T]{Call: _e.mock.On("Insert", document)}
}

func (_c *IMongoORM_Insert_Call[T]) Run(run func(document T)) *IMongoORM_Insert_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(T))
	})
	return _c
}

func (_c *IMongoORM_Insert_Call[T]) Return(_a0 uuid.UUID, _a1 error) *IMongoORM_Insert_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_Insert_Call[T]) RunAndReturn(run func(T) (uuid.UUID, error)) *IMongoORM_Insert_Call[T] {
	_c.Call.Return(run)
	return _c
}

// RawQueryParseRsql provides a mock function with given fields: filter
func (_m *IMongoORM[T]) RawQueryParseRsql(filter rsql.QueryFilter) ([]T, int64, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for RawQueryParseRsql")
	}

	var r0 []T
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(rsql.QueryFilter) ([]T, int64, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(rsql.QueryFilter) []T); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func(rsql.QueryFilter) int64); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(rsql.QueryFilter) error); ok {
		r2 = rf(filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// IMongoORM_RawQueryParseRsql_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RawQueryParseRsql'
type IMongoORM_RawQueryParseRsql_Call[T interface{}] struct {
	*mock.Call
}

// RawQueryParseRsql is a helper method to define mock.On call
//   - filter rsql.QueryFilter
func (_e *IMongoORM_Expecter[T]) RawQueryParseRsql(filter interface{}) *IMongoORM_RawQueryParseRsql_Call[T] {
	return &IMongoORM_RawQueryParseRsql_Call[T]{Call: _e.mock.On("RawQueryParseRsql", filter)}
}

func (_c *IMongoORM_RawQueryParseRsql_Call[T]) Run(run func(filter rsql.QueryFilter)) *IMongoORM_RawQueryParseRsql_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(rsql.QueryFilter))
	})
	return _c
}

func (_c *IMongoORM_RawQueryParseRsql_Call[T]) Return(_a0 []T, _a1 int64, _a2 error) *IMongoORM_RawQueryParseRsql_Call[T] {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *IMongoORM_RawQueryParseRsql_Call[T]) RunAndReturn(run func(rsql.QueryFilter) ([]T, int64, error)) *IMongoORM_RawQueryParseRsql_Call[T] {
	_c.Call.Return(run)
	return _c
}

// SetCollection provides a mock function with given fields: collectionName
func (_m *IMongoORM[T]) SetCollection(collectionName string) {
	_m.Called(collectionName)
}

// IMongoORM_SetCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCollection'
type IMongoORM_SetCollection_Call[T interface{}] struct {
	*mock.Call
}

// SetCollection is a helper method to define mock.On call
//   - collectionName string
func (_e *IMongoORM_Expecter[T]) SetCollection(collectionName interface{}) *IMongoORM_SetCollection_Call[T] {
	return &IMongoORM_SetCollection_Call[T]{Call: _e.mock.On("SetCollection", collectionName)}
}

func (_c *IMongoORM_SetCollection_Call[T]) Run(run func(collectionName string)) *IMongoORM_SetCollection_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IMongoORM_SetCollection_Call[T]) Return() *IMongoORM_SetCollection_Call[T] {
	_c.Call.Return()
	return _c
}

func (_c *IMongoORM_SetCollection_Call[T]) RunAndReturn(run func(string)) *IMongoORM_SetCollection_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: document
func (_m *IMongoORM[T]) Update(document T) (bool, error) {
	ret := _m.Called(document)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(T) (bool, error)); ok {
		return rf(document)
	}
	if rf, ok := ret.Get(0).(func(T) bool); ok {
		r0 = rf(document)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(T) error); ok {
		r1 = rf(document)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IMongoORM_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type IMongoORM_Update_Call[T interface{}] struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - document T
func (_e *IMongoORM_Expecter[T]) Update(document interface{}) *IMongoORM_Update_Call[T] {
	return &IMongoORM_Update_Call[T]{Call: _e.mock.On("Update", document)}
}

func (_c *IMongoORM_Update_Call[T]) Run(run func(document T)) *IMongoORM_Update_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(T))
	})
	return _c
}

func (_c *IMongoORM_Update_Call[T]) Return(_a0 bool, _a1 error) *IMongoORM_Update_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IMongoORM_Update_Call[T]) RunAndReturn(run func(T) (bool, error)) *IMongoORM_Update_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewIMongoORM creates a new instance of IMongoORM. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMongoORM[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *IMongoORM[T] {
	mock := &IMongoORM[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
