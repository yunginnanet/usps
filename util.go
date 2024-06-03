package usps

// ZipExtract is just an adapter for LookupAllZipCodesInString to not break things, this was a duplicate function
func ZipExtract(data string) []*ZipCode {
	zips, ok := LookupAllZipCodesInString(data)
	if !ok {
		return []*ZipCode{}
	}
	return zips
}
