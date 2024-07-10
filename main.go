package main

import (
	"car-comparison-service/appcontext"
	"car-comparison-service/config"
	"car-comparison-service/service/api"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	commandArgs := os.Args
	config.Load(commandArgs)

	//logger.SetupLogger()

	err := appcontext.Initiate(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	app := &cli.App{
		Name:  "Car comparison service",
		Usage: "Car comparison service",
		Commands: []*cli.Command{
			{
				Name:   "api",
				Usage:  "run api",
				Action: api.StartAPI,
			},
		},
	}

	if e := app.Run(os.Args); e != nil {
		fmt.Println(e)
	}
}
