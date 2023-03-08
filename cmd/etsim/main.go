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
}

type Alien struct {
	id      int
	current *City
	alive   bool
}

type World struct {
	cities          map[string]*City
	aliens          []*Alien
	citiesToDestroy chan *City
	sync.RWMutex
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
		RWMutex:         sync.RWMutex{},
	}
	go world.destroyCities(&wg)
	go world.startSimulation(&wg)
	wg.Wait()

	// Print out the remaining cities and their neighboring cities
	for cityName, city := range world.cities {
		neighborNames := make([]string, 0, len(city.neighbors))
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
	alien.current.occupants[idx] = nil

	// invade the new city
	c.occupants[idx] = alien
	alien.current = c
	return

}

func (w *World) startSimulation(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(w.citiesToDestroy)
	for i := 0; i < 10000 && len(w.aliens) > 0; i++ {
		// Move each alien to a neighboring city. If the city has more than one alien occupant, start a battle
		for _, alien := range w.aliens {
			if alien.current != nil && alien.alive {
				neighborDirections := make([]string, 0, len(alien.current.neighbors))
				for direction := range alien.current.neighbors {
					neighborDirections = append(neighborDirections, direction)
				}
				if len(neighborDirections) == 0 {
					// fmt.Println(fmt.Sprintf("Alien %d is trapped in %s", alien.id, alien.current.name))
					continue
				}
				newCityDirection := neighborDirections[rand.Intn(len(neighborDirections))]
				city := w.cities[alien.current.neighbors[newCityDirection]]
				if city.occupants[0] == nil || city.occupants[1] == nil {
					w.cities[alien.current.neighbors[newCityDirection]].invade(alien)
				} else {
					w.initiateBattle(city)
				}
			}
		}
	}
}

func createAliens(numAliens int, cities map[string]*City) []*Alien {
	aliens := make([]*Alien, numAliens)
	for i := 0; i < numAliens; i++ {
		cityNames := make([]string, 0, len(cities))
		for cityName := range cities {
			cityNames = append(cityNames, cityName)
		}
		cityName := cityNames[rand.Intn(len(cityNames))]
		currCity := cities[cityName]
		aliens[i] = &Alien{id: i + 1, current: currCity, alive: true}
	}
	return aliens
}

func (w *World) initiateBattle(city *City) {
	w.citiesToDestroy <- city
	city.occupants[0].alive = false
	city.occupants[1].alive = false
	fmt.Printf("%s has been destroyed by alien %d and alien %d!\n",
		city.name, city.occupants[0].id, city.occupants[1].id)
}

func (w *World) destroyCities(wg *sync.WaitGroup) {
	defer wg.Done()
	for city := range w.citiesToDestroy {
		w.Lock()
		delete(w.cities, city.name)
		for _, neighborCityName := range city.neighbors {
			for direction, neighborDirectionCityName := range w.cities[neighborCityName].neighbors {
				if neighborDirectionCityName == city.name {
					delete(w.cities[neighborCityName].neighbors, direction)
				}
			}
		}
		// Remove destroyed cities from aliens' current city
		for _, alien := range w.aliens {
			if alien.current != nil && alien.current.name == city.name {
				alien.current = nil
			}
		}
		w.Unlock()
	}
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
