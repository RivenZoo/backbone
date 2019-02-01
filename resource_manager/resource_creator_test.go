package resource_manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type closableCreator struct {
	Name string
}

func (c closableCreator) Close() error {
	return nil
}

func newCreator() (closableCreator, error) {
	return closableCreator{"new creator"}, nil
}

func TestResourceCreator(t *testing.T) {
	var cc closableCreator
	c := ResourceCreator{
		CreateFunc: newCreator,
		Receiver:   &cc,
	}
	c.validate()
	obj, err := c.create()
	assert.Nil(t, err)
	c.setReceiver(obj)
	retCC := obj.(closableCreator)
	assert.EqualValues(t, retCC.Name, cc.Name)
	t.Log(cc, retCC)

	c = ResourceCreator{
		CreateFunc: func() int { return 0 },
	}
	assert.Panics(t, func() {
		defer func() {
			if e := recover(); e != nil {
				t.Log(e)
				panic(e)
			}
		}()
		c.validate()
		c.create()
	}, "should panic because return T must implement Closable")
}
