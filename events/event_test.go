package events

import "testing"

func TestEventsTypes(t *testing.T) {

	tt := []struct {
		name         string
		eventType    string
		eventPayload []byte
		want         string
	}{
		{
			name:         "Is a valid PullRequestEvent",
			eventType:    "pr:opened",
			eventPayload: newPullRequestPayload(),
			want:         "PullRequestEvent",
		},
		{
			name:         "Is a valid PushEvent",
			eventType:    "repo:refs_changed",
			eventPayload: newPushRequestPayload(),
			want:         "PushEvent",
		},
		{
			name:         "Is a valid ModifiedEvent",
			eventType:    "repo:modified",
			eventPayload: newModifiedRequestPayload(),
			want:         "ModifiedEvent",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			event, err := NewBitbucketEvent(tc.eventType, tc.eventPayload)
			if err != nil {
				t.Errorf("Want 'nil', got '%s'", err.Error())
			}
			eventType := GetType(event)
			eventValid := event.IsValid()

			if eventType != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, eventType)
			}
			if eventValid != nil {
				t.Errorf("Want 'nil', got '%s'", eventValid.Error())
			}
		})
	}
}

func TestNewEvent(t *testing.T) {

	tt := []struct {
		name         string
		eventType    string
		eventPayload []byte
		want         string
	}{
		{
			name:         "has this",
			eventType:    "not:valid",
			eventPayload: []byte(`{}`),
			want:         "not:valid is not a supported eventKey",
		},
		{
			name:         "test test",
			eventType:    "repo:comment:deleted payload",
			eventPayload: []byte(`{}`),
			want:         "not implemented",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewBitbucketEvent(tc.eventType, tc.eventPayload)
			if err == nil {
				t.Errorf("Want '%s', got 'nil'", tc.want)
			} else {
				if err.Error() != tc.want {
					t.Errorf("Want '%s', got '%s'", tc.want, err.Error())
				}
			}
		})
	}

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
