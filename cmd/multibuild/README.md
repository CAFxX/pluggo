# multibuild
A small tool to build and link together multiple go packages

## Purpose
This tool was designed to work with [pluggo](https://github.com/CAFxX/pluggo),
to perform linking of plugins with the application/library they're intended for.

It can also be used in other cases, where you're importing a package only for
its side-effects, e.g.:

- To link additional `database/sql` drivers
- To link `pprof` at compile-time
- ...

## Installation
go install github.com/cafxx/pluggo/cmd/multibuild

## Usage
`multibuild <mainPkg> <importPkg1> <importPkg2> ...` will build `<mainPkg>` (as
if `go build <mainPkg>` had been invoked) and link it together with the
additional `<importPkg>`s.

`multibuild -h` prints usage instructions.

## FAQ

### Why is this not part of `go build`?
Good question. It is literally beyond me.
