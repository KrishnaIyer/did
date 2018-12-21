package main // import "go.htdvisser.nl/did/cmd/did"

import (
	"fmt"
	"os"

	"go.htdvisser.nl/did"
)

func main() {
	if err := did.Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
