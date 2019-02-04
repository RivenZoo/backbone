package configutils

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

const (
	ConfigTypeJSON ConfigType = "json"
	ConfigTypeTOML ConfigType = "toml"
	ConfigTypeYAML ConfigType = "yaml"
)

type ConfigType string

var (
	errUnsupportFormat = errors.New("unsupport input format")
)

// Unmarshal support json/toml/yaml.
// Notice: yaml should use tag `yaml:"key"`
func Unmarshal(r io.Reader, v interface{}, tp ConfigType) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	switch tp {
	case ConfigTypeJSON:
		return unmarshalJson(data, v)
	case ConfigTypeTOML:
		return unmarshalToml(data, v)
	case ConfigTypeYAML:
		return unmarshalYaml(data, v)
	default:
	}
	return errUnsupportFormat
}

// UnmarshalFile support json/toml/yaml.
// Notice: yaml should use tag `yaml:"key"`
func UnmarshalFile(fPath string, v interface{}, tp ConfigType) error {
	f, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v, tp)
}

func unmarshalJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func unmarshalToml(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func unmarshalYaml(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
