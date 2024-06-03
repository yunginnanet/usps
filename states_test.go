package usps

import "testing"

func TestStateByShort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		short string
		want  string
	}{
		{"ny", "NY", "New York"},
		{"ca", "CA", "California"},
		{"tx", "TX", "Texas"},
		{"il", "IL", "Illinois"},
		{"fl", "FL", "Florida"},
		{"wa", "WA", "Washington"},
		{"or", "OR", "Oregon"},
		{"az", "AZ", "Arizona"},
		{"nv", "NV", "Nevada"},
		{"invalid", "XX", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StateByShort(tt.short); got.State != tt.want {
				t.Errorf("StateByShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateByLong(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		long string
		want string
	}{
		{"ny", "New York", "NY"},
		{"ny", "New-York", "NY"},
		{"ny", "NeW yOrK", "NY"},
		{"ca", "California", "CA"},
		{"tx", "Texas", "TX"},
		{"il", "Illinois", "IL"},
		{"fl", "Florida", "FL"},
		{"wa", "Washington", "WA"},
		{"or", "Oregon", "OR"},
		{"az", "Arizona", "AZ"},
		{"nv", "Nevada", "NV"},
		{"invalid", "Invalid", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StateByLong(tt.long); got.Short != tt.want {
				t.Errorf("StateByLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateList(t *testing.T) {
	t.Parallel()
	longList := StateListLong()
	shortList := StateListShort()
	list := StateList()
	// codeList := StateListCode()
	if len(longList) != len(shortList) || len(list) != len(shortList) || len(shortList) != 51 {
		t.Errorf("StateList() = %v, want %v", len(longList), len(shortList))
	}
}

func TestStateToShort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		long string
		want string
	}{
		{"ny", "New York", "NY"},
		{"ny2", "NY", "NY"},
		{"ca", "California", "CA"},
		{"ca2", "CA", "CA"},
		{"tx", "Texas", "TX"},
		{"tx2", "TX", "TX"},
		{"il", "Illinois", "IL"},
		{"il2", "IL", "IL"},
		{"fl", "Florida", "FL"},
		{"fl2", "FL", "FL"},
		{"wa", "Washington", "WA"},
		{"wa2", "WA", "WA"},
		{"or", "Oregon", "OR"},
		{"or2", "OR", "OR"},
		{"az", "Arizona", "AZ"},
		{"az2", "AZ", "AZ"},
		{"nv", "Nevada", "NV"},
		{"nv2", "NV", "NV"},
		{"bogus", "Yeet", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StateToShort(tt.long); got != tt.want {
				t.Errorf("StateToShort() = %v, want %v", got, tt.want)
			}
		})
	}
}
