package usps

import "testing"

func TestShortHand(t *testing.T) {
	t.Parallel()
	if len(shortToLong) != len(longToShort) {
		t.Fatal("shortHand and longToShort are not the same length")
	}

}
