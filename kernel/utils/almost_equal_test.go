package utils

import "testing"

func TestAlmostEqual(t *testing.T) {
	tests := []struct {
		a, b float64
		want bool
	}{
		{1.0, 1.0, true},
		{1.0, 1.01, false},
		{1.0, 1.0000000000000001, true},
	}

	for _, test := range tests {
		got := AlmostEqual(test.a, test.b)
		if got != test.want {
			t.Errorf("AlmostEqual(%f, %f) = %t, want %t", test.a, test.b, got, test.want)
		}
	}
}
