package physics

import "math"

type Vector2 struct {
	X float64
	Y float64
}

// Normalize returns a unit vector in the same direction as the original vector.
// If the original vector is a zero vector, it returns a zero vector.
func (vector *Vector2) Normalize() Vector2 {
	magnitude := math.Sqrt(vector.X*vector.X + vector.Y*vector.Y)
	if magnitude == 0 {
		return Vector2{0, 0}
	}
	return Vector2{vector.X / magnitude, vector.Y / magnitude}
}

// Dot returns the dot product of the vector and another vector.
// It is a measure of how much one vector goes in the direction of another.
func (vector *Vector2) Dot(other Vector2) float64 {
	return vector.X*other.X + vector.Y*other.Y
}

// Cross returns the cross product of the vector and another vector.
// The cross product of two vectors is a vector that is perpendicular to both of them.
func (vector *Vector2) Cross(other Vector2) float64 {
	return vector.X*other.Y - vector.Y*other.X
}

// Magnitude returns the magnitude of the vector.
func (vector *Vector2) Magnitude() float64 {
	return math.Sqrt(vector.X*vector.X + vector.Y*vector.Y)
}

// Clamp limits the magnitude of the vector to the given maximum length.
func (vector *Vector2) Clamp(maxLength float64) Vector2 {
	magnitude := vector.Magnitude()
	if magnitude > maxLength {
		return vector.Multiply(maxLength / magnitude)
	}
	return *vector
}

func (vector *Vector2) Subtract(other Vector2) Vector2 {
	return Vector2{vector.X - other.X, vector.Y - other.Y}
}

func (vector *Vector2) Add(other Vector2) Vector2 {
	return Vector2{vector.X + other.X, vector.Y + other.Y}
}

func (vector *Vector2) Multiply(scalar float64) Vector2 {
	return Vector2{vector.X * scalar, vector.Y * scalar}
}

func (vector *Vector2) Distance(other Vector2) float64 {
	return math.Sqrt(math.Pow(vector.X-other.X, 2) + math.Pow(vector.Y-other.Y, 2))
}

func (vector *Vector2) Rotate(radians float64) Vector2 {
	cos := math.Cos(radians)
	sin := math.Sin(radians)

	return Vector2{
		X: vector.X*cos - vector.Y*sin,
		Y: vector.X*sin + vector.Y*cos,
	}
}

func (vector *Vector2) RotateAround(origin Vector2, radians float64) Vector2 {
	cos := math.Cos(radians)
	sin := math.Sin(radians)

	x := vector.X - origin.X
	y := vector.Y - origin.Y

	return Vector2{
		X: x*cos - y*sin + origin.X,
		Y: x*sin + y*cos + origin.Y,
	}
}
