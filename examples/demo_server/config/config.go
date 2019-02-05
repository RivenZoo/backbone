package config

import (
	"github.com/RivenZoo/backbone/bootconfig"
	"github.com/RivenZoo/backbone/configutils"
)

var config = &Config{}

type bootConfigFileSet map[string]string

type Config struct {
	BootConfigs bootConfigFileSet `json:"boot_configs"`
}

func GetConfig() *Config {
	return config
}

func MustLoadConfig(fname string) {
	if err := configutils.UnmarshalFile(fname, config, configutils.ConfigTypeJSON); err != nil {
		panic(err)
	}
	bootCfgGetter := bootconfig.NewFileBootConfigGetter(bootconfig.FileBootConfig{
		BootConfigFileSet: config.BootConfigs,
	})
	bootconfig.RegisterConfigGetter(bootCfgGetter)
}
