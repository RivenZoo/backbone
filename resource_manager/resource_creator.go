package resource_manager

import (
	"errors"
	"github.com/RivenZoo/backbone/closure"
)

var errUnSupportCreatorFunc  = errors.New("resource creator must be func() T or func() (T, error), T must implements Closable")

// ResourceCreator used to create resource object.
//
// CreateFunc must be func() T or func() (T, error)
// T must implements Closable.
// If T implements Initializable, Init will be called.
//
// Receiver is optional
// If set, it must be *T
type ResourceCreator struct {
	closure.ObjectCreator
}

// NewResourceCreator parameter CreateFunc must be func() T or func() (T, error).
// Receiver is optional, If set, it must be *T.
func NewResourceCreator(createFunc interface{}, receiver interface{}) ResourceCreator {
	return ResourceCreator{
		ObjectCreator: closure.NewObjectCreator(createFunc, receiver),
	}
}

func (c ResourceCreator) createResource() (Closable, error) {
	obj, err := c.ObjectCreator.Create()
	if err != nil {
		return nil, err
	}
	ret, ok := obj.(Closable)
	if !ok {
		return nil, errUnSupportCreatorFunc
	}
	return ret, nil
}
