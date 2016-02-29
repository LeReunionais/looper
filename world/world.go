package world

import (
	"math"
	"time"
)

type Vector3 struct {
	X, Y, Z float64
}

type Particle struct {
	Position, velocity Vector3
	inverseMass        float64
}

type World struct {
	Particles  []Particle
	lastUpdate time.Time
}

func Init() World {
	p := Vector3{0, 0, 0}
	v := Vector3{0, 0, 0}
	point := Particle{p, v, 1}

	w := World{
		Particles:  []Particle{point},
		lastUpdate: time.Now(),
	}

	return w
}

func (w *World) Update(delta time.Duration) {
	w.lastUpdate.Add(delta)
	x := math.Cos(float64(w.lastUpdate.Second()))
	p := Vector3{x, 0, 0}
	v := Vector3{0, 0, 0}
	point := Particle{p, v, 1}
	w.Particles = []Particle{point}
}
