package resource_manager

import (
	"errors"
	"io"
)

var (
	errDuplicateResourceName = errors.New("duplicate resource name")
	errUnSupportResourceType = errors.New("resource object must implements Closable")
	errUnSupportCreatorFunc  = errors.New("resource creator must be func() T or func() (T, error), T must implements Closable")
)

type Closable interface {
	io.Closer
}

type Initializable interface {
	Init() error
}

type namedResourceCreator struct {
	name string
	ResourceCreator
}

type ResourceContainer struct {
	resourceMap       map[string]Closable
	needToInitObjects []Initializable
	creators          []namedResourceCreator
}

func NewResourceContainer() *ResourceContainer {
	return &ResourceContainer{
		resourceMap:       map[string]Closable{},
		needToInitObjects: []Initializable{},
		creators:          []namedResourceCreator{},
	}
}

// RegisterResource register resource object by name.
// Parameter object should be Closable.
// If object implements Initializable, Init will be called.
// Panics if name is duplicated.
func (rc *ResourceContainer) RegisterResource(name string, object interface{}) {
	if _, ok := rc.resourceMap[name]; ok {
		panic(errDuplicateResourceName)
	}
	switch obj := object.(type) {
	case Closable:
		rc.resourceMap[name] = obj
		if initObj, ok := object.(Initializable); ok {
			rc.needToInitObjects = append(rc.needToInitObjects, initObj)
		}
	default:
		panic(errUnSupportResourceType)
	}
}

// RegisterResource register resource creator by name.
// Resource object created by CreateFunc when `ResourceContainer.Init` called.
// If Receiver not nil, resource object will be assigned to Receiver.
// Panics if name is duplicated.
func (rc *ResourceContainer) RegisterCreator(name string, creator ResourceCreator) {
	if _, ok := rc.resourceMap[name]; ok {
		panic(errDuplicateResourceName)
	}
	creator.validate()
	rc.creators = append(rc.creators, namedResourceCreator{
		ResourceCreator: creator,
		name:            name,
	})
}

// GetResource return nil if no such resource provided.
func (rc *ResourceContainer) GetResource(name string) Closable {
	return rc.resourceMap[name]
}

// Init first create all resource objects by calling create func.
// Then call Init if resource object is Initializable.
// Panics if any error occurs.
func (rc *ResourceContainer) Init() {
	// first: create all resource objects
	for _, creator := range rc.creators {
		obj, err := creator.create()
		if err != nil {
			panic(err)
		}
		creator.setReceiver(obj)
		rc.resourceMap[creator.name] = obj
		if initObj, ok := obj.(Initializable); ok {
			rc.needToInitObjects = append(rc.needToInitObjects, initObj)
		}
	}
	// second: init all objects
	for _, initObj := range rc.needToInitObjects {
		if err := initObj.Init(); err != nil {
			panic(err)
		}
	}
	return
}

// Close all resource objects.
func (rc *ResourceContainer) Close() {
	for _, obj := range rc.resourceMap {
		if err := obj.Close(); err != nil {
			panic(err)
		}
	}
	return
}
