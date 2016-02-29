package main

import (
	"github.com/LeReunionais/looper/world"
	"log"
	"time"
)

func main() {
	w := world.Init()
	log.Printf("looper")
	start := time.Now()
	for {
		time.Sleep(100 * time.Millisecond)
		end := time.Now()
		delta := end.Sub(start)
		w.Update(delta)
		start = end
	}
}
