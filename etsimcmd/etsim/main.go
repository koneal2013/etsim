package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/koneal2013/etsim/etsimcmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	const (
		exitCodeOK    = 0
		exitCodeError = 1
	)
	// Parse command-line arguments
	if len(os.Args) < 2 {
		os.Exit(exitCodeError)
	}

	numAliens, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error converting %s to int: %s\n", os.Args[1], err)
		os.Exit(exitCodeError)
	}

	world := etsimcmd.New(numAliens, "map.txt")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go world.StartSimulation()
	go world.DestroyCitiesAndAliens(&wg)
	fmt.Printf("simulation started %c %c ...\n", '👽', '🛸')
	wg.Wait()
	fmt.Printf("simulation complete %c %c ...\n", '👽', '🔥')
	world.PrintRemainingCities()
	return exitCodeOK
}
