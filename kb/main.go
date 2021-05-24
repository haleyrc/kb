package main

import (
	"os"

	"github.com/haleyrc/kb/kb/app"
)

func main() {
	if err := app.Run(os.Args[1:]...); err != nil {
		panic(err)
	}
}
