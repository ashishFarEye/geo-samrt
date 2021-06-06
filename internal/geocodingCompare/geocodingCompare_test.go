package geocodingCompare

import (
	"context"
	"geo-smart/api/gen/models"
	"go.uber.org/zap"
	"testing"
)

type geoDbMock struct {
}

func (p geoDbMock) SaveGeocodeCompareData(m map[string]models.Success, m2 map[string]models.Success, m3 map[string]models.Address, s string) error {
	return nil
}

func (p geoDbMock) Save(result map[string]models.Success, addrMap map[string]models.Address) error {
	return nil
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

func TestGeocodingCompareGeocodeCompare(t *testing.T) {
	peliasSuccessMock := map[string]models.Success{
		"test1": {
			Lat:         40.749864,
			Lng:         -73.987683,
			Accuracy:    "0.800000",
			Source:      "FarEye Maps",
			Address:     "West 34th Street, Manhattan, New York, NY, USA",
			AddressComp: "{\"id\":\"polyline:22203914\",\"gid\":\"openstreetmap:street:polyline:22203914\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:22203914\",\"name\":\"West 34th Street\",\"street\":\"West 34th Street\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"United States\",\"country_gid\":\"whosonfirst:country:85633793\",\"country_a\":\"USA\",\"region\":\"New York\",\"region_gid\":\"whosonfirst:region:85688543\",\"region_a\":\"NY\",\"county\":\"New York County\",\"county_gid\":\"whosonfirst:county:102081863\",\"county_a\":\"NE\",\"locality\":\"New York\",\"locality_gid\":\"whosonfirst:locality:85977539\",\"locality_a\":\"NYC\",\"borough\":\"Manhattan\",\"borough_gid\":\"whosonfirst:borough:421205771\",\"neighbourhood\":\"Hell's Kitchen\",\"neighbourhood_gid\":\"whosonfirst:neighbourhood:85868039\",\"continent\":\"North America\",\"continent_gid\":\"whosonfirst:continent:102191575\",\"label\":\"West 34th Street, Manhattan, New York, NY, USA\"}",
		},
	}
	peliasFailMock := map[string]models.Fail{
		"test2": {Code: 1, Message: "Not Found By Pelias"},
	}
	p := peliasMock{
		succ: peliasSuccessMock,
		fail: peliasFailMock,
	}
	log, _ := zap.NewProduction()
	g := New(log, &p, &geoDbMock{})
	addr := map[string]models.Address{
		"test1": {
			Line1:       "20 W",
			Line2:       "34th St",
			Landmark:    "New York",
			Pincode:     "10001",
			CountryCode: "US",
		},
		"test2": {
			Line1:       "FarEye",
			Line2:       "Sec 127",
			Landmark:    "Noida",
			Pincode:     "201302",
			CountryCode: "IN",
		},
	}

	succMap := map[string]models.Success{
		"test1": {
			Lat:         40.7487836,
			Lng:         -73.98615769999999,
			Accuracy:    "Accurate",
			Source:      "Google",
			Address:     "20 W 34th St, New York, NY 10001, USA",
			AddressComp: "[{\"types\":[\"street_number\"],\"short_name\":\"20\",\"long_name\":\"20\"},{\"types\":[\"route\"],\"short_name\":\"W 34th St\",\"long_name\":\"West 34th Street\"},{\"types\":[\"political\",\"sublocality\",\"sublocality_level_1\"],\"short_name\":\"Manhattan\",\"long_name\":\"Manhattan\"},{\"types\":[\"locality\",\"political\"],\"short_name\":\"New York\",\"long_name\":\"New York\"},{\"types\":[\"administrative_area_level_2\",\"political\"],\"short_name\":\"New York County\",\"long_name\":\"New York County\"},{\"types\":[\"administrative_area_level_1\",\"political\"],\"short_name\":\"NY\",\"long_name\":\"New York\"},{\"types\":[\"country\",\"political\"],\"short_name\":\"US\",\"long_name\":\"United States\"},{\"types\":[\"postal_code\"],\"short_name\":\"10001\",\"long_name\":\"10001\"}]",
		},
		"test2": {
			Lat:         28.5821195,
			Lng:         77.3266991,
			Accuracy:    "Accurate",
			Source:      "Google",
			Address:     "Noida, Uttar Pradesh 201301, India",
			AddressComp: "[{\"long_name\":\"201301\",\"short_name\":\"201301\",\"types\":[\"postal_code\"]},{\"long_name\":\"Noida\",\"short_name\":\"Noida\",\"types\":[\"locality\",\"political\"]},{\"long_name\":\"Uttar Pradesh\",\"short_name\":\"UP\",\"types\":[\"administrative_area_level_1\",\"political\"]},{\"long_name\":\"India\",\"short_name\":\"IN\",\"types\":[\"country\",\"political\"]}]",
		},
	}
	err := g.GeocodeCompare(context.Background(), addr, succMap, "geocoding.fareye.co:1")
	if err != nil {
		t.Errorf("Error Nor Desired")
	}
}
