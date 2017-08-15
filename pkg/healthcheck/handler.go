package healthcheck

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

var Version string

// healthCheck return 200 current pid
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pid := os.Getpid()

	j := []string{Version, strconv.Itoa(pid)}

	if err := json.NewEncoder(w).Encode(j); err != nil {
		panic(err)
	}
	return
}
