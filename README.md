# usps

[![GoDoc](https://godoc.org/github.com/yunginnanet/usps?status.svg)](https://pkg.go.dev/github.com/yunginnanet/usps) [![codecov](https://codecov.io/gh/yunginnanet/usps/branch/main/graph/badge.svg)](https://codecov.io/gh/yunginnanet/usps)

#### package usps is a U.S (phone number|state|area|zip) code toolkit.

this toolkit is offline, and the datasets are mostly acquired from U.S.
government, all data is compiled into the library.

because of this design choice, almost all functions end up executing in <100
nano seconds. run the benchmarks yourself to see what I mean.

### limitations

at this time, there is no methodology for updating the datasets, and they are
bound to be inaccurate.

this should be used for ultra-fast assessment potentially pending further online
data queries.

## simple example

    import "github.com/yunginnanet/usps"

    package main

    import "usps"

    func main() {
    	entries, _ := usps.LookupCityAreaCodes("CHICAGO")
    	for _, entry := range entries {
    		println(entry.City)
    		println(entry.Code)
    		println(entry.State)
    		println(entry.Country)
    		println(entry.Latitude, entry.Longitude)
    	}
    }

### output:

    EAST CHICAGO
    219
    INDIANA
    US
    +4.163920e+001 -8.745476e+001
    WEST CHICAGO
    312
    ILLINOIS
    US
    +4.188475e+001 -8.820396e+001
    WEST CHICAGO
    331
    ILLINOIS
    US
    +4.188475e+001 -8.820396e+001
    WEST CHICAGO
    630
    ILLINOIS
    US
    +4.188475e+001 -8.820396e+001
    WEST CHICAGO
    773
    ILLINOIS
    US
    +4.188475e+001 -8.820396e+001
    CHICAGO HEIGHTS
    708
    ILLINOIS
    US
    +4.150615e+001 -8.763560e+001
    NORTH CHICAGO
    224
    ILLINOIS
    US
    +4.232558e+001 -8.784118e+001
    NORTH CHICAGO
    847
    ILLINOIS
    US
    +4.232558e+001 -8.784118e+001


```go
const AreaCodeLength = 3
```

#### func  ExtractZipCodesFromString

```go
func ExtractZipCodesFromString(s string) ([]string, bool)
```

#### func  GetAreaCodeInt

```go
func GetAreaCodeInt(number string) int
```
GetAreaCodeInt returns the area code of the given phone number or string as an
int. If invalid, returns 0.

#### func  GetCityZipStrings

```go
func GetCityZipStrings(city string) []string
```
GetCityZipStrings returns a slice of zip codes in string form associated with a
city when given a city in string form. Data is loaded from the USPS dataset on
first call and cached in memory.

#### func  GetStateZipStrings

```go
func GetStateZipStrings(state string) []string
```
GetStateZipStrings returns a slice of zip codes in string form associated with a
state when given a state in string form. Data is loaded from the USPS dataset on
first call and cached in memory.

#### func  InitAreaCodes

```go
func InitAreaCodes()
```

#### func  InitStates

```go
func InitStates()
```

#### func  InitZipCodes

```go
func InitZipCodes()
```

#### func  LongHand

```go
func LongHand(short string) string
```

#### func  PhoneNumberAreaCode

```go
func PhoneNumberAreaCode(number string) string
```
PhoneNumberAreaCode returns the area code of the given phone number. This
function will return an empty string if the number is invalid or isn't a known
USA area code.

#### func  PhoneNumberCity

```go
func PhoneNumberCity(number string) string
```

#### func  PhoneNumberState

```go
func PhoneNumberState(number string) string
```
PhoneNumberState returns the state that the given phone number is in. This is
probably the most useful function in this package, and the one you're looking
for. This function will return an empty string if the number is invalid.

#### func  ShortHand

```go
func ShortHand(long string) string
```

#### func  StateListLong

```go
func StateListLong() []string
```

#### func  StateListShort

```go
func StateListShort() []string
```

#### func  StateToShort

```go
func StateToShort(long string) string
```

#### type AreaCode

```go
type AreaCode struct {
	Code, State, City, Country string
	Latitude, Longitude        float64
}
```

AreaCode is a struct representing an area code.

#### func  LookupAreaCode

```go
func LookupAreaCode(code string) ([]*AreaCode, bool)
```
LookupAreaCode returns a slice of AreaCode types that match the given area code.

#### func  LookupCityAreaCodes

```go
func LookupCityAreaCodes(city string) ([]*AreaCode, bool)
```
LookupCityAreaCodes returns a slice of AreaCode types that match the given city.
Note! This is not a perfect match. If a city is not found, it will attempt to
find a city that contains the query..

#### func  LookupStateAreaCodes

```go
func LookupStateAreaCodes(state string) ([]*AreaCode, bool)
```
LookupStateAreaCodes returns a slice of AreaCode types that match the given
state.

#### type GeoPoint2D

```go
type GeoPoint2D struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
```


#### type State

```go
type State struct {
	State  string `json:"state"`
	Abbrev string `json:"abbrev"`
	Short  string `json:"code"`
}
```


#### func  StateByLong

```go
func StateByLong(long string) State
```

#### func  StateByShort

```go
func StateByShort(short string) State
```

#### func  StateList

```go
func StateList() []State
```

#### type ZipCode

```go
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
```

ZipCode is a struct representing a zipcode, this is based on the USPS dataset.

#### func  GetCityZips

```go
func GetCityZips(city string) []*ZipCode
```
GetCityZips returns a slice of pointers to ZipCode structs associated with a
city when given a city in string form. Data is loaded from the USPS dataset on
first call and cached in memory.

#### func  GetStateZips

```go
func GetStateZips(state string) []*ZipCode
```
GetStateZips returns a slice of pointers to ZipCode structs associated with a
state when given a state in string form. Data is loaded from the USPS dataset on
first call and cached in memory.

#### func  GetZipCode

```go
func GetZipCode(zip string) (*ZipCode, bool)
```
GetZipCode returns a pointer to a ZipCode struct when given a zip code in string
form. Data is loaded from the USPS dataset on first call and cached in memory.

#### func  LookupAllZipCodesInString

```go
func LookupAllZipCodesInString(s string) ([]*ZipCode, bool)
```

#### func  ZipExtract

```go
func ZipExtract(data string) []*ZipCode
```
ZipExtract is just an adapter for LookupAllZipCodesInString to not break things,
this was a duplicate function

#### func (*ZipCode) Bytes

```go
func (z *ZipCode) Bytes() []byte
```

#### func (*ZipCode) MustMarshalJSON

```go
func (z *ZipCode) MustMarshalJSON() []byte
```

#### func (*ZipCode) Pretty

```go
func (z *ZipCode) Pretty() []byte
```

#### func (*ZipCode) String

```go
func (z *ZipCode) String() string
```

---
