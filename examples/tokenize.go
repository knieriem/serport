package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/knieriem/g/text"
)

// Read a line from stdin, and split it into
// fields using strings.Tokenize.

func main() {
	r := bufio.NewReader(os.Stdin)
	l, _, _ := r.ReadLine()

	for _, s := range text.Tokenize(string(l)) {
		fmt.Printf("«%s»\n", s)
	}
}
