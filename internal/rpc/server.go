package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
)

type Server struct {
	Chain *blockchain.Chain
	Port  string
}

func NewServer(chain *blockchain.Chain, port string) *Server {
	return &Server{
		Chain: chain,
		Port:  port,
	}
}

func (s *Server) Start() error {

	http.HandleFunc("/status", s.statusHandler)
	http.HandleFunc("/height", s.heightHandler)

	return http.ListenAndServe(
		s.Port,
		nil,
	)
}

func (s *Server) statusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	response := map[string]string{
		"network": "Sayanox Chain",
		"status":  "running",
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) heightHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	response := map[string]int{
		"height": s.Chain.GetHeight(),
	}

	json.NewEncoder(w).Encode(response)
}
