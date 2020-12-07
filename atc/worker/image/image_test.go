package image_test

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/baggageclaim"
	"github.com/concourse/baggageclaim/baggageclaimfakes"
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db/dbfakes"
	"github.com/pf-qiu/concourse/v6/atc/worker"
	"github.com/pf-qiu/concourse/v6/atc/worker/image"
	"github.com/pf-qiu/concourse/v6/atc/worker/workerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image", func() {
	var (
		imageFactory     worker.ImageFactory
		img              worker.Image
		ctx              context.Context
		logger           *lagertest.TestLogger
		fakeWorker       *workerfakes.FakeWorker
		fakeVolumeClient *workerfakes.FakeVolumeClient
		fakeContainer    *dbfakes.FakeCreatingContainer
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("image-tests")
		fakeWorker = new(workerfakes.FakeWorker)
		fakeWorker.TagsReturns(atc.Tags{"worker", "tags"})

		ctx = context.Background()
		fakeVolumeClient = new(workerfakes.FakeVolumeClient)
		fakeContainer = new(dbfakes.FakeCreatingContainer)
		imageFactory = image.NewImageFactory()
	})

	Describe("imageProvidedByPreviousStepOnSameWorker", func() {
		var fakeArtifactVolume *workerfakes.FakeVolume
		var cowStrategy baggageclaim.COWStrategy

		BeforeEach(func() {
			fakeArtifactVolume = new(workerfakes.FakeVolume)
			cowStrategy = baggageclaim.COWStrategy{
				Parent: new(baggageclaimfakes.FakeVolume),
			}
			fakeArtifactVolume.COWStrategyReturns(cowStrategy)

			fakeImageArtifactSource := new(workerfakes.FakeStreamableArtifactSource)
			fakeImageArtifactSource.ExistsOnReturns(fakeArtifactVolume, true, nil)
			metadataReader := ioutil.NopCloser(strings.NewReader(
				`{"env": ["A=1", "B=2"], "user":"image-volume-user"}`,
			))
			fakeImageArtifactSource.StreamFileReturns(metadataReader, nil)

			fakeContainerRootfsVolume := new(workerfakes.FakeVolume)
			fakeContainerRootfsVolume.PathReturns("some-path")
			fakeVolumeClient.FindOrCreateCOWVolumeForContainerReturns(fakeContainerRootfsVolume, nil)

			var err error
			img, err = imageFactory.GetImage(
				logger,
				fakeWorker,
				fakeVolumeClient,
				worker.ImageSpec{
					ImageArtifactSource: fakeImageArtifactSource,
					Privileged:          true,
				},
				42,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("finds or creates cow volume", func() {
			_, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeVolumeClient.FindOrCreateCOWVolumeForContainerCallCount()).To(Equal(1))
			_, volumeSpec, container, volume, teamID, path := fakeVolumeClient.FindOrCreateCOWVolumeForContainerArgsForCall(0)
			Expect(volumeSpec).To(Equal(worker.VolumeSpec{
				Strategy:   cowStrategy,
				Privileged: true,
			}))
			Expect(container).To(Equal(fakeContainer))
			Expect(volume).To(Equal(fakeArtifactVolume))
			Expect(teamID).To(Equal(42))
			Expect(path).To(Equal("/"))
		})

		It("returns fetched image", func() {
			fetchedImage, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())

			Expect(fetchedImage).To(Equal(worker.FetchedImage{
				Metadata: worker.ImageMetadata{
					Env:  []string{"A=1", "B=2"},
					User: "image-volume-user",
				},
				URL:        "raw://some-path/rootfs",
				Privileged: true,
			}))
		})
	})

	Describe("imageProvidedByPreviousStepOnDifferentWorker", func() {
		var (
			fakeArtifactVolume        *workerfakes.FakeVolume
			fakeImageArtifactSource   *workerfakes.FakeStreamableArtifactSource
			fakeContainerRootfsVolume *workerfakes.FakeVolume
		)

		BeforeEach(func() {
			fakeArtifactVolume = new(workerfakes.FakeVolume)
			fakeImageArtifactSource = new(workerfakes.FakeStreamableArtifactSource)
			fakeImageArtifactSource.ExistsOnReturns(fakeArtifactVolume, false, nil)
			metadataReader := ioutil.NopCloser(strings.NewReader(
				`{"env": ["A=1", "B=2"], "user":"image-volume-user"}`,
			))
			fakeImageArtifactSource.StreamFileReturns(metadataReader, nil)

			fakeContainerRootfsVolume = new(workerfakes.FakeVolume)
			fakeContainerRootfsVolume.PathReturns("some-path")
			fakeVolumeClient.FindOrCreateVolumeForContainerReturns(fakeContainerRootfsVolume, nil)

			var err error
			img, err = imageFactory.GetImage(
				logger,
				fakeWorker,
				fakeVolumeClient,
				worker.ImageSpec{
					ImageArtifactSource: fakeImageArtifactSource,
					Privileged:          true,
				},
				42,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("finds or creates volume", func() {
			_, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeVolumeClient.FindOrCreateVolumeForContainerCallCount()).To(Equal(1))
			_, volumeSpec, container, teamID, path := fakeVolumeClient.FindOrCreateVolumeForContainerArgsForCall(0)
			Expect(volumeSpec).To(Equal(worker.VolumeSpec{
				Strategy:   baggageclaim.EmptyStrategy{},
				Privileged: true,
			}))
			Expect(container).To(Equal(fakeContainer))
			Expect(teamID).To(Equal(42))
			Expect(path).To(Equal("/"))
		})

		Context("when VolumeClient fails to create a volume", func() {
			BeforeEach(func() {
				fakeVolumeClient.FindOrCreateVolumeForContainerReturns(nil, errors.New("some error"))
			})

			It("returns an error", func() {
				_, err := img.FetchForContainer(ctx, logger, fakeContainer)
				Expect(err).To(HaveOccurred())
				Expect(fakeVolumeClient.FindOrCreateVolumeForContainerCallCount()).To(Equal(1))
			})
		})

		It("streams the volume from another worker", func() {
			_, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeImageArtifactSource.StreamToCallCount()).To(Equal(1))

			_, artifactDestination := fakeImageArtifactSource.StreamToArgsForCall(0)
			artifactDestination.StreamIn(context.TODO(), "fake-path", baggageclaim.GzipEncoding, strings.NewReader("fake-tar-stream"))
			Expect(fakeContainerRootfsVolume.StreamInCallCount()).To(Equal(1))
		})

		Context("when streamTo fails", func() {
			BeforeEach(func() {
				fakeImageArtifactSource.StreamFileReturns(nil, errors.New("some error"))
			})

			It("returns an error", func() {
				_, err := img.FetchForContainer(ctx, logger, fakeContainer)
				Expect(err).To(HaveOccurred())
				Expect(fakeImageArtifactSource.StreamToCallCount()).To(Equal(1))
			})
		})

		It("returns fetched image", func() {
			fetchedImage, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())

			Expect(fetchedImage).To(Equal(worker.FetchedImage{
				Metadata: worker.ImageMetadata{
					Env:  []string{"A=1", "B=2"},
					User: "image-volume-user",
				},
				URL:        "raw://some-path/rootfs",
				Privileged: true,
			}))
		})
	})

	Describe("imageFromBaseResourceType", func() {
		var cowStrategy baggageclaim.COWStrategy
		var workerResourceType atc.WorkerResourceType
		var fakeContainerRootfsVolume *workerfakes.FakeVolume
		var fakeImportVolume *workerfakes.FakeVolume

		BeforeEach(func() {
			fakeContainerRootfsVolume = new(workerfakes.FakeVolume)
			fakeContainerRootfsVolume.PathReturns("some-path")
			fakeVolumeClient.FindOrCreateCOWVolumeForContainerReturns(fakeContainerRootfsVolume, nil)

			fakeImportVolume = new(workerfakes.FakeVolume)
			cowStrategy = baggageclaim.COWStrategy{
				Parent: new(baggageclaimfakes.FakeVolume),
			}
			fakeImportVolume.COWStrategyReturns(cowStrategy)
			fakeVolumeClient.FindOrCreateVolumeForBaseResourceTypeReturns(fakeImportVolume, nil)

			workerResourceType = atc.WorkerResourceType{
				Type:       "some-base-resource-type",
				Image:      "some-base-image-path",
				Version:    "some-base-version",
				Privileged: false,
			}

			fakeWorker.ResourceTypesReturns([]atc.WorkerResourceType{workerResourceType})

			fakeWorker.NameReturns("some-worker-name")

			var err error
			img, err = imageFactory.GetImage(
				logger,
				fakeWorker,
				fakeVolumeClient,
				worker.ImageSpec{
					ResourceType: "some-base-resource-type",
				},
				42,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("finds or creates unprivileged import volume", func() {
			_, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeVolumeClient.FindOrCreateVolumeForBaseResourceTypeCallCount()).To(Equal(1))
			_, volumeSpec, teamID, resourceTypeName := fakeVolumeClient.FindOrCreateVolumeForBaseResourceTypeArgsForCall(0)
			Expect(volumeSpec).To(Equal(worker.VolumeSpec{
				Strategy: baggageclaim.ImportStrategy{
					Path: "some-base-image-path",
				},
				Privileged: false,
			}))
			Expect(teamID).To(Equal(42))
			Expect(resourceTypeName).To(Equal("some-base-resource-type"))
		})

		It("finds or creates unprivileged cow volume", func() {
			_, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeVolumeClient.FindOrCreateCOWVolumeForContainerCallCount()).To(Equal(1))
			_, volumeSpec, container, volume, teamID, path := fakeVolumeClient.FindOrCreateCOWVolumeForContainerArgsForCall(0)
			Expect(volumeSpec).To(Equal(worker.VolumeSpec{
				Strategy:   cowStrategy,
				Privileged: false,
			}))
			Expect(teamID).To(Equal(42))
			Expect(container).To(Equal(fakeContainer))
			Expect(volume).To(Equal(fakeImportVolume))
			Expect(path).To(Equal("/"))
		})

		It("returns fetched image", func() {
			fetchedImage, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())

			Expect(fetchedImage).To(Equal(worker.FetchedImage{
				Metadata:   worker.ImageMetadata{},
				URL:        "raw://some-path",
				Version:    atc.Version{"some-base-resource-type": "some-base-version"},
				Privileged: false,
			}))
		})

		Context("when the worker base resource type is privileged", func() {
			BeforeEach(func() {
				workerResourceType.Privileged = true
				fakeWorker.ResourceTypesReturns([]atc.WorkerResourceType{workerResourceType})
			})

			It("finds or creates privileged import volume", func() {
				_, err := img.FetchForContainer(ctx, logger, fakeContainer)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeVolumeClient.FindOrCreateVolumeForBaseResourceTypeCallCount()).To(Equal(1))
				_, volumeSpec, teamID, resourceTypeName := fakeVolumeClient.FindOrCreateVolumeForBaseResourceTypeArgsForCall(0)
				Expect(volumeSpec).To(Equal(worker.VolumeSpec{
					Strategy: baggageclaim.ImportStrategy{
						Path: "some-base-image-path",
					},
					Privileged: true,
				}))
				Expect(teamID).To(Equal(42))
				Expect(resourceTypeName).To(Equal("some-base-resource-type"))
			})

			It("finds or creates privileged cow volume", func() {
				_, err := img.FetchForContainer(ctx, logger, fakeContainer)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeVolumeClient.FindOrCreateCOWVolumeForContainerCallCount()).To(Equal(1))
				_, volumeSpec, container, volume, teamID, path := fakeVolumeClient.FindOrCreateCOWVolumeForContainerArgsForCall(0)
				Expect(volumeSpec).To(Equal(worker.VolumeSpec{
					Strategy:   cowStrategy,
					Privileged: true,
				}))
				Expect(teamID).To(Equal(42))
				Expect(container).To(Equal(fakeContainer))
				Expect(volume).To(Equal(fakeImportVolume))
				Expect(path).To(Equal("/"))
			})

			It("returns privileged fetched image", func() {
				fetchedImage, err := img.FetchForContainer(ctx, logger, fakeContainer)
				Expect(err).NotTo(HaveOccurred())

				Expect(fetchedImage).To(Equal(worker.FetchedImage{
					Metadata:   worker.ImageMetadata{},
					URL:        "raw://some-path",
					Version:    atc.Version{"some-base-resource-type": "some-base-version"},
					Privileged: true,
				}))
			})
		})
	})

	Describe("imageFromRootfsURI", func() {
		BeforeEach(func() {
			var err error
			img, err = imageFactory.GetImage(
				logger,
				fakeWorker,
				fakeVolumeClient,
				worker.ImageSpec{
					ImageURL: "some-image-url",
				},
				42,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the fetched image", func() {
			fetchedImage, err := img.FetchForContainer(ctx, logger, fakeContainer)
			Expect(err).NotTo(HaveOccurred())

			Expect(fetchedImage).To(Equal(worker.FetchedImage{
				URL: "some-image-url",
			}))
		})
	})
})
