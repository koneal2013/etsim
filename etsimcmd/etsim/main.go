package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/koneal2013/etsim/etsimcmd"
)

const (
	defaultNumOfAliens = 10
)

var (
	numOfAliens  uint16
	worldMapPath string
)

func main() {
	cmd := &cobra.Command{
		Use: "etsim [flags]",
		Long: "ETSim is a simple command-line simulation game that allows you to simulate an alien invasion on a " +
			"\nmap of cities. The game is written in Go and is run entirely from the command line.",
		RunE: run,
	}
	cmd.Flags().StringVarP(&worldMapPath, "worldMapPath", "m", "map.txt", "path to world map")
	cmd.Flags().Uint16VarP(&numOfAliens, "numOfAliens", "n", defaultNumOfAliens, "number of aliens")

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
	if numOfAliens <= 0 {
		return fmt.Errorf("%c numAliens must be greater than 0 to run simulation", errorEmoji)
	}

	world := etsimcmd.New(numOfAliens, worldMapPath)

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
