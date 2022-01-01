package bitbucketwebhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tc := []struct {
		Name         string
		Body         io.Reader
		EventKey     string
		Header       map[string][]string
		ExpectedErr  bool
		Secret       string
		ExpectedType interface{}
	}{
		{
			Name:         "Valid pr:opened",
			Body:         NewPullRequestOpened(),
			ExpectedErr:  false,
			Secret:       "",
			ExpectedType: PullRequestOpenedPayload{},
			Header: map[string][]string{
				"X-Event-Key": {"pr:opened"},
			},
		},
		{
			Name:         "Valid pr:declined",
			Body:         NewPullRequestOpened(),
			ExpectedErr:  false,
			Secret:       "",
			ExpectedType: PullRequestDeclinedPayload{},
			Header: map[string][]string{
				"X-Event-Key": {"pr:declined"},
			},
		},
		{
			Name:         "invalid pr:opened, secret not set",
			Body:         NewPullRequestOpened(),
			ExpectedErr:  true,
			ExpectedType: PullRequestOpenedPayload{},
			Secret:       "",
			Header: map[string][]string{
				"X-Event-Key":     {"pr:opened"},
				"X-Hub-Signature": {"sha256=f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8"},
			},
		},
		{
			Name:         "valid pr:comment:added",
			Body:         NewPullRequestOpened(),
			ExpectedErr:  false,
			Secret:       "i am groot test",
			ExpectedType: PullRequestCommentAddedPayload{},
			Header: map[string][]string{
				"X-Event-Key": {"pr:comment:added"},
			},
		},
		{
			Name:         "valid pr:comment:deleted",
			Body:         NewPullRequestOpened(),
			ExpectedErr:  false,
			Secret:       "i am groot test",
			ExpectedType: PullRequestCommentDeletedPayload{},
			Header: map[string][]string{
				"X-Event-Key": {"pr:comment:deleted"},
			},
		},
	}

	for _, tt := range tc {
		fmt.Println("Test:", tt.Name)
		req := httptest.NewRequest(http.MethodPost, "/", tt.Body)
		req.Header = tt.Header

		w := New(WithSecret(tt.Secret), PreserveBody())

		event, err := w.Parse(req)

		if tt.Secret != w.secret {
			t.Errorf("Expected: %s, Got: %s", tt.Secret, w.secret)
		}

		if tt.ExpectedErr && err == nil || !tt.ExpectedErr && err != nil {
			t.Errorf("Expected: %v, Got: %v", tt.ExpectedErr, reflect.TypeOf(err))
			return
		}

		if !tt.ExpectedErr && reflect.TypeOf(event) != reflect.TypeOf(tt.ExpectedType) {
			t.Errorf("Expected: %v, Got: %v", reflect.TypeOf(tt.ExpectedType).Name(), reflect.TypeOf(event).Name())
		}

	}
}

func NewPullRequestOpened() io.Reader {
	jsonStr := `{
		"eventKey": "pr:opened"
		}`

	_, _ = json.Marshal(jsonStr)
	return strings.NewReader(jsonStr)
}
