package geocodingCompare

import (
	"context"
	"reflect"

	"geo-smart/api/gen/models"
	"geo-smart/internal/geosmart"
	"go.uber.org/zap"
)

func New(log *zap.Logger, geocoder geosmart.MapsGeocoder, db geosmart.GeoDB) *GeocodingCompare {
	return &GeocodingCompare{
		geocoder: geocoder,
		log:    log,
		db:     db,
	}
}

type GeocodingCompare struct {
	geocoder geosmart.MapsGeocoder
	log    *zap.Logger
	db     geosmart.GeoDB
}

func (g *GeocodingCompare) GeocodeCompare(ctx context.Context, addresses map[string]models.Address, prodSuccessMap map[string]models.Success, envCompanyId string) error {
	adrrMap := map[string]models.GeoAddress{}
	for key, addr := range addresses {
		adrrMap[key] = models.GeoAddress{
			Address: Stringify(addr,"Pincode", "CountryCode"),
			Pin:     addr.Pincode,
			Country: addr.CountryCode,
		}
	}
	success, fail := g.geocode(ctx, adrrMap)
	for _, v := range fail {
		g.log.Info("", zap.String("Error: ", v.Message))
	}
	err := g.db.SaveGeocodeCompareData(prodSuccessMap, success, addresses, envCompanyId)
	if err != nil {
		return err
	}
	return nil
}

func (g *GeocodingCompare) geocode(ctx context.Context, addr map[string]models.GeoAddress) (map[string]models.Success, map[string]models.Fail) {
	success := map[string]models.Success{}
	fail := map[string]models.Fail{}
	failaddr := addr
	if len(failaddr) == 0 {
		return success, fail
	}
	s, f, failReason := g.geocoder.MapsGeocoder(ctx, addr, false)
	failaddr = *f
	for k, v := range *s {
		success[k] = v
	}
	for k, v := range *failReason {
		fail[k] = models.Fail{
			Code:    v.Code,
			Message: v.Message,
		}
	}

	return success, fail
}

func Stringify(addr models.Address, ignore ...string) string {
	str := ""
	v := reflect.ValueOf(addr)
OUTER:
	for i := 0; i < v.NumField(); i++ {
		for _, ign := range ignore {
			if v.Type().Field(i).Name == ign {
				continue OUTER
			}
		}

		val := v.Field(i).Interface().(string)
		if val != "" {
			if str != "" {
				str += ", "
			}
			str += val
		}
	}
	return str
}