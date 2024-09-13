package physics

type Edge struct {
	Start, End Vector2
}

func (edge *Edge) ClosestPoint(point Vector2) Vector2 {
	offset := ((point.X-edge.Start.X)*(edge.End.X-edge.Start.X) + (point.Y-edge.Start.Y)*(edge.End.Y-edge.Start.Y)) / ((edge.End.X-edge.Start.X)*(edge.End.X-edge.Start.X) + (edge.End.Y-edge.Start.Y)*(edge.End.Y-edge.Start.Y))
	if offset <= 0 {
		return edge.Start
	}
	if offset >= 1 {
		return edge.End
	}

	return Vector2{
		X: edge.Start.X + offset*(edge.End.X-edge.Start.X),
		Y: edge.Start.Y + offset*(edge.End.Y-edge.Start.Y),
	}
}
