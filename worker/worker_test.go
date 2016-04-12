package worker_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/garden"
	gfakes "github.com/cloudfoundry-incubator/garden/fakes"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	. "github.com/concourse/atc/worker"
	wfakes "github.com/concourse/atc/worker/fakes"
	"github.com/concourse/baggageclaim"
	bfakes "github.com/concourse/baggageclaim/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/clock/fakeclock"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Worker", func() {
	var (
		logger                 *lagertest.TestLogger
		fakeGardenClient       *gfakes.FakeClient
		fakeBaggageclaimClient *bfakes.FakeClient
		fakeVolumeFactory      *wfakes.FakeVolumeFactory
		fakeImageFetcher       *wfakes.FakeImageFetcher
		fakeGardenWorkerDB     *wfakes.FakeGardenWorkerDB
		fakeWorkerProvider     *wfakes.FakeWorkerProvider
		fakeClock              *fakeclock.FakeClock
		activeContainers       int
		resourceTypes          []atc.WorkerResourceType
		platform               string
		tags                   atc.Tags
		workerName             string
		httpProxyURL           string
		httpsProxyURL          string
		noProxy                string

		gardenWorker Worker
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeGardenClient = new(gfakes.FakeClient)
		fakeBaggageclaimClient = new(bfakes.FakeClient)
		fakeVolumeFactory = new(wfakes.FakeVolumeFactory)
		fakeImageFetcher = new(wfakes.FakeImageFetcher)
		fakeGardenWorkerDB = new(wfakes.FakeGardenWorkerDB)
		fakeWorkerProvider = new(wfakes.FakeWorkerProvider)
		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))
		activeContainers = 42
		resourceTypes = []atc.WorkerResourceType{
			{Type: "some-resource", Image: "some-resource-image"},
		}
		platform = "some-platform"
		tags = atc.Tags{"some", "tags"}
		workerName = "some-worker"

		gardenWorker = NewGardenWorker(
			fakeGardenClient,
			fakeBaggageclaimClient,
			fakeVolumeFactory,
			fakeImageFetcher,
			fakeGardenWorkerDB,
			fakeWorkerProvider,
			fakeClock,
			activeContainers,
			resourceTypes,
			platform,
			tags,
			workerName,
			httpProxyURL,
			httpsProxyURL,
			noProxy,
		)
	})

	Describe("CreateVolume", func() {
		var baggageclaimClient baggageclaim.Client

		var volumeSpec VolumeSpec

		var createdVolume Volume
		var createErr error

		BeforeEach(func() {
			volumeSpec = VolumeSpec{
				Properties: VolumeProperties{
					"some": "property",
				},
				Privileged: true,
				TTL:        6 * time.Minute,
			}
		})

		JustBeforeEach(func() {
			createdVolume, createErr = NewGardenWorker(
				fakeGardenClient,
				baggageclaimClient,
				fakeVolumeFactory,
				fakeImageFetcher,
				fakeGardenWorkerDB,
				fakeWorkerProvider,
				fakeClock,
				activeContainers,
				resourceTypes,
				platform,
				tags,
				workerName,
				httpProxyURL,
				httpsProxyURL,
				noProxy,
			).CreateVolume(logger, volumeSpec)
		})

		Context("when there is no baggageclaim client", func() {
			BeforeEach(func() {
				baggageclaimClient = nil
			})

			It("returns ErrNoVolumeManager", func() {
				Expect(createErr).To(Equal(ErrNoVolumeManager))
			})
		})

		Context("when there is a baggageclaim client", func() {
			var fakeBaggageclaimVolume *bfakes.FakeVolume
			var builtVolume *wfakes.FakeVolume

			BeforeEach(func() {
				baggageclaimClient = fakeBaggageclaimClient

				fakeBaggageclaimVolume = new(bfakes.FakeVolume)
				fakeBaggageclaimVolume.HandleReturns("created-volume")

				fakeBaggageclaimClient.CreateVolumeReturns(fakeBaggageclaimVolume, nil)

				builtVolume = new(wfakes.FakeVolume)
				fakeVolumeFactory.BuildReturns(builtVolume, true, nil)
			})

			Context("when creating a ResourceCacheStrategy volume", func() {
				BeforeEach(func() {
					volumeSpec.Strategy = ResourceCacheStrategy{
						ResourceHash:    "some-resource-hash",
						ResourceVersion: atc.Version{"some": "resource-version"},
					}
				})

				It("succeeds", func() {
					Expect(createErr).ToNot(HaveOccurred())
				})

				It("creates the volume via BaggageClaim", func() {
					Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(1))

					_, spec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
					Expect(spec).To(Equal(baggageclaim.VolumeSpec{
						Strategy:   baggageclaim.EmptyStrategy{},
						Properties: baggageclaim.VolumeProperties(volumeSpec.Properties),
						TTL:        volumeSpec.TTL,
						Privileged: volumeSpec.Privileged,
					}))
				})

				It("inserts the volume into the database", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(1))

					dbVolume := fakeGardenWorkerDB.InsertVolumeArgsForCall(0)
					Expect(dbVolume).To(Equal(db.Volume{
						Handle:     "created-volume",
						WorkerName: workerName,
						TTL:        volumeSpec.TTL,
						Identifier: db.VolumeIdentifier{
							ResourceCache: &db.ResourceCacheIdentifier{
								ResourceHash:    "some-resource-hash",
								ResourceVersion: atc.Version{"some": "resource-version"},
							},
						},
					}))
				})

				It("builds the baggageclaim.Volume into a worker.Volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(1))

					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume))

					Expect(createdVolume).To(Equal(builtVolume))
				})

				Context("when creating the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeBaggageclaimClient.CreateVolumeReturns(nil, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when inserting the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeGardenWorkerDB.InsertVolumeReturns(disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume cannot find the volume in the database", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, nil)
					})

					It("returns ErrMissingVolume", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})
			})

			Context("when creating a OutputStrategy volume", func() {
				BeforeEach(func() {
					volumeSpec.Strategy = OutputStrategy{
						Name: "some-output",
					}
				})

				It("succeeds", func() {
					Expect(createErr).ToNot(HaveOccurred())
				})

				It("creates the volume via BaggageClaim", func() {
					Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(1))

					_, spec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
					Expect(spec).To(Equal(baggageclaim.VolumeSpec{
						Strategy:   baggageclaim.EmptyStrategy{},
						Properties: baggageclaim.VolumeProperties(volumeSpec.Properties),
						TTL:        volumeSpec.TTL,
						Privileged: volumeSpec.Privileged,
					}))
				})

				It("inserts the volume into the database", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(1))

					dbVolume := fakeGardenWorkerDB.InsertVolumeArgsForCall(0)
					Expect(dbVolume).To(Equal(db.Volume{
						Handle:     "created-volume",
						WorkerName: workerName,
						TTL:        volumeSpec.TTL,
						Identifier: db.VolumeIdentifier{
							Output: &db.OutputIdentifier{
								Name: "some-output",
							},
						},
					}))
				})

				It("builds the baggageclaim.Volume into a worker.Volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(1))

					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume))

					Expect(createdVolume).To(Equal(builtVolume))
				})

				Context("when creating the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeBaggageclaimClient.CreateVolumeReturns(nil, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when inserting the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeGardenWorkerDB.InsertVolumeReturns(disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume cannot find the volume in the database", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, nil)
					})

					It("returns ErrMissingVolume", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})
			})

			Context("when creating an HostRootFSStrategy volume", func() {
				BeforeEach(func() {
					volumeSpec.Strategy = HostRootFSStrategy{
						Path:       "some-image-path",
						WorkerName: workerName,
					}
				})

				It("succeeds", func() {
					Expect(createErr).ToNot(HaveOccurred())
				})

				It("creates the volume via BaggageClaim", func() {
					Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(1))

					_, spec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
					Expect(spec).To(Equal(baggageclaim.VolumeSpec{
						Strategy:   baggageclaim.ImportStrategy{Path: "some-image-path"},
						Properties: baggageclaim.VolumeProperties(volumeSpec.Properties),
						TTL:        volumeSpec.TTL,
						Privileged: volumeSpec.Privileged,
					}))
				})

				It("inserts the volume into the database", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(1))

					dbVolume := fakeGardenWorkerDB.InsertVolumeArgsForCall(0)
					Expect(dbVolume).To(Equal(db.Volume{
						Handle:     "created-volume",
						WorkerName: workerName,
						TTL:        volumeSpec.TTL,
						Identifier: db.VolumeIdentifier{
							Import: &db.ImportIdentifier{
								WorkerName: "some-worker",
								Path:       "some-image-path",
							},
						},
					}))
				})

				It("builds the baggageclaim.Volume into a worker.Volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(1))

					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume))

					Expect(createdVolume).To(Equal(builtVolume))
				})

				Context("when creating the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeBaggageclaimClient.CreateVolumeReturns(nil, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when inserting the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeGardenWorkerDB.InsertVolumeReturns(disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when building the volume cannot find the volume in the database", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, nil)
					})

					It("returns ErrMissingVolume", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})
			})
		})
	})

	Describe("LookupVolume", func() {
		var baggageclaimClient baggageclaim.Client

		var handle string

		var foundVolume Volume
		var found bool
		var lookupErr error

		BeforeEach(func() {
			handle = "some-handle"
		})

		JustBeforeEach(func() {
			foundVolume, found, lookupErr = NewGardenWorker(
				fakeGardenClient,
				baggageclaimClient,
				fakeVolumeFactory,
				fakeImageFetcher,
				fakeGardenWorkerDB,
				fakeWorkerProvider,
				fakeClock,
				activeContainers,
				resourceTypes,
				platform,
				tags,
				workerName,
				httpProxyURL,
				httpsProxyURL,
				noProxy,
			).LookupVolume(logger, handle)
		})

		Context("when there is no baggageclaim client", func() {
			BeforeEach(func() {
				baggageclaimClient = nil
			})

			It("returns false", func() {
				Expect(found).To(BeFalse())
			})
		})

		Context("when there is a baggageclaim client", func() {
			BeforeEach(func() {
				baggageclaimClient = fakeBaggageclaimClient
			})

			Context("when the volume can be found on baggageclaim", func() {
				var fakeBaggageclaimVolume *bfakes.FakeVolume
				var builtVolume *wfakes.FakeVolume

				BeforeEach(func() {
					fakeBaggageclaimVolume = new(bfakes.FakeVolume)
					fakeBaggageclaimVolume.HandleReturns(handle)

					fakeBaggageclaimClient.LookupVolumeReturns(fakeBaggageclaimVolume, true, nil)

					builtVolume = new(wfakes.FakeVolume)
					fakeVolumeFactory.BuildReturns(builtVolume, true, nil)
				})

				It("succeeds", func() {
					Expect(lookupErr).ToNot(HaveOccurred())
				})

				It("looks up the volume via BaggageClaim", func() {
					Expect(fakeBaggageclaimClient.LookupVolumeCallCount()).To(Equal(1))

					_, lookedUpHandle := fakeBaggageclaimClient.LookupVolumeArgsForCall(0)
					Expect(lookedUpHandle).To(Equal(handle))
				})

				It("builds the baggageclaim.Volume into a worker.Volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(1))

					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume))

					Expect(foundVolume).To(Equal(builtVolume))
				})

				Context("when building the volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(lookupErr).To(Equal(disaster))
					})
				})

				Context("when building the volume cannot find the volume in the database", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, nil)
					})

					It("returns false", func() {
						Expect(found).To(BeFalse())
					})
				})
			})

			Context("when the volume cannot be found on baggageclaim", func() {
				BeforeEach(func() {
					fakeBaggageclaimClient.LookupVolumeReturns(nil, false, nil)
				})

				It("succeeds", func() {
					Expect(lookupErr).ToNot(HaveOccurred())
				})

				It("returns false", func() {
					Expect(found).To(BeFalse())
				})
			})

			Context("when looking up the volume fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeBaggageclaimClient.LookupVolumeReturns(nil, false, disaster)
				})

				It("returns the error", func() {
					Expect(lookupErr).To(Equal(disaster))
				})
			})
		})
	})

	Describe("ListVolumes", func() {
		var baggageclaimClient baggageclaim.Client

		var properties VolumeProperties

		var foundVolumes []Volume
		var listErr error

		BeforeEach(func() {
			properties = VolumeProperties{
				"some": "properties",
			}
		})

		JustBeforeEach(func() {
			foundVolumes, listErr = NewGardenWorker(
				fakeGardenClient,
				baggageclaimClient,
				fakeVolumeFactory,
				fakeImageFetcher,
				fakeGardenWorkerDB,
				fakeWorkerProvider,
				fakeClock,
				activeContainers,
				resourceTypes,
				platform,
				tags,
				workerName,
				httpProxyURL,
				httpsProxyURL,
				noProxy,
			).ListVolumes(logger, properties)
		})

		Context("when there is no baggageclaim client", func() {
			BeforeEach(func() {
				baggageclaimClient = nil
			})

			It("succeeds", func() {
				Expect(listErr).ToNot(HaveOccurred())
			})

			It("returns no volumes", func() {
				Expect(foundVolumes).To(BeEmpty())
			})
		})

		Context("when there is a baggageclaim client", func() {
			BeforeEach(func() {
				baggageclaimClient = fakeBaggageclaimClient
			})

			Context("when the volume can be found on baggageclaim", func() {
				var fakeBaggageclaimVolume1 *bfakes.FakeVolume
				var fakeBaggageclaimVolume2 *bfakes.FakeVolume
				var fakeBaggageclaimVolume3 *bfakes.FakeVolume

				var builtVolume1 *wfakes.FakeVolume
				var builtVolume3 *wfakes.FakeVolume

				BeforeEach(func() {
					fakeBaggageclaimVolume1 = new(bfakes.FakeVolume)
					fakeBaggageclaimVolume1.HandleReturns("found-volume-1")

					fakeBaggageclaimVolume2 = new(bfakes.FakeVolume)
					fakeBaggageclaimVolume2.HandleReturns("found-volume-2")

					fakeBaggageclaimVolume3 = new(bfakes.FakeVolume)
					fakeBaggageclaimVolume3.HandleReturns("found-volume-3")

					fakeBaggageclaimClient.ListVolumesReturns([]baggageclaim.Volume{
						fakeBaggageclaimVolume1,
						fakeBaggageclaimVolume2,
						fakeBaggageclaimVolume3,
					}, nil)

					builtVolume1 = new(wfakes.FakeVolume)
					builtVolume3 = new(wfakes.FakeVolume)

					fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
						switch volume.Handle() {
						case "found-volume-1":
							return builtVolume1, true, nil
						case "found-volume-2":
							return nil, false, nil
						case "found-volume-3":
							return builtVolume3, true, nil
						default:
							panic("unknown volume: " + volume.Handle())
						}
					}
				})

				It("succeeds", func() {
					Expect(listErr).ToNot(HaveOccurred())
				})

				It("lists up the volumes via BaggageClaim", func() {
					Expect(fakeBaggageclaimClient.ListVolumesCallCount()).To(Equal(1))

					_, listedProperties := fakeBaggageclaimClient.ListVolumesArgsForCall(0)
					Expect(listedProperties).To(Equal(baggageclaim.VolumeProperties(properties)))
				})

				It("builds the baggageclaim.Volumes into a worker.Volume, omitting those who are not found in the database", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(3))

					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume1))

					_, volume = fakeVolumeFactory.BuildArgsForCall(1)
					Expect(volume).To(Equal(fakeBaggageclaimVolume2))

					_, volume = fakeVolumeFactory.BuildArgsForCall(2)
					Expect(volume).To(Equal(fakeBaggageclaimVolume3))

					Expect(foundVolumes).To(Equal([]Volume{builtVolume1, builtVolume3}))
				})

				Context("when building a volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(listErr).To(Equal(disaster))
					})
				})
			})

			Context("when looking up the volume fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeBaggageclaimClient.ListVolumesReturns(nil, disaster)
				})

				It("returns the error", func() {
					Expect(listErr).To(Equal(disaster))
				})
			})
		})
	})

	Describe("FindVolume", func() {
		var (
			foundVolume Volume
			found       bool
			err         error
		)

		JustBeforeEach(func() {
			foundVolume, found, err = gardenWorker.FindVolume(logger, VolumeSpec{
				Strategy: HostRootFSStrategy{
					Path:       "/some/path",
					WorkerName: "worker-name",
				},
			})
		})

		Context("when there is no baggageclaim client", func() {
			BeforeEach(func() {
				gardenWorker = NewGardenWorker(
					fakeGardenClient,
					nil, // baggageclaimClient
					fakeVolumeFactory,
					fakeImageFetcher,
					fakeGardenWorkerDB,
					fakeWorkerProvider,
					fakeClock,
					activeContainers,
					resourceTypes,
					platform,
					tags,
					workerName,
					httpProxyURL,
					httpsProxyURL,
					noProxy,
				)
			})

			It("returns ErrNoVolumeManager", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrNoVolumeManager))
				Expect(found).To(BeFalse())
			})
		})

		It("tries to find the volume in the db", func() {
			Expect(fakeGardenWorkerDB.GetVolumeByIdentifierCallCount()).To(Equal(1))
			Expect(fakeGardenWorkerDB.GetVolumeByIdentifierArgsForCall(0)).To(Equal(db.VolumeIdentifier{
				Import: &db.ImportIdentifier{
					Path:       "/some/path",
					WorkerName: "worker-name",
				},
			}))
		})

		Context("when the volume is found in the db", func() {
			BeforeEach(func() {
				fakeGardenWorkerDB.GetVolumeByIdentifierReturns(db.SavedVolume{
					Volume: db.Volume{
						Handle: "db-vol-handle",
					},
				}, true, nil)
			})

			It("tries to find the db volume in baggageclaim", func() {
				Expect(fakeBaggageclaimClient.LookupVolumeCallCount()).To(Equal(1))
				_, actualHandle := fakeBaggageclaimClient.LookupVolumeArgsForCall(0)
				Expect(actualHandle).To(Equal("db-vol-handle"))
			})

			Context("when the volume can be found in baggageclaim", func() {
				var fakeBaggageclaimVolume *bfakes.FakeVolume

				BeforeEach(func() {
					fakeBaggageclaimVolume = new(bfakes.FakeVolume)
					fakeBaggageclaimVolume.HandleReturns("bg-vol-handle")

					fakeBaggageclaimClient.LookupVolumeReturns(fakeBaggageclaimVolume, true, nil)

				})

				It("tries to build the worker volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(1))
					_, volume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(volume).To(Equal(fakeBaggageclaimVolume))
				})

				Context("when building the worker volume succeeds", func() {
					var builtVolume *wfakes.FakeVolume

					BeforeEach(func() {
						builtVolume = new(wfakes.FakeVolume)
						fakeVolumeFactory.BuildReturns(builtVolume, true, nil)
					})

					It("returns the worker volume", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeTrue())
						Expect(foundVolume).To(Equal(builtVolume))
					})
				})

				Context("when building the worker volume fails", func() {
					disaster := errors.New("nope")

					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, disaster)
					})

					It("returns the error", func() {
						Expect(err).To(Equal(disaster))
						Expect(found).To(BeFalse())
					})
				})

				Context("when the volume ttl cannot be found in the database", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildReturns(nil, false, nil)
					})

					It("does not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(found).To(BeFalse())
					})
				})
			})

			Context("when the volume cannot be found in baggageclaim", func() {
				BeforeEach(func() {
					fakeBaggageclaimClient.LookupVolumeReturns(nil, false, nil)
				})

				It("does not return an error", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeFalse())
				})
			})

			Context("when looking up the volume in baggageclaim fails", func() {
				disaster := errors.New("nope")

				BeforeEach(func() {
					fakeBaggageclaimClient.LookupVolumeReturns(nil, false, disaster)
				})

				It("returns the error", func() {
					Expect(err).To(Equal(disaster))
					Expect(found).To(BeFalse())
				})
			})
		})

		Context("when the volume is not found in the db", func() {
			BeforeEach(func() {
				fakeGardenWorkerDB.GetVolumeByIdentifierReturns(db.SavedVolume{}, false, nil)
			})

			It("returns an error", func() {
				Expect(err).To(Equal(ErrMissingVolume))
				Expect(found).To(BeFalse())
			})
		})

		Context("when finding the volume in the db results in an error", func() {
			var dbErr error

			BeforeEach(func() {
				dbErr = errors.New("an-error")
				fakeGardenWorkerDB.GetVolumeByIdentifierReturns(db.SavedVolume{}, false, dbErr)
			})

			It("returns an error", func() {
				Expect(err).To(Equal(dbErr))
				Expect(found).To(BeFalse())
			})
		})

	})

	Describe("CreateContainer", func() {
		var (
			logger                    lager.Logger
			signals                   <-chan os.Signal
			fakeImageFetchingDelegate *wfakes.FakeImageFetchingDelegate
			containerID               Identifier
			containerMetadata         Metadata
			customTypes               atc.ResourceTypes
			resourceTypeContainerSpec *ResourceTypeContainerSpec
			taskContainerSpec         *TaskContainerSpec

			createdContainer Container
			createErr        error
		)

		BeforeEach(func() {
			logger = lagertest.NewTestLogger("test")

			signals = make(chan os.Signal)
			fakeImageFetchingDelegate = new(wfakes.FakeImageFetchingDelegate)

			containerID = Identifier{
				BuildID: 42,
			}

			containerMetadata = Metadata{
				BuildName: "lol",
			}

			customTypes = atc.ResourceTypes{
				{
					Name:   "custom-type-b",
					Type:   "custom-type-a",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-a",
					Type:   "some-resource",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-c",
					Type:   "custom-type-b",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-d",
					Type:   "custom-type-b",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "unknown-custom-type",
					Type:   "unknown-base-type",
					Source: atc.Source{"some": "source"},
				},
			}
		})

		JustBeforeEach(func() {
			var spec ContainerSpec
			if resourceTypeContainerSpec != nil {
				spec = *resourceTypeContainerSpec
			} else if taskContainerSpec != nil {
				spec = *taskContainerSpec
			}
			createdContainer, createErr = gardenWorker.CreateContainer(logger, signals, fakeImageFetchingDelegate, containerID, containerMetadata, spec, customTypes)
		})

		Context("when the spec is a TaskContainerSpec", func() {
			BeforeEach(func() {
				taskContainerSpec = &TaskContainerSpec{
					Privileged: true,
					Image:      "some-image",
				}
				fakeGardenClient.CreateStub = func(garden.ContainerSpec) (garden.Container, error) {
					return new(gfakes.FakeContainer), nil
				}
			})

			Context("when the spec specifies Inputs", func() {
				var (
					volume1                *wfakes.FakeVolume
					volume2                *wfakes.FakeVolume
					baggageclaimCOWVolume1 *wfakes.FakeVolume
					baggageclaimCOWVolume2 *wfakes.FakeVolume
					cowVolume1             *wfakes.FakeVolume
					cowVolume2             *wfakes.FakeVolume
				)

				BeforeEach(func() {
					volume1 = new(wfakes.FakeVolume)
					volume1.HandleReturns("vol-1-handle")
					volume2 = new(wfakes.FakeVolume)
					volume2.HandleReturns("vol-2-handle")

					taskContainerSpec.Inputs = []VolumeMount{
						{
							volume1,
							"vol-1-mount-path",
						},
						{
							volume2,
							"vol-2-mount-path",
						},
					}

					baggageclaimCOWVolume1 = new(wfakes.FakeVolume)
					baggageclaimCOWVolume1.HandleReturns("bc-cow-vol-1-handle")
					baggageclaimCOWVolume2 = new(wfakes.FakeVolume)
					baggageclaimCOWVolume2.HandleReturns("bc-cow-vol-2-handle")

					fakeBaggageclaimClient.CreateVolumeStub = func(lager.Logger, baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
						switch fakeBaggageclaimClient.CreateVolumeCallCount() {
						case 1:
							return baggageclaimCOWVolume1, nil
						case 2:
							return baggageclaimCOWVolume2, nil
						default:
							return nil, nil
						}
					}

					cowVolume1 = new(wfakes.FakeVolume)
					cowVolume1.HandleReturns("cow-vol-1-handle")
					cowVolume2 = new(wfakes.FakeVolume)
					cowVolume2.HandleReturns("cow-vol-2-handle")

					fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
						switch fakeVolumeFactory.BuildCallCount() {
						case 1:
							return cowVolume1, true, nil
						case 2:
							return cowVolume2, true, nil
						default:
							return new(wfakes.FakeVolume), true, nil
						}
					}
				})

				It("creates a COW volume for each mount", func() {
					Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(2))
					_, volumeSpec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
					Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
						Strategy: baggageclaim.COWStrategy{
							Parent: volume1,
						},
						Privileged: true,
						TTL:        VolumeTTL,
					}))

					_, volumeSpec = fakeBaggageclaimClient.CreateVolumeArgsForCall(1)
					Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
						Strategy: baggageclaim.COWStrategy{
							Parent: volume2,
						},
						Privileged: true,
						TTL:        VolumeTTL,
					}))
				})

				Context("when creating any baggageclaim COW volume fails", func() {
					var baggageclaimCOWVolume1 *wfakes.FakeVolume
					var bccCreateVolumeErr error

					BeforeEach(func() {
						bccCreateVolumeErr = errors.New("an-error")
						baggageclaimCOWVolume1 = new(wfakes.FakeVolume)
						fakeBaggageclaimClient.CreateVolumeStub = func(lager.Logger, baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
							if fakeBaggageclaimClient.CreateVolumeCallCount() != 1 {
								return nil, bccCreateVolumeErr
							}

							return baggageclaimCOWVolume1, nil
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(HaveOccurred())
						Expect(createErr).To(Equal(bccCreateVolumeErr))
					})
				})

				It("inserts each baggageclaim COW volume into the db", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(2))
					Expect(fakeGardenWorkerDB.InsertVolumeArgsForCall(0)).To(Equal(db.Volume{
						Handle:     "bc-cow-vol-1-handle",
						WorkerName: "some-worker",
						TTL:        VolumeTTL,
						Identifier: db.VolumeIdentifier{
							COW: &db.COWIdentifier{
								ParentVolumeHandle: "vol-1-handle",
							},
						},
					}))
				})

				Context("when inserting any volume into the db fails", func() {
					var insertErr error

					BeforeEach(func() {
						insertErr = errors.New("an-error")
						fakeGardenWorkerDB.InsertVolumeStub = func(db.Volume) error {
							if fakeGardenWorkerDB.InsertVolumeCallCount() == 1 {
								return nil
							}
							return insertErr
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(insertErr))
					})
				})

				It("tries to build a COW volume with the volume factory for each baggageclaim COW volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(2))
					_, actualBCCOWVolume1 := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(actualBCCOWVolume1).To(Equal(baggageclaimCOWVolume1))
					_, actualBCCOWVolume2 := fakeVolumeFactory.BuildArgsForCall(1)
					Expect(actualBCCOWVolume2).To(Equal(baggageclaimCOWVolume2))
				})

				Context("when building any COW volume fails", func() {
					var buildErr error

					BeforeEach(func() {
						buildErr = errors.New("an-error")
						fakeVolumeFactory.BuildStub = func(lager.Logger, baggageclaim.Volume) (Volume, bool, error) {
							switch fakeVolumeFactory.BuildCallCount() {
							case 2:
								return nil, false, buildErr
							default:
								return new(wfakes.FakeVolume), true, nil
							}
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(buildErr))
					})
				})

				Context("when any COW volume cannot be found", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildStub = func(lager.Logger, baggageclaim.Volume) (Volume, bool, error) {
							switch fakeVolumeFactory.BuildCallCount() {
							case 2:
								return nil, false, nil
							default:
								return new(wfakes.FakeVolume), true, nil
							}
						}
					})

					It("returns an error", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})

				It("releases each cow volume after attempting to create the container", func() {
					Expect(cowVolume1.ReleaseCallCount()).To(Equal(1))
					Expect(cowVolume1.ReleaseArgsForCall(0)).To(BeNil())
					Expect(cowVolume2.ReleaseCallCount()).To(Equal(1))
					Expect(cowVolume2.ReleaseArgsForCall(0)).To(BeNil())
				})

				It("adds each cow volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					concourseVolumes := []string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
					Expect(err).NotTo(HaveOccurred())
					Expect(concourseVolumes).To(ContainElement("cow-vol-1-handle"))
					Expect(concourseVolumes).To(ContainElement("cow-vol-2-handle"))
				})

				It("adds each cow volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					volumeMountProperties := map[string]string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volume-mounts"]), &volumeMountProperties)
					Expect(err).NotTo(HaveOccurred())
					Expect(volumeMountProperties["cow-vol-1-handle"]).To(Equal("vol-1-mount-path"))
					Expect(volumeMountProperties["cow-vol-2-handle"]).To(Equal("vol-2-mount-path"))
				})
			})

			Context("when the spec specifies Outputs", func() {
				var (
					volume1 *wfakes.FakeVolume
					volume2 *wfakes.FakeVolume
				)

				BeforeEach(func() {
					volume1 = new(wfakes.FakeVolume)
					volume1.HandleReturns("vol-1-handle")
					volume1.PathReturns("vol-1-path")
					volume2 = new(wfakes.FakeVolume)
					volume2.HandleReturns("vol-2-handle")
					volume2.PathReturns("vol-2-path")

					taskContainerSpec.Outputs = []VolumeMount{
						{
							volume1,
							"vol-1-mount-path",
						},
						{
							volume2,
							"vol-2-mount-path",
						},
					}
				})

				It("creates a bind mount for each output volume", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.BindMounts).To(ConsistOf([]garden.BindMount{
						{
							SrcPath: "vol-1-path",
							DstPath: "vol-1-mount-path",
							Mode:    garden.BindMountModeRW,
						},
						{
							SrcPath: "vol-2-path",
							DstPath: "vol-2-mount-path",
							Mode:    garden.BindMountModeRW,
						},
					}))
				})

				It("adds each output volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					concourseVolumes := []string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
					Expect(err).NotTo(HaveOccurred())
					Expect(concourseVolumes).To(ConsistOf([]string{"vol-1-handle", "vol-2-handle"}))
				})
			})

			It("tries to create a container", func() {
				Expect(createErr).NotTo(HaveOccurred())
				Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
				actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
				Expect(actualGardenSpec.Properties["user"]).To(Equal(""))
				Expect(actualGardenSpec.Privileged).To(BeTrue())
			})

			Context("when the spec specifies ImageResource", func() {
				var image *wfakes.FakeImage

				BeforeEach(func() {
					taskContainerSpec.ImageResource = &atc.ImageResource{
						Type:   "some-resource",
						Source: atc.Source{"some": "source"},
					}

					image = new(wfakes.FakeImage)

					imageVolume := new(wfakes.FakeVolume)
					imageVolume.HandleReturns("image-volume")
					imageVolume.PathReturns("/some/image/path")
					image.VolumeReturns(imageVolume)
					image.MetadataReturns(ImageMetadata{
						Env:  []string{"A=1", "B=2"},
						User: "image-volume-user",
					})
					image.VersionReturns(atc.Version{"image": "version"})

					fakeImageFetcher.FetchImageReturns(image, nil)

					fakeGardenClient.CreateStub = func(garden.ContainerSpec) (garden.Container, error) {
						Expect(image.ReleaseCallCount()).To(Equal(0))
						fakeContainer := new(gfakes.FakeContainer)
						return fakeContainer, nil
					}
				})

				It("tries to fetch the image for the resource type", func() {
					Expect(fakeImageFetcher.FetchImageCallCount()).To(Equal(1))
				})

				It("releases the image after creating the container", func() {
					// see fakeGardenClient.CreateStub for the rest of this assertion
					Expect(image.ReleaseCallCount()).To(Equal(1))
				})

				It("creates the container with raw://volume/path/rootfs as the rootfs", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.RootFSPath).To(Equal("raw:///some/image/path/rootfs"))
				})

				It("adds the image volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					concourseVolumes := []string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
					Expect(err).NotTo(HaveOccurred())
					Expect(concourseVolumes).To(ContainElement("image-volume"))
				})

				It("adds the image user to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.Properties["user"]).To(Equal("image-volume-user"))
				})

				Context("when fetching the image fails", func() {
					BeforeEach(func() {
						fakeImageFetcher.FetchImageReturns(nil, errors.New("fetch-err"))
					})

					It("returns an error", func() {
						Expect(createErr).To(HaveOccurred())
						Expect(createErr.Error()).To(Equal("fetch-err"))
					})
				})
			})
		})

		Context("when the spec is a ResourceTypeContainerSpec", func() {
			var image *wfakes.FakeImage
			BeforeEach(func() {
				resourceTypeContainerSpec = &ResourceTypeContainerSpec{
					Type: "custom-type-a",
					Env:  []string{"env-1", "env-2"},
				}

				image = new(wfakes.FakeImage)

				imageVolume := new(wfakes.FakeVolume)
				imageVolume.HandleReturns("image-volume")
				imageVolume.PathReturns("/some/image/path")
				image.VolumeReturns(imageVolume)
				image.MetadataReturns(ImageMetadata{
					Env:  []string{"A=1", "B=2"},
					User: "image-volume-user",
				})
				image.VersionReturns(atc.Version{"image": "version"})

				fakeImageFetcher.FetchImageReturns(image, nil)
				fakeGardenClient.CreateStub = func(garden.ContainerSpec) (garden.Container, error) {
					Expect(image.ReleaseCallCount()).To(Equal(0))
					fakeContainer := new(gfakes.FakeContainer)
					return fakeContainer, nil
				}
			})

			It("tries to fetch the image for the resource type", func() {
				Expect(fakeImageFetcher.FetchImageCallCount()).To(Equal(1))
				_, fetchImageConfig, fetchSignals, fetchID, fetchMetadata, fetchDelegate, fetchWorker, fetchTags, fetchCustomTypes, fetchPrivileged := fakeImageFetcher.FetchImageArgsForCall(0)
				Expect(fetchImageConfig).To(Equal(atc.ImageResource{
					Type:   "some-resource",
					Source: atc.Source{"some": "source"},
				}))
				Expect(fetchSignals).To(Equal(signals))
				Expect(fetchID).To(Equal(containerID))
				Expect(fetchMetadata).To(Equal(containerMetadata))
				Expect(fetchDelegate).To(Equal(fakeImageFetchingDelegate))
				Expect(fetchWorker).To(Equal(gardenWorker))
				Expect(fetchTags).To(Equal(atc.Tags{"some", "tags"}))
				Expect(fetchCustomTypes).To(Equal(customTypes.Without("custom-type-a")))
				Expect(fetchPrivileged).To(Equal(true))
			})

			It("releases the image after creating the container", func() {
				// see fakeGardenClient.CreateStub for the rest of this assertion
				Expect(image.ReleaseCallCount()).To(Equal(1))
			})

			It("adds the image volume to the garden spec properties", func() {
				Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
				actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
				concourseVolumes := []string{}
				err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
				Expect(err).NotTo(HaveOccurred())
				Expect(concourseVolumes).To(ContainElement("image-volume"))
			})

			It("adds the image user to the garden spec properties", func() {
				Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
				actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
				Expect(actualGardenSpec.Properties["user"]).To(Equal("image-volume-user"))
			})

			Context("when fetching the image fails", func() {
				BeforeEach(func() {
					fakeImageFetcher.FetchImageReturns(nil, errors.New("fetch-err"))
				})

				It("returns an error", func() {
					Expect(createErr).To(HaveOccurred())
					Expect(createErr.Error()).To(Equal("fetch-err"))
				})
			})

			It("tries to create a container", func() {
				Expect(createErr).NotTo(HaveOccurred())
				Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
				actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
				expectedEnv := resourceTypeContainerSpec.Env
				expectedEnv = append([]string{"A=1", "B=2"}, expectedEnv...)
				Expect(actualGardenSpec.Env).To(Equal(expectedEnv))
				Expect(actualGardenSpec.Properties["user"]).To(Equal("image-volume-user"))
				Expect(actualGardenSpec.Privileged).To(BeTrue())
				Expect(actualGardenSpec.RootFSPath).To(Equal("raw:///some/image/path/rootfs"))
			})

			Context("when the worker has a HTTPProxyURL", func() {
				BeforeEach(func() {
					gardenWorker = NewGardenWorker(
						fakeGardenClient,
						fakeBaggageclaimClient,
						fakeVolumeFactory,
						fakeImageFetcher,
						fakeGardenWorkerDB,
						fakeWorkerProvider,
						fakeClock,
						activeContainers,
						resourceTypes,
						platform,
						tags,
						workerName,
						"http://example.com",
						httpsProxyURL,
						noProxy,
					)
				})

				It("adds the proxy url to the garden spec env", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.Env).To(ContainElement("http_proxy=http://example.com"))
				})
			})

			Context("when the worker has NoProxy", func() {
				BeforeEach(func() {
					gardenWorker = NewGardenWorker(
						fakeGardenClient,
						fakeBaggageclaimClient,
						fakeVolumeFactory,
						fakeImageFetcher,
						fakeGardenWorkerDB,
						fakeWorkerProvider,
						fakeClock,
						activeContainers,
						resourceTypes,
						platform,
						tags,
						workerName,
						httpProxyURL,
						httpsProxyURL,
						"localhost",
					)
				})

				It("adds the proxy url to the garden spec env", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.Env).To(ContainElement("no_proxy=localhost"))
				})
			})

			Context("when the worker has a HTTPSProxyURL", func() {
				BeforeEach(func() {
					gardenWorker = NewGardenWorker(
						fakeGardenClient,
						fakeBaggageclaimClient,
						fakeVolumeFactory,
						fakeImageFetcher,
						fakeGardenWorkerDB,
						fakeWorkerProvider,
						fakeClock,
						activeContainers,
						resourceTypes,
						platform,
						tags,
						workerName,
						httpProxyURL,
						"https://example.com",
						noProxy,
					)
				})

				It("adds the proxy url to the garden spec env", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.Env).To(ContainElement("https_proxy=https://example.com"))
				})
			})

			Context("when the spec specifies Ephemeral", func() {
				BeforeEach(func() {
					resourceTypeContainerSpec.Ephemeral = true
				})

				It("creates the container with ephemeral = true", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.Properties["concourse:ephemeral"]).To(Equal("true"))
				})
			})

			Context("when the spec specifies Mounts", func() {
				var (
					volume1                *wfakes.FakeVolume
					volume2                *wfakes.FakeVolume
					baggageclaimCOWVolume1 *wfakes.FakeVolume
					baggageclaimCOWVolume2 *wfakes.FakeVolume
					cowVolume1             *wfakes.FakeVolume
					cowVolume2             *wfakes.FakeVolume
				)

				BeforeEach(func() {
					volume1 = new(wfakes.FakeVolume)
					volume1.HandleReturns("vol-1-handle")
					volume2 = new(wfakes.FakeVolume)
					volume2.HandleReturns("vol-2-handle")

					resourceTypeContainerSpec.Mounts = []VolumeMount{
						{
							volume1,
							"vol-1-mount-path",
						},
						{
							volume2,
							"vol-2-mount-path",
						},
					}

					baggageclaimCOWVolume1 = new(wfakes.FakeVolume)
					baggageclaimCOWVolume1.HandleReturns("bc-cow-vol-1-handle")
					baggageclaimCOWVolume2 = new(wfakes.FakeVolume)
					baggageclaimCOWVolume2.HandleReturns("bc-cow-vol-2-handle")

					fakeBaggageclaimClient.CreateVolumeStub = func(lager.Logger, baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
						switch fakeBaggageclaimClient.CreateVolumeCallCount() {
						case 1:
							return baggageclaimCOWVolume1, nil
						case 2:
							return baggageclaimCOWVolume2, nil
						default:
							return nil, nil
						}
					}

					cowVolume1 = new(wfakes.FakeVolume)
					cowVolume1.HandleReturns("cow-vol-1-handle")
					cowVolume2 = new(wfakes.FakeVolume)
					cowVolume2.HandleReturns("cow-vol-2-handle")

					fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
						switch fakeVolumeFactory.BuildCallCount() {
						case 1:
							return cowVolume1, true, nil
						case 2:
							return cowVolume2, true, nil
						default:
							return new(wfakes.FakeVolume), true, nil
						}
					}
				})

				It("creates a COW volume for each mount", func() {
					Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(2))
					_, volumeSpec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
					Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
						Strategy: baggageclaim.COWStrategy{
							Parent: volume1,
						},
						Privileged: true,
						TTL:        VolumeTTL,
					}))

					_, volumeSpec = fakeBaggageclaimClient.CreateVolumeArgsForCall(1)
					Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
						Strategy: baggageclaim.COWStrategy{
							Parent: volume2,
						},
						Privileged: true,
						TTL:        VolumeTTL,
					}))
				})

				Context("when creating any baggageclaim COW volume fails", func() {
					var baggageclaimCOWVolume1 *wfakes.FakeVolume
					var bccCreateVolumeErr error

					BeforeEach(func() {
						bccCreateVolumeErr = errors.New("an-error")
						baggageclaimCOWVolume1 = new(wfakes.FakeVolume)
						fakeBaggageclaimClient.CreateVolumeStub = func(lager.Logger, baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
							if fakeBaggageclaimClient.CreateVolumeCallCount() != 1 {
								return nil, bccCreateVolumeErr
							}

							return baggageclaimCOWVolume1, nil
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(HaveOccurred())
						Expect(createErr).To(Equal(bccCreateVolumeErr))
					})
				})

				It("inserts each baggageclaim COW volume into the db", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(2))
					Expect(fakeGardenWorkerDB.InsertVolumeArgsForCall(0)).To(Equal(db.Volume{
						Handle:     "bc-cow-vol-1-handle",
						WorkerName: "some-worker",
						TTL:        VolumeTTL,
						Identifier: db.VolumeIdentifier{
							COW: &db.COWIdentifier{
								ParentVolumeHandle: "vol-1-handle",
							},
						},
					}))
				})

				Context("when inserting any volume into the db fails", func() {
					var insertErr error

					BeforeEach(func() {
						insertErr = errors.New("an-error")
						fakeGardenWorkerDB.InsertVolumeStub = func(db.Volume) error {
							if fakeGardenWorkerDB.InsertVolumeCallCount() == 1 {
								return nil
							}
							return insertErr
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(insertErr))
					})
				})

				It("tries to build a COW volume with the volume factory for each baggageclaim COW volume", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(2))
					_, actualBCCOWVolume1 := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(actualBCCOWVolume1).To(Equal(baggageclaimCOWVolume1))
					_, actualBCCOWVolume2 := fakeVolumeFactory.BuildArgsForCall(1)
					Expect(actualBCCOWVolume2).To(Equal(baggageclaimCOWVolume2))
				})

				Context("when building any COW volume fails", func() {
					var buildErr error

					BeforeEach(func() {
						buildErr = errors.New("an-error")
						fakeVolumeFactory.BuildStub = func(lager.Logger, baggageclaim.Volume) (Volume, bool, error) {
							switch fakeVolumeFactory.BuildCallCount() {
							case 2:
								return nil, false, buildErr
							default:
								return new(wfakes.FakeVolume), true, nil
							}
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(buildErr))
					})
				})

				Context("when any COW volume cannot be found", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildStub = func(lager.Logger, baggageclaim.Volume) (Volume, bool, error) {
							switch fakeVolumeFactory.BuildCallCount() {
							case 2:
								return nil, false, nil
							default:
								return new(wfakes.FakeVolume), true, nil
							}
						}
					})

					It("returns an error", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})

				It("releases each cow volume after attempting to create the container", func() {
					Expect(cowVolume1.ReleaseCallCount()).To(Equal(1))
					Expect(cowVolume1.ReleaseArgsForCall(0)).To(BeNil())
					Expect(cowVolume2.ReleaseCallCount()).To(Equal(1))
					Expect(cowVolume2.ReleaseArgsForCall(0)).To(BeNil())
				})

				It("does not release the volumes that were passed in", func() {
					Expect(volume1.ReleaseCallCount()).To(BeZero())
					Expect(volume2.ReleaseCallCount()).To(BeZero())
				})

				It("adds each cow volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					concourseVolumes := []string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
					Expect(err).NotTo(HaveOccurred())
					Expect(concourseVolumes).To(ContainElement("cow-vol-1-handle"))
					Expect(concourseVolumes).To(ContainElement("cow-vol-2-handle"))
				})

				It("adds each cow volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					volumeMountProperties := map[string]string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volume-mounts"]), &volumeMountProperties)
					Expect(err).NotTo(HaveOccurred())
					Expect(volumeMountProperties["cow-vol-1-handle"]).To(Equal("vol-1-mount-path"))
					Expect(volumeMountProperties["cow-vol-2-handle"]).To(Equal("vol-2-mount-path"))
				})
			})

			Context("when the spec specifies a resource type that a worker provides", func() {
				var (
					importBCVolume *bfakes.FakeVolume
					importVolume   *wfakes.FakeVolume
					cowBCVolume    *bfakes.FakeVolume
					cowVolume      *wfakes.FakeVolume
				)

				BeforeEach(func() {
					resourceTypeContainerSpec.Type = "some-resource"

					fakeBaggageclaimClient.CreateVolumeStub = func(logger lager.Logger, volumeSpec baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
						switch volumeSpec.Strategy.(type) {
						case baggageclaim.ImportStrategy:
							return importBCVolume, nil
						case baggageclaim.COWStrategy:
							return cowBCVolume, nil
						default:
							return nil, nil
						}
					}

					importBCVolume = new(bfakes.FakeVolume)
					importBCVolume.HandleReturns("import-bc-vol")

					importVolume = new(wfakes.FakeVolume)
					importVolume.HandleReturns("import-vol")

					cowBCVolume = new(bfakes.FakeVolume)
					cowBCVolume.HandleReturns("cow-bc-vol")

					cowVolume = new(wfakes.FakeVolume)
					cowVolume.HandleReturns("cow-vol")
					cowVolume.PathReturns("cow-vol-path")

					fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
						switch volume.Handle() {
						case "import-bc-vol":
							return importVolume, true, nil
						case "cow-bc-vol":
							return cowVolume, true, nil
						default:
							return new(wfakes.FakeVolume), true, nil
						}
					}
				})

				It("tries to find an existing import volume", func() {
					Expect(fakeGardenWorkerDB.GetVolumeByIdentifierCallCount()).To(Equal(1))
					Expect(fakeGardenWorkerDB.GetVolumeByIdentifierArgsForCall(0)).To(Equal(db.VolumeIdentifier{
						Import: &db.ImportIdentifier{
							WorkerName: "some-worker",
							Path:       "some-resource-image",
						},
					}))
				})

				It("releases the import volume", func() {
					Expect(importVolume.ReleaseCallCount()).To(Equal(1))
					Expect(importVolume.ReleaseArgsForCall(0)).To(BeNil())
				})

				It("releases the cow volume after attempting to create the container", func() {
					Expect(cowVolume.ReleaseCallCount()).To(Equal(1))
					Expect(cowVolume.ReleaseArgsForCall(0)).To(BeNil())
				})

				It("inserts import and COW volumes into the db", func() {
					Expect(fakeGardenWorkerDB.InsertVolumeCallCount()).To(Equal(2))
					Expect(fakeGardenWorkerDB.InsertVolumeArgsForCall(0)).To(Equal(db.Volume{
						Handle:     "import-bc-vol",
						WorkerName: "some-worker",
						TTL:        0,
						Identifier: db.VolumeIdentifier{
							Import: &db.ImportIdentifier{
								WorkerName: "some-worker",
								Path:       "some-resource-image",
							},
						},
					}))

					Expect(fakeGardenWorkerDB.InsertVolumeArgsForCall(1)).To(Equal(db.Volume{
						Handle:     "cow-bc-vol",
						WorkerName: "some-worker",
						TTL:        VolumeTTL,
						Identifier: db.VolumeIdentifier{
							COW: &db.COWIdentifier{
								ParentVolumeHandle: "import-vol",
							},
						},
					}))
				})

				It("adds the cow volume to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					concourseVolumes := []string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volumes"]), &concourseVolumes)
					Expect(err).NotTo(HaveOccurred())
					Expect(concourseVolumes).To(ContainElement("cow-vol"))
				})

				It("adds the cow volume mount to the garden spec properties", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					volumeMountProperties := map[string]string{}
					err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volume-mounts"]), &volumeMountProperties)
					Expect(err).NotTo(HaveOccurred())
					Expect(volumeMountProperties["cow-vol"]).To(Equal("cow-vol-path"))
				})

				It("uses the path of the cow volume as the rootfs", func() {
					Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
					actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
					Expect(actualGardenSpec.RootFSPath).To(Equal("raw://cow-vol-path"))
				})

				It("tries to build an import and COW volume with the volume factory", func() {
					Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(2))
					_, actualBCImportVolume := fakeVolumeFactory.BuildArgsForCall(0)
					Expect(actualBCImportVolume).To(Equal(importBCVolume))
					_, actualBCCOWVolume := fakeVolumeFactory.BuildArgsForCall(1)
					Expect(actualBCCOWVolume).To(Equal(cowBCVolume))
				})

				Context("when the import volume can be retrieved", func() {
					BeforeEach(func() {
						fakeGardenWorkerDB.GetVolumeByIdentifierReturns(db.SavedVolume{
							Volume: db.Volume{
								Handle:     "imported-volume-handle",
								WorkerName: "worker-name",
								TTL:        1 * time.Millisecond,
								Identifier: db.VolumeIdentifier{},
							},
							ID: 5,
						}, true, nil)

						fakeBaggageclaimClient.LookupVolumeReturns(importBCVolume, true, nil)
					})

					It("tries to find the volume in baggageclaim", func() {
						Expect(fakeBaggageclaimClient.LookupVolumeCallCount()).To(Equal(1))
						_, handle := fakeBaggageclaimClient.LookupVolumeArgsForCall(0)
						Expect(handle).To(Equal("imported-volume-handle"))
					})

					It("builds a worker volume from the baggageclaim volume", func() {
						Expect(fakeVolumeFactory.BuildCallCount()).To(Equal(2))
						_, actualBCVolume := fakeVolumeFactory.BuildArgsForCall(0)
						Expect(actualBCVolume).To(Equal(importBCVolume))
					})

					It("creates a COW volume with the import volume as its parent", func() {
						Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(1))
						_, volumeSpec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
						Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
							Strategy: baggageclaim.COWStrategy{
								Parent: importVolume,
							},
							Privileged: true,
							TTL:        VolumeTTL,
							Properties: baggageclaim.VolumeProperties{},
						}))
					})
				})

				Context("when the import volume cannot be retrieved", func() {
					BeforeEach(func() {
						fakeGardenWorkerDB.GetVolumeByIdentifierReturns(db.SavedVolume{}, false, nil)
					})

					It("creates import and COW volumes for resource image", func() {
						Expect(fakeBaggageclaimClient.CreateVolumeCallCount()).To(Equal(2))
						_, volumeSpec := fakeBaggageclaimClient.CreateVolumeArgsForCall(0)
						Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
							Strategy: baggageclaim.ImportStrategy{
								Path: "some-resource-image",
							},
							Privileged: true,
							TTL:        0,
							Properties: baggageclaim.VolumeProperties{},
						}))

						_, volumeSpec = fakeBaggageclaimClient.CreateVolumeArgsForCall(1)
						Expect(volumeSpec).To(Equal(baggageclaim.VolumeSpec{
							Strategy: baggageclaim.COWStrategy{
								Parent: importVolume,
							},
							Privileged: true,
							TTL:        VolumeTTL,
							Properties: baggageclaim.VolumeProperties{},
						}))
					})

					Context("when creating an import volume fails", func() {
						var disaster error
						BeforeEach(func() {
							disaster = errors.New("failed-to-create-volume")
							fakeBaggageclaimClient.CreateVolumeReturns(nil, disaster)
						})

						It("returns the error", func() {
							Expect(createErr).To(Equal(disaster))
						})
					})
				})

				Context("when inserting the import volume into the db fails", func() {
					var insertErr error

					BeforeEach(func() {
						insertErr = errors.New("an-error")
						fakeGardenWorkerDB.InsertVolumeReturns(insertErr)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(insertErr))
					})
				})

				Context("when building the import volume fails", func() {
					var buildErr error

					BeforeEach(func() {
						buildErr = errors.New("an-error")
						fakeVolumeFactory.BuildReturns(nil, false, buildErr)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(buildErr))
					})
				})

				Context("when creating the COW volume fails", func() {
					var disaster error
					BeforeEach(func() {
						disaster = errors.New("failed-to-create-volume")
						fakeBaggageclaimClient.CreateVolumeStub = func(logger lager.Logger, volumeSpec baggageclaim.VolumeSpec) (baggageclaim.Volume, error) {
							switch volumeSpec.Strategy.(type) {
							case baggageclaim.ImportStrategy:
								return importBCVolume, nil
							case baggageclaim.COWStrategy:
								return nil, disaster
							default:
								return nil, nil
							}
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(disaster))
					})
				})

				Context("when inserting the COW volume into the db fails", func() {
					var insertErr error

					BeforeEach(func() {
						insertErr = errors.New("an-error")
						fakeGardenWorkerDB.InsertVolumeStub = func(volume db.Volume) error {
							switch volume.Handle {
							case "cow-bc-vol":
								return insertErr
							default:
								return nil
							}
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(insertErr))
					})
				})

				Context("when building the COW volume fails", func() {
					var buildErr error

					BeforeEach(func() {
						buildErr = errors.New("an-error")
						fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
							switch volume.Handle() {
							case "import-bc-vol":
								return importVolume, true, nil
							case "cow-bc-vol":
								return nil, false, buildErr
							default:
								return nil, false, nil
							}
						}
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(buildErr))
					})
				})

				Context("when the COW volume cannot be found", func() {
					BeforeEach(func() {
						fakeVolumeFactory.BuildStub = func(logger lager.Logger, volume baggageclaim.Volume) (Volume, bool, error) {
							switch volume.Handle() {
							case "import-bc-vol":
								return importVolume, true, nil
							default:
								return nil, false, nil
							}
						}
					})

					It("returns an error", func() {
						Expect(createErr).To(Equal(ErrMissingVolume))
					})
				})
			})

			Context("when the spec specifies a resource type that is unknown", func() {
				BeforeEach(func() {
					resourceTypeContainerSpec.Type = "some-bogus-resource"
				})

				It("returns ErrUnsupportedResourceType", func() {
					Expect(createErr).To(Equal(ErrUnsupportedResourceType))
				})
			})

			Context("when creating the container succeeds", func() {
				var fakeContainer *gfakes.FakeContainer
				BeforeEach(func() {
					fakeContainer = new(gfakes.FakeContainer)
					fakeContainer.HandleReturns("some-container-handle")
					fakeGardenClient.CreateReturns(fakeContainer, nil)
				})

				It("tries to create the container in the db", func() {
					Expect(fakeGardenWorkerDB.CreateContainerCallCount()).To(Equal(1))
					c, ttl := fakeGardenWorkerDB.CreateContainerArgsForCall(0)

					expectedContainerID := containerID
					expectedContainerID.ResourceTypeVersion = atc.Version{"image": "version"}

					expectedContainerMetadata := containerMetadata
					expectedContainerMetadata.Handle = "some-container-handle"
					expectedContainerMetadata.User = "image-volume-user"
					expectedContainerMetadata.WorkerName = "some-worker"

					Expect(c).To(Equal(db.Container{
						ContainerIdentifier: db.ContainerIdentifier(expectedContainerID),
						ContainerMetadata:   db.ContainerMetadata(expectedContainerMetadata),
					}))

					Expect(ttl).To(Equal(ContainerTTL))
				})

				It("returns a container that be destroyed", func() {
					err := createdContainer.Destroy()
					Expect(err).NotTo(HaveOccurred())

					By("destroying via garden")
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(1))
					Expect(fakeGardenClient.DestroyArgsForCall(0)).To(Equal("some-container-handle"))

					By("no longer heartbeating")
					fakeClock.Increment(30 * time.Second)
					Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(1))
				})

				It("performs an initial heartbeat synchronously on the returned container", func() {
					Expect(fakeContainer.SetGraceTimeCallCount()).To(Equal(1))
					Expect(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount()).To(Equal(1))
				})

				It("heartbeats to the database and the container", func() {
					fakeClock.Increment(30 * time.Second)

					Eventually(fakeContainer.SetGraceTimeCallCount).Should(Equal(2))
					Expect(fakeContainer.SetGraceTimeArgsForCall(1)).To(Equal(5 * time.Minute))

					Eventually(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(2))
					handle, interval := fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(1)
					Expect(handle).To(Equal("some-container-handle"))
					Expect(interval).To(Equal(5 * time.Minute))

					fakeClock.Increment(30 * time.Second)

					Eventually(fakeContainer.SetGraceTimeCallCount).Should(Equal(3))
					Expect(fakeContainer.SetGraceTimeArgsForCall(2)).To(Equal(5 * time.Minute))

					Eventually(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(3))
					handle, interval = fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(2)
					Expect(handle).To(Equal("some-container-handle"))
					Expect(interval).To(Equal(5 * time.Minute))
				})

				It("sets a final ttl on the container and stops heartbeating when the container is released", func() {
					createdContainer.Release(FinalTTL(30 * time.Minute))

					Expect(fakeContainer.SetGraceTimeCallCount()).Should(Equal(2))
					Expect(fakeContainer.SetGraceTimeArgsForCall(1)).To(Equal(30 * time.Minute))

					Expect(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount()).Should(Equal(2))
					handle, interval := fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(1)
					Expect(handle).To(Equal("some-container-handle"))
					Expect(interval).To(Equal(30 * time.Minute))

					fakeClock.Increment(30 * time.Second)

					Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(2))
					Consistently(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(2))
				})

				It("does not perform a final heartbeat when there is no final ttl", func() {
					createdContainer.Release(nil)

					Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(1))
					Consistently(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(1))
				})

				Context("when creating the container in the db fails", func() {
					var gardenWorkerDBCreateContainerErr error
					BeforeEach(func() {
						gardenWorkerDBCreateContainerErr = errors.New("an-error")
						fakeGardenWorkerDB.CreateContainerReturns(db.SavedContainer{}, gardenWorkerDBCreateContainerErr)
					})

					It("returns the error", func() {
						Expect(createErr).To(Equal(gardenWorkerDBCreateContainerErr))
					})
				})

				Context("when creating the container in the db succeeds", func() {
					BeforeEach(func() {
						fakeGardenWorkerDB.CreateContainerReturns(db.SavedContainer{}, nil)
					})

					It("returns a Container", func() {
						Expect(createdContainer).NotTo(BeNil())
					})
				})
			})

			Context("when creating the container fails", func() {
				var gardenCreateErr error

				BeforeEach(func() {
					gardenCreateErr = errors.New("an-error")
					fakeGardenClient.CreateReturns(nil, gardenCreateErr)
				})

				It("returns the error", func() {
					Expect(createErr).To(HaveOccurred())
					Expect(createErr).To(Equal(gardenCreateErr))
				})
			})

			Context("when the spec defines a cache volume", func() {
				BeforeEach(func() {
					v := new(wfakes.FakeVolume)
					v.PathReturns("cache-volume-src-path")
					v.HandleReturns("cache-volume-handle")
					resourceTypeContainerSpec.Cache.Volume = v
				})

				Context("when the spec defines a cache mount path", func() {
					BeforeEach(func() {
						resourceTypeContainerSpec.Cache.MountPath = "cache-volume-mount-path"
					})

					It("creates a bind mount for the cache volume", func() {
						Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
						actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
						Expect(actualGardenSpec.BindMounts).To(ContainElement(garden.BindMount{
							SrcPath: "cache-volume-src-path",
							DstPath: "cache-volume-mount-path",
							Mode:    garden.BindMountModeRW,
						}))
					})

					It("adds the cache volume to the garden spec properties", func() {
						Expect(fakeGardenClient.CreateCallCount()).To(Equal(1))
						actualGardenSpec := fakeGardenClient.CreateArgsForCall(0)
						volumeMountProperties := map[string]string{}
						err := json.Unmarshal([]byte(actualGardenSpec.Properties["concourse:volume-mounts"]), &volumeMountProperties)
						Expect(err).NotTo(HaveOccurred())
						Expect(volumeMountProperties["cache-volume-handle"]).To(Equal("cache-volume-mount-path"))
					})
				})

				Context("when the spec also defines mounts", func() {
					BeforeEach(func() {
						resourceTypeContainerSpec.Mounts = []VolumeMount{
							{
								new(wfakes.FakeVolume),
								"mount-path",
							},
						}
					})

					It("returns an error", func() {
						Expect(createErr).To(HaveOccurred())
						Expect(createErr.Error()).To(Equal("a container may not have mounts and a cache"))
					})
				})
			})
		})
	})

	Describe("LookupContainer", func() {
		var handle string

		BeforeEach(func() {
			handle = "we98lsv"
		})

		Context("when the gardenClient returns a container and no error", func() {
			var (
				fakeContainer *gfakes.FakeContainer
			)

			BeforeEach(func() {
				fakeContainer = new(gfakes.FakeContainer)
				fakeContainer.HandleReturns("some-handle")

				fakeGardenClient.LookupReturns(fakeContainer, nil)
			})

			It("returns the container and no error", func() {
				foundContainer, found, err := gardenWorker.LookupContainer(logger, handle)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())

				Expect(foundContainer.Handle()).To(Equal(fakeContainer.Handle()))
			})

			Describe("the container", func() {
				var foundContainer Container
				var findErr error

				JustBeforeEach(func() {
					foundContainer, _, findErr = gardenWorker.LookupContainer(logger, handle)
				})

				Context("when the concourse:volumes property is present", func() {
					var (
						handle1Volume         *wfakes.FakeVolume
						handle2Volume         *wfakes.FakeVolume
						expectedHandle1Volume *wfakes.FakeVolume
						expectedHandle2Volume *wfakes.FakeVolume
					)

					BeforeEach(func() {
						handle1Volume = new(wfakes.FakeVolume)
						handle2Volume = new(wfakes.FakeVolume)
						expectedHandle1Volume = new(wfakes.FakeVolume)
						expectedHandle2Volume = new(wfakes.FakeVolume)

						fakeContainer.PropertiesReturns(garden.Properties{
							"concourse:volumes":       `["handle-1","handle-2"]`,
							"concourse:volume-mounts": `{"handle-1":"/handle-1/path","handle-2":"/handle-2/path"}`,
						}, nil)

						fakeBaggageclaimClient.LookupVolumeStub = func(logger lager.Logger, handle string) (baggageclaim.Volume, bool, error) {
							if handle == "handle-1" {
								return handle1Volume, true, nil
							} else if handle == "handle-2" {
								return handle2Volume, true, nil
							} else {
								panic("unknown handle: " + handle)
							}
						}

						fakeVolumeFactory.BuildStub = func(logger lager.Logger, vol baggageclaim.Volume) (Volume, bool, error) {
							if vol == handle1Volume {
								return expectedHandle1Volume, true, nil
							} else if vol == handle2Volume {
								return expectedHandle2Volume, true, nil
							} else {
								panic("unknown volume: " + vol.Handle())
							}
						}
					})

					Describe("Volumes", func() {
						It("returns all bound volumes based on properties on the container", func() {
							Expect(foundContainer.Volumes()).To(Equal([]Volume{
								expectedHandle1Volume,
								expectedHandle2Volume,
							}))
						})

						Context("when LookupVolume returns an error", func() {
							disaster := errors.New("nope")

							BeforeEach(func() {
								fakeBaggageclaimClient.LookupVolumeReturns(nil, false, disaster)
							})

							It("returns the error on lookup", func() {
								Expect(findErr).To(Equal(disaster))
							})
						})

						Context("when LookupVolume cannot find the volume", func() {
							BeforeEach(func() {
								fakeBaggageclaimClient.LookupVolumeReturns(nil, false, nil)
							})

							It("returns ErrMissingVolume", func() {
								Expect(findErr).To(Equal(ErrMissingVolume))
							})
						})

						Context("when Build cannot find the volume", func() {
							BeforeEach(func() {
								fakeVolumeFactory.BuildReturns(nil, false, nil)
							})

							It("returns ErrMissingVolume", func() {
								Expect(findErr).To(Equal(ErrMissingVolume))
							})
						})

						Context("when Build returns an error", func() {
							disaster := errors.New("nope")

							BeforeEach(func() {
								fakeVolumeFactory.BuildReturns(nil, false, disaster)
							})

							It("returns the error on lookup", func() {
								Expect(findErr).To(Equal(disaster))
							})
						})

						Context("when there is no baggageclaim", func() {
							BeforeEach(func() {
								gardenWorker = NewGardenWorker(
									fakeGardenClient,
									nil,
									nil,
									fakeImageFetcher,
									fakeGardenWorkerDB,
									fakeWorkerProvider,
									fakeClock,
									activeContainers,
									resourceTypes,
									platform,
									tags,
									workerName,
									httpProxyURL,
									httpsProxyURL,
									noProxy,
								)
							})

							It("returns an empty slice", func() {
								Expect(foundContainer.Volumes()).To(BeEmpty())
							})
						})
					})

					Describe("VolumeMounts", func() {
						It("returns all bound volumes based on properties on the container", func() {
							Expect(foundContainer.VolumeMounts()).To(ConsistOf([]VolumeMount{
								{Volume: expectedHandle1Volume, MountPath: "/handle-1/path"},
								{Volume: expectedHandle2Volume, MountPath: "/handle-2/path"},
							}))
						})

						Context("when LookupVolume returns an error", func() {
							disaster := errors.New("nope")

							BeforeEach(func() {
								fakeBaggageclaimClient.LookupVolumeReturns(nil, false, disaster)
							})

							It("returns the error on lookup", func() {
								Expect(findErr).To(Equal(disaster))
							})
						})

						Context("when Build returns an error", func() {
							disaster := errors.New("nope")

							BeforeEach(func() {
								fakeVolumeFactory.BuildReturns(nil, false, disaster)
							})

							It("returns the error on lookup", func() {
								Expect(findErr).To(Equal(disaster))
							})
						})

						Context("when there is no baggageclaim", func() {
							BeforeEach(func() {
								gardenWorker = NewGardenWorker(
									fakeGardenClient,
									nil,
									nil,
									fakeImageFetcher,
									fakeGardenWorkerDB,
									fakeWorkerProvider,
									fakeClock,
									activeContainers,
									resourceTypes,
									platform,
									tags,
									workerName,
									httpProxyURL,
									httpsProxyURL,
									noProxy,
								)
							})

							It("returns an empty slice", func() {
								Expect(foundContainer.Volumes()).To(BeEmpty())
							})
						})
					})

					Describe("Release", func() {
						It("releases the container's volumes once and only once", func() {
							foundContainer.Release(FinalTTL(time.Minute))
							Expect(expectedHandle1Volume.ReleaseCallCount()).To(Equal(1))
							Expect(expectedHandle1Volume.ReleaseArgsForCall(0)).To(Equal(FinalTTL(time.Minute)))
							Expect(expectedHandle2Volume.ReleaseCallCount()).To(Equal(1))
							Expect(expectedHandle2Volume.ReleaseArgsForCall(0)).To(Equal(FinalTTL(time.Minute)))

							foundContainer.Release(FinalTTL(time.Hour))
							Expect(expectedHandle1Volume.ReleaseCallCount()).To(Equal(1))
							Expect(expectedHandle2Volume.ReleaseCallCount()).To(Equal(1))
						})
					})
				})

				Context("when the concourse:volumes property is not present", func() {
					BeforeEach(func() {
						fakeContainer.PropertiesReturns(garden.Properties{}, nil)
					})

					Describe("Volumes", func() {
						It("returns an empty slice", func() {
							Expect(foundContainer.Volumes()).To(BeEmpty())
						})
					})
				})

				Context("when the user property is present", func() {
					var (
						actualSpec garden.ProcessSpec
						actualIO   garden.ProcessIO
					)

					BeforeEach(func() {
						actualSpec = garden.ProcessSpec{
							Path: "some-path",
							Args: []string{"some", "args"},
							Env:  []string{"some=env"},
							Dir:  "some-dir",
						}

						actualIO = garden.ProcessIO{}

						fakeContainer.PropertiesReturns(garden.Properties{"user": "maverick"}, nil)
					})

					JustBeforeEach(func() {
						foundContainer.Run(actualSpec, actualIO)
					})

					Describe("Run", func() {
						It("calls Run() on the garden container and injects the user", func() {
							Expect(fakeContainer.RunCallCount()).To(Equal(1))
							spec, io := fakeContainer.RunArgsForCall(0)
							Expect(spec).To(Equal(garden.ProcessSpec{
								Path: "some-path",
								Args: []string{"some", "args"},
								Env:  []string{"some=env"},
								Dir:  "some-dir",
								User: "maverick",
							}))
							Expect(io).To(Equal(garden.ProcessIO{}))
						})
					})
				})

				Context("when the user property is not present", func() {
					var (
						actualSpec garden.ProcessSpec
						actualIO   garden.ProcessIO
					)

					BeforeEach(func() {
						actualSpec = garden.ProcessSpec{
							Path: "some-path",
							Args: []string{"some", "args"},
							Env:  []string{"some=env"},
							Dir:  "some-dir",
						}

						actualIO = garden.ProcessIO{}

						fakeContainer.PropertiesReturns(garden.Properties{"user": ""}, nil)
					})

					JustBeforeEach(func() {
						foundContainer.Run(actualSpec, actualIO)
					})

					Describe("Run", func() {
						It("calls Run() on the garden container and injects the default user", func() {
							Expect(fakeContainer.RunCallCount()).To(Equal(1))
							spec, io := fakeContainer.RunArgsForCall(0)
							Expect(spec).To(Equal(garden.ProcessSpec{
								Path: "some-path",
								Args: []string{"some", "args"},
								Env:  []string{"some=env"},
								Dir:  "some-dir",
								User: "root",
							}))
							Expect(io).To(Equal(garden.ProcessIO{}))
							Expect(fakeContainer.RunCallCount()).To(Equal(1))
						})
					})
				})
			})
		})

		Context("when the gardenClient returns garden.ContainerNotFoundError", func() {
			BeforeEach(func() {
				fakeGardenClient.LookupReturns(nil, garden.ContainerNotFoundError{Handle: "some-handle"})
			})

			It("returns false and no error", func() {
				_, found, err := gardenWorker.LookupContainer(logger, handle)
				Expect(err).ToNot(HaveOccurred())

				Expect(found).To(BeFalse())
			})
		})

		Context("when the gardenClient returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = fmt.Errorf("container not found")
				fakeGardenClient.LookupReturns(nil, expectedErr)
			})

			It("returns nil and forwards the error", func() {
				foundContainer, _, err := gardenWorker.LookupContainer(logger, handle)
				Expect(err).To(Equal(expectedErr))

				Expect(foundContainer).To(BeNil())
			})
		})
	})

	Describe("FindContainerForIdentifier", func() {
		var (
			id Identifier

			foundContainer Container
			found          bool
			lookupErr      error
		)

		BeforeEach(func() {
			id = Identifier{
				ResourceID: 1234,
			}
		})

		JustBeforeEach(func() {
			foundContainer, found, lookupErr = gardenWorker.FindContainerForIdentifier(logger, id)
		})

		Context("when the container can be found", func() {
			var (
				fakeContainer *gfakes.FakeContainer
			)

			BeforeEach(func() {
				fakeContainer = new(gfakes.FakeContainer)
				fakeContainer.HandleReturns("provider-handle")

				fakeWorkerProvider.FindContainerForIdentifierReturns(db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							Handle: "provider-handle",
						},
					},
				}, true, nil)

				fakeGardenClient.LookupReturns(fakeContainer, nil)
			})

			It("succeeds", func() {
				Expect(lookupErr).NotTo(HaveOccurred())
			})

			It("looks for containers with matching properties via the Garden client", func() {
				Expect(fakeWorkerProvider.FindContainerForIdentifierCallCount()).To(Equal(1))
				Expect(fakeWorkerProvider.FindContainerForIdentifierArgsForCall(0)).To(Equal(id))

				Expect(fakeGardenClient.LookupCallCount()).To(Equal(1))
				lookupHandle := fakeGardenClient.LookupArgsForCall(0)

				Expect(lookupHandle).To(Equal("provider-handle"))
			})

			Describe("the found container", func() {
				It("can be destroyed", func() {
					err := foundContainer.Destroy()
					Expect(err).NotTo(HaveOccurred())

					By("destroying via garden")
					Expect(fakeGardenClient.DestroyCallCount()).To(Equal(1))
					Expect(fakeGardenClient.DestroyArgsForCall(0)).To(Equal("provider-handle"))

					By("no longer heartbeating")
					fakeClock.Increment(30 * time.Second)
					Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(1))
				})

				It("performs an initial heartbeat synchronously", func() {
					Expect(fakeContainer.SetGraceTimeCallCount()).To(Equal(1))
					Expect(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount()).To(Equal(1))
				})

				Describe("every 30 seconds", func() {
					It("heartbeats to the database and the container", func() {
						fakeClock.Increment(30 * time.Second)

						Eventually(fakeContainer.SetGraceTimeCallCount).Should(Equal(2))
						Expect(fakeContainer.SetGraceTimeArgsForCall(1)).To(Equal(5 * time.Minute))

						Eventually(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(2))
						handle, interval := fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(1)
						Expect(handle).To(Equal("provider-handle"))
						Expect(interval).To(Equal(5 * time.Minute))

						fakeClock.Increment(30 * time.Second)

						Eventually(fakeContainer.SetGraceTimeCallCount).Should(Equal(3))
						Expect(fakeContainer.SetGraceTimeArgsForCall(2)).To(Equal(5 * time.Minute))

						Eventually(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(3))
						handle, interval = fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(2)
						Expect(handle).To(Equal("provider-handle"))
						Expect(interval).To(Equal(5 * time.Minute))
					})
				})

				Describe("releasing", func() {
					It("sets a final ttl on the container and stops heartbeating", func() {
						foundContainer.Release(FinalTTL(30 * time.Minute))

						Expect(fakeContainer.SetGraceTimeCallCount()).Should(Equal(2))
						Expect(fakeContainer.SetGraceTimeArgsForCall(1)).To(Equal(30 * time.Minute))

						Expect(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount()).Should(Equal(2))
						handle, interval := fakeGardenWorkerDB.UpdateExpiresAtOnContainerArgsForCall(1)
						Expect(handle).To(Equal("provider-handle"))
						Expect(interval).To(Equal(30 * time.Minute))

						fakeClock.Increment(30 * time.Second)

						Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(2))
						Consistently(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(2))
					})

					Context("with no final ttl", func() {
						It("does not perform a final heartbeat", func() {
							foundContainer.Release(nil)

							Consistently(fakeContainer.SetGraceTimeCallCount).Should(Equal(1))
							Consistently(fakeGardenWorkerDB.UpdateExpiresAtOnContainerCallCount).Should(Equal(1))
						})
					})
				})

				It("can be released multiple times", func() {
					foundContainer.Release(nil)

					Expect(func() {
						foundContainer.Release(nil)
					}).NotTo(Panic())
				})
			})
		})

		Context("when no containers are found", func() {
			BeforeEach(func() {
				fakeWorkerProvider.FindContainerForIdentifierReturns(db.SavedContainer{}, false, nil)
			})

			It("returns that the container could not be found", func() {
				Expect(found).To(BeFalse())
			})
		})

		Context("when finding the containers fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				fakeWorkerProvider.FindContainerForIdentifierReturns(db.SavedContainer{}, false, disaster)
			})

			It("returns the error", func() {
				Expect(lookupErr).To(Equal(disaster))
			})
		})

		Context("when the container cannot be found", func() {
			BeforeEach(func() {
				containerToReturn := db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							Handle: "handle",
						},
					},
				}

				fakeWorkerProvider.FindContainerForIdentifierReturns(containerToReturn, true, nil)
				fakeGardenClient.LookupReturns(nil, garden.ContainerNotFoundError{Handle: "handle"})
			})

			It("expires the container and returns false and no error", func() {
				Expect(lookupErr).ToNot(HaveOccurred())
				Expect(found).To(BeFalse())
				Expect(foundContainer).To(BeNil())

				expiredHandle := fakeWorkerProvider.ReapContainerArgsForCall(0)
				Expect(expiredHandle).To(Equal("handle"))
			})
		})

		Context("when looking up the container fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				containerToReturn := db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							Handle: "handle",
						},
					},
				}

				fakeWorkerProvider.FindContainerForIdentifierReturns(containerToReturn, true, nil)
				fakeGardenClient.LookupReturns(nil, disaster)
			})

			It("returns the error", func() {
				Expect(lookupErr).To(Equal(disaster))
			})
		})
	})

	Describe("Satisfying", func() {
		var (
			spec WorkerSpec

			satisfyingWorker Worker
			satisfyingErr    error

			customTypes atc.ResourceTypes
		)

		BeforeEach(func() {
			spec = WorkerSpec{}
			customTypes = atc.ResourceTypes{
				{
					Name:   "custom-type-b",
					Type:   "custom-type-a",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-a",
					Type:   "some-resource",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-c",
					Type:   "custom-type-b",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "custom-type-d",
					Type:   "custom-type-b",
					Source: atc.Source{"some": "source"},
				},
				{
					Name:   "unknown-custom-type",
					Type:   "unknown-base-type",
					Source: atc.Source{"some": "source"},
				},
			}
		})

		JustBeforeEach(func() {
			gardenWorker = NewGardenWorker(
				fakeGardenClient,
				fakeBaggageclaimClient,
				fakeVolumeFactory,
				fakeImageFetcher,
				fakeGardenWorkerDB,
				fakeWorkerProvider,
				fakeClock,
				activeContainers,
				resourceTypes,
				platform,
				tags,
				workerName,
				httpProxyURL,
				httpsProxyURL,
				noProxy,
			)

			satisfyingWorker, satisfyingErr = gardenWorker.Satisfying(spec, customTypes)
		})

		Context("when the platform is compatible", func() {
			BeforeEach(func() {
				spec.Platform = "some-platform"
			})

			Context("when no tags are specified", func() {
				BeforeEach(func() {
					spec.Tags = nil
				})

				It("returns ErrIncompatiblePlatform", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})

			Context("when the worker has no tags", func() {
				BeforeEach(func() {
					tags = []string{}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when all of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some", "tags"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when some of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when any of the requested tags are not present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"bogus", "tags"}
				})

				It("returns ErrMismatchedTags", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})
		})

		Context("when the platform is incompatible", func() {
			BeforeEach(func() {
				spec.Platform = "some-bogus-platform"
			})

			It("returns ErrIncompatiblePlatform", func() {
				Expect(satisfyingErr).To(Equal(ErrIncompatiblePlatform))
			})
		})

		Context("when the resource type is supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "some-resource"
			})

			Context("when all of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some", "tags"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when some of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when any of the requested tags are not present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"bogus", "tags"}
				})

				It("returns ErrMismatchedTags", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})
		})

		Context("when the resource type is a custom type supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "custom-type-c"
				spec.Tags = []string{"some", "tags"}
			})

			It("returns the worker", func() {
				Expect(satisfyingWorker).To(Equal(gardenWorker))
			})

			It("returns no error", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})
		})

		Context("when the resource type is a custom type that overrides one supported by the worker", func() {
			BeforeEach(func() {
				customTypes = append(customTypes, atc.ResourceType{
					Name:   "some-resource",
					Type:   "some-resource",
					Source: atc.Source{"some": "source"},
				})

				spec.ResourceType = "some-resource"
				spec.Tags = []string{"some", "tags"}
			})

			It("returns the worker", func() {
				Expect(satisfyingWorker).To(Equal(gardenWorker))
			})

			It("returns no error", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})
		})

		Context("when the resource type is a custom type that results in a circular dependency", func() {
			BeforeEach(func() {
				customTypes = append(customTypes, atc.ResourceType{
					Name:   "circle-a",
					Type:   "circle-b",
					Source: atc.Source{"some": "source"},
				}, atc.ResourceType{
					Name:   "circle-b",
					Type:   "circle-c",
					Source: atc.Source{"some": "source"},
				}, atc.ResourceType{
					Name:   "circle-c",
					Type:   "circle-a",
					Source: atc.Source{"some": "source"},
				})

				spec.ResourceType = "circle-a"
				spec.Tags = []string{"some", "tags"}
			})

			It("returns ErrUnsupportedResourceType", func() {
				Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
			})
		})

		Context("when the resource type is a custom type not supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "unknown-custom-type"
				spec.Tags = []string{"some", "tags"}
			})

			It("returns ErrUnsupportedResourceType", func() {
				Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
			})
		})

		Context("when the type is not supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "some-other-resource"
			})

			Context("when all of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some", "tags"}
				})

				It("returns ErrUnsupportedResourceType", func() {
					Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
				})
			})

			Context("when some of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some"}
				})

				It("returns ErrUnsupportedResourceType", func() {
					Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
				})
			})

			Context("when any of the requested tags are not present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"bogus", "tags"}
				})

				It("returns ErrUnsupportedResourceType", func() {
					Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
				})
			})
		})
	})
})
