// Code generated by counterfeiter. DO NOT EDIT.
package releasefakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/release"
)

type FakeWriter struct {
	WriteStub        func(release.Release, []string) (string, error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		arg1 release.Release
		arg2 []string
	}
	writeReturns struct {
		result1 string
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWriter) Write(arg1 release.Release, arg2 []string) (string, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.writeMutex.Lock()
	ret, specificReturn := fake.writeReturnsOnCall[len(fake.writeArgsForCall)]
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		arg1 release.Release
		arg2 []string
	}{arg1, arg2Copy})
	fake.recordInvocation("Write", []interface{}{arg1, arg2Copy})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.writeReturns.result1, fake.writeReturns.result2
}

func (fake *FakeWriter) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *FakeWriter) WriteArgsForCall(i int) (release.Release, []string) {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return fake.writeArgsForCall[i].arg1, fake.writeArgsForCall[i].arg2
}

func (fake *FakeWriter) WriteReturns(result1 string, result2 error) {
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeWriter) WriteReturnsOnCall(i int, result1 string, result2 error) {
	fake.WriteStub = nil
	if fake.writeReturnsOnCall == nil {
		fake.writeReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.writeReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeWriter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeWriter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ release.Writer = new(FakeWriter)
