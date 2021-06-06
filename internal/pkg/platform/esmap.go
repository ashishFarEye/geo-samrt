package platform

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"time"

	"geo-smart/api/gen/models"
	elastic "github.com/olivere/elastic/v7"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"go.uber.org/zap"
)

type ES struct {
	client *elastic.Client
	log    *zap.Logger
}

type esResp struct {
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	Landmark    string `json:"landmark"`
	Pincode     string `json:"pincode"`
	CountryCode string `json:"countryCode"`
	UUID        string `json:"uuid"`
	Address     string `json:"address"`
	Accuracy    string `json:"accuracy"`
	Count       int64  `json:"count"`
	Location    struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func NewES(host string, user string, pwd string, log *zap.Logger) *ES {
	es, err := elastic.NewClient(
		elastic.SetBasicAuth(user, pwd),
		elastic.SetURL(host), elastic.SetHealthcheckTimeout(time.Second*10), elastic.SetSniff(false))
	if err != nil {
		log.Fatal("---Error While Making Connections With ES---", zap.Error(err))
	}
	return &ES{
		client: es,
		log:    log,
	}
}

func (e *ES) MapsGeocoder(ctx context.Context, addr map[string]models.GeoAddress, partial bool) (*map[string]models.Success, *map[string]models.GeoAddress, *map[string]models.Fail) {
	success, fail, succMutex, failMutex, f := map[string]models.Success{}, map[string]models.GeoAddress{}, sync.RWMutex{}, sync.RWMutex{}, map[string]models.Fail{}
	e.log.Info("---Sending Request To FarEye Maps---", zap.Reflect("Address", &addr))
	var wg sync.WaitGroup
	for key, val := range addr {
		wg.Add(1)
		go func(k string, v models.GeoAddress, wg *sync.WaitGroup) {
			defer wg.Done()
			successful, failure := e.get(ctx, v)
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

func (e *ES) get(ctx context.Context, addr models.GeoAddress) (models.Success, models.Fail) {
	res, err := e.geocode(ctx, addr)
	if err == nil && !reflect.DeepEqual(esResp{}, *res) {
		return *e.parse(res), models.Fail{}
	} else {
		e.log.Debug("----Error While Fetching From Elasticsearch----", zap.Error(err))
		return models.Success{}, models.Fail{
			Code:    1,
			Message: "Not Found By FarEye Maps",
		}
	}
}

func (e *ES) getDataFromES(ctx context.Context, addr models.GeoAddress, q elastic.BoolQuery) (*elastic.SearchResult, error) {
	ss := elastic.NewSearchSource()
	ss.Query(&q).From(0).Size(100)
	// ssSource, err := ss.Source()
	// data, _ := json.MarshalIndent(ssSource, "", "  ")
	// fmt.Printf("%s\n", string(data))

	resp, err := e.client.Search("dna-" + strings.ToLower(addr.Country)).SearchSource(ss).Pretty(true).Do(context.Background())
	return resp, err
}

func (e *ES) geocode(ctx context.Context, addr models.GeoAddress) (*esResp, error) {
	result := esResp{}
	q := e.getExactMatchQuery(addr)
	resp, err := e.getDataFromES(ctx, addr, q)
	if err != nil {
		e.log.Error("---Error While Fetching Exact Match from ES---", zap.Error(err))
		return &result, err
	}
	if resp.Hits.TotalHits.Value <= 0 {
		if addr.UUID != "" {
			q = e.getUUIDSearchQuery(addr)
			resp, err = e.getDataFromES(ctx, addr, q)
			if err != nil {
				e.log.Error("---Error While Fetching Based on UUID From ES---", zap.Error(err))
				return &result, err
			}
			if resp.Hits.TotalHits.Value > 0 {
				return e.getFuzzyMatchResult(resp, addr)
			}
		}
		q = e.getAllDataSearchQuery(addr)
		resp, err = e.getDataFromES(ctx, addr, q)
		if err != nil {
			e.log.Error("---Error While Fetching With Fuzzy logic From ES---", zap.Error(err))
			return &result, err
		}

	}
	if resp.Hits.TotalHits.Value > 0 {
		json.Unmarshal(resp.Hits.Hits[0].Source, &result)
		if result.Count > 1 {
			return &result, nil
		} else {
			return &esResp{}, nil
		}
	}
	return &result, nil

}

func (e *ES) getFuzzyMatchResult(resp *elastic.SearchResult, addr models.GeoAddress) (*esResp, error) {
	//var ttyp esResp
	maxScore := 0
	resAddr := esResp{}
	for _, hit := range resp.Hits.Hits {
		var t esResp
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			e.log.Error("Error:", zap.Error(err))
		} else {
			score := fuzzy.TokenSortRatio(t.Address, addr.Address)
			if score > maxScore {
				maxScore = score
				resAddr = t
			}
		}
	}

	return &resAddr, nil
}

func (e *ES) getExactMatchQuery(a models.GeoAddress) elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if a.Address != "" {
		query.Must(elastic.NewTermQuery("address.keyword", a.Address))
	}
	if a.Pin != "" {
		query.Must(elastic.NewMatchQuery("pincode", a.Pin)).MinimumShouldMatch("100%")
	}
	if a.UUID != "" {
		query.Must(elastic.NewMatchQuery("uuid", a.UUID)).MinimumShouldMatch("100%")
	}
	return *query
}

func (e *ES) getAllDataSearchQuery(a models.GeoAddress) elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if a.Address != "" {
		query.Must(elastic.NewMatchQuery("address", a.Address).Fuzziness("AUTO").Operator("OR")).MinimumShouldMatch("80%")
	}
	if a.Pin != "" {
		query.Must(elastic.NewMatchQuery("pincode", a.Pin)).MinimumShouldMatch("80%")
	}
	return *query
}

func (e *ES) getUUIDSearchQuery(a models.GeoAddress) elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewMatchQuery("uuid", a.UUID)).MinimumShouldMatch("100%")
	return *query
}

func (e *ES) parse(resp *esResp) *models.Success {
	return &models.Success{
		Lat:         resp.Location.Lat,
		Lng:         resp.Location.Lon,
		Accuracy:    e.getAccuraccy(resp),
		Source:      "FarEye Maps",
		Address:     resp.Address,
		AddressComp: "",
	}
}

//Todo Figure Out How To Mark Accurate & Partial
func (e *ES) getAccuraccy(resp *esResp) string {
	return "Accurate"
}
