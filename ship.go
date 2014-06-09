/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"

	"github.com/go-gl/gl"
)

type Ship struct {
	Entity
	Friction float64
}

func NewShip(x float64, y float64, angle float64, friction float64) *Ship {
	shape := Polygon{
		[]Vector{
			Vector{0, 5},
			Vector{0, -4},
			Vector{4, -5},
			Vector{0, 5},
			Vector{0, -4},
			Vector{-4, -5},
		},
		[]Color{
			Color{0.0, 0.1, 1.0},
			Color{0.0, 0.1, 0.7},
			Color{0.0, 0.1, 0.7},
			Color{0.1, 0.2, 1.0},
			Color{0.1, 0.2, 0.7},
			Color{0.1, 0.2, 0.7},
		},
	}
	return &Ship{*NewEntity(shape, x, y, angle, 0.5, 0, 0, 0.0025, 0.25), friction}
}

func (ship *Ship) Shoot() *Bullet {
	var rad float64 = ((ship.Angle) * math.Pi) / 180
	x, y := RotateVector(&Vector{0, 5}, ship.Angle)
	return NewBullet(
		ship.PosX+x,
		ship.PosY+y,
		ship.MaxVelocity*math.Sin(rad)*2,
		ship.MaxVelocity*math.Cos(rad)*2,
	)
}

func (ship *Ship) Update() {
	ship.Entity.Update()
	ship.AddFrictionToVelocity(ship.Friction)
}

func (ship *Ship) Draw() {
	if ship.IsAlive() {
		gl.LoadIdentity()
		gl.Begin(gl.TRIANGLES)

		for v := range ship.Shape.Vectors {
			gl.Color3d(ship.Shape.Colors[v].R, ship.Shape.Colors[v].G, ship.Shape.Colors[v].B)
			ship.GlVertex2d(ship.Shape.Vectors[v])
		}

		gl.End()
	}
}

func (ship *Ship) Destroy() {
	ship.Entity.Destroy()
	explosions = append(explosions, NewExplosion(ship.PosX, ship.PosY, 5))
}
