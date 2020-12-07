// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

type FakeTeamFactory struct {
	CreateDefaultTeamIfNotExistsStub        func() (db.Team, error)
	createDefaultTeamIfNotExistsMutex       sync.RWMutex
	createDefaultTeamIfNotExistsArgsForCall []struct {
	}
	createDefaultTeamIfNotExistsReturns struct {
		result1 db.Team
		result2 error
	}
	createDefaultTeamIfNotExistsReturnsOnCall map[int]struct {
		result1 db.Team
		result2 error
	}
	CreateTeamStub        func(atc.Team) (db.Team, error)
	createTeamMutex       sync.RWMutex
	createTeamArgsForCall []struct {
		arg1 atc.Team
	}
	createTeamReturns struct {
		result1 db.Team
		result2 error
	}
	createTeamReturnsOnCall map[int]struct {
		result1 db.Team
		result2 error
	}
	FindTeamStub        func(string) (db.Team, bool, error)
	findTeamMutex       sync.RWMutex
	findTeamArgsForCall []struct {
		arg1 string
	}
	findTeamReturns struct {
		result1 db.Team
		result2 bool
		result3 error
	}
	findTeamReturnsOnCall map[int]struct {
		result1 db.Team
		result2 bool
		result3 error
	}
	GetByIDStub        func(int) db.Team
	getByIDMutex       sync.RWMutex
	getByIDArgsForCall []struct {
		arg1 int
	}
	getByIDReturns struct {
		result1 db.Team
	}
	getByIDReturnsOnCall map[int]struct {
		result1 db.Team
	}
	GetTeamsStub        func() ([]db.Team, error)
	getTeamsMutex       sync.RWMutex
	getTeamsArgsForCall []struct {
	}
	getTeamsReturns struct {
		result1 []db.Team
		result2 error
	}
	getTeamsReturnsOnCall map[int]struct {
		result1 []db.Team
		result2 error
	}
	NotifyCacherStub        func() error
	notifyCacherMutex       sync.RWMutex
	notifyCacherArgsForCall []struct {
	}
	notifyCacherReturns struct {
		result1 error
	}
	notifyCacherReturnsOnCall map[int]struct {
		result1 error
	}
	NotifyResourceScannerStub        func() error
	notifyResourceScannerMutex       sync.RWMutex
	notifyResourceScannerArgsForCall []struct {
	}
	notifyResourceScannerReturns struct {
		result1 error
	}
	notifyResourceScannerReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTeamFactory) CreateDefaultTeamIfNotExists() (db.Team, error) {
	fake.createDefaultTeamIfNotExistsMutex.Lock()
	ret, specificReturn := fake.createDefaultTeamIfNotExistsReturnsOnCall[len(fake.createDefaultTeamIfNotExistsArgsForCall)]
	fake.createDefaultTeamIfNotExistsArgsForCall = append(fake.createDefaultTeamIfNotExistsArgsForCall, struct {
	}{})
	fake.recordInvocation("CreateDefaultTeamIfNotExists", []interface{}{})
	fake.createDefaultTeamIfNotExistsMutex.Unlock()
	if fake.CreateDefaultTeamIfNotExistsStub != nil {
		return fake.CreateDefaultTeamIfNotExistsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createDefaultTeamIfNotExistsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTeamFactory) CreateDefaultTeamIfNotExistsCallCount() int {
	fake.createDefaultTeamIfNotExistsMutex.RLock()
	defer fake.createDefaultTeamIfNotExistsMutex.RUnlock()
	return len(fake.createDefaultTeamIfNotExistsArgsForCall)
}

func (fake *FakeTeamFactory) CreateDefaultTeamIfNotExistsCalls(stub func() (db.Team, error)) {
	fake.createDefaultTeamIfNotExistsMutex.Lock()
	defer fake.createDefaultTeamIfNotExistsMutex.Unlock()
	fake.CreateDefaultTeamIfNotExistsStub = stub
}

func (fake *FakeTeamFactory) CreateDefaultTeamIfNotExistsReturns(result1 db.Team, result2 error) {
	fake.createDefaultTeamIfNotExistsMutex.Lock()
	defer fake.createDefaultTeamIfNotExistsMutex.Unlock()
	fake.CreateDefaultTeamIfNotExistsStub = nil
	fake.createDefaultTeamIfNotExistsReturns = struct {
		result1 db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) CreateDefaultTeamIfNotExistsReturnsOnCall(i int, result1 db.Team, result2 error) {
	fake.createDefaultTeamIfNotExistsMutex.Lock()
	defer fake.createDefaultTeamIfNotExistsMutex.Unlock()
	fake.CreateDefaultTeamIfNotExistsStub = nil
	if fake.createDefaultTeamIfNotExistsReturnsOnCall == nil {
		fake.createDefaultTeamIfNotExistsReturnsOnCall = make(map[int]struct {
			result1 db.Team
			result2 error
		})
	}
	fake.createDefaultTeamIfNotExistsReturnsOnCall[i] = struct {
		result1 db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) CreateTeam(arg1 atc.Team) (db.Team, error) {
	fake.createTeamMutex.Lock()
	ret, specificReturn := fake.createTeamReturnsOnCall[len(fake.createTeamArgsForCall)]
	fake.createTeamArgsForCall = append(fake.createTeamArgsForCall, struct {
		arg1 atc.Team
	}{arg1})
	fake.recordInvocation("CreateTeam", []interface{}{arg1})
	fake.createTeamMutex.Unlock()
	if fake.CreateTeamStub != nil {
		return fake.CreateTeamStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createTeamReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTeamFactory) CreateTeamCallCount() int {
	fake.createTeamMutex.RLock()
	defer fake.createTeamMutex.RUnlock()
	return len(fake.createTeamArgsForCall)
}

func (fake *FakeTeamFactory) CreateTeamCalls(stub func(atc.Team) (db.Team, error)) {
	fake.createTeamMutex.Lock()
	defer fake.createTeamMutex.Unlock()
	fake.CreateTeamStub = stub
}

func (fake *FakeTeamFactory) CreateTeamArgsForCall(i int) atc.Team {
	fake.createTeamMutex.RLock()
	defer fake.createTeamMutex.RUnlock()
	argsForCall := fake.createTeamArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTeamFactory) CreateTeamReturns(result1 db.Team, result2 error) {
	fake.createTeamMutex.Lock()
	defer fake.createTeamMutex.Unlock()
	fake.CreateTeamStub = nil
	fake.createTeamReturns = struct {
		result1 db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) CreateTeamReturnsOnCall(i int, result1 db.Team, result2 error) {
	fake.createTeamMutex.Lock()
	defer fake.createTeamMutex.Unlock()
	fake.CreateTeamStub = nil
	if fake.createTeamReturnsOnCall == nil {
		fake.createTeamReturnsOnCall = make(map[int]struct {
			result1 db.Team
			result2 error
		})
	}
	fake.createTeamReturnsOnCall[i] = struct {
		result1 db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) FindTeam(arg1 string) (db.Team, bool, error) {
	fake.findTeamMutex.Lock()
	ret, specificReturn := fake.findTeamReturnsOnCall[len(fake.findTeamArgsForCall)]
	fake.findTeamArgsForCall = append(fake.findTeamArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindTeam", []interface{}{arg1})
	fake.findTeamMutex.Unlock()
	if fake.FindTeamStub != nil {
		return fake.FindTeamStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.findTeamReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeTeamFactory) FindTeamCallCount() int {
	fake.findTeamMutex.RLock()
	defer fake.findTeamMutex.RUnlock()
	return len(fake.findTeamArgsForCall)
}

func (fake *FakeTeamFactory) FindTeamCalls(stub func(string) (db.Team, bool, error)) {
	fake.findTeamMutex.Lock()
	defer fake.findTeamMutex.Unlock()
	fake.FindTeamStub = stub
}

func (fake *FakeTeamFactory) FindTeamArgsForCall(i int) string {
	fake.findTeamMutex.RLock()
	defer fake.findTeamMutex.RUnlock()
	argsForCall := fake.findTeamArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTeamFactory) FindTeamReturns(result1 db.Team, result2 bool, result3 error) {
	fake.findTeamMutex.Lock()
	defer fake.findTeamMutex.Unlock()
	fake.FindTeamStub = nil
	fake.findTeamReturns = struct {
		result1 db.Team
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamFactory) FindTeamReturnsOnCall(i int, result1 db.Team, result2 bool, result3 error) {
	fake.findTeamMutex.Lock()
	defer fake.findTeamMutex.Unlock()
	fake.FindTeamStub = nil
	if fake.findTeamReturnsOnCall == nil {
		fake.findTeamReturnsOnCall = make(map[int]struct {
			result1 db.Team
			result2 bool
			result3 error
		})
	}
	fake.findTeamReturnsOnCall[i] = struct {
		result1 db.Team
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTeamFactory) GetByID(arg1 int) db.Team {
	fake.getByIDMutex.Lock()
	ret, specificReturn := fake.getByIDReturnsOnCall[len(fake.getByIDArgsForCall)]
	fake.getByIDArgsForCall = append(fake.getByIDArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("GetByID", []interface{}{arg1})
	fake.getByIDMutex.Unlock()
	if fake.GetByIDStub != nil {
		return fake.GetByIDStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getByIDReturns
	return fakeReturns.result1
}

func (fake *FakeTeamFactory) GetByIDCallCount() int {
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	return len(fake.getByIDArgsForCall)
}

func (fake *FakeTeamFactory) GetByIDCalls(stub func(int) db.Team) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = stub
}

func (fake *FakeTeamFactory) GetByIDArgsForCall(i int) int {
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	argsForCall := fake.getByIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTeamFactory) GetByIDReturns(result1 db.Team) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = nil
	fake.getByIDReturns = struct {
		result1 db.Team
	}{result1}
}

func (fake *FakeTeamFactory) GetByIDReturnsOnCall(i int, result1 db.Team) {
	fake.getByIDMutex.Lock()
	defer fake.getByIDMutex.Unlock()
	fake.GetByIDStub = nil
	if fake.getByIDReturnsOnCall == nil {
		fake.getByIDReturnsOnCall = make(map[int]struct {
			result1 db.Team
		})
	}
	fake.getByIDReturnsOnCall[i] = struct {
		result1 db.Team
	}{result1}
}

func (fake *FakeTeamFactory) GetTeams() ([]db.Team, error) {
	fake.getTeamsMutex.Lock()
	ret, specificReturn := fake.getTeamsReturnsOnCall[len(fake.getTeamsArgsForCall)]
	fake.getTeamsArgsForCall = append(fake.getTeamsArgsForCall, struct {
	}{})
	fake.recordInvocation("GetTeams", []interface{}{})
	fake.getTeamsMutex.Unlock()
	if fake.GetTeamsStub != nil {
		return fake.GetTeamsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getTeamsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTeamFactory) GetTeamsCallCount() int {
	fake.getTeamsMutex.RLock()
	defer fake.getTeamsMutex.RUnlock()
	return len(fake.getTeamsArgsForCall)
}

func (fake *FakeTeamFactory) GetTeamsCalls(stub func() ([]db.Team, error)) {
	fake.getTeamsMutex.Lock()
	defer fake.getTeamsMutex.Unlock()
	fake.GetTeamsStub = stub
}

func (fake *FakeTeamFactory) GetTeamsReturns(result1 []db.Team, result2 error) {
	fake.getTeamsMutex.Lock()
	defer fake.getTeamsMutex.Unlock()
	fake.GetTeamsStub = nil
	fake.getTeamsReturns = struct {
		result1 []db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) GetTeamsReturnsOnCall(i int, result1 []db.Team, result2 error) {
	fake.getTeamsMutex.Lock()
	defer fake.getTeamsMutex.Unlock()
	fake.GetTeamsStub = nil
	if fake.getTeamsReturnsOnCall == nil {
		fake.getTeamsReturnsOnCall = make(map[int]struct {
			result1 []db.Team
			result2 error
		})
	}
	fake.getTeamsReturnsOnCall[i] = struct {
		result1 []db.Team
		result2 error
	}{result1, result2}
}

func (fake *FakeTeamFactory) NotifyCacher() error {
	fake.notifyCacherMutex.Lock()
	ret, specificReturn := fake.notifyCacherReturnsOnCall[len(fake.notifyCacherArgsForCall)]
	fake.notifyCacherArgsForCall = append(fake.notifyCacherArgsForCall, struct {
	}{})
	fake.recordInvocation("NotifyCacher", []interface{}{})
	fake.notifyCacherMutex.Unlock()
	if fake.NotifyCacherStub != nil {
		return fake.NotifyCacherStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.notifyCacherReturns
	return fakeReturns.result1
}

func (fake *FakeTeamFactory) NotifyCacherCallCount() int {
	fake.notifyCacherMutex.RLock()
	defer fake.notifyCacherMutex.RUnlock()
	return len(fake.notifyCacherArgsForCall)
}

func (fake *FakeTeamFactory) NotifyCacherCalls(stub func() error) {
	fake.notifyCacherMutex.Lock()
	defer fake.notifyCacherMutex.Unlock()
	fake.NotifyCacherStub = stub
}

func (fake *FakeTeamFactory) NotifyCacherReturns(result1 error) {
	fake.notifyCacherMutex.Lock()
	defer fake.notifyCacherMutex.Unlock()
	fake.NotifyCacherStub = nil
	fake.notifyCacherReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTeamFactory) NotifyCacherReturnsOnCall(i int, result1 error) {
	fake.notifyCacherMutex.Lock()
	defer fake.notifyCacherMutex.Unlock()
	fake.NotifyCacherStub = nil
	if fake.notifyCacherReturnsOnCall == nil {
		fake.notifyCacherReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.notifyCacherReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTeamFactory) NotifyResourceScanner() error {
	fake.notifyResourceScannerMutex.Lock()
	ret, specificReturn := fake.notifyResourceScannerReturnsOnCall[len(fake.notifyResourceScannerArgsForCall)]
	fake.notifyResourceScannerArgsForCall = append(fake.notifyResourceScannerArgsForCall, struct {
	}{})
	fake.recordInvocation("NotifyResourceScanner", []interface{}{})
	fake.notifyResourceScannerMutex.Unlock()
	if fake.NotifyResourceScannerStub != nil {
		return fake.NotifyResourceScannerStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.notifyResourceScannerReturns
	return fakeReturns.result1
}

func (fake *FakeTeamFactory) NotifyResourceScannerCallCount() int {
	fake.notifyResourceScannerMutex.RLock()
	defer fake.notifyResourceScannerMutex.RUnlock()
	return len(fake.notifyResourceScannerArgsForCall)
}

func (fake *FakeTeamFactory) NotifyResourceScannerCalls(stub func() error) {
	fake.notifyResourceScannerMutex.Lock()
	defer fake.notifyResourceScannerMutex.Unlock()
	fake.NotifyResourceScannerStub = stub
}

func (fake *FakeTeamFactory) NotifyResourceScannerReturns(result1 error) {
	fake.notifyResourceScannerMutex.Lock()
	defer fake.notifyResourceScannerMutex.Unlock()
	fake.NotifyResourceScannerStub = nil
	fake.notifyResourceScannerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTeamFactory) NotifyResourceScannerReturnsOnCall(i int, result1 error) {
	fake.notifyResourceScannerMutex.Lock()
	defer fake.notifyResourceScannerMutex.Unlock()
	fake.NotifyResourceScannerStub = nil
	if fake.notifyResourceScannerReturnsOnCall == nil {
		fake.notifyResourceScannerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.notifyResourceScannerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTeamFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createDefaultTeamIfNotExistsMutex.RLock()
	defer fake.createDefaultTeamIfNotExistsMutex.RUnlock()
	fake.createTeamMutex.RLock()
	defer fake.createTeamMutex.RUnlock()
	fake.findTeamMutex.RLock()
	defer fake.findTeamMutex.RUnlock()
	fake.getByIDMutex.RLock()
	defer fake.getByIDMutex.RUnlock()
	fake.getTeamsMutex.RLock()
	defer fake.getTeamsMutex.RUnlock()
	fake.notifyCacherMutex.RLock()
	defer fake.notifyCacherMutex.RUnlock()
	fake.notifyResourceScannerMutex.RLock()
	defer fake.notifyResourceScannerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTeamFactory) recordInvocation(key string, args []interface{}) {
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

var _ db.TeamFactory = new(FakeTeamFactory)
