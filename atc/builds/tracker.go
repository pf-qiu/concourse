package builds

import (
	"context"
	"sync"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/engine"
	"github.com/pf-qiu/concourse/v6/atc/metric"
	"github.com/pf-qiu/concourse/v6/atc/util"
)

func NewTracker(
	buildFactory db.BuildFactory,
	engine engine.Engine,
) *Tracker {
	return &Tracker{
		buildFactory: buildFactory,
		engine:       engine,
		running:      &sync.Map{},
	}
}

type Tracker struct {
	logger lager.Logger

	buildFactory db.BuildFactory
	engine       engine.Engine

	running *sync.Map
}

func (bt *Tracker) Run(ctx context.Context) error {
	logger := lagerctx.FromContext(ctx)

	logger.Debug("start")
	defer logger.Debug("done")

	builds, err := bt.buildFactory.GetAllStartedBuilds()
	if err != nil {
		logger.Error("failed-to-lookup-started-builds", err)
		return err
	}

	for _, b := range builds {
		if _, exists := bt.running.LoadOrStore(b.ID(), true); !exists {
			go func(build db.Build) {
				loggerData := build.LagerData()
				defer func() {
					err := util.DumpPanic(recover(), "tracking build %d", build.ID())
					if err != nil {
						logger.Error("panic-in-tracker-build-run", err)

						build.Finish(db.BuildStatusErrored)
					}
				}()

				defer bt.running.Delete(build.ID())

				metric.Metrics.BuildsRunning.Inc()
				defer metric.Metrics.BuildsRunning.Dec()

				bt.engine.NewBuild(build).Run(
					lagerctx.NewContext(
						context.Background(),
						logger.Session("run", loggerData),
					),
				)
			}(b)
		}
	}

	return nil
}

func (bt *Tracker) Drain(ctx context.Context) {
	bt.engine.Drain(ctx)
}
