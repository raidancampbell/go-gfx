package perlin

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/raidancampbell/go-gfx/internal"
	"math/rand"
)

func Perlin2DTemporal(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw, offset float64) []pixel.Vec {
	cfg.Title = "perlin-waves"

	points := make([]pixel.Vec, internal.WINDOW_WIDTH)
	for x := 0.; x < float64(len(points)); x+= 1 {

		points[int(x)] = pixel.Vec {
			X: x,
			Y: float64(internal.WINDOW_HEIGHT)/2 +
				(doPerlin2D(internal.Coord2D{X: (x + offset*5)/140, Y: offset * 5/140}) * float64(internal.WINDOW_HEIGHT) /4) +
				(doPerlin2D(internal.Coord2D{X: (x + offset*10)/50, Y: offset * 5/50}) * float64(internal.WINDOW_HEIGHT)/4)/4,
		}
	}
	imd.Push(points...)
	imd.Line(1)
	return points
}

func Perlin2DSpatial(cfg *pixelgl.WindowConfig, win *pixelgl.Window) {
	cfg.Title = "perlin-2D"
	var px = make([]uint8, len(win.Canvas().Pixels()))
	idx := 0
	baseNoise := GenerateWhiteNoise(internal.WINDOW_WIDTH, internal.WINDOW_HEIGHT)
	perlin := GeneratePerlinNoise(baseNoise, 7)
	for y := 0; y < int(win.Bounds().H()); y++ {
		for x := 0; x < int(win.Bounds().W()); x++ {
			n := perlin[x][y]
			px[idx] = uint8(255 * n) //R
			idx++
			px[idx] = uint8(255 * n) //G
			idx++
			px[idx] = uint8(255 * n) // B
			idx++
			px[idx] = uint8(255) //A
			idx++
		}
	}

	win.Canvas().SetPixels(px)
}

func GenerateSmoothNoise(baseNoise [][]float64, octave int) [][]float64 {
	width := len(baseNoise)
	height := len(baseNoise[0])

	smoothNoise := make([][]float64, width)

	for i := range smoothNoise {
		smoothNoise[i] = make([]float64, height)
	}

	samplePeriod := 1 << octave // calculates 2 ^ k
	sampleFrequency := 1.0 / float64(samplePeriod)

	for i := 0; i < width; i++ {
		//calculate the horizontal sampling indices
		sampleI0 := (i / samplePeriod) * samplePeriod
		sampleI1 := (sampleI0 + samplePeriod) % width //wrap around
		horizontalBlend := float64(i - sampleI0) * sampleFrequency

		for j := 0; j < height; j++ {
			//calculate the vertical sampling indices
			sampleJ0 := (j / samplePeriod) * samplePeriod
			sampleJ1 := (sampleJ0 + samplePeriod) % height //wrap around
			verticalBlend := float64(j - sampleJ0) * sampleFrequency

			//blend the top two corners
			top := Interpolate(baseNoise[sampleI0][sampleJ0], baseNoise[sampleI1][sampleJ0], float64(horizontalBlend))

			//blend the bottom two corners
			bottom := Interpolate(baseNoise[sampleI0][sampleJ1], baseNoise[sampleI1][sampleJ1], float64(horizontalBlend))

			//final blend
			smoothNoise[i][j] = Interpolate(top, bottom, float64(verticalBlend))
		}
	}

	return smoothNoise
}

func Interpolate(x0, x1, alpha float64) float64 {
	return x0*(1-alpha) + alpha*x1
}

func GeneratePerlinNoise(baseNoise [][]float64, octaveCount int) [][]float64 {
	width := len(baseNoise)
	height := len(baseNoise[0])

	smoothNoise := make([][][]float64, octaveCount)

	persistance := 0.8

	//generate smooth noise
	for i := 0; i < octaveCount; i++ {
		smoothNoise[i] = GenerateSmoothNoise(baseNoise, i)
	}

	perlinNoise := make([][]float64, width)

	for i := range perlinNoise {
		perlinNoise[i] = make([]float64, height)
	}
	amplitude := 1.0
	totalAmplitude := 0.0

	//blend noise together
	for octave := octaveCount - 1; octave >= 0; octave-- {
		amplitude *= persistance
		totalAmplitude += amplitude

		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				perlinNoise[i][j] += smoothNoise[octave][i][j] * amplitude
			}
		}
	}

	//normalisation
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			perlinNoise[i][j] /= totalAmplitude
		}
	}

	return perlinNoise
}

func GenerateWhiteNoise(width, height int) [][]float64 {
	noise := make([][]float64, width)

	for i := range noise {
		noise[i] = make([]float64, height)
	}
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			noise[i][j] = rand.Float64()
		}
	}

	return noise
}

func doPerlin2D(p internal.Coord2D) float64 {
	/* Calculate lattice points. */
	p0 := p.Floor()
	p1 := p0.Add(internal.New2D(1., 0.))
	p2 := p0.Add(internal.New2D(0., 1.))
	p3 := p0.Add(internal.New2D(1., 1.))

	/* Look up gradients at lattice points. */
	g0 := grad2D(p0)
	g1 := grad2D(p1)
	g2 := grad2D(p2)
	g3 := grad2D(p3)

	fadeT0 := fade(p.X - p0.X) /* Used for interpolation in horizontal direction */
	fadeT1 := fade(p.Y - p0.Y) /* Used for interpolation in vertical direction. */

	/* Calculate dot products and interpolate.*/
	p0p1 := (1.0-fadeT0)*g0.Dot(p.Sub(p0)) + fadeT0*g1.Dot(p.Sub(p1)) // between upper two lattice points
	p2p3 := (1.0-fadeT0)*g2.Dot(p.Sub(p2)) + fadeT0*g3.Dot(p.Sub(p3)) // between lower two lattice points

	/* Calculate final result */
	return (1.0-fadeT1)*p0p1 + fadeT1*p2p3
}


func grad2D(c internal.Coord2D) internal.Coord2D {

	ret := internal.Coord2D{
		X: perlinLUT[int(c.X)%len(perlinLUT)],
		Y: perlinLUT[int(c.Y)%len(perlinLUT)],
	}.
		Scale(2.).
		Sub(internal.New2D(1, 1))
	return ret.Normalized()
}
