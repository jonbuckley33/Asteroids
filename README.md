Asteroids
=========

OpenGL Asteroids in Golang..

### Installation instructions

* `pacman -S glew`
* `pacman -S glfw`
* `pacman -S geos`

---

* `go get github.com/go-gl/gl`
* `go get github.com/go-gl/glfw3`
* `go get github.com/paulsmith/gogeos/geos`
* `go get github.com/JamesClonk/asteroids`
* `go build`

### Run

`./asteroids`

### Todo

* add statusbar (or HUD?) at top (or bottom?), with inverted background color (displaying lives, score, keybindings, etc..? or display all this without statusbar? as HUD?)
* add stars / starfield background
* add score
* add ship lives (3?)
* add increasing difficulty after each clearing of all asteroids (+1 asteroid more?)
* add proper game over screen
* display score and past highscore at game over (all lives used up)
* add restart game functionality (backspace?)
* seamless zone clipping
* add torpedos with blast radius (propelled timed bombs, they explode after timer, not upon contact (blast radius affects player ship too))
* add mines (they explode upon contact (even with player ship))
* add enemy ships (random appearance?)
