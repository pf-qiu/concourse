package exec_test

import (
	"context"

	. "github.com/pf-qiu/concourse/v6/atc/exec"
	"github.com/pf-qiu/concourse/v6/atc/exec/execfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Identity", func() {
	var (
		state *execfakes.FakeRunState

		step IdentityStep

		stepOk  bool
		stepErr error
	)

	BeforeEach(func() {
		state = new(execfakes.FakeRunState)
	})

	JustBeforeEach(func() {
		stepOk, stepErr = step.Run(context.Background(), state)
	})

	Describe("Run", func() {
		It("is a no-op", func() {
			Expect(stepErr).ToNot(HaveOccurred())
			Expect(stepOk).To(BeTrue())
		})
	})
})
