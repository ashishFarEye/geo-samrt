package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"geo-smart/api/gen/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

)

type Pelias struct {
	url    string
	client *resty.Client
	log    *zap.Logger
}

type peliasResponse struct {
	Features []struct {
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Id            string  `json:"id"`
			Gid           string  `json:"gid"`
			Layer         string  `json:"layer"`
			Source        string  `json:"source"`
			Source_id     string  `json:"source_id"`
			Name          string  `json:"name"`
			Confidence    float64 `json:"confidence"`
			Match_type    string  `json:"match_type"`
			Accuracy      string  `json:"accuracy"`
			Country       string  `json:"country"`
			Country_gid   string  `json:"country_gid"`
			Country_a     string  `json:"country_a"`
			Region        string  `json:"region"`
			Region_gid    string  `json:"region_gid"`
			Region_a      string  `json:"region_a"`
			Continent     string  `json:"continent"`
			Continent_gid string  `json:"continent_gid"`
			Label         string  `json:"label"`
		} `json:"properties"`
	} `json:"features"`
}

func NewPelias(url string, client *resty.Client, log *zap.Logger) *Pelias {
	return &Pelias{
		url:    url,
		client: client,
		log:    log,
	}
}

func (p *Pelias) MapsGeocoder(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail) {
	success, fail, succMutex, failMutex, f := map[string]models.Success{}, map[string]models.GeoAddress{}, sync.RWMutex{}, sync.RWMutex{}, map[string]models.Fail{}
	p.log.Debug("---Sending Request To Pelias---", zap.Reflect("Address", &addr))
	var wg sync.WaitGroup
	for key, val := range addr {
		wg.Add(1)
		go func(k string, v models.GeoAddress, wg *sync.WaitGroup) {
			defer wg.Done()
			successful, failure := p.get(ctx, v)
			if !reflect.DeepEqual(models.Success{}, successful) && !(!partial && successful.Accuracy == "Partial") {
				succMutex.Lock()
				success[k] = successful
				succMutex.Unlock()
			} else {
				if !reflect.DeepEqual(models.Success{}, successful) {
					failure = models.Fail{Code: 2, Message: "Partial Accuraccy Not Enabled"}
				}
				failMutex.Lock()
				fail[k] = v
				f[k] = failure
				failMutex.Unlock()
			}
		}(key, val, &wg)
	}
	wg.Wait()
	return &success, &fail, &f
}

func (p *Pelias) get(ctx context.Context, addr models.GeoAddress) (models.Success, models.Fail) {
	res, err := p.geocode(ctx, addr, 1)
	if err == nil && len(res.Features) > 0 {
		return *p.parse(res), models.Fail{}
	} else {
		p.log.Debug("----Error While Fetching From Pelias----", zap.Error(err))
		return models.Success{}, models.Fail{
			Code:    1,
			Message: "Not Found By Pelias",
		}
	}

}

func (p *Pelias) geocode(ctx context.Context, addr models.GeoAddress, counter int8) (*peliasResponse, error) {
	counter++
	req := p.client.R()
	req.SetContext(ctx)
	req.SetQueryParams(p.getParams(addr))
	resp, err := req.Get(p.url + "/v1/search")
	if err != nil {
		p.log.Error("failed to call Pelias API", zap.Error(err))
		return nil, err
	}
	var r peliasResponse
	err = json.Unmarshal(resp.Body(), &r)
	if err != nil {
		return nil, err
	}
	if (r.Features == nil || len(r.Features) == 0) && strings.Contains(string(resp.Body()), "errors") && counter < 3 {
		return p.geocode(ctx, addr, counter)
	}
	return &r, nil
}

func (p *Pelias) getParams(addr models.GeoAddress) map[string]string {
	params := map[string]string{}
	address := addr.Address
	if addr.Pin != "" {
		if address != "" {
			address += ", "
		}
		address += addr.Pin
	}
	if address != "" {
		params["text"] = address
		params["size"] = "1"
	}
	if addr.Country != "" {
		params["boundary.country"] = addr.Country
	}
	return params
}

//Todo Pelias Accuraccy Should either be accurate or partial
func (p *Pelias) parse(r *peliasResponse) *models.Success {
	prop, err := json.Marshal(r.Features[0].Properties)
	if err != nil {
		p.log.Error("", zap.Error(err))
	}
	return &models.Success{
		Lat:         r.Features[0].Geometry.Coordinates[1],
		Lng:         r.Features[0].Geometry.Coordinates[0],
		Accuracy:    fmt.Sprintf("%f", r.Features[0].Properties.Confidence),
		Source:      "FarEye Maps",
		Address:     r.Features[0].Properties.Label,
		AddressComp: string(prop),
	}
}
