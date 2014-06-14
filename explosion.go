/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Explosion struct {
	Entity
	MaxLifetime float64
	Size        float64
	Lines       []*ExplosionLine
}

type ExplosionLine struct {
	Entity
	Size float64
}

func NewExplosion(x, y, size float64) *Explosion {
	var lines []*ExplosionLine
	for i := 0; i < 7; i++ {
		lines = append(lines, NewExplosionLine(x, y, 0.05, size))
	}

	explosion := &Explosion{*NewEntity(Polygon{}, x, y, 0, 0, 0, 0, 0, 0), 0.3, size, lines}
	return explosion
}

func NewExplosionLine(x, y, velocity, size float64) *ExplosionLine {
	shape := Polygon{
		[]Vector{
			Vector{0, 1},
			Vector{0, 1 + (1.5 * size)},
		},
		[]Color{
			Color{1, 0.6, 0},
			Color{1, 0.5, 0},
		},
	}

	angle := rng.Float64() * 360
	rad := ((angle) * math.Pi) / 180
	vX := velocity * math.Sin(rad)
	vY := velocity * math.Cos(rad)

	return &ExplosionLine{*NewEntity(shape, x, y, angle, 0, vX, vY, 0, 5), size}
}

func (explosion *Explosion) Update() {
	if paused {
		timediff := (glfw.GetTime() - explosion.Entity.lastUpdatedTime)
		explosion.MaxLifetime = explosion.MaxLifetime + timediff
	}
	explosion.Entity.Update()
	for l, _ := range explosion.Lines {
		explosion.Lines[l].Update()
	}
}

func (explosion *Explosion) Draw() {
	if explosion.IsAlive() {
		for l, _ := range explosion.Lines {
			gl.LoadIdentity()
			gl.Begin(gl.LINES)

			for v, _ := range explosion.Lines[l].Shape.Vectors {
				gl.Color3d(explosion.Lines[l].Shape.Colors[v].R, explosion.Lines[l].Shape.Colors[v].G, explosion.Lines[l].Shape.Colors[v].B)
				explosion.Lines[l].GlVertex2d(explosion.Lines[l].Shape.Vectors[v])
			}

			gl.End()
		}
	}
}

func (explosion *Explosion) IsAlive() bool {
	if glfw.GetTime() > explosion.createdTime+explosion.MaxLifetime {
		return false
	}
	return explosion.Entity.IsAlive()
}
