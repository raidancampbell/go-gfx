package bouncyball

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/raidancampbell/go-gfx/internal"
	"github.com/raidancampbell/go-gfx/perlin"
	"github.com/raidancampbell/libraidan/pkg/roper"
)

type state struct {
	bLocation pixel.Vec
	bVelocity pixel.Vec
}

var s state

func applyGravity() {
	s.bVelocity = s.bVelocity.Add(pixel.Vec{X: 0., Y:-0.1})
}

func getBallRadius() float64 {
	return 0.05 * float64(internal.WINDOW_HEIGHT)
}

func drawBall(imd *imdraw.IMDraw) {
	imd.Push(s.bLocation)

	imd.Circle(getBallRadius(), 0)
}

func checkCollision(ground []pixel.Vec) bool {
	// collide with left/right walls
	if s.bLocation.X - getBallRadius() < 0 {
		s.bVelocity.X *= -1
		s.bLocation.X = 0 + getBallRadius()
		return true
	} else if s.bLocation.X + getBallRadius() > float64(internal.WINDOW_WIDTH) {
		s.bVelocity.X *= -1
		s.bLocation.X = float64(internal.WINDOW_WIDTH) - getBallRadius()
		return true
	}
	// ceil
	if s.bLocation.Y + getBallRadius() > float64(internal.WINDOW_HEIGHT) {
		s.bVelocity.Y *= -1
		s.bLocation.Y = float64(internal.WINDOW_HEIGHT) - getBallRadius()
		return true
	}
	// irregular ground
	// only think about the ground underneath the ball
	// 2-pass collision detection. first pass to find how much of the arc intersects
	// second pass to calculate collision from middle of arc
	var xIntersect []float64
	for x := s.bLocation.X - getBallRadius(); x <= s.bLocation.X + getBallRadius(); x++ {
		circ := pixel.C(s.bLocation, getBallRadius())
		if circ.Contains(ground[int(x)%len(ground)]) {
			xIntersect = append(xIntersect, x)
		}
	}
	if xIntersect != nil {
		x := xIntersect[len(xIntersect)/2]
		circ := pixel.C(s.bLocation, getBallRadius())
		if circ.Contains(ground[int(x)%len(ground)]) {
			// the velocity is reversed in the direction of the angle of impact
			angle := circ.Center.To(ground[int(x)])
			angle = angle.Scaled(-1.)
			s.bVelocity = angle.Unit().Scaled(s.bVelocity.Len())

			// back the circle out until it no longer intersects
			for ;circ.Contains(ground[int(x) % len(ground)]);circ = circ.Moved(s.bVelocity.Unit()) {}
			s.bLocation = circ.Center
			return true
		}
	}

	// placeholder to check math
	if s.bLocation.Y - getBallRadius() < 0. {
		s.bVelocity.Y *= -1
		s.bLocation.Y = 0 + getBallRadius()
		return true
	}
	return false
}

func Do(cfg *pixelgl.WindowConfig, imd *imdraw.IMDraw, offset float64) {
	// always ensure the ball is drawn,
	defer drawBall(imd)

	//draw an interesting surface
	points := perlin.Perlin2DTemporal(cfg, imd, offset)

	// initialize state
	if roper.IsDefaultValue(s) {
		s.bLocation = pixel.Vec{X: float64(internal.WINDOW_WIDTH)/2, Y:float64(internal.WINDOW_HEIGHT) - getBallRadius()}
		applyGravity()
		s.bVelocity.X = 10
		return
	}
	applyGravity()
	checkCollision(points)
	s.bLocation = s.bLocation.Add(s.bVelocity)
}