package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-core/src/services/user_service/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	helpers "go-core/src/helpers"
	handlers "go-core/src/services/user_service/handlers"
)

type contextKey string

func main() {
	startService()
}

type Service struct {
	ctx    context.Context
	hs     *handlers.HandlerService
	router *mux.Router
}

func startService() {
	config.LoadConfig()
	app := config.ServConfig.App
	name := config.ServConfig.Name
	port := config.ServConfig.Port
	env := config.ServConfig.Env

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	context := context.WithValue(context.Background(), contextKey(app), name)
	service := &Service{
		ctx:    context,
		hs:     &handlers.HandlerService{},
		router: r,
	}

	service.initHandlers()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		log.Printf("🚀 Service is running at [%s] http://localhost:%d", env, port)
		err := http.ListenAndServe(fmt.Sprintf(":%v", port), service.router)
		if err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	<-signalChan
	log.Printf("❌ Server %s stop at [%s] port: %d", name, env, port)
}

func (s *Service) initHandlers() {
	s.router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}).Methods("GET")

	s.router.Use(helpers.LoggerMiddleware)
	s.privateHandlers()
	s.publicHandlers()
}

func (s *Service) privateHandlers() {
	s.router.HandleFunc("/get-permissions", wrapperHandler(s.ctx, s.hs.HandlerGetPermission)).Methods("GET")
}

func (s *Service) publicHandlers() {
	s.router.HandleFunc("/create-user", wrapperHandler(s.ctx, s.hs.HandlerCreateUser)).Methods("POST")
}

func wrapperHandler[TReq any, TResp any](
	ctx context.Context,
	handler func(ctx context.Context, req *TReq) (*TResp, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req TReq
		if r.Method == "GET" {
			err := convertUrlToString(r.URL.Query(), &req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
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

func convertUrlToString[TReq any](urlMap map[string][]string, req *TReq) error {
	mapString := make(map[string]string)

	for k, v := range urlMap {
		mapAlias := mapAlias[k]
		mapString[mapAlias] = v[0]
	}

	data, err := json.Marshal(mapString)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	return nil
}

var (
	mapAlias = map[string]string{
		"u":    "username",
		"code": "appCode",
	}
)
