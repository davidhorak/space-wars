package utils

import "testing"

func TestAlmostEqual(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected bool
	}{
		{1.0, 1.0, true},
		{1.0, 1.01, false},
		{1.0, 1.0000000000000001, true},
	}

	for _, test := range tests {
		got := AlmostEqual(test.a, test.b)
		if got != test.expected {
			t.Errorf("AlmostEqual(%f, %f) = %t, expected %t", test.a, test.b, got, test.expected)
		}
	}
}

func TestAlmostEqualVector2(t *testing.T) {
	tests := []struct {
		a, b     struct{ X, Y float64 }
		expected bool
	}{
		{struct{ X, Y float64 }{1.0, 1.0}, struct{ X, Y float64 }{1.0, 1.0}, true},
		{struct{ X, Y float64 }{1.0, 1.0}, struct{ X, Y float64 }{1.0, 1.01}, false},
		{struct{ X, Y float64 }{1.0, 1.0}, struct{ X, Y float64 }{1.0000000000000001, 1.0}, true},
	}

	for _, test := range tests {
		got := AlmostEqualVector2(test.a, test.b)
		if got != test.expected {
			t.Errorf("AlmostEqualVector2(%v, %v) = %t, expected %t", test.a, test.b, got, test.expected)
		}
	}
}
