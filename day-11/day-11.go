package main

import (
	"fmt"
	"github.com/DomBlack/advent-of-code-2018/lib/vectors"
)

func main() {
	part1P, _ := FindBestFuelCellPatch(1133, 3, 3)
	fmt.Printf("Part 1: %d,%d\n", part1P.X, part1P.Y)

	part2P, size := FindBestFuelCellPatch(1133, 1, 300)
	fmt.Printf("Part 2: %d,%d,%d\n", part2P.X, part2P.Y, size)
}

// The power level in a given fuel cell is based on it's position and the grid serial number
func GetPowerLevel(position vectors.Vec2, gridSerialNum int) int {
	rackID := position.X + 10
	powerLevel := position.Y * rackID
	powerLevel += gridSerialNum
	powerLevel *= rackID

	hundredsDigit := (powerLevel / 100) % 10

	return hundredsDigit - 5
}

func FindBestFuelCellPatch(gridSerialNum, minPatchSize, maxPatchSize int) (vectors.Vec2, int) {
	grid := make([]int, 300*300)

	// Build the grid of known values
	for x := 0; x < 300; x++ {
		for y := 0; y < 300; y++ {
			index := y*300 + x

			grid[index] = GetPowerLevel(vectors.NewVec2(x, y), gridSerialNum)
		}
	}

	highestPower := 0
	position := vectors.NewVec2(0, 0)
	size := 0

	// Loop over the grid and build it for the min patch size
	powerGrid := make([]int, 300*300)
	for x := 0; x <= 300-minPatchSize; x++ {
		for y := 0; y <= 300-minPatchSize; y++ {
			baseIndex := y*300 + x
			patchPower := 0

			for xA := 0; xA < minPatchSize; xA++ {
				for yA := 0; yA < minPatchSize; yA++ {
					patchPower += grid[baseIndex+xA+(yA*300)]
				}
			}

			powerGrid[baseIndex] = patchPower

			if patchPower > highestPower {
				highestPower = patchPower
				position = vectors.NewVec2(x, y)
				size = minPatchSize
			}
		}
	}

	// Now loop over the adjustable patch size to account for it, building upon the previous loops result
	for patchSize := minPatchSize + 1; patchSize <= maxPatchSize; patchSize++ {
		for x := 0; x <= 300-patchSize; x++ {
			for y := 0; y <= 300-patchSize; y++ {
				baseIndex := y*300 + x
				patchPower := powerGrid[baseIndex] // Read the existing power for the last size grid

				patchOffset := patchSize - 1
				patchPower -= grid[baseIndex + patchOffset + patchOffset * 300] // We are going to be adding this twice
				for i := 0; i < patchSize; i++ {
					patchPower += grid[baseIndex + patchOffset + (i * 300)] // Add the right edge of this patch
					patchPower += grid[baseIndex + (patchOffset * 300) + i] // Add the bottom edge of this patch
				}

				// Save the grid size
				powerGrid[baseIndex] = patchPower
				if patchPower > highestPower {
					highestPower = patchPower
					position = vectors.NewVec2(x, y)
					size = patchSize
				}
			}
		}
	}

	return position, size
}
