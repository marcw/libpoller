package poller

import (
	"net"
	"testing"
	"time"
)

func TestUDPSuccessfullTest(t *testing.T) {
	probe := NewUdpProbe(10 * time.Second)

	c, _ := NewCheck("foobar", "10s", false, "", false, make(map[string]interface{}))
	c.Config.Set("host", "localhost")
	c.Config.Set("port", 4321)
	c.Config.Set("send", "foobar")
	c.Config.Set("receive", "foobar")

	conn, err := net.ListenPacket("udp", "localhost:4321")
	if err != nil {
		t.Error(err.Error())
	}
	go func() {
		sent := make([]byte, 16)
		_, addr, err := conn.ReadFrom(sent)
		if err != nil {
			t.Error(err.Error())
		}
		if _, err := conn.WriteTo(sent, addr); err != nil {
			t.Error(err.Error())
		}
	}()
	event := probe.Test(c)
	if event.IsUp() != true {
		t.Error("IsUp() should be true")
	}
	if event.Duration.Nanoseconds() == 0 {
		t.Error("Duration can't be equals to 0 nanosecond")
	}
}
