package configutils

import (
	"path"
	"strings"
)

const (
	ConfigTypeJSON ConfigType = "json"
	ConfigTypeTOML ConfigType = "toml"
	ConfigTypeYAML ConfigType = "yaml"
)

type ConfigType string

func DetechFileConfigType(fname string) ConfigType {
	ext := strings.ToLower(path.Ext(fname))
	switch ext {
	case ".json":
		return ConfigTypeJSON
	case ".toml":
		return ConfigTypeTOML
	case ".yaml", ".yml":
		return ConfigTypeYAML
	}
	return ConfigType("")
}
