package commands

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
	"github.com/jessevdk/go-flags"
)

type PinResourceCommand struct {
	Resource flaghelpers.ResourceFlag `short:"r" long:"resource" required:"true" value-name:"PIPELINE/RESOURCE" description:"Name of the resource"`
	Version  *atc.Version             `short:"v" long:"version" description:"Version of the resource to pin. The given key value pair(s) has to be an exact match but not all fields are needed. In the case of multiple resource versions matched, it will pin the latest one."`
	Comment  string                   `short:"c" long:"comment" description:"Message to be saved to the pinned resource. Resource has to be pinned otherwise --version should be specified to pin the resource first."`
}

func (command *PinResourceCommand) Execute([]string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	team := target.Team()

	pipelineRef := command.Resource.PipelineRef

	if command.Version == nil && command.Comment == "" {
		shortDelim := "-"
		longDelim := "--"
		if runtime.GOOS == "windows" {
			shortDelim = "/"
			longDelim = "/"
		}
		return &flags.Error{
			Type: flags.ErrRequired,
			Message: fmt.Sprintf(
				"the required flag `%sv, %sversion' was not specified",
				shortDelim,
				longDelim,
			),
		}
	}

	if command.Version != nil {
		latestResourceVersion, err := GetLatestResourceVersion(team, command.Resource, *command.Version)
		if err != nil {
			return err
		}

		pinned, err := team.PinResourceVersion(pipelineRef, command.Resource.ResourceName, latestResourceVersion.ID)

		if err != nil {
			return err
		}

		if pinned {
			versionBytes, err := json.Marshal(latestResourceVersion.Version)
			if err != nil {
				return err
			}

			fmt.Printf("pinned '%s/%s' with version %s\n", pipelineRef.String(), command.Resource.ResourceName, string(versionBytes))
		} else {
			displayhelpers.Failf("could not pin '%s/%s', make sure the resource exists\n", pipelineRef.String(), command.Resource.ResourceName)
		}
	}

	if command.Comment != "" {
		saved, err := team.SetPinComment(pipelineRef, command.Resource.ResourceName, command.Comment)

		if err != nil {
			return err
		}

		if saved {
			fmt.Printf("pin comment '%s' is saved\n", command.Comment)
		} else {
			displayhelpers.Failf("could not save comment, make sure '%s/%s' is pinned\n", pipelineRef.String(), command.Resource.ResourceName)
		}
	}

	return nil
}
