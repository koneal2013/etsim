package etsimcmd

import (
	"sync"
	"testing"
)

func TestNewWorld(t *testing.T) {
	numAliens := uint16(10)
	worldMapPath := "test_map.txt"
	world := New(numAliens, worldMapPath)

	if uint16(len(world.aliens)+len(world.alienShip)) != numAliens {
		t.Errorf("Expected %d aliens, but got %d", numAliens, len(world.aliens))
	}

	if len(world.cities) != 4 {
		t.Errorf("Expected 4 cities, but got %d", len(world.cities))
	}
}

func TestReadMapFile(t *testing.T) {
	worldMapPath := "test_map.txt"
	cities := readMapFile(worldMapPath)

	if len(cities) != 4 {
		t.Errorf("Expected 4 cities, but got %d", len(cities))
	}

	testCity := cities["Bar"]
	if len(testCity.neighbors) != 1 {
		t.Errorf("Expected Test City to have 1 neighbor, but got %d", len(testCity.neighbors))
	}
}

func TestCreateAliens(t *testing.T) {
	numAliens := uint16(10)
	cities := map[string]*City{
		"City 1": {name: "City 1", occupants: [2]*Alien{nil, nil}},
		"City 2": {name: "City 2", occupants: [2]*Alien{nil, nil}},
	}

	aliens, reserveAliens := createAliens(numAliens, cities)

	if len(aliens)+len(reserveAliens) != int(numAliens) && len(aliens) > 0 {
		t.Errorf("Expected %d aliens, but got %d", numAliens, len(aliens)+len(reserveAliens))
	}

	for alien := range aliens {
		if alien.current == nil {
			t.Error("Alien was not successfully placed in a city")
		}
	}
}

func TestOppositeDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"north", "south"},
		{"south", "north"},
		{"east", "west"},
		{"west", "east"},
		{"", ""},
		{"invalid", ""},
	}

	for _, test := range tests {
		result := oppositeDirection(test.input)
		if result != test.expected {
			t.Errorf("oppositeDirection(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestNew(t *testing.T) {
	// Test creating a new world with 2 aliens
	w := New(2, "test_map.txt")
	if len(w.aliens) != 2 {
		t.Errorf("Expected 2 aliens, but got %d", len(w.aliens))
	}

	if len(w.cities) != 4 {
		t.Errorf("Expected 4 cities, but got %d", len(w.cities))
	}
}

func TestWorld_StartSimulation(t *testing.T) {
	world := New(2, "test_map.txt")
	wg := sync.WaitGroup{}
	wg.Add(3)

	// run simulation in separate goroutine
	go func() {
		world.StartSimulation()
		wg.Done()
	}()

	go func() {
		world.DestroyCitiesAndAliens(&wg)
	}()

	go func() {
		world.DeployReserveAliens()
		wg.Done()
	}()

	// wait for simulation to finish
	wg.Wait()

	// assert that all reserve aliens have been deployed
	if len(world.alienShip) > 0 {
		t.Errorf("Expected all reserve aliens to be deployed, but %d reserve aliens remain", len(world.alienShip))
	}

	// assert that all deployed aliens are dead
	if len(world.deadAliens) != 2 {
		t.Errorf("Expected all deployed aliens to be dead, but %d "+
			"deployed aliens remain alive", len(world.aliens)-len(world.deadAliens))
	}
}
