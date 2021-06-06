package postgres

import (
	"geo-smart/api/gen/models"
	"geo-smart/internal/geosmart"
	"github.com/go-pg/pg/v9"
	"github.com/umahmood/haversine"
)

func NewGeoDB(db *pg.DB) *GeoDB {
	return &GeoDB{
		db: db,
	}
}

type GeoDB struct {
	db *pg.DB
}

func (p *GeoDB) Save(result map[string]models.Success, addrMap map[string]models.Address) error {
	var address []geosmart.Address
	for key, val := range result {
		address = append(address, geosmart.Address{
			Lat:           val.Lat,
			Lng:           val.Lng,
			InputLine1:    addrMap[key].Line1,
			InputLine2:    addrMap[key].Line2,
			InputLandmark: addrMap[key].Landmark,
			InputPin:      addrMap[key].Pincode,
			InputCountry:  addrMap[key].CountryCode,
			Accuracy:      val.Accuracy,
			Source:        val.Source,
			Address:       val.Address,
			AddressComp:   val.AddressComp,
		})
	}
	err := p.db.Insert(&address)
	if err != nil {
		return err
	}
	return nil
}

func (p *GeoDB) SaveGeocodeCompareData(prodServiceResponse map[string]models.Success, alphaServiceResponse map[string]models.Success,
	geoAddress map[string]models.Address, envCompanyId string) error {

	var addressCompare []geosmart.AddressCompare
	for key, val := range geoAddress {
		_, distKm := haversine.Distance(haversine.Coord{prodServiceResponse[key].Lat, prodServiceResponse[key].Lng}, haversine.Coord{alphaServiceResponse[key].Lat, alphaServiceResponse[key].Lng})
		addressCompare = append(addressCompare, geosmart.AddressCompare{
			InputLine1:       val.Line1,
			InputLine2:       val.Line2,
			InputLandmark:    val.Landmark,
			InputPin:         val.Pincode,
			InputCountry:     val.CountryCode,
			EnvCompanyId:     envCompanyId,
			ProdAccuracy:     prodServiceResponse[key].Accuracy,
			ProdSource:       prodServiceResponse[key].Source,
			ProdLat:          prodServiceResponse[key].Lat,
			ProdLng:          prodServiceResponse[key].Lng,
			ProdAddress:      prodServiceResponse[key].Address,
			ProdAddressComp:  prodServiceResponse[key].AddressComp,
			AlphaAccuracy:    alphaServiceResponse[key].Accuracy,
			AlphaSource:      alphaServiceResponse[key].Source,
			AlphaLat:         alphaServiceResponse[key].Lat,
			AlphaLng:         alphaServiceResponse[key].Lng,
			AlphaAddress:     alphaServiceResponse[key].Address,
			AlphaAddressComp: alphaServiceResponse[key].AddressComp,
			AerialDistance:   distKm,
		})
	}
	err := p.db.Insert(&addressCompare)
	if err != nil {
		return err
	}
	return nil
}
