package check

import (
	_ "fmt"
	"testing"
)

func TestNewCheck(t *testing.T) {
	_, err := NewCheck("https://google.com", "foobar", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("http://google.com", "foobar", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("https://google.com:8444", "foobar", "1s", make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}
}

func TestLoadBasicConfig(t *testing.T) {
	json := `[
    {
        "key": "symfony_com",
        "url": "http://symfony.com",
        "interval": "60s"
    },
    {
    "key": "connect_sensiolabs_com_api",
    "url": "https://connect.sensiolabs.com/api/",
    "timeout": "10s",
    "interval": "60s",
    "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"
}}]`

	c := ChecksList{}
	err := c.AddFromJson([]byte(json))
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
	}

	check, ok := c["connect_sensiolabs_com_api"]
	if !ok {
		t.Log("Checkslist should contain key check")
		t.FailNow()
	}

	if check.Header.Get("Accept") != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Check headers does not contain Accept header")
	}

	if check.Interval.Seconds() != 60 {
		t.Errorf("Check interval should be equal to 60s.")
	}

	data, err := c.JSON()
	if err != nil {
		t.Log(err)
		t.Error("Marshaling to JSON should not fail")
	}
	c.Wipe()
	if len(c) > 0 {
		t.Error("After Wipe, length should be 0")
	}

	err = c.AddFromJson(data)
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
	}

	check, ok = c["connect_sensiolabs_com_api"]
	if !ok {
		t.Log("Checkslist should contain key check")
		t.FailNow()
	}

	if check.Header.Get("Accept") != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Check headers does not contain Accept header")
	}

	if check.Interval.Seconds() != 60 {
		t.Errorf("Check interval should be equal to 60s.")
	}
}
