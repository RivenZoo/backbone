package service_manager

// Runnable is required to implement by service.
type Runnable interface {
	Run() error
}

// Stoppable is optional to implement.
type Stoppable interface {
	Stop() error
}
