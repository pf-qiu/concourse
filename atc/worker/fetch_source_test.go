package worker_test

import (
	"context"
	"errors"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/garden/gardenfakes"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/db/dbfakes"
	"github.com/pf-qiu/concourse/v6/atc/resource"
	"github.com/pf-qiu/concourse/v6/atc/resource/resourcefakes"
	"github.com/pf-qiu/concourse/v6/atc/runtime"
	"github.com/pf-qiu/concourse/v6/atc/worker"
	"github.com/pf-qiu/concourse/v6/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchSource", func() {
	var (
		fetchSourceFactory worker.FetchSourceFactory
		fetchSource        worker.FetchSource

		fakeContainer            *workerfakes.FakeContainer
		fakeVolume               *workerfakes.FakeVolume
		fakeWorker               *workerfakes.FakeWorker
		fakeResourceCacheFactory *dbfakes.FakeResourceCacheFactory
		fakeUsedResourceCache    *dbfakes.FakeUsedResourceCache
		fakeResource             *resourcefakes.FakeResource
		metadata                 db.ContainerMetadata
		owner                    db.ContainerOwner

		ctx    context.Context
		cancel func()
	)

	BeforeEach(func() {
		logger := lagertest.NewTestLogger("test")
		fakeContainer = new(workerfakes.FakeContainer)

		ctx, cancel = context.WithCancel(context.Background())

		fakeContainer.PropertyReturns("", errors.New("nope"))
		inProcess := new(gardenfakes.FakeProcess)
		inProcess.IDReturns("process-id")
		inProcess.WaitStub = func() (int, error) {
			return 0, nil
		}

		fakeContainer.AttachReturns(nil, errors.New("process not found"))

		fakeContainer.RunStub = func(ctx context.Context, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
			_, err := io.Stdout.Write([]byte("{}"))
			Expect(err).NotTo(HaveOccurred())

			return inProcess, nil
		}

		fakeVolume = new(workerfakes.FakeVolume)
		fakeVolume.HandleReturns("some-handle")
		fakeContainer.VolumeMountsReturns([]worker.VolumeMount{
			{
				Volume:    fakeVolume,
				MountPath: resource.ResourcesDir("get"),
			},
		})

		fakeWorker = new(workerfakes.FakeWorker)
		fakeWorker.FindOrCreateContainerReturns(fakeContainer, nil)

		fakeResourceCacheFactory = new(dbfakes.FakeResourceCacheFactory)
		fakeUsedResourceCache = new(dbfakes.FakeUsedResourceCache)
		fakeUsedResourceCache.IDReturns(42)
		fakeResourceCacheFactory.FindOrCreateResourceCacheReturns(fakeUsedResourceCache, nil)

		owner = db.NewBuildStepContainerOwner(43, atc.PlanID("some-plan-id"), 42)
		metadata = db.ContainerMetadata{Type: db.ContainerTypeGet}

		fakeResource = new(resourcefakes.FakeResource)
		fakeResourceCacheFactory = new(dbfakes.FakeResourceCacheFactory)
		fakeResourceCacheFactory.UpdateResourceCacheMetadataReturns(nil)
		fakeResourceCacheFactory.ResourceCacheMetadataReturns([]db.ResourceConfigMetadataField{
			{Name: "some", Value: "metadata"},
		}, nil)

		getProcessSpec := runtime.ProcessSpec{
			Path: "/opt/resource/in",
			Args: []string{resource.ResourcesDir("get")},
		}

		fetchSourceFactory = worker.NewFetchSourceFactory(fakeResourceCacheFactory)
		fetchSource = fetchSourceFactory.NewFetchSource(
			logger,
			fakeWorker,
			owner,
			fakeUsedResourceCache,
			fakeResource,
			worker.ContainerSpec{
				TeamID: 42,
				ImageSpec: worker.ImageSpec{
					ResourceType: "fake-resource-type",
				},
				Outputs: map[string]string{
					"resource": resource.ResourcesDir("get"),
				},
			},
			getProcessSpec,
			metadata,
		)
	})

	AfterEach(func() {
		cancel()
	})

	Describe("Find", func() {
		var expectedGetResult worker.GetResult

		Context("when there is volume", func() {
			BeforeEach(func() {
				fakeWorker.FindVolumeForResourceCacheReturns(fakeVolume, true, nil)

				expectedMetadata := []atc.MetadataField{
					{Name: "some", Value: "metadata"},
				}
				expectedGetResult = worker.GetResult{
					ExitStatus:    0,
					VersionResult: runtime.VersionResult{Metadata: expectedMetadata},
					GetArtifact:   runtime.GetArtifact{VolumeHandle: fakeVolume.Handle()},
				}
			})

			It("finds the resource cache volume and returns the correct result", func() {
				getResult, volume, found, err := fetchSource.Find()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				Expect(volume).To(Equal(fakeVolume))
				Expect(getResult).To(Equal(expectedGetResult))
			})
		})

		Context("when there is no volume", func() {
			BeforeEach(func() {
				fakeWorker.FindVolumeForResourceCacheReturns(nil, false, nil)
				expectedGetResult = worker.GetResult{}
			})

			It("does not find volume", func() {
				getResult, volume, found, err := fetchSource.Find()
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(volume).To(BeNil())
				Expect(getResult).To(Equal(expectedGetResult))
			})
		})
	})

	Describe("Create", func() {
		var (
			err               error
			getResult         worker.GetResult
			expectedGetResult worker.GetResult
			volume            worker.Volume
		)

		JustBeforeEach(func() {
			getResult, volume, err = fetchSource.Create(ctx)
		})

		Context("when there is initialized volume", func() {
			BeforeEach(func() {
				fakeWorker.FindVolumeForResourceCacheReturns(fakeVolume, true, nil)

				expectedMetadata := []atc.MetadataField{
					{Name: "some", Value: "metadata"},
				}
				expectedGetResult = worker.GetResult{
					ExitStatus:    0,
					VersionResult: runtime.VersionResult{Metadata: expectedMetadata},
					GetArtifact:   runtime.GetArtifact{VolumeHandle: fakeVolume.Handle()},
				}
			})

			It("does not fetch resource", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeWorker.FindOrCreateContainerCallCount()).To(Equal(0))
				Expect(fakeResource.GetCallCount()).To(Equal(0))
			})

			It("finds initialized volume and returns a successful GetResult", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(getResult).To(Equal(expectedGetResult))
				Expect(volume).ToNot(BeNil())
			})
		})

		Context("when there is no initialized volume", func() {
			var atcMetadata []atc.MetadataField
			BeforeEach(func() {
				atcMetadata = []atc.MetadataField{{Name: "foo", Value: "bar"}}
				fakeWorker.FindVolumeForResourceCacheReturns(nil, false, nil)
				fakeResource.GetReturns(runtime.VersionResult{Metadata: atcMetadata}, nil)
			})

			It("creates container with volume and worker", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeWorker.FindOrCreateContainerCallCount()).To(Equal(1))
				_, logger, owner, actualMetadata, containerSpec := fakeWorker.FindOrCreateContainerArgsForCall(0)
				Expect(owner).To(Equal(db.NewBuildStepContainerOwner(43, atc.PlanID("some-plan-id"), 42)))
				Expect(actualMetadata).To(Equal(metadata))
				Expect(containerSpec).To(Equal(
					worker.ContainerSpec{
						TeamID: 42,
						ImageSpec: worker.ImageSpec{
							ResourceType: "fake-resource-type",
						},
						BindMounts: []worker.BindMountSource{&worker.CertsVolumeMount{Logger: logger}},
						Outputs: map[string]string{
							"resource": resource.ResourcesDir("get"),
						},
					}))
			})

			It("executes the get script", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeResource.GetCallCount()).To(Equal(1))
			})

			It("initializes cache", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeVolume.InitializeResourceCacheCallCount()).To(Equal(1))
				rc := fakeVolume.InitializeResourceCacheArgsForCall(0)
				Expect(rc).To(Equal(fakeUsedResourceCache))
			})

			It("updates resource cache metadata", func() {
				Expect(fakeResourceCacheFactory.UpdateResourceCacheMetadataCallCount()).To(Equal(1))
				passedResourceCache, versionResultMetadata := fakeResourceCacheFactory.UpdateResourceCacheMetadataArgsForCall(0)
				Expect(passedResourceCache).To(Equal(fakeUsedResourceCache))
				Expect(versionResultMetadata).To(Equal(atcMetadata))
			})

			Context("when getting resource fails with other error", func() {
				var disaster error

				BeforeEach(func() {
					disaster = errors.New("failed")
					fakeResource.GetReturns(runtime.VersionResult{}, disaster)
				})

				It("returns the error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(disaster))
				})
			})

			It("returns a successful GetResult and volume with fetched bits", func() {
				Expect(getResult.ExitStatus).To(BeZero())
				Expect(getResult.GetArtifact.VolumeHandle).To(Equal(fakeVolume.Handle()))
				Expect(volume).ToNot(BeNil())
			})
		})
	})
})
