package zbxsender

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	Timeout        int    `mapstructure:"Timeout"`
	Hostname       string `mapstructure:"Hostname"`
	ServerActive   string `mapstructure:"ServerActive"`
	TLSConnect     string `mapstructure:"TLSConnect"`
	TLSPSKIdentity string `mapstructure:"TLSPSKIdentity"`
	TLSPSKFile     string `mapstructure:"TLSPSKFile"`
	ServerAddr     string
}

func NewConfig(filename string) *Config {
	conf := Config{}

	viper.AddConfigPath(filepath.Dir(filename))
	viper.SetConfigName(filepath.Base(filename))
	viper.SetConfigType("env")

	viper.SetDefault("Timeout", 3)
	viper.SetDefault("Hostname", "ZabbixServer")
	viper.SetDefault("ServerActive", "localhost")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("Fatal error unmarshalling file: %s \n", err))
	}

	conf.setServerAddr()
	return &conf
}

func (c *Config) setServerAddr() {
	if strings.Index(c.ServerActive, ":") > -1 {
		c.ServerAddr = c.ServerActive
	} else {
		c.ServerAddr = c.ServerActive + ":10051"
	}
}
