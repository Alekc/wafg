package matcher

import "testing"

func TestStringContains_Match(t *testing.T) {
	obj := &stringContains{searchValue: "ab"}
	if !obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain ab")
	}
	obj.searchValue = "fg"
	if !obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain fg")
	}
	obj.searchValue = "cd"
	if !obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain cd")
	}
	if !obj.Match("abcdefg") {
		t.Errorf("abcdefg doesn't contain abcdefg")
	}
}

func TestStringContains(t *testing.T) {
	searchKey := "alksjdlkj"
	obj := StringContains(searchKey)
	var m *stringContains
	var ok bool
	if m, ok = obj.(*stringContains); !ok {
		t.Errorf("object is not string contains")
		return
	}
	if m.searchValue != searchKey {
		t.Errorf("Search value is wrong")
	}
}
