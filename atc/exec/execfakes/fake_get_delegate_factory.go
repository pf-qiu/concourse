// Code generated by counterfeiter. DO NOT EDIT.
package execfakes

import (
	"sync"

	"github.com/pf-qiu/concourse/v6/atc/exec"
)

type FakeGetDelegateFactory struct {
	GetDelegateStub        func(exec.RunState) exec.GetDelegate
	getDelegateMutex       sync.RWMutex
	getDelegateArgsForCall []struct {
		arg1 exec.RunState
	}
	getDelegateReturns struct {
		result1 exec.GetDelegate
	}
	getDelegateReturnsOnCall map[int]struct {
		result1 exec.GetDelegate
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGetDelegateFactory) GetDelegate(arg1 exec.RunState) exec.GetDelegate {
	fake.getDelegateMutex.Lock()
	ret, specificReturn := fake.getDelegateReturnsOnCall[len(fake.getDelegateArgsForCall)]
	fake.getDelegateArgsForCall = append(fake.getDelegateArgsForCall, struct {
		arg1 exec.RunState
	}{arg1})
	fake.recordInvocation("GetDelegate", []interface{}{arg1})
	fake.getDelegateMutex.Unlock()
	if fake.GetDelegateStub != nil {
		return fake.GetDelegateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getDelegateReturns
	return fakeReturns.result1
}

func (fake *FakeGetDelegateFactory) GetDelegateCallCount() int {
	fake.getDelegateMutex.RLock()
	defer fake.getDelegateMutex.RUnlock()
	return len(fake.getDelegateArgsForCall)
}

func (fake *FakeGetDelegateFactory) GetDelegateCalls(stub func(exec.RunState) exec.GetDelegate) {
	fake.getDelegateMutex.Lock()
	defer fake.getDelegateMutex.Unlock()
	fake.GetDelegateStub = stub
}

func (fake *FakeGetDelegateFactory) GetDelegateArgsForCall(i int) exec.RunState {
	fake.getDelegateMutex.RLock()
	defer fake.getDelegateMutex.RUnlock()
	argsForCall := fake.getDelegateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGetDelegateFactory) GetDelegateReturns(result1 exec.GetDelegate) {
	fake.getDelegateMutex.Lock()
	defer fake.getDelegateMutex.Unlock()
	fake.GetDelegateStub = nil
	fake.getDelegateReturns = struct {
		result1 exec.GetDelegate
	}{result1}
}

func (fake *FakeGetDelegateFactory) GetDelegateReturnsOnCall(i int, result1 exec.GetDelegate) {
	fake.getDelegateMutex.Lock()
	defer fake.getDelegateMutex.Unlock()
	fake.GetDelegateStub = nil
	if fake.getDelegateReturnsOnCall == nil {
		fake.getDelegateReturnsOnCall = make(map[int]struct {
			result1 exec.GetDelegate
		})
	}
	fake.getDelegateReturnsOnCall[i] = struct {
		result1 exec.GetDelegate
	}{result1}
}

func (fake *FakeGetDelegateFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getDelegateMutex.RLock()
	defer fake.getDelegateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGetDelegateFactory) recordInvocation(key string, args []interface{}) {
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

var _ exec.GetDelegateFactory = new(FakeGetDelegateFactory)
