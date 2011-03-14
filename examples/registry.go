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
	for k, v := range key.Values() {
		fmt.Println(k, v.String())
	}
}
