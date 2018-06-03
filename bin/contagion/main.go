package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	contagion "github.com/kentwait/contagiongo"
)

func main() {
	numCPUPtr := flag.Int("threads", runtime.NumCPU(), "number of CPU threads")
	loggerType := flag.String("logger", "csv", "data logger type (csv|sqlite)")
	seedNum := flag.Int64("seed", time.Now().UTC().UnixNano(), "random seed. Uses Unix time in nanoseconds as default")
	// benchmarkPtr := flag.String("benchmark", "", "Benchmark mode. Logs memory and wall time and saves to the specified path")
	flag.Parse()

	// Set random number
	rand.Seed(*seedNum)

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
	firstStart := time.Now()
	for i := 1; i <= conf.NumInstances(); i++ {
		log.Printf("starting instance %03d\n\n", i)
		start := time.Now()
		// Create a new logger for every realization
		var logger contagion.DataLogger
		switch *loggerType {
		case "csv":
			logger = contagion.NewCSVLogger(conf.LogPath(), i)
		case "sqlite":
			logger = contagion.NewSQLiteLogger(conf.LogPath(), i)
		default:
			log.Fatalf("%s is not a valid logger type (csv|sqlite)", *loggerType)
		}
		// Create a new simulation based on the epidemic model
		var sim contagion.EpidemicSimulation
		switch conf.SimParams.EpidemicModel {
		case "si":
			sim, err = contagion.NewSISimulation(conf, logger)
		case "sir":
			sim, err = contagion.NewSIRSimulation(conf, logger)
		case "endtrans":
			sim, err = contagion.NewEndTransSimulation(conf, logger)
		default:
			err = fmt.Errorf("epidemic model %s has not yet been implemented", conf.SimParams.EpidemicModel)
		}
		if err != nil {
			log.Fatalf("error creating a new simulation from the configuration file: %s", err)
		}
		sim.Run(i)
		log.Printf("Finished instance %03d in %s.\n\n", i, time.Since(start))
	}
	log.Printf("Completed all runs in %s.", time.Since(firstStart))
}
