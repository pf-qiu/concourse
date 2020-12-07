package gc_test

import (
	"context"

	"github.com/pf-qiu/concourse/v6/atc/db/dbfakes"
	"github.com/pf-qiu/concourse/v6/atc/gc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PipelineCollector", func() {
	var collector GcCollector
	var fakePipelineLifecycle *dbfakes.FakePipelineLifecycle

	BeforeEach(func() {
		fakePipelineLifecycle = new(dbfakes.FakePipelineLifecycle)

		collector = gc.NewPipelineCollector(fakePipelineLifecycle)
	})

	Describe("Run", func() {
		It("tells the pipeline lifecycle to remove abandoned pipelines", func() {
			err := collector.Run(context.TODO())
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePipelineLifecycle.ArchiveAbandonedPipelinesCallCount()).To(Equal(1))
		})
	})
})
