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
	return strings.Contains(value.(string), sm.searchValue)
}
