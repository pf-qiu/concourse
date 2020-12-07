package resourceserver

import (
	"encoding/json"
	"net/http"

	"github.com/pf-qiu/concourse/v6/atc/api/present"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

func (s *Server) ListVersionedResourceTypes(pipeline db.Pipeline) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.Session("list-versioned-resource-types")

		resourceTypes, err := pipeline.ResourceTypes()
		if err != nil {
			logger.Error("failed-to-get-resources-types", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		versionedResourceTypes := present.VersionedResourceTypes(resourceTypes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(versionedResourceTypes)
		if err != nil {
			logger.Error("failed-to-encode-versioned-resource-types", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
