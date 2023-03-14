# ETSim

ETSim is a simple command-line simulation game that allows you to simulate an alien invasion on a map of cities. The game is written in Go and is run entirely from the command line.

## Installation

To install ETSim, you need to have Go installed on your system. Once you have installed Go, you can download and install ETSim with the following command:

```bash
go get github.com/koneal2013/etsim/etsimcmd
```

This will download and install the ETSim command-line tool.

## Usage

Once you have installed ETSim, you can run it from the command line with the following command:

```bash
etsim [flags]
```

The following flags are available:

* `-n`: The number of aliens to be spawned. (default 10)
* `-m`: The path to the map file. (default "map.txt")

## Gameplay

The game begins by spawning a number of aliens on the map. Each alien is randomly placed in a city. The aliens then move around the map, randomly choosing neighboring cities to move to. If two aliens end up in the same city, they will fight, and both aliens will be destroyed, along with the city they were in.

As the game progresses, the aliens will destroy cities and fight each other, until either all of the aliens have moved 10,000 times, have been destroyed or there are no more cities left on the map.

## Contact

Please feel free to email me at [koneal2013@gmail.com](mailto:koneal2013@gmail.com) with questions regarding my solution.

