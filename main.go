package main

import (
	"github.com/irpanzy/Task-Forge/config"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
}
