package db_test

import (
	"context"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/time/rate"
)

var _ = Describe("ResourceCheckRateLimiter", func() {
	var (
		checkInterval   time.Duration
		checksPerSecond int
		refreshInterval time.Duration
		fakeClock       *fakeclock.FakeClock

		checkableCount int

		ctx context.Context

		limiter *db.ResourceCheckRateLimiter
	)

	BeforeEach(func() {
		checkInterval = time.Minute
		checksPerSecond = 0
		refreshInterval = 5 * time.Minute
		fakeClock = fakeclock.NewFakeClock(time.Now())

		checkableCount = 0

		ctx = context.Background()
	})

	JustBeforeEach(func() {
		limiter = db.NewResourceCheckRateLimiter(
			rate.Limit(checksPerSecond),
			checkInterval,
			dbConn,
			refreshInterval,
			fakeClock,
		)
	})

	wait := func(limiter *db.ResourceCheckRateLimiter) <-chan error {
		errs := make(chan error)
		go func() {
			errs <- limiter.Wait(ctx)
		}()
		return errs
	}

	createCheckable := func() {
		config, err := resourceConfigFactory.FindOrCreateResourceConfig(
			defaultWorkerResourceType.Type,
			atc.Source{"some": "source", "count": checkableCount},
			atc.VersionedResourceTypes{},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = config.FindOrCreateScope(nil)
		Expect(err).ToNot(HaveOccurred())

		checkableCount++
	}

	Context("with no static limit provided", func() {
		BeforeEach(func() {
			checksPerSecond = 0
		})

		It("rate limits while adjusting to the amount of checkables", func() {
			By("immediately returning with 0 checkables")
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Inf))

			By("creating one checkable")
			createCheckable()

			By("continuing to return immediately, as the refresh interval has not elapsed")
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Inf))

			By("waiting for the refresh interval")
			fakeClock.Increment(refreshInterval)

			By("adjusting the limit but returning immediately for the first time")
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Limit(float64(checkableCount) / checkInterval.Seconds())))

			done := wait(limiter)
			select {
			case <-done:
				Fail("should not have returned yet")
			case <-time.After(100 * time.Millisecond):
			}

			By("unblocking after the rate limit elapses")
			fakeClock.Increment(checkInterval / time.Duration(checkableCount))
			Expect(<-done).To(Succeed())

			By("creating more checkables")
			for i := 0; i < 10; i++ {
				createCheckable()
			}

			By("waiting for the refresh interval")
			fakeClock.Increment(refreshInterval)

			By("adjusting the limit but returning immediately for the first time")
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Limit(float64(checkableCount) / checkInterval.Seconds())))

			done = wait(limiter)
			select {
			case <-done:
				Fail("should not have returned yet")
			case <-time.After(100 * time.Millisecond):
			}

			By("unblocking after the the new rate limit elapses")
			fakeClock.Increment(checkInterval / time.Duration(checkableCount))
			Expect(<-done).To(Succeed())
		})
	})

	Context("when a static checks per second value is provided", func() {
		BeforeEach(func() {
			checksPerSecond = 42
		})

		It("respects the value rather than determining it dynamically", func() {
			Expect(limiter.Limit()).To(Equal(rate.Limit(checksPerSecond)))
		})
	})

	Context("when a negative static checks per second value is provided", func() {
		BeforeEach(func() {
			checksPerSecond = -1
		})

		It("results in an infinite rate limit that ignores checkable count", func() {
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Limit(rate.Inf)))

			By("creating a few (ignored) checkables")
			for i := 0; i < 10; i++ {
				createCheckable()
			}

			By("waiting for the (ignored) refresh interval")
			fakeClock.Increment(refreshInterval)

			By("still returning immediately and retaining the infinite rate")
			Expect(<-wait(limiter)).To(Succeed())
			Expect(limiter.Limit()).To(Equal(rate.Limit(rate.Inf)))
		})
	})
})
