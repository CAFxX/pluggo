package pluggo

import (
	"fmt"
	"sync"
)

// Factory is the signature of functions that returns instances of the interface
// to be used at the extension points.
type Factory func() interface{}

var plugins = make(map[string]Factory)
var lock = sync.RWMutex{}

// Register should be called from a plugin init() function to register a Factory
// for the named extension point. It is illegal to register more than one
// Factory for the same extension point.
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
// registered for the specified extension point.
func Get(name string) interface{} {
	lock.RLock()
	factory := plugins[name]
	lock.RUnlock()
	if factory == nil {
		return nil
	}
	return factory()
}
