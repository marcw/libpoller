package poller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
        "alert": true,
        "alertDelay": "1h",
        "interval": "60s",
        "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"
}}]`

	c := NewConfig(NewInMemoryStore(), NewScheduler())
	err := c.AddFromJSON([]byte(json))
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
	}

	check := c.Get("connect_sensiolabs_com_api")
	if check == nil {
		t.Log("Checkslist should contain key check")
		t.FailNow()
	}

	if check.Header.Get("Accept") != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Check headers does not contain Accept header")
	}

	if check.Interval.Seconds() != 60 {
		t.Errorf("Check interval should be equal to 60s.")
	}

	if check.Alert != true {
		t.Errorf("Check Alert should be true.")
	}

	data, err := c.JSON()
	if err != nil {
		t.Log(err)
		t.Error("Marshaling to JSON should not fail")
	}
	c.Clear()
	if c.Len() > 0 {
		t.Error("After Clear, length should be 0")
	}

	err = c.AddFromJSON(data)
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
	}

	check = c.Get("connect_sensiolabs_com_api")
	if check == nil {
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
	s := NewConfig(NewInMemoryStore(), NewScheduler())

	server := httptest.NewServer(s)
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
	check := s.Get("connect_sensiolabs_com_api")
	if check == nil {
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
	s.Clear()

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
	check = s.Get("connect_sensiolabs_com_api")
	if check == nil {
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
