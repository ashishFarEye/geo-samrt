package platform

import (
	"context"
	"geo-smart/api/gen/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

var peliasSuccessMock = `
{
    "geocoding": {
        "version": "0.2",
        "attribution": "http://api.pelias.fareye.ml/attribution",
        "query": {
            "text": "20 W 34th St, New York, 10001, US",
            "size": 10,
            "private": false,
            "lang": {
                "name": "English",
                "iso6391": "en",
                "iso6393": "eng",
                "defaulted": false
            },
            "querySize": 20,
            "parser": "libpostal",
            "parsed_text": {
                "number": "20",
                "street": "w 34th st",
                "city": "new york",
                "postalcode": "10001",
                "country": "USA"
            }
        },
        "engine": {
            "name": "Pelias",
            "author": "Mapzen",
            "version": "1.0"
        },
        "timestamp": 1596200065523
    },
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -74.001432,
                    40.755535
                ]
            },
            "properties": {
                "id": "polyline:22203914",
                "gid": "openstreetmap:street:polyline:22203914",
                "layer": "street",
                "source": "openstreetmap",
                "source_id": "polyline:22203914",
                "name": "West 34th Street",
                "street": "West 34th Street",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "centroid",
                "country": "United States",
                "country_gid": "whosonfirst:country:85633793",
                "country_a": "USA",
                "region": "New York",
                "region_gid": "whosonfirst:region:85688543",
                "region_a": "NY",
                "county": "New York County",
                "county_gid": "whosonfirst:county:102081863",
                "county_a": "NE",
                "locality": "New York",
                "locality_gid": "whosonfirst:locality:85977539",
                "locality_a": "NYC",
                "borough": "Manhattan",
                "borough_gid": "whosonfirst:borough:421205771",
                "neighbourhood": "Hell's Kitchen",
                "neighbourhood_gid": "whosonfirst:neighbourhood:85868039",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "West 34th Street, Manhattan, New York, NY, USA"
            },
            "bbox": [
                -74.004807,
                40.748439,
                -73.984559,
                40.756954
            ]
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -73.987683,
                    40.749864
                ]
            },
            "properties": {
                "id": "polyline:10528652",
                "gid": "openstreetmap:street:polyline:10528652",
                "layer": "street",
                "source": "openstreetmap",
                "source_id": "polyline:10528652",
                "name": "West 34th Street",
                "street": "West 34th Street",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "centroid",
                "country": "United States",
                "country_gid": "whosonfirst:country:85633793",
                "country_a": "USA",
                "region": "New York",
                "region_gid": "whosonfirst:region:85688543",
                "region_a": "NY",
                "county": "New York County",
                "county_gid": "whosonfirst:county:102081863",
                "county_a": "NE",
                "locality": "New York",
                "locality_gid": "whosonfirst:locality:85977539",
                "locality_a": "NYC",
                "borough": "Manhattan",
                "borough_gid": "whosonfirst:borough:421205771",
                "neighbourhood": "Midtown West",
                "neighbourhood_gid": "whosonfirst:neighbourhood:85882233",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "West 34th Street, Manhattan, New York, NY, USA"
            },
            "bbox": [
                -73.987786,
                40.749824,
                -73.98758,
                40.749904
            ]
        }
    ],
    "bbox": [
        -74.004807,
        40.748439,
        -73.984559,
        40.756954
    ]
}`

var peliasPinMock = `
{
    "geocoding": {
        "version": "0.2",
        "attribution": "http://api.pelias.fareye.ml/attribution",
        "query": {
            "text": "10001, US",
            "size": 10,
            "private": false,
            "lang": {
                "name": "English",
                "iso6391": "en",
                "iso6393": "eng",
                "defaulted": false
            },
            "querySize": 20,
            "parser": "pelias",
            "parsed_text": {
                "subject": "10001",
                "postcode": "10001",
                "country": "US",
                "admin": "US"
            }
        },
        "engine": {
            "name": "Pelias",
            "author": "Mapzen",
            "version": "1.0"
        },
        "timestamp": 1596200098176
    },
    "type": "FeatureCollection",
    "features": [
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    16.032165,
                    45.896259
                ]
            },
            "properties": {
                "id": "polyline:20169928",
                "gid": "openstreetmap:street:polyline:20169928",
                "layer": "street",
                "source": "openstreetmap",
                "source_id": "polyline:20169928",
                "name": "10001",
                "street": "10001",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "centroid",
                "country": "Croatia",
                "country_gid": "whosonfirst:country:85633229",
                "country_a": "HRV",
                "region": "Grad Zagreb",
                "region_gid": "whosonfirst:region:85684811",
                "region_a": "GZ",
                "continent": "Europe",
                "continent_gid": "whosonfirst:continent:102191581",
                "label": "10001, Croatia"
            },
            "bbox": [
                16.031923,
                45.896057,
                16.03235,
                45.89643
            ]
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    11.43679,
                    44.459802
                ]
            },
            "properties": {
                "id": "polyline:8011304",
                "gid": "openstreetmap:street:polyline:8011304",
                "layer": "street",
                "source": "openstreetmap",
                "source_id": "polyline:8011304",
                "name": "10001",
                "street": "10001",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "centroid",
                "country": "Italy",
                "country_gid": "whosonfirst:country:85633253",
                "country_a": "ITA",
                "macroregion": "Emilia-Romagna",
                "macroregion_gid": "whosonfirst:macroregion:404227509",
                "region": "Bologna",
                "region_gid": "whosonfirst:region:85685351",
                "region_a": "BO",
                "localadmin": "San Lazzaro Di Savena",
                "localadmin_gid": "whosonfirst:localadmin:404472345",
                "continent": "Europe",
                "continent_gid": "whosonfirst:continent:102191581",
                "label": "10001, San Lazzaro Di Savena, Italy"
            },
            "bbox": [
                11.436697,
                44.459754,
                11.436883,
                44.45985
            ]
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -87.932209,
                    42.862893
                ]
            },
            "properties": {
                "id": "way/435916303",
                "gid": "openstreetmap:venue:way/435916303",
                "layer": "venue",
                "source": "openstreetmap",
                "source_id": "way/435916303",
                "name": "10001",
                "housenumber": "10001",
                "street": "South Hampton Drive",
                "postalcode": "53154",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "United States",
                "country_gid": "whosonfirst:country:85633793",
                "country_a": "USA",
                "region": "Wisconsin",
                "region_gid": "whosonfirst:region:85688517",
                "region_a": "WI",
                "county": "Milwaukee County",
                "county_gid": "whosonfirst:county:102081837",
                "county_a": "MU",
                "localadmin": "Oak Creek",
                "localadmin_gid": "whosonfirst:localadmin:404494177",
                "locality": "Oak Creek",
                "locality_gid": "whosonfirst:locality:101733079",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001, Oak Creek, WI, USA"
            },
            "bbox": [
                -87.9323186,
                42.8628251,
                -87.9321238,
                42.8629808
            ]
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    47.811249,
                    56.897647
                ]
            },
            "properties": {
                "id": "polyline:19222524",
                "gid": "openstreetmap:street:polyline:19222524",
                "layer": "street",
                "source": "openstreetmap",
                "source_id": "polyline:19222524",
                "name": "88Н-10001",
                "street": "88Н-10001",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "centroid",
                "country": "Russia",
                "country_gid": "whosonfirst:country:85632685",
                "country_a": "RUS",
                "region": "Mari El",
                "region_gid": "whosonfirst:region:85688043",
                "region_a": "ME",
                "county": "Orshanskiy",
                "county_gid": "whosonfirst:county:1108735711",
                "continent": "Europe",
                "continent_gid": "whosonfirst:continent:102191581",
                "label": "88Н-10001, Russia"
            },
            "bbox": [
                47.735839,
                56.877437,
                47.864437,
                56.917919
            ]
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -79.021464,
                    35.853881
                ]
            },
            "properties": {
                "id": "node/5785582951",
                "gid": "openstreetmap:address:node/5785582951",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/5785582951",
                "name": "10001 Fountain",
                "housenumber": "10001",
                "street": "Fountain",
                "postalcode": "27517",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "United States",
                "country_gid": "whosonfirst:country:85633793",
                "country_a": "USA",
                "region": "North Carolina",
                "region_gid": "whosonfirst:region:85688773",
                "region_a": "NC",
                "county": "Chatham County",
                "county_gid": "whosonfirst:county:102080923",
                "county_a": "CH",
                "locality": "Chapel Hill",
                "locality_gid": "whosonfirst:locality:85980563",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 Fountain, Chapel Hill, NC, USA"
            }
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -76.850262,
                    39.149806
                ]
            },
            "properties": {
                "id": "node/4301480482",
                "gid": "openstreetmap:address:node/4301480482",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/4301480482",
                "name": "10001 Anise Court",
                "housenumber": "10001",
                "street": "Anise Court",
                "postalcode": "20723",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "United States",
                "country_gid": "whosonfirst:country:85633793",
                "country_a": "USA",
                "region": "Maryland",
                "region_gid": "whosonfirst:region:85688501",
                "region_a": "MD",
                "county": "Howard County",
                "county_gid": "whosonfirst:county:102084263",
                "county_a": "HO",
                "locality": "Laurel",
                "locality_gid": "whosonfirst:locality:85949303",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 Anise Court, Laurel, MD, USA"
            }
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -117.288322,
                    56.228904
                ]
            },
            "properties": {
                "id": "node/2784176973",
                "gid": "openstreetmap:address:node/2784176973",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/2784176973",
                "name": "10001 104 Avenue",
                "housenumber": "10001",
                "street": "104 Avenue",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "Canada",
                "country_gid": "whosonfirst:country:85633041",
                "country_a": "CAN",
                "region": "Alberta",
                "region_gid": "whosonfirst:region:85682091",
                "region_a": "AB",
                "county": "Peace No. 135",
                "county_gid": "whosonfirst:county:1158863217",
                "locality": "Peace River",
                "locality_gid": "whosonfirst:locality:890457589",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 104 Avenue, Peace River, AB, Canada"
            }
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -117.322917,
                    56.232646
                ]
            },
            "properties": {
                "id": "node/2784176985",
                "gid": "openstreetmap:address:node/2784176985",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/2784176985",
                "name": "10001 85 Street",
                "housenumber": "10001",
                "street": "85 Street",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "Canada",
                "country_gid": "whosonfirst:country:85633041",
                "country_a": "CAN",
                "region": "Alberta",
                "region_gid": "whosonfirst:region:85682091",
                "region_a": "AB",
                "county": "Peace No. 135",
                "county_gid": "whosonfirst:county:1158863217",
                "locality": "Peace River",
                "locality_gid": "whosonfirst:locality:890457589",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 85 Street, Peace River, AB, Canada"
            }
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -117.308504,
                    56.232901
                ]
            },
            "properties": {
                "id": "node/2784176995",
                "gid": "openstreetmap:address:node/2784176995",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/2784176995",
                "name": "10001 89 Street",
                "housenumber": "10001",
                "street": "89 Street",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "Canada",
                "country_gid": "whosonfirst:country:85633041",
                "country_a": "CAN",
                "region": "Alberta",
                "region_gid": "whosonfirst:region:85682091",
                "region_a": "AB",
                "county": "Peace No. 135",
                "county_gid": "whosonfirst:county:1158863217",
                "locality": "Peace River",
                "locality_gid": "whosonfirst:locality:890457589",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 89 Street, Peace River, AB, Canada"
            }
        },
        {
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    -118.795145,
                    55.166037
                ]
            },
            "properties": {
                "id": "node/2079162097",
                "gid": "openstreetmap:address:node/2079162097",
                "layer": "address",
                "source": "openstreetmap",
                "source_id": "node/2079162097",
                "name": "10001 95 Avenue",
                "housenumber": "10001",
                "street": "95 Avenue",
                "confidence": 0.8,
                "match_type": "fallback",
                "accuracy": "point",
                "country": "Canada",
                "country_gid": "whosonfirst:country:85633041",
                "country_a": "CAN",
                "region": "Alberta",
                "region_gid": "whosonfirst:region:85682091",
                "region_a": "AB",
                "county": "Grande Prairie No. 1",
                "county_gid": "whosonfirst:county:1158863179",
                "locality": "Grande Prairie",
                "locality_gid": "whosonfirst:locality:890458665",
                "continent": "North America",
                "continent_gid": "whosonfirst:continent:102191575",
                "label": "10001 95 Avenue, Grande Prairie, AB, Canada"
            }
        }
    ],
    "bbox": [
        -118.795145,
        35.853881,
        47.864437,
        56.917919
    ]
}`

var PeliasnotFoundMock = `
{
    "geocoding": {
        "version": "0.2",
        "attribution": "http://api.pelias.fareye.ml/attribution",
        "query": {
            "text": "FarEye Sector 127, 201301, Noida, IN",
            "size": 10,
            "private": false,
            "lang": {
                "name": "English",
                "iso6391": "en",
                "iso6393": "eng",
                "defaulted": false
            },
            "querySize": 20,
            "parser": "pelias",
            "parsed_text": {
                "subject": "FarEye",
                "locality": "Sector",
                "postcode": "201301",
                "admin": "Sector 127, , Noida, IN"
            }
        },
        "engine": {
            "name": "Pelias",
            "author": "Mapzen",
            "version": "1.0"
        },
        "timestamp": 1596200270666
    },
    "type": "FeatureCollection",
    "features": []
}`

func TestNewPelias(t *testing.T) {
	log, _ := zap.NewProduction()
	p := NewPelias("", resty.New(), log)
	if p == nil {
		t.Errorf("Unable to initialize Pelias")
	}
}

func TestPeliasGeocoder(t *testing.T) {
	addr := map[string]models.GeoAddress{}
	addr["test"] = models.GeoAddress{
		Address: "20 W 34th St, New York",
		Pin:     "10001",
		Country: "US",
	}
	s, _ := peliasMockResp(peliasSuccessMock, addr, false)
	if !reflect.DeepEqual(models.Success{
		Lat:         40.755535,
		Lng:         -74.001432,
		Accuracy:    "0.800000",
		Source:      "FarEye Maps",
		Address:     "West 34th Street, Manhattan, New York, NY, USA",
		AddressComp: "{\"id\":\"polyline:22203914\",\"gid\":\"openstreetmap:street:polyline:22203914\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:22203914\",\"name\":\"West 34th Street\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"United States\",\"country_gid\":\"whosonfirst:country:85633793\",\"country_a\":\"USA\",\"region\":\"New York\",\"region_gid\":\"whosonfirst:region:85688543\",\"region_a\":\"NY\",\"continent\":\"North America\",\"continent_gid\":\"whosonfirst:continent:102191575\",\"label\":\"West 34th Street, Manhattan, New York, NY, USA\"}",
	}, s["test"]) {
		t.Errorf("Test Case Is Failing For Accurate Response With Partial False")
	}
	s, _ = peliasMockResp(peliasSuccessMock, addr, true)
	if !reflect.DeepEqual(models.Success{
		Lat:         40.755535,
		Lng:         -74.001432,
		Accuracy:    "0.800000",
		Source:      "FarEye Maps",
		Address:     "West 34th Street, Manhattan, New York, NY, USA",
		AddressComp: "{\"id\":\"polyline:22203914\",\"gid\":\"openstreetmap:street:polyline:22203914\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:22203914\",\"name\":\"West 34th Street\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"United States\",\"country_gid\":\"whosonfirst:country:85633793\",\"country_a\":\"USA\",\"region\":\"New York\",\"region_gid\":\"whosonfirst:region:85688543\",\"region_a\":\"NY\",\"continent\":\"North America\",\"continent_gid\":\"whosonfirst:continent:102191575\",\"label\":\"West 34th Street, Manhattan, New York, NY, USA\"}",
	}, s["test"]) {
		t.Errorf("Test Case Is Failing For Accurate Response With Partial True")
	}

	a := addr["test"]
	a.Address = ""
	addr["test"] = a
	s, _ = peliasMockResp(peliasPinMock, addr, true)
	if !reflect.DeepEqual(models.Success{
		Lat:         45.896259,
		Lng:         16.032165,
		Accuracy:    "0.800000",
		Source:      "FarEye Maps",
		Address:     "10001, Croatia",
		AddressComp: "{\"id\":\"polyline:20169928\",\"gid\":\"openstreetmap:street:polyline:20169928\",\"layer\":\"street\",\"source\":\"openstreetmap\",\"source_id\":\"polyline:20169928\",\"name\":\"10001\",\"confidence\":0.8,\"match_type\":\"fallback\",\"accuracy\":\"centroid\",\"country\":\"Croatia\",\"country_gid\":\"whosonfirst:country:85633229\",\"country_a\":\"HRV\",\"region\":\"Grad Zagreb\",\"region_gid\":\"whosonfirst:region:85684811\",\"region_a\":\"GZ\",\"continent\":\"Europe\",\"continent_gid\":\"whosonfirst:continent:102191581\",\"label\":\"10001, Croatia\"}",
	}, s["test"]) {
		t.Errorf("Test Case Is Failing For Partial Response With Partial True")
	}
	/*_,f:= peliasMockResp(peliasPinMock, addr,false)
	if !reflect.DeepEqual(models.Fail{Code:2,Message: "Partial Accuraccy Not Enabled"}, f["test"]){
		t.Errorf("Test Case Is Failing For Accurate Response With Partial True")
	}*/

	addr["test"] = models.GeoAddress{
		Address: "FarEye, Sec 127, Noida",
		Pin:     "201301",
		Country: "IN",
	}
	_, f := peliasMockResp(PeliasnotFoundMock, addr, true)
	if !reflect.DeepEqual(models.Fail{Code: "1", Message: "Not Found By Pelias"}, f["test"]) {
		t.Errorf("Test Case Is Failing For Not Found Response With Partial True")
	}
}

func peliasMockResp(resp string, addr map[string]models.GeoAddress, partial bool) (map[string]models.Success, map[string]models.Fail) {
	log, _ := zap.NewProduction()
	server := mockServer(200, resp)
	defer server.Close()
	p := NewPelias(server.URL, resty.New(), log)
	s, _, f := p.MapsGeocoder(context.Background(), addr, partial)
	return *s, *f
}
