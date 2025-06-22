package handlers

import "net/http"

func HandlerAdminSetPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte("Create"))
}
