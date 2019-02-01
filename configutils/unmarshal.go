package configutils

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"reflect"
)

var (
	errWrongUnmarshalReceiver = errors.New("unmarshal output should be pointer")
	errUnsupportFormat        = errors.New("unsupport input format")
)

// Unmarshal support json/toml/yaml.
// Decode order: json,toml,yaml.
// Notice: yaml should use tag `yaml:"key"`
func Unmarshal(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	unmarshalFn := [](func([]byte, interface{}) error){unmarshalJson, unmarshalToml, unmarshalYaml}
	for i := range unmarshalFn {
		tp := reflect.TypeOf(v)
		var newVal interface{}
		switch tp.Kind() {
		case reflect.Ptr:
			newVal = reflect.New(tp.Elem()).Interface()
		default:
			return errWrongUnmarshalReceiver
		}
		err = unmarshalFn[i](data, newVal)
		if err == nil {
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(newVal).Elem())
			return nil
		}
	}
	return errUnsupportFormat
}

// UnmarshalFile support json/toml/yaml.
// Decode order: json,toml,yaml.
// Notice: yaml should use tag `yaml:"key"`
func UnmarshalFile(fPath string, v interface{}) error {
	f, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v)
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
