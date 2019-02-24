package services

import (
	"github.com/RivenZoo/backbone/service_manager"
)

var serviceContainer = service_manager.NewServiceContainer()

func GetServiceContainer() *service_manager.ServiceContainer {
	return serviceContainer
}

func Init() {
	serviceContainer.Init()
}

func Close() {
	serviceContainer.Close()
}

func RunServices() {
	serviceContainer.RunServices()
}
