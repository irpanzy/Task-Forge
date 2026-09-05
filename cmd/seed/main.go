package main

import (
	"log"

	"github.com/irpanzy/Task-Forge/database/seed"
	"github.com/irpanzy/Task-Forge/internal/config"
)

func main() {
	log.Println("Memulai proses seeding database...")
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()

	log.Println("Proses seeding selesai.")
}
