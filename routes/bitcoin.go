package routes

import (
	"fmt"
	"io"
	"net/http"

	"../bitcoin"
	"github.com/ordishs/gocore"
)

var logger = gocore.Log("APIServer")

var bHost, _ = gocore.Config().Get("RSV_host")
var bPort, _ = gocore.Config().GetInt("RSV_port")
var bUser, _ = gocore.Config().Get("RSV_user")
var bPassword, _ = gocore.Config().Get("RSV_password")

var bitcoind, _ = bitcoin.New(bHost, bPort, bUser, bPassword, false)

// GetDifficulty returns the proof-of-work difficulty from Bitcoin
func GetDifficulty(w http.ResponseWriter, r *http.Request) {
	d, err := bitcoind.GetDifficulty()
	if err != nil {
		logger.Errorf("Error getting difficulty: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"difficuly": %f}`, d))

}
