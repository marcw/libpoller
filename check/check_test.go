package check

import (
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
