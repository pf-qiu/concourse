// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"
	"time"

	"github.com/pf-qiu/concourse/v6/atc/db"
)

type FakeCreatedContainer struct {
	DestroyingStub        func() (db.DestroyingContainer, error)
	destroyingMutex       sync.RWMutex
	destroyingArgsForCall []struct {
	}
	destroyingReturns struct {
		result1 db.DestroyingContainer
		result2 error
	}
	destroyingReturnsOnCall map[int]struct {
		result1 db.DestroyingContainer
		result2 error
	}
	HandleStub        func() string
	handleMutex       sync.RWMutex
	handleArgsForCall []struct {
	}
	handleReturns struct {
		result1 string
	}
	handleReturnsOnCall map[int]struct {
		result1 string
	}
	IDStub        func() int
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 int
	}
	iDReturnsOnCall map[int]struct {
		result1 int
	}
	LastHijackStub        func() time.Time
	lastHijackMutex       sync.RWMutex
	lastHijackArgsForCall []struct {
	}
	lastHijackReturns struct {
		result1 time.Time
	}
	lastHijackReturnsOnCall map[int]struct {
		result1 time.Time
	}
	MetadataStub        func() db.ContainerMetadata
	metadataMutex       sync.RWMutex
	metadataArgsForCall []struct {
	}
	metadataReturns struct {
		result1 db.ContainerMetadata
	}
	metadataReturnsOnCall map[int]struct {
		result1 db.ContainerMetadata
	}
	StateStub        func() string
	stateMutex       sync.RWMutex
	stateArgsForCall []struct {
	}
	stateReturns struct {
		result1 string
	}
	stateReturnsOnCall map[int]struct {
		result1 string
	}
	UpdateLastHijackStub        func() error
	updateLastHijackMutex       sync.RWMutex
	updateLastHijackArgsForCall []struct {
	}
	updateLastHijackReturns struct {
		result1 error
	}
	updateLastHijackReturnsOnCall map[int]struct {
		result1 error
	}
	WorkerNameStub        func() string
	workerNameMutex       sync.RWMutex
	workerNameArgsForCall []struct {
	}
	workerNameReturns struct {
		result1 string
	}
	workerNameReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCreatedContainer) Destroying() (db.DestroyingContainer, error) {
	fake.destroyingMutex.Lock()
	ret, specificReturn := fake.destroyingReturnsOnCall[len(fake.destroyingArgsForCall)]
	fake.destroyingArgsForCall = append(fake.destroyingArgsForCall, struct {
	}{})
	fake.recordInvocation("Destroying", []interface{}{})
	fake.destroyingMutex.Unlock()
	if fake.DestroyingStub != nil {
		return fake.DestroyingStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.destroyingReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCreatedContainer) DestroyingCallCount() int {
	fake.destroyingMutex.RLock()
	defer fake.destroyingMutex.RUnlock()
	return len(fake.destroyingArgsForCall)
}

func (fake *FakeCreatedContainer) DestroyingCalls(stub func() (db.DestroyingContainer, error)) {
	fake.destroyingMutex.Lock()
	defer fake.destroyingMutex.Unlock()
	fake.DestroyingStub = stub
}

func (fake *FakeCreatedContainer) DestroyingReturns(result1 db.DestroyingContainer, result2 error) {
	fake.destroyingMutex.Lock()
	defer fake.destroyingMutex.Unlock()
	fake.DestroyingStub = nil
	fake.destroyingReturns = struct {
		result1 db.DestroyingContainer
		result2 error
	}{result1, result2}
}

func (fake *FakeCreatedContainer) DestroyingReturnsOnCall(i int, result1 db.DestroyingContainer, result2 error) {
	fake.destroyingMutex.Lock()
	defer fake.destroyingMutex.Unlock()
	fake.DestroyingStub = nil
	if fake.destroyingReturnsOnCall == nil {
		fake.destroyingReturnsOnCall = make(map[int]struct {
			result1 db.DestroyingContainer
			result2 error
		})
	}
	fake.destroyingReturnsOnCall[i] = struct {
		result1 db.DestroyingContainer
		result2 error
	}{result1, result2}
}

func (fake *FakeCreatedContainer) Handle() string {
	fake.handleMutex.Lock()
	ret, specificReturn := fake.handleReturnsOnCall[len(fake.handleArgsForCall)]
	fake.handleArgsForCall = append(fake.handleArgsForCall, struct {
	}{})
	fake.recordInvocation("Handle", []interface{}{})
	fake.handleMutex.Unlock()
	if fake.HandleStub != nil {
		return fake.HandleStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.handleReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) HandleCallCount() int {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	return len(fake.handleArgsForCall)
}

func (fake *FakeCreatedContainer) HandleCalls(stub func() string) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = stub
}

func (fake *FakeCreatedContainer) HandleReturns(result1 string) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	fake.handleReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) HandleReturnsOnCall(i int, result1 string) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	if fake.handleReturnsOnCall == nil {
		fake.handleReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.handleReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) ID() int {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
	}{})
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.iDReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *FakeCreatedContainer) IDCalls(stub func() int) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *FakeCreatedContainer) IDReturns(result1 int) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeCreatedContainer) IDReturnsOnCall(i int, result1 int) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeCreatedContainer) LastHijack() time.Time {
	fake.lastHijackMutex.Lock()
	ret, specificReturn := fake.lastHijackReturnsOnCall[len(fake.lastHijackArgsForCall)]
	fake.lastHijackArgsForCall = append(fake.lastHijackArgsForCall, struct {
	}{})
	fake.recordInvocation("LastHijack", []interface{}{})
	fake.lastHijackMutex.Unlock()
	if fake.LastHijackStub != nil {
		return fake.LastHijackStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.lastHijackReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) LastHijackCallCount() int {
	fake.lastHijackMutex.RLock()
	defer fake.lastHijackMutex.RUnlock()
	return len(fake.lastHijackArgsForCall)
}

func (fake *FakeCreatedContainer) LastHijackCalls(stub func() time.Time) {
	fake.lastHijackMutex.Lock()
	defer fake.lastHijackMutex.Unlock()
	fake.LastHijackStub = stub
}

func (fake *FakeCreatedContainer) LastHijackReturns(result1 time.Time) {
	fake.lastHijackMutex.Lock()
	defer fake.lastHijackMutex.Unlock()
	fake.LastHijackStub = nil
	fake.lastHijackReturns = struct {
		result1 time.Time
	}{result1}
}

func (fake *FakeCreatedContainer) LastHijackReturnsOnCall(i int, result1 time.Time) {
	fake.lastHijackMutex.Lock()
	defer fake.lastHijackMutex.Unlock()
	fake.LastHijackStub = nil
	if fake.lastHijackReturnsOnCall == nil {
		fake.lastHijackReturnsOnCall = make(map[int]struct {
			result1 time.Time
		})
	}
	fake.lastHijackReturnsOnCall[i] = struct {
		result1 time.Time
	}{result1}
}

func (fake *FakeCreatedContainer) Metadata() db.ContainerMetadata {
	fake.metadataMutex.Lock()
	ret, specificReturn := fake.metadataReturnsOnCall[len(fake.metadataArgsForCall)]
	fake.metadataArgsForCall = append(fake.metadataArgsForCall, struct {
	}{})
	fake.recordInvocation("Metadata", []interface{}{})
	fake.metadataMutex.Unlock()
	if fake.MetadataStub != nil {
		return fake.MetadataStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.metadataReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) MetadataCallCount() int {
	fake.metadataMutex.RLock()
	defer fake.metadataMutex.RUnlock()
	return len(fake.metadataArgsForCall)
}

func (fake *FakeCreatedContainer) MetadataCalls(stub func() db.ContainerMetadata) {
	fake.metadataMutex.Lock()
	defer fake.metadataMutex.Unlock()
	fake.MetadataStub = stub
}

func (fake *FakeCreatedContainer) MetadataReturns(result1 db.ContainerMetadata) {
	fake.metadataMutex.Lock()
	defer fake.metadataMutex.Unlock()
	fake.MetadataStub = nil
	fake.metadataReturns = struct {
		result1 db.ContainerMetadata
	}{result1}
}

func (fake *FakeCreatedContainer) MetadataReturnsOnCall(i int, result1 db.ContainerMetadata) {
	fake.metadataMutex.Lock()
	defer fake.metadataMutex.Unlock()
	fake.MetadataStub = nil
	if fake.metadataReturnsOnCall == nil {
		fake.metadataReturnsOnCall = make(map[int]struct {
			result1 db.ContainerMetadata
		})
	}
	fake.metadataReturnsOnCall[i] = struct {
		result1 db.ContainerMetadata
	}{result1}
}

func (fake *FakeCreatedContainer) State() string {
	fake.stateMutex.Lock()
	ret, specificReturn := fake.stateReturnsOnCall[len(fake.stateArgsForCall)]
	fake.stateArgsForCall = append(fake.stateArgsForCall, struct {
	}{})
	fake.recordInvocation("State", []interface{}{})
	fake.stateMutex.Unlock()
	if fake.StateStub != nil {
		return fake.StateStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.stateReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) StateCallCount() int {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return len(fake.stateArgsForCall)
}

func (fake *FakeCreatedContainer) StateCalls(stub func() string) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = stub
}

func (fake *FakeCreatedContainer) StateReturns(result1 string) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = nil
	fake.stateReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) StateReturnsOnCall(i int, result1 string) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = nil
	if fake.stateReturnsOnCall == nil {
		fake.stateReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.stateReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) UpdateLastHijack() error {
	fake.updateLastHijackMutex.Lock()
	ret, specificReturn := fake.updateLastHijackReturnsOnCall[len(fake.updateLastHijackArgsForCall)]
	fake.updateLastHijackArgsForCall = append(fake.updateLastHijackArgsForCall, struct {
	}{})
	fake.recordInvocation("UpdateLastHijack", []interface{}{})
	fake.updateLastHijackMutex.Unlock()
	if fake.UpdateLastHijackStub != nil {
		return fake.UpdateLastHijackStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.updateLastHijackReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) UpdateLastHijackCallCount() int {
	fake.updateLastHijackMutex.RLock()
	defer fake.updateLastHijackMutex.RUnlock()
	return len(fake.updateLastHijackArgsForCall)
}

func (fake *FakeCreatedContainer) UpdateLastHijackCalls(stub func() error) {
	fake.updateLastHijackMutex.Lock()
	defer fake.updateLastHijackMutex.Unlock()
	fake.UpdateLastHijackStub = stub
}

func (fake *FakeCreatedContainer) UpdateLastHijackReturns(result1 error) {
	fake.updateLastHijackMutex.Lock()
	defer fake.updateLastHijackMutex.Unlock()
	fake.UpdateLastHijackStub = nil
	fake.updateLastHijackReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCreatedContainer) UpdateLastHijackReturnsOnCall(i int, result1 error) {
	fake.updateLastHijackMutex.Lock()
	defer fake.updateLastHijackMutex.Unlock()
	fake.UpdateLastHijackStub = nil
	if fake.updateLastHijackReturnsOnCall == nil {
		fake.updateLastHijackReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateLastHijackReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCreatedContainer) WorkerName() string {
	fake.workerNameMutex.Lock()
	ret, specificReturn := fake.workerNameReturnsOnCall[len(fake.workerNameArgsForCall)]
	fake.workerNameArgsForCall = append(fake.workerNameArgsForCall, struct {
	}{})
	fake.recordInvocation("WorkerName", []interface{}{})
	fake.workerNameMutex.Unlock()
	if fake.WorkerNameStub != nil {
		return fake.WorkerNameStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.workerNameReturns
	return fakeReturns.result1
}

func (fake *FakeCreatedContainer) WorkerNameCallCount() int {
	fake.workerNameMutex.RLock()
	defer fake.workerNameMutex.RUnlock()
	return len(fake.workerNameArgsForCall)
}

func (fake *FakeCreatedContainer) WorkerNameCalls(stub func() string) {
	fake.workerNameMutex.Lock()
	defer fake.workerNameMutex.Unlock()
	fake.WorkerNameStub = stub
}

func (fake *FakeCreatedContainer) WorkerNameReturns(result1 string) {
	fake.workerNameMutex.Lock()
	defer fake.workerNameMutex.Unlock()
	fake.WorkerNameStub = nil
	fake.workerNameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) WorkerNameReturnsOnCall(i int, result1 string) {
	fake.workerNameMutex.Lock()
	defer fake.workerNameMutex.Unlock()
	fake.WorkerNameStub = nil
	if fake.workerNameReturnsOnCall == nil {
		fake.workerNameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.workerNameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeCreatedContainer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.destroyingMutex.RLock()
	defer fake.destroyingMutex.RUnlock()
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.lastHijackMutex.RLock()
	defer fake.lastHijackMutex.RUnlock()
	fake.metadataMutex.RLock()
	defer fake.metadataMutex.RUnlock()
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	fake.updateLastHijackMutex.RLock()
	defer fake.updateLastHijackMutex.RUnlock()
	fake.workerNameMutex.RLock()
	defer fake.workerNameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCreatedContainer) recordInvocation(key string, args []interface{}) {
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

var _ db.CreatedContainer = new(FakeCreatedContainer)
