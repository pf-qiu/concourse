package gc

import (
	"context"
	"time"

	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/metric"
	multierror "github.com/hashicorp/go-multierror"
)

type resourceConfigCheckSessionCollector struct {
	configCheckSessionLifecycle db.ResourceConfigCheckSessionLifecycle
}

func NewResourceConfigCheckSessionCollector(
	configCheckSessionLifecycle db.ResourceConfigCheckSessionLifecycle,
) *resourceConfigCheckSessionCollector {
	return &resourceConfigCheckSessionCollector{
		configCheckSessionLifecycle: configCheckSessionLifecycle,
	}
}

func (rccsc *resourceConfigCheckSessionCollector) Run(ctx context.Context) error {
	logger := lagerctx.FromContext(ctx).Session("resource-config-check-session-collector")

	logger.Debug("start")
	defer logger.Debug("done")

	start := time.Now()
	defer func() {
		metric.ResourceConfigCheckSessionCollectorDuration{
			Duration: time.Since(start),
		}.Emit(logger)
	}()

	var errs error

	err := rccsc.configCheckSessionLifecycle.CleanExpiredResourceConfigCheckSessions()
	if err != nil {
		errs = multierror.Append(errs, err)
		logger.Error("failed-to-clean-up-expired-resource-config-check-sessions", err)
	}

	err = rccsc.configCheckSessionLifecycle.CleanInactiveResourceConfigCheckSessions()
	if err != nil {
		errs = multierror.Append(errs, err)
		logger.Error("failed-to-clean-up-inactive-resource-config-check-sessions", err)
	}

	return errs
}
