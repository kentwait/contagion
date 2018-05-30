package main

import (
	"flag"
	"log"
	"runtime"

	contagion "github.com/kentwait/contagiongo"
)

func main() {
	numCPUPtr := flag.Int("threads", runtime.NumCPU(), "Number of CPU threads")
	// logIntervalPtr := flag.Int("interval", 5, "Log every n seconds")
	// benchmarkPtr := flag.String("benchmark", "", "Benchmark mode. Logs memory and wall time and saves to the specified path")
	flag.Parse()

	// Set number of CPUs to be used
	runtime.GOMAXPROCS(*numCPUPtr)

	// Load config file
	tomlPath := flag.Arg(0)
	spec, err := contagion.LoadSingleHostConfig(tomlPath)
	if err != nil {
		log.Fatal(err)
	}
	// Validate
	err = spec.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// for r := 1; r <= spec.NumReplicates; r++ {
	// 	// Init recorder for each replicate
	// 	recorder := contagion.NewCSVCache(fmt.Sprintf("%s.%02d", spec.Simulation.CachePath, r))
	// 	recorder.DeleteAll()

	// 	// Create simulation
	// 	sim, err := spec.Create()
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	err = contagion.InitSim(sim)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	err = contagion.RunSim(sim, spec.Simulation.Generations, r, recorder)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	fmt.Println(r)
	// }
}
