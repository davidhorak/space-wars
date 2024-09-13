package physics

import (
	"math"
	"testing"

	"github.com/davidhorak/space-wars/kernel/utils"
)

func TestVector2_Normalize(t *testing.T) {

	tests := []struct {
		vector   Vector2
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 0.6, Y: 0.8}},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}},
		{Vector2{X: -3, Y: -4}, Vector2{X: -0.6, Y: -0.8}},
	}

	for _, test := range tests {
		result := test.vector.Normalize()
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Dot(t *testing.T) {
	tests := []struct {
		vector1  Vector2
		vector2  Vector2
		expected float64
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 5, Y: 12}, 63},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}, 0},
		{Vector2{X: -3, Y: -4}, Vector2{X: -5, Y: -12}, 63},
	}

	for _, test := range tests {
		result := test.vector1.Dot(test.vector2)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Cross(t *testing.T) {
	tests := []struct {
		vector1  Vector2
		vector2  Vector2
		expected float64
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 5, Y: 12}, 16},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}, 0},
		{Vector2{X: -3, Y: -4}, Vector2{X: -5, Y: -12}, 16},
	}

	for _, test := range tests {
		result := test.vector1.Cross(test.vector2)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Magnitude(t *testing.T) {
	tests := []struct {
		vector   Vector2
		expected float64
	}{
		{Vector2{X: 3, Y: 4}, 5},
		{Vector2{X: 0, Y: 0}, 0},
		{Vector2{X: -3, Y: -4}, 5},
	}

	for _, test := range tests {
		result := test.vector.Magnitude()
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Clamp(t *testing.T) {
	tests := []struct {
		vector    Vector2
		maxLength float64
		expected  Vector2
	}{
		{Vector2{X: 3, Y: 4}, 5, Vector2{X: 3, Y: 4}},
		{Vector2{X: 6, Y: 8}, 5, Vector2{X: 3, Y: 4}},
		{Vector2{X: -3, Y: -4}, 5, Vector2{X: -3, Y: -4}},
	}

	for _, test := range tests {
		result := test.vector.Clamp(test.maxLength)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Subtract(t *testing.T) {
	tests := []struct {
		vector1  Vector2
		vector2  Vector2
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 5, Y: 12}, Vector2{X: -2, Y: -8}},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}},
		{Vector2{X: -3, Y: -4}, Vector2{X: -5, Y: -12}, Vector2{X: 2, Y: 8}},
	}

	for _, test := range tests {
		result := test.vector1.Subtract(test.vector2)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Add(t *testing.T) {
	tests := []struct {
		vector1  Vector2
		vector2  Vector2
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 5, Y: 12}, Vector2{X: 8, Y: 16}},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}},
		{Vector2{X: -3, Y: -4}, Vector2{X: -5, Y: -12}, Vector2{X: -8, Y: -16}},
	}

	for _, test := range tests {
		result := test.vector1.Add(test.vector2)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Multiply(t *testing.T) {
	tests := []struct {
		vector   Vector2
		scalar   float64
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, 2, Vector2{X: 6, Y: 8}},
		{Vector2{X: 0, Y: 0}, 0, Vector2{X: 0, Y: 0}},
		{Vector2{X: -3, Y: -4}, -2, Vector2{X: 6, Y: 8}},
	}

	for _, test := range tests {
		result := test.vector.Multiply(test.scalar)
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Distance(t *testing.T) {
	tests := []struct {
		vector1  Vector2
		vector2  Vector2
		expected float64
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 5, Y: 12}, 8.2462112512},
		{Vector2{X: 0, Y: 0}, Vector2{X: 0, Y: 0}, 0},
		{Vector2{X: -3, Y: -4}, Vector2{X: -5, Y: -12}, 8.2462112512},
	}

	for _, test := range tests {
		result := test.vector1.Distance(test.vector2)
		if !utils.AlmostEqual(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_Rotate(t *testing.T) {
	tests := []struct {
		vector   Vector2
		radians  float64
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, math.Pi / 2, Vector2{X: -4, Y: 3}},
		{Vector2{X: 0, Y: 0}, math.Pi / 2, Vector2{X: 0, Y: 0}},
		{Vector2{X: -3, Y: -4}, math.Pi / 2, Vector2{X: 4, Y: -3}},
	}

	for _, test := range tests {
		result := test.vector.Rotate(test.radians)
		if !utils.AlmostEqualVector2(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestVector2_RotateAround(t *testing.T) {
	tests := []struct {
		vector   Vector2
		origin   Vector2
		radians  float64
		expected Vector2
	}{
		{Vector2{X: 3, Y: 4}, Vector2{X: 0, Y: 0}, math.Pi / 2, Vector2{X: -4, Y: 3}},
		{Vector2{X: 3, Y: 4}, Vector2{X: 3, Y: 5}, math.Pi / 2, Vector2{X: 4, Y: 5}},
	}

	for _, test := range tests {
		result := test.vector.RotateAround(test.origin, test.radians)
		if !utils.AlmostEqualVector2(result, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}
