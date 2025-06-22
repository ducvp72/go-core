package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	handlers "go-core/src/services/user_service/handlers"

	"github.com/gorilla/mux"
)

const (
	PORT    = 8080
	SERVICE = "user"
)

type Service struct {
	r   *mux.Router
	ctx context.Context
}

type contextKey string

func main() {
	startService()
}

func startService() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	context := context.WithValue(context.Background(), contextKey("go-core"), SERVICE)
	service := &Service{
		r:   router,
		ctx: context,
	}

	service.initHandlers()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Printf("🚀 Service is running at http://localhost:%d", PORT)
		err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), service.r)
		if err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	<-signalChan
	log.Printf("❌ Server %s stop at port: %d", SERVICE, PORT)
}

func (s *Service) initHandlers() {
	s.r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("abc"))
	}).Methods("GET")

	// privateHandlers(r, ctx)
	s.publicHandlers()
}

// func privateHandlers(r *mux.Router, ctx context.Context) {
// 	r.HandleFunc("/get-permissions", handlers.HandlerGetPermission).Methods("GET")
// }

func (s *Service) publicHandlers() {
	s.r.HandleFunc("/create-user", wrapperHandler(s.ctx, handlers.HandlerCreateUser)).Methods("POST")
}

func wrapperHandler[TReq any, TResp any](
	ctx context.Context,
	handler func(ctx context.Context, req *TReq) (*TResp, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		resp, err := handler(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		prettyJSON, err := json.MarshalIndent(resp, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}
}
