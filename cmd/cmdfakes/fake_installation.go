// Code generated by counterfeiter. DO NOT EDIT.
package cmdfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/v7/installation"
)

type FakeInstallation struct {
	JobsStub        func() []installation.InstalledJob
	jobsMutex       sync.RWMutex
	jobsArgsForCall []struct {
	}
	jobsReturns struct {
		result1 []installation.InstalledJob
	}
	jobsReturnsOnCall map[int]struct {
		result1 []installation.InstalledJob
	}
	TargetStub        func() installation.Target
	targetMutex       sync.RWMutex
	targetArgsForCall []struct {
	}
	targetReturns struct {
		result1 installation.Target
	}
	targetReturnsOnCall map[int]struct {
		result1 installation.Target
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInstallation) Jobs() []installation.InstalledJob {
	fake.jobsMutex.Lock()
	ret, specificReturn := fake.jobsReturnsOnCall[len(fake.jobsArgsForCall)]
	fake.jobsArgsForCall = append(fake.jobsArgsForCall, struct {
	}{})
	stub := fake.JobsStub
	fakeReturns := fake.jobsReturns
	fake.recordInvocation("Jobs", []interface{}{})
	fake.jobsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeInstallation) JobsCallCount() int {
	fake.jobsMutex.RLock()
	defer fake.jobsMutex.RUnlock()
	return len(fake.jobsArgsForCall)
}

func (fake *FakeInstallation) JobsCalls(stub func() []installation.InstalledJob) {
	fake.jobsMutex.Lock()
	defer fake.jobsMutex.Unlock()
	fake.JobsStub = stub
}

func (fake *FakeInstallation) JobsReturns(result1 []installation.InstalledJob) {
	fake.jobsMutex.Lock()
	defer fake.jobsMutex.Unlock()
	fake.JobsStub = nil
	fake.jobsReturns = struct {
		result1 []installation.InstalledJob
	}{result1}
}

func (fake *FakeInstallation) JobsReturnsOnCall(i int, result1 []installation.InstalledJob) {
	fake.jobsMutex.Lock()
	defer fake.jobsMutex.Unlock()
	fake.JobsStub = nil
	if fake.jobsReturnsOnCall == nil {
		fake.jobsReturnsOnCall = make(map[int]struct {
			result1 []installation.InstalledJob
		})
	}
	fake.jobsReturnsOnCall[i] = struct {
		result1 []installation.InstalledJob
	}{result1}
}

func (fake *FakeInstallation) Target() installation.Target {
	fake.targetMutex.Lock()
	ret, specificReturn := fake.targetReturnsOnCall[len(fake.targetArgsForCall)]
	fake.targetArgsForCall = append(fake.targetArgsForCall, struct {
	}{})
	stub := fake.TargetStub
	fakeReturns := fake.targetReturns
	fake.recordInvocation("Target", []interface{}{})
	fake.targetMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeInstallation) TargetCallCount() int {
	fake.targetMutex.RLock()
	defer fake.targetMutex.RUnlock()
	return len(fake.targetArgsForCall)
}

func (fake *FakeInstallation) TargetCalls(stub func() installation.Target) {
	fake.targetMutex.Lock()
	defer fake.targetMutex.Unlock()
	fake.TargetStub = stub
}

func (fake *FakeInstallation) TargetReturns(result1 installation.Target) {
	fake.targetMutex.Lock()
	defer fake.targetMutex.Unlock()
	fake.TargetStub = nil
	fake.targetReturns = struct {
		result1 installation.Target
	}{result1}
}

func (fake *FakeInstallation) TargetReturnsOnCall(i int, result1 installation.Target) {
	fake.targetMutex.Lock()
	defer fake.targetMutex.Unlock()
	fake.TargetStub = nil
	if fake.targetReturnsOnCall == nil {
		fake.targetReturnsOnCall = make(map[int]struct {
			result1 installation.Target
		})
	}
	fake.targetReturnsOnCall[i] = struct {
		result1 installation.Target
	}{result1}
}

func (fake *FakeInstallation) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.jobsMutex.RLock()
	defer fake.jobsMutex.RUnlock()
	fake.targetMutex.RLock()
	defer fake.targetMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInstallation) recordInvocation(key string, args []interface{}) {
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

var _ installation.Installation = new(FakeInstallation)
