package server

type Message struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Request struct {
	Type    string   `json:"type"`
	Users   []string `json:"users"`
	Message string   `json:"message"`
}
