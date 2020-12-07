package engine

import (
	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/atc/event"
	"github.com/pf-qiu/concourse/v6/atc/exec"
)

func NewSetPipelineStepDelegate(
	build db.Build,
	planID atc.PlanID,
	state exec.RunState,
	clock clock.Clock,
) *setPipelineStepDelegate {
	return &setPipelineStepDelegate{
		buildStepDelegate{
			build:  build,
			planID: planID,
			clock:  clock,
			state:  state,
			stdout: nil,
			stderr: nil,
		},
	}
}

type setPipelineStepDelegate struct {
	buildStepDelegate
}

func (delegate *setPipelineStepDelegate) SetPipelineChanged(logger lager.Logger, changed bool) {
	err := delegate.build.SaveEvent(event.SetPipelineChanged{
		Origin: event.Origin{
			ID: event.OriginID(delegate.planID),
		},
		Changed: changed,
	})
	if err != nil {
		logger.Error("failed-to-save-set-pipeline-changed-event", err)
		return
	}

	logger.Debug("set pipeline changed")
}
