package main

import (
	"os"

	"github.com/pacificporter/goshim"
)

func main() {
	os.Exit(goshim.Run(os.Args[1:]))
}
