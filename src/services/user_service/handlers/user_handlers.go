package handlers

import (
	"context"
	"net/http"
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

type GetUserPermissionRequest struct {
	Username string `json:"username"`
	AppCode  string `json:"appCode"`
}

type GetUserPermissionResponse struct {
	Permissions []string `json:"permissions"`
}

func (s *HandlerService) HandlerCreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
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

func (s *HandlerService) HandlerGetPermission(ctx context.Context, req *GetUserPermissionRequest) (*GetUserPermissionResponse, error) {
	return nil, nil
}
