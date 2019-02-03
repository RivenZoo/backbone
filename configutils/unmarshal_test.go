package configutils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type config struct {
		Name string            `yaml:"name"`
		ID   int               `yaml:"id"`
		KV   map[string]string `yaml:"kv"`
	}
	name := "config"
	id := 1234

	jsonData := []byte(`{"Name": "config", "ID": 1234, "KV": {"v1": "1"}}`)
	tomlData := []byte(`
Name = "config"
ID = 1234
[KV]
v1 = "1"
`)
	yamlData := []byte(`
name : config 
id : 1234
kv:
  v1: "1"
`)
	dataSets := [][]byte{jsonData, tomlData, yamlData}
	tp := []ConfigType{ConfigTypeJSON, ConfigTypeTOML, ConfigTypeYAML}
	for i := range dataSets {
		v := config{}
		r := bytes.NewReader(dataSets[i])
		err := Unmarshal(r, &v, tp[i])
		if assert.Nil(t, err) {
			t.Log(v)
		}
		assert.Equal(t, name, v.Name)
		assert.Equal(t, id, v.ID)
	}

	{
		r := bytes.NewReader(jsonData)
		var v *config
		err := Unmarshal(r, &v, ConfigTypeJSON)
		assert.Nil(t, err)
		t.Log(*v)
	}

	{
		r := bytes.NewReader(jsonData)
		v := config{}
		err := Unmarshal(r, v, ConfigTypeJSON)
		assert.NotNil(t, err)
	}

	{
		r := bytes.NewReader(jsonData[:1])
		v := config{}
		err := Unmarshal(r, &v, ConfigTypeJSON)
		assert.NotNil(t, err)
	}
}

func TestUnmarshal2(t *testing.T) {
	type config struct {
		Name string            `yaml:"name"`
		ID   int               `yaml:"id"`
		KV   map[string]string `yaml:"kv"`
	}

	// wrong json format
	jsonData := []byte(`{"Name": "config", "ID": 1234, "KV": {"v1": "1"}, }`)

	v := config{}
	r := bytes.NewReader(jsonData)
	err := Unmarshal(r, &v, ConfigTypeJSON)
	assert.NotNil(t, err)
	t.Log(v, err)
}
