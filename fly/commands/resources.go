package commands

import (
	"os"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
	"github.com/pf-qiu/concourse/v6/fly/ui"
	"github.com/fatih/color"
)

type ResourcesCommand struct {
	Pipeline flaghelpers.PipelineFlag `short:"p" long:"pipeline" required:"true" description:"Get resources in this pipeline"`
	Json     bool                     `long:"json" description:"Print command result as JSON"`
}

func (command *ResourcesCommand) Execute([]string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	var headers []string
	var resources []atc.Resource

	resources, err = target.Team().ListResources(command.Pipeline.Ref())
	if err != nil {
		return err
	}

	if command.Json {
		err = displayhelpers.JsonPrint(resources)
		if err != nil {
			return err
		}
		return nil
	}

	headers = []string{"name", "type", "pinned", "check status"}
	table := ui.Table{Headers: ui.TableRow{}}
	for _, h := range headers {
		table.Headers = append(table.Headers, ui.TableCell{Contents: h, Color: color.New(color.Bold)})
	}

	for _, resource := range resources {
		var pinnedColumn ui.TableCell
		if resource.PinnedVersion != nil {
			pinnedColumn.Contents = ui.PresentVersion(resource.PinnedVersion)
		} else {
			pinnedColumn.Contents = "n/a"
		}

		var statusColumn ui.TableCell
		if resource.Build != nil {
			statusColumn = ui.BuildStatusCell(resource.Build.Status)
		} else {
			statusColumn.Contents = "n/a"
			statusColumn.Color = ui.OffColor
		}

		table.Data = append(table.Data, ui.TableRow{
			ui.TableCell{Contents: resource.Name},
			ui.TableCell{Contents: resource.Type},
			pinnedColumn,
			statusColumn,
		})
	}

	return table.Render(os.Stdout, Fly.PrintTableHeaders)
}
