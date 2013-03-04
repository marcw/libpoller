package check

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewCheck(t *testing.T) {
	_, err := NewCheck("https://google.com", "foobar", "1s", false, "", false, make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("http://google.com", "foobar", "1s", false, "", false, make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("https://google.com:8444", "foobar", "1s", false, "", false, make(map[string]string))
	if err != nil {
		t.Error("NewCheck should not returns an error here")
	}

	_, err = NewCheck("https://google.com:8444", "foobar", "1s", true, "fadsfs", false, make(map[string]string))
	if err == nil {
		t.Error("NewCheck should returns an error here")
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
        "alert": true,
        "alertDelay": "1h",
        "interval": "60s",
        "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"
}}]`

	s := NewScheduler()
	err := s.AddFromJSON([]byte(json))
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
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

	if check.Alert != true {
		t.Errorf("Check Alert should be true.")
	}

	data, err := s.JSON()
	if err != nil {
		t.Log(err)
		t.Error("Marshaling to JSON should not fail")
	}
	s.Wipe()
	if s.Len() > 0 {
		t.Error("After Wipe, length should be 0")
	}

	err = s.AddFromJSON(data)
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
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

func TestServeHTTP(t *testing.T) {
	s := NewScheduler()

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
	s.Wipe()

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

func TestShouldAlert(t *testing.T) {
	c, _ := NewCheck("foo", "foo", "10s", false, "", false, make(map[string]string))

	c.Alert = true
	c.Alerted = false
	c.DownSince = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	if c.ShouldAlert() == false {
		t.Errorf("Should be true")
	}

	c.Alert = true
	c.Alerted = true
	c.DownSince = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	if c.ShouldAlert() == true {
		t.Errorf("Should be false")
	}

	c.Alert = false
	c.Alerted = false
	c.DownSince = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	if c.ShouldAlert() == true {
		t.Errorf("Should be false")
	}

	c.Alert = true
	c.Alerted = false
	c.DownSince = time.Now()
	c.AlertDelay = time.Hour
	if c.ShouldAlert() == true {
		t.Errorf("Should be false")
	}

	c.Alert = true
	c.Alerted = false
	c.DownSince = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	c.AlertDelay = time.Hour
	if c.ShouldAlert() == false {
		t.Errorf("Should be true")
	}
}

func TestShouldNotifyFix(t *testing.T) {
	c, _ := NewCheck("foo", "foo", "10s", false, "", true, make(map[string]string))

	c.NotifyFix = false
	c.WasDownFor, _ = time.ParseDuration("10s")
	if c.ShouldNotifyFix() == true {
		t.Errorf("Should be false")
	}

	c.NotifyFix = true
	c.Alerted = false
	if c.ShouldNotifyFix() == true {
		t.Errorf("Should be false")
	}

	c.NotifyFix = true
	c.Alerted = true
	c.WasDownFor, _ = time.ParseDuration("10s")
	if c.ShouldNotifyFix() == false {
		t.Errorf("Should be true")
	}
}
