package service_manager

import (
	"errors"
	"github.com/RivenZoo/backbone/closure"
)

var errUnSupportCreatorFunc = errors.New("service creator must be func() T or func() (T, error), T must implements Runnable")

// ServiceCreator used to create service object.
//
// CreateFunc must be func() T or func() (T, error)
// T must implements Runnable.
// If T implements Stoppable, Stop will be called.
//
// Receiver is optional
// If set, it must be *T
type ServiceCreator struct {
	closure.ObjectCreator
}

// NewServiceCreator parameter CreateFunc must be func() T or func() (T, error).
// Receiver is optional, If set, it must be *T.
func NewServiceCreator(createFunc interface{}, receiver interface{}) ServiceCreator {
	return ServiceCreator{
		ObjectCreator: closure.NewObjectCreator(createFunc, receiver),
	}
}

func (c ServiceCreator) createService() (Runnable, error) {
	obj, err := c.ObjectCreator.Create()
	if err != nil {
		return nil, err
	}
	ret, ok := obj.(Runnable)
	if !ok {
		return nil, errUnSupportCreatorFunc
	}
	return ret, nil
}
