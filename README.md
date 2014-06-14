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
* `go get github.com/go-gl/gltext`
* `go get github.com/go-gl/glu`
* `go get github.com/paulsmith/gogeos/geos`
* `go get github.com/JamesClonk/asteroids`
* `go build`

### Run

`./asteroids`

### Todo

* add stars / starfield background
* add ship lives (3? show also as icons on HUD?)
* grant 1 extra live every x points scored
* add torpedos with blast radius (propelled timed bombs, they explode after timer, not upon contact (blast radius affects player ship too))
* add enemy ships (random appearance?)
