package usps

import (
	"bufio"
	"strconv"
	"strings"
	"sync"

	"git.tcp.direct/kayos/common/squish"
	"github.com/google/shlex"

	"usps/internal/data"
)

const AreaCodeLength = 3

// AreaCode is a struct representing an area code.
type AreaCode struct {
	Code, State, City, Country string
	Latitude, Longitude        float64
}

var (
	areaCodeSetupOnce      = &sync.Once{}
	areaCodeToAreaCodes    = make(map[string][]*AreaCode)
	intAreaCodeToAreaCodes = make(map[int][]*AreaCode)
	intAreaCodeToState     = make(map[int]string)
	areaCodeToState        = make(map[string]string)
	stateToAreaCodes       = make(map[string][]*AreaCode)
	cityToAreaCodes        = make(map[string][]*AreaCode)
)

func resetAllAreaCodes() {
	areaCodeToAreaCodes = make(map[string][]*AreaCode)
	intAreaCodeToAreaCodes = make(map[int][]*AreaCode)
	intAreaCodeToState = make(map[int]string)
	areaCodeToState = make(map[string]string)
	stateToAreaCodes = make(map[string][]*AreaCode)
	cityToAreaCodes = make(map[string][]*AreaCode)
	areaCodeSetupOnce = &sync.Once{}
}

var cleaner = strings.NewReplacer(
	"'", "",
	"\"", "", "-", "", ",",
	"", ".", "", "+", "",
	"*", "", "#", "", "(",
	"", ")", "", "!", "",
)

func clean(s string) string {
	return cleaner.Replace(
		strings.ToUpper(
			strings.TrimSpace(s),
		),
	)
}

func areaCodeInit() {
	var dat string

	// dataset is a constant, cannot fail
	// so we do a lot of error yeeting here

	dat, _ = squish.UnpackStr(data.XAreacodedata)
	xerox := bufio.NewScanner(strings.NewReader(dat))
	for xerox.Scan() {

		var tokens []string

		tokens, _ = shlex.Split(strings.ReplaceAll(xerox.Text(), ",", " "))

		latitude, _ := strconv.ParseFloat(tokens[4], 32)

		longitude, _ := strconv.ParseFloat(tokens[5], 32)

		acode := &AreaCode{
			Code: clean(tokens[0]), City: clean(tokens[1]),
			State: clean(tokens[2]), Country: clean(tokens[3]),
			Latitude: latitude, Longitude: longitude,
		}

		intAc, _ := strconv.Atoi(acode.Code)

		intAreaCodeToAreaCodes[intAc] = append(intAreaCodeToAreaCodes[intAc], acode)
		intAreaCodeToState[intAc] = acode.State
		areaCodeToAreaCodes[acode.Code] = append(areaCodeToAreaCodes[acode.Code], acode)
		stateToAreaCodes[acode.State] = append(stateToAreaCodes[acode.State], acode)
		cityToAreaCodes[acode.City] = append(cityToAreaCodes[acode.City], acode)
		areaCodeToState[acode.Code] = acode.State
	}
	// if we haven't panic'd at this point then our dataset is good
}

func InitAreaCodes() {
	areaCodeSetupOnce.Do(areaCodeInit)
}

// LookupAreaCode returns a slice of AreaCode types that match the given area code.
func LookupAreaCode(code string) ([]*AreaCode, bool) {
	InitAreaCodes()
	res, ok := areaCodeToAreaCodes[code]
	if !ok {
		res, ok = areaCodeToAreaCodes[clean(code)]
	}
	return res, ok
}

// LookupStateAreaCodes returns a slice of AreaCode types that match the given state.
func LookupStateAreaCodes(state string) ([]*AreaCode, bool) {
	InitAreaCodes()
	if len(state) == 2 {
		state = LongHand(state)
	}
	res, ok := stateToAreaCodes[clean(state)]
	return res, ok
}

// LookupCityAreaCodes returns a slice of AreaCode types that match the given city.
// Note! This is not a perfect match. If a city is not found, it will attempt to find a city that contains the query..
func LookupCityAreaCodes(city string) ([]*AreaCode, bool) {
	InitAreaCodes()
	city = clean(city)
	res, ok := cityToAreaCodes[city]
	if !ok {
		for k, v := range cityToAreaCodes {
			if strings.Contains(k, city) {
				res = append(res, v...)
				ok = true
			}
		}
	}
	return res, ok
}

// PhoneNumberState returns the state that the given phone number is in.
// This is probably the most useful function in this package, and the one you're looking for.
// This function will return an empty string if the number is invalid.
func PhoneNumberState(number string) string {
	InitAreaCodes()
	ac := PhoneNumberAreaCode(clean(number))
	if as, ok := areaCodeToState[ac]; ok {
		return clean(as)
	}
	return ""
}

func PhoneNumberCity(number string) string {
	InitAreaCodes()
	ac := PhoneNumberAreaCode(clean(number))
	a, ok := LookupAreaCode(ac)
	if len(a) > 0 && ok {
		return clean(a[0].City)
	}
	return ""
}

// GetAreaCodeInt returns the area code of the given phone number or string as an int. If invalid, returns 0.
func GetAreaCodeInt(number string) int {
	switch {
	case len(number) > AreaCodeLength:
		number = PhoneNumberAreaCode(number)
	case len(number) < AreaCodeLength:
		return 0
	case len(number) == AreaCodeLength:
		//
	}
	n, ok := strconv.Atoi(number)
	if ok != nil {
		return 0
	}
	return n
}

// PhoneNumberAreaCode returns the area code of the given phone number.
// This function will return an empty string if the number is invalid or isn't a known USA area code.
func PhoneNumberAreaCode(number string) string {
	if len(number) < 10 {
		return ""
	}
	if len(number) > 10 {
		number = clean(strings.ReplaceAll(number, " ", ""))
	}
	if len(number) == 11 && string(number[0]) == "1" {
		number = number[1:]
	}
	if len(number) != 10 {
		return ""
	}
	return number[:AreaCodeLength]
}
