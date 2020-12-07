package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/raidancampbell/go-gfx/bouncyball"
	"github.com/raidancampbell/go-gfx/internal"
	"golang.org/x/image/colornames"
	"time"
)
func init() {
	// elevating into main package for readability
	internal.WINDOW_WIDTH = 1024
	internal.WINDOW_HEIGHT = 512
}

func run() {
	cfg := &pixelgl.WindowConfig{
		Title:  "demo",
		Bounds: pixel.R(0, 0, float64(internal.WINDOW_WIDTH), float64(internal.WINDOW_HEIGHT)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(*cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	win.Clear(colornames.Black)

	frames:=0
	second := time.Tick(time.Second)
	frameNum := 0.
	for !win.Closed() {
		// safe to remove: vsync is set.  60fps is too fast for viewing though.
		//time.Sleep(50 * time.Millisecond)
		win.Clear(colornames.Black)
		imd.Clear()

		//correlatedNoiseLine(cfg, imd)
		//perlin.Perlin2DTemporal(cfg, imd, frameNum)
		//perlin.Perlin2DSpatial(cfg, win)
		bouncyball.Do(cfg, imd, frameNum)

		imd.Draw(win)
		win.Update()
		frames++
		frameNum++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}