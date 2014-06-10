/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Entity struct {
	Shape            Polygon
	PosX             float64
	PosY             float64
	Angle            float64
	TurnRate         float64
	VelocityX        float64
	VelocityY        float64
	AccelerationRate float64
	MaxVelocity      float64
	rotateLeft       bool
	rotateRight      bool
	accelerate       bool
	decelerate       bool
	isAlive          bool
	createdTime      float64
	lastUpdatedTime  float64
}

func NewEntity(shape Polygon, x float64, y float64, angle float64, turnrate float64, vX float64, vY float64, accel float64, maxvel float64) *Entity {
	return &Entity{
		Shape:            shape,
		PosX:             x,
		PosY:             y,
		Angle:            angle,
		TurnRate:         turnrate,
		VelocityX:        vX,
		VelocityY:        vY,
		AccelerationRate: accel,
		MaxVelocity:      maxvel,
		rotateLeft:       false,
		rotateRight:      false,
		accelerate:       false,
		decelerate:       false,
		isAlive:          true,
		createdTime:      glfw.GetTime(),
		lastUpdatedTime:  glfw.GetTime(),
	}
}

func (ent *Entity) Draw(invertColors bool) {
	if ent.IsAlive() {
		gl.LoadIdentity()
		gl.Begin(gl.POLYGON)

		for v := range ent.Shape.Vectors {
			if invertColors {
				ent.Color3d(ent.Shape.Colors[v])
			} else {
				gl.Color3d(ent.Shape.Colors[v].R, ent.Shape.Colors[v].G, ent.Shape.Colors[v].B)
			}
			ent.GlVertex2d(ent.Shape.Vectors[v])
		}

		gl.End()
	}
}

func (ent *Entity) Color3d(c Color) {
	gl.Color3d(Colorize(c.R), Colorize(c.G), Colorize(c.B))
}

func (ent *Entity) GlVertex2d(v Vector) {
	x, y := v.Rotate(ent.Angle)
	gl.Vertex2d(ent.PosX+x, ent.PosY+y)
}

func (ent *Entity) RotateLeft(flag bool) {
	ent.rotateLeft = flag
}

func (ent *Entity) RotateRight(flag bool) {
	ent.rotateRight = flag
}

func (ent *Entity) Accelerate(flag bool) {
	ent.accelerate = flag
}

func (ent *Entity) Decelerate(flag bool) {
	//ent.decelerate = flag
}

func (ent *Entity) Update() {
	if paused {
		ent.lastUpdatedTime = glfw.GetTime()
	} else {
		timediff := (glfw.GetTime() - ent.lastUpdatedTime) * 500
		ent.lastUpdatedTime = glfw.GetTime()
		var rad float64 = ((ent.Angle) * math.Pi) / 180

		// rotation
		if ent.rotateLeft {
			ent.Angle = ent.Angle - (ent.TurnRate * timediff)
			if ent.Angle < 0 {
				ent.Angle += 360
			}
		} else if ent.rotateRight {
			ent.Angle = ent.Angle + (ent.TurnRate * timediff)
			if ent.Angle > 360 {
				ent.Angle -= 360
			}
		}

		/*
			0°		Sin(0), Cos(1)
			90°		Sin(1), Cos(0)
			180°	Sin(0), Cos(-1)
			270°	Sin(-1), Cos(0)
		*/
		if ent.accelerate {
			ent.VelocityX = ent.VelocityX + (ent.AccelerationRate * math.Sin(rad))
			ent.VelocityY = ent.VelocityY + (ent.AccelerationRate * math.Cos(rad))
		} else if ent.decelerate {
			ent.VelocityX = ent.VelocityX - (ent.AccelerationRate * math.Sin(rad))
			ent.VelocityY = ent.VelocityY - (ent.AccelerationRate * math.Cos(rad))
		}

		// max velocity
		totalVelocity := math.Sqrt(ent.VelocityX*ent.VelocityX + ent.VelocityY*ent.VelocityY)
		if totalVelocity > ent.MaxVelocity {
			ent.VelocityX = ent.VelocityX / totalVelocity
			ent.VelocityY = ent.VelocityY / totalVelocity
			ent.VelocityX = ent.VelocityX * ent.MaxVelocity
			ent.VelocityY = ent.VelocityY * ent.MaxVelocity
		}

		// move
		ent.PosX = ent.VelocityX*timediff + ent.PosX
		ent.PosY = ent.VelocityY*timediff + ent.PosY

		// crude zone clipping
		// TODO: for now it works, but needs to be updated for seamless clipping..
		if ent.PosX > gameWidth {
			ent.PosX -= gameWidth
		} else if ent.PosX < 0 {
			ent.PosX += gameWidth
		}
		if ent.PosY > gameHeight {
			ent.PosY -= gameHeight
		} else if ent.PosY < 0 {
			ent.PosY += gameHeight
		}
	}
}

func (ent *Entity) AddFrictionToVelocity(friction float64) {
	frict := (friction / 100)

	newX := ent.VelocityX * (1 - frict)
	if math.Abs(newX) < frict/10 {
		newX = 0
	}
	ent.VelocityX = newX

	newY := ent.VelocityY * (1 - frict)
	if math.Abs(newY) < frict/10 {
		newY = 0
	}
	ent.VelocityY = newY
}

func (ent *Entity) IsAlive() bool {
	return ent.isAlive
}

func (ent *Entity) Destroy() {
	ent.isAlive = false
}
