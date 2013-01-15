package poll

import (
	"github.com/marcw/poller/check"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type successPollHandler struct {
}
type errorPollHandler struct {
}
type timeoutPollHandler struct {
}

func (p successPollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (p errorPollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error", 500)
}

func (p timeoutPollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(200 * time.Millisecond)
}

func TestSuccessfullPoll(t *testing.T) {
	server := httptest.NewServer(successPollHandler{})
	defer server.Close()

	poll := NewHttpPoller("foobar", 10*time.Second)

	c, _ := check.NewCheck(server.URL, "foobar", "10s", make(map[string]string))
	event := poll.Poll(c)
	if event.StatusCode != 200 {
		t.Error("statusCode should be 200")
	}
	if event.Up != true {
		t.Error("Up should be true")
	}
	if event.Duration.Nanoseconds() == 0 {
		t.Error("Duration can't be equals to 0 nanosecond")
	}
}

func TestFailedPoll(t *testing.T) {
	server := httptest.NewServer(errorPollHandler{})
	defer server.Close()

	poll := NewHttpPoller("foobar", 10*time.Second)

	c, _ := check.NewCheck(server.URL, "foobar", "10s", make(map[string]string))
	event := poll.Poll(c)
	if event.StatusCode != 500 {
		t.Error("statusCode should be 500")
	}
	if event.Up != false {
		t.Error("Up should be false")
	}
}

func TestTimeoutedPoll(t *testing.T) {
	server := httptest.NewServer(timeoutPollHandler{})
	defer server.Close()

	poll := NewHttpPoller("foobar", 100*time.Millisecond)

	c, _ := check.NewCheck(server.URL, "foobar", "10s", make(map[string]string))
	event := poll.Poll(c)
	if event.StatusCode != 0 {
		t.Error("statusCode should be 0")
	}
	if event.Up != false {
		t.Error("Up should be false")
	}
}
