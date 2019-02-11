package main

type Identity struct {
	Name string `json:"name"`
}

type FileScopeMarker struct {
	Identity
}
