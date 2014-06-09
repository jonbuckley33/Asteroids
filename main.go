/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

var (
	ship             *Ship
	bullets          []*Bullet
	asteroids        []*Asteroid
	lastBulletFired  float64 = -1
	bulletsPerSecond float64 = 5
	gameWidth        float64
	gameHeight       float64
	fieldSize        float64 = 400
	rng                      = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("can't init glfw!")
	}
	defer glfw.Terminate()

	monitor, err := glfw.GetPrimaryMonitor()
	if err != nil {
		panic(err)
	}

	videomode, err := monitor.GetVideoMode()
	if err != nil {
		panic(err)
	}
	if videomode.Height < 480 || videomode.Width < 640 {
		panic("unsupported resolution!")
	}

	window, err := glfw.CreateWindow(videomode.Width, videomode.Height, "Golang Asteroids!", monitor, nil)
	//window, err := glfw.CreateWindow(640, 480, "Golang Asteroids!", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.SetKeyCallback(keyCallback)
	window.MakeContextCurrent()

	gl.Init()
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	width, height := window.GetFramebufferSize()
	ratio := float64(width) / float64(height)
	gameWidth = ratio * fieldSize
	gameHeight = fieldSize
	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, gameWidth, 0, gameHeight, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)

	// init
	ship = NewShip(gameWidth/2, gameHeight/2, 0, 0.01)
	// create a couple of random asteroids
	for i := float64(1); i <= 7; i++ {
		CreateAsteroid(2+rng.Float64()*8, 3)
	}

	for !window.ShouldClose() {
		// game logic

		//check if objects are still alive
		if !ship.IsAlive() {
			fmt.Println("You lost!")
			window.SetShouldClose(true)
		}

		var bullets2 []*Bullet
		for _, bullet := range bullets {
			if bullet.IsAlive() {
				bullets2 = append(bullets2, bullet)
			}
		}
		bullets = bullets2

		var asteroids2 []*Asteroid
		for _, asteroid := range asteroids {
			if asteroid.IsAlive() {
				asteroids2 = append(asteroids2, asteroid)
			}
		}
		asteroids = asteroids2

		if len(asteroids) == 0 {
			fmt.Println("You won!")
			window.SetShouldClose(true)
		}

		// update objects
		ship.Update()
		for _, bullet := range bullets {
			bullet.Update()
		}
		for _, asteroid := range asteroids {
			asteroid.Update()
		}

		// hit detection
		for _, asteroid := range asteroids {
			for _, bullet := range bullets {
				if IsColliding(&asteroid.Entity, &bullet.Entity) {
					asteroid.Destroy()
					bullet.Destroy()
				}
			}
			if IsColliding(&asteroid.Entity, &ship.Entity) {
				asteroid.Destroy()
				ship.Destroy()
			}
		}

		// ---------------------------------------------------------------
		// draw calls
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Color3d(1.0, 1.0, 1.0)

		ship.Draw()
		for _, bullet := range bullets {
			bullet.Draw()
		}
		for _, asteroid := range asteroids {
			asteroid.Draw()
		}

		gl.Flush()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//fmt.Printf("%v, %v, %v, %v\n", key, scancode, action, mods)

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	if key == glfw.KeyLeft {
		if action == glfw.Press {
			ship.RotateLeft(true)
		} else if action == glfw.Release {
			ship.RotateLeft(false)
		}
	} else if key == glfw.KeyRight {
		if action == glfw.Press {
			ship.RotateRight(true)
		} else if action == glfw.Release {
			ship.RotateRight(false)
		}
	}

	if key == glfw.KeyUp {
		if action == glfw.Press {
			ship.Accelerate(true)
		} else if action == glfw.Release {
			ship.Accelerate(false)
		}
	} else if key == glfw.KeyDown {
		if action == glfw.Press {
			ship.Decelerate(true)
		} else if action == glfw.Release {
			ship.Decelerate(false)
		}
	}

	if key == glfw.KeySpace && action == glfw.Press && glfw.GetTime() > lastBulletFired+(1/bulletsPerSecond) {
		bullet := ship.Shoot()
		bullets = append(bullets, bullet)
		lastBulletFired = glfw.GetTime()
	}
}
