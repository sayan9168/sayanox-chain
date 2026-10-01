package rpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
)

type Server struct {
	Chain *blockchain.Chain
	Port  string
}

func NewServer(chain *blockchain.Chain, port string) *Server {
	return &Server{Chain: chain, Port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.statusHandler)
	mux.HandleFunc("/height", s.heightHandler)
	mux.HandleFunc("/block/latest", s.latestBlockHandler)
	mux.HandleFunc("/block/", s.blockHandler)
	mux.HandleFunc("/balance/", s.balanceHandler)
	return http.ListenAndServe(s.Port, mux)
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]interface{}{
		"network": "Sayanox Chain",
		"status": "running",
		"height": s.Chain.GetHeight(),
	})
}

func (s *Server) heightHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]int{"height": s.Chain.GetHeight()})
}

func (s *Server) latestBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	writeJSON(w, s.Chain.LatestBlock())
}

func (s *Server) blockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	raw := strings.TrimPrefix(r.URL.Path, "/block/")
	height, err := strconv.ParseUint(raw, 10, 64)
	if err != nil { http.Error(w, "invalid height", http.StatusBadRequest); return }
	block := s.Chain.GetBlock(height)
	if block == nil { http.Error(w, "block not found", http.StatusNotFound); return }
	writeJSON(w, block)
}

func (s *Server) balanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	address := strings.TrimPrefix(r.URL.Path, "/balance/")
	if address == "" { http.Error(w, "missing address", http.StatusBadRequest); return }
	writeJSON(w, map[string]string{"address": address, "balance": s.Chain.State.GetBalance(address).String()})
}

func writeJSON(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
