package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
	"time"
)



func Perlin2DTemporal(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw, offset float64) {
	cfg.Title = "perlin-waves"

	points := make([]pixel.Vec, WINDOW_WIDTH)
	for x := 0.; x < float64(len(points)); x+= 1 {

		//o := offset
		points[int(x)] = pixel.Vec {
			X: x,
			Y: WINDOW_HEIGHT/2 +
				(doPerlin2D(coord2D{(x + offset*5)/140, offset * 5/140}) * WINDOW_HEIGHT/4) +
				(doPerlin2D(coord2D{(x + offset*10)/50, offset * 5/50}) * WINDOW_HEIGHT/4)/4,
		}
	}
	imd.Push(points...)
	imd.Line(1)
}

func Perlin2DSpatial(cfg *pixelgl.WindowConfig, win *pixelgl.Window) {
	cfg.Title = "perlin-2D"
	var px = make([]uint8, len(win.Canvas().Pixels()))
	idx := 0
	for y := 0; y < int(win.Bounds().H()); y++ {
		for x := 0; x < int(win.Bounds().W()); x++ {
			c := new2D(x, y)
			n := doPerlin2D(c.scale(1./64)) * 1.0 +
				doPerlin2D(c.scale(1./32)) * 0.5 +
				doPerlin2D(c.scale(1./16)) * 0.25 +
				doPerlin2D(c.scale(1./8)) * 0.125
			n = (n + 1) / 2 // from [-1,1] to [0,1]
			px[idx] = uint8(255 * n) //R
			idx++
			px[idx] = uint8(255 * n) //G
			idx++
			px[idx] = uint8(255 * n)// B
			idx++
			px[idx] = uint8(255) //A
			idx++
		}
	}

	win.Canvas().SetPixels(px)
}


func doPerlin2D(p coord2D) float64 {
	/* Calculate lattice points. */
	p0 := p.floor()
	p1 := p0.add(new2D(1., 0.))
	p2 := p0.add(new2D(0., 1.))
	p3 := p0.add(new2D(1., 1.))

	/* Look up gradients at lattice points. */
	g0 := grad2D(p0)
	g1 := grad2D(p1)
	g2 := grad2D(p2)
	g3 := grad2D(p3)

	fadeT0 := fade(p.x - p0.x) /* Used for interpolation in horizontal direction */
	fadeT1 := fade(p.y - p0.y)     /* Used for interpolation in vertical direction. */

	/* Calculate dot products and interpolate.*/
	p0p1 := (1.0 - fadeT0) * g0.dot(p.add(p0.scale(-1.))) + fadeT0 * g1.dot(p.add(p1.scale(-1.)))// between upper two lattice points
	p2p3 := (1.0 - fadeT0) * g2.dot(p.add(p2.scale(-1.))) + fadeT0 * g3.dot(p.add(p3.scale(-1.)))// between lower two lattice points

	/* Calculate final result */
	return (1.0 - fadeT1) * p0p1 + fadeT1 * p2p3
}

func dot(x0, y0, x1, y1 float64) float64 {
	return x0 * x1 + y0 * y1
}

func grad2D(c coord2D) coord2D {
	//const float texture_width = 256.0;
	//vec4 v = texture2D(iChannel0, vec2(p.x / texture_width, p.y / texture_width));
	//return normalize(v.xy*2.0 - vec2(1.0)); /* remap sampled value to [-1; 1] and normalize */

	ret := coord2D{
		x: perlinLUT[int(c.x) % len(perlinLUT)] * 2 - 1,
		y: perlinLUT2[int(c.y) % len(perlinLUT2)] * 2 - 1,
	}
	return ret.normalized()
}

// Generate a 2D gradient vector from a given input vector
func grad2(x, y float64) (float64, float64) {
	retX := perlinLUT[int(x) % len(perlinLUT)] * 2 - 1 // map FROM [0,1] TO [-1,1]
	retY := perlinLUT2[int(y) % len(perlinLUT2)] * 2 - 1

	// normalize
	mag := math.Sqrt(retX * retX + retY * retY)
	return retX / mag, retY / mag
}

var perlinLUT2 = make([]float64, 512)

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := range perlinLUT2 {
		perlinLUT2[i] = rand.Float64()
	}
}