package main

import "math"

type coord2D struct {
	x float64
	y float64
}

func new2D(x, y int) coord2D {
	return coord2D{
		x: float64(x),
		y: float64(y),
	}
}

func (c coord2D) floor() coord2D {
	return coord2D{
		x: math.Floor(c.x),
		y: math.Floor(c.y),
	}
}

func (c coord2D) add(other coord2D) coord2D {
	return coord2D{
		x: c.x + other.x,
		y: c.y + other.y,
	}
}

func (c coord2D) sub(other coord2D) coord2D {
	return coord2D{
		x: c.x - other.x,
		y: c.y - other.y,
	}
}

func (c coord2D) scale(factor float64) coord2D {
	return coord2D{
		x: c.x * factor,
		y: c.y * factor,
	}
}

func (c coord2D) normalized() coord2D {
	mag := math.Sqrt(c.x * c.x + c.y *c.y)
	return coord2D{
		x: c.x / mag,
		y: c.y / mag,
	}
}

func (c coord2D) dot(other coord2D) float64 {
	return c.x * other.x + c.y * other.y
}