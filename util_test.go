package usps

import "testing"

const (
	testData1 = `hey its me long john johnson and my zip code isn't 123 it's not 560 man it's 95959 and we out here, you already know we are NOT from the 86406 on 5 son`
	testData2 = `hey its me long john my zip is 9-5-9-5 9`
	testData3 = `hey 9 its 5 me 9 long 5 john 9 johnson`
)

func TestZipExtract(t *testing.T) {
	t.Parallel()
	t.Run("testData1", func(t *testing.T) {
		zips, ok := LookupAllZipCodesInString(testData1)
		if !ok {
			t.Fatal("no zips found in test string")
		}
		if len(zips) != 2 {
			t.Fatalf("expected 2 zip codes, got %d", len(zips))
		}
		if zips[0].State != "California" {
			t.Error("expected first zip state to be california")
		}
		if zips[1].State != "Arizona" {
			t.Error("expected first zip state to be Arizona")
		}
		zips2 := ZipExtract(testData1)
		if zips2[0] != zips[0] || zips2[1] != zips[1] {
			t.Error("ZipExtract function is not equal to LookupAllZipCodesInString")
		}
	})
	t.Run("testData2", func(t *testing.T) {
		zips, ok := LookupAllZipCodesInString(testData2)
		if !ok {
			t.Fatal("no zips found in test string")
		}
		if len(zips) != 1 {
			t.Fatalf("expected 1 zip codes, got %d", len(zips))
		}
		if zips[0].State != "California" {
			t.Error("expected first zip state to be california")
		}
		zips2 := ZipExtract(testData2)
		if zips2[0] != zips[0] {
			t.Error("ZipExtract function is not equal to LookupAllZipCodesInString")
		}
	})
	t.Run("testData3", func(t *testing.T) {
		zips, ok := LookupAllZipCodesInString(testData3)
		if ok {
			t.Fatal("zips found in test string")
		}
		if len(zips) != 0 {
			t.Fatalf("expected 0 zip codes, got %d", len(zips))
		}
		zips2 := ZipExtract(testData3)
		if len(zips2) != 0 {
			t.Error("ZipExtract function is not equal to LookupAllZipCodesInString")
		}
	})
}
