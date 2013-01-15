package check

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestServeHTTP(t *testing.T) {
	c := ChecksList{}

	server := httptest.NewServer(c)
	defer server.Close()

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

	r, err := http.NewRequest("POST", server.URL, strings.NewReader(json))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("Status code should be 201. Got %d", resp.StatusCode)
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

	resp, err = http.DefaultClient.Get(server.URL + "/checks")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	json = string(data)
	resp.Body.Close()
	c.Wipe()

	r, err = http.NewRequest("POST", server.URL, strings.NewReader(json))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("Status code should be 201. Got %d", resp.StatusCode)
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
