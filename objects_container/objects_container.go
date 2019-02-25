/*
* Objects Container will take care of all provided objects init and close.
*/
package objects_container

import (
	"github.com/RivenZoo/injectgo"
)

var container = injectgo.NewContainer()

func Init() {
	// populate injected objects
	container.Populate(nil)
}

func Close() {
	container.Close()
}

func GetObjectContainer() *injectgo.Container {
	return container
}
