package main

import (
	"bufio"
	"fmt"
	"github.com/knieriem/g/text"
	"os"
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
