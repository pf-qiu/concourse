// Code generated by counterfeiter. DO NOT EDIT.
package workerfakes

import (
	"context"
	"sync"

	"github.com/pf-qiu/concourse/v6/tsa"
	"github.com/pf-qiu/concourse/v6/worker"
)

type FakeTSAClient struct {
	ContainersToDestroyStub        func(context.Context) ([]string, error)
	containersToDestroyMutex       sync.RWMutex
	containersToDestroyArgsForCall []struct {
		arg1 context.Context
	}
	containersToDestroyReturns struct {
		result1 []string
		result2 error
	}
	containersToDestroyReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	DeleteStub        func(context.Context) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 context.Context
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	LandStub        func(context.Context) error
	landMutex       sync.RWMutex
	landArgsForCall []struct {
		arg1 context.Context
	}
	landReturns struct {
		result1 error
	}
	landReturnsOnCall map[int]struct {
		result1 error
	}
	RegisterStub        func(context.Context, tsa.RegisterOptions) error
	registerMutex       sync.RWMutex
	registerArgsForCall []struct {
		arg1 context.Context
		arg2 tsa.RegisterOptions
	}
	registerReturns struct {
		result1 error
	}
	registerReturnsOnCall map[int]struct {
		result1 error
	}
	ReportContainersStub        func(context.Context, []string) error
	reportContainersMutex       sync.RWMutex
	reportContainersArgsForCall []struct {
		arg1 context.Context
		arg2 []string
	}
	reportContainersReturns struct {
		result1 error
	}
	reportContainersReturnsOnCall map[int]struct {
		result1 error
	}
	ReportVolumesStub        func(context.Context, []string) error
	reportVolumesMutex       sync.RWMutex
	reportVolumesArgsForCall []struct {
		arg1 context.Context
		arg2 []string
	}
	reportVolumesReturns struct {
		result1 error
	}
	reportVolumesReturnsOnCall map[int]struct {
		result1 error
	}
	RetireStub        func(context.Context) error
	retireMutex       sync.RWMutex
	retireArgsForCall []struct {
		arg1 context.Context
	}
	retireReturns struct {
		result1 error
	}
	retireReturnsOnCall map[int]struct {
		result1 error
	}
	VolumesToDestroyStub        func(context.Context) ([]string, error)
	volumesToDestroyMutex       sync.RWMutex
	volumesToDestroyArgsForCall []struct {
		arg1 context.Context
	}
	volumesToDestroyReturns struct {
		result1 []string
		result2 error
	}
	volumesToDestroyReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTSAClient) ContainersToDestroy(arg1 context.Context) ([]string, error) {
	fake.containersToDestroyMutex.Lock()
	ret, specificReturn := fake.containersToDestroyReturnsOnCall[len(fake.containersToDestroyArgsForCall)]
	fake.containersToDestroyArgsForCall = append(fake.containersToDestroyArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("ContainersToDestroy", []interface{}{arg1})
	fake.containersToDestroyMutex.Unlock()
	if fake.ContainersToDestroyStub != nil {
		return fake.ContainersToDestroyStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.containersToDestroyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTSAClient) ContainersToDestroyCallCount() int {
	fake.containersToDestroyMutex.RLock()
	defer fake.containersToDestroyMutex.RUnlock()
	return len(fake.containersToDestroyArgsForCall)
}

func (fake *FakeTSAClient) ContainersToDestroyCalls(stub func(context.Context) ([]string, error)) {
	fake.containersToDestroyMutex.Lock()
	defer fake.containersToDestroyMutex.Unlock()
	fake.ContainersToDestroyStub = stub
}

func (fake *FakeTSAClient) ContainersToDestroyArgsForCall(i int) context.Context {
	fake.containersToDestroyMutex.RLock()
	defer fake.containersToDestroyMutex.RUnlock()
	argsForCall := fake.containersToDestroyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTSAClient) ContainersToDestroyReturns(result1 []string, result2 error) {
	fake.containersToDestroyMutex.Lock()
	defer fake.containersToDestroyMutex.Unlock()
	fake.ContainersToDestroyStub = nil
	fake.containersToDestroyReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeTSAClient) ContainersToDestroyReturnsOnCall(i int, result1 []string, result2 error) {
	fake.containersToDestroyMutex.Lock()
	defer fake.containersToDestroyMutex.Unlock()
	fake.ContainersToDestroyStub = nil
	if fake.containersToDestroyReturnsOnCall == nil {
		fake.containersToDestroyReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.containersToDestroyReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeTSAClient) Delete(arg1 context.Context) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeTSAClient) DeleteCalls(stub func(context.Context) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeTSAClient) DeleteArgsForCall(i int) context.Context {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTSAClient) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) Land(arg1 context.Context) error {
	fake.landMutex.Lock()
	ret, specificReturn := fake.landReturnsOnCall[len(fake.landArgsForCall)]
	fake.landArgsForCall = append(fake.landArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Land", []interface{}{arg1})
	fake.landMutex.Unlock()
	if fake.LandStub != nil {
		return fake.LandStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.landReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) LandCallCount() int {
	fake.landMutex.RLock()
	defer fake.landMutex.RUnlock()
	return len(fake.landArgsForCall)
}

func (fake *FakeTSAClient) LandCalls(stub func(context.Context) error) {
	fake.landMutex.Lock()
	defer fake.landMutex.Unlock()
	fake.LandStub = stub
}

func (fake *FakeTSAClient) LandArgsForCall(i int) context.Context {
	fake.landMutex.RLock()
	defer fake.landMutex.RUnlock()
	argsForCall := fake.landArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTSAClient) LandReturns(result1 error) {
	fake.landMutex.Lock()
	defer fake.landMutex.Unlock()
	fake.LandStub = nil
	fake.landReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) LandReturnsOnCall(i int, result1 error) {
	fake.landMutex.Lock()
	defer fake.landMutex.Unlock()
	fake.LandStub = nil
	if fake.landReturnsOnCall == nil {
		fake.landReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.landReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) Register(arg1 context.Context, arg2 tsa.RegisterOptions) error {
	fake.registerMutex.Lock()
	ret, specificReturn := fake.registerReturnsOnCall[len(fake.registerArgsForCall)]
	fake.registerArgsForCall = append(fake.registerArgsForCall, struct {
		arg1 context.Context
		arg2 tsa.RegisterOptions
	}{arg1, arg2})
	fake.recordInvocation("Register", []interface{}{arg1, arg2})
	fake.registerMutex.Unlock()
	if fake.RegisterStub != nil {
		return fake.RegisterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.registerReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) RegisterCallCount() int {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	return len(fake.registerArgsForCall)
}

func (fake *FakeTSAClient) RegisterCalls(stub func(context.Context, tsa.RegisterOptions) error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = stub
}

func (fake *FakeTSAClient) RegisterArgsForCall(i int) (context.Context, tsa.RegisterOptions) {
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	argsForCall := fake.registerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTSAClient) RegisterReturns(result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	fake.registerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) RegisterReturnsOnCall(i int, result1 error) {
	fake.registerMutex.Lock()
	defer fake.registerMutex.Unlock()
	fake.RegisterStub = nil
	if fake.registerReturnsOnCall == nil {
		fake.registerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.registerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) ReportContainers(arg1 context.Context, arg2 []string) error {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.reportContainersMutex.Lock()
	ret, specificReturn := fake.reportContainersReturnsOnCall[len(fake.reportContainersArgsForCall)]
	fake.reportContainersArgsForCall = append(fake.reportContainersArgsForCall, struct {
		arg1 context.Context
		arg2 []string
	}{arg1, arg2Copy})
	fake.recordInvocation("ReportContainers", []interface{}{arg1, arg2Copy})
	fake.reportContainersMutex.Unlock()
	if fake.ReportContainersStub != nil {
		return fake.ReportContainersStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.reportContainersReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) ReportContainersCallCount() int {
	fake.reportContainersMutex.RLock()
	defer fake.reportContainersMutex.RUnlock()
	return len(fake.reportContainersArgsForCall)
}

func (fake *FakeTSAClient) ReportContainersCalls(stub func(context.Context, []string) error) {
	fake.reportContainersMutex.Lock()
	defer fake.reportContainersMutex.Unlock()
	fake.ReportContainersStub = stub
}

func (fake *FakeTSAClient) ReportContainersArgsForCall(i int) (context.Context, []string) {
	fake.reportContainersMutex.RLock()
	defer fake.reportContainersMutex.RUnlock()
	argsForCall := fake.reportContainersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTSAClient) ReportContainersReturns(result1 error) {
	fake.reportContainersMutex.Lock()
	defer fake.reportContainersMutex.Unlock()
	fake.ReportContainersStub = nil
	fake.reportContainersReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) ReportContainersReturnsOnCall(i int, result1 error) {
	fake.reportContainersMutex.Lock()
	defer fake.reportContainersMutex.Unlock()
	fake.ReportContainersStub = nil
	if fake.reportContainersReturnsOnCall == nil {
		fake.reportContainersReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.reportContainersReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) ReportVolumes(arg1 context.Context, arg2 []string) error {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.reportVolumesMutex.Lock()
	ret, specificReturn := fake.reportVolumesReturnsOnCall[len(fake.reportVolumesArgsForCall)]
	fake.reportVolumesArgsForCall = append(fake.reportVolumesArgsForCall, struct {
		arg1 context.Context
		arg2 []string
	}{arg1, arg2Copy})
	fake.recordInvocation("ReportVolumes", []interface{}{arg1, arg2Copy})
	fake.reportVolumesMutex.Unlock()
	if fake.ReportVolumesStub != nil {
		return fake.ReportVolumesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.reportVolumesReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) ReportVolumesCallCount() int {
	fake.reportVolumesMutex.RLock()
	defer fake.reportVolumesMutex.RUnlock()
	return len(fake.reportVolumesArgsForCall)
}

func (fake *FakeTSAClient) ReportVolumesCalls(stub func(context.Context, []string) error) {
	fake.reportVolumesMutex.Lock()
	defer fake.reportVolumesMutex.Unlock()
	fake.ReportVolumesStub = stub
}

func (fake *FakeTSAClient) ReportVolumesArgsForCall(i int) (context.Context, []string) {
	fake.reportVolumesMutex.RLock()
	defer fake.reportVolumesMutex.RUnlock()
	argsForCall := fake.reportVolumesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTSAClient) ReportVolumesReturns(result1 error) {
	fake.reportVolumesMutex.Lock()
	defer fake.reportVolumesMutex.Unlock()
	fake.ReportVolumesStub = nil
	fake.reportVolumesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) ReportVolumesReturnsOnCall(i int, result1 error) {
	fake.reportVolumesMutex.Lock()
	defer fake.reportVolumesMutex.Unlock()
	fake.ReportVolumesStub = nil
	if fake.reportVolumesReturnsOnCall == nil {
		fake.reportVolumesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.reportVolumesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) Retire(arg1 context.Context) error {
	fake.retireMutex.Lock()
	ret, specificReturn := fake.retireReturnsOnCall[len(fake.retireArgsForCall)]
	fake.retireArgsForCall = append(fake.retireArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Retire", []interface{}{arg1})
	fake.retireMutex.Unlock()
	if fake.RetireStub != nil {
		return fake.RetireStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.retireReturns
	return fakeReturns.result1
}

func (fake *FakeTSAClient) RetireCallCount() int {
	fake.retireMutex.RLock()
	defer fake.retireMutex.RUnlock()
	return len(fake.retireArgsForCall)
}

func (fake *FakeTSAClient) RetireCalls(stub func(context.Context) error) {
	fake.retireMutex.Lock()
	defer fake.retireMutex.Unlock()
	fake.RetireStub = stub
}

func (fake *FakeTSAClient) RetireArgsForCall(i int) context.Context {
	fake.retireMutex.RLock()
	defer fake.retireMutex.RUnlock()
	argsForCall := fake.retireArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTSAClient) RetireReturns(result1 error) {
	fake.retireMutex.Lock()
	defer fake.retireMutex.Unlock()
	fake.RetireStub = nil
	fake.retireReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) RetireReturnsOnCall(i int, result1 error) {
	fake.retireMutex.Lock()
	defer fake.retireMutex.Unlock()
	fake.RetireStub = nil
	if fake.retireReturnsOnCall == nil {
		fake.retireReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.retireReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTSAClient) VolumesToDestroy(arg1 context.Context) ([]string, error) {
	fake.volumesToDestroyMutex.Lock()
	ret, specificReturn := fake.volumesToDestroyReturnsOnCall[len(fake.volumesToDestroyArgsForCall)]
	fake.volumesToDestroyArgsForCall = append(fake.volumesToDestroyArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("VolumesToDestroy", []interface{}{arg1})
	fake.volumesToDestroyMutex.Unlock()
	if fake.VolumesToDestroyStub != nil {
		return fake.VolumesToDestroyStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.volumesToDestroyReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTSAClient) VolumesToDestroyCallCount() int {
	fake.volumesToDestroyMutex.RLock()
	defer fake.volumesToDestroyMutex.RUnlock()
	return len(fake.volumesToDestroyArgsForCall)
}

func (fake *FakeTSAClient) VolumesToDestroyCalls(stub func(context.Context) ([]string, error)) {
	fake.volumesToDestroyMutex.Lock()
	defer fake.volumesToDestroyMutex.Unlock()
	fake.VolumesToDestroyStub = stub
}

func (fake *FakeTSAClient) VolumesToDestroyArgsForCall(i int) context.Context {
	fake.volumesToDestroyMutex.RLock()
	defer fake.volumesToDestroyMutex.RUnlock()
	argsForCall := fake.volumesToDestroyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTSAClient) VolumesToDestroyReturns(result1 []string, result2 error) {
	fake.volumesToDestroyMutex.Lock()
	defer fake.volumesToDestroyMutex.Unlock()
	fake.VolumesToDestroyStub = nil
	fake.volumesToDestroyReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeTSAClient) VolumesToDestroyReturnsOnCall(i int, result1 []string, result2 error) {
	fake.volumesToDestroyMutex.Lock()
	defer fake.volumesToDestroyMutex.Unlock()
	fake.VolumesToDestroyStub = nil
	if fake.volumesToDestroyReturnsOnCall == nil {
		fake.volumesToDestroyReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.volumesToDestroyReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeTSAClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.containersToDestroyMutex.RLock()
	defer fake.containersToDestroyMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.landMutex.RLock()
	defer fake.landMutex.RUnlock()
	fake.registerMutex.RLock()
	defer fake.registerMutex.RUnlock()
	fake.reportContainersMutex.RLock()
	defer fake.reportContainersMutex.RUnlock()
	fake.reportVolumesMutex.RLock()
	defer fake.reportVolumesMutex.RUnlock()
	fake.retireMutex.RLock()
	defer fake.retireMutex.RUnlock()
	fake.volumesToDestroyMutex.RLock()
	defer fake.volumesToDestroyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTSAClient) recordInvocation(key string, args []interface{}) {
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

var _ worker.TSAClient = new(FakeTSAClient)
