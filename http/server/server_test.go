package server

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSimpleServer_Run(t *testing.T) {
	cfg := &ServerConfig{
		Addr: "127.0.0.1:8000",
	}
	s, e := NewSimpleServer(cfg)
	assert.Nil(t, e)
	go func() {
		err := s.Run()
		assert.Nil(t, err)
	}()
	resp, err := http.Get("http://" + cfg.Addr)
	assert.Nil(t, err)

	data, err := ioutil.ReadAll(resp.Body)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	defer resp.Body.Close()

	t.Log(string(data))
	err = s.Stop()
	assert.Nil(t, err)
}
