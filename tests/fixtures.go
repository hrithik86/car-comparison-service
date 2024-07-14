package tests

import (
	"car-comparison-service/config"
	"flag"
)

func SetupFixtures() {
	flag.Parse()
	commandArgs := []string{"/app/main", "api", "/profiles/development.yml"}
	config.Load(commandArgs)
}
