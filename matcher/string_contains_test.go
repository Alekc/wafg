package matcher

import "testing"

func TestStringContains_Match(t *testing.T) {
	obj := &stringContains{searchValue:"ab"}
	if !obj.Match("abcdefg") {
		t.Errorf("String is not matched at the beginning")
	}
	obj.searchValue = "fg"
	if !obj.Match("abcdefg") {
		t.Errorf("String is not matched at the end")
	}
	obj.searchValue = "cd"
	if !obj.Match("abcdefg") {
		t.Errorf("String is not matched in the middle")
	}
	if !obj.Match("abcdefg"){
		t.Errorf("Full string is not matched")
	}
}

func TestStringContains(t *testing.T){
	searchKey := "alksjdlkj"
	obj := StringContains(searchKey)
	var m *stringContains
	var ok bool
	if m,ok = obj.(*stringContains); !ok {
		t.Errorf("object is not string contains")
	}
	if m.searchValue != searchKey {
		t.Errorf("Search value is wrong")
	}
}