package checker

import (
	"testing"
)

func TestNewChecker(t *testing.T) {
	var c Checker
	c.host = "127.0.0.1"
	c.influxDb = "redfactor"
	c.port = 8086
}

func BenchmarkNewCheker(b *testing.B) {
	var c Checker
	c.host = "127.0.0.1"
	c.influxDb = "redfactor"
	c.port = 8086
}