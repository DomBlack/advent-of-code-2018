package xcom

// A cell of the map
type Cell struct {
	IsWall bool
	Unit *Unit // The unit on this cell, or null
}

// Is this cell a floor and has not got a unit on it
func (c Cell) IsEmpty() bool {
	return !c.IsWall && c.Unit == nil
}

// Convert the cell to a string
func (c Cell) String() string {
	if c.Unit != nil {
		if c.Unit.IsElf {
			return "E"
		} else {
			return "G"
		}
	} else if c.IsWall {
		return "#"
	} else {
		return "."
	}
}
