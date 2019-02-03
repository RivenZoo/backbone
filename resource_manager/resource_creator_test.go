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
	c := NewResourceCreator(newCreator, &cc)
	c.Validate()

	obj, err := c.createResource()
	assert.Nil(t, err)
	c.SetReceiver(obj)
	retCC := obj.(closableCreator)
	assert.EqualValues(t, retCC.Name, cc.Name)
	t.Log(cc, retCC)

	c = NewResourceCreator(func() int { return 0 }, nil)
	c.Validate()
	_, err = c.createResource()
	assert.EqualValues(t, errUnSupportCreatorFunc, err)
}
