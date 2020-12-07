package gc

import (
	"context"
	"time"

	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/metric"
)

type artifactCollector struct {
	artifactLifecycle db.WorkerArtifactLifecycle
}

func NewArtifactCollector(artifactLifecycle db.WorkerArtifactLifecycle) *artifactCollector {
	return &artifactCollector{
		artifactLifecycle: artifactLifecycle,
	}
}

func (a *artifactCollector) Run(ctx context.Context) error {
	logger := lagerctx.FromContext(ctx).Session("artifact-collector")

	logger.Debug("start")
	defer logger.Debug("done")

	start := time.Now()
	defer func() {
		metric.ArtifactCollectorDuration{
			Duration: time.Since(start),
		}.Emit(logger)
	}()

	return a.artifactLifecycle.RemoveExpiredArtifacts()
}
