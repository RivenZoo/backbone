package bootconfig

import (
	"errors"
	"github.com/RivenZoo/backbone/configutils"
	"io/ioutil"
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
			tp = configutils.DetechFileConfigType(fname)
		}
		return RawConfigData(data), tp
	}
	panic(errNoSuchBootConfig)
}
