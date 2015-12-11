/* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"fmt"
	"flag"
	"math/rand"
	"runtime"
	"time"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	glfw "github.com/go-gl/glfw3/v3.0/glfw"
)

var (
	ship           *Ship
	Ships 		   []*Ship
	bullets        []*Bullet
	torpedos       []*Torpedo
	mines          []*Mine
	asteroids      []*Asteroid
	explosions     []*Explosion
	bigExplosions  []*BigExplosion
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
	gameNode	   *GameNode
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
    host := flag.String("server", "", "the host:port of the game server")
    myHostPort := flag.String("hostAt", "", "port at which to start a game server")
    flag.Parse()

    // Default
    if *host == "" && *myHostPort == "" {
    	log.Fatal("Please specify a host or a port at which to serve a game.")
    }

    // We are a client of a game server. 
    if *host != "" {
    	*myHostPort = ":10029"
    	gn, err := NewGameClient(*myHostPort, *host)
    	gameNode = gn
    	if err != nil {
    		panic("Could not make game client")
    	}
    } else {
    	println("GOT THIS FAR")
    	gn, err := NewGameServer(*myHostPort)
    	println("GOT THIS FAR TOO!")
    	gameNode = gn
    	if err != nil {
    		panic("Could not start game server")
    	} 
    }

	runtime.LockOSThread()
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


	//create randome ship
	if key == glfw.KeyU && action == glfw.Press { //&& mods == glfw.ModAlt {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		x:=r.Float64()
		y:=r.Float64()
		fmt.Println("RANDOM R,",x,y)
		shipNew:=new(Ship)
		shipNew = NewShip(gameWidth/x, gameHeight/y, 0, 0.01)

		Ships=append(Ships,shipNew)
		fmt.Println("SHIPS:", Ships)
	}



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
			ship.ShootTorpedo()
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
	gl.Viewport(0, 0, int32(width), int32(height))
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
		println("Trying to create window")
		window, err = glfw.CreateWindow(int(ratio*480), 480, "Golang Asteroids!", nil, nil)
		if err != nil {
			return nil, err
		}
		println("Created window!")
		window.SetPosition(videomode.Width/2-320, videomode.Height/2-240)
	}

	window.SetKeyCallback(keyCallback)
	window.SetFramebufferSizeCallback(reshapeWindow)
	window.MakeContextCurrent()

	gl.Init()
	gl.ClearColor(GLclampf(Colorize(0)), GLclampf(Colorize(0)), GLclampf(Colorize(0)), 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	width, height := window.GetFramebufferSize()
	reshapeWindow(window, width, height)

	return window, nil
}

func GLclampf(f float64) float32 {
	if f > 1.0 {
		return 1.0
	} else if f < 0 {
		return -1.0
	}

	return float32(f)
}

func switchHighscore() {
	showHighscore = !showHighscore
}

func switchColors() {
	colorsInverted = !colorsInverted
	gl.ClearColor(GLclampf(Colorize(0)), GLclampf(Colorize(0)), GLclampf(Colorize(0)), 0.0)
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
	return len(asteroids) == 0 && len(Ships)>0
}

func isGameLost() bool {
	return len(Ships)==0
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

	Ships=Ships[:0]


	ship = NewShip(gameWidth/2,gameHeight/2, 0, 0.01)
	Ships=append(Ships,ship)

	// create a couple of random asteroids
	asteroids = nil
	for i := 1; i <= difficulty; i++ {
		CreateAsteroid(2+rng.Float64()*8, 3)
	}

	bullets = nil
	mines = nil
	explosions = nil
	torpedos = nil
	bigExplosions = nil
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

func shareGameState() {
	gameNode.SharePlayerLocation(ship.PosX, ship.PosY)
	//share velocity
}


//update ship locations functions




func runGameLoop(window *glfw.Window) {
	for !window.ShouldClose() {
		// update objects
		updateObjects()

		// hit detection
		hitDetection()

		shareGameState()

		//update ship locations function


		// println("------------------")
		// for k, v := range(gameNode.GetPlayerLocations()) {
		// 	fmt.Printf("Player %v: (%v, %v)", k, v[0], v[1])
		// }

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

				drawObjects()

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

func drawObjects() {

	for _,ships:=range Ships{
		ships.Draw(false)
	}
	for _, bullet := range bullets {
		bullet.Draw(false)
	}
	for _, torpedo := range torpedos {
		torpedo.Draw(false)
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
	for _, bigExplosion := range bigExplosions {
		bigExplosion.Draw(false)
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

	var torpedos2 []*Torpedo
	for _, torpedo := range torpedos {
		if torpedo.IsAlive() {
			torpedos2 = append(torpedos2, torpedo)
		}
	}
	torpedos = torpedos2

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

	var bigExplosions2 []*BigExplosion
	for _, bigExplosion := range bigExplosions {
		if bigExplosion.IsAlive() {
			bigExplosions2 = append(bigExplosions2, bigExplosion)
		}
	}
	bigExplosions = bigExplosions2

	for _,ships:=range Ships{
		ships.Update()
	}
	// call their update func
	for _, bullet := range bullets {
		bullet.Update()
	}
	for _, torpedo := range torpedos {
		torpedo.Update()
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
	for _, bigExplosion := range bigExplosions {
		bigExplosion.Update()
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
		// for _, torpedo := range torpedos {
		// 	if IsColliding(&asteroid.Entity, &torpedo.Entity) {
		// 		asteroid.Destroy()
		// 		torpedo.Destroy()
		// 	}
		// }
		for _, mine := range mines {
			if IsColliding(&asteroid.Entity, &mine.Entity) {
				asteroid.Destroy()
				mine.Destroy()
			}
		}
		for _, bigExplosion := range bigExplosions {
			if IsColliding(&asteroid.Entity, &bigExplosion.Entity) {
				asteroid.Destroy()
			}
		}
		for i,ships:=range Ships{
			if ships.IsAlive() && IsColliding(&asteroid.Entity, &ships.Entity) {
				asteroid.Destroy()
				ships.Destroy()
				Ships= append(Ships[:i], Ships[i+1:]...)
			}
		}	
	}
	// for _, torpedo := range torpedos {
	// 	if ship.IsAlive() && IsColliding(&torpedo.Entity, &ship.Entity) {
	// 		torpedo.Destroy()
	// 		ship.Destroy()
	// 	}
	// }
	for _, bigExplosion := range bigExplosions {

		for i,ships:=range Ships{
			if ships.IsAlive() && IsColliding(&bigExplosion.Entity, &ships.Entity) {
				ships.Destroy()
				Ships= append(Ships[:i], Ships[i+1:]...)

			}
		}
		for _, mine := range mines {
			if ship.IsAlive() && IsColliding(&mine.Entity, &ship.Entity) {
				mine.Destroy()
				ship.Destroy()
			}
		}
	}
}