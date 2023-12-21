package graph

import (
	"context"

	"cubawheeler.io/pkg/cubawheeler"
)

type voiceInstructionsResolver struct{ *Resolver }

// DistanceAlongGeometry is the resolver for the distanceAlongGeometry field.
func (r *voiceInstructionsResolver) DistanceAlongGeometry(ctx context.Context, obj *cubawheeler.VoiceInstructions) (int, error) {
	return int(obj.DistanceAlongGeometry), nil
}
