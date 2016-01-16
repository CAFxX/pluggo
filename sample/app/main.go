package main

import (
	"github.com/cafxx/pluggo"
	"github.com/cafxx/pluggo/sample/app/extensions"
)

func main() {
	h := pluggo.Get("hello")
	if hello, ok := h.(extensions.Hello); ok {
		hello.Say()
	}
}
