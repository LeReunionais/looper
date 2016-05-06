package interfaces

import (
	"encoding/json"
	c "github.com/LeReunionais/looper/common"
	zmq "github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      string `json:"id"`
}

type readyRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Name    string `json:"params"`
	Id      string `json:"id"`
}

type readyReply struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  readyResult `json:"result"`
	Id      string      `json:"id"`
}

type readyResult struct {
	Particle c.Particle `json:"particle"`
	Delta    float64    `json:"delta"`
	WorkId   uuid.UUID  `json:"workId"`
}

type resultRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Particle c.Particle `json:"particle"`
		WorkId   uuid.UUID  `json:"workId"`
	} `json:"params"`
	Id string `json:"id"`
}

type resultReply struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      string `json:"id"`
}

type work struct {
	uuid.UUID
	c.Particle
}

func Init(endpoint string) *zmq.Socket {
	replier, errSock := zmq.NewSocket(zmq.REP)
	if errSock != nil {
		log.Fatal(errSock)
	}
	log.Println("replier created")
	errBind := replier.Bind(endpoint)
	if errBind != nil {
		log.Fatal(errBind)
	}
	log.Println("replier bound to", endpoint)

	return replier
}

func find_next_index(to_integrate []work, integrated map[uuid.UUID]c.Particle, index int) (int, bool) {
	total := len(to_integrate)
	i := 0
	for {
		if _, ok := integrated[to_integrate[index%total].UUID]; !ok {
			break
		}
		if i == total {
			break
		}
		index++
		i++
	}
	return index % total, i != total
}

func Update(replier zmq.Socket, particles []c.Particle, delta time.Duration) []c.Particle {
	log.Println("Init data stucture")
	to_integrate := make([]work, len(particles))
	integrated := make(map[uuid.UUID]c.Particle)

	for _, p := range particles {
		work_id := uuid.NewV4()
		w := work{work_id, p}
		to_integrate = append(to_integrate, w)
	}

	log.Println("Loop and request for integration")
	iterator_index, work_pending := find_next_index(to_integrate, integrated, 0)
	for work_pending {
		req_msg, _ := replier.Recv(0)
		log.Println("request received", req_msg)
		req := new(request)
		json.Unmarshal([]byte(req_msg), req)
		log.Println("Received", req.Method, "request")
		if req.Method == "ready" {
			worker_ready_req := new(readyRequest)
			json.Unmarshal([]byte(req_msg), worker_ready_req)
			log.Println("Worker", worker_ready_req.Name, "ready.")

			next_p := to_integrate[iterator_index].Particle
			next_UUID := to_integrate[iterator_index].UUID
			worker_ready_rep := readyReply{
				"2.0",
				readyResult{next_p, delta.Seconds(), next_UUID},
				worker_ready_req.Id,
			}
			worker_ready_rep_json, _ := json.Marshal(worker_ready_rep)
			replier.Send(string(worker_ready_rep_json), 0)
		} else if req.Method == "result" {
			result_req := new(resultRequest)
			json.Unmarshal([]byte(req_msg), result_req)
		}
		iterator_index, work_pending = find_next_index(to_integrate, integrated, iterator_index)
	}

	results := make([]c.Particle, len(particles))
	for _, p := range integrated {
		results = append(results, p)
	}
	return results
}

/*
func Integrate(replier zmq.Socket, works []common.Particle, delta time.Duration) []common.Particle {

	r := ring.New(len(works))
	work_map := make(map[uuid.UUID]*common.Particle)
	for _, p := range works {
		work_id := uuid.NewV4()
		rr := reply_result{p, delta.Seconds(), work_id}
		var p_to_integrate common.Particle
		r.Value = work{rr, &p_to_integrate}
		work_map[work_id] = &p_to_integrate
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
			integrated_p := result_msg.Work.Integrated
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
*/
