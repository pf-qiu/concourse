package gc

import (
	"context"
	"time"

	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/metric"
)

type resourceCacheCollector struct {
	cacheLifecycle db.ResourceCacheLifecycle
}

func NewResourceCacheCollector(cacheLifecycle db.ResourceCacheLifecycle) *resourceCacheCollector {
	return &resourceCacheCollector{
		cacheLifecycle: cacheLifecycle,
	}
}

func (rcc *resourceCacheCollector) Run(ctx context.Context) error {
	logger := lagerctx.FromContext(ctx).Session("resource-cache-collector")

	logger.Debug("start")
	defer logger.Debug("done")

	start := time.Now()
	defer func() {
		metric.ResourceCacheCollectorDuration{
			Duration: time.Since(start),
		}.Emit(logger)
	}()

	return rcc.cacheLifecycle.CleanUpInvalidCaches(logger)
}
