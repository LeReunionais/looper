package main

import (
	"github.com/LeReunionais/looper/interface"
	"github.com/LeReunionais/looper/world"
	"github.com/LeReunionais/service"
	"log"
	"os"
	"time"
)

func main() {
	w := world.Init()
	log.Printf("Looper - world initialized")
	start := time.Now()
	go interfaces.Publish("tcp", 6000, &w)
	register()
	for {
		time.Sleep(100 * time.Millisecond)
		end := time.Now()
		delta := end.Sub(start)
		w.Update(delta)
		start = end
	}
}

func register() {
	host, registry := extract_env()
	endpoint := "tcp://" + registry + ":3001"
	port := 6000
	world := service.Service{
		"world",
		host,
		"PUB",
		port,
	}
	service.Register(endpoint, world)
}

func extract_env() (host, registry string) {
	host = os.Getenv("SERVICE_HOST")
	if host == "" {
		log.Fatal("Environment variable SERVICE_HOST is not defined. Please set it up. This variable defined the host on which is published the service")
	}

	registry = os.Getenv("REGISTRY_HOST")
	if registry == "" {
		log.Fatal("Environment variable REGISTRY_HOST is not defined. Please set it up. Variable defined registry host.")
	}

	return
}
