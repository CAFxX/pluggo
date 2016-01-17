# multibuild
A small tool to build and link together multiple go packages

## Purpose
This tool was designed to work with [pluggo](../..), to perform linking of
plugins with the application/library they're intended for.

It can also be used in other cases, where you're importing a package only for
its side-effects, e.g.:

- To link additional `database/sql` drivers
- To link `pprof` at compile-time
- ...
