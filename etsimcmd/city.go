package etsimcmd

const (
	cityOccupantMax = 2
)

type City struct {
	name      string
	neighbors map[string]string
	occupants [cityOccupantMax]*Alien
	full      bool
}

func (c *City) invade(alien *Alien) {
	// get the index of the alien to remove
	idx := alien.id % cityOccupantMax
	// remove the alien from the city it currently occupies
	if alien.current != nil {
		alien.current.occupants[idx] = nil
		alien.current.full = false
	}

	// invade the new city
	c.occupants[idx] = alien
	alien.current = c
	c.full = c.occupants[0] != nil && c.occupants[1] != nil

	return
}
