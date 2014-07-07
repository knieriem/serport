// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/knieriem/g/registry"
)

func main() {
	key, err := registry.KeyLocalMachine.Subkey(os.Args[1:]...)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	values, err := key.Values()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for k, v := range values {
		fmt.Println(k, v.String())
	}
}
