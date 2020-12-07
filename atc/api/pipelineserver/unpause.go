package pipelineserver

import (
	"net/http"

	"github.com/pf-qiu/concourse/v6/atc/db"
)

func (s *Server) UnpausePipeline(pipelineDB db.Pipeline) http.Handler {
	logger := s.logger.Session("unpause-pipeline")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := pipelineDB.Unpause()

		if err != nil {
			logger.Error("failed-to-unpause-pipeline", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.teamFactory.NotifyResourceScanner()
		if err != nil {
			logger.Error("failed-to-notify-resource-scanner", err)
		}

		w.WriteHeader(http.StatusOK)
	})
}
