package handlers

import (
	"fmt"

	"geo-smart/api/gen/models"
	"geo-smart/api/gen/restapi/operations"
	"geo-smart/internal/geocodingCompare"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

func NewGeocoderCompare(svc geocodingCompare.Service, log *zap.Logger) GeocodingCompare {
	return GeocodingCompare{
		svc: svc,
		log: log,
	}
}

type GeocodingCompare struct {
	svc geocodingCompare.Service
	log *zap.Logger
}

func (g GeocodingCompare) Handle(params operations.CompareParams) middleware.Responder {
	addrMap := map[string]models.Address{}
	successMap := map[string]models.Success{}
	for k, prodResponse := range params.Data.Response {
		successMap[k] = models.Success{
			Lat:         prodResponse.Lat,
			Lng:         prodResponse.Lng,
			Accuracy:    prodResponse.Accuracy,
			Source:      prodResponse.Source,
			Address:     prodResponse.Address,
			AddressComp: prodResponse.AddressComp,
		}
	}
	for k, addr := range params.Data.Request {
		addrMap[k] = models.Address{
			Line1:       addr.Line1,
			Line2:       addr.Line2,
			Landmark:    addr.Landmark,
			Pincode:     addr.Pincode,
			CountryCode: addr.CountryCode,
		}
	}
	err := g.svc.GeocodeCompare(params.HTTPRequest.Context(), addrMap, successMap, params.Data.EnvCompanyID)
	if err != nil {
		g.log.Error("Error GeoCoding", zap.Error(err))
		return operations.NewCompareDefault(500).WithPayload(fmt.Sprintf("Data Recording failed. Error: %s", err.Error()))
	}
	g.log.Info("Data recorded")
	return operations.NewCompareOK().WithPayload("Data Recorded")
}
