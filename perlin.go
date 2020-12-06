package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
	"time"
)

// n(p) = (1 - F(p-p0))g(p0)(p-p0) + F(p-p0)g(p1)(p-p1)

// implementation from https://gpfault.net/posts/perlin-noise.txt.html

var perlinLUT = make([]float64, 512)

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := range perlinLUT {
		perlinLUT[i] = rand.Float64()
	}
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6.0-15.0) + 10.0)
}

func grad(p float64) float64 {
	v := perlinLUT[int(p) % len(perlinLUT)]
	if v > 0.5 {
		return 1.
	}
	return -1.
}

func doPerlin(p float64) float64 {
	p0 := math.Floor(p)
	p1 := p0 + 1.0 // always 1?

	t := p - p0
	fadeT := fade(t)

	g0 := grad(p0)
	g1 := grad(p1)

	return (1.0-fadeT)*g0*(p-p0) + fadeT*g1*(p-p1)
}


func Perlin(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw) {
	cfg.Title = "perlin-line"
	frequency := 1.0 / 50.0

	points := make([]pixel.Vec, WINDOW_WIDTH)
	for x := 0.; x < float64(len(points)); x+= 1 {
		points[int(x)] = pixel.Vec {
			X: x,
			Y: WINDOW_HEIGHT/2 + (doPerlin(x * frequency) * float64(WINDOW_HEIGHT)),
		}
	}
	imd.Push(points...)
	imd.Line(1)
}