package world

import (
	"encoding/json"
	"github.com/LeReunionais/looper/common"
	"sync"
	"time"
)

type World struct {
	sync.RWMutex
	Particles  []common.Particle `json:"particles"`
	lastUpdate time.Time
}

func Init() World {
	p := common.Vector3{0, 0, 0}
	v := common.Vector3{0, 0, 0}
	point := common.Particle{p, v, 1}

	w := World{
		Particles:  []common.Particle{point},
		lastUpdate: time.Now(),
	}

	return w
}

func (w *World) StateJson() string {
	w.RLock()
	defer w.RUnlock()
	json, _ := json.Marshal(w)
	return string(json)
}
