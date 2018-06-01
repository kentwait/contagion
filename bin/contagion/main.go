package main

import (
	"flag"
	"log"
	"runtime"

	contagion "github.com/kentwait/contagiongo"
)

func main() {
	numCPUPtr := flag.Int("threads", runtime.NumCPU(), "number of CPU threads")
	loggerType := flag.String("logger", "csv", "data logger type (csv|sqlite)")
	// benchmarkPtr := flag.String("benchmark", "", "Benchmark mode. Logs memory and wall time and saves to the specified path")
	flag.Parse()

	// Set number of CPUs to be used
	runtime.GOMAXPROCS(*numCPUPtr)

	// Load config file
	configPath := flag.Arg(0)
	conf, err := contagion.LoadEvoEpiConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	// Validate configuration
	err = conf.Validate()
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= conf.NumInstances(); i++ {
		// Create a new logger for every realization
		var logger contagion.DataLogger
		switch *loggerType {
		case "csv":
			logger = contagion.NewCSVLogger(conf.LogPath(), i)
		case "sqlite":
			log.Fatalf("sqlite logger not yet implemented")
		default:
			log.Fatalf("%s is not a valid logger type (csv|sqlite)", *loggerType)
		}
		sim, err := contagion.NewSISimulation(conf, logger)
		if err != nil {
			log.Fatalf("error creating a new simulation from the configuration file: %s", err)
		}
		sim.Run(i)
	}
}
