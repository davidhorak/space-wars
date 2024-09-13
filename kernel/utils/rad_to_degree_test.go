package utils

import (
	"math"
	"testing"
)

func TestRadToDegree(t *testing.T) {
	tests := []struct {
		radians float64
		degree  float64
	}{
		{0, 0},
		{math.Pi / 2, 90},
		{math.Pi, 180},
		{3 * math.Pi / 2, 270},
		{2 * math.Pi, 360},
	}

	for _, test := range tests {
		got := RadToDegree(test.radians)
		if !AlmostEqual(got, test.degree) {
			t.Errorf("RadToDegree(%f) = %f, expected %f", test.radians, got, test.degree)
		}
	}
}
