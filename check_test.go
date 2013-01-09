package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
)

type successPollHandler struct {
}

func (p successPollHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestPollIsContactingHttpServer(t *testing.T) {
	server := httptest.NewServer(successPollHandler{})
    defer server.Close()

	check, _ := NewCheck(server.URL, "foobar", "10s", make(map[string]string))
	statusCode, _, err := check.Poll(&http.Client{})
    if err != nil {
        t.Error(err)
	}
    if statusCode != 200 {
        t.Error("statusCode should be 200")
	}
}
