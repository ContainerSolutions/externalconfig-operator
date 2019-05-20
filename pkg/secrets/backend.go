package secrets

import "fmt"

//Backend is a secret store backend
type Backend struct{}

//BackendIface is an interface to a Backend
type BackendIface interface {
	Init(...interface{}) error
	Get(string) (string, error)
}

// BackendInstances are instantiated backends
var BackendInstances map[string]BackendIface

// BackendFunctions is a map of functions that return Backends
var BackendFunctions map[string]func() BackendIface

// BackendInstantiate instantiates a Backend of type `backendType`
func BackendInstantiate(name string, backendType string) error {
	if BackendInstances == nil {
		BackendInstances = make(map[string]BackendIface)
	}

	function, found := BackendFunctions[backendType]
	if !found {
		return fmt.Errorf("Unkown backend type: %v", backendType)
	}

	BackendInstances[name] = function()

	return nil
}

// BackendRegister registers a new backend type with name `name`
// function is a function that returns a backend of that type
func BackendRegister(name string, function func() BackendIface) {
	if BackendFunctions == nil {
		BackendFunctions = make(map[string]func() BackendIface)
	}

	BackendFunctions[name] = function
}
