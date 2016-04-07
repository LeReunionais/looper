package interfaces

import (
	"github.com/LeReunionais/looper/world"
	zmq "github.com/pebbe/zmq4"
	"log"
	"strconv"
	"time"
)

func Publish(protocol string, port int, w *world.World) {
	publisher, errSock := zmq.NewSocket(zmq.PUB)
	defer publisher.Close()
	if errSock != nil {
		log.Fatal(errSock)
	}
	log.Println("PUB socket created")

	endpoint := protocol + "://*:" + strconv.Itoa(port)
	errBind := publisher.Bind(endpoint)
	if errBind != nil {
		log.Fatal(errBind)
	}
	log.Println("PUB socket bound to", endpoint)

	for {
		time.Sleep(250 * time.Millisecond)
		state := w.StateJson()
		log.Println("Publishing", state)
		publisher.Send(state, 0)
	}
}
