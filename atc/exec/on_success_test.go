package exec_test

import (
	"context"
	"errors"

	"github.com/pf-qiu/concourse/v6/atc/exec"
	"github.com/pf-qiu/concourse/v6/atc/exec/build"
	"github.com/pf-qiu/concourse/v6/atc/exec/execfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("On Success Step", func() {
	var (
		ctx    context.Context
		cancel func()

		step *execfakes.FakeStep
		hook *execfakes.FakeStep

		repo  *build.Repository
		state *execfakes.FakeRunState

		onSuccessStep exec.Step

		stepOk  bool
		stepErr error
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		step = &execfakes.FakeStep{}
		hook = &execfakes.FakeStep{}

		repo = build.NewRepository()
		state = new(execfakes.FakeRunState)
		state.ArtifactRepositoryReturns(repo)

		onSuccessStep = exec.OnSuccess(step, hook)

		stepOk = false
		stepErr = nil
	})

	JustBeforeEach(func() {
		stepOk, stepErr = onSuccessStep.Run(ctx, state)
	})

	AfterEach(func() {
		cancel()
	})

	Context("when the step succeeds", func() {
		BeforeEach(func() {
			step.RunReturns(true, nil)
		})

		It("runs the hook", func() {
			Expect(step.RunCallCount()).To(Equal(1))
			Expect(hook.RunCallCount()).To(Equal(1))
		})

		It("runs the hook with the run state", func() {
			Expect(hook.RunCallCount()).To(Equal(1))

			_, argsState := hook.RunArgsForCall(0)
			Expect(argsState).To(Equal(state))
		})

		It("propagates the context to the hook", func() {
			runCtx, _ := hook.RunArgsForCall(0)
			Expect(runCtx).To(Equal(ctx))
		})

		It("returns nil", func() {
			Expect(stepErr).ToNot(HaveOccurred())
		})
	})

	Context("when the step errors", func() {
		disaster := errors.New("disaster")

		BeforeEach(func() {
			step.RunReturns(false, disaster)
		})

		It("does not run the hook", func() {
			Expect(step.RunCallCount()).To(Equal(1))
			Expect(hook.RunCallCount()).To(Equal(0))
		})

		It("returns the error", func() {
			Expect(stepErr).To(Equal(disaster))
		})
	})

	Context("when the step fails", func() {
		BeforeEach(func() {
			step.RunReturns(false, nil)
		})

		It("does not run the hook", func() {
			Expect(step.RunCallCount()).To(Equal(1))
			Expect(hook.RunCallCount()).To(Equal(0))
		})

		It("returns nil", func() {
			Expect(stepErr).To(BeNil())
		})
	})

	It("propagates the context to the step", func() {
		runCtx, _ := step.RunArgsForCall(0)
		Expect(runCtx).To(Equal(ctx))
	})

	Context("when step fails", func() {
		BeforeEach(func() {
			step.RunReturns(false, nil)
		})

		It("returns false", func() {
			Expect(stepOk).To(BeFalse())
		})
	})

	Context("when step succeeds and hook succeeds", func() {
		BeforeEach(func() {
			step.RunReturns(true, nil)
			hook.RunReturns(true, nil)
		})

		It("returns true", func() {
			Expect(stepOk).To(BeTrue())
		})
	})

	Context("when step succeeds and hook fails", func() {
		BeforeEach(func() {
			step.RunReturns(false, nil)
			hook.RunReturns(false, nil)
		})

		It("returns false", func() {
			Expect(stepOk).To(BeFalse())
		})
	})
})
