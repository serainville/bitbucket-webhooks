package main

import (
	"log"
	"net/http"
)

// Using https://github.com/go-playground/webhooks/blob/master/bitbucket-server/bitbucketserver.go as a reference

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/pullrequest", handlePullRequest)
	mux.HandleFunc("/v1/repo", handleRepoRequest)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handlePullRequest(res http.ResponseWriter, req *http.Request) {
}

func handleRepoRequest(res http.ResponseWriter, req *http.Request) {
}
