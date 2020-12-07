package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
)

type DisableResourceVersionCommand struct {
	Resource flaghelpers.ResourceFlag `short:"r" long:"resource" required:"true" value-name:"PIPELINE/RESOURCE" description:"Name of the resource"`
	Version  *atc.Version             `short:"v" long:"version" required:"true" value-name:"KEY:VALUE" description:"Version of the resource to disable. The given key value pair(s) has to be an exact match but not all fields are needed. In the case of multiple resource versions matched, it will disable the latest one."`
}

func (command *DisableResourceVersionCommand) Execute([]string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	team := target.Team()

	if command.Version != nil {
		latestResourceVer, err := GetLatestResourceVersion(team, command.Resource, *command.Version)
		if err != nil {
			return err
		}

		disabled := !latestResourceVer.Enabled

		if !disabled {
			disabled, err = team.DisableResourceVersion(command.Resource.PipelineRef, command.Resource.ResourceName, latestResourceVer.ID)
			if err != nil {
				return err
			}
		}

		if disabled {
			disableVersionBytes, err := json.Marshal(latestResourceVer.Version)
			if err != nil {
				return err
			}

			fmt.Printf("disabled '%s/%s' with version %s\n", command.Resource.PipelineRef.String(), command.Resource.ResourceName, string(disableVersionBytes))
		} else {
			displayhelpers.Failf("could not disable '%s/%s', make sure the resource version exists\n", command.Resource.PipelineRef.String(), command.Resource.ResourceName)
		}
	}

	return nil
}
