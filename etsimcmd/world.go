package etsimcmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
)

type World struct {
	cities          map[string]*City
	aliens          map[*Alien]struct{}
	alienShip       map[*Alien]struct{}
	citiesToDestroy chan *City
	destroyedCities sync.Map
	sync.Mutex
}

// New creates a new World struct and populates the cities defined in the 'worldMapPath'
// parameter with the number of aliens specified by the 'numAliens' parameter.
func New(numAliens uint16, worldMapPath string) *World {
	cities := readMapFile(worldMapPath)
	deployedAliens, reserveAliens := createAliens(numAliens, cities)

	return &World{
		cities:          cities,
		aliens:          deployedAliens,
		alienShip:       reserveAliens,
		citiesToDestroy: make(chan *City),
	}
}

// StartSimulation simulates alien invasions by randomly moving each alien to a new city up to 10k times or until
// there are no aliens left. If two aliens land in one city, the aliens destroy each other as well as the city.
func (w *World) StartSimulation() {
	const (
		maxAlienMoves = 10000
	)

	defer close(w.citiesToDestroy)

	for i := 0; len(w.aliens)+len(w.alienShip) > 0; i++ {
		var allMoved bool
		// Move each alien to a neighboring city. If the city has more than one alien occupant, start a battle
		for alien := range w.aliens {
			if alien.current != nil {
				w.Lock()

				if !alien.alive {
					delete(w.aliens, alien)
					w.Unlock()

					continue
				}

				w.Unlock()

				if alien.current.full {
					w.citiesToDestroy <- alien.current

					continue
				}

				currAlienNeighbors := alien.current.neighbors

				var neighborDirections []string

				for neighborDirection := range currAlienNeighbors {
					neighborDirections = append(neighborDirections, neighborDirection)
				}

				if len(neighborDirections) == 0 {
					alien.alive = false

					continue
				}

				newCityDirection := neighborDirections[rand.Intn(len(neighborDirections))]
				newCityName, _ := currAlienNeighbors[newCityDirection]

				w.Lock()
				city, ok := w.cities[newCityName]
				w.Unlock()

				if !ok && len(neighborDirections) == 1 {
					continue
				}

				if city.occupants[0] == nil || city.occupants[1] == nil {
					city.invade(alien)

					alien.moves++
					if alien.moves >= maxAlienMoves {
						allMoved = true
					}

					continue
				} else {
					w.citiesToDestroy <- city
				}
			}

			w.Lock()

			delete(w.aliens, alien)

			w.Unlock()
		}

		if allMoved || len(w.cities) == 0 || len(w.aliens) == 0 {
			break
		}
	}
}

func (w *World) DestroyCitiesAndAliens(wg *sync.WaitGroup) {
	defer wg.Done()

	for city := range w.citiesToDestroy {
		if city == nil {
			continue
		}

		_, alreadyDestroyed := w.destroyedCities.LoadOrStore(city.name, true)

		if alreadyDestroyed {
			// City has already been destroyed, skip it
			continue
		}

		if city.occupants[0] == nil || city.occupants[1] == nil {
			continue
		}

		fmt.Printf("%s has been destroyed by alien %d and alien %d!\n",
			city.name, city.occupants[0].id, city.occupants[1].id)
		// function to range over neighboring cities and remove roads
		for direction, neighborName := range city.neighbors {
			if neighbor, ok := w.cities[neighborName]; ok {
				delete(neighbor.neighbors, oppositeDirection(direction))
			}
		}

		w.Lock()
		// kill the aliens
		city.occupants[0].alive = false
		city.occupants[1].alive = false
		// destroy the city
		delete(w.cities, city.name)

		w.Unlock()
	}
}

func (w *World) DeployReserveAliens() {
	for {
		w.Lock()
		for alien := range w.alienShip {
			var currCity *City

			var cityNames []string

			if len(w.cities) == 0 {
				w.Unlock()

				return
			}

			for city := range w.cities {
				cityNames = append(cityNames, city)
			}

			currCity = findCityToInvade(cityNames, w.cities)

			if currCity != nil {
				w.aliens[alien] = struct{}{}

				currCity.invade(alien)

				delete(w.alienShip, alien)
			}
		}
		w.Unlock()

		if len(w.alienShip) == 0 {
			break
		}
	}
}

func oppositeDirection(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	default:
		return ""
	}
}

// PrintRemainingCities prints the cities that are remaining after the simulation has completed to standard output.
func (w *World) PrintRemainingCities() {
	// Print out the remaining cities and their neighboring cities
	fmt.Println("Remaining cities: ")

	if len(w.cities) == 0 {
		fmt.Println("none")
	}

	for cityName, city := range w.cities {
		var neighborNames []string
		for direction, neighborName := range city.neighbors {
			neighborNames = append(neighborNames, fmt.Sprintf("%s=%s", direction, neighborName))
		}

		fmt.Printf("%s %s\n", cityName, strings.Join(neighborNames, " "))
	}
}

// readMapFile Reads in the map file using the worldMapPath parameter passed to the New function.
func readMapFile(worldMapPath string) map[string]*City {
	const (
		partsLen = 2
	)

	cities := make(map[string]*City)

	file, err := os.Open(worldMapPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		neighbors := make(map[string]string)

		if len(parts) < partsLen {
			continue
		}

		for _, part := range parts[1:] {
			dirCity := strings.Split(part, "=")

			if len(dirCity) != partsLen {
				continue
			}

			direction := dirCity[0]
			neighborName := dirCity[1]
			neighbors[direction] = neighborName
		}

		cityName := parts[0]
		cities[cityName] = &City{name: cityName, neighbors: neighbors}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return cities
}

// createAliens creates the aliens and place them randomly on the map.
func createAliens(numAliens uint16, cities map[string]*City) (map[*Alien]struct{}, map[*Alien]struct{}) {
	ship := make(map[*Alien]struct{})
	aliens := make(map[*Alien]struct{}, numAliens)

	cityNames := make([]string, 0, len(cities))
	for cityName := range cities {
		cityNames = append(cityNames, cityName)
	}

	rand.Shuffle(len(cityNames), func(i, j int) {
		cityNames[i], cityNames[j] = cityNames[j], cityNames[i]
	})

	for i := uint16(0); i < numAliens; i++ {
		currCity := findCityToInvade(cityNames, cities)
		alien := &Alien{id: i + 1, alive: true}

		if currCity == nil {
			ship[alien] = struct{}{}
		} else {
			currCity.invade(alien)
			aliens[alien] = struct{}{}
		}
	}

	return aliens, ship
}

func findCityToInvade(cityNames []string, cities map[string]*City) *City {
	for _, cityName := range cityNames {
		city := cities[cityName]
		if !city.full {
			return city
		}
	}

	return nil
}
