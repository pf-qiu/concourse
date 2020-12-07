package commands

import (
	"errors"
	"fmt"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse"
	"sort"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
)

var ErrMissingPipelineName = errors.New("Need to specify atleast one pipeline name")

type OrderPipelinesCommand struct {
	Alphabetical bool                       `short:"a"  long:"alphabetical" description:"Order all pipelines alphabetically"`
	Pipelines    []flaghelpers.PipelineFlag `short:"p" long:"pipeline" description:"Name of pipeline to order"`
	Team         string                     `long:"team" description:"Name of the team to which the pipelines belong, if different from the target default"`
}

func (command *OrderPipelinesCommand) Validate() (atc.OrderPipelinesRequest, error) {
	var pipelineRefs atc.OrderPipelinesRequest

	for _, p := range command.Pipelines {
		_, err := p.Validate()
		if err != nil {
			return nil, err
		}
		pipelineRefs = append(pipelineRefs, p.Ref())
	}
	return pipelineRefs, nil
}

func (command *OrderPipelinesCommand) Execute(args []string) error {
	var pipelineRefs atc.OrderPipelinesRequest

	if !command.Alphabetical && command.Pipelines == nil {
		displayhelpers.Failf("error: either --pipeline or --alphabetical are required")
	}

	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	if command.Alphabetical {
		ps, err := target.Team().ListPipelines()
		if err != nil {
			return err
		}

		for _, p := range ps {
			pipelineRefs = append(pipelineRefs, p.Ref())
		}
		sort.Sort(pipelineRefs)
	} else {
		pipelineRefs, err = command.Validate()
		if err != nil {
			return err
		}
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

	err = team.OrderingPipelines(pipelineRefs)
	if err != nil {
		displayhelpers.FailWithErrorf("failed to order pipelines", err)
	}

	fmt.Printf("ordered pipelines \n")
	for _, p := range pipelineRefs {
		fmt.Printf("  - %s \n", p.String())
	}

	return nil
}
