package worker

import (
	"context"
	"errors"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/db/lock"
	"github.com/pf-qiu/concourse/v6/atc/resource"
	"github.com/pf-qiu/concourse/v6/atc/runtime"
)

const GetResourceLockInterval = 5 * time.Second

var ErrFailedToGetLock = errors.New("failed to get lock")

//go:generate counterfeiter . Fetcher

type Fetcher interface {
	Fetch(
		ctx context.Context,
		logger lager.Logger,
		containerMetadata db.ContainerMetadata,
		gardenWorker Worker,
		containerSpec ContainerSpec,
		processSpec runtime.ProcessSpec,
		resource resource.Resource,
		owner db.ContainerOwner,
		cache db.UsedResourceCache,
		lockName string,
	) (GetResult, Volume, error)
}

func NewFetcher(
	clock clock.Clock,
	lockFactory lock.LockFactory,
	fetchSourceFactory FetchSourceFactory,
) Fetcher {
	return &fetcher{
		clock:              clock,
		lockFactory:        lockFactory,
		fetchSourceFactory: fetchSourceFactory,
	}
}

type fetcher struct {
	clock              clock.Clock
	lockFactory        lock.LockFactory
	fetchSourceFactory FetchSourceFactory
}

func (f *fetcher) Fetch(
	ctx context.Context,
	logger lager.Logger,
	containerMetadata db.ContainerMetadata,
	gardenWorker Worker,
	containerSpec ContainerSpec,
	processSpec runtime.ProcessSpec,
	resource resource.Resource,
	owner db.ContainerOwner,
	cache db.UsedResourceCache,
	lockName string,
) (GetResult, Volume, error) {
	result := GetResult{}
	var volume Volume
	// TODO: resource_instance_fetch_source.go already knows which volume to use for the resource output, can this be consolidated
	containerSpec.Outputs = map[string]string{
		"resource": processSpec.Args[0],
	}

	// TODO: just pass in imageFetcherSpec not its contents
	//		 this might be a bad idea, don't want the fetchsource to know about images?
	fetchSource := f.fetchSourceFactory.NewFetchSource(
		logger,
		gardenWorker,
		owner,
		cache,
		resource,
		containerSpec,
		processSpec,
		containerMetadata,
	)

	ticker := f.clock.NewTicker(GetResourceLockInterval)
	defer ticker.Stop()

	result, volume, err := f.fetchUnderLock(
		ctx,
		logger,
		fetchSource,
		cache,
		lockName,
	)
	if err == nil || err != ErrFailedToGetLock {
		return result, volume, err
	}

	for {
		select {
		case <-ticker.C():
			result, volume, err = f.fetchUnderLock(
				ctx,
				logger,
				fetchSource,
				cache,
				lockName,
			)
			if err != nil {
				if err == ErrFailedToGetLock {
					break
				}
				return result, nil, err
			}

			return result, volume, nil

		case <-ctx.Done():
			return GetResult{}, nil, ctx.Err()
		}
	}
}

func (f *fetcher) fetchUnderLock(
	ctx context.Context,
	logger lager.Logger,
	source FetchSource,
	cache db.UsedResourceCache,
	lockName string,
) (GetResult, Volume, error) {
	result := GetResult{}
	findResult, volume, found, err := source.Find()
	if err != nil {
		return result, nil, err
	}

	if found {
		return findResult, volume, nil
	}

	lockLogger := logger.Session("lock-task", lager.Data{"lock-name": lockName})

	lock, acquired, err := f.lockFactory.Acquire(lockLogger, lock.NewTaskLockID(lockName))
	if err != nil {
		lockLogger.Error("failed-to-get-lock", err)
		return result, nil, ErrFailedToGetLock
	}

	if !acquired {
		lockLogger.Debug("did-not-get-lock")
		return result, nil, ErrFailedToGetLock
	}

	defer lock.Release()

	return source.Create(ctx)
}
