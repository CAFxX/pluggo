This directory contains a sample application (in `app/`) and plugin (in
`plugin/`). The app will call the Say() method on the interface instance
provided by the plugin.

# Try it out

Let's build just the application first and check that it prints nothing because
the plugin is not loaded:

```
$ go build github.com/cafxx/pluggo/sample/app
$ ./app
```

Let's build and link the application and the plugin: the plugin gets called and
it prints "Hello pluggo":

```
$ multibuild github.com/cafxx/pluggo/sample/app github.com/cafxx/pluggo/sample/plugin
$ ./app
Hello pluggo
```
