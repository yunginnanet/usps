/* [![GoDoc](https://godoc.org/github.com/yunginnanet/usps?status.svg)](https://pkg.go.dev/github.com/yunginnanet/usps) [![codecov](https://codecov.io/gh/yunginnanet/usps/graph/badge.svg?token=YbAagRnkvt)](https://codecov.io/gh/yunginnanet/usps)

#### package usps is a U.S (phone number|state|area|zip) code toolkit.

this toolkit is offline, and the datasets are mostly acquired from U.S. government, all data is compiled into the library.

because of this design choice, almost all functions end up executing in <100 nano seconds. run the benchmarks yourself to see what I mean.

### limitations

at this time, there is no methodology for updating the datasets, and they are bound to become out of date, and may be incomplete.

this should be used for ultra-fast assessment potentially pending further online data queries.

## simple example

	package main

	import "github.com/yunginnanet/usps"

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


*/
package usps
