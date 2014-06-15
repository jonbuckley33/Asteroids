/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import glfw "github.com/go-gl/glfw3"

type Torpedo struct {
	Entity
	MaxLifetime float64
}

func NewTorpedo(x, y, angle, vX, vY float64) *Torpedo {
	shape := Polygon{
		[]Vector{
			Vector{0, 1},
			Vector{-1, -3},
			Vector{1, -3},
		},
		[]Color{
			Color{1, 0, 1},
			Color{1, 0, 1},
			Color{1, 0, 1},
		},
	}
	return &Torpedo{*NewEntity(shape, x, y, angle, 0, vX, vY, 0, 5), 1.5}
}

func (torpedo *Torpedo) Update() {
	if paused {
		timediff := (glfw.GetTime() - torpedo.Entity.lastUpdatedTime)
		torpedo.MaxLifetime = torpedo.MaxLifetime + timediff
	}
	torpedo.Entity.Update()
}

func (torpedo *Torpedo) IsAlive() bool {
	if glfw.GetTime() > torpedo.createdTime+torpedo.MaxLifetime {
		torpedo.Destroy()
	}
	return torpedo.Entity.IsAlive()
}

func (torpedo *Torpedo) Destroy() {
	torpedo.Entity.Destroy()
	explosions = append(explosions, NewExplosion(torpedo.PosX, torpedo.PosY, 10))
	bigExplosions = append(bigExplosions, NewBigExplosion(torpedo.PosX, torpedo.PosY, 2))
}
