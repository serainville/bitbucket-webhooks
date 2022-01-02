package bitbucket

// EventKey stores the key for an event received by Bitbucket
type EventKey string

// EventDate stores the date an event was trigger by Bitbucket
type EventDate string

type commonBitbucketEventFields struct {
	// EventKey is the event key of a Bitbucket Webhook
	EventKey string `json:"eventKey"`
	// EventDate is the date the event occurred
	EventDate string `json:"date"`
}

// DiagnosticPingEvent maps to "diagnostic:ping" Bitbucekt webhook events
type DiagnosticPingEvent struct {
	Test bool `json:"test"`
}

// PullRequestOpenedPayload maps to "pr:opened" Bitbucket Webook events
type PullRequestOpenedPayload struct {
	commonBitbucketEventFields
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

// PullRequestModifiedPayload maps to "pr:modified" Bitbucket Webhook events
type PullRequestModifiedPayload struct {
	commonBitbucketEventFields
	Actor               `json:"actor"`
	PullRequest         `json:"pullRequest"`
	PreviousTitle       string `json:"previousTitle"`
	PreviousDescription string `json:"previousDescription"`
	PreviousTarget      `json:"previousTarget"`
}

// PullRequestDeletedPayload maps to "pr:deleted" Bitbucket Webhook events
type PullRequestDeletedPayload struct {
	commonBitbucketEventFields
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

// PullRequestMergedPayload maps to "pr:merged" Bitbucket Webhook events
type PullRequestMergedPayload struct {
	commonBitbucketEventFields
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

// PullRequestDeclinedPayload maps to 'pr:declined' Bitbucket Webhook events
type PullRequestDeclinedPayload struct {
	commonBitbucketEventFields
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

// PullRequestReviewerPayload maps to "pr:reviewer:approved", "pr:reviewer:needs_work", and "pr:reviewer:unapproved" Bitbucket events
type PullRequestReviewerPayload struct {
	commonBitbucketEventFields
	Actor          `json:"actor"`
	PullRequest    `json:"pullRequest"`
	Participant    `json:"participant"`
	PreviousStatus string `json:"previousStatus"`
}

// PullRequestReviewerUpdatedPayload maps to a "pr:reviewer:updated" Bitbucket event
type PullRequestReviewerUpdatedPayload struct {
	commonBitbucketEventFields
	Actor            `json:"actor"`
	PullRequest      `json:"pullRequest"`
	AddedReviewers   []Actor `json:"addedReviewers"`
	RemovedReviewers []Actor `json:"removedReviewers"`
}

// PullRequestCommentAddedPayload maps to a 'pr:comment:added' Bitbucket webhook event
type PullRequestCommentAddedPayload struct {
	commonBitbucketEventFields
	Actor           `json:"actor"`
	PullRequest     `json:"pullRequest"`
	Comment         `json:"comment"`
	CommentParentID uint `json:"commentParentId"`
}

// PullRequestCommentEditedPayload maps to a 'pr:comment:edited' Bitbucket webhook event
type PullRequestCommentEditedPayload struct {
	commonBitbucketEventFields

	// Actor is the user that edited the comment
	Actor `json:"actor"`

	// PullRequest is the pull request where the comment exists
	PullRequest `json:"pullRequest"`

	// Comment is the comment edited
	Comment `json:"comment"`

	// CommentParentID is the ID of the parent comment if one exists.
	CommentParentID uint `json:"commentParentId"`

	// PreviousComment is the text of the previous comment.
	PreviousComment string `json:"previousComment"`
}

// PullRequestCommentDeletedPayload maps to a 'pr:comment:deleted' Bitbucket webhook event
type PullRequestCommentDeletedPayload struct {
	commonBitbucketEventFields
	Actor           `json:"actor"`
	PullRequest     `json:"pullRequest"`
	Comment         `json:"comment"`
	CommentParentID uint `json:"commentParentId"`
}

// RepoRefsChangedPayload maps to 'repo:refs_changed' Bitbucket Webhook events
type RepoRefsChangedPayload struct {
	commonBitbucketEventFields
	Actor      `json:"actor"`
	Repository `json:"repository"`
	Changes    []Changes `json:"changes"`
}

// RepoModifiedPayload maps to 'repo:modified' Bitbucket Webhook events
type RepoModifiedPayload struct {
	commonBitbucketEventFields
	Actor      `json:"actor"`
	OldVersion RepoVersion `json:"old"`
	NewVersion RepoVersion `json:"new"`
}

// FromRefUpdatedPayload maps to 'from_ref_updated' Bitbucket Webhook events
type FromRefUpdatedPayload struct {
	commonBitbucketEventFields
	Actor            `json:"actor"`
	PullRequest      `json:"pullRequest"`
	PreviousFromHash string `json:"previousFramHash"`
}

// Actor represents the actor field of a Bitbucket Webhook request
type Actor struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           uint64 `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

// PullRequest represents the pullRequest field of a Bitbucket Webhook request
type PullRequest struct {
	ID          uint64 `json:"id"`
	Version     uint64 `json:"version"`
	Title       string `json:"title"`
	State       string `json:"state"`
	Open        bool   `json:"open"`
	Closed      bool   `json:"closed"`
	CreatedDate uint64 `json:"createdDate"`
	UpdatedDate uint64 `json:"updatedDate"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
}

// Ref represents the fromRef field of a Bitbucket Webhook request
type Ref struct {
	ID           string `json:"id"`
	DisplayID    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   `json:"repository"`
}

// Repository maps to the repository key from a Bitbucket event
type Repository struct {
	Slug          string `json:"slug"`
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       `json:"project"`
	Public        bool `json:"public"`
}

// Project maps to the project key from a Bitbucket event
type Project struct {
	Key    string `json:"key"`
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
}

// Changes maps to the changes key from a Bitbucket event
type Changes struct {
	Ref struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
		Type      string `json:"type"`
	} `json:"ref"`
	RefID  string `json:"refId"`
	ToHash string `json:"toHash"`
	Type   string `json:"type"`
}

// RepoVersion maps to the version key of a Bitbucket event
type RepoVersion struct {
	Slug          string `json:"slug"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       `json:"project"`
	Public        bool `json:"public"`
}

// Participant maps to the participant key of a Bitbucket event
type Participant struct {
	Actor              `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Approved           string `json:"approved"`
	Status             string `json:"status"`
}

// PreviousTarget maps to the previousTarget key of a Bitbucket event
type PreviousTarget struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
}

// Comment maps to the `comment` key of a Bitbucket event
type Comment struct {
	Properties struct {
		RepositoryID uint `json:"repositoryId"`
	} `json:"properties"`
	ID          uint   `json:"id"`
	Version     uint   `json:"version"`
	Text        string `json:"text"`
	Actor       `json:"author"`
	CreatedDate uint                     `json:"createdDate"`
	UpdatedDate uint                     `json:"updatedDate"`
	Comments    []Comment                `json:"comments"`
	Tasks       []map[string]interface{} `json:"tasks"`
}
