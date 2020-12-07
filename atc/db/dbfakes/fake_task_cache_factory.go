// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"

	"github.com/pf-qiu/concourse/v6/atc/db"
)

type FakeTaskCacheFactory struct {
	FindStub        func(int, string, string) (db.UsedTaskCache, bool, error)
	findMutex       sync.RWMutex
	findArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 string
	}
	findReturns struct {
		result1 db.UsedTaskCache
		result2 bool
		result3 error
	}
	findReturnsOnCall map[int]struct {
		result1 db.UsedTaskCache
		result2 bool
		result3 error
	}
	FindOrCreateStub        func(int, string, string) (db.UsedTaskCache, error)
	findOrCreateMutex       sync.RWMutex
	findOrCreateArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 string
	}
	findOrCreateReturns struct {
		result1 db.UsedTaskCache
		result2 error
	}
	findOrCreateReturnsOnCall map[int]struct {
		result1 db.UsedTaskCache
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTaskCacheFactory) Find(arg1 int, arg2 string, arg3 string) (db.UsedTaskCache, bool, error) {
	fake.findMutex.Lock()
	ret, specificReturn := fake.findReturnsOnCall[len(fake.findArgsForCall)]
	fake.findArgsForCall = append(fake.findArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("Find", []interface{}{arg1, arg2, arg3})
	fake.findMutex.Unlock()
	if fake.FindStub != nil {
		return fake.FindStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.findReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeTaskCacheFactory) FindCallCount() int {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return len(fake.findArgsForCall)
}

func (fake *FakeTaskCacheFactory) FindCalls(stub func(int, string, string) (db.UsedTaskCache, bool, error)) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = stub
}

func (fake *FakeTaskCacheFactory) FindArgsForCall(i int) (int, string, string) {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	argsForCall := fake.findArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeTaskCacheFactory) FindReturns(result1 db.UsedTaskCache, result2 bool, result3 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	fake.findReturns = struct {
		result1 db.UsedTaskCache
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTaskCacheFactory) FindReturnsOnCall(i int, result1 db.UsedTaskCache, result2 bool, result3 error) {
	fake.findMutex.Lock()
	defer fake.findMutex.Unlock()
	fake.FindStub = nil
	if fake.findReturnsOnCall == nil {
		fake.findReturnsOnCall = make(map[int]struct {
			result1 db.UsedTaskCache
			result2 bool
			result3 error
		})
	}
	fake.findReturnsOnCall[i] = struct {
		result1 db.UsedTaskCache
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTaskCacheFactory) FindOrCreate(arg1 int, arg2 string, arg3 string) (db.UsedTaskCache, error) {
	fake.findOrCreateMutex.Lock()
	ret, specificReturn := fake.findOrCreateReturnsOnCall[len(fake.findOrCreateArgsForCall)]
	fake.findOrCreateArgsForCall = append(fake.findOrCreateArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 string
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

func (fake *FakeTaskCacheFactory) FindOrCreateCallCount() int {
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	return len(fake.findOrCreateArgsForCall)
}

func (fake *FakeTaskCacheFactory) FindOrCreateCalls(stub func(int, string, string) (db.UsedTaskCache, error)) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = stub
}

func (fake *FakeTaskCacheFactory) FindOrCreateArgsForCall(i int) (int, string, string) {
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	argsForCall := fake.findOrCreateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeTaskCacheFactory) FindOrCreateReturns(result1 db.UsedTaskCache, result2 error) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = nil
	fake.findOrCreateReturns = struct {
		result1 db.UsedTaskCache
		result2 error
	}{result1, result2}
}

func (fake *FakeTaskCacheFactory) FindOrCreateReturnsOnCall(i int, result1 db.UsedTaskCache, result2 error) {
	fake.findOrCreateMutex.Lock()
	defer fake.findOrCreateMutex.Unlock()
	fake.FindOrCreateStub = nil
	if fake.findOrCreateReturnsOnCall == nil {
		fake.findOrCreateReturnsOnCall = make(map[int]struct {
			result1 db.UsedTaskCache
			result2 error
		})
	}
	fake.findOrCreateReturnsOnCall[i] = struct {
		result1 db.UsedTaskCache
		result2 error
	}{result1, result2}
}

func (fake *FakeTaskCacheFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	fake.findOrCreateMutex.RLock()
	defer fake.findOrCreateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTaskCacheFactory) recordInvocation(key string, args []interface{}) {
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

var _ db.TaskCacheFactory = new(FakeTaskCacheFactory)
