/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
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
	explosions       []*Explosion
	lastBulletFired  float64 = -1
	bulletsPerSecond float64 = 5
	gameWidth        float64
	gameHeight       float64
	fieldSize        float64 = 400
	fullscreen       bool    = false
	altEnter         bool    = false
	colorsInverted   bool    = false
	wireframe        bool    = true
	paused           bool    = false
	rng                      = rand.New(rand.NewSource(time.Now().UnixNano()))
	score            int     = 0
	highscore        int     = 0
	showHighscore    bool    = true
	difficulty       int     = 6
	debug            bool    = true
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

	var window *glfw.Window = initGame()
	runGameLoop(window)

	fmt.Printf("Your highscore was %d points!\n", highscore)
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//fmt.Printf("%v, %v, %v, %v\n", key, scancode, action, mods)

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	if !paused {
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

		if key == glfw.KeySpace && action == glfw.Press && glfw.GetTime() > lastBulletFired+(1/bulletsPerSecond) && ship.IsAlive() {
			bullet := ship.Shoot()
			bullets = append(bullets, bullet)
			lastBulletFired = glfw.GetTime()
		}
	}

	if key == glfw.KeyEnter && action == glfw.Press { //&& mods == glfw.ModAlt {
		altEnter = true
	}

	if key == glfw.KeyF1 && action == glfw.Press {
		switchHighscore()
	}

	if key == glfw.KeyF2 && action == glfw.Press {
		switchColors()
	}

	if key == glfw.KeyF3 && action == glfw.Press {
		switchWireframe()
	}

	if (key == glfw.KeyF9 || key == glfw.KeyR || key == glfw.KeyBackspace) && action == glfw.Press {
		resetGame()
	}

	if (key == glfw.KeyPause || key == glfw.KeyP) && action == glfw.Press {
		paused = !paused
	}

	if key == glfw.KeyN && action == glfw.Press && len(asteroids) == 0 && ship.IsAlive() {
		difficulty += 3
		resetGame()
	}

	if debug && key == glfw.KeyF10 && action == glfw.Press {
		asteroids = nil
	}
}

func reshapeWindow(window *glfw.Window, width, height int) {
	ratio := float64(width) / float64(height)
	gameWidth = ratio * fieldSize
	gameHeight = fieldSize
	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, gameWidth, 0, gameHeight, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
	if wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}
}

func initWindow() (window *glfw.Window, err error) {
	monitor, err := glfw.GetPrimaryMonitor()
	if err != nil {
		return nil, err
	}
	videomode, err := monitor.GetVideoMode()
	if err != nil {
		return nil, err
	}
	if videomode.Height < 480 || videomode.Width < 640 {
		return nil, errors.New("unsupported resolution!")
	}

	ratio := float64(videomode.Width) / float64(videomode.Height)

	if fullscreen {
		glfw.WindowHint(glfw.Decorated, 0)
		window, err = glfw.CreateWindow(videomode.Width, videomode.Height, "Golang Asteroids!", nil, nil)
		if err != nil {
			return nil, err
		}
		window.SetPosition(0, 0)
	} else {
		glfw.WindowHint(glfw.Decorated, 1)
		window, err = glfw.CreateWindow(int(ratio*480), 480, "Golang Asteroids!", nil, nil)
		if err != nil {
			return nil, err
		}
		window.SetPosition(videomode.Width/2-320, videomode.Height/2-240)
	}

	window.SetKeyCallback(keyCallback)
	window.SetFramebufferSizeCallback(reshapeWindow)
	window.MakeContextCurrent()

	gl.Init()
	gl.ClearColor(gl.GLclampf(Colorize(0)), gl.GLclampf(Colorize(0)), gl.GLclampf(Colorize(0)), 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	width, height := window.GetFramebufferSize()
	reshapeWindow(window, width, height)

	return window, nil
}

func switchHighscore() {
	showHighscore = !showHighscore
}

func switchColors() {
	colorsInverted = !colorsInverted
	gl.ClearColor(gl.GLclampf(Colorize(0)), gl.GLclampf(Colorize(0)), gl.GLclampf(Colorize(0)), 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func switchWireframe() {
	wireframe = !wireframe
	if wireframe {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func initGame() *glfw.Window {
	window, err := initWindow()
	if err != nil {
		panic(err)
	}

	resetGame()

	return window
}

func resetGame() {
	score = 0

	// init ship
	ship = NewShip(gameWidth/2, gameHeight/2, 0, 0.01)

	// create a couple of random asteroids
	asteroids = nil
	for i := 1; i <= difficulty; i++ {
		CreateAsteroid(2+rng.Float64()*8, 3)
	}

	bullets = nil
	explosions = nil
}

func drawHighScore() {
	if score > highscore {
		highscore = score
	}
	if showHighscore {
		DrawString(10, fieldSize-32, 1, Color{0.5, 0.5, 0.5}, fmt.Sprintf("highscore: %d", highscore))
	}
}

func drawCurrentScore() {
	DrawString(10, fieldSize-20, 1, Color{1, 1, 1}, fmt.Sprintf("score: %d", score))
}

func drawWinningScreen() {
	DrawString(fieldSize/2-20, fieldSize/2+10, 5, Color{1, 1, 1}, fmt.Sprintf("You won!"))
	DrawString(fieldSize/2-120, fieldSize/2-20, 2, Color{1, 1, 1}, fmt.Sprintf("Press R to restart current level"))
	DrawString(fieldSize/2-120, fieldSize/2-50, 2, Color{1, 1, 1}, fmt.Sprintf("Press N to advance to next difficulty level"))
}

func addScore(value int) {
	score = score + value
}

func runGameLoop(window *glfw.Window) {
	for !window.ShouldClose() {
		//check if objects are still alive
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

		var explosions2 []*Explosion
		for _, explosion := range explosions {
			if explosion.IsAlive() {
				explosions2 = append(explosions2, explosion)
			}
		}
		explosions = explosions2

		// update objects
		ship.Update()
		for _, bullet := range bullets {
			bullet.Update()
		}
		for _, asteroid := range asteroids {
			asteroid.Update()
		}
		for _, explosion := range explosions {
			explosion.Update()
		}

		// hit detection
		for _, asteroid := range asteroids {
			for _, bullet := range bullets {
				if IsColliding(&asteroid.Entity, &bullet.Entity) {
					asteroid.Destroy()
					bullet.Destroy()
				}
			}
			if ship.IsAlive() && IsColliding(&asteroid.Entity, &ship.Entity) {
				asteroid.Destroy()
				ship.Destroy()
			}
		}

		// ---------------------------------------------------------------
		// draw calls
		gl.Clear(gl.COLOR_BUFFER_BIT)

		ship.Draw(false)
		for _, bullet := range bullets {
			bullet.Draw()
		}
		for _, asteroid := range asteroids {
			asteroid.Draw(true)
		}
		for _, explosion := range explosions {
			explosion.Draw()
		}

		drawCurrentScore()
		drawHighScore()

		if len(asteroids) == 0 && ship.IsAlive() {
			drawWinningScreen()
		}

		gl.Flush()
		window.SwapBuffers()
		glfw.PollEvents()

		// switch resolution
		if altEnter {
			window.Destroy()

			fullscreen = !fullscreen
			var err error
			window, err = initWindow()
			if err != nil {
				panic(err)
			}

			altEnter = false

			gl.LineWidth(1)
			if fullscreen {
				gl.LineWidth(2)
			}
		}
	}
}
