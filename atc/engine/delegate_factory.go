package engine

import (
	"code.cloudfoundry.org/clock"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/exec"
	"github.com/pf-qiu/concourse/v6/atc/policy"
)

type DelegateFactory struct {
	build         db.Build
	plan          atc.Plan
	rateLimiter   RateLimiter
	policyChecker policy.Checker
}

func buildDelegateFactory(
	build db.Build,
	plan atc.Plan,
	rateLimiter RateLimiter,
	policyChecker policy.Checker,
) DelegateFactory {
	return DelegateFactory{
		build:         build,
		plan:          plan,
		rateLimiter:   rateLimiter,
		policyChecker: policyChecker,
	}
}

func (delegate DelegateFactory) GetDelegate(state exec.RunState) exec.GetDelegate {
	return NewGetDelegate(delegate.build, delegate.plan.ID, state, clock.NewClock(), delegate.policyChecker)
}

func (delegate DelegateFactory) PutDelegate(state exec.RunState) exec.PutDelegate {
	return NewPutDelegate(delegate.build, delegate.plan.ID, state, clock.NewClock(), delegate.policyChecker)
}

func (delegate DelegateFactory) TaskDelegate(state exec.RunState) exec.TaskDelegate {
	return NewTaskDelegate(delegate.build, delegate.plan.ID, state, clock.NewClock(), delegate.policyChecker)
}

func (delegate DelegateFactory) CheckDelegate(state exec.RunState) exec.CheckDelegate {
	return NewCheckDelegate(delegate.build, delegate.plan, state, clock.NewClock(), delegate.rateLimiter, delegate.policyChecker)
}

func (delegate DelegateFactory) BuildStepDelegate(state exec.RunState) exec.BuildStepDelegate {
	return NewBuildStepDelegate(delegate.build, delegate.plan.ID, state, clock.NewClock(), delegate.policyChecker)
}

func (delegate DelegateFactory) SetPipelineStepDelegate(state exec.RunState) exec.SetPipelineStepDelegate {
	return NewSetPipelineStepDelegate(delegate.build, delegate.plan.ID, state, clock.NewClock())
}
