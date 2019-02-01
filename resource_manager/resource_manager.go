package resource_manager

import (
	"errors"
	"io"
)

var (
	errDuplicateResourceName = errors.New("duplicate resource name")
	errUnsupportResourceType = errors.New("unsupport resource object type")
)

type Closable interface {
	io.Closer
}

type Initializable interface {
	Closable
	Init() error
}

type ResourceCreator struct {
	// CreateFunc must be func() T or func() (T, error)
	// T must be Closable or Initializable
	CreateFunc interface{}
	// Receiver is optional
	// If set, it must be *T
	Receiver interface{}
}

type ResourceContainer struct {
	resourceMap   map[string]Closable
	initializable []Initializable
	creators      []ResourceCreator
}

// RegisterResource register resource object by name.
// Parameter object should be Closable or Initializable.
// Panics if name is duplicated.
func (rc *ResourceContainer) RegisterResource(name string, object interface{}) {
	if _, ok := rc.resourceMap[name]; ok {
		panic(errDuplicateResourceName)
	}
	switch obj := object.(type) {
	case Closable:
		rc.resourceMap[name] = obj
	case Initializable:
		rc.resourceMap[name] = nil
		rc.initializable = append(rc.initializable, obj)
	default:
		panic(errUnsupportResourceType)
	}
}

// RegisterResource register resource creator by name.
// Resource object created by CreateFunc when `Init` called.
// If Receiver not nil, resource object will be assigned to Receiver.
// Panics if name is duplicated.
func (rc *ResourceContainer) RegisterCreator(name string, creator ResourceCreator) {
	if _, ok := rc.resourceMap[name]; ok {
		panic(errDuplicateResourceName)
	}

}

// GetResource return nil if no such resource provided.
func (rc *ResourceContainer) GetResource(name string) interface{} {
	return nil
}

// Init first create all resource objects by calling create func.
// Then call Init if resource object is Initializable.
// Panics if any error occurs.
func (rc *ResourceContainer) Init() {
	return
}

// Close all resource objects.
func (rc *ResourceContainer) Close() {
	return
}
