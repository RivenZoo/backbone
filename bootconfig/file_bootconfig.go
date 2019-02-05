package bootconfig

import (
	"errors"
	"github.com/RivenZoo/backbone/configutils"
	"io/ioutil"
	"path"
	"strings"
)

var errNoSuchBootConfig = errors.New("no such boot config provided")

type FileBootConfig struct {
	BootConfigFileSet map[string]string
	Type              configutils.ConfigType
}

type FileBootConfigGetter struct {
	config FileBootConfig
}

func NewFileBootConfigGetter(cfg FileBootConfig) FileBootConfigGetter {
	ret := FileBootConfigGetter{
		config: cfg,
	}
	return ret
}

func (g FileBootConfigGetter) GetConfig(key string) (RawConfigData, configutils.ConfigType) {
	if fname, ok := g.config.BootConfigFileSet[key]; ok {
		tp := g.config.Type
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			panic(err)
		}
		if tp == "" {
			tp = detechFileConfigType(fname)
		}
		return RawConfigData(data), tp
	}
	panic(errNoSuchBootConfig)
}

func detechFileConfigType(fname string) configutils.ConfigType {
	ext := strings.ToLower(path.Ext(fname))
	switch ext {
	case ".json":
		return configutils.ConfigTypeJSON
	case ".toml":
		return configutils.ConfigTypeTOML
	case ".yaml", ".yml":
		return configutils.ConfigTypeYAML
	}
	return configutils.ConfigType("")
}
