# Pluggo
Golang compile-time, in-process plugin framework

## Purpose
Pluggo allows you to define interface-based extension points in your code, so
that users of your code can plug in their modifications at compile time without
modifying your code.

## How it works
Similarly as how `database/sql` drivers registers themselves: there's a
"extension point registry" (just a `map[string]func() interface{}`) where
plugins `Register` their factories for the appropriate extension points.
Application code at each extension point requests to the registry instances of
the plugin using `Get`.

## Example
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

You can have a look at the [sample](sample/README.md) directory for a
ready-to-run example.

## FAQ

### Can it load plugins at runtime?
Short answer: no.

Longer answer: something could be hacked together using a mixture of this
approach, CGO and LD_PRELOAD. But nothing of the sort has been implemented here:
for now you can only load plugins at compile time.

That being said, you can link multiple plugins at compile times and then choose
which ones to use at runtime.

### Can I supply some parameters/configurations to a plugin?
Not right now, but I'm considering it. I'd like to avoid feature creep if
possible, so I'm taking some time to come up with a minimal yet flexible design.

## License
MIT
