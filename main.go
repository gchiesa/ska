package main

import "github.com/gchiesa/ska/cmd"

var version = "development"

func main() {
	_ = cmd.Execute(version)
}
