package poller

import (
	"encoding/json"
)

// Used for marshalling / unmarshalling
type jsonCheck struct {
	Url        string            `json:"url"`
	Key        string            `json:"key"`
	Interval   string            `json:"interval"`
	Alert      bool              `json:"alert"`
	AlertDelay string            `json:"alertDelay"`
	NotifyFix  bool              `json:"notifyFix"`
	Headers    map[string]string `json:"headers"`
}

func (c *jsonCheck) toCheck() (*Check, error) {
	return NewCheck(c.Url, c.Key, c.Interval, c.Alert, c.AlertDelay, c.NotifyFix, c.Headers)
}

// Returns a jsonCheck object used internally before marshalling check to JSON
func (c *Check) json() *jsonCheck {
	header := make(map[string]string)
	for k, h := range c.Header {
		header[k] = h[0]
	}

	check := &jsonCheck{
		Url:        c.Url.String(),
		Key:        c.Key,
		Interval:   c.Interval.String(),
		NotifyFix:  c.NotifyFix,
		Alert:      c.Alert,
		AlertDelay: c.AlertDelay.String(),
		Headers:    header}

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
	check := &jsonCheck{Headers: make(map[string]string)}
	if err := json.Unmarshal(data, check); err != nil {
		return nil, err
	}

	return check.toCheck()
}
