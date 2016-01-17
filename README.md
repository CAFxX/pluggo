# pluggo
Golang compile-time, in-process plugin framework

[![GoDoc](https://godoc.org/github.com/CAFxX/pluggo?status.svg)](https://godoc.org/github.com/CAFxX/pluggo)

## Purpose
Pluggo allows you to define interface-based extension points in your code, so
that users of your code can plug in their modifications at compile time while
keeping the application code and plugin code in completely separated packages
and repositories.

Compared to RPC/IPC approaches, plugins using this framework run in the same
process as the application with no IPC/RPC overhead.

## How it works
Similarly as how `database/sql` drivers registers themselves: there's a
"extension point registry" (just a `map[string]func() interface{}`) where
plugins `Register` their factories for the appropriate extension points.
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
if userGreeter == nil {
  fmt.Printf("Hello %s!", username) // behavior in case no plugin was defined
} else {
  fmt.Print(userGreeter(username))
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
at runtime.

## FAQ

### Why did you write this?
Sometimes the functionalities you need in a upstream application or library are
not useful in the context of the upstream project. Sure, you can fork the
upstream but that creates a great deal of maintenance burden.

Extension points have negligible performance overhead unless used in very tight
loops (where, most likely, you shouldn't use them anyway). Having extensions
points defined upstream may prove less controversial and will allow any
interested users to easily provide their own extensions.

### Can it load plugins at runtime?
Short answer: no.

Longer answer: something could be hacked together using a mixture of this
approach, CGO and LD_PRELOAD. But nothing of the sort has been implemented here:
for now you can only load plugins at compile time.

That being said, you can link multiple plugins at compile time and then choose
which ones to use at runtime. See the factory pattern above for a generic
example of how to do it.

While not being able to load plugins at runtime is a limitation, it has some
clear upsides: you are effectively vendoring plugins so you sidestep all kind
of version incompatibility issues (DLL hell, anyone?), you maintain the "single
binary" nature of compiled Go programs and, most importantly, all Go tooling
keep working correctly.

### Can I supply some parameters/configurations to a plugin?
Not right now, but I'm considering it. I'd like to avoid feature creep if
possible, so I'm taking some time to come up with a minimal yet flexible design.

## License
[MIT](LICENSE)
