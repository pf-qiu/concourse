package worker_test

import (
	"context"
	"errors"
	"time"

	"github.com/pf-qiu/concourse/v6/atc/db/dbfakes"
	"github.com/pf-qiu/concourse/v6/atc/resource"
	"github.com/pf-qiu/concourse/v6/atc/runtime"
	"github.com/pf-qiu/concourse/v6/atc/worker"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/db/lock"
	"github.com/pf-qiu/concourse/v6/atc/db/lock/lockfakes"
	"github.com/pf-qiu/concourse/v6/atc/resource/resourcefakes"
	"github.com/pf-qiu/concourse/v6/atc/worker/workerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher", func() {
	var (
		fakeClock       *fakeclock.FakeClock
		fakeLockFactory *lockfakes.FakeLockFactory
		fetcher         worker.Fetcher
		ctx             context.Context
		cancel          func()

		fakeWorker             *workerfakes.FakeWorker
		fakeVolume             *workerfakes.FakeVolume
		fakeFetchSourceFactory *workerfakes.FakeFetchSourceFactory
		fakeResource           *resourcefakes.FakeResource
		fakeUsedResourceCache  *dbfakes.FakeUsedResourceCache

		getResult worker.GetResult
		fetchErr  error
		teamID    = 123

		volume worker.Volume
	)

	BeforeEach(func() {
		fakeClock = fakeclock.NewFakeClock(time.Unix(0, 123))
		fakeLockFactory = new(lockfakes.FakeLockFactory)
		fakeFetchSourceFactory = new(workerfakes.FakeFetchSourceFactory)

		fakeWorker = new(workerfakes.FakeWorker)
		fakeWorker.NameReturns("some-worker")

		fakeVolume = new(workerfakes.FakeVolume)
		fakeVolume.HandleReturns("some-handle")
		fakeResource = new(resourcefakes.FakeResource)
		fakeUsedResourceCache = new(dbfakes.FakeUsedResourceCache)

		fetcher = worker.NewFetcher(
			fakeClock,
			fakeLockFactory,
			fakeFetchSourceFactory,
		)

		ctx, cancel = context.WithCancel(context.Background())
	})

	JustBeforeEach(func() {
		getResult, volume, fetchErr = fetcher.Fetch(
			ctx,
			lagertest.NewTestLogger("test"),
			db.ContainerMetadata{},
			fakeWorker,
			worker.ContainerSpec{
				TeamID: teamID,
			},
			runtime.ProcessSpec{
				Args: []string{resource.ResourcesDir("get")},
			},
			fakeResource,
			db.NewBuildStepContainerOwner(0, "some-plan-id", 0),
			fakeUsedResourceCache,
			"fake-lock-name",
		)
	})

	Context("when getting source", func() {
		var fakeFetchSource *workerfakes.FakeFetchSource

		BeforeEach(func() {
			fakeFetchSource = new(workerfakes.FakeFetchSource)
			fakeFetchSourceFactory.NewFetchSourceReturns(fakeFetchSource)

			fakeFetchSource.FindReturns(worker.GetResult{}, fakeVolume, false, nil)
		})

		Describe("failing to get a lock", func() {
			Context("when did not get a lock", func() {
				BeforeEach(func() {
					fakeLock := new(lockfakes.FakeLock)
					callCount := 0
					fakeLockFactory.AcquireStub = func(lager.Logger, lock.LockID) (lock.Lock, bool, error) {
						callCount++
						fakeClock.Increment(worker.GetResourceLockInterval)
						if callCount == 1 {
							return nil, false, nil
						}
						return fakeLock, true, nil
					}
				})

				It("retries until it gets the lock", func() {
					Expect(fakeLockFactory.AcquireCallCount()).To(Equal(2))
				})

				It("creates fetch source after lock is acquired", func() {
					Expect(fakeFetchSource.CreateCallCount()).To(Equal(1))
				})
			})

			Context("when acquiring lock returns error", func() {
				BeforeEach(func() {
					fakeLock := new(lockfakes.FakeLock)
					callCount := 0
					fakeLockFactory.AcquireStub = func(lager.Logger, lock.LockID) (lock.Lock, bool, error) {
						callCount++
						fakeClock.Increment(worker.GetResourceLockInterval)
						if callCount == 1 {
							return nil, false, errors.New("disaster")
						}
						return fakeLock, true, nil
					}
				})

				It("retries until it gets the lock", func() {
					Expect(fakeLockFactory.AcquireCallCount()).To(Equal(2))
				})

				It("creates fetch source after lock is acquired", func() {
					Expect(fakeFetchSource.CreateCallCount()).To(Equal(1))
				})
			})
		})

		Context("when getting lock succeeds", func() {
			var fakeLock *lockfakes.FakeLock

			BeforeEach(func() {
				fakeLock = new(lockfakes.FakeLock)
				fakeLockFactory.AcquireReturns(fakeLock, true, nil)
				fakeFetchSource.CreateReturns(worker.GetResult{}, fakeVolume, nil)
			})

			It("acquires a lock with source lock name", func() {
				Expect(fakeLockFactory.AcquireCallCount()).To(Equal(1))
				_, lockID := fakeLockFactory.AcquireArgsForCall(0)
				Expect(lockID).To(Equal(lock.NewTaskLockID("fake-lock-name")))
			})

			It("releases the lock", func() {
				Expect(fakeLock.ReleaseCallCount()).To(Equal(1))
			})

			It("creates source", func() {
				Expect(fakeFetchSource.CreateCallCount()).To(Equal(1))
			})

			It("returns the result and volume", func() {
				Expect(getResult).To(Equal(worker.GetResult{}))
				Expect(volume).ToNot(BeNil())
			})
		})

		Context("when finding fails", func() {
			var disaster error

			BeforeEach(func() {
				disaster = errors.New("fail")
				fakeFetchSource.FindReturns(worker.GetResult{}, nil, false, disaster)
			})

			It("returns an error", func() {
				Expect(fetchErr).To(HaveOccurred())
				Expect(fetchErr).To(Equal(disaster))
			})
		})

		Context("when cancelled", func() {
			BeforeEach(func() {
				cancel()
			})

			It("returns the context err", func() {
				Expect(fetchErr).To(Equal(context.Canceled))
			})
		})
	})
})
