package geosmart

import (
	"context"
	"geo-smart/api/gen/models"
	"go.uber.org/zap"
)
func New(log *zap.Logger, maps Geocoder) *GeoSmart {
	return &GeoSmart{
		maps: maps,
		log:  log,
	}
}

type GeoSmart struct {
	maps  Geocoder
	log   *zap.Logger
}

func (g *GeoSmart) Geocode(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*models.GeoSmartResponse, error){
	s, f, failReason := g.maps.Geocoder(ctx, addr, partial)

	return &models.GeoSmartResponse{
		Address: *f,
		Fail:    *failReason,
		Success: *s,
	}, nil
}