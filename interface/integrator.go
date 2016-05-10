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
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  resultRequestParams `json:"params"`
	Id      string              `json:"id"`
}
type resultRequestParams struct {
	Particle c.Particle `json:"particle"`
	WorkId   uuid.UUID  `json:"workId"`
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

func find_next_index(to_integrate []work, integrated map[uuid.UUID]c.Particle, index int) (work_index int, work_pending bool) {
	total := len(to_integrate)
	if total == 0 {
		return 0, false
	}
	i := 0
	index++
	for {
		if i == total {
			break
		}
		if _, ok := integrated[to_integrate[index%total].UUID]; !ok {
			break
		}
		index++
		i++
	}
	return index % total, i != total
}

func Update(replier *zmq.Socket, particles []c.Particle, delta time.Duration) []c.Particle {
	log.Println("Init data stucture")
	to_integrate := make([]work, 0)
	integrated := make(map[uuid.UUID]c.Particle)

	for _, p := range particles {
		work_id := uuid.NewV4()
		w := work{work_id, p}
		to_integrate = append(to_integrate, w)
	}

	log.Println("Loop and request for integration")
	iterator_index, work_pending := find_next_index(to_integrate, integrated, 0)
	for work_pending {
		log.Println("Waiting for request")
		req_msg, _ := replier.Recv(0)
		req := new(request)
		json.Unmarshal([]byte(req_msg), req)
		log.Println("Received", req_msg)
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
			integrated[result_req.Params.WorkId] = result_req.Params.Particle
			replier.Send("Thanks", 0)
		} else {
			replier.Send("??", 0)
		}
		iterator_index, work_pending = find_next_index(to_integrate, integrated, iterator_index)
	}

	results := make([]c.Particle, len(particles))
	for _, w := range to_integrate {
		results = append(results, integrated[w.UUID])
	}
	return results
}
