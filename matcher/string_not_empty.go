package matcher

type stringNotEmpty struct {
}

func StringNotEmpty() Generic {
	obj := new(stringNotEmpty)
	return obj
}

func (sm *stringNotEmpty) Match(value interface{}) bool {
	return value.(string) != ""
}

