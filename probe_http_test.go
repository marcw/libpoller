package poller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type successTestHandler struct {
}
type errorTestHandler struct {
}
type timeoutTestHandler struct {
}

func (p successTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (p errorTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error", 500)
}

func (p timeoutTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(200 * time.Millisecond)
}

func TestSuccessfullTest(t *testing.T) {
	server := httptest.NewServer(successTestHandler{})
	defer server.Close()

	probe := NewHttpProbe("foobar", 10*time.Second)

	c, _ := NewCheck(server.URL, "foobar", "10s", false, "", false, make(map[string]string))
	event := probe.Test(c)
	if event.StatusCode != 200 {
		t.Error("statusCode should be 200")
	}
	if event.IsUp() != true {
		t.Error("IsUp() should be true")
	}
	if event.Duration.Nanoseconds() == 0 {
		t.Error("Duration can't be equals to 0 nanosecond")
	}
}

func TestFailedTest(t *testing.T) {
	server := httptest.NewServer(errorTestHandler{})
	defer server.Close()

	probe := NewHttpProbe("foobar", 10*time.Second)

	c, _ := NewCheck(server.URL, "foobar", "10s", false, "", false, make(map[string]string))
	event := probe.Test(c)
	if event.StatusCode != 500 {
		t.Error("statusCode should be 500")
	}
	if event.IsUp() != false {
		t.Error("IsUp() should be false")
	}
}

func TestTimeoutedTest(t *testing.T) {
	server := httptest.NewServer(timeoutTestHandler{})
	defer server.Close()

	probe := NewHttpProbe("foobar", 100*time.Millisecond)

	c, _ := NewCheck(server.URL, "foobar", "10s", false, "", false, make(map[string]string))
	event := probe.Test(c)
	if event.StatusCode != 0 {
		t.Error("statusCode should be 0")
	}
	if event.IsUp() != false {
		t.Error("IsUp() should be false")
	}
}
