package interfaces

import (
	"encoding/json"
	c "github.com/LeReunionais/looper/common"
	zmq "github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"testing"
	"time"
)

func worker(name, endpoint string) {
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect(endpoint)
	rand.Seed(time.Now().Unix())
	work_done := 0

	for {
		ready := readyRequest{"2.0", "ready", name, "1"}
		readyJson, _ := json.Marshal(ready)
		requester.Send(string(readyJson), 0)
		log.Println(name, "ready for work")

		message, _ := requester.Recv(0)
		log.Println(message, "received by", name)
		work_request := new(readyReply)
		json.Unmarshal([]byte(message), work_request)
		time.Sleep(time.Duration(rand.Int63n(50)) * time.Millisecond)
		log.Println(work_request)
		resultParams := resultRequestParams{c.Particle{}, work_request.Result.WorkId}
		result := resultRequest{"2.0", "result", resultParams, "id"}
		resultJson, _ := json.Marshal(result)
		log.Println("Send result", string(resultJson))
		requester.Send(string(resultJson), 0)
		work_done++
		log.Println(name, "sent result")
		requester.Recv(0)
		log.Println(name, "got a thank you note")
		log.Println(name, "has done", work_done, "work")
	}
}

func TestNormalRun(t *testing.T) {
	endpoint := "ipc://testing.integrator"
	go worker("Ringo", endpoint)
	go worker("Jimmy", endpoint)
	go worker("Joe", endpoint)
	go worker("George", endpoint)
	go worker("Robert", endpoint)
	socket := Init(endpoint)
	Update(socket, make([]c.Particle, 500), time.Second)
}

func TestFindNextIndexFindCorrectIndex(t *testing.T) {
	a := uuid.NewV4()
	b := uuid.NewV4()
	d := uuid.NewV4()
	work_1 := work{a, c.Particle{}}
	work_2 := work{b, c.Particle{}}
	work_3 := work{d, c.Particle{}}
	works := []work{work_1, work_2, work_3}

	cases := []struct {
		want, in int
	}{
		{1, 0},
		{2, 1},
		{0, 2},
	}

	log.Println("Testing different index starting point")
	for _, cas := range cases {
		actual_index, _ := find_next_index(works, map[uuid.UUID]c.Particle{}, cas.in)
		if cas.want != actual_index {
			t.Errorf("Test faild, expected '%d', got: '%d'", cas.want, actual_index)
		}
	}
}

func TestFindNextIndexDetectThatAllWorkIsDone(t *testing.T) {
	a := uuid.NewV4()
	b := uuid.NewV4()
	d := uuid.NewV4()
	work_1 := work{a, c.Particle{}}
	work_2 := work{b, c.Particle{}}
	works := []work{work_1, work_2}

	full_map := map[uuid.UUID]c.Particle{
		a: c.Particle{},
		b: c.Particle{},
	}

	half_map := map[uuid.UUID]c.Particle{
		a: c.Particle{},
	}
	not_correct_map := map[uuid.UUID]c.Particle{
		a: c.Particle{},
		d: c.Particle{},
	}
	cases := []struct {
		want bool
		in   map[uuid.UUID]c.Particle
	}{
		{true, map[uuid.UUID]c.Particle{}},
		{false, full_map},
		{true, half_map},
		{true, not_correct_map},
	}

	log.Println("Testing different result map")
	for _, cas := range cases {
		_, more_work := find_next_index(works, cas.in, 0)
		if cas.want != more_work {
			t.Errorf("Test failed, expected '%t', got: '%v'", cas.want, more_work)
		}
	}
}

func TestFindNextWithEmptyListRequest(t *testing.T) {
	cases := []struct {
		want bool
		in   []work
	}{
		{false, []work{}},
	}

	log.Println("Testing different request list")
	for _, cas := range cases {
		_, more_work := find_next_index(cas.in, map[uuid.UUID]c.Particle{}, 0)
		if cas.want != more_work {
			t.Errorf("Test failed, expected '%t', got: '%v'", cas.want, more_work)
		}
	}
}
