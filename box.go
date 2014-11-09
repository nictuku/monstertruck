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
		forces:          []lmath.Vec3{},
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

const (
	floorPos = -200
	gravity  = 9
)

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
	// Gravity.
	v.Z -= gravity

	if v.X > 0 {
		v.X -= b.kineticFriction
		v.X = math.Max(0, v.X)

		v.X = math.Min(v.X, b.maxSpeed)
	}

	fmt.Println("result", v)
	fmt.Println("applying forces")

	// Finding collisions
	for _, obj := range b.layer {
		pos := obj.Pos().Add(v)
		if pos.Z < floorPos { // collision with the floor.
			v.Z = floorPos - obj.Pos().Z
		}
	}
	// resulting vector
	b.forces = []lmath.Vec3{v}

	for _, obj := range b.layer {
		obj.SetPos(obj.Pos().Add(v))
	}
	return v
}
