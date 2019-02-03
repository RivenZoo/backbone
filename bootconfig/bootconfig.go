package bootconfig

import "github.com/RivenZoo/backbone/configutils"

type ConfigGetter interface {
	GetConfig(key string) (RawConfigData, configutils.ConfigType)
}

var configGetter ConfigGetter

func RegisterConfigGetter(g ConfigGetter) {
	configGetter = g
}

func GetConfigGetter() ConfigGetter {
	return configGetter
}
