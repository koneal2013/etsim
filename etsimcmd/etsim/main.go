package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/koneal2013/etsim/etsimcmd"
)

const (
	numOfArgsExpected = 2
)

func main() {
	cmd := &cobra.Command{
		Use:   "etsim <numAliens> <worldMapPath>",
		Short: "Simulate alien invasions on a map",
		Args:  cobra.ExactArgs(numOfArgsExpected),
		RunE:  run,
	}
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	const (
		alienEmoji        = '👽'
		flyingSaucerEmoji = '🛸'
		fireEmoji         = '🔥'
		errorEmoji        = '🚨'
	)
	// Parse command-line arguments
	numAliens, err := strconv.Atoi(args[0])

	if err != nil {
		return fmt.Errorf("Error converting %s to int: %s\n", args[0], err)
	}

	if numAliens <= 0 {
		return fmt.Errorf("%c numAliens must be greater than 0 to run simulation", errorEmoji)
	}

	worldMapPath := args[1]

	world := etsimcmd.New(uint16(numAliens), worldMapPath)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go world.StartSimulation()
	go world.DestroyCitiesAndAliens(&wg)
	go world.DeployReserveAliens()

	fmt.Printf("simulation started %c %c ...\n", alienEmoji, flyingSaucerEmoji)
	wg.Wait()
	fmt.Printf("simulation complete %c %c ...\n", alienEmoji, fireEmoji)
	world.PrintRemainingCities()

	return nil
}
