package algos

// An item which can tick it's internal state
type Tickable interface {
	Tick()                  // Move forward in time
	CopyTickable() Tickable // Create a copy of the object
	String() string         // Convert to a string for comparision
}

// Detects the cycle within the x0 tickable object.
func FloydCycleDetection(x0 Tickable) (cycleLength, cycleStart int) {
	// Copy original state
	tortoise := x0.CopyTickable()
	hare := x0.CopyTickable()

	// Set up the hare / tortoise
	tortoise.Tick()
	hare.Tick()
	hare.Tick()

	// Main phase of the algorithm, finding a repetition
	// the hare moves twice as fast as the tortoise and the distance increases by 1 each step
	// eventually they will both be inside the cycle and then at some point, the distance
	// between them will be divisible by the period λ.
	for tortoise.String() != hare.String() {
		tortoise.Tick()
		hare.Tick()
		hare.Tick()
	}

	// At this point the tortoise position, v, which is also equal to the distance and the tortoise
	// is divisible by period λ. So hare moving in cricle one step at a time, and the tortoise (reset to x0)
	// moving towards the circle, will intersect at the beginning of the circle. Because the distance
	// between them is constant at 2v, a multiple of λ, they will agree as soon as the tortoise reaches
	// index μ.
	tortoise = x0.CopyTickable()
	for tortoise.String() != hare.String() {
		tortoise.Tick()
		hare.Tick()
		cycleStart++
	}

	// Find the length of the shortest cycle starting from x_μ
	// The hare moves one step at a time while the tortoise is still.
	// cycleLength is increases until λ is found.
	cycleLength = 1
	hare = tortoise.CopyTickable()
	hare.Tick()
	tortoiseStr := tortoise.String()
	for tortoiseStr != hare.String() {
		hare.Tick()
		cycleLength++
	}

	return
}
