package utils

import (
	"math"
)

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
func AlmostEqualVector2(a, b struct{ X, Y float64 }) bool {
	return AlmostEqual(a.X, b.X) && AlmostEqual(a.Y, b.Y)
}
