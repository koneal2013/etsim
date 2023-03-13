package etsimcmd

import "testing"

func TestCityInvade(t *testing.T) {
	city := &City{
		name: "testCity",
		neighbors: map[string]string{
			"north": "city1",
			"south": "city2",
		},
		full: false,
	}
	alien1 := &Alien{id: 1}
	alien2 := &Alien{id: 2}

	// Test invading an empty city
	city.invade(alien1)
	if city.occupants[1] != alien1 {
		t.Errorf("Expected alien1 to occupy the city, but got %v", city.occupants[0])
	}
	if alien1.current != city {
		t.Errorf("Expected alien1 to be in testCity, but it's in %v", alien1.current)
	}

	// Test invading a city with one alien already present
	city.invade(alien2)
	if city.occupants[0] == nil || city.occupants[1] == nil {
		t.Errorf("Expected city to be full, but it's not")
	}
}
