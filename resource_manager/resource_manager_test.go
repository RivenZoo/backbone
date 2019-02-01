package resource_manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceContainer_RegisterResource(t *testing.T) {
	rc := NewResourceContainer()
	key := "test-resource"
	obj := closableCreator{"test-object"}

	rc.RegisterResource(key, obj)
	rc.Init()

	ret := rc.GetResource(key)
	assert.NotNil(t, ret)

	retCC, ok := ret.(closableCreator)
	assert.True(t, ok)
	assert.EqualValues(t, obj.Name, retCC.Name)

	assert.Panics(t, func() {
		rc.RegisterResource(key, obj)
	}, "should panic because duplicate key")
}

func TestResourceContainer_RegisterCreator(t *testing.T) {
	rc := NewResourceContainer()
	key := "test-resource"

	var expect closableCreator
	var obj closableCreator
	rc.RegisterCreator(key, ResourceCreator{
		CreateFunc: func() (closableCreator, error) {
			var err error
			expect, err = newCreator()
			return expect, err
		},
		Receiver: &obj,
	})
	rc.Init()

	ret := rc.GetResource(key)
	assert.NotNil(t, ret)

	retCC, ok := ret.(closableCreator)
	assert.True(t, ok)
	assert.EqualValues(t, expect.Name, retCC.Name)
	assert.EqualValues(t, expect.Name, obj.Name)
}

type initializableResource struct {
	Count int
}

const defaultTestCount = 100

func (ir *initializableResource) Close() error {
	return nil
}

func (ir *initializableResource) Init() error {
	ir.Count = defaultTestCount
	return nil
}

func TestResourceContainer_Init(t *testing.T) {
	rc := NewResourceContainer()
	key := "test-resource"

	var expect *initializableResource
	var obj *initializableResource
	rc.RegisterCreator(key, ResourceCreator{
		CreateFunc: func() (*initializableResource, error) {
			var err error
			expect = &initializableResource{}
			return expect, err
		},
		Receiver: &obj,
	})
	rc.Init()

	ret := rc.GetResource(key)
	assert.NotNil(t, ret)

	retCC, ok := ret.(*initializableResource)
	assert.True(t, ok)
	assert.EqualValues(t, defaultTestCount, expect.Count)
	assert.EqualValues(t, expect.Count, retCC.Count)
	assert.EqualValues(t, expect.Count, obj.Count)
}
