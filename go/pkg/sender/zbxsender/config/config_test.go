package config

import (
	"testing"
)

func compare(t *testing.T, expected Config, actual Config) {
	if expected.Timeout != actual.Timeout {
		t.Errorf("Expected Timeout %v\nGot: %v", expected.Timeout, actual.Timeout)
	}

	if expected.ServerAddr != actual.ServerAddr {
		t.Errorf("Expected ServerAddr %v\nGot: %v", expected.ServerAddr, actual.ServerAddr)
	}

	if expected.Hostname != actual.Hostname {
		t.Errorf("Expected Hostname %v\nGot: %v", expected.Hostname, actual.Hostname)
	}
}

func TestConfig_WithOptions(t *testing.T) {
	expected := Config{Timeout: 3, Hostname: "ZabbixClient", ServerAddr: "somehost:12345"}
	compare(t, WithOptions(WithServerAddr("somehost:12345")), expected)

	expected.Timeout = 123
	expected.ServerAddr = "somehost:10051"
	expected.Hostname = "MyHost"
	compare(t, WithOptions(WithTimeout(123), WithServerAddr("somehost"), WithHostname("MyHost")), expected)
}
