package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/koneal2013/etsim/etsimcmd"
)

func main() {
	cmd := &cobra.Command{
		Use:   "etsim <numAliens> <worldMapPath>",
		Short: "Simulate alien invasions on a map",
		Args:  cobra.ExactArgs(2),
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
	)
	// Parse command-line arguments
	numAliens, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Error converting %s to int: %s\n", args[0], err)
	}
	worldMapPath := args[1]

	world := etsimcmd.New(numAliens, worldMapPath)
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
