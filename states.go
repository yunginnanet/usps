package usps

import (
	"encoding/json"
	"runtime"
	"slices"
	"strings"
	"sync"
	"unicode"

	"git.tcp.direct/kayos/common/squish"

	"github.com/yunginnanet/usps/internal/data"
)

var allStates []State

type State struct {
	State  string `json:"state"`
	Abbrev string `json:"abbrev"`
	Short  string `json:"code"`
}

var (
	stateShortToState = make(map[string]State)
	stateLongToState  = make(map[string]State)
	stateSetupOnce    = &sync.Once{}
)

func stateSetup() {
	// dataset is a constant, cannot fail
	// so we do a lot of error yeeting here

	dat, _ := squish.Gunzip(squish.B64d(data.X_statesjson))

	defer runtime.GC()
	_ = json.Unmarshal(dat, &allStates)

	for _, s := range allStates {
		stateShortToState[s.Short] = s
		stateLongToState[s.State] = s
		stateLongToState[strings.ToUpper(s.State)] = s
		stateLongToState[strings.ToLower(s.State)] = s
		// stateCodeToState[s.Code] = s
	}
}

func InitStates() {
	stateSetupOnce.Do(stateSetup)
}

var emptyState = State{
	State:  "",
	Abbrev: "",
	Short:  "",
}

func StateByShort(short string) State {
	InitStates()

	res, ok := stateShortToState[short]
	if !ok {
		res, ok = stateShortToState[strings.ToUpper(short)]
		if !ok {
			return emptyState
		}
	}

	return res
}

func StateByLong(long string) State {
	InitStates()

	res, ok := stateLongToState[long]
	if !ok {
		// remove all non-alphanumeric characters
		long = strings.ToUpper(
			strings.TrimSpace(
				strings.Map(func(r rune) rune {
					if r == '-' || r == '_' {
						return ' '
					}
					if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
						return -1
					}

					return r
				}, long),
			),
		)
		res, ok = stateLongToState[long]
		if !ok {
			return emptyState
		}
	}

	return res
}

func StateToShort(long string) string {
	if s := StateByLong(long); s.State != "" && s != emptyState {
		return s.Short
	}
	if s := StateByShort(long); s.Short != "" && s != emptyState {
		return s.Short
	}
	return ""
}

func StateList() []State {
	InitStates()

	return slices.Clone(allStates)
}

func StateListShort() []string {
	InitStates()
	var res []string
	for _, s := range allStates {
		res = append(res, s.Short)
	}

	return res
}

func StateListLong() []string {
	InitStates()
	var res []string
	for _, s := range allStates {
		res = append(res, s.State)
	}

	return res
}
