package interfaces

import (
	"container/ring"
	"encoding/json"
	"github.com/LeReunionais/looper/common"
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

type request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      string `json:"id"`
}

type ready_request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Name    string `json:"params"`
	Id      string `json:"id"`
}

type result_request struct {
	Jsonrpc    string          `json:"jsonrpc"`
	Method     string          `json:"method"`
	Integrated common.Particle `json:"params"`
	Id         string          `json:"id"`
}

type reply struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  reply_result `json:"result"`
	Id      string       `json:"id"`
}

type reply_result struct {
	Particle common.Particle `json:"particle"`
	Delta    float64         `json:"delta"`
}

func Init(endpoint string) *zmq.Socket {
	replier, errSock := zmq.NewSocket(zmq.REP)
	if errSock != nil {
		log.Fatal(errSock)
	}
	log.Println("Socket created")
	replier.Bind(endpoint)
	log.Println("replier bound to", endpoint)

	return replier
}
func Integrate(replier zmq.Socket, works []common.Particle, delta time.Duration) []common.Particle {

	r := ring.New(len(works))
	for _, p := range works {
		rr := reply_result{p, delta.Seconds()}
		r.Value = work{rr, nil}
		r = r.Next()
	}

	r, no_work_remaining := find_next_work(r)
	for !no_work_remaining {
		msg := new(request)
		received, _ := replier.Recv(0)
		json.Unmarshal([]byte(received), msg)
		if msg.Method == "ready" {
			log.Println("received ready")
			ready_msg := new(ready_request)
			json.Unmarshal([]byte(received), ready_msg)
			log.Println("worker ready, we send him some work")
			p_to_integrate, _ := r.Value.(work)
			work := reply{"2.0", p_to_integrate.p, msg.Id}
			workJson, _ := json.Marshal(work)
			replier.Send(string(workJson), 0)
		} else if msg.Method == "result" {
			log.Println("received result")
			result_msg := new(result_request)
			json.Unmarshal([]byte(received), result_msg)
			integrated_p := result_msg.Integrated
			current_work, _ := r.Value.(work)
			current_work.p_next = &integrated_p
			r.Value = current_work
			replier.Send("thanks", 0)
		}
		r, no_work_remaining = find_next_work(r)
	}

	log.Println("Finish all work")

	integrated_particules := make([]common.Particle, r.Len())
	for i := 0; i < r.Len(); i++ {
		integrated_particules[i] = *r.Value.(work).p_next
	}
	return integrated_particules
}

type work struct {
	p      reply_result
	p_next *common.Particle
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
