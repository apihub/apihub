// This file was generated by counterfeiter
package apihubfakes

import (
	"sync"

	"github.com/apihub/apihub"
)

type FakeStorage struct {
	UpsertServiceStub        func(apihub.ServiceSpec) error
	upsertServiceMutex       sync.RWMutex
	upsertServiceArgsForCall []struct {
		arg1 apihub.ServiceSpec
	}
	upsertServiceReturns struct {
		result1 error
	}
	FindServiceByHandleStub        func(string) (apihub.ServiceSpec, error)
	findServiceByHandleMutex       sync.RWMutex
	findServiceByHandleArgsForCall []struct {
		arg1 string
	}
	findServiceByHandleReturns struct {
		result1 apihub.ServiceSpec
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStorage) UpsertService(arg1 apihub.ServiceSpec) error {
	fake.upsertServiceMutex.Lock()
	fake.upsertServiceArgsForCall = append(fake.upsertServiceArgsForCall, struct {
		arg1 apihub.ServiceSpec
	}{arg1})
	fake.recordInvocation("UpsertService", []interface{}{arg1})
	fake.upsertServiceMutex.Unlock()
	if fake.UpsertServiceStub != nil {
		return fake.UpsertServiceStub(arg1)
	} else {
		return fake.upsertServiceReturns.result1
	}
}

func (fake *FakeStorage) UpsertServiceCallCount() int {
	fake.upsertServiceMutex.RLock()
	defer fake.upsertServiceMutex.RUnlock()
	return len(fake.upsertServiceArgsForCall)
}

func (fake *FakeStorage) UpsertServiceArgsForCall(i int) apihub.ServiceSpec {
	fake.upsertServiceMutex.RLock()
	defer fake.upsertServiceMutex.RUnlock()
	return fake.upsertServiceArgsForCall[i].arg1
}

func (fake *FakeStorage) UpsertServiceReturns(result1 error) {
	fake.UpsertServiceStub = nil
	fake.upsertServiceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStorage) FindServiceByHandle(arg1 string) (apihub.ServiceSpec, error) {
	fake.findServiceByHandleMutex.Lock()
	fake.findServiceByHandleArgsForCall = append(fake.findServiceByHandleArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindServiceByHandle", []interface{}{arg1})
	fake.findServiceByHandleMutex.Unlock()
	if fake.FindServiceByHandleStub != nil {
		return fake.FindServiceByHandleStub(arg1)
	} else {
		return fake.findServiceByHandleReturns.result1, fake.findServiceByHandleReturns.result2
	}
}

func (fake *FakeStorage) FindServiceByHandleCallCount() int {
	fake.findServiceByHandleMutex.RLock()
	defer fake.findServiceByHandleMutex.RUnlock()
	return len(fake.findServiceByHandleArgsForCall)
}

func (fake *FakeStorage) FindServiceByHandleArgsForCall(i int) string {
	fake.findServiceByHandleMutex.RLock()
	defer fake.findServiceByHandleMutex.RUnlock()
	return fake.findServiceByHandleArgsForCall[i].arg1
}

func (fake *FakeStorage) FindServiceByHandleReturns(result1 apihub.ServiceSpec, result2 error) {
	fake.FindServiceByHandleStub = nil
	fake.findServiceByHandleReturns = struct {
		result1 apihub.ServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeStorage) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.upsertServiceMutex.RLock()
	defer fake.upsertServiceMutex.RUnlock()
	fake.findServiceByHandleMutex.RLock()
	defer fake.findServiceByHandleMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeStorage) recordInvocation(key string, args []interface{}) {
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

var _ apihub.Storage = new(FakeStorage)
