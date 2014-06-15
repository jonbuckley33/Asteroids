/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import glfw "github.com/go-gl/glfw3"

type BigExplosion struct {
	Entity
	MaxLifetime float64
	Size        float64
}

func NewBigExplosion(x, y, size float64) *BigExplosion {
	shape := Polygon{
		[]Vector{
			Vector{-1 * size, 2 * size},
			Vector{1 * size, 2 * size},
			Vector{2 * size, 1 * size},
			Vector{2 * size, -1 * size},
			Vector{1 * size, -2 * size},
			Vector{-1 * size, -2 * size},
			Vector{-2 * size, -1 * size},
			Vector{-2 * size, 1 * size},
		},
		[]Color{
			Color{1, 0.5, 0},
			Color{1, 0.2, 0},
			Color{1, 0.3, 0},
			Color{1, 0.1, 0},
			Color{1, 0.5, 0},
			Color{1, 0.2, 0},
			Color{1, 0.3, 0},
			Color{1, 0.1, 0},
		},
	}

	explosion := &BigExplosion{*NewEntity(shape, x, y, 0, 0, 0, 0, 0, 0), 1, size}
	return explosion
}

func (explosion *BigExplosion) Update() {
	timediff := (glfw.GetTime() - explosion.Entity.lastUpdatedTime)
	if paused {
		explosion.MaxLifetime = explosion.MaxLifetime + timediff
	} else {
		addSize := 50 * timediff
		for v := range explosion.Entity.Shape.Vectors {
			explosion.Entity.Shape.Vectors[v].X = explosion.Entity.Shape.Vectors[v].X / explosion.Size * (explosion.Size + addSize)
			explosion.Entity.Shape.Vectors[v].Y = explosion.Entity.Shape.Vectors[v].Y / explosion.Size * (explosion.Size + addSize)
		}
		explosion.Size += addSize
	}
	explosion.Entity.Update()
}

func (explosion *BigExplosion) IsAlive() bool {
	if glfw.GetTime() > explosion.createdTime+explosion.MaxLifetime {
		return false
	}
	return explosion.Entity.IsAlive()
}
