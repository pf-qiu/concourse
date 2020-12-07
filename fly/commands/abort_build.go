package commands

import (
	"fmt"
	"strconv"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
)

type AbortBuildCommand struct {
	Job   flaghelpers.JobFlag `short:"j" long:"job" value-name:"PIPELINE/JOB"   description:"Name of a job to cancel"`
	Build string              `short:"b" long:"build" required:"true" description:"If job is specified: build number to cancel. If job not specified: build id"`
}

func (command *AbortBuildCommand) Execute([]string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	var build atc.Build
	var exists bool
	if command.Job.PipelineRef.Name == "" && command.Job.JobName == "" {
		build, exists, err = target.Client().Build(command.Build)
	} else {
		build, exists, err = target.Team().JobBuild(command.Job.PipelineRef, command.Job.JobName, command.Build)
	}
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("build does not exist")
	}

	if err := target.Client().AbortBuild(strconv.Itoa(build.ID)); err != nil {
		return err
	}

	fmt.Println("build successfully aborted")
	return nil
}
