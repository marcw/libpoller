package config

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadBasicConfig(t *testing.T) {
	json := `{
    "backends": ["stdout"],
    "checks": [
        {
            "key": "connect_sensiolabs_com_api",
            "url": "https://connect.sensiolabs.com/api/",
            "interval": "60s",
            "headers": {
                "Accept": "application/vnd.com.sensiolabs.connect+xml"
            }
        }
    ]
}`

	config := &Config{}
	err := config.Load([]byte(json))
	if err != nil {
		t.Error("Config failed to load with a valid json file")
	}

	if config.UserAgent != DEFAULT_USER_AGENT {
		t.Error("Config should have default User-Agent if none is provided")
	}

	timeout, _ := time.ParseDuration(DEFAULT_TIMEOUT)
	if config.Timeout.Seconds() != timeout.Seconds() {
		t.Error("Config should have default timeout")
	}

	check := config.Checks[0]
	if check.Header.Get("Accept") != "application/vnd.com.sensiolabs.connect+xml" {
		t.Error("Check headers does not contain Accept header")
	}

	if check.Header.Get("User-Agent") != DEFAULT_USER_AGENT {
		t.Error("Check headers does not contain User-Agent header")
	}

	if check.Interval.Seconds() != 60 {
		t.Error("Check interval should be equal to 60s")
	}
}

func TestLoadCustomizedConfig(t *testing.T) {
	json := `{
    "userAgent": "foobar",
    "backends": ["stdout"],
    "timeout": "5s",
    "checks": [
        {
            "key": "symfony_com",
            "url": "http://symfony.com",
            "interval": "60s"
        },
        {
            "key": "connect_sensiolabs_com_api",
            "url": "https://connect.sensiolabs.com/api/",
            "interval": "30s",
            "headers": {
                "Accept": "application/vnd.com.sensiolabs.connect+xml"
            }
        }
    ]
}`

	config := &Config{}
	err := config.Load([]byte(json))
	if err != nil {
		t.Log(err)
		t.Error("Config failed to load with a valid json file")
	}

	if config.Timeout.Seconds() != 5 {
		t.Error("Config should have customized global timeout")
	}

	if config.UserAgent != "foobar" {
		fmt.Println(config.UserAgent)
		t.Error("Config should have customized user-agent")
	}
}
