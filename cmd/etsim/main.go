package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type City struct {
	name      string
	neighbors sync.Map
	occupants [2]atomic.Pointer[Alien]
}

type Alien struct {
	id      int
	current *City
}

type World struct {
	cities           sync.Map
	aliens           sync.Map
	numOfAliensAlive atomic.Int32
	citiesToDestroy  chan *City
}

func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	numAliens := atoi(os.Args[1])

	// wait group for goroutines
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Read in the map file
	// Create the aliens and place them randomly on the map
	// Run the simulation until all aliens are destroyed or have moved 10,000 times
	cities := readMapFile("map.txt")
	world := &World{
		numOfAliensAlive: atomic.Int32{},
		cities:           cities,
		aliens:           createAliens(numAliens, cities),
		citiesToDestroy:  make(chan *City),
	}
	world.numOfAliensAlive.Add(int32(numAliens))
	go world.destroyCities(&wg)
	world.startSimulation()
	wg.Wait()

	// Print out the remaining cities and their neighboring cities
	world.cities.Range(func(cityName, city any) bool {
		var neighborNames []string
		city.(*City).neighbors.Range(func(direction, neighborName any) bool {
			neighborNames = append(neighborNames, fmt.Sprintf("%s=%s", direction, neighborName))
			return true
		})
		fmt.Printf("%s %s\n", cityName, strings.Join(neighborNames, " "))
		return true
	})
}

func (c *City) invade(alien *Alien) {
	// get the index of the alien to remove
	idx := alien.id % 2
	// remove the alien from the city it currently occupies
	if alien.current != nil {
		alien.current.occupants[idx].Store(nil)
	}

	// invade the new city
	c.occupants[idx].Store(alien)
	alien.current = c
	return

}

func (w *World) startSimulation() {
	defer close(w.citiesToDestroy)
	for i := 0; w.numOfAliensAlive.Load() > 0 && i < 10000; i++ {
		// Move each alien to a neighboring city. If the city has more than one alien occupant, start a battle
		w.aliens.Range(func(key, value any) bool {
			alien := value.(*Alien)
			if alien.current != nil {
				currAlienNeighbors := alien.current.neighbors
				var neighborDirections []string
				currAlienNeighbors.Range(func(key, value any) bool {
					neighborDirections = append(neighborDirections, key.(string))
					return true
				})
				if len(neighborDirections) == 0 {
					// fmt.Println(fmt.Sprintf("Alien %d is trapped in %s", alien.id, alien.current.name))
					return true
				}
				newCityDirection := neighborDirections[rand.Intn(len(neighborDirections))]
				newCityName, _ := currAlienNeighbors.Load(newCityDirection)
				c, _ := w.cities.Load(newCityName)
				if c == nil {
					return true
				}
				city, _ := c.(*City)
				if city.occupants[0].Load() == nil || city.occupants[1].Load() == nil {
					city.invade(alien)
					return true
				}
				w.initiateBattle(city)
			}
			return true
		})
	}
}

func createAliens(numAliens int, cities sync.Map) sync.Map {
	aliens := sync.Map{}
	for i := 0; i < numAliens; i++ {
		var cityNames []string
		cities.Range(func(key, value any) bool {
			cityNames = append(cityNames, key.(string))
			return true
		})
		cityName := cityNames[rand.Intn(len(cityNames))]
		currCity, _ := cities.Load(cityName)
		aliens.Store(i, &Alien{id: i, current: currCity.(*City)})
	}
	return aliens
}

func (w *World) initiateBattle(city *City) {
	w.aliens.Delete(city.occupants[0].Load().id)
	w.aliens.Delete(city.occupants[1].Load().id)
	w.citiesToDestroy <- city
	w.numOfAliensAlive.Add(-2)
}

func (w *World) destroyCities(wg *sync.WaitGroup) {
	defer wg.Done()
	for city := range w.citiesToDestroy {
		// function to range over neighboring cities and remove roads
		f := func(key any, value any) bool {
			nCity, ok := w.cities.Load(value)
			if !ok {
				return true
			}
			neighboringCity := nCity.(*City)

			neighboringCity.neighbors.Range(func(key, value any) bool {
				if value == city.name {
					neighboringCity.neighbors.Delete(key)
				}
				return true
			})
			return true
		}
		w.cities.Delete(city.name)
		city.neighbors.Range(f)
		fmt.Printf("%s has been destroyed by alien %d and alien %d!\n",
			city.name, city.occupants[0].Load().id, city.occupants[1].Load().id)

		// Remove destroyed cities from aliens' current city
		w.aliens.Range(func(key, value any) bool {
			alien := value.(*Alien)
			if alien.current != nil && alien.current.name == city.name {
				alien.current = nil
			}
			return true
		})
	}
}

func readMapFile(filename string) sync.Map {
	cities := sync.Map{}
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
		neighbors := sync.Map{}
		for _, part := range parts[1:] {
			dirCity := strings.Split(part, "=")
			if len(dirCity) != 2 {
				continue
			}
			direction := dirCity[0]
			neighborName := dirCity[1]
			neighbors.Store(direction, neighborName)
		}
		cities.Store(cityName, &City{name: cityName, neighbors: neighbors})
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
