package main

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
	"fmt"
)

func newBox(layer map[string]*gfx.Object) *box {
	return &box{
		layer:           layer,
		forces:          map[string]lmath.Vec3{"gravity": {Z: -9}},
		surfaceFriction: 1,
	}
}

type box struct {
	layer  map[string]*gfx.Object
	forces map[string]lmath.Vec3
	// used for one dimensional friction.
	surfaceFriction float64
}

const floor = -200

// applyPhysics to the box and returns its final movement vector.
func (b *box) applyPhysics() lmath.Vec3 {
	var v lmath.Vec3
	for _, f := range b.forces {
		v = v.Add(f)
	}
	if v.X <= b.surfaceFriction {
		v.X = 0
	}
	fmt.Println("result", v)
	fmt.Println("applying forces")

	// Finding collisions
	for _, obj := range b.layer {
		pos := obj.Pos().Add(v)
		if pos.Z < floor { // collision with the floor.
			v.Z = floor - obj.Pos().Z
		}
	}
	for _, obj := range b.layer {
		obj.SetPos(obj.Pos().Add(v))
	}
	return v
}
