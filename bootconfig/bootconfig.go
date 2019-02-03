package bootconfig

type ConfigGetter interface {
	GetConfig(key string) RawConfigData
}

var configGetter ConfigGetter

func RegisterConfigGetter(g ConfigGetter) {
	configGetter = g
}

func GetConfigGetter() ConfigGetter {
	return configGetter
}
