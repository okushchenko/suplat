package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	Version   string
	BuildTime string
)

type network struct {
	CIDRs []string // CIDRs of a networks
}

type config struct {
	Networks    map[string]network
	InfluxDBUrl string // metrics storage url
}

var C config
var err error

func InitConfig() error {
	_, err = toml.DecodeFile("./config.toml", &C)
	if err != nil {
		return fmt.Errorf("Failed to decode config: %s", err.Error())
	}
	return nil
}
