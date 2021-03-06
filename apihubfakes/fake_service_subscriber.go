// This file was generated by counterfeiter
package apihubfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/apihub/apihub"
)

type FakeServiceSubscriber struct {
	SubscribeStub        func(logger lager.Logger, prefix string, servicesCh chan apihub.ServiceSpec, stop <-chan struct{}) error
	subscribeMutex       sync.RWMutex
	subscribeArgsForCall []struct {
		logger     lager.Logger
		prefix     string
		servicesCh chan apihub.ServiceSpec
		stop       <-chan struct{}
	}
	subscribeReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeServiceSubscriber) Subscribe(logger lager.Logger, prefix string, servicesCh chan apihub.ServiceSpec, stop <-chan struct{}) error {
	fake.subscribeMutex.Lock()
	fake.subscribeArgsForCall = append(fake.subscribeArgsForCall, struct {
		logger     lager.Logger
		prefix     string
		servicesCh chan apihub.ServiceSpec
		stop       <-chan struct{}
	}{logger, prefix, servicesCh, stop})
	fake.recordInvocation("Subscribe", []interface{}{logger, prefix, servicesCh, stop})
	fake.subscribeMutex.Unlock()
	if fake.SubscribeStub != nil {
		return fake.SubscribeStub(logger, prefix, servicesCh, stop)
	} else {
		return fake.subscribeReturns.result1
	}
}

func (fake *FakeServiceSubscriber) SubscribeCallCount() int {
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	return len(fake.subscribeArgsForCall)
}

func (fake *FakeServiceSubscriber) SubscribeArgsForCall(i int) (lager.Logger, string, chan apihub.ServiceSpec, <-chan struct{}) {
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	return fake.subscribeArgsForCall[i].logger, fake.subscribeArgsForCall[i].prefix, fake.subscribeArgsForCall[i].servicesCh, fake.subscribeArgsForCall[i].stop
}

func (fake *FakeServiceSubscriber) SubscribeReturns(result1 error) {
	fake.SubscribeStub = nil
	fake.subscribeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeServiceSubscriber) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.subscribeMutex.RLock()
	defer fake.subscribeMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeServiceSubscriber) recordInvocation(key string, args []interface{}) {
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

var _ apihub.ServiceSubscriber = new(FakeServiceSubscriber)
