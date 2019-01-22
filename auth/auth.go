package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var users = map[string]string{
	"admin":   "supersecretpassword",
	"person1": "aPassword1",
	"person2": "aPassword2",
	"person3": "aPassword3",
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticate will check the username and password
// returning a JWT token if the user is authenticated
func Authenticate(w http.ResponseWriter, r *http.Request) {
	var u user

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return
	}

	defer r.Body.Close()

	if len(u.Username) == 0 || len(u.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must provide username and Password"))
		return
	}

	password, ok := users[u.Username]
	if !ok || password != u.Password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect username or password"))
		return
	}

	token, err := getToken(u.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating JWT token: " + err.Error()))
	} else {
		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Token: " + token))
	}
}

// Verify checks the JWT token
func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		username := claims.(jwt.MapClaims)["username"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		r.Header.Set("username", username)
		r.Header.Set("role", role)

		next.ServeHTTP(w, r)
	})
}
