package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
)

type Server struct {
	Chain      *blockchain.Chain
	Port       string
	httpServer *http.Server
}

func NewServer(chain *blockchain.Chain, port string) *Server {
	return &Server{
		Chain: chain,
		Port:  port,
	}
}

func (s *Server) Name() string {
	return "rpc"
}

func (s *Server) Start(ctx context.Context) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/status", s.statusHandler)
	mux.HandleFunc("/height", s.heightHandler)

	s.httpServer = &http.Server{
		Addr:    s.Port,
		Handler: mux,
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

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

func (s *Server) statusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	response := map[string]string{
		"network": "Sayanox Chain",
		"status":  "running",
	}

	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) heightHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	response := map[string]int{
		"height": s.Chain.GetHeight(),
	}

	_ = json.NewEncoder(w).Encode(response)
}
