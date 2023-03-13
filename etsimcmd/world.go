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
	citiesToDestroy chan *City
	destroyedCities sync.Map
	sync.Mutex
}

// New creates a new World struct and populates the cities defined in the 'worldMapPath' parameter with the number of aliens specified
// by the 'numAliens' parameter
func New(numAliens int, worldMapPath string) *World {
	cities := readMapFile(worldMapPath)
	return &World{
		cities:          cities,
		aliens:          createAliens(numAliens, cities),
		citiesToDestroy: make(chan *City),
	}
}

// StartSimulation simulates alien invasions by randomly moving each alien to a new city up to 10k times or until
// there are no aliens left. If two aliens land in one city, the aliens destroy each other as well as the city.
func (w *World) StartSimulation() {
	defer close(w.citiesToDestroy)
	for i := 0; len(w.aliens) > 0 && i < 10000; i++ {
		// Move each alien to a neighboring city. If the city has more than one alien occupant, start a battle
		// w.Lock()
		for alien := range w.aliens {
			if alien.current != nil {
				currAlienNeighbors := alien.current.neighbors
				var neighborDirections []string
				for neighborDirection := range currAlienNeighbors {
					neighborDirections = append(neighborDirections, neighborDirection)
				}
				if len(neighborDirections) == 0 {
					continue
				}
				newCityDirection := neighborDirections[rand.Intn(len(neighborDirections))]
				newCityName, _ := currAlienNeighbors[newCityDirection]
				city, ok := w.cities[newCityName]
				if !ok {
					continue
				}
				if city.occupants[0] == nil || city.occupants[1] == nil {
					city.invade(alien)
					continue
				}
				w.citiesToDestroy <- city
			}
		}
		// w.Unlock()
	}
}

func (w *World) DestroyCitiesAndAliens(wg *sync.WaitGroup) {
	defer wg.Done()
	for city := range w.citiesToDestroy {
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
		for _, neighborName := range city.neighbors {
			if neighbor, ok := w.cities[neighborName]; ok {
				delete(neighbor.neighbors, getDirectionToNeighbor(neighbor, neighborName))
			}
		}
		// Remove destroyed cities from aliens' current city
		for alien := range w.aliens {
			if alien.current != nil && alien.current.name == city.name {
				alien.current = nil
			}
		}
		delete(w.aliens, city.occupants[0])
		delete(w.aliens, city.occupants[1])
		delete(w.cities, city.name)
	}
}

// PrintRemainingCities prints the cities that are remaining after the simulation has completed to standard output
func (w *World) PrintRemainingCities() {
	// Print out the remaining cities and their neighboring cities
	fmt.Println("Remaining cities: ")
	for cityName, city := range w.cities {
		var neighborNames []string
		for direction, neighborName := range city.neighbors {
			neighborNames = append(neighborNames, fmt.Sprintf("%s=%s", direction, neighborName))
		}
		fmt.Printf("%s %s\n", cityName, strings.Join(neighborNames, " "))
	}
}

// readMapFile Reads in the map file using the worldMapPath parameter passed to the New function
func readMapFile(worldMapPath string) map[string]*City {
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
		if len(parts) < 2 {
			continue
		}
		cityName := parts[0]
		neighbors := make(map[string]string)
		for _, part := range parts[1:] {
			dirCity := strings.Split(part, "=")
			if len(dirCity) != 2 {
				continue
			}
			direction := dirCity[0]
			neighborName := dirCity[1]
			neighbors[direction] = neighborName
		}
		cities[cityName] = &City{name: cityName, neighbors: neighbors}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	return cities
}

func getDirectionToNeighbor(current *City, neighborName string) string {
	for direction, name := range current.neighbors {
		if name == neighborName {
			return direction
		}
	}
	return ""
}

// createAliens creates the aliens and place them randomly on the map
func createAliens(numAliens int, cities map[string]*City) map[*Alien]struct{} {
	aliens := make(map[*Alien]struct{}, numAliens)
	cityNames := make([]string, 0, len(cities))
	for cityName := range cities {
		cityNames = append(cityNames, cityName)
	}
	rand.Shuffle(len(cityNames), func(i, j int) {
		cityNames[i], cityNames[j] = cityNames[j], cityNames[i]
	})
	for i := 0; i < numAliens; i++ {
		var currCity *City
		for _, cityName := range cityNames {
			city := cities[cityName]
			if !city.full {
				currCity = city
				break
			}
		}
		if currCity == nil {
			break
		}
		idx := rand.Intn(2)
		alien := &Alien{id: i + 1, current: currCity}
		currCity.occupants[idx] = alien
		aliens[alien] = struct{}{}
		currCity.full = currCity.occupants[0] != nil && currCity.occupants[1] != nil
	}
	return aliens
}
