package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/raidancampbell/go-gfx/internal"
	"math"
	"math/rand"
)

const (
	POINT_RATIO = 1./2  // how many points to choose: line will be drawn between points
)

func correlatedNoiseLine(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw) {
	cfg.Title = "correlated-noise-line"

	tmp := float64(internal.WINDOW_WIDTH) * POINT_RATIO // patch for golang's special treatment of type conversion on constants
	points := make([]pixel.Vec, int(tmp))
	for x := 0.; x < float64(len(points)); x+= 1 {
		if x == 0 {
			points[int(x)] = pixel.Vec {
				X: x,
				Y: correlatedNoiseNextPoint(math.NaN(), float64(internal.WINDOW_HEIGHT)),
			}
		} else {
			points[int(x)] = pixel.Vec {
				X: x / POINT_RATIO,
				Y: correlatedNoiseNextPoint(points[int(x)-1].Y, float64(internal.WINDOW_HEIGHT)),
			}
		}
	}
	imd.Push(points...)
	imd.Line(1)
}

func correlatedNoiseNextPoint(prevy, lim float64) float64 {
	if math.IsNaN(prevy) {
		return lim / 2
	}
	coeff := 1.
	if rand.Int() % 2 == 0 {
		coeff *= -1
	}

	// next position is a -1 - 1 walk from the previous position.
	newY := rand.Float64() * coeff + prevy
	// if we're walking out of bounds, invert to keep within bounds
	if newY > lim || newY < 0 {
		coeff *= -1
		newY = rand.Float64() * coeff + prevy
	}
	return newY
}