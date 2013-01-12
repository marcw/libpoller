package check

import (
	"testing"
)

func TestNewCheck(t *testing.T) {
	_, err := NewCheck("https://google.com", "foobar", "1s", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("http://google.com", "foobar", "1s", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("https://google.com:8444", "foobar", "1s", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}
}
