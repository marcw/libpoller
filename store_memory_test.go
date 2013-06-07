package poller

import (
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	cl := NewInMemoryStore()

	l, _ := cl.Len()
	if l != 0 {
		t.Error("CheckList's length should be 0")
	}

	check, _ := NewCheck("foobar", "10s", false, "0s", false, make(map[string]interface{}))
	cl.Add(check)

	l, _ = cl.Len()
	if l != 1 {
		t.Error("CheckList's length should be 1")
	}

	check, _ = cl.Get("foobar")
	if check == nil {
		t.Error("Get should have returned the checks instance")
	}
}
