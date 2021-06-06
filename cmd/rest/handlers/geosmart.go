package handlers

import (
	"geo-smart/api/gen/restapi/operations"
	"geo-smart/internal/geosmart"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

func NewGeocoder(svc geosmart.Service, log *zap.Logger) GeoSmart {
	return GeoSmart{
		svc: svc,
		log: log,
	}
}

type GeoSmart struct {
	svc geosmart.Service
	log *zap.Logger
}

func (g GeoSmart) Handle(params operations.GeosmartParams) middleware.Responder {
	addrMap := params.Geosmart.Address
	partial := params.Geosmart.Partial
	resp, err := g.svc.Geocode(params.HTTPRequest.Context(), addrMap, partial)
	if err != nil {
		g.log.Error("Error GeoCoding", zap.Error(err))
		return operations.NewGeosmartGoInternalServerError().WithPayload(resp)
	}
	return operations.NewGeosmartOK().WithPayload(resp)
}
