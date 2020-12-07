package resource_test

import (
	"testing"

	"github.com/pf-qiu/concourse/v6/atc/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	resourceFactory resource.ResourceFactory
)

var _ = BeforeEach(func() {

	resourceFactory = resource.NewResourceFactory()
})

func TestResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resource Suite")
}
