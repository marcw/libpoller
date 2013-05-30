package poller

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"time"
)

type checkConfigurator func(*Check, *simplejson.Json) error

var checkConfigurators = map[CheckType]checkConfigurator{
	CheckTypeUDP:  readUDPConfig,
	CheckTypeHTTP: readHTTPConfig}

// Used for marshalling / unmarshalling
type jsonCheck struct {
	Type       string                 `json:"type"`
	Key        string                 `json:"key"`
	Interval   string                 `json:"interval"`
	Alert      bool                   `json:"alert"`
	AlertDelay string                 `json:"alertDelay"`
	NotifyFix  bool                   `json:"notifyFix"`
	Config     map[string]interface{} `json:"config"`
}

func (c *jsonCheck) toCheck() (*Check, error) {
	return NewCheck(c.Key, c.Interval, c.Alert, c.AlertDelay, c.NotifyFix, c.Config)
}

// Returns a jsonCheck object used internally before marshalling check to JSON
func (c *Check) json() *jsonCheck {
	check := &jsonCheck{
		Type:       string(c.Type()),
		Key:        c.Key,
		Interval:   c.Interval.String(),
		NotifyFix:  c.NotifyFix,
		Alert:      c.Alert,
		AlertDelay: c.AlertDelay.String(),
		Config:     c.Config.Map()}

	return check
}

// Returns a JSON representation of the Check.
func (c *Check) JSON() ([]byte, error) {
	data, err := json.Marshal(c.json())
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewCheckFromJSON() instantiates a new Check from a JSON representation.
func NewCheckFromJSON(data []byte) (*Check, error) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}

	check := newCheck()
	var checkType, key, alertDelay, interval string
	var alert, notifyFix bool

	if checkType, err = js.Get("type").String(); err != nil {
		return nil, err
	}

	if key, err = js.Get("key").String(); err != nil {
		return nil, err
	}

	if alert, err = js.Get("alert").Bool(); err != nil {
		return nil, err
	}
	if notifyFix, err = js.Get("notifyFix").Bool(); err != nil {
		return nil, err
	}

	if alertDelay, err = js.Get("alertDelay").String(); err != nil {
		return nil, err
	}

	if interval, err = js.Get("interval").String(); err != nil {
		return nil, err
	}

	check.checkType = CheckType(checkType)
	check.Key = key
	check.Alert = alert
	check.NotifyFix = notifyFix
	if delay, err := time.ParseDuration(alertDelay); err != nil {
		return nil, err
	} else {
		check.AlertDelay = delay
	}
	if interval, err := time.ParseDuration(interval); err != nil {
		return nil, err
	} else {
		check.Interval = interval
	}

	configurator, ok := checkConfigurators[check.Type()]
	if !ok {
		// TODO: Be nice to the user and try to guess what he meant
		return nil, fmt.Errorf("Unknown check type %s", checkType)
	}
	if err := configurator(check, js); err != nil {
		return nil, err
	}

	return check, nil
}

func readHTTPConfig(check *Check, js *simplejson.Json) error {
	if url, err := js.Get("config").Get("url").String(); err != nil {
		return err
	} else {
		check.Config.Set("url", url)
	}

	headers := make(map[string]string)
	for k, v := range js.Get("config").Get("headers").MustMap() {
		value, ok := v.(string)
		if !ok {
			return fmt.Errorf("Headers can only accept string.")
		}
		headers[k] = value
	}
	check.Config.Set("headers", headers)

	return nil
}

func readUDPConfig(check *Check, js *simplejson.Json) error {
	if host, err := js.Get("config").Get("host").String(); err != nil {
		return err
	} else {
		check.Config.Set("host", host)
	}
	if port, err := js.Get("config").Get("port").Int(); err != nil {
		return err
	} else {
		check.Config.Set("port", port)
	}
	if send, err := js.Get("config").Get("send").String(); err != nil {
		return err
	} else {
		check.Config.Set("send", send)
	}
	if receive, err := js.Get("config").Get("receive").String(); err != nil {
		return err
	} else {
		check.Config.Set("receive", receive)
	}

	return nil
}
