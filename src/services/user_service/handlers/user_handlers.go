package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	Name        string   `json:"name"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
}

func HandlerCreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{
		Name:        req.Name,
		Token:       "abc",
		Permissions: []string{"create", "update", "del"},
		Role:        "user",
	}, nil
}

func GetToken(w http.ResponseWriter, r *http.Request) {
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
}

func HandlerGetPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte("Hello"))
}

func privateHandlers(r *mux.Router) {
	r.HandleFunc("/get-permissions", HandlerGetPermission).Methods("GET")
}
