package physics

import (
	"testing"
)

func TestEdge_ClosestPoint(t *testing.T) {
	var tests = []struct {
		edge     Edge
		point    Vector2
		expected Vector2
	}{
		{Edge{Start: Vector2{X: 0, Y: 0}, End: Vector2{X: 1, Y: 0}}, Vector2{X: 0.5, Y: 0}, Vector2{X: 0.5, Y: 0}},
		{Edge{Start: Vector2{X: 0, Y: 0}, End: Vector2{X: 1, Y: 0}}, Vector2{X: 0.33, Y: 0}, Vector2{X: 0.33, Y: 0}},
		{Edge{Start: Vector2{X: 0, Y: 0}, End: Vector2{X: 1, Y: 0}}, Vector2{X: -10, Y: -10}, Vector2{X: 0, Y: 0}},
		{Edge{Start: Vector2{X: 0, Y: 0}, End: Vector2{X: 1, Y: 0}}, Vector2{X: 10, Y: 10}, Vector2{X: 1, Y: 0}},
	}

	for _, test := range tests {
		result := test.edge.ClosestPoint(test.point)
		if result != test.expected {
			t.Errorf("ClosestPoint(%v) = %v; expected %v", test.point, result, test.expected)
		}
	}
}
