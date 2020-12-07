package jobserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/api/auth"
	"github.com/pf-qiu/concourse/v6/atc/creds"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

type Server struct {
	logger lager.Logger

	externalURL   string
	rejector      auth.Rejector
	secretManager creds.Secrets
	jobFactory    db.JobFactory
	checkFactory  db.CheckFactory
}

func NewServer(
	logger lager.Logger,
	externalURL string,
	secretManager creds.Secrets,
	jobFactory db.JobFactory,
	checkFactory db.CheckFactory,
) *Server {
	return &Server{
		logger:        logger,
		externalURL:   externalURL,
		rejector:      auth.UnauthorizedRejector{},
		secretManager: secretManager,
		jobFactory:    jobFactory,
		checkFactory:  checkFactory,
	}
}
