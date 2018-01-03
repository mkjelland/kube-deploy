// Code generated by counterfeiter. DO NOT EDIT.
package indexfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/releasedir/index"
)

type FakeReporter struct {
	IndexEntryStartedAddingStub        func(type_, desc string)
	indexEntryStartedAddingMutex       sync.RWMutex
	indexEntryStartedAddingArgsForCall []struct {
		type_ string
		desc  string
	}
	IndexEntryFinishedAddingStub        func(type_, desc string, err error)
	indexEntryFinishedAddingMutex       sync.RWMutex
	indexEntryFinishedAddingArgsForCall []struct {
		type_ string
		desc  string
		err   error
	}
	IndexEntryDownloadStartedStub        func(type_, desc string)
	indexEntryDownloadStartedMutex       sync.RWMutex
	indexEntryDownloadStartedArgsForCall []struct {
		type_ string
		desc  string
	}
	IndexEntryDownloadFinishedStub        func(type_, desc string, err error)
	indexEntryDownloadFinishedMutex       sync.RWMutex
	indexEntryDownloadFinishedArgsForCall []struct {
		type_ string
		desc  string
		err   error
	}
	IndexEntryUploadStartedStub        func(type_, desc string)
	indexEntryUploadStartedMutex       sync.RWMutex
	indexEntryUploadStartedArgsForCall []struct {
		type_ string
		desc  string
	}
	IndexEntryUploadFinishedStub        func(type_, desc string, err error)
	indexEntryUploadFinishedMutex       sync.RWMutex
	indexEntryUploadFinishedArgsForCall []struct {
		type_ string
		desc  string
		err   error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReporter) IndexEntryStartedAdding(type_ string, desc string) {
	fake.indexEntryStartedAddingMutex.Lock()
	fake.indexEntryStartedAddingArgsForCall = append(fake.indexEntryStartedAddingArgsForCall, struct {
		type_ string
		desc  string
	}{type_, desc})
	fake.recordInvocation("IndexEntryStartedAdding", []interface{}{type_, desc})
	fake.indexEntryStartedAddingMutex.Unlock()
	if fake.IndexEntryStartedAddingStub != nil {
		fake.IndexEntryStartedAddingStub(type_, desc)
	}
}

func (fake *FakeReporter) IndexEntryStartedAddingCallCount() int {
	fake.indexEntryStartedAddingMutex.RLock()
	defer fake.indexEntryStartedAddingMutex.RUnlock()
	return len(fake.indexEntryStartedAddingArgsForCall)
}

func (fake *FakeReporter) IndexEntryStartedAddingArgsForCall(i int) (string, string) {
	fake.indexEntryStartedAddingMutex.RLock()
	defer fake.indexEntryStartedAddingMutex.RUnlock()
	return fake.indexEntryStartedAddingArgsForCall[i].type_, fake.indexEntryStartedAddingArgsForCall[i].desc
}

func (fake *FakeReporter) IndexEntryFinishedAdding(type_ string, desc string, err error) {
	fake.indexEntryFinishedAddingMutex.Lock()
	fake.indexEntryFinishedAddingArgsForCall = append(fake.indexEntryFinishedAddingArgsForCall, struct {
		type_ string
		desc  string
		err   error
	}{type_, desc, err})
	fake.recordInvocation("IndexEntryFinishedAdding", []interface{}{type_, desc, err})
	fake.indexEntryFinishedAddingMutex.Unlock()
	if fake.IndexEntryFinishedAddingStub != nil {
		fake.IndexEntryFinishedAddingStub(type_, desc, err)
	}
}

func (fake *FakeReporter) IndexEntryFinishedAddingCallCount() int {
	fake.indexEntryFinishedAddingMutex.RLock()
	defer fake.indexEntryFinishedAddingMutex.RUnlock()
	return len(fake.indexEntryFinishedAddingArgsForCall)
}

func (fake *FakeReporter) IndexEntryFinishedAddingArgsForCall(i int) (string, string, error) {
	fake.indexEntryFinishedAddingMutex.RLock()
	defer fake.indexEntryFinishedAddingMutex.RUnlock()
	return fake.indexEntryFinishedAddingArgsForCall[i].type_, fake.indexEntryFinishedAddingArgsForCall[i].desc, fake.indexEntryFinishedAddingArgsForCall[i].err
}

func (fake *FakeReporter) IndexEntryDownloadStarted(type_ string, desc string) {
	fake.indexEntryDownloadStartedMutex.Lock()
	fake.indexEntryDownloadStartedArgsForCall = append(fake.indexEntryDownloadStartedArgsForCall, struct {
		type_ string
		desc  string
	}{type_, desc})
	fake.recordInvocation("IndexEntryDownloadStarted", []interface{}{type_, desc})
	fake.indexEntryDownloadStartedMutex.Unlock()
	if fake.IndexEntryDownloadStartedStub != nil {
		fake.IndexEntryDownloadStartedStub(type_, desc)
	}
}

func (fake *FakeReporter) IndexEntryDownloadStartedCallCount() int {
	fake.indexEntryDownloadStartedMutex.RLock()
	defer fake.indexEntryDownloadStartedMutex.RUnlock()
	return len(fake.indexEntryDownloadStartedArgsForCall)
}

func (fake *FakeReporter) IndexEntryDownloadStartedArgsForCall(i int) (string, string) {
	fake.indexEntryDownloadStartedMutex.RLock()
	defer fake.indexEntryDownloadStartedMutex.RUnlock()
	return fake.indexEntryDownloadStartedArgsForCall[i].type_, fake.indexEntryDownloadStartedArgsForCall[i].desc
}

func (fake *FakeReporter) IndexEntryDownloadFinished(type_ string, desc string, err error) {
	fake.indexEntryDownloadFinishedMutex.Lock()
	fake.indexEntryDownloadFinishedArgsForCall = append(fake.indexEntryDownloadFinishedArgsForCall, struct {
		type_ string
		desc  string
		err   error
	}{type_, desc, err})
	fake.recordInvocation("IndexEntryDownloadFinished", []interface{}{type_, desc, err})
	fake.indexEntryDownloadFinishedMutex.Unlock()
	if fake.IndexEntryDownloadFinishedStub != nil {
		fake.IndexEntryDownloadFinishedStub(type_, desc, err)
	}
}

func (fake *FakeReporter) IndexEntryDownloadFinishedCallCount() int {
	fake.indexEntryDownloadFinishedMutex.RLock()
	defer fake.indexEntryDownloadFinishedMutex.RUnlock()
	return len(fake.indexEntryDownloadFinishedArgsForCall)
}

func (fake *FakeReporter) IndexEntryDownloadFinishedArgsForCall(i int) (string, string, error) {
	fake.indexEntryDownloadFinishedMutex.RLock()
	defer fake.indexEntryDownloadFinishedMutex.RUnlock()
	return fake.indexEntryDownloadFinishedArgsForCall[i].type_, fake.indexEntryDownloadFinishedArgsForCall[i].desc, fake.indexEntryDownloadFinishedArgsForCall[i].err
}

func (fake *FakeReporter) IndexEntryUploadStarted(type_ string, desc string) {
	fake.indexEntryUploadStartedMutex.Lock()
	fake.indexEntryUploadStartedArgsForCall = append(fake.indexEntryUploadStartedArgsForCall, struct {
		type_ string
		desc  string
	}{type_, desc})
	fake.recordInvocation("IndexEntryUploadStarted", []interface{}{type_, desc})
	fake.indexEntryUploadStartedMutex.Unlock()
	if fake.IndexEntryUploadStartedStub != nil {
		fake.IndexEntryUploadStartedStub(type_, desc)
	}
}

func (fake *FakeReporter) IndexEntryUploadStartedCallCount() int {
	fake.indexEntryUploadStartedMutex.RLock()
	defer fake.indexEntryUploadStartedMutex.RUnlock()
	return len(fake.indexEntryUploadStartedArgsForCall)
}

func (fake *FakeReporter) IndexEntryUploadStartedArgsForCall(i int) (string, string) {
	fake.indexEntryUploadStartedMutex.RLock()
	defer fake.indexEntryUploadStartedMutex.RUnlock()
	return fake.indexEntryUploadStartedArgsForCall[i].type_, fake.indexEntryUploadStartedArgsForCall[i].desc
}

func (fake *FakeReporter) IndexEntryUploadFinished(type_ string, desc string, err error) {
	fake.indexEntryUploadFinishedMutex.Lock()
	fake.indexEntryUploadFinishedArgsForCall = append(fake.indexEntryUploadFinishedArgsForCall, struct {
		type_ string
		desc  string
		err   error
	}{type_, desc, err})
	fake.recordInvocation("IndexEntryUploadFinished", []interface{}{type_, desc, err})
	fake.indexEntryUploadFinishedMutex.Unlock()
	if fake.IndexEntryUploadFinishedStub != nil {
		fake.IndexEntryUploadFinishedStub(type_, desc, err)
	}
}

func (fake *FakeReporter) IndexEntryUploadFinishedCallCount() int {
	fake.indexEntryUploadFinishedMutex.RLock()
	defer fake.indexEntryUploadFinishedMutex.RUnlock()
	return len(fake.indexEntryUploadFinishedArgsForCall)
}

func (fake *FakeReporter) IndexEntryUploadFinishedArgsForCall(i int) (string, string, error) {
	fake.indexEntryUploadFinishedMutex.RLock()
	defer fake.indexEntryUploadFinishedMutex.RUnlock()
	return fake.indexEntryUploadFinishedArgsForCall[i].type_, fake.indexEntryUploadFinishedArgsForCall[i].desc, fake.indexEntryUploadFinishedArgsForCall[i].err
}

func (fake *FakeReporter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.indexEntryStartedAddingMutex.RLock()
	defer fake.indexEntryStartedAddingMutex.RUnlock()
	fake.indexEntryFinishedAddingMutex.RLock()
	defer fake.indexEntryFinishedAddingMutex.RUnlock()
	fake.indexEntryDownloadStartedMutex.RLock()
	defer fake.indexEntryDownloadStartedMutex.RUnlock()
	fake.indexEntryDownloadFinishedMutex.RLock()
	defer fake.indexEntryDownloadFinishedMutex.RUnlock()
	fake.indexEntryUploadStartedMutex.RLock()
	defer fake.indexEntryUploadStartedMutex.RUnlock()
	fake.indexEntryUploadFinishedMutex.RLock()
	defer fake.indexEntryUploadFinishedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeReporter) recordInvocation(key string, args []interface{}) {
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

var _ index.Reporter = new(FakeReporter)
