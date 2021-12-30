package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	bitbucket "github.com/serainville/bitbucket-webhooks"
)

// Using https://github.com/go-playground/webhooks/blob/master/bitbucket-server/bitbucketserver.go as a reference

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/endpoint", handleV1Endpoint)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleV1Endpoint(res http.ResponseWriter, req *http.Request) {
	log.Println("Processing incoming request")

	a := bitbucket.CreateWebookHandler()
	event, err := a.WebhookEvent(res, req)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(event)

}

func Test1(resp http.ResponseWriter, req *http.Request) {
	var msg string
	var processError error

	wb, _ := bitbucket.New()
	result, _ := wb.Parse(req, "PullRequestOpened", "PullRequestDeclined")

	switch result.(type) {
	case bitbucket.PullRequest:
		msg, processError = processPullRequestOpened(result)
	}

	if processError != nil {
		resp.WriteHeader(403)
	}

	log.Println(msg)
}

func processPullRequestOpened(event interface{}) (string, error) {
	var pl bitbucket.PullRequestEvent

	pl, ok := event.(bitbucket.PullRequestEvent)
	if !ok {
		return "", errors.New("could not process pr:opened event")
	}

	msg := fmt.Sprintf("Event: Pull Request Opened, Actor: %s", pl.Actor.DisplayName)
	return msg, nil
}
