package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type City struct {
	name      string
	neighbors map[string]string
	occupants [2]*Alien
	full      bool
}

type Alien struct {
	id      int
	current *City
}

type World struct {
	cities          map[string]*City
	aliens          map[*Alien]struct{}
	citiesToDestroy chan *City
	destroyedCities sync.Map
	sync.Mutex
}

func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	numAliens := atoi(os.Args[1])

	// wait group for goroutines
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Read in the map file
	// Create the aliens and place them randomly on the map
	// Run the simulation until all aliens are destroyed or have moved 10,000 times
	cities := readMapFile("map.txt")
	world := &World{
		cities:          cities,
		aliens:          createAliens(numAliens, cities),
		citiesToDestroy: make(chan *City),
	}
	go world.destroyCitiesAndAliens(&wg)
	go world.startSimulation(&wg)
	wg.Wait()

	// Print out the remaining cities and their neighboring cities
	fmt.Println("Remaining cities: ")
	for cityName, city := range world.cities {
		var neighborNames []string
		for direction, neighborName := range city.neighbors {
			neighborNames = append(neighborNames, fmt.Sprintf("%s=%s", direction, neighborName))
		}
		fmt.Printf("%s %s\n", cityName, strings.Join(neighborNames, " "))
	}
}

func (c *City) invade(alien *Alien) {
	// get the index of the alien to remove
	idx := alien.id % 2
	// remove the alien from the city it currently occupies
	if alien.current != nil {
		alien.current.occupants[idx] = nil
	}

	// invade the new city
	c.occupants[idx] = alien
	alien.current = c
	return

}

func (w *World) startSimulation(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(w.citiesToDestroy)
	for i := 0; len(w.aliens) > 0 && i < 10000; i++ {
		// Move each alien to a neighboring city. If the city has more than one alien occupant, start a battle
		w.Lock()
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
		w.Unlock()
	}
}

func createAliens(numAliens int, cities map[string]*City) map[*Alien]struct{} {
	aliens := make(map[*Alien]struct{}, numAliens)
	cityNames := make([]string, 0, len(cities))
	for cityName := range cities {
		cityNames = append(cityNames, cityName)
	}
	for i := 0; i < numAliens; i++ {
		cityIdx := rand.Intn(len(cityNames))
		cityName := cityNames[cityIdx]
		currCity := cities[cityName]
		alien := &Alien{id: i, current: currCity}
		aliens[alien] = struct{}{}
	}
	return aliens
}

func (w *World) destroyCitiesAndAliens(wg *sync.WaitGroup) {
	defer wg.Done()
	for city := range w.citiesToDestroy {
		_, alreadyDestroyed := w.destroyedCities.LoadOrStore(city.name, true)
		if alreadyDestroyed {
			// City has already been destroyed, skip it
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

func getDirectionToNeighbor(current *City, neighborName string) string {
	for direction, name := range current.neighbors {
		if name == neighborName {
			return direction
		}
	}
	return ""
}

func readMapFile(filename string) map[string]*City {
	cities := make(map[string]*City)
	file, err := os.Open(filename)
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

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting %s to int: %s\n", s, err)
		os.Exit(1)
	}
	return i
}
