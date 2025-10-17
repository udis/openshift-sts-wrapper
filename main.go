package main

import (
	"os"

	"github.com/clobrano/ccoctl-sso/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
