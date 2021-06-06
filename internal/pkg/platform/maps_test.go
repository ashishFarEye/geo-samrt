package platform

import (
	"context"
	"geo-smart/api/gen/models"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewMaps(t *testing.T) {
	log, _ := zap.NewProduction()
	m := NewMaps(log, &ES{}, &Pelias{})
	if m == nil {
		t.Errorf("Unable to initialize Map")
	}
}

type peliasMock struct {
	succ map[string]models.Success
	fail map[string]models.Fail
}

func (p *peliasMock) MapsGeocoder(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail) {
	address := map[string]models.GeoAddress{}
	address["test2"] = addr["test2"]
	return &p.succ, &address, &p.fail
}

type esMock struct {
	succ map[string]models.Success
	fail map[string]models.Fail
}

func (e *esMock) MapsGeocoder(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail) {
	address := map[string]models.GeoAddress{}
	address["test2"] = addr["test2"]
	return &e.succ, &address, &e.fail
}

func TestMapsGeocoder(t *testing.T) {
	m, addr := mockMaps()
	s, a, f := mapsMockResponse(addr, false, m)
	if !reflect.DeepEqual(models.Success{
		Lat:         45.896259,
		Lng:         16.032165,
		Accuracy:    "0.800000",
		Source:      "FarEye Maps",
		Address:     "10001, Croatia",
		AddressComp: "{\"id\":\"polyline:20169928\",\"gid\":\"openstreetmap:street:polyline:20169928\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:20169928\",\"name\":\"10001\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"Croatia\",\"country_gid\":\"whosonfirst:country:85633229\",\"country_a\":\"HRV\",\"region\":\"Grad Zagreb\",\"region_gid\":\"whosonfirst:region:85684811\",\"region_a\":\"GZ\",\"continent\":\"Europe\",\"continent_gid\":\"whosonfirst:continent:102191581\",\"label\":\"10001, Croatia\"}",
	}, s["test1"]) {
		t.Errorf("Test Case Is Failing For Accurate Response With Partial False")
	}
	if !reflect.DeepEqual(models.GeoAddress{
		Address: "20 W 34th St",
		Pin:     "",
		Country: "US",
		UUID:    "",
	}, a["test2"]) {
		t.Errorf("Test Case Is Failing For Accurate Response With Partial False")
	}

	if !reflect.DeepEqual(models.Fail{Code: "1", Message: "Not Found By Pelias"}, f["test2"]) {
		t.Errorf("Test Case Is Failing For Accurate Response With Partial False")
	}

}

func mockMaps() (*Maps, map[string] models.GeoAddress) {
	log, _ := zap.NewProduction()
	addr := map[string]models.GeoAddress{
		"test1": models.GeoAddress{
			Address: "20 W 34th St, New York",
			Pin:     "10001",
			Country: "US",
		},
		"test2": models.GeoAddress{
			Address: "20 W 34th St",
			Pin:     "",
			Country: "US",
		},
	}
	peliasSuccessMock := map[string]models.Success{
		"test1": models.Success{
			Lat:         45.896259,
			Lng:         16.032165,
			Accuracy:    "0.800000",
			Source:      "FarEye Maps",
			Address:     "10001, Croatia",
			AddressComp: "{\"id\":\"polyline:20169928\",\"gid\":\"openstreetmap:street:polyline:20169928\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:20169928\",\"name\":\"10001\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"Croatia\",\"country_gid\":\"whosonfirst:country:85633229\",\"country_a\":\"HRV\",\"region\":\"Grad Zagreb\",\"region_gid\":\"whosonfirst:region:85684811\",\"region_a\":\"GZ\",\"continent\":\"Europe\",\"continent_gid\":\"whosonfirst:continent:102191581\",\"label\":\"10001, Croatia\"}",
		},
	}
	peliasFailMock := map[string]models.Fail{
		"test2": models.Fail{Code: "1", Message: "Not Found By Pelias"},
	}
	p := peliasMock{
		succ: peliasSuccessMock,
		fail: peliasFailMock,
	}
	return NewMaps(log, &esMock{}, &p), addr
}

func mapsMockResponse(addr map[string]models.GeoAddress, partial bool, m *Maps) (map[string]models.Success, map[string]models.GeoAddress, map[string]models.Fail) {
	s, a, f := m.Geocoder(context.Background(), addr, partial)
	return *s, *a, *f
}
