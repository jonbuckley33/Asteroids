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
	ship           *Ship
	bullets        []*Bullet
	mines          []*Mine
	asteroids      []*Asteroid
	explosions     []*Explosion
	gameWidth      float64
	gameHeight     float64
	fieldSize      float64 = 400
	fullscreen     bool    = false
	altEnter       bool    = false
	colorsInverted bool    = false
	wireframe      bool    = true
	paused         bool    = false
	rng                    = rand.New(rand.NewSource(time.Now().UnixNano()))
	score          int     = 0
	highscore      int     = 0
	showHighscore  bool    = true
	difficulty     int     = 6
	debug          bool    = true
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

		if key == glfw.KeySpace || key == glfw.KeyX {
			if action == glfw.Press {
				ship.Shoot(true)
			} else if action == glfw.Release {
				ship.Shoot(false)
			}
		}

		if (key == glfw.KeyY || key == glfw.KeyZ || key == glfw.KeyLeftShift || key == glfw.KeyRightShift) && action == glfw.Press {
			ship.DropMine()
		}

		if (key == glfw.KeyC || key == glfw.KeyLeftControl || key == glfw.KeyRightControl) && action == glfw.Press {
			//ship.ShootTorpedo()
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
		score = 0
		resetGame()
	}

	if (key == glfw.KeyPause || key == glfw.KeyP) && action == glfw.Press {
		paused = !paused
	}

	if key == glfw.KeyN && action == glfw.Press && isGameWon() {
		difficulty += 3
		resetGame()
	}

	if debug && key == glfw.KeyF10 && action == glfw.Press {
		for _, asteroid := range asteroids {
			if asteroid.IsAlive() {
				asteroid.Destroy()
			}
		}
		for _, mine := range mines {
			if mine.IsAlive() {
				mine.Destroy()
			}
		}
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

func isGameWon() bool {
	return len(asteroids) == 0 && ship.IsAlive()
}

func isGameLost() bool {
	return !ship.IsAlive()
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
	// init ship
	ship = NewShip(gameWidth/2, gameHeight/2, 0, 0.01)

	// create a couple of random asteroids
	asteroids = nil
	for i := 1; i <= difficulty; i++ {
		CreateAsteroid(2+rng.Float64()*8, 3)
	}

	bullets = nil
	mines = nil
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
	DrawString(fieldSize/2-120, fieldSize/2-20, 1.5, Color{1, 1, 1}, fmt.Sprintf("Press R to restart current level"))
	DrawString(fieldSize/2-120, fieldSize/2-50, 1.5, Color{1, 1, 1}, fmt.Sprintf("Press N to advance to next difficulty level"))
}

func drawGameOverScreen() {
	DrawString(fieldSize/2-20, fieldSize/2+10, 5, Color{1, 1, 1}, fmt.Sprintf("Game Over!"))
	DrawString(fieldSize/2-120, fieldSize/2-20, 1.5, Color{1, 1, 1}, fmt.Sprintf("Press R to restart current level"))
}

func addScore(value int) {
	if ship.IsAlive() {
		score = score + value
	}
}

func runGameLoop(window *glfw.Window) {
	for !window.ShouldClose() {
		// update objects
		updateObjects()

		// hit detection
		hitDetection()

		// ---------------------------------------------------------------
		// draw calls
		gl.Clear(gl.COLOR_BUFFER_BIT)

		drawCurrentScore()
		drawHighScore()

		if isGameWon() {
			drawWinningScreen()
		} else if isGameLost() {
			drawGameOverScreen()
		}

		// draw everything 9 times in a 3x3 grid stitched together for seamless clipping
		for x := -1.0; x < 2.0; x++ {
			for y := -1.0; y < 2.0; y++ {
				gl.MatrixMode(gl.MODELVIEW)
				gl.PushMatrix()
				gl.Translated(gameWidth*x, gameHeight*y, 0)

				ship.Draw(false)
				for _, bullet := range bullets {
					bullet.Draw()
				}
				for _, mine := range mines {
					mine.Draw(false)
				}
				for _, asteroid := range asteroids {
					asteroid.Draw(true)
				}
				for _, explosion := range explosions {
					explosion.Draw()
				}

				gl.PopMatrix()
			}
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

func updateObjects() {
	//check if objects are still alive
	var bullets2 []*Bullet
	for _, bullet := range bullets {
		if bullet.IsAlive() {
			bullets2 = append(bullets2, bullet)
		}
	}
	bullets = bullets2

	var mines2 []*Mine
	for _, mine := range mines {
		if mine.IsAlive() {
			mines2 = append(mines2, mine)
		}
	}
	mines = mines2

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

	// call their update func
	ship.Update()
	for _, bullet := range bullets {
		bullet.Update()
	}
	for _, mine := range mines {
		mine.Update()
	}
	for _, asteroid := range asteroids {
		asteroid.Update()
	}
	for _, explosion := range explosions {
		explosion.Update()
	}
}

func hitDetection() {
	for _, asteroid := range asteroids {
		for _, bullet := range bullets {
			if IsColliding(&asteroid.Entity, &bullet.Entity) {
				asteroid.Destroy()
				bullet.Destroy()
			}
		}
		for _, mine := range mines {
			if IsColliding(&asteroid.Entity, &mine.Entity) {
				asteroid.Destroy()
				mine.Destroy()
			}
		}
		if ship.IsAlive() && IsColliding(&asteroid.Entity, &ship.Entity) {
			asteroid.Destroy()
			ship.Destroy()
		}
	}
	for _, mine := range mines {
		if ship.IsAlive() && IsColliding(&mine.Entity, &ship.Entity) {
			mine.Destroy()
			ship.Destroy()
		}
	}
}
