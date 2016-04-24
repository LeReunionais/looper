package common

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Particle struct {
	Position    Vector3 `json:"position"`
	Velocity    Vector3 `json:"velocity"`
	InverseMass float64 `json:"inverseMass"`
}
