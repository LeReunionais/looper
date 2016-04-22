package interfaces

import (
	"container/ring"
	"encoding/json"
	"github.com/LeReunionais/looper/world"
	zmq "github.com/pebbe/zmq4"
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
		ready := ready_request{"2.0", "ready", "", "1"}
		readyJson, _ := json.Marshal(ready)
		requester.Send(string(readyJson), 0)
		log.Println(name, "ready for work")

		message, _ := requester.Recv(0)
		log.Println(message, "received by", name)
		time.Sleep(time.Duration(rand.Int63n(1500)) * time.Millisecond)

		integrated_p := world.Particle{
			world.Vector3{1, 0, 0},
			world.Vector3{1, 0, 0},
			2.0,
		}
		result := result_request{"2.0", "result", integrated_p, "1"}
		resultJson, _ := json.Marshal(result)
		requester.Send(string(resultJson), 0)
		work_done++
		log.Println(name, "sent result")
		requester.Recv(0)
		log.Println(name, "got a thank you note")
		log.Println(name, "has done", work_done, "work")
	}
}

func TestIntegrate(t *testing.T) {
	particule_1 := world.Particle{Position: world.Vector3{1, 0, 0}}
	particule_2 := world.Particle{Position: world.Vector3{2, 0, 0}}
	particule_3 := world.Particle{Position: world.Vector3{3, 0, 0}}
	particule_4 := world.Particle{Position: world.Vector3{4, 0, 0}}
	particule_5 := world.Particle{Position: world.Vector3{5, 0, 0}}
	particules := []world.Particle{
		particule_1,
		particule_2,
		particule_3,
		particule_4,
		particule_5,
	}
	go worker("Joe", "ipc://testing.integrator")
	go worker("Ringo", "ipc://testing.integrator")
	go worker("Harry", "ipc://testing.integrator")
	go worker("George", "ipc://testing.integrator")
	result := Integrate("ipc://testing.integrator", particules)
	for _, p := range result {
		log.Println(p)
	}
}

func TestFindNextWorkRing(t *testing.T) {
	particule_1 := world.Particle{Position: world.Vector3{1, 0, 0}}
	particule_2 := world.Particle{Position: world.Vector3{2, 0, 0}}
	particule_3 := world.Particle{Position: world.Vector3{3, 0, 0}}
	particule_4 := world.Particle{Position: world.Vector3{4, 0, 0}}
	particule_5 := world.Particle{Position: world.Vector3{5, 0, 0}}
	particules := []work{
		work{particule_1, nil},
		work{particule_2, nil},
		work{particule_3, &particule_3},
		work{particule_4, nil},
		work{particule_5, nil},
	}

	r := ring.New(len(particules))
	for _, p := range particules {
		r.Value = p
		r = r.Next()
	}

	cases := []struct {
		want int
	}{
		{2}, {4}, {5}, {1}, {2},
	}
	for _, c := range cases {
		r, _ = find_next_work(r)
		todo, _ := r.Value.(work)
		if todo.p.Position.X != float64(c.want) {
			t.Errorf("find_next_work(r).p.Position.X return %f, want %d", todo.p.Position.X, c.want)
		}
	}

}

func TestFindNextWorkRingNoMoreWorkd(t *testing.T) {
	particule_1 := world.Particle{Position: world.Vector3{1, 0, 0}}
	particule_2 := world.Particle{Position: world.Vector3{2, 0, 0}}
	particule_3 := world.Particle{Position: world.Vector3{3, 0, 0}}
	particule_4 := world.Particle{Position: world.Vector3{4, 0, 0}}
	particule_5 := world.Particle{Position: world.Vector3{5, 0, 0}}
	particules := []work{
		work{particule_1, &particule_1},
		work{particule_2, &particule_2},
		work{particule_3, &particule_3},
		work{particule_4, &particule_4},
		work{particule_5, &particule_5},
		work{particule_5, nil},
	}
	r := ring.New(len(particules))
	for _, p := range particules {
		r.Value = p
		r = r.Next()
	}
	r, no_more_work := find_next_work(r)
	if no_more_work {
		t.Errorf("There should be some work available")
	}
	todo, _ := r.Value.(work)
	todo.p_next = &particule_5
	r.Value = todo

	_, no_more_work = find_next_work(r)
	if !no_more_work {
		t.Errorf("There should be no work available")
	}
}
