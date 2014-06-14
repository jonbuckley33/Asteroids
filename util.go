/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"

	"github.com/go-gl/gl"
	"github.com/paulsmith/gogeos/geos"
)

type Polygon struct {
	Vectors []Vector
	Colors  []Color
}

type Vector struct {
	X float64
	Y float64
}

type Color struct {
	R float64
	G float64
	B float64
}

func Colorize(c float64) float64 {
	if colorsInverted {
		c = math.Abs(c - 1)
	}
	return c
}

func RotateVector(v *Vector, angle float64) (float64, float64) {
	return v.Rotate(angle)
}

func (v *Vector) Rotate(angle float64) (float64, float64) {
	var rad float64 = ((angle + 90) * math.Pi) / 180
	x := (v.X * math.Sin(rad)) - (v.Y * math.Cos(rad))
	y := (v.X * math.Cos(rad)) + (v.Y * math.Sin(rad))
	return x, y
}

func IsColliding(a *Entity, b *Entity) bool {
	// (ab)use GEOS C library for intersection detection between polygons.. ;)
	intersects, err := getGeometry(a).Intersects(getGeometry(b))
	if err != nil {
		panic(err)
	}

	return intersects
}

func getGeometry(ent *Entity) *geos.Geometry {
	var shell []geos.Coord
	for v, _ := range ent.Shape.Vectors {
		x, y := ent.Shape.Vectors[v].Rotate(ent.Angle)
		shell = append(shell, geos.Coord{x + ent.PosX, y + ent.PosY, 1})
	}
	x, y := ent.Shape.Vectors[0].Rotate(ent.Angle)
	shell = append(shell, geos.Coord{x + ent.PosX, y + ent.PosY, 1})
	geometry, err := geos.NewPolygon(shell)
	if err != nil {
		panic(err)
	}

	return geometry
}

func drawLoss(x, y float64) {
	drawCharacter(x, y, 1, "L")
	drawCharacter(x+7, y, 1, "O")
	drawCharacter(x+14, y, 1, "S")
	drawCharacter(x+21, y, 1, "E")
}

func drawCharacter(x, y, size float64, char string) {
	gl.LoadIdentity()
	gl.Begin(gl.LINES)

	switch char {
	case "L":
		vertex2d(x, y, size, 0, 8)
		vertex2d(x, y, size, 0, 0)
		vertex2d(x, y, size, 0, 0)
		vertex2d(x, y, size, 4, 0)
	case "O":
		vertex2d(x, y, size, 0, 0)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 3, 5)
		vertex2d(x, y, size, 3, 5)
		vertex2d(x, y, size, 3, 0)
		vertex2d(x, y, size, 3, 0)
		vertex2d(x, y, size, 0, 0)
	case "S":
		vertex2d(x, y, size, 4, 5)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 0, 2)
		vertex2d(x, y, size, 0, 2)
		vertex2d(x, y, size, 4, 2)
		vertex2d(x, y, size, 4, 2)
		vertex2d(x, y, size, 4, 0)
		vertex2d(x, y, size, 4, 0)
		vertex2d(x, y, size, 0, 0)
	case "E":
		vertex2d(x, y, size, 0, 0)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 0, 5)
		vertex2d(x, y, size, 4, 5)
		vertex2d(x, y, size, 0, 2)
		vertex2d(x, y, size, 3, 2)
		vertex2d(x, y, size, 0, 0)
		vertex2d(x, y, size, 4, 0)
	}

	gl.End()
}

func vertex2d(x, y, size, xv, yv float64) {
	gl.Vertex2d(x+(xv*size), y+(yv*size))
}
