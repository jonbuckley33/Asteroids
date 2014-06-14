/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math"
	"strings"

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

type Char struct {
	X    float64
	Y    float64
	Size float64
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

func DrawString(x, y, size float64, color Color, text string) {
	text = strings.ToUpper(text)
	for i, c := range text {
		drawCharacter(x+(7*float64(i)*size), y, size, color, string(c))
	}
}

// this is silly, but oh well.. ;)
func drawCharacter(x, y, size float64, color Color, char string) {
	gl.LoadIdentity()
	gl.Begin(gl.LINES)

	gl.Color3d(Colorize(color.R), Colorize(color.G), Colorize(color.B))

	c := Char{x, y, size}
	switch char {
	case "0":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(2, 4)
		c.glVertex2d(2, 5)
	case "1":
		c.glVertex2d(2, 0)
		c.glVertex2d(2, 8)
		c.glVertex2d(2, 8)
		c.glVertex2d(0, 6)
	case "2":
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
	case "3":
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
	case "4":
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
	case "6":
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(0, 4)
	case "7":
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
	case "8":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
	case "9":
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
	case "C":
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
	case "E":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 4)
		c.glVertex2d(3, 4)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
	case "G":
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(2, 4)
	case "H":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
	case "I":
		c.glVertex2d(2, 0)
		c.glVertex2d(2, 8)
	case "L":
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
	case "N":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 8)
	case "O":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 0)
	case "P":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(0, 4)
	case "R":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 8)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(2, 4)
		c.glVertex2d(4, 0)
	case "S", "5":
		c.glVertex2d(4, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 4)
		c.glVertex2d(0, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 4)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(0, 0)
	case "U":
		c.glVertex2d(0, 0)
		c.glVertex2d(0, 8)
		c.glVertex2d(0, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 0)
		c.glVertex2d(4, 8)
	case "W":
		c.glVertex2d(0, 8)
		c.glVertex2d(1, 0)
		c.glVertex2d(1, 0)
		c.glVertex2d(2, 3)
		c.glVertex2d(2, 3)
		c.glVertex2d(3, 0)
		c.glVertex2d(3, 0)
		c.glVertex2d(4, 8)
	case "Y":
		c.glVertex2d(0, 8)
		c.glVertex2d(2, 4)
		c.glVertex2d(2, 4)
		c.glVertex2d(4, 8)
		c.glVertex2d(2, 4)
		c.glVertex2d(2, 0)
	case " ":
	case "-":
		c.glVertex2d(1, 4)
		c.glVertex2d(3, 4)
	case ":":
		c.glVertex2d(1, 2)
		c.glVertex2d(1, 3)
		c.glVertex2d(1, 5)
		c.glVertex2d(1, 6)
	case "!":
		c.glVertex2d(1, 0)
		c.glVertex2d(1, 1)
		c.glVertex2d(1, 3)
		c.glVertex2d(1, 8)

	}

	gl.End()
}

func (char *Char) glVertex2d(x, y float64) {
	gl.Vertex2d(char.X+(x*char.Size), char.Y+(y*char.Size))
}
