package exec_test

import (
	"context"
	"errors"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db/dbfakes"
	"github.com/pf-qiu/concourse/v6/atc/exec"
	"github.com/pf-qiu/concourse/v6/atc/exec/build"
	"github.com/pf-qiu/concourse/v6/atc/worker/workerfakes"
	"github.com/pf-qiu/concourse/v6/vars"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArtifactInputStep", func() {
	var (
		ctx    context.Context
		cancel func()

		state exec.RunState

		step             exec.Step
		stepOk           bool
		stepErr          error
		plan             atc.Plan
		fakeBuild        *dbfakes.FakeBuild
		fakeWorkerClient *workerfakes.FakeClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		state = exec.NewRunState(noopStepper, vars.StaticVariables{}, false)

		fakeBuild = new(dbfakes.FakeBuild)
		fakeWorkerClient = new(workerfakes.FakeClient)

		plan = atc.Plan{ArtifactInput: &atc.ArtifactInputPlan{34, "some-input-artifact-name"}}
		step = exec.NewArtifactInputStep(plan, fakeBuild, fakeWorkerClient)
	})

	AfterEach(func() {
		cancel()
	})

	JustBeforeEach(func() {
		stepOk, stepErr = step.Run(ctx, state)
	})

	Context("when looking up the build artifact errors", func() {
		BeforeEach(func() {
			fakeBuild.ArtifactReturns(nil, errors.New("nope"))
		})
		It("returns the error", func() {
			Expect(stepErr).To(HaveOccurred())
		})
	})

	Context("when looking up the build artifact succeeds", func() {
		var fakeWorkerArtifact *dbfakes.FakeWorkerArtifact

		BeforeEach(func() {
			fakeWorkerArtifact = new(dbfakes.FakeWorkerArtifact)
			fakeBuild.ArtifactReturns(fakeWorkerArtifact, nil)
		})

		Context("when looking up the db volume fails", func() {
			BeforeEach(func() {
				fakeWorkerArtifact.VolumeReturns(nil, false, errors.New("nope"))
			})
			It("returns the error", func() {
				Expect(stepErr).To(HaveOccurred())
			})
		})

		Context("when the db volume does not exist", func() {
			BeforeEach(func() {
				fakeWorkerArtifact.VolumeReturns(nil, false, nil)
			})
			It("returns the error", func() {
				Expect(stepErr).To(HaveOccurred())
			})
		})

		Context("when the db volume does exist", func() {
			var fakeVolume *dbfakes.FakeCreatedVolume

			BeforeEach(func() {
				fakeVolume = new(dbfakes.FakeCreatedVolume)
				fakeWorkerArtifact.VolumeReturns(fakeVolume, true, nil)
			})

			Context("when looking up the worker volume fails", func() {
				BeforeEach(func() {
					fakeWorkerClient.FindVolumeReturns(nil, false, errors.New("nope"))
				})
				It("returns the error", func() {
					Expect(stepErr).To(HaveOccurred())
				})
			})

			Context("when the worker volume does not exist", func() {
				BeforeEach(func() {
					fakeWorkerClient.FindVolumeReturns(nil, false, nil)
				})
				It("returns the error", func() {
					Expect(stepErr).To(HaveOccurred())
				})
			})

			Context("when the volume does exist", func() {
				var fakeWorkerVolume *workerfakes.FakeVolume
				var fakeDBWorkerArtifact *dbfakes.FakeWorkerArtifact
				var fakeDBCreatedVolume *dbfakes.FakeCreatedVolume

				BeforeEach(func() {
					fakeWorkerVolume = new(workerfakes.FakeVolume)
					fakeWorkerClient.FindVolumeReturns(fakeWorkerVolume, true, nil)

					fakeDBWorkerArtifact = new(dbfakes.FakeWorkerArtifact)
					fakeDBCreatedVolume = new(dbfakes.FakeCreatedVolume)
					fakeDBCreatedVolume.HandleReturns("some-volume-handle")
					fakeDBWorkerArtifact.VolumeReturns(fakeDBCreatedVolume, true, nil)
					fakeBuild.ArtifactReturns(fakeDBWorkerArtifact, nil)
				})

				It("registers the artifact", func() {
					artifact, found := state.ArtifactRepository().ArtifactFor(build.ArtifactName("some-input-artifact-name"))

					Expect(stepErr).NotTo(HaveOccurred())
					Expect(found).To(BeTrue())

					Expect(artifact.ID()).To(Equal("some-volume-handle"))
				})

				It("succeeds", func() {
					Expect(stepOk).To(BeTrue())
				})
			})
		})
	})
})
