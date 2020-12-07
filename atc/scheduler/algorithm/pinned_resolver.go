package algorithm

import (
	"context"

	"github.com/pf-qiu/concourse/v6/atc/db"
	"github.com/pf-qiu/concourse/v6/tracing"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type pinnedResolver struct {
	vdb         db.VersionsDB
	inputConfig db.InputConfig
}

func NewPinnedResolver(vdb db.VersionsDB, inputConfig db.InputConfig) Resolver {
	return &pinnedResolver{
		vdb:         vdb,
		inputConfig: inputConfig,
	}
}

func (r *pinnedResolver) InputConfigs() db.InputConfigs {
	return db.InputConfigs{r.inputConfig}
}

func (r *pinnedResolver) Resolve(ctx context.Context) (map[string]*versionCandidate, db.ResolutionFailure, error) {
	ctx, span := tracing.StartSpan(ctx, "pinnedResolver.Resolve", tracing.Attrs{
		"input": r.inputConfig.Name,
	})
	defer span.End()

	version, found, err := r.vdb.FindVersionOfResource(ctx, r.inputConfig.ResourceID, r.inputConfig.PinnedVersion)
	if err != nil {
		tracing.End(span, err)
		return nil, "", err
	}

	if !found {
		span.AddEvent(ctx, "pinned version not found")
		span.SetStatus(codes.NotFound, "pinned version not found")
		return nil, db.PinnedVersionNotFound{PinnedVersion: r.inputConfig.PinnedVersion}.String(), nil
	}

	span.AddEvent(ctx, "found via pin", label.String("version", string(version)))

	versionCandidate := map[string]*versionCandidate{
		r.inputConfig.Name: newCandidateVersion(version),
	}

	span.SetStatus(codes.OK, "found via pin")
	return versionCandidate, "", nil
}
