package poller

import (
	"testing"
)

func TestCheckListOperations(t *testing.T) {
	cl := NewCheckList()

	if 0 != cl.Len() {
		t.Error("CheckList's length should be 0")
	}

	check, _ := NewCheck("http://google.com", "foobar", "10s", false, "0s", false, make(map[string]string))
	cl.Add(check)

	if 1 != cl.Len() {
		t.Error("CheckList's length should be 0")
	}

	if nil == cl.Get("foobar") {
		t.Error("Get should have returned the checks instance")
	}

	cl.Each(func(check *Check) {
		if check.Key != "foobar" {
			t.Error("Each should receive have received the check")
		}
	})

	cl.Clear()
	if 0 != cl.Len() {
		t.Error("CheckList's length should be 0")
	}
}

func TestNewCheckListFromJSON(t *testing.T) {
	cl := NewCheckList()
	check, _ := NewCheck("http://google.com", "foobar", "10s", false, "0s", false, make(map[string]string))
	cl.Add(check)

	json, _ := cl.JSON()
	cl, _ = NewCheckListFromJSON(json)

	if 1 != cl.Len() {
		t.Error("CheckList's length should be 0")
	}

	if nil == cl.Get("foobar") {
		t.Error("Get should have returned the checks instance")
	}
}
