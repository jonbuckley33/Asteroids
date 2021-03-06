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
	ship           *Ship // Ship for this player.
	shipMap		   map[int]*Ship // Holds all the players/ships on the board.
	shipId		   int // Used to store player/ship info in paxos.
	PlayerId	   int
	isClient       bool
  
	bullets        []*Bullet
	torpedos       []*Torpedo
	mines          []*Mine
	asteroids      map[int]*Asteroid
	AsteroidCounter int
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
   	clientPort := flag.String("myNodeAt", "", "port at which to start game client on")
    flag.Parse()

    // Complain if flags weren't set.
    if *host == "" && *myHostPort == "" {
    	log.Fatal("Please specify a host or a port at which to serve a game.")
    } else if *host != "" && *clientPort == "" {
    	log.Fatal("You must specify a local port to host your client on with the -myNodeAt flag")
    }

	runtime.LockOSThread()
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("can't init glfw!")
	}
	defer glfw.Terminate()

	// Client or server of game?
 	isClient = *host != ""

	var window *glfw.Window
	window, err := initWindow()
	if err != nil {
		panic(err)
	}

	// Attempt to construct GameNode.
    if isClient {
    	gn, err := NewGameClient(*clientPort, *host)
    	gameNode = gn
    	if err != nil {
    		panic("Could not make game client")
    	}
    } else {
    	gn, err := NewGameServer(*myHostPort)
    	gameNode = gn
    	if err != nil {
    		panic("Could not start game server")
    	} 
    }
	PlayerId = gameNode.PlayerId    

	// Initializes data structures.
    resetGame(!isClient)

    // Start the main game loop.
	runGameLoop(window)

	fmt.Printf("Your highscore was %d points!\n", highscore)
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//create random ship
	if key == glfw.KeyU && action == glfw.Press { //&& mods == glfw.ModAlt {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		x:=r.Float64()
		y:=r.Float64()
		shipNew := NewShip(gameWidth/x, gameHeight/y, 0, 0.01)
		shipId += 1
		shipMap[shipId]=shipNew
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
		resetGame(!isClient)
	}

	if (key == glfw.KeyPause || key == glfw.KeyP) && action == glfw.Press {
		paused = !paused
	}

	if key == glfw.KeyN && action == glfw.Press && isGameWon() {
		difficulty += 3
		resetGame(!isClient)
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
	return len(asteroids) == 0 && len(shipMap)>0
}

func isGameLost() bool {
	return len(shipMap)==0
}

/* BEGIN CUSTOM CODE */

// Gets a unique ID (paxos-wide) to assign to a new Asteroid.
func NextAsteroidId() int {
	id := (AsteroidCounter << 5) | PlayerId
	AsteroidCounter += 1
	return id
}

// Initializes a game. Generates new asteroids if
// generateAsteroids is set to true. Note that clients
// shouldn't generateAsteroids, only the master should.
func resetGame(generateAsteroids bool) {
	// Init ship.
	shipId=PlayerId

	// Create ship/player map.
	shipMap=make(map[int]*Ship)

	// Create new ship.
	ship = NewShip(gameWidth/2,gameHeight/2, 0, 0.01)
	
	// Add to player list.
	shipMap[shipId] = ship

	asteroids = make(map[int]*Asteroid)

	if generateAsteroids {
		// Create a couple of random asteroids
		for i := 1; i <= difficulty; i++ {
			CreateAsteroid(2+rng.Float64()*8, 3)
		}
	}

	bullets = nil
	mines = nil
	explosions = nil
	torpedos = nil
	bigExplosions = nil
}

// Share's current user information such as player position
// and asteroid information. These calls propose values in
// the Paxos ring.
func shareGameState() {
	gameNode.SharePlayer(ship)
	gameNode.ShareAsteroids(asteroids)
}

// Queries paxos for the current asteroids and updates our
// local state to reflect this.
func updateAsteroids() {
	// Get all asteroids from Paxos.
	asteroids2 := gameNode.GetAsteroids()
	for i, v := range(asteroids2) {
		asteroid, ok := asteroids[i]
		if !ok && v.Lives > 0 {
			// New asteroid.
			asteroid = NewAsteroid(v.PosX, v.PosY, v.Angle, v.TurnRate, 
				v.VelocityX, v.VelocityY, v.SizeRatio, v.Lives)
			asteroid.Id = i
			asteroids[i] = asteroid
		} else if v.Lives > 0 {
			// Update existing asteroid.
			asteroids[i].PosX = v.PosX
			asteroids[i].PosY = v.PosY
			asteroids[i].Angle = v.Angle
			asteroids[i].TurnRate = v.TurnRate
			asteroids[i].VelocityY = v.VelocityY
			asteroids[i].VelocityX = v.VelocityX
			asteroids[i].SizeRatio = v.SizeRatio
			asteroids[i].AccelerationRate = v.AccelerationRate
			asteroids[i].Lives = v.Lives
		} else if v.Lives == 0 {
			// Delete asteroid.
			delete(asteroids, i)
		}
	}
}

// Queries paxos for the current players and updates local
// state to reflect this.
func updatePlayers(){
	paxosShips:=gameNode.GetPlayers()
	for shipId, ship := range paxosShips {
		existingShip, ok:=shipMap[shipId]
    	if ok && !ship.IsAlive() {
    		// Existing player died.
    		existingShip.Destroy()
    		delete(shipMap,shipId)
    	} else if ok {
    		// Existing playe update.
    		shipMap[shipId].PosX=ship.PosX
    		shipMap[shipId].PosY=ship.PosY
    		shipMap[shipId].Angle=ship.Angle
    		shipMap[shipId].VelocityX=ship.VelocityX
    		shipMap[shipId].VelocityY=ship.VelocityY
    		shipMap[shipId].TurnRate=ship.TurnRate
    		shipMap[shipId].AccelerationRate=ship.AccelerationRate
    	} else if ship.IsAlive() {
    		// New player added.
			shipMap[shipId] = NewShip(gameWidth/2, gameHeight/2, 0, 0.01)
    		shipMap[shipId].PosX=ship.PosX
    		shipMap[shipId].PosY=ship.PosY
    		shipMap[shipId].Angle=ship.Angle
    		shipMap[shipId].VelocityX=ship.VelocityX
    		shipMap[shipId].VelocityY=ship.VelocityY
    		shipMap[shipId].TurnRate=ship.TurnRate
    		shipMap[shipId].AccelerationRate=ship.AccelerationRate
    	}
	}
}

// Main game loop of code. Called once per game step.
func runGameLoop(window *glfw.Window) {
	for !window.ShouldClose() {
		updateObjects()
		hitDetection()

		// Upload data to Paxos.
		shareGameState()
		// Pull data from Paxos.
		updateAsteroids()
		updatePlayers()

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

/* END CUSTOM CODE */

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
};


func drawObjects() {

	for _,ships:=range shipMap{
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

	asteroids2 := make(map[int]*Asteroid)
	for _, asteroid := range asteroids {
		if asteroid.IsAlive() {
			asteroids2[asteroid.Id] = asteroid
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

	for _,ships:=range shipMap{
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
		for i,ships:=range shipMap{
			if ships.IsAlive() && IsColliding(&asteroid.Entity, &ships.Entity) {
				asteroid.Destroy()
				ships.Destroy()
				delete(shipMap,i)
				//Ships= append(Ships[:i], Ships[i+1:]...)
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

		for i,ships:=range shipMap{
			if ships.IsAlive() && IsColliding(&bigExplosion.Entity, &ships.Entity) {
				ships.Destroy()
				delete(shipMap,i)

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