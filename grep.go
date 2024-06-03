package usps

import (
	"regexp"
	"strings"
)

var rgx = regexp.MustCompile(`[0-9]{5}(?:-[0-9]{4})?`)

func ExtractZipCodesFromString(s string) ([]string, bool) {
	zips := rgx.FindAllString(s, -1)
	if len(zips) == 0 {
		s = strings.ReplaceAll(s, "-", "")
		s = strings.ReplaceAll(s, " ", "")
		zips = rgx.FindAllString(s, -1)
	}
	return zips, len(zips) > 0
}

func LookupAllZipCodesInString(s string) ([]*ZipCode, bool) {
	zips, ok := ExtractZipCodesFromString(s)
	if !ok {
		return nil, false
	}
	var results []*ZipCode
	for _, z := range zips {
		if zc, ok := zipToZip[z]; ok {
			results = append(results, zc)
		}
	}
	return results, true
}
