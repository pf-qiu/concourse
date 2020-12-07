// Code generated by counterfeiter. DO NOT EDIT.
package credsfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/creds"
)

type FakeVarSourcePool struct {
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	FindOrCreateStub        func(lager.Logger, map[string]interface{}, creds.ManagerFactory) (creds.Secrets, error)
	findOrCreateMutex       sync.RWMutex
	findOrCreateArgsForCall []struct {
		arg1 lager.Logger
		arg2 map[string]interface{}
		arg3 creds.ManagerFactory
	}
	findOrCreateReturns struct {
		result1 creds.Secrets
		result2 error
	}
	findOrCreateReturnsOnCall map[int]struct {
		result1 creds.Secrets
		result2 error
	}
	SizeStub        func() int
	sizeMutex       sync.RWMutex
	sizeArgsForCall []struct {
	}
	sizeReturns struct {
		result1 int
	}
	sizeReturnsOnCall map[int]struct {
		result1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVarSourcePool) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeVarSourcePool) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeVarSourcePool) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *FakeVarSourcePool) FindOrCreate(arg1 lager.Logger, arg2 map[string]interface{}, arg3 creds.ManagerFactory) (creds.Secrets, error) {
	fake.findOrCreateMutex.Lock()
	ret, specificReturn := fake.findOrCreateReturnsOnCall[len(fake.findOrCreateArgsForCall)]
	fake.findOrCreateArgsForCall = append(fake.findOrCreateArgsForCall, struct {
		arg1 lager.Logger
		arg2 map[string]interface{}
		arg3 creds.ManagerFactory
	}{arg1, arg2, arg3})
	fake.recordInvocation("FindOrCreate", []interface{}{arg1, arg2, arg3})
	fake.findOrCreateMutex.Unlock()
	if fake.FindOrCreateStub != nil {
		return fake.FindOrCreateStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findOrCreateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVarSourcePool) FindOrCreateCallCount() int {
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	return len(fake.findOrCreateArgsForCall)
}

func (fake *FakeVarSourcePool) FindOrCreateCalls(stub func(lager.Logger, map[string]interface{}, creds.ManagerFactory) (creds.Secrets, error)) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = stub
}

func (fake *FakeVarSourcePool) FindOrCreateArgsForCall(i int) (lager.Logger, map[string]interface{}, creds.ManagerFactory) {
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	argsForCall := fake.findOrCreateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeVarSourcePool) FindOrCreateReturns(result1 creds.Secrets, result2 error) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = nil
	fake.findOrCreateReturns = struct {
		result1 creds.Secrets
		result2 error
	}{result1, result2}
}

func (fake *FakeVarSourcePool) FindOrCreateReturnsOnCall(i int, result1 creds.Secrets, result2 error) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = nil
	if fake.findOrCreateReturnsOnCall == nil {
		fake.findOrCreateReturnsOnCall = make(map[int]struct {
			result1 creds.Secrets
			result2 error
		})
	}
	fake.findOrCreateReturnsOnCall[i] = struct {
		result1 creds.Secrets
		result2 error
	}{result1, result2}
}

func (fake *FakeVarSourcePool) Size() int {
	fake.sizeMutex.Lock()
	ret, specificReturn := fake.sizeReturnsOnCall[len(fake.sizeArgsForCall)]
	fake.sizeArgsForCall = append(fake.sizeArgsForCall, struct {
	}{})
	fake.recordInvocation("Size", []interface{}{})
	fake.sizeMutex.Unlock()
	if fake.SizeStub != nil {
		return fake.SizeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sizeReturns
	return fakeReturns.result1
}

func (fake *FakeVarSourcePool) SizeCallCount() int {
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	return len(fake.sizeArgsForCall)
}

func (fake *FakeVarSourcePool) SizeCalls(stub func() int) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = stub
}

func (fake *FakeVarSourcePool) SizeReturns(result1 int) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = nil
	fake.sizeReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeVarSourcePool) SizeReturnsOnCall(i int, result1 int) {
	fake.sizeMutex.Lock()
	defer fake.sizeMutex.Unlock()
	fake.SizeStub = nil
	if fake.sizeReturnsOnCall == nil {
		fake.sizeReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.sizeReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeVarSourcePool) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	fake.sizeMutex.RLock()
	defer fake.sizeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVarSourcePool) recordInvocation(key string, args []interface{}) {
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

var _ creds.VarSourcePool = new(FakeVarSourcePool)
