package interfaces

import (
	"container/ring"
	"github.com/LeReunionais/looper/world"
	zmq "github.com/pebbe/zmq4"
	"log"
	"strconv"
)

func Integrate(endpoint string) {
	replier, errSock := zmq.NewSocket(zmq.REP)
	defer replier.Close()
	if errSock != nil {
		log.Fatal(errSock)
	}
	log.Println("Socket created")
	replier.Bind(endpoint)
	log.Println("replier bound to", endpoint)

	r := ring.New(100)
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}
	for r.Len() > 0 {
		message, _ := replier.Recv(0)
		if message == "ready" {
			log.Println("worker ready, we send him some work")
		} else {
			work_index, _ := strconv.Atoi(message)
			log.Println("one more work done")
			r.Unlink(work_index)
			replier.Send("thanks", 0)
		}
		log.Println(r.Len())
	}
	log.Println("Finish all work")
}

type work struct {
	p      world.Particle
	p_next *world.Particle
}

func find_next_work(r *ring.Ring) (*ring.Ring, bool) {
	r = r.Next()
	todo, _ := r.Value.(work)
	counter := 0
	for todo.p_next != nil && counter < r.Len() {
		counter++
		r = r.Next()
		todo, _ = r.Value.(work)
	}
	return r, counter == r.Len()
}
