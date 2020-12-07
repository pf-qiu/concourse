package commands

import (
	"fmt"

	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse"
)

type HidePipelineCommand struct {
	Pipeline flaghelpers.PipelineFlag `short:"p"   long:"pipeline" required:"true" description:"Pipeline to hide"`
	Team     string                   `long:"team"                                 description:"Name of the team to which the pipeline belongs, if different from the target default"`
}

func (command *HidePipelineCommand) Validate() error {
	_, err := command.Pipeline.Validate()
	return err
}

func (command *HidePipelineCommand) Execute(args []string) error {
	err := command.Validate()
	if err != nil {
		return err
	}

	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	var team concourse.Team
	if command.Team != "" {
		team, err = target.FindTeam(command.Team)
		if err != nil {
			return err
		}
	} else {
		team = target.Team()
	}

	pipelineRef := command.Pipeline.Ref()
	found, err := team.HidePipeline(pipelineRef)
	if err != nil {
		return err
	}

	if found {
		fmt.Printf("hid '%s'\n", pipelineRef.String())
	} else {
		displayhelpers.Failf("pipeline '%s' not found\n", pipelineRef.String())
	}

	return nil
}
