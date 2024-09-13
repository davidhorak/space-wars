package utils

import (
	"math"
	"testing"
)

func TestDegreeToRad(t *testing.T) {
	tests := []struct {
		degree  float64
		radians float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{270, 3 * math.Pi / 2},
		{360, 2 * math.Pi},
	}

	for _, test := range tests {
		got := DegreeToRad(test.degree)
		if !AlmostEqual(got, test.radians) {
			t.Errorf("DegreeToRad(%f) = %f, expected %f", test.degree, got, test.radians)
		}
	}
}
