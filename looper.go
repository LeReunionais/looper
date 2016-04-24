package main

import (
	"github.com/LeReunionais/looper/interface"
	"github.com/LeReunionais/looper/world"
	"github.com/LeReunionais/service"
	"log"
	"os"
	"strconv"
	"time"
)

const REGISTRY_PORT = 3001
const WORLD_PUBLICATION_PORT = 6000
const INTEGRATOR_PORT = 6001

func main() {
	w := world.Init()
	log.Printf("Looper - world initialized")
	start := time.Now()
	go interfaces.Publish("tcp", WORLD_PUBLICATION_PORT, &w)
	register()
	INTEGRATOR_ENDPOINT := "tcp://*:" + strconv.Itoa(INTEGRATOR_PORT)
	replier := interfaces.Init(INTEGRATOR_ENDPOINT)
	for {
		time.Sleep(100 * time.Millisecond)
		end := time.Now()
		delta := end.Sub(start)
		updated_particles := interfaces.Integrate(*replier, w.Particles, delta)
		w.Particles = updated_particles
		start = end
	}
}

func register() {
	host, registry := extract_env()
	registry_endpoint := "tcp://" + registry + ":" + strconv.Itoa(REGISTRY_PORT)

	world := service.Service{
		"world",
		host,
		"PUB",
		WORLD_PUBLICATION_PORT,
	}
	service.Register(registry_endpoint, world)

	integrator := service.Service{
		"integrator",
		host,
		"REP",
		INTEGRATOR_PORT,
	}
	service.Register(registry_endpoint, integrator)
}

func extract_env() (host, registry string) {
	host = os.Getenv("SERVICE_HOST")
	if host == "" {
		log.Fatal("Environment variable SERVICE_HOST is not defined. Please set it up. This variable defined the host on which is published the service.")
	}

	registry = os.Getenv("REGISTRY_HOST")
	if registry == "" {
		log.Fatal("Environment variable REGISTRY_HOST is not defined. Please set it up. Variable defined registry host.")
	}

	return
}
