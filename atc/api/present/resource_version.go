package present

import (
	"github.com/pf-qiu/concourse/v6/atc"
)

func ResourceVersions(hideMetadata bool, resourceVersions []atc.ResourceVersion) []atc.ResourceVersion {
	presented := []atc.ResourceVersion{}

	for _, resourceVersion := range resourceVersions {
		if hideMetadata {
			resourceVersion.Metadata = nil
		}

		presented = append(presented, resourceVersion)
	}

	return presented
}
