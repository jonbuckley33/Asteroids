/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"
	"math/rand"
	"time"

	glfw "github.com/go-gl/glfw3/v3.0/glfw"
)

type Ship struct {
	Entity
	Friction         float64
	shooting         bool
	lastBulletFired  float64
	bulletsPerSecond float64
	mines            int
	torpedos         int
}

func NewShip(x, y, angle, friction float64) *Ship {
	var tip Color
	if PlayerId == 0 {
		tip = Color{1.0, 0.0, 0.0}
	} else if PlayerId == 1 {
		tip = Color{0.0, 1.0, 0.0}
	} else if PlayerId == 2 {
		tip = Color{0.0, 0.0, 1.0}
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		f1 := r.Float64()
		f2 := r.Float64()
		f3 := r.Float64()
		tip = Color{f1, f2, f3}
	}
	shape := Polygon{
		[]Vector{
			Vector{0, 5},
			Vector{4, -5},
			Vector{-4, -5},
		},
		[]Color{
			tip,
			Color{1.0, 1.0, 1.0},
			Color{1.0, 1.0, 1.0},
		},
	}
	return &Ship{*NewEntity(shape, x, y, angle, 0.5, 0, 0, 0.0025, 0.25), friction, false, 0, 5, 3, 1}
}

func (ship *Ship) DropMine() {
	if ship.IsAlive() && ship.mines > 0 {
		x, y := RotateVector(&Vector{0, -10}, ship.Angle)

		mine := NewMine(ship.PosX+x, ship.PosY+y)
		mines = append(mines, mine)

		ship.mines -= 1
	}
}

func (ship *Ship) Shoot(flag bool) {
	ship.shooting = flag
	ship.shoot()
}

func (ship *Ship) shoot() {
	if ship.shooting && glfw.GetTime() > ship.lastBulletFired+(1/ship.bulletsPerSecond) && ship.IsAlive() {
		var rad float64 = ((ship.Angle) * math.Pi) / 180
		x, y := RotateVector(&Vector{0, 5}, ship.Angle)

		bullet := NewBullet(
			ship.PosX+x,
			ship.PosY+y,
			ship.MaxVelocity*math.Sin(rad)*2,
			ship.MaxVelocity*math.Cos(rad)*2,
		)
		bullets = append(bullets, bullet)
		ship.lastBulletFired = glfw.GetTime()
	}
}

func (ship *Ship) ShootTorpedo() {
	if ship.IsAlive() && ship.torpedos > 0 {
		var rad float64 = ((ship.Angle) * math.Pi) / 180
		x, y := RotateVector(&Vector{0, 8}, ship.Angle)

		torpedo := NewTorpedo(
			ship.PosX+x,
			ship.PosY+y,
			ship.Angle,
			ship.MaxVelocity*math.Sin(rad)*1.5,
			ship.MaxVelocity*math.Cos(rad)*1.5,
		)
		torpedos = append(torpedos, torpedo)

		ship.torpedos -= 1
	}
}

func (ship *Ship) Update() {
	ship.shoot()
	ship.Entity.Update()
	if !paused {
		ship.AddFrictionToVelocity(ship.Friction)
	}
}

func (ship *Ship) Destroy() {
	ship.shooting = false
	ship.Entity.Destroy()
	explosions = append(explosions, NewExplosion(ship.PosX, ship.PosY, 5))
}
