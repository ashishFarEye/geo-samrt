package platform

import (
	"context"
	"geo-smart/api/gen/models"
	"geo-smart/internal/geosmart"
	"go.uber.org/zap"
)

func NewMaps(log *zap.Logger, es geosmart.MapsGeocoder, p geosmart.MapsGeocoder) *Maps {
	return &Maps{
		log: log,
		es:  es,
		p:   p,
	}
}

type Maps struct {
	log *zap.Logger
	es  geosmart.MapsGeocoder
	p   geosmart.MapsGeocoder
}

func (m *Maps) Geocoder(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail) {
	success, fail, f, es := map[string]models.Success{}, map[string]models.GeoAddress{}, map[string]models.Fail{}, map[string]models.GeoAddress{}
	for key, val := range addr {
		es[key] = val
		//if peliasCountry[val.Country] {
		//	pelias[key] = val
		//} else {
		//	es[key] = val
		//}
	}
	m.GeocodeFromEs(ctx, es, partial, success, fail, f)
	//m.GeocodeFromPelias(ctx, pelias, partial, success, fail, f)
	return &success, &fail, &f
}

func (m *Maps) GeocodeFromEs(ctx context.Context, addr map[string]models.GeoAddress, partial bool, success map[string]models.Success, fail map[string]models.GeoAddress, f map[string]models.Fail) {
	succesEs, failEs, fEs := m.es.MapsGeocoder(ctx, addr, partial)
	for k, v := range *succesEs {
		success[k] = v
	}
	for k, v := range *failEs {
		fail[k] = v
	}
	for k, v := range *fEs {
		f[k] = v
	}
}

func (m *Maps) GeocodeFromPelias(ctx context.Context, addr map[string]models.GeoAddress, partial bool, success map[string]models.Success, fail map[string]models.GeoAddress, f map[string]models.Fail) {
	succesPelias, failPelias, fPelias := m.p.MapsGeocoder(ctx, addr, partial)
	for k, v := range *succesPelias {
		success[k] = v
	}
	for k, v := range *failPelias {
		fail[k] = v
	}
	for k, v := range *fPelias {
		f[k] = v
	}
}

var peliasCountry = map[string]bool{
	"US": true,
}
