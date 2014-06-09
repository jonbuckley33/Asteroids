/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Bullet struct {
	Entity
	MaxLifetime float64
}

func NewBullet(x float64, y float64, vX float64, vY float64) *Bullet {
	shape := Polygon{
		[]Vector{
			Vector{0, 1},
			Vector{-1, -1},
			Vector{1, -1},
		},
		[]Color{
			Color{1, 0, 0},
			Color{1, 0, 0},
			Color{1, 0, 0},
		},
	}
	bullet := &Bullet{*NewEntity(shape, x, y, 0, 2, vX, vY, 0, 5), 1.8}
	if rng.Float64() > 0.5 {
		bullet.RotateRight(true)
	} else {
		bullet.RotateLeft(true)
	}
	return bullet
}

func (bullet *Bullet) Draw() {
	if bullet.IsAlive() {
		gl.LoadIdentity()
		gl.Begin(gl.TRIANGLES)

		for v := range bullet.Shape.Vectors {
			gl.Color3d(bullet.Shape.Colors[v].R, bullet.Shape.Colors[v].G, bullet.Shape.Colors[v].B)
			bullet.GlVertex2d(bullet.Shape.Vectors[v])
		}

		gl.End()
	}
}

func (bullet *Bullet) IsAlive() bool {
	if glfw.GetTime() > bullet.createdTime+bullet.MaxLifetime {
		return false
	}
	return bullet.Entity.IsAlive()
}
