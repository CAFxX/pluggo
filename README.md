# pluggo
Golang compile-time, in-process plugin framework

[![GoDoc](https://godoc.org/github.com/CAFxX/pluggo?status.svg)](https://godoc.org/github.com/CAFxX/pluggo)
[![Build Status](https://travis-ci.org/CAFxX/pluggo.svg?branch=master)](https://travis-ci.org/CAFxX/pluggo)
[![Coverage](http://gocover.io/_badge/github.com/CAFxX/pluggo)](http://gocover.io/github.com/CAFxX/pluggo)

## Purpose
Pluggo allows you to define interface-based extension points in your code, so
that users of your code can plug in their modifications at compile time while
keeping the application code and plugin code in completely separated packages
and repositories.

Compared to RPC/IPC approaches, plugins using this framework run in the same
process as the application with no IPC/RPC overhead.

## How it works
Similarly as how `database/sql` drivers registers themselves: pluggo keeps an
internal factory registry (just a `map[string]func() interface{}`) where
plugins `Register` their factories.
Application code at each extension point requests to the registry instances of
the plugin using `Get`. Application and plugins are then compiled and linked
together in the same executable with the [`multibuild`](cmd/multibuild) tool.

## Examples

### Extension point pattern
Let's say you have a place in your code where you greet the user:

```
fmt.Printf("Hello %s!\n", username)
```

Now, let's pretend you want plugins to be able to override what gets printed.

First, you define a public interface that plugins can implement:

```
type Greeter interface {
  Greet(who string) string
}
```

then you change your greeting code to:

```
userGreeter := pluggo.Get("userGreeter").(Greeter)

if userGreeter != nil {
  fmt.Print(userGreeter.Greet(username))
} else {
  fmt.Printf("Hello %s!", username) // behavior in case no plugin was defined
}
```

Now if somebody else wants to implement a plugin to change the greeting they
just have to implement the interface and register it:

```
func init() {
	pluggo.Register("userGreeter", func() interface{} {
		return &shouter{}
	})
}

type shouter struct {
}

func (*shouter) Greet(who string) string {
	return fmt.Sprintf("WASSUP %s!!!", strings.ToUpper(who))
}
```

To enable the plugin you have to link it in at compile time. To simplify this
process a very rudimental tool called `multibuild` is provided. This tool builds
and links together a main package (the application) with one or more additional
packages (the plugins).

```
$ multibuild appMainPkg pluginPkg1 pluginPkg2 ...
```

You can have a look at the [sample](sample) directory for a ready-to-run
example.

Note that if a unknown extension point name (or `""`) is requested, `Get` will
simply return `nil`.

### Factory pattern

Whereas in the extension point pattern plugin code is responsible to register
for a specific extension point, in the factory pattern plugins register
themselves with a unique name and choosing which one to use for a certain
extension point is delegated to the application.

In the previous example, this would mean replacing the plugin instantiation code
from:

```
userGreeter := pluggo.Get("userGreeter").(Greeter)
```

to

```
userGreeter := pluggo.Get(conf.plugins.userGreeter).(Greeter)
```

where `conf.plugins.userGreeter` will likely come from the configuration
mechanism in use by your application. This allows to choose which plugin to use
at runtime. Note that if a unknown plugin name (or `""`) is requested, `Get`
will simply return `nil`.

## FAQ

### Why did you write this?
Sometimes the functionalities you need in a upstream application or library are
not useful in the context of the upstream project. Sure, you can fork the
upstream but that creates a great deal of maintenance burden.

Extension points have negligible performance overhead unless used in very tight
loops (where, most likely, you shouldn't use them anyway). Having extensions
points defined upstream may prove less controversial and will allow any
interested users to easily provide their own extensions.

### Can it dynamically link/unlink plugins at runtime?
Short answer: no. Plugins are linked at compile-time and can't be unlinked.

Longer answer: something could be hacked together using a mixture of this
approach, CGO and LD_PRELOAD. But nothing of the sort has been implemented here:
for now you can only load plugins at compile time.

That being said, you can link multiple plugins at compile time and then choose
which ones to load at runtime. See the "Can it dynamically load/unload plugins
at runtime?" question below.

While not being able to link plugins at runtime is a limitation, it has some
clear upsides: you are effectively vendoring plugins so you sidestep all kind
of version incompatibility issues (DLL hell, anyone?), you maintain the "single
binary" nature of compiled Go programs and, most importantly, all Go tools keep
working correctly (test, pprof, etc.).

### Can it dynamically load/unload plugins at runtime?
Plugins all only registered at process initialization time, but normally this is
really quick (it boils down to inserting one entry per factory in a map).

Plugin instantiation ("loading") on the other side is delegated to the interface
between the calling code (the application) and plugins. Because instantiation
is triggered by the calling code, it may never happen for a certain plugin if
the calling code decides that it is not needed (normally in response to
configuration or user input). See the factory pattern above for a generic
example of how to do it.

Similarly, plugin unloading is delegated to the interface between the calling
code and the plugins: if a plugin is able to shutdown and clean after itself
(e.g. terminate its goroutines, remove global state, etc.) all memory used by
the plugin instance should be eventually reclaimed by the go GC.

### Can I supply some parameters/configurations to a plugin?
It's not part of the framework right now, but I'm considering it. I'd like to
avoid feature creep if possible, so I'm taking some time to come up with a
minimal yet flexible design.

Keep in mind that nothing prevents you from designing the interface that plugins
should implement to enable providing configuration to plugin implementations.

Recycling the `Greeter` example above, you could add a `Init` functions that has
to be called once after instantiation:

```
type Greeter interface {
  Init(conf GreeterConfig)
  Greet(who string) string
}

type GreeterConfig struct {
  // ...
}
```

```
userGreeter := pluggo.Get("userGreeter").(Greeter)

if userGreeter != nil {
  err := userGreeter.Init(userGreeterConfig)
  if err != nil {
    // handle error
  }

  fmt.Print(userGreeter.Greet(username))
} else {
  fmt.Printf("Hello %s!", username) // behavior in case no plugin was defined
}
```

### Do I have to instantiate the plugin every time I use it?
It depends on the contract between calling code and plugins. Pluggo does not
dictate any convention about this: if you define the contract to be "the plugin
is instantiated only once per process" you can definitely instantiate it once
and reuse it multiple times.

### Isn't all of this very low-level?
Yes, it is. Pluggo is just the foundation of a plugin system right now. But
because every additional feature risks being very opinionated I'm inclined to
keep this library small and well-scoped and build a higher-level framework on
top of it.

A possible way forward is to define a set of strictly optional interfaces that
plugins can choose to implement to respond to standard plugin lifecycle events.
This would be very generic and completely optional.

## Potential directions
- Implement lifecycle interfaces
  - start plugin
  - configure plugin
  - enumerate plugin interfaces
  - start plugin instance
  - stop plugin instance
  - stop plugin
- Implement plugin enumeration

## License
[MIT](LICENSE)
