package poller

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"testing"
)

var testJsonHttpCheck = `
{
    "type": "http",
    "key": "connect_sensiolabs_com_api",
    "interval": "1m0s",
    "alert": true,
    "alertDelay": "1h0m0s",
    "notifyFix": true,
    "config": {
        "headers": {
            "Accept": "application/vnd.com.sensiolabs.connect+xml"
        },
        "url": "https://connect.sensiolabs.com/api/"
    }
}`

func TestNewCheckFromJSON(t *testing.T) {
	check, err := NewCheckFromJSON([]byte(testJsonHttpCheck))
	if err != nil {
		t.Error(err)
	}

	headers := check.Config.GetMapStringString("headers")
	if headers["Accept"] != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Accept header is incorrect.")
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

	if check.Config.GetString("url") != "https://connect.sensiolabs.com/api/" {
		t.Errorf("delay is wrong.")
	}

	if check.NotifyFix != true {
		t.Errorf("NotifyFix is wrong")
	}
}

func TestCheckJSON(t *testing.T) {
	buffer := new(bytes.Buffer)
	check, _ := NewCheckFromJSON([]byte(testJsonHttpCheck))
	marshaled, _ := check.JSON()

	json.Compact(buffer, []byte(testJsonHttpCheck))
	if string(marshaled) != buffer.String() {
		fmt.Println(buffer.String())
		fmt.Println(string(marshaled))
		t.Errorf("JSON() do not output correct representation of Check")
	}
}
