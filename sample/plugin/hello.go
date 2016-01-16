package helloplugin

import (
	"fmt"

	"github.com/cafxx/pluggo"
)

func init() {
	pluggo.Register("hello", func() interface{} {
		return &hello{}
	})
}

type hello struct {
}

func (*hello) Say() {
	fmt.Println("Hello pluggo")
}
