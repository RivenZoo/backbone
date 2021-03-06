package bootconfig

import (
	"bytes"
	"github.com/RivenZoo/backbone/configutils"
)

type RawConfigData []byte

// Unmarshal support json/toml/yaml.
// Notice: yaml should use tag `yaml:"key"`
func (d RawConfigData) Unmarshal(v interface{}, tp configutils.ConfigType) error {
	r := bytes.NewReader(d)
	return configutils.Unmarshal(r, v, tp)
}
