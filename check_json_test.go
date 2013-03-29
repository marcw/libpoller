package poller

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNewCheckFromJSON(t *testing.T) {
	json := `
    {
        "key": "connect_sensiolabs_com_api",
        "url": "https://connect.sensiolabs.com/api/",
        "alert": true,
        "alertDelay": "1h",
        "interval": "60s",
        "notifyFix": true,
        "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"
}}`
	check, err := NewCheckFromJSON([]byte(json))

	if err != nil {
		t.Error(err)
	}

	if check.Header.Get("Accept") != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Headers do not contain Accept header")
	}

	if check.Interval.Seconds() != 60 {
		t.Errorf("Interval should be equal to 60s.")
	}

	if check.Alert != true {
		t.Errorf("Alert should be true.")
	}

	if check.AlertDelay.Seconds() != 3600 {
		t.Errorf("Alert delay is wrong.")
	}

	if check.Url.String() != "https://connect.sensiolabs.com/api/" {
		t.Errorf("delay is wrong.")
	}

	if check.NotifyFix != true {
		t.Errorf("NotifyFix is wrong")
	}
}

func TestCheckJSON(t *testing.T) {
	jsn := `{
        "url": "https://connect.sensiolabs.com/api/",
        "key": "connect_sensiolabs_com_api",
        "interval": "1m0s",
        "alert": true,
        "alertDelay": "1h0m0s",
        "notifyFix": true,
        "headers": {
        "Accept": "application/vnd.com.sensiolabs.connect+xml"
}}`

	buffer := new(bytes.Buffer)
	check, _ := NewCheckFromJSON([]byte(jsn))
	marshaled, _ := check.JSON()

	json.Compact(buffer, []byte(jsn))
	if string(marshaled) != buffer.String() {
		t.Errorf("JSON() do not output correct representation of Check")
	}
}
