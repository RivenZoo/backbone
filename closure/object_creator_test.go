package closure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	name string
}

func TestObjectCreator_Validate(t *testing.T) {
	var p *person
	oc := ObjectCreator{
		CreateFunc: func() (person, error) { return person{"test"}, nil },
		Receiver:   p,
	}
	// validate success
	oc.Validate()

	oc = ObjectCreator{
		CreateFunc: func() person { return person{"test"} },
		Receiver:   p,
	}
	// validate success
	oc.Validate()

	oc = ObjectCreator{
		CreateFunc: func() {},
	}
	assert.Panics(t, func() {
		oc.Validate()
	}, "should panic because CreateFunc must be func() T or func() (T,error)")

	oc = ObjectCreator{
		CreateFunc: func() (int, person, error) {
			return 0, person{""}, nil
		},
	}
	assert.Panics(t, func() {
		oc.Validate()
	}, "should panic because CreateFunc must be func() T or func() (T,error)")
}

func TestObjectCreator_Create(t *testing.T) {
	target := person{"test"}
	var p *person
	oc := ObjectCreator{
		CreateFunc: func() (person, error) { return target, nil },
		Receiver:   p,
	}
	oc.Validate()
	ret, err := oc.Create()
	assert.Nil(t, err)

	retObj, ok := ret.(person)
	assert.True(t, ok)
	assert.EqualValues(t, target, retObj)
}

func TestObjectCreator_SetReceiver(t *testing.T) {
	target := person{"test"}
	p := &person{}
	oc := ObjectCreator{
		CreateFunc: func() (person, error) { return target, nil },
		Receiver:   p,
	}
	oc.Validate()
	ret, err := oc.Create()
	assert.Nil(t, err)

	oc.SetReceiver(ret)
	assert.EqualValues(t, target, ret)
	assert.EqualValues(t, target, *p)
}
