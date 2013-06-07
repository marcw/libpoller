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

	r, err := http.NewRequest("POST", server.URL, strings.NewReader(testJsonHttpCheck))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Status code should be 201. Got %d\n", resp.StatusCode)
		t.Errorf(string(body))
	}
	check, err := c.store.Get("connect_sensiolabs_com_api")
	if check == nil {
		t.Log("Store should contain check")
		t.FailNow()
	}

	if check.Interval.Seconds() != 60 {
		t.Errorf("Check interval should be equal to 60s.")
	}
}
