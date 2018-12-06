package lib

// ABS taken from Hackers Delight
func Abs(a int) int {
	b := a >> 31

	return (a ^ b) - b
}

// Min of two numbers
func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// Max of two numbers
func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
