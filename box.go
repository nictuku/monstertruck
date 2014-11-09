package main

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
	"fmt"
	"math"
)

func newBox(layer map[string]*gfx.Object) *box {
	return &box{
		layer:           layer,
		forces:          []lmath.Vec3{{Z: -9}},
		surfaceFriction: 2,
		kineticFriction: 0.5,
		maxSpeed:        15,
	}
}

type box struct {
	layer  map[string]*gfx.Object
	forces []lmath.Vec3
	// used for one dimensional friction.
	surfaceFriction float64
	kineticFriction float64
	maxSpeed        float64 // or add wind resistance.
}

const floor = -200

// applyPhysics to the box and returns its final movement vector.
// assumes that velocity on X is always non-negative (can't move to the left).
func (b *box) applyPhysics() lmath.Vec3 {
	var v lmath.Vec3
	for _, f := range b.forces {
		v = v.Add(f)
	}
	if v.X <= b.surfaceFriction {
		v.X = 0
	}
	if v.X > 0 {
		v.X -= b.kineticFriction
		v.X = math.Max(0, v.X)

		v.X = math.Min(v.X, b.maxSpeed)
	}
	// resulting vector
	b.forces = []lmath.Vec3{v}

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
