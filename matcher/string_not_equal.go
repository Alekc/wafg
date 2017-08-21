package matcher

type stringNotEqual struct {
	searchValue string
}

func StringNotEqual(search string) Generic {
	obj := new(stringNotEqual)
	obj.searchValue = search
	return obj
}

func (sm *stringNotEqual) Match(value interface{}) bool {
	return value.(string) != sm.searchValue
}
