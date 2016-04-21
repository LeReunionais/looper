package world

import (
	"encoding/json"
	"math"
	"sync"
	"time"
)

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Particle struct {
	Position    Vector3 `json:"position"`
	velocity    Vector3
	inverseMass float64
}

type World struct {
	sync.RWMutex
	Particles  []Particle `json:"particles"`
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
	w.Lock()
	w.lastUpdate = w.lastUpdate.Add(delta)
	t := float64(w.lastUpdate.UnixNano()) * 0.0000000001 * math.Pi / 30
	y := math.Cos(t)
	p := Vector3{t, y, 0}
	v := Vector3{0, 0, 0}
	point := Particle{p, v, 1}
	w.Particles = []Particle{point}
	w.Unlock()
}

func (w *World) StateJson() string {
	w.RLock()
	defer w.RUnlock()
	json, _ := json.Marshal(w)
	return string(json)
}
