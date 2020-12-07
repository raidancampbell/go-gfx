package internal

import "math"

type Coord2D struct {
	X float64
	Y float64
}

func New2D(x, y int) Coord2D {
	return Coord2D{
		X: float64(x),
		Y: float64(y),
	}
}

func (c Coord2D) Floor() Coord2D {
	return Coord2D{
		X: math.Floor(c.X),
		Y: math.Floor(c.Y),
	}
}

func (c Coord2D) Add(other Coord2D) Coord2D {
	return Coord2D{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Coord2D) Sub(other Coord2D) Coord2D {
	return Coord2D{
		X: c.X - other.X,
		Y: c.Y - other.Y,
	}
}

func (c Coord2D) Scale(factor float64) Coord2D {
	return Coord2D{
		X: c.X * factor,
		Y: c.Y * factor,
	}
}

func (c Coord2D) Normalized() Coord2D {
	mag := math.Sqrt(c.X* c.X + c.Y*c.Y)
	return Coord2D{
		X: c.X / mag,
		Y: c.Y / mag,
	}
}

func (c Coord2D) Dot(other Coord2D) float64 {
	return c.X* other.X + c.Y* other.Y
}