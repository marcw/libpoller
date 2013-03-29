package poller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeHTTPPost(t *testing.T) {
	c := NewConfig(NewInMemoryStore(), NewSimpleScheduler())

	server := httptest.NewServer(NewConfigHttpHandler(c))
	defer server.Close()

	json := `
    {
    "key": "connect_sensiolabs_com_api",
    "url": "https://connect.sensiolabs.com/api/",
    "timeout": "10s",
    "interval": "60s",
    "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"}}`

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
	check := c.checks.Get("connect_sensiolabs_com_api")
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
	resp.Body.Close()
	json = string(data)
	c.Clear()

	r, err = http.NewRequest("PUT", server.URL, strings.NewReader(json))
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
	check = c.checks.Get("connect_sensiolabs_com_api")
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
