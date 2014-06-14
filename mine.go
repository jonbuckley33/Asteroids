/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

type Mine struct {
	Entity
}

func NewMine(x, y float64) *Mine {
	shape := Polygon{
		[]Vector{
			Vector{-2, 2},
			Vector{2, 2},
			Vector{-2, -2},
			Vector{2, -2},
		},
		[]Color{
			Color{0.5, 1, 0},
			Color{0.5, 1, 0},
			Color{0.5, 1, 0},
			Color{0.5, 1, 0},
		},
	}
	mine := &Mine{*NewEntity(shape, x, y, 0, 0.5, 0, 0, 0, 5)}
	if rng.Float64() > 0.5 {
		mine.RotateRight(true)
	} else {
		mine.RotateLeft(true)
	}
	return mine
}

func (mine *Mine) Destroy() {
	mine.Entity.Destroy()
	explosions = append(explosions, NewExplosion(mine.PosX, mine.PosY, 10))
}
