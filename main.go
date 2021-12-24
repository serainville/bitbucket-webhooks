package main

import (
	"fmt"
	"log"

	"github.com/serainville/bitbucket-webhooks/events"
)

func main() {

	event1, err := events.NewBitbucketEvent("pr:opened", newPullRequestPayload())
	if err != nil {
		log.Println("event1", err)
	}

	event2, err := events.NewBitbucketEvent("repo:refs_changed", newPushRequestPayload())
	if err != nil {
		log.Println("event2", err)
	}

	event3, err := events.NewBitbucketEvent("repo:modified", newModifiedRequestPayload())
	if err != nil {
		log.Println("event3", err)
	}

	fmt.Println("EventType:", events.GetType(event1), " IsValid: ", event1.IsValid())
	fmt.Println("EventType:", events.GetType(event2), " IsValid: ", event2.IsValid())
	fmt.Println("EventType:", events.GetType(event3), " IsValid: ", event3.IsValid())

}

func newPullRequestPayload() []byte {
	return []byte(`{"eventKey": "pr:open", "date": "1234567890","actor": {"name": "Shane Rainville"},"pullRequest": {"title": "making a new pr!","open": true,"closed": false}}`)
}

func newPushRequestPayload() []byte {
	return []byte(`{  
		"eventKey":"repo:refs_changed",
		"date":"2017-09-19T09:45:32+1000",
		"actor":{  
		  "name":"admin",
		  "emailAddress":"admin@example.com",
		  "id":1,
		  "displayName":"Administrator",
		  "active":true,
		  "slug":"admin",
		  "type":"NORMAL"
		},
		"repository":{  
		  "slug":"repository",
		  "id":84,
		  "name":"repository",
		  "scmId":"git",
		  "state":"AVAILABLE",
		  "statusMessage":"Available",
		  "forkable":true,
		  "project":{  
			"key":"PROJ",
			"id":84,
			"name":"project",
			"public":false,
			"type":"NORMAL"
		  },
		  "public":false
		},
		"changes":[  
		  {  
			"ref":{  
			  "id":"refs/heads/master",
			  "displayId":"master",
			  "type":"BRANCH"
			},
			"refId":"refs/heads/master",
			"fromHash":"ecddabb624f6f5ba43816f5926e580a5f680a932",
			"toHash":"178864a7d521b6f5e720b386b2c2b0ef8563e0dc",
			"type":"UPDATE"
		  }
		]
	  }`)
}

func newModifiedRequestPayload() []byte {
	return []byte(`{  
		"eventKey":"repo:modified",
		"date":"2017-09-19T09:51:20+1000",
		"actor":{  
		  "name":"admin",
		  "emailAddress":"admin@example.com",
		  "id":1,
		  "displayName":"Administrator",
		  "active":true,
		  "slug":"admin",
		  "type":"NORMAL"
		},
		"old":{  
		  "slug":"repository",
		  "id":84,
		  "name":"repository",
		  "scmId":"git",
		  "state":"AVAILABLE",
		  "statusMessage":"Available",
		  "forkable":true,
		  "project":{  
			"key":"PROJ",
			"id":84,
			"name":"project",
			"public":false,
			"type":"NORMAL"
		  },
		  "public":false
		},
		"new":{  
		  "slug":"repository2",
		  "id":84,
		  "name":"repository2",
		  "scmId":"git",
		  "state":"AVAILABLE",
		  "statusMessage":"Available",
		  "forkable":true,
		  "project":{  
			"key":"PROJ",
			"id":84,
			"name":"project",
			"public":false,
			"type":"NORMAL"
		  },
		  "public":false
		}
	  }`)
}
