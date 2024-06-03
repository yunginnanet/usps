package usps

import (
	"strings"
	"testing"

	"git.tcp.direct/kayos/common/entropy"
)

func TestAreaCodes(t *testing.T) {
	t.Parallel()
	resetAllAreaCodes()
	entries, ok := LookupAreaCode("702")
	if !ok {
		t.Fatal("failed area code lookup")
	}
	for _, entry := range entries {
		switch {
		case entry.Code != "702":
			t.Fatalf("fetched 702, but code is not 702: %v", entry.Code)
		case entry.City != "LAS VEGAS" && entry.City != "NORTH LAS VEGAS" &&
			entry.City != "HENDERSON" && entry.City != "PARADISE" &&
			entry.City != "SPRING VALLEY" && entry.City != "SUNRISE MANOR" &&
			entry.City != "WINCHESTER":
			t.Fatalf("fetched 702, but city is not Las Vegas: %v", entry.City)
		case entry.State != "NEVADA":
			t.Fatalf("fetched 702, but state is not Nevada: %v", entry.State)
		}
	}
	var res string
	testNum := "702-005-5555"
	if res = PhoneNumberState(testNum); res != "NEVADA" {
		t.Fatalf("failed phone number state lookup, got: %s", res)
	}
	if short := ShortHand(res); short != "NV" {
		t.Fatalf("shorthand failed, got: %s", ShortHand(res))
	}
	stateRes := stateToAreaCodes[res]
	entropy.GetOptimizedRand().Shuffle(len(stateRes), func(i, j int) {
		stateRes[i], stateRes[j] = stateRes[j], stateRes[i]
	})
	s := &strings.Builder{}
	s.WriteString(stateRes[0].Code)
	s.WriteString(testNum[3:])

	if tried := ShortHand(PhoneNumberState(s.String())); tried != "NV" {
		t.Fatalf("shorthand failed, got: %s", tried)
	}

	if PhoneNumberState("7020055555") != LongHand("NV") {
		t.Fatalf("failed phone number state lookup, got area code: %s", PhoneNumberAreaCode("7020055555"))
	}

	if PhoneNumberState("yeeties") != "" {
		t.Fatalf("somehow got a state from some bogus string")
	}

	if PhoneNumberCity("yeeties") != "" {
		t.Fatalf("somehow got a city from some bogus string")
	}

	if GetAreaCodeInt("5") != 0 {
		t.Fatalf("somehow got an area code from a single digit string")
	}

	if PhoneNumberAreaCode("+1 (123) 456-7890") != "123" {
		t.Errorf("failed to clean long phone number string")
	}

	if PhoneNumberAreaCode(strings.Repeat("5", 55)) != "" {
		t.Errorf("somehow got an area code from a 55 digit string")
	}

	if PhoneNumberAreaCode("12345678901") != "234" {
		t.Errorf("failed to clean preceeding '1' in phone number string")
	}
}

func TestLookupCityAreaCodes(t *testing.T) {
	t.Parallel()
	city := "LAS VEGAS"
	entries, ok := LookupCityAreaCodes(city)
	if !ok {
		t.Fatal("failed city area code lookup")
	}
	for _, entry := range entries {
		switch {
		case !strings.Contains(entry.City, city):
			t.Fatalf("fetched %s, but city is not %s: %v", city, city, entry.City)
		case entry.State != "NEVADA":
			t.Fatalf("fetched %s, but state is not Nevada: %v", city, entry.State)
		}
	}
}

func TestGetAreaCodeInt(t *testing.T) {
	t.Parallel()
	if GetAreaCodeInt("702") != 702 {
		t.Errorf("failed to get area code int, got %d", GetAreaCodeInt("702"))
	}
}

func TestGetAreaCodeBogus(t *testing.T) {
	t.Parallel()
	if GetAreaCodeInt("bogus") != 0 {
		t.Errorf("somehow got an area code from some bogus string")
	}
	if acodes, err := LookupAreaCode("yeet"); len(acodes) != 0 || err != false {
		t.Errorf("somehow got an area code from some bogus string")
	}
}

func TestLookupStateAreaCodes(t *testing.T) {
	t.Parallel()
	state := "UT"
	entries, ok := LookupStateAreaCodes(state)
	if !ok {
		t.Fatal("failed state area code lookup")
	}
	for _, entry := range entries {
		switch {
		case entry.State != "UTAH":
			t.Fatalf("fetched %s, but state is not Utah: %v", state, entry.State)
		}
	}
}

func TestPhoneNumberCity(t *testing.T) {
	t.Parallel()
	if PhoneNumberCity("7020055555") != "HENDERSON" {
		t.Fatalf("failed phone number city lookup, got: %s", PhoneNumberCity("7020055555"))
	}
}

func BenchmarkAreaCodeInitialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resetAllAreaCodes()
		InitAreaCodes()
	}
}

func BenchmarkAreaCodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LookupAreaCode("702")
	}
}
