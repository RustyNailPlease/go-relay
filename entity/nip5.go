package entity

type Nip5Response struct {
	Names  map[string]string   `json:"names"`
	Relays map[string][]string `json:"relays"`
}
