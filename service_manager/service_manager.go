package service_manager

import (
	"context"
	"errors"
	"sync"
)

var (
	errDuplicateServiceName = errors.New("duplicate service name")
)

type namedServiceCreator struct {
	name string
	ServiceCreator
}

type ServiceContainer struct {
	serviceMap        map[string]Runnable
	creators          []namedServiceCreator
	stoppableServices []Stoppable
	stopCtx           context.Context
	cancelFunc        context.CancelFunc
}

func NewServiceContainer() *ServiceContainer {
	c := &ServiceContainer{
		serviceMap:        map[string]Runnable{},
		creators:          []namedServiceCreator{},
		stoppableServices: []Stoppable{},
	}
	c.stopCtx, c.cancelFunc = context.WithCancel(context.Background())
	return c
}

// RegisterService register service by name.
// Panics if duplicate name.
func (c *ServiceContainer) RegisterService(name string, svc Runnable) {
	if _, ok := c.serviceMap[name]; ok {
		panic(errDuplicateServiceName)
	}
	c.serviceMap[name] = svc
	if stoppableSvc, ok := svc.(Stoppable); ok {
		c.stoppableServices = append(c.stoppableServices, stoppableSvc)
	}
}

// RegisterCreator register service creator by name.
// Panics if duplicate name or creator is not valid.
func (c *ServiceContainer) RegisterCreator(name string, creator ServiceCreator) {
	if _, ok := c.serviceMap[name]; ok {
		panic(errDuplicateServiceName)
	}
	creator.Validate()
	c.creators = append(c.creators, namedServiceCreator{
		name:           name,
		ServiceCreator: creator,
	})
}

// GetService return nil if no such service provided.
func (c *ServiceContainer) GetService(name string) Runnable {
	return c.serviceMap[name]
}

// Init create all service object by calling create function.
// Panics if error occurs.
func (c *ServiceContainer) Init() {
	for _, creator := range c.creators {
		obj, err := creator.createService()
		if err != nil {
			panic(err)
		}
		c.serviceMap[creator.name] = obj
		if stoppableSvc, ok := obj.(Stoppable); ok {
			c.stoppableServices = append(c.stoppableServices, stoppableSvc)
		}
	}
}

// RunServices run service in seprate routine then wait all service to stop.
// Panics if Run return error.
func (c *ServiceContainer) RunServices() {
	wg := &sync.WaitGroup{}

	for _, svc := range c.serviceMap {
		c.runService(svc, wg)
	}
	<-c.stopCtx.Done()

	// wait stoppable service exit
	wg.Wait()
}

func (c *ServiceContainer) runService(svc Runnable, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := svc.Run(); err != nil {
			panic(err)
		}
	}()
}

// Close stop all stoppable service.
// Panics if error occurs.
func (c *ServiceContainer) Close() {
	for _, svc := range c.stoppableServices {
		if err := svc.Stop(); err != nil {
			panic(err)
		}
	}
	c.cancelFunc()
}
