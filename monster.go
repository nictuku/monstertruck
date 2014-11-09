// Copyright 2014 Yves Junqueira. All rights reserved.
//
// Based on azul3d examples.
// Copyright 2014 The Azul3D Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Displays a TMX tiled map.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
	"azul3d.org/mouse.v1"
	"azul3d.org/tmx.dev"
)

// absPath the absolute path to an file given one relative to the examples
// directory:
//  $GOPATH/src/azul3d.org/examples.dev
var examplesDir string

func absPath(relPath string) string {
	if len(examplesDir) == 0 {
		// Find assets directory.
		for _, path := range filepath.SplitList(build.Default.GOPATH) {
			path = filepath.Join(path, "src/github.com/nictuku/monstertruck")
			if _, err := os.Stat(path); err == nil {
				examplesDir = path
				break
			}
		}
	}
	return filepath.Join(examplesDir, relPath)
}

// setOrthoScale sets the camera's projection matrix to an orthographic one
// using the given viewing rectangle. It performs scaling with the viewing
// rectangle.
func setOrthoScale(c *gfx.Camera, view image.Rectangle, scale, near, far float64) {
	w := float64(view.Dx())
	w *= scale
	w = float64(int((w / 2.0)))

	h := float64(view.Dy())
	h *= scale
	h = float64(int((h / 2.0)))

	m := lmath.Mat4Ortho(-w, w, -h, h, near, far)
	c.Projection = gfx.ConvertMat4(m)
}

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 10000.0
	camZoom := 2.0       // 1x zoom
	camZoomSpeed := 0.01 // 0.01x zoom for each scroll wheel click.
	camMinZoom := 0.1

	// updateCamera simply locks the camera, and calls setOrthoScale with the
	// values above, and then unlocks the camera.
	updateCamera := func() {
		if camZoom < camMinZoom {
			camZoom = camMinZoom
		}
		camera.Lock()
		fmt.Println("camNear", camNear)
		fmt.Println("camFar", camFar)
		setOrthoScale(camera, r.Bounds(), camZoom, camNear, camFar)
		camera.Unlock()
	}

	// Update the camera now.
	updateCamera()

	// Move the camera back two units away from the card.
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Load TMX map file.
	tmxMap, layers, err := tmx.LoadFile(*mapFile, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create an event mask for the events we are interested in.
	evMask := window.FramebufferResizedEvents
	evMask |= window.CursorMovedEvents
	evMask |= window.MouseEvents
	evMask |= window.KeyboardTypedEvents

	// Create a channel of events.
	events := make(chan window.Event, 256)

	// Have the window notify our channel whenever events occur.
	w.Notify(events, evMask)

	handleEvents := func() {
		limit := len(events)
		for i := 0; i < limit; i++ {
			e := <-events
			switch ev := e.(type) {
			case window.FramebufferResized:
				// Update the camera's to account for the new width and height.
				updateCamera()

			case mouse.Event:
				if ev.Button == mouse.Left && ev.State == mouse.Up {
					// Toggle mouse grab.
					props := w.Props()
					props.SetCursorGrabbed(!props.CursorGrabbed())
					w.Request(props)
				}
				if ev.Button == mouse.Wheel && ev.State == mouse.ScrollForward {
					// Zoom in, and update the camera.
					camZoom -= camZoomSpeed
					updateCamera()
				}
				if ev.Button == mouse.Wheel && ev.State == mouse.ScrollBack {
					// Zoom out, and update the camera.
					camZoom += camZoomSpeed
					updateCamera()
				}

			case window.CursorMoved:
				if ev.Delta {
					p := lmath.Vec3{ev.X, 0, -ev.Y}
					camera.Lock()
					camera.SetPos(camera.Pos().Add(p))
					camera.Unlock()
					fmt.Println("newPost", camera.Pos())
				}

			case keyboard.TypedEvent:
				fmt.Println("event", ev)
				switch ev.Rune {
				case 'm':
					// Toggle MSAA now.
					msaa := !r.MSAA()
					r.SetMSAA(msaa)
					fmt.Println("MSAA Enabled?", msaa)
				case 'r':
					camera.Lock()
					camera.SetPos(lmath.Vec3{0, -2, 0})
					camera.Unlock()
				}
			}
		}
	}
	truckBox := newBox(layers["Truck"])
	handleKey := func() {
		var v lmath.Vec3
		// Depending on keyboard state, transform the triangle.
		kb := w.Keyboard()
		if kb.Down(keyboard.Space) {
			v.X += 3
		}
		if kb.Down(keyboard.ArrowDown) {
			v.Z -= 1
		}
		if kb.Down(keyboard.ArrowUp) {
			v.Z += 20
		}
		// Move the truck.
		truckBox.forces = append(truckBox.forces, v)
		v = truckBox.applyPhysics()

		camera.Lock()
		camera.SetPos(camera.Pos().Add(v))
		camera.Unlock()
	}

	// 400 == 2 tiles. Get this from TMX.
	camera.SetPos(lmath.Vec3{windowWidth, -2, 400})

	for {
		// Handle events.
		handleEvents()

		handleKey()

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw the TMX map to the screen.
		for _, layer := range tmxMap.Layers {
			objects, ok := layers[layer.Name]
			if !ok {
				continue
			}
			for _, obj := range objects {
				r.Draw(image.Rect(0, 0, 0, 0), obj, camera)
			}

		}

		// Render the whole frame.
		r.Render()
	}
}

const (
	windowWidth  = 800
	windowHeight = 200
)

func windowProps() *window.Props {
	p := window.NewProps()
	p.SetSize(windowWidth, windowHeight)
	// p.SetPos(0, 0)
	return p
}

var (
	defaultMapFile = absPath("assets/monstermap.tmx")
	mapFile        = flag.String("file", defaultMapFile, "tmx map file to load")
)

func init() {
	flag.Parse()
}

func main() {

	window.Run(gfxLoop, windowProps())
}
