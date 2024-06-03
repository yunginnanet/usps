package usps

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"git.tcp.direct/kayos/common/squish"

	"github.com/yunginnanet/usps/internal/data"
)

/*
	  {
		"zip_code": "99632",
		"usps_city": "Mountain Village",
		"stusps_code": "AK",
		"ste_name": "Alaska",
		"zcta": "TRUE",
		"parent_zcta": null,
		"population": 784,
		"density": 12,
		"primary_coty_code": "2158",
		"primary_coty_name": "Kusilvak",
		"county_weights": "{\"02158\": 100}",
		"coty_name": [
		  "Kusilvak"
		],
		"cty_code": [
		  "2158"
		],
		"imprecise": "FALSE",
		"military": "FALSE",
		"timezone": "America/Nome",
		"geo_point_2d": {
		  "lon": -163.66453,
		  "lat": 62.11232
		}
	  },

*/

// map of state to slice of pointers to ZipCodes
var stateToZips = make(map[string]map[string]*ZipCode)

// map of city to slice of pointers to ZipCodes
var cityToZips = make(map[string]map[string]*ZipCode)

// map of zip to pointer to ZipCode
var zipToZip = make(map[string]*ZipCode)

// ZipCode is a struct representing a zipcode, this is based on the USPS dataset.
type ZipCode struct {
	ZipCode         string      `json:"zip_code"`
	City            string      `json:"usps_city"`
	StateShort      string      `json:"stusps_code"`
	State           string      `json:"ste_name"`
	Zcta            string      `json:"zcta"`
	ParentZcta      interface{} `json:"parent_zcta,omitempty"`
	Population      float64     `json:"population,omitempty"`
	Density         float64     `json:"density,omitempty"`
	PrimaryCotyCode string      `json:"primary_coty_code,omitempty"`
	PrimaryCotyName string      `json:"primary_coty_name,omitempty"`
	CountyWeights   string      `json:"county_weights,omitempty"`
	CotyNames       []string    `json:"coty_name,omitempty"`
	CtyCodes        []string    `json:"cty_code,omitempty"`
	Imprecise       string      `json:"imprecise"`
	Military        string      `json:"military"`
	Timezone        string      `json:"timezone"`
	GeoPoint2D      GeoPoint2D  `json:"geo_point_2d,omitempty"`
}

func (z *ZipCode) String() string {
	return z.ZipCode
}

func (z *ZipCode) MustMarshalJSON() []byte {
	res, err := json.Marshal(z)
	if err != nil || z == nil {
		if err == nil {
			err = fmt.Errorf("zip code is nil")
		}
		panic(err)
	}
	return res
}

func (z *ZipCode) Pretty() []byte {
	res, err := json.MarshalIndent(z, "", "\t")
	if err != nil || z == nil {
		if err == nil {
			err = fmt.Errorf("zip code is nil")
		}
		return []byte(`{"error":"` + err.Error() + `"}`)
	}
	return res
}

func (z *ZipCode) Bytes() []byte {
	return z.MustMarshalJSON()
}

type GeoPoint2D struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

var (
	setupOnce = &sync.Once{}
)

func zipSetup() {
	// unmarshal our zipcode data into our zipToZip map
	// and populate our stateToZips and cityToZips maps
	//
	// data is a constant so we yeet errors

	rawZipJSON, _ := squish.UnpackStr(data.X_zipdata)

	var zipSlice []*ZipCode
	defer runtime.GC()
	_ = json.Unmarshal([]byte(rawZipJSON), &zipSlice)

	for _, z := range zipSlice {
		if _, ok := cityToZips[z.City]; !ok {
			cityToZips[z.City] = make(map[string]*ZipCode)
		}
		if _, ok := stateToZips[z.StateShort]; !ok {
			stateToZips[z.StateShort] = make(map[string]*ZipCode)
		}
		cityToZips[z.City][z.ZipCode] = z
		stateToZips[z.StateShort][z.ZipCode] = z
		if _, ok := zipToZip[z.ZipCode]; !ok {
			zipToZip[z.ZipCode] = z
		}
	}
}

func InitZipCodes() {
	setupOnce.Do(zipSetup)
}

// GetZipCode returns a pointer to a ZipCode struct when given a zip code in string form.
// Data is loaded from the USPS dataset on first call and cached in memory.
func GetZipCode(zip string) (*ZipCode, bool) {
	InitZipCodes()
	z, ok := zipToZip[zip]
	return z, ok
}

// GetStateZips returns a slice of pointers to ZipCode structs associated with a state when given a state in string form.
// Data is loaded from the USPS dataset on first call and cached in memory.
func GetStateZips(state string) []*ZipCode {
	InitZipCodes()
	var zs = make([]*ZipCode, len(stateToZips[state]))
	for _, zip := range stateToZips[state] {
		zs = append(zs, zip)
	}
	return zs
}

// GetStateZipStrings returns a slice of zip codes in string form associated with a state when given a state in string form.
// Data is loaded from the USPS dataset on first call and cached in memory.
func GetStateZipStrings(state string) []string {
	InitZipCodes()
	results := stateToZips[state]
	if len(results) == 0 {
		return nil
	}
	var zs = make([]string, 0, len(results))
	for zip := range results {
		zs = append(zs, zip)
	}
	return zs
}

// GetCityZips returns a slice of pointers to ZipCode structs associated with a city when given a city in string form.
// Data is loaded from the USPS dataset on first call and cached in memory.
func GetCityZips(city string) []*ZipCode {
	InitZipCodes()
	var zs = make([]*ZipCode, 0, len(cityToZips[city]))
	for _, zip := range cityToZips[city] {
		zs = append(zs, zip)
	}
	return zs
}

// GetCityZipStrings returns a slice of zip codes in string form associated with a city when given a city in string form.
// Data is loaded from the USPS dataset on first call and cached in memory.
func GetCityZipStrings(city string) []string {
	InitZipCodes()
	var zs = make([]string, len(cityToZips[city]))
	for zip := range cityToZips[city] {
		zs = append(zs, zip)
	}
	return zs
}
