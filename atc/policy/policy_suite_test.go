package policy_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/pf-qiu/concourse/v6/atc/policy"
	"github.com/pf-qiu/concourse/v6/atc/policy/policyfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPolicy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Policy Suite")
}

var (
	testLogger = lagertest.NewTestLogger("test")

	fakeAgent        *policyfakes.FakeAgent
	fakeAgentFactory *policyfakes.FakeAgentFactory
)

var _ = BeforeSuite(func() {
	fakeAgentFactory = new(policyfakes.FakeAgentFactory)
	fakeAgentFactory.IsConfiguredReturns(true)
	fakeAgentFactory.DescriptionReturns("fakeAgent")
	policy.RegisterAgent(fakeAgentFactory)
})
