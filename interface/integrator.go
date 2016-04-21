package interfaces

import (
	"container/ring"
	"encoding/json"
	"github.com/LeReunionais/looper/world"
	zmq "github.com/pebbe/zmq4"
	"log"
)

type request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  string `json:"params"`
	Id      string `json:"id"`
}

type reply struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"jsonrpc"`
	Id      string `json:"id"`
}

func Integrate(endpoint string) {
	replier, errSock := zmq.NewSocket(zmq.REP)
	defer replier.Close()
	if errSock != nil {
		log.Fatal(errSock)
	}
	log.Println("Socket created")
	replier.Bind(endpoint)
	log.Println("replier bound to", endpoint)

	workCount := 0
	for workCount == 0 {
		msg := new(request)
		received, _ := replier.Recv(0)
		json.Unmarshal([]byte(received), msg)

		if msg.Method == "ready" {
			log.Println("worker ready, we send him some work")
			work := reply{"2.0", "some work", msg.Id}
			workJson, _ := json.Marshal(work)
			replier.Send(string(workJson), 0)
		} else if msg.Method == "result" {
			log.Println("result", msg.Params)
			replier.Send("thanks", 0)
			workCount++
		}
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
