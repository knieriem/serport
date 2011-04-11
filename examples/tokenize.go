package main

import (
	"fmt"
	"encoding/line"
	"os"
	"github.com/knieriem/g/util"
)

// Read a line from stdin, and split it into
// fields using strings.Tokenize.

func main() {
	r := line.NewReader(os.Stdin, 256)
	l, _, _ := r.ReadLine()

	for _, s := range util.Tokenize(string(l)) {
		fmt.Printf("«%s»\n", s)
	}
}
