package resource_manager

import "io"

// Closable is required to implement.
type Closable interface {
	io.Closer
}

// Initializable is optional to implement.
type Initializable interface {
	Init() error
}
