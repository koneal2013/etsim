package etsimcmd

type City struct {
	name      string
	neighbors map[string]string
	occupants [2]*Alien
	full      bool
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
