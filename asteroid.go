/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import "github.com/go-gl/gl"

type Asteroid struct {
	Entity
	SizeRatio float64
	Lives     int
}

func NewAsteroid(x float64, y float64, angle float64, turnrate float64, vX float64, vY float64, size float64, lives int) *Asteroid {
	shape := Polygon{
		[]Vector{
			Vector{0 * size, 5.0 * size},
			Vector{5.0 * size, 2.0 * size},
			Vector{4.0 * size, -3.0 * size},
			Vector{1.0 * size, -3.0 * size},
			Vector{-1.0 * size, -5.0 * size},
			Vector{-3.0 * size, -5.0 * size},
			Vector{-4.0 * size, -1.0 * size},
			Vector{-6.0 * size, 2.0 * size},
		},
		[]Color{
			Color{0.4, 0.3, 0.1},
			Color{0.9, 0.9, 0.4},
			Color{1, 1, 0.5},
			Color{0.8, 0.8, 0.2},
			Color{0.7, 0.7, 0.2},
			Color{0.5, 0.5, 0.1},
			Color{0.6, 0.6, 0.1},
			Color{0.5, 0.5, 0.1},
		},
	}
	return &Asteroid{*NewEntity(shape, x, y, angle, turnrate, vX, vY, 0, 5), size, lives}
}

func (ast *Asteroid) Draw() {
	if ast.IsAlive() {
		gl.LoadIdentity()
		gl.Begin(gl.POLYGON)

		for v := range ast.Shape.Vectors {
			gl.Color3d(ast.Shape.Colors[v].R, ast.Shape.Colors[v].G, ast.Shape.Colors[v].B)
			ast.GlVertex2d(ast.Shape.Vectors[v])
		}

		gl.End()
	}
}

func (ast *Asteroid) Destroy() {
	ast.Entity.Destroy()
	if ast.Lives > 0 {
		ast.CreateChild()
		ast.CreateChild()
	}
}

func (ast *Asteroid) CreateChild() {
	asteroid := NewAsteroid(ast.PosX, ast.PosY, rng.Float64()*360, rng.Float64()/10, (rng.Float64()-0.5)/4, (rng.Float64()-0.5)/4, ast.SizeRatio/1.5, ast.Lives-1)
	if rng.Float64() > 0.5 {
		asteroid.RotateRight(true)
	} else {
		asteroid.RotateLeft(true)
	}
	asteroids = append(asteroids, asteroid)
}

func CreateAsteroid(size float64, lives int) {
	// avoid creating asteroid too close to ship/player starting position..
	var x float64 = 0
	if rng.Float64() > 0.5 {
		x = gameWidth / 3 * 2
	}
	var y float64 = 0
	if rng.Float64() > 0.5 {
		y = gameHeight / 3 * 2
	}

	asteroid := NewAsteroid((rng.Float64()*gameWidth/3)+x, (rng.Float64()*gameHeight/3)+y, rng.Float64()*360, rng.Float64()/10, (rng.Float64()-0.5)/2, (rng.Float64()-0.5)/2, size, lives)
	if rng.Float64() > 0.5 {
		asteroid.RotateRight(true)
	} else {
		asteroid.RotateLeft(true)
	}
	asteroids = append(asteroids, asteroid)
}
