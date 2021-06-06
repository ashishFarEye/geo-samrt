package geosmart

import (
	"context"
	"geo-smart/api/gen/models"
)

type MapsGeocoder interface {
	MapsGeocoder(context.Context,map[string]models.GeoAddress, bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail)
}

type Service interface {
	Geocode(context.Context, map[string]models.GeoAddress, bool) (*models.GeoSmartResponse, error)
}

type Geocoder interface {
	Geocoder(context.Context, map[string]models.GeoAddress, bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail)
}

type GeoDB interface {
	Save(map[string]models.Success, map[string]models.Address) error
	SaveGeocodeCompareData(map[string]models.Success, map[string]models.Success, map[string]models.Address, string) error
}

type Address struct {
	Id            int64
	Lat           float64
	Lng           float64
	InputLine1    string
	InputLine2    string
	InputLandmark string
	InputPin      string
	InputCountry  string
	Accuracy      string
	Source        string
	Address       string
	AddressComp   string
	CreatedAt     string
}

type AddressCompare struct {
	Id               int64
	InputLine1       string
	InputLine2       string
	InputLandmark    string
	InputPin         string
	InputCountry     string
	EnvCompanyId     string
	ProdAccuracy     string
	ProdSource       string
	ProdLat          float64
	ProdLng          float64
	ProdAddress      string
	ProdAddressComp  string
	AlphaAccuracy    string
	AlphaSource      string
	AlphaLat         float64
	AlphaLng         float64
	AlphaAddress     string
	AlphaAddressComp string
	AerialDistance   float64
}
