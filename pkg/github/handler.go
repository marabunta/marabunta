package github

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Commit struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

type Webhook struct {
	Ref        string `json:"ref"`
	Commits    []Commit
	Repository struct {
		Name           string `json:"name"`
		Url            string `json:"url"`
		Default_branch string `json:"default_branch"`
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	secret := []byte(os.Getenv("GITHUB_TOKEN"))
	hook, err := Parse(secret, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	wh := new(Webhook)
	decoder := json.NewDecoder(bytes.NewReader(hook.Payload))
	if err := decoder.Decode(wh); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("wh = %+v\n", wh)
	log.Printf("msg = %s\n", wh.Commits[0].Message)
}
