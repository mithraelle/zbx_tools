package config

import (
	"strings"
)

type Option interface {
	set(*Config)
}

type HostnameOption string

func (hostnameOption HostnameOption) set(c *Config) {
	c.Hostname = string(hostnameOption)
}

func WithHostname(hostname string) HostnameOption {
	return HostnameOption(hostname)
}

type ServerAddrOption string

func (opt ServerAddrOption) set(c *Config) {
	c.ServerAddr = opt.String()
}

func (opt ServerAddrOption) String() string {
	s := string(opt)

	if strings.Index(s, ":") == -1 {
		s += ":10051"
	}

	return s
}

func WithServerAddr(serverAddr string) ServerAddrOption {
	return ServerAddrOption(serverAddr)
}

type TimeoutOption int

func (timeoutOption TimeoutOption) set(c *Config) {
	c.Timeout = int(timeoutOption)
}

func WithTimeout(timeout int) TimeoutOption {
	return TimeoutOption(timeout)
}
