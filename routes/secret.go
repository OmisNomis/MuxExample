package routes

import (
	"fmt"
	"net/http"
)

// SecretHandler comment
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	username := r.Header.Get("username")
	role := r.Header.Get("role")

	resp := fmt.Sprintf("Ye are a priviledged one %s! You have a role of %s", username, role)

	w.Write([]byte(resp))
	return
}
