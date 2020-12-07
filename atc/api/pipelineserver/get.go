package pipelineserver

import (
	"encoding/json"
	"net/http"

	"github.com/pf-qiu/concourse/v6/atc/api/present"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

func (s *Server) GetPipeline(pipeline db.Pipeline) http.Handler {
	logger := s.logger.Session("get-pipeline")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(present.Pipeline(pipeline))
		if err != nil {
			logger.Error("failed-to-encode-pipeline", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
