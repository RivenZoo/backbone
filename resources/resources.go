package resources

import "github.com/RivenZoo/backbone/resource_manager"

var resourceContainer = resource_manager.NewResourceContainer()

func GetResourceContainer() *resource_manager.ResourceContainer {
	return resourceContainer
}
