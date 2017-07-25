package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

// healthCheck return 200 current pid
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pid := os.Getpid()

	j := []string{version, strconv.Itoa(pid)}

	if err := json.NewEncoder(w).Encode(j); err != nil {
		panic(err)
	}
	return
}
