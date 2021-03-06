// Code generated by counterfeiter. DO NOT EDIT.
package workerfakes

import (
	"context"
	"io"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/worker"
)

type FakeStreamableArtifactSource struct {
	ExistsOnStub        func(lager.Logger, worker.Worker) (worker.Volume, bool, error)
	existsOnMutex       sync.RWMutex
	existsOnArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.Worker
	}
	existsOnReturns struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	existsOnReturnsOnCall map[int]struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}
	StreamFileStub        func(context.Context, string) (io.ReadCloser, error)
	streamFileMutex       sync.RWMutex
	streamFileArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	streamFileReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	streamFileReturnsOnCall map[int]struct {
		result1 io.ReadCloser
		result2 error
	}
	StreamToStub        func(context.Context, worker.ArtifactDestination) error
	streamToMutex       sync.RWMutex
	streamToArgsForCall []struct {
		arg1 context.Context
		arg2 worker.ArtifactDestination
	}
	streamToReturns struct {
		result1 error
	}
	streamToReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStreamableArtifactSource) ExistsOn(arg1 lager.Logger, arg2 worker.Worker) (worker.Volume, bool, error) {
	fake.existsOnMutex.Lock()
	ret, specificReturn := fake.existsOnReturnsOnCall[len(fake.existsOnArgsForCall)]
	fake.existsOnArgsForCall = append(fake.existsOnArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.Worker
	}{arg1, arg2})
	fake.recordInvocation("ExistsOn", []interface{}{arg1, arg2})
	fake.existsOnMutex.Unlock()
	if fake.ExistsOnStub != nil {
		return fake.ExistsOnStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.existsOnReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeStreamableArtifactSource) ExistsOnCallCount() int {
	fake.existsOnMutex.RLock()
	defer fake.existsOnMutex.RUnlock()
	return len(fake.existsOnArgsForCall)
}

func (fake *FakeStreamableArtifactSource) ExistsOnCalls(stub func(lager.Logger, worker.Worker) (worker.Volume, bool, error)) {
	fake.existsOnMutex.Lock()
	defer fake.existsOnMutex.Unlock()
	fake.ExistsOnStub = stub
}

func (fake *FakeStreamableArtifactSource) ExistsOnArgsForCall(i int) (lager.Logger, worker.Worker) {
	fake.existsOnMutex.RLock()
	defer fake.existsOnMutex.RUnlock()
	argsForCall := fake.existsOnArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStreamableArtifactSource) ExistsOnReturns(result1 worker.Volume, result2 bool, result3 error) {
	fake.existsOnMutex.Lock()
	defer fake.existsOnMutex.Unlock()
	fake.ExistsOnStub = nil
	fake.existsOnReturns = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeStreamableArtifactSource) ExistsOnReturnsOnCall(i int, result1 worker.Volume, result2 bool, result3 error) {
	fake.existsOnMutex.Lock()
	defer fake.existsOnMutex.Unlock()
	fake.ExistsOnStub = nil
	if fake.existsOnReturnsOnCall == nil {
		fake.existsOnReturnsOnCall = make(map[int]struct {
			result1 worker.Volume
			result2 bool
			result3 error
		})
	}
	fake.existsOnReturnsOnCall[i] = struct {
		result1 worker.Volume
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeStreamableArtifactSource) StreamFile(arg1 context.Context, arg2 string) (io.ReadCloser, error) {
	fake.streamFileMutex.Lock()
	ret, specificReturn := fake.streamFileReturnsOnCall[len(fake.streamFileArgsForCall)]
	fake.streamFileArgsForCall = append(fake.streamFileArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("StreamFile", []interface{}{arg1, arg2})
	fake.streamFileMutex.Unlock()
	if fake.StreamFileStub != nil {
		return fake.StreamFileStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.streamFileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStreamableArtifactSource) StreamFileCallCount() int {
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	return len(fake.streamFileArgsForCall)
}

func (fake *FakeStreamableArtifactSource) StreamFileCalls(stub func(context.Context, string) (io.ReadCloser, error)) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = stub
}

func (fake *FakeStreamableArtifactSource) StreamFileArgsForCall(i int) (context.Context, string) {
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	argsForCall := fake.streamFileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStreamableArtifactSource) StreamFileReturns(result1 io.ReadCloser, result2 error) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = nil
	fake.streamFileReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeStreamableArtifactSource) StreamFileReturnsOnCall(i int, result1 io.ReadCloser, result2 error) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = nil
	if fake.streamFileReturnsOnCall == nil {
		fake.streamFileReturnsOnCall = make(map[int]struct {
			result1 io.ReadCloser
			result2 error
		})
	}
	fake.streamFileReturnsOnCall[i] = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeStreamableArtifactSource) StreamTo(arg1 context.Context, arg2 worker.ArtifactDestination) error {
	fake.streamToMutex.Lock()
	ret, specificReturn := fake.streamToReturnsOnCall[len(fake.streamToArgsForCall)]
	fake.streamToArgsForCall = append(fake.streamToArgsForCall, struct {
		arg1 context.Context
		arg2 worker.ArtifactDestination
	}{arg1, arg2})
	fake.recordInvocation("StreamTo", []interface{}{arg1, arg2})
	fake.streamToMutex.Unlock()
	if fake.StreamToStub != nil {
		return fake.StreamToStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.streamToReturns
	return fakeReturns.result1
}

func (fake *FakeStreamableArtifactSource) StreamToCallCount() int {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	return len(fake.streamToArgsForCall)
}

func (fake *FakeStreamableArtifactSource) StreamToCalls(stub func(context.Context, worker.ArtifactDestination) error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = stub
}

func (fake *FakeStreamableArtifactSource) StreamToArgsForCall(i int) (context.Context, worker.ArtifactDestination) {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	argsForCall := fake.streamToArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeStreamableArtifactSource) StreamToReturns(result1 error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = nil
	fake.streamToReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStreamableArtifactSource) StreamToReturnsOnCall(i int, result1 error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = nil
	if fake.streamToReturnsOnCall == nil {
		fake.streamToReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.streamToReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStreamableArtifactSource) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.existsOnMutex.RLock()
	defer fake.existsOnMutex.RUnlock()
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStreamableArtifactSource) recordInvocation(key string, args []interface{}) {
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

var _ worker.StreamableArtifactSource = new(FakeStreamableArtifactSource)
