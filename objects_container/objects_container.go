/*
* Objects Container will take care of all resource and service init and close.
*/
package objects_container

import (
	"github.com/RivenZoo/backbone/resources"
	"github.com/RivenZoo/backbone/services"
	"github.com/RivenZoo/injectgo"
)

var container = injectgo.NewContainer()

func Init() {
	resources.Init()
	services.Init()
	// populate injected objects
	container.Populate(nil)
}

func Close() {
	services.Close()
	resources.Close()
}

func GetObjectContainer() *injectgo.Container {
	return container
}
