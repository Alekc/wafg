package matcher

type stringMatcher struct {
	searchValue string
}

func StringMatch(search string) Generic {
	obj := new(stringMatcher)
	obj.searchValue = search
	return obj
}

func (sm *stringMatcher) Match(value interface{}) bool {
	return value.(string) == sm.searchValue
}
