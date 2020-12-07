package ccserver

import (
	"code.cloudfoundry.org/lager"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

type Server struct {
	logger      lager.Logger
	teamFactory db.TeamFactory
	externalURL string
}

func NewServer(
	logger lager.Logger,
	teamFactory db.TeamFactory,
	externalURL string,
) *Server {
	return &Server{
		logger:      logger,
		teamFactory: teamFactory,
		externalURL: externalURL,
	}
}
