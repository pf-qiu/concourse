package main

import (
	"github.com/pf-qiu/concourse/v6/atc/atccmd"
	"github.com/pf-qiu/concourse/v6/worker/land"
	"github.com/pf-qiu/concourse/v6/worker/retire"
	"github.com/pf-qiu/concourse/v6/worker/workercmd"
	"github.com/jessevdk/go-flags"
)

type ConcourseCommand struct {
	Version func() `short:"v" long:"version" description:"Print the version of Concourse and exit"`

	Web     WebCommand              `command:"web"     description:"Run the web UI and build scheduler."`
	Worker  workercmd.WorkerCommand `command:"worker"  description:"Run and register a worker."`
	Migrate atccmd.Migration        `command:"migrate" description:"Run database migrations."`

	Quickstart QuickstartCommand `command:"quickstart" description:"Run both 'web' and 'worker' together, auto-wired. Not recommended for production."`

	LandWorker   land.LandWorkerCommand     `command:"land-worker" description:"Safely drain a worker's assignments for temporary downtime."`
	RetireWorker retire.RetireWorkerCommand `command:"retire-worker" description:"Safely remove a worker from the cluster permanently."`

	GenerateKey GenerateKeyCommand `command:"generate-key" description:"Generate RSA key for use with Concourse components."`
}

func (cmd ConcourseCommand) LessenRequirements(parser *flags.Parser) {
	cmd.Quickstart.LessenRequirements(parser.Find("quickstart"))
	cmd.Web.LessenRequirements(parser.Find("web"))
	cmd.Worker.LessenRequirements("", parser.Find("worker"))
}
