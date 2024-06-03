package usps

import (
	"bytes"
	"sync"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
	"github.com/davecgh/go-spew/spew"
)

func TestZipCodes(t *testing.T) {
	t.Parallel()
	InitZipCodes()
	if len(zipToZip) == 0 {
		t.Fatal("zipToZip is empty")
	}
	if len(stateToZips) == 0 {
		t.Fatal("stateToZips is empty")
	}
	if len(cityToZips) == 0 {
		t.Fatal("cityToZips is empty")
	}
	zip, ok := GetZipCode("10001")
	if zip == nil || !ok {
		t.Fatal("failed zip lookup")
	}
	spew.Dump(zip)
	if zip.ZipCode != "10001" {
		t.Fatalf("fetched 10001, but zip is not 10001: %v", zip.ZipCode)
	}
	if zip.City != "New York" {
		t.Fatalf("fetched 10001, but city is not New York: %v", zip.City)
	}
	if zip.StateShort != "NY" {
		t.Fatalf("fetched 10001, but state is not New York: %v", zip.StateShort)
	}
	t.Log(zip.Pretty())
	t.Log(string(zip.MustMarshalJSON()))
	t.Log(string(zip.Bytes()))
	czips := GetCityZips("New York")
	if len(czips) == 0 {
		t.Fatal("failed city zip lookup")
	}
	czipStrings := GetCityZipStrings("New York")
	if len(czipStrings) == 0 {
		t.Fatal("failed city zip string lookup")
	}
	szips := GetStateZips("NY")
	if len(szips) == 0 {
		t.Fatal("failed state zip lookup")
	}
}

func TestZipCodesMustFail(t *testing.T) {
	t.Parallel()
	InitZipCodes()
	zip, ok := GetZipCode("00000")
	if zip != nil || ok {
		t.Fatal("failed zip lookup")
	}
	t.Run("panic1", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		asdf := zip.MustMarshalJSON()
		t.Log(asdf)
	})
	t.Run("panic2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		asdf := zip.Bytes()
		t.Log(asdf)
	})
	t.Run("json_error", func(t *testing.T) {
		if res := zip.Pretty(); !bytes.Contains(res, []byte("error")) {
			t.Fatal("expected error in pretty print of nil zip")
		}
	})
	t.Run("bogus_longform", func(t *testing.T) {
		state := StateByLong(string([]byte{0x00, 0x05, 0xa0, 0x55, 0xcc}))
		if state != emptyState {
			t.Fatal("expected nil state")
		}
	})
}

func TestGetStateZips_EmptyMap(t *testing.T) {
	stateToZips = make(map[string]map[string]*ZipCode) // Ensure the map is empty
	zips := GetStateZips("NY")
	if len(zips) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(zips))
	}
}

func TestGetStateZips_NonExistentState(t *testing.T) {
	stateToZips = map[string]map[string]*ZipCode{
		"CA": {},
	}
	zips := GetStateZips("NY") // NY does not exist in the map
	if len(zips) != 0 {
		t.Errorf("expected empty slice for non-existent state, got %d elements", len(zips))
	}
}

func TestGetStateZips_StateExistsNoZips(t *testing.T) {
	stateToZips = map[string]map[string]*ZipCode{
		"NY": {},
	}
	zips := GetStateZips("NY") // NY exists but has no zip codes
	if len(zips) != 0 {
		t.Errorf("expected empty slice for state with no zips, got %d elements", len(zips))
	}
}

func TestGetStateZips_EmptyStateString(t *testing.T) {
	stateToZips = map[string]map[string]*ZipCode{
		"NY": {
			"10001": &ZipCode{ZipCode: "10001", City: "New York", StateShort: "NY"},
		},
	}
	zips := GetStateZips("")
	if len(zips) != 0 {
		t.Errorf("expected empty slice for empty state string, got %d elements", len(zips))
	}
}

func TestGetStateZipStrings_NoZipCodes(t *testing.T) {
	InitZipCodes()
	stateToZips["TestState"] = make(map[string]*ZipCode)
	zips := GetStateZipStrings("TestState")
	if zips != nil {
		t.Errorf("Expected nil, got %v", zips)
	}
}

func TestGetStateZipStrings_MultipleZipCodes(t *testing.T) {
	InitZipCodes()
	stateToZips["TestState"] = map[string]*ZipCode{
		"12345": {},
		"67890": {},
	}
	expected := map[string]bool{"12345": false, "67890": false}
	testZips := GetStateZipStrings("TestState")
	if testZips == nil {
		t.Errorf("Expected non-nil, got nil")
	}
	if len(testZips) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(testZips))
	}
	for _, z := range testZips {
		if found, ok := expected[z]; !ok {
			t.Errorf("Unexpected zip code %s", z)
		} else if found {
			t.Errorf("Duplicate zip code %s", z)
		}
		expected[z] = true
	}
}

func TestGetStateZipStrings_EmptyStateString(t *testing.T) {
	InitZipCodes()
	zips := GetStateZipStrings("")
	if zips != nil {
		t.Errorf("Expected nil for empty state string, got %v", zips)
	}
}

func TestGetStateZipStrings_InvalidState(t *testing.T) {
	InitZipCodes()
	zips := GetStateZipStrings("InvalidState")
	if zips != nil {
		t.Errorf("Expected nil for invalid state, got %v", zips)
	}
}

var once = &sync.Once{}
var zips []string

func init() {
	once.Do(func() {
		InitZipCodes()
		for _, z := range zipToZip {
			zips = append(zips, z.ZipCode)
		}
		// shuffle the zips
		for i := range zips {
			j := entropy.GetSharedRand().Intn(i + 1)
			zips[i], zips[j] = zips[j], zips[i]
		}
	})
}

func BenchmarkZipCodes(b *testing.B) {
	if len(zips) == 0 {
		panic("zips is empty")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		wg.Add(len(zips))
		for zipIndex := range zips {
			go func(i int) {
				_, ok := GetZipCode(zips[i])
				if !ok {
					b.Error("failed zip lookup")
				}
				wg.Done()
			}(zipIndex)
		}
		wg.Wait()
	}
	b.StopTimer()
	b.Logf("fetched the details for %d zipcodes %d times in %s", b.Elapsed(), b.N, b.Elapsed())
	b.Run("GetSingleZipCode", func(b *testing.B) {
		if len(zips) == 0 {
			panic("zips is empty")
		}
		for i := 0; i < b.N; i++ {
			_, ok := GetZipCode(zips[0])
			if !ok {
				b.Fatal("failed zip lookup")
			}
		}
	})
}
