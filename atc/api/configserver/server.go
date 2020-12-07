package configserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/creds"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

type Server struct {
	logger        lager.Logger
	teamFactory   db.TeamFactory
	secretManager creds.Secrets
}

func NewServer(
	logger lager.Logger,
	teamFactory db.TeamFactory,
	secretManager creds.Secrets,
) *Server {
	return &Server{
		logger:        logger,
		teamFactory:   teamFactory,
		secretManager: secretManager,
	}
}
