package commands

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
	"github.com/pf-qiu/concourse/v6/fly/ui"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse"
	"github.com/fatih/color"
)

type ResourceVersionsCommand struct {
	Count    int                      `short:"c" long:"count" default:"50" description:"Number of versions you want to limit the return to"`
	Resource flaghelpers.ResourceFlag `short:"r" long:"resource" required:"true" value-name:"PIPELINE/RESOURCE" description:"Name of a resource to get versions for"`
	Json     bool                     `long:"json" description:"Print command result as JSON"`
}

func (command *ResourceVersionsCommand) Execute([]string) error {
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	page := concourse.Page{Limit: command.Count}

	team := target.Team()

	versions, _, _, err := team.ResourceVersions(command.Resource.PipelineRef, command.Resource.ResourceName, page, atc.Version{})
	if err != nil {
		return err
	}

	if command.Json {
		err = displayhelpers.JsonPrint(versions)
		if err != nil {
			return err
		}
		return nil
	}

	table := ui.Table{
		Headers: ui.TableRow{
			{Contents: "id", Color: color.New(color.Bold)},
			{Contents: "version", Color: color.New(color.Bold)},
			{Contents: "enabled", Color: color.New(color.Bold)},
		},
	}

	var rangeUntil int
	if command.Count < len(versions) {
		rangeUntil = command.Count
	} else {
		rangeUntil = len(versions)
	}

	for _, version := range versions[:rangeUntil] {
		var enabledCell ui.TableCell
		if version.Enabled {
			enabledCell.Color = ui.OnColor
			enabledCell.Contents = "yes"
		} else {
			enabledCell.Contents = "no"
		}

		fields := []string{}
		for k, v := range version.Version {
			fields = append(fields, k+":"+v)
		}

		sort.Strings(fields)

		table.Data = append(table.Data, []ui.TableCell{
			{Contents: strconv.Itoa(version.ID)},
			{Contents: strings.Join(fields, ",")},
			enabledCell,
		})
	}

	return table.Render(os.Stdout, Fly.PrintTableHeaders)
}
