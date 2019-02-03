package resource_manager

import (
	"fmt"
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
	c := NewResourceCreator(func() (closableCreator, error) {
		var err error
		expect, err = newCreator()
		return expect, err
	}, &obj)
	rc.RegisterCreator(key, c)
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
	fmt.Printf("close call")
	ir.Count = 0
	return nil
}

func (ir *initializableResource) Init() error {
	fmt.Println("init call")
	ir.Count = defaultTestCount
	return nil
}

func TestResourceContainer_Init(t *testing.T) {
	rc := NewResourceContainer()
	key := "test-resource"

	var expect *initializableResource
	var obj *initializableResource
	c := NewResourceCreator(func() (*initializableResource, error) {
		var err error
		expect = &initializableResource{}
		return expect, err
	}, &obj)
	rc.RegisterCreator(key, c)
	rc.Init()

	ret := rc.GetResource(key)
	assert.NotNil(t, ret)

	retCC, ok := ret.(*initializableResource)
	assert.True(t, ok)
	assert.EqualValues(t, defaultTestCount, expect.Count)
	assert.EqualValues(t, expect.Count, retCC.Count)
	assert.EqualValues(t, expect.Count, obj.Count)

	rc.Close()
	assert.EqualValues(t, 0, obj.Count)
}
