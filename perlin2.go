package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

func Perlin2(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw, offset float64) {
	cfg.Title = "perlin-waves"

	points := make([]pixel.Vec, WINDOW_WIDTH)
	for x := 0.; x < float64(len(points)); x+= 1 {

		//o := offset
		points[int(x)] = pixel.Vec {
			X: x,
			Y: WINDOW_HEIGHT/2 +
				(doPerlin2((x + offset*5)/140, offset * 5/140) * WINDOW_HEIGHT/4) +
				(doPerlin2((x + offset*10)/50, offset * 5/50) * WINDOW_HEIGHT/4)/4,
		}
	}
	imd.Push(points...)
	imd.Line(1)
}

func doPerlin2 (x, y float64) float64 {
	/* Calculate lattice points. */
	x0 := math.Floor(x)
	x1 := x0 + 1.0
	y0 := math.Floor(y)
	y1 := y0 + 1.0

	/* Look up gradients at lattice points. */
	gxx0y0, gyx0y0 := grad2(x0, y0)
	gxx1y0, gyx1y0 := grad2(x1, y0)
	gxx0y1, gyx0y1 := grad2(x0, y1)
	gxx1y1, gyx1y1 := grad2(x1, y1)

	fadeT0 := fade(x - x0) /* Used for interpolation in horizontal direction */

	fadeT1 := fade(y - y0) /* Used for interpolation in vertical direction. */

	/* Calculate dot products and interpolate.*/
	p0p1 := (1.0 - fadeT0) * dot(gxx0y0, gyx0y0, x-x0, y-y0) + fadeT0 * dot(gxx1y0, gyx1y0, x-x1, y-y0)// between upper two lattice points
	p2p3 := (1.0 - fadeT0) * dot(gxx0y1, gyx0y1, x-x0, y-y1) + fadeT0 * dot(gxx1y1, gyx1y1, x-x1, y-y1)// between lower two lattice points

	/* Calculate final result */
	return (1.0 - fadeT1) * p0p1 + fadeT1 * p2p3
}

func dot(x0, y0, x1, y1 float64) float64 {
	return x0 * x1 + y0 * y1
}

// Generate a 2D gradient vector from a given input vector
func grad2(x, y float64) (float64, float64) {
	retX := perlinLUT[int(x) % len(perlinLUT)] * 2 - 1 // map FROM [0,1] TO [-1,1]
	retY := perlinLUT[int(y) % len(perlinLUT)] * 2 - 1

	// normalize
	mag := math.Sqrt(retX * retX + retY * retY)
	return retX / mag, retY / mag
}
