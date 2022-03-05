package bouncyball

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/raidancampbell/go-gfx/internal"
	"github.com/raidancampbell/go-gfx/perlin"
	"github.com/raidancampbell/libraidan/pkg/roper"
	"math/rand"
)

type ballState struct {
	bLocation pixel.Vec
	bVelocity pixel.Vec
}

var s []ballState

func applyGravity() {
	for i := range s {
		s[i].bVelocity = s[i].bVelocity.Add(pixel.Vec{X: 0., Y:-0.1})
	}
}

func getBallRadius() float64 {
	return 0.03 * float64(internal.WINDOW_HEIGHT)
}

func drawBalls(imd *imdraw.IMDraw) {
	for i := range s {
		imd.Push(s[i].bLocation)
	}
	imd.Circle(getBallRadius(), 0)
}

func checkCollision(ground []pixel.Vec, i int) bool {
	// collide with left/right walls
	if s[i].bLocation.X - getBallRadius() < 0 {
		s[i].bVelocity.X *= -1
		s[i].bLocation.X = 0 + getBallRadius()
		return true
	} else if s[i].bLocation.X + getBallRadius() > float64(internal.WINDOW_WIDTH) {
		s[i].bVelocity.X *= -1
		s[i].bLocation.X = float64(internal.WINDOW_WIDTH) - getBallRadius()
		return true
	}
	// ceil
	if s[i].bLocation.Y + getBallRadius() > float64(internal.WINDOW_HEIGHT) {
		s[i].bVelocity.Y *= -1
		s[i].bLocation.Y = float64(internal.WINDOW_HEIGHT) - getBallRadius()
		return true
	}
	// irregular ground
	// only think about the ground underneath the ball
	// 2-pass collision detection. first pass to find how much of the arc intersects
	// second pass to calculate collision from middle of arc
	var xIntersect []float64
	for x := s[i].bLocation.X - getBallRadius(); x <= s[i].bLocation.X + getBallRadius(); x++ {
		circ := pixel.C(s[i].bLocation, getBallRadius())
		if circ.Contains(ground[int(x)%len(ground)]) {
			xIntersect = append(xIntersect, x)
		}
	}
	if xIntersect != nil {
		x := xIntersect[len(xIntersect)/2]
		circ := pixel.C(s[i].bLocation, getBallRadius())
		if circ.Contains(ground[int(x)%len(ground)]) {
			// the velocity is reversed in the direction of the angle of impact
			angle := circ.Center.To(ground[int(x)])
			angle = angle.Scaled(-1.)
			s[i].bVelocity = angle.Unit().Scaled(s[i].bVelocity.Len())

			// back the circle out until it no longer intersects
			for ;circ.Contains(ground[int(x) % len(ground)]);circ = circ.Moved(s[i].bVelocity.Unit()) {}
			s[i].bLocation = circ.Center
			return true
		}
	}

	// placeholder to check math
	if s[i].bLocation.Y - getBallRadius() < 0. {
		s[i].bVelocity.Y *= -1
		s[i].bLocation.Y = 0 + getBallRadius()
		return true
	}
	return false
}

func Do(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw, offset float64) {
	// always ensure the ball is drawn,
	defer drawBalls(imd)

	//draw an interesting surface
	points := perlin.Perlin2DTemporal(cfg, imd, offset)

	// initialize ballState
	if roper.IsDefaultValue(s) {
		s = make([]ballState, 30)
		for i := range s {
			s[i] = ballState{}
			s[i].bLocation = pixel.Vec{X: float64(internal.WINDOW_WIDTH)/2, Y:float64(internal.WINDOW_HEIGHT) - getBallRadius()}
			s[i].bVelocity.X = rand.Float64() * 5. - 5.
		}
		applyGravity()
		return
	}
	applyGravity()
	for i := range s {
		checkCollision(points, i)
		s[i].bLocation = s[i].bLocation.Add(s[i].bVelocity)
	}
}