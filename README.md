![example workflow](https://github.com/serainville/bitbucket-webhooks/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/serainville/bitbucket-webhooks)](https://goreportcard.com/report/github.com/serainville/bitbucket-webhooks)
# Go Bitbucket Webhook Module
## Overview
This module is used for handling Bitbucket Server Webhook events. Combined with an `http.Server`, events can be validated and processed.

## Parsing Webhook Events
Events are parsed using the the `Parse(*http.Request)` Webhook method. The method accepts `http.Request` variables.

```golang

import (
    webhook "github.com/serainville/bitbucket-webhooks"
)

func handleEvent(w http.ResponseWriter, r *http.Request) {
    hook := webhook.New()

    err := hook.Parse(r)
    if err != nil {
        r.Write([]byte("failed to parse webhook"))
        log.Println(err)
    }
}

```

## HMAC Signature
This module features HMAC signature validation to ensure webhook requests are authentic and sent from a trusted source. HMAC signature hashes are generated using the request body and a secret key. To verify the integrity of the received event, the a signature is generated by this module using the request body and a key.

In order for HMAC signatures to be validated, a webhook must have a secret set.

```golang
hook := webhook.New()
hook.Secret("WEBHOOK_SECRET")
```

## Examples
### Handling Events
The `Parse(*http.Request)` does not return a struct. Rather, an `interface{}` is returned instead. By doing so, `Parse()` is capable of returning a variety of event types.

The event type of a returned event can be reflected back using `event.(type)`. This allows further processing of returned events based on its event type.



The following example shows how s how to 

```golang
func handlePullRequests(resp http.ResponseWriter, req *http.Request) {
	hook := v1.New()
	hook.Secret("WEBHOOK_SECRET")

	event, err := hook.Parse(req)
	if err != nil {
		resp.Write([]byte(err.Error()))
		resp.WriteHeader(403)
		log.Printf("Error: %v", err)
		return
	}

	var rMessage ResponseMessage

	resp.Header().Add("Content-Type", "application/json")

	switch evt := event.(type) {
	case v1.DiagnosticPingEvent:
		rMessage.Msg = "Ignoring diagnostic:ping event"
	case v1.PullRequestOpenedPayload:
		rMessage.Msg = fmt.Sprintf("Processing event '%s', submitter: %s", evt.EventKey, evt.Actor.DisplayName)
	case v1.PullRequestDeclinedPayload:
		rMessage.Msg = fmt.Sprintf("Processing event '%s', submitter: %s", evt.EventKey, evt.Actor.DisplayName)
	default:
		rMessage.Msg = fmt.Sprintf("Event not processed. '%s' events are not supported", req.Header.Get("X-Event-Key"))
		resp.WriteHeader(501)
	}

	jsonBody, _ := json.Marshal(rMessage)
	resp.Write(jsonBody)
	log.Println(rMessage.Msg)
}
```

### Slack
By incorporating the Slack module, events processed by this module can be used to trigger Slack messages, for example.

The following examples sends a message to a Slack channel whenever a Bitbucket Webhook event is processed.

```golang
import (
	"fmt"

    bitbucket "github.com/serainville/bitbucket-webhooks"
	"github.com/slack-go/slack"
)

func handlePullRequests(resp http.ResponseWriter, req *http.Request) {
	hook := bitbucket.New()
	hook.Secret("WEBHOOK_SECRET")

	event, err := hook.Parse(req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

    var msg string
    var err error

	switch evt := event.(type) {
	case v1.PullRequestOpenedPayload:
        msg = fmt.Sprintf("Pull Request Event '%s' submitted by %s -- Repo: %s/%s, PR: %d", 
            evt.EventKey, evt.Actor.DisplayName, evt.PullRequest.FromRef.Project.Name, evt.PullRequest.FromRef.Repository.Name, evt.PullRequest.ID)
        err = sendSlackMessage(msg, slackChannelID)
	default:
        // Nothing to do
        return
	}

    if err != nil {
        log.Println(err)
    }
}


func sendSlackMessage(message string, channelID string) error {
    api := slack.New("TOKEN")

    channelID, timestamp, err := api.PostMessage(
        channelID,
        slack.MsgOptionText(message, false),
    )

    if err != nil {
        return fmt.Errorf("failed to send Slack message to channel: %w", err, channelID)
    }

    return nil
}

```