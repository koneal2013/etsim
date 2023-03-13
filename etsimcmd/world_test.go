package etsimcmd

import (
	"testing"
)

func TestNewWorld(t *testing.T) {
	numAliens := 10
	worldMapPath := "test_map.txt"
	world := New(numAliens, worldMapPath)

	if len(world.aliens)+len(world.alienShip) != numAliens {
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
	numAliens := 10
	cities := map[string]*City{
		"City 1": {name: "City 1", occupants: [2]*Alien{nil, nil}},
		"City 2": {name: "City 2", occupants: [2]*Alien{nil, nil}},
	}

	aliens, reserveAliens := createAliens(numAliens, cities)

	if (len(aliens)+len(reserveAliens) != numAliens) && len(aliens) > 0 {
		t.Errorf("Expected %d aliens, but got %d", numAliens, len(aliens)+len(reserveAliens))
	}

	for alien := range aliens {
		if alien.current == nil {
			t.Error("Alien was not successfully placed in a city")
		}
	}
}

func TestGetDirectionToNeighbor(t *testing.T) {
	current := &City{
		name: "Test City",
		neighbors: map[string]string{
			"north": "Neighbor City",
			"south": "Another City",
		},
	}

	direction := getDirectionToNeighbor(current, "Neighbor City")
	if direction != "north" {
		t.Errorf("Expected direction 'north', but got '%s'", direction)
	}

	direction = getDirectionToNeighbor(current, "Nonexistent City")
	if direction != "" {
		t.Errorf("Expected empty string, but got '%s'", direction)
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
	// Test simulation with 2 aliens and a small map
	w := New(2, "test_map.txt")
	w.StartSimulation()
	// Verify that all aliens are in a city
	for alien := range w.aliens {
		if alien.current == nil {
			t.Error("Alien is not in a city")
		}
	}
	// Verify that there are no cities with more than one occupant
	for _, city := range w.cities {
		if city.occupants[0] != nil && city.occupants[1] != nil {
			t.Errorf("City %s has more than one occupant", city.name)
		}
	}
}
