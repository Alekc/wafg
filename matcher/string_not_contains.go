package matcher

import "strings"

type stringDoesntContain struct {
	searchValue string
}

func StringDoesntContain(search string) Generic {
	obj := new(stringDoesntContain)
	obj.searchValue = search
	return obj
}

func (sm *stringDoesntContain) Match(value interface{}) bool {
	return !strings.Contains(value.(string), sm.searchValue)
}
