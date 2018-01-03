// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/cloudfoundry/bosh-cli/state/pkg (interfaces: Compiler,CompiledPackageRepo)

package mocks

import (
	pkg "github.com/cloudfoundry/bosh-cli/release/pkg"
	pkg0 "github.com/cloudfoundry/bosh-cli/state/pkg"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Compiler interface
type MockCompiler struct {
	ctrl     *gomock.Controller
	recorder *_MockCompilerRecorder
}

// Recorder for MockCompiler (not exported)
type _MockCompilerRecorder struct {
	mock *MockCompiler
}

func NewMockCompiler(ctrl *gomock.Controller) *MockCompiler {
	mock := &MockCompiler{ctrl: ctrl}
	mock.recorder = &_MockCompilerRecorder{mock}
	return mock
}

func (_m *MockCompiler) EXPECT() *_MockCompilerRecorder {
	return _m.recorder
}

func (_m *MockCompiler) Compile(_param0 pkg.Compilable) (pkg0.CompiledPackageRecord, bool, error) {
	ret := _m.ctrl.Call(_m, "Compile", _param0)
	ret0, _ := ret[0].(pkg0.CompiledPackageRecord)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCompilerRecorder) Compile(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Compile", arg0)
}

// Mock of CompiledPackageRepo interface
type MockCompiledPackageRepo struct {
	ctrl     *gomock.Controller
	recorder *_MockCompiledPackageRepoRecorder
}

// Recorder for MockCompiledPackageRepo (not exported)
type _MockCompiledPackageRepoRecorder struct {
	mock *MockCompiledPackageRepo
}

func NewMockCompiledPackageRepo(ctrl *gomock.Controller) *MockCompiledPackageRepo {
	mock := &MockCompiledPackageRepo{ctrl: ctrl}
	mock.recorder = &_MockCompiledPackageRepoRecorder{mock}
	return mock
}

func (_m *MockCompiledPackageRepo) EXPECT() *_MockCompiledPackageRepoRecorder {
	return _m.recorder
}

func (_m *MockCompiledPackageRepo) Find(_param0 pkg.Compilable) (pkg0.CompiledPackageRecord, bool, error) {
	ret := _m.ctrl.Call(_m, "Find", _param0)
	ret0, _ := ret[0].(pkg0.CompiledPackageRecord)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCompiledPackageRepoRecorder) Find(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Find", arg0)
}

func (_m *MockCompiledPackageRepo) Save(_param0 pkg.Compilable, _param1 pkg0.CompiledPackageRecord) error {
	ret := _m.ctrl.Call(_m, "Save", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCompiledPackageRepoRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Save", arg0, arg1)
}
