# ETSim 👽 🛸

ETSim is a simple command-line simulation game that allows you to simulate an alien invasion on a map of cities. The game is written in Go and is run entirely from the command line.

## Installation

To install ETSim, you need to have Go installed on your system. Once you have installed Go, you can download and install ETSim with the following command:

```bash
git clone https://github.com/koneal2013/etsim.git
```

This will download and install the ETSim command-line tool.

Change to the `etsim` (project root) directory, then run ETSim with the following command:

```bash
make start
```

This will start ETSim with the default values for `numOfAliens` and `worldMapPath`

Run the following command to produce an executable in the current working directory:

```bash
make build
```

## Usage

Once you have installed ETSim, you can run it from the command line with the following command:

```bash
./etsim [flags]
```

The following flags are available:

* `-n`: The number of aliens to be spawned. (default 10)
* `-m`: The path to the map file. (default "map.txt")
* `-h`: Help for etsim

## Makefile

* `make test`: Runs all tests in the project using the `go test` command.
* `make start`: Cleans the project, builds it, and runs the main program with any arguments passed in using the `./etsim $(ARGS)` command. The `ARGS` variable can be set to pass any desired arguments.
* `make build`: Runs the tests using the `make test` command, then runs the `golangci-lint run` command to perform a static analysis of the code, and finally builds the program using the `go build` command.
* `make clean`: Cleans the project using the `go clean -i` command.
* `make cover`: Runs all tests and generates a coverage report in the `.coverage` directory using the `go test -v -cover ./... -coverprofile .coverage/coverage.out` and `go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html` commands.
* `make covero`: Runs the `make cover` command and opens the generated coverage report in the default browser using the `open .coverage/coverage.html` command.

## Gameplay

The game begins by spawning a number of aliens on the map. Each alien is randomly placed in a city. The aliens then move around the map, randomly choosing neighboring cities to move to. If two aliens end up in the same city, they will fight, and both aliens will be destroyed, along with the city they were in.

As the game progresses, the aliens will destroy cities and fight each other, until either all of the aliens have moved 10,000 times, all the aliens have been destroyed or there are no more cities left on the map.

## Assumptions

I made the following assumptions regarding the assignment:

* When spawning the aliens, if all the cities in the map are currently occupied and there are aliens that haven't been placed in a city, those aliens stay aboard the "mother ship" until there is a city with space available.
* City names will contain only string characters.
* The number of aliens that can be spawned for a simulation must be greater than 0 and less than or equal to 65535.
* Trapped aliens are considered dead and removed from active deployment.
* Only two aliens can be engaged in a fight in a particular city.



## Contact

Please feel free to email me at [koneal2013@gmail.com](mailto:koneal2013@gmail.com) with questions regarding my solution.

