package commands

import (
	"os"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/displayhelpers"
	"github.com/pf-qiu/concourse/v6/fly/commands/internal/flaghelpers"
	"github.com/pf-qiu/concourse/v6/fly/rc"
	"github.com/pf-qiu/concourse/v6/fly/ui"
	"github.com/pf-qiu/concourse/v6/go-concourse/concourse"
	"github.com/fatih/color"
)

type JobsCommand struct {
	Pipeline flaghelpers.PipelineFlag `short:"p" long:"pipeline" required:"true" description:"Get jobs in this pipeline"`
	Json     bool                     `long:"json" description:"Print command result as JSON"`
	Team     string                   `long:"team" description:"Name of the team to which the pipeline belongs, if different from the target default"`
}

func (command *JobsCommand) Execute([]string) error {
	var (
		headers []string
		team    concourse.Team
	)

	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	if command.Team != "" {
		team, err = target.FindTeam(command.Team)
		if err != nil {
			return err
		}
	} else {
		team = target.Team()
	}

	var jobs []atc.Job
	jobs, err = team.ListJobs(command.Pipeline.Ref())
	if err != nil {
		return err
	}

	if command.Json {
		err = displayhelpers.JsonPrint(jobs)
		if err != nil {
			return err
		}
		return nil
	}

	headers = []string{"name", "paused", "status", "next"}
	table := ui.Table{Headers: ui.TableRow{}}
	for _, h := range headers {
		table.Headers = append(table.Headers, ui.TableCell{Contents: h, Color: color.New(color.Bold)})
	}

	for _, p := range jobs {
		var pausedColumn ui.TableCell
		if p.Paused {
			pausedColumn.Contents = "yes"
			pausedColumn.Color = color.New(color.FgCyan)
		} else {
			pausedColumn.Contents = "no"
		}

		row := ui.TableRow{}
		row = append(row, ui.TableCell{Contents: p.Name})

		row = append(row, pausedColumn)

		var statusColumn ui.TableCell
		if p.FinishedBuild != nil {
			statusColumn = ui.BuildStatusCell(p.FinishedBuild.Status)
		} else {
			statusColumn.Contents = "n/a"
		}
		row = append(row, statusColumn)

		var nextColumn ui.TableCell
		if p.NextBuild != nil {
			nextColumn = ui.BuildStatusCell(p.NextBuild.Status)
		} else {
			nextColumn.Contents = "n/a"
		}
		row = append(row, nextColumn)

		table.Data = append(table.Data, row)
	}

	return table.Render(os.Stdout, Fly.PrintTableHeaders)
}
