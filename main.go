package main

import (
	"tommy/cmd"
	"tommy/pkg/config"
)

func main() {
	config.LoadEnv()
	cmd.Execute()
}
