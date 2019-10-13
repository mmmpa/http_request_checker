package main

import (
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	defaultPort            = "8080"
	disabledSlackChallenge = true
)

func init() {
	pp.ColoringEnabled = false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	disabledSlackChallenge = isAsFalse(os.Getenv("FOR_SLACK"))

	log.Fatal(http.ListenAndServe(":"+port, build()))
}

func isAsFalse(s string) bool {
	return s == "" || s == "false" || s == "0"
}

func isChallenge(req *http.Request) bool {
	return !disabledSlackChallenge && req.Method == http.MethodPost
}

func toRoughBody(data []byte) interface{} {
	hash := map[string]interface{}{}
	if err := json.Unmarshal(data, &hash); err != nil {
		return string(data)
	} else {
		return hash
	}
}

func build() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		pp.Println(req)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			_, _ = pp.Println(err)
		}

		if isChallenge(req) {
			challenge, err := handleChallenge(w, body)
			if err != nil {
				_, _ = pp.Println(err)
			} else {
				pp.Println(challenge)
				return
			}
		}

		data := map[string]interface{}{
			"method": req.Method,
			"url":    req.URL,
			"header": req.Header,
			"body":   toRoughBody(body),
		}
		w.WriteHeader(http.StatusOK)
		r := pp.Sprint(data)
		fmt.Println(r)
		fmt.Fprintln(w, r)
	})

	return mux
}

type SlackChallenge struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func handleChallenge(w http.ResponseWriter, body []byte) (challenge SlackChallenge, err error) {
	if err = json.Unmarshal(body, &challenge); err != nil {
		return
	}

	if challenge.Type != "url_verification" {
		err = fmt.Errorf("error: %s", "not url_verification")
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// ignore the error because it cannot recover
	_, _ = fmt.Fprint(w, challenge.Challenge)

	return
}
