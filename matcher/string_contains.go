package matcher

import "strings"

type stringContains struct {
	searchValue string
}

func StringContains(search string) Generic {
	obj := new(stringContains)
	obj.searchValue = search
	return obj
}

func (sm *stringContains) Match(value interface{}) bool {
	if castedString,ok := value.(string); !ok {
		return false //not a string
	} else {
		return strings.Contains(castedString, sm.searchValue)
	}
}
