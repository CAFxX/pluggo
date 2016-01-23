// Package pluggo provides a compile-time, in-process plugin framework. It
// allows to define interface-based extension points in your code, so that users
// of your code can plug in their modifications at compile time while keeping
// the application code and plugin code in completely separated packages or
// repositories.
package pluggo

import (
	"fmt"
	"sync"
)

// Factory functions should return instances of the interface appropriate for
// the extension point they have been registered for. The pluggo framework
// enforces no rules regarding the kind of interface factories should return:
// this is delegated to the contract between calling code and plugins.
type Factory func() interface{}

var plugins = make(map[string]Factory)
var lock = sync.RWMutex{}

// Register should be called from a plugin init() function to register a Factory
// for the named extension point. Naming of the extension points is delegated to
// the contract between calling code and plugins. It is illegal to register more
// than one Factory for the same extension point.
//
// It is strongly recommended that plugins perform no initialization in init()
// aside from registering with pluggo. Any initialization procedures should be
// deferred at least until the first call to the Factory function.
func Register(name string, f Factory) error {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := plugins[name]; exists {
		return fmt.Errorf("extension point %s has already a registered factory", name)
	}
	plugins[name] = f
	return nil
}

// Get is used to get an instance of the interface returned by the Factory
// registered for the specified extension point. The instantiation semantics
// are not defined by pluggo: they should be defined by the contract between
// calling code and plugins. It is safe to call Get concurrently from multiple
// goroutines.
func Get(name string) interface{} {
	lock.RLock()
	factory := plugins[name]
	lock.RUnlock()
	if factory == nil {
		return nil
	}
	return factory()
}
