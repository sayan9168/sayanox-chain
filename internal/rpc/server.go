package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/mempool"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Server struct {
	Chain      *blockchain.Chain
	Mempool    *mempool.Mempool
	Port       string
	httpServer *http.Server
}

func NewServer(chain *blockchain.Chain, pool *mempool.Mempool, port string) *Server {
	return &Server{Chain: chain, Mempool: pool, Port: port}
}

func (s *Server) Name() string { return "rpc" }

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.statusHandler)
	mux.HandleFunc("/height", s.heightHandler)
	mux.HandleFunc("/mempool", s.mempoolHandler)
	mux.HandleFunc("/transaction", s.transactionHandler)

	s.httpServer = &http.Server{
		Addr:              s.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	fmt.Println("RPC server running on", s.Port)

	go func() {
		<-ctx.Done()
		_ = s.httpServer.Shutdown(context.Background())
	}()

	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"network": "Sayanox Chain",
		"status":  "running",
	})
}

func (s *Server) heightHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]int{"height": s.Chain.GetHeight()})
}

func (s *Server) mempoolHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"count": s.Mempool.Size(),
		"txs":   s.Mempool.GetAll(),
	})
}

func (s *Server) transactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx transaction.Transaction
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&tx); err != nil {
		http.Error(w, "invalid transaction payload", http.StatusBadRequest)
		return
	}

	if err := s.Mempool.Add(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "accepted",
		"hash":   tx.Hash,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
