package resourceserver

import (
	"encoding/json"
	"net/http"

	"github.com/pf-qiu/concourse/v6/atc"
	"github.com/pf-qiu/concourse/v6/atc/api/accessor"
	"github.com/pf-qiu/concourse/v6/atc/api/present"
	"github.com/pf-qiu/concourse/v6/atc/db"
)

func (s *Server) ListAllResources(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-all-resources")

	acc := accessor.GetAccessor(r)

	var (
		dbResources []db.Resource
		err         error
	)
	if acc.IsAdmin() {
		dbResources, err = s.resourceFactory.AllResources()
	} else {
		dbResources, err = s.resourceFactory.VisibleResources(acc.TeamNames())
	}
	if err != nil {
		logger.Error("failed-to-get-all-visible-resources", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resources := []atc.Resource{}

	for _, resource := range dbResources {
		resources = append(
			resources,
			present.Resource(resource),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resources)
	if err != nil {
		logger.Error("failed-to-encode-resources", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
