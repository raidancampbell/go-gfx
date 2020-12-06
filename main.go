package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)
const (
	WINDOW_WIDTH = 1024
	WINDOW_HEIGHT = 128
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "demo",
		Bounds: pixel.R(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	win.Clear(colornames.Black)

	frames:=0
	second := time.Tick(time.Second)
	for !win.Closed() {
		// safe to remove: vsync is set.  60fps is too fast for viewing though.
		time.Sleep(50 * time.Millisecond)
		win.Clear(colornames.Black)
		imd.Clear()

		correlatedNoiseLine(cfg, imd)

		imd.Draw(win)
		win.Update()
		frames++
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