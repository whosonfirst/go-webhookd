package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	gogithub "github.com/google/go-github/github"
)

func GenerateSignature(body string, secret string) (string, error) {

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(body))

	sum := mac.Sum(nil)
	enc := hex.EncodeToString(sum)

	sig := fmt.Sprintf("sha1=%s", enc)

	return sig, nil
}

// this doesn't really work as I'd like it only returns an interface{}
// and not a typed thing that the compiler knows what to do with after
// the fact... (20161017/thisisaaronland)

func UnmarshalEvent(event_type string, body []byte) (interface{}, error) {

	var event interface{}
	ok := true

	switch event_type {
	case "commit_comment":
		event = gogithub.CommitCommentEvent{}
	case "create":
		event = gogithub.CreateEvent{}
	case "delete":
		event = gogithub.DeleteEvent{}
	case "deployment":
		event = gogithub.DeploymentEvent{}
	case "deployment_status":
		event = gogithub.DeploymentStatusEvent{}
	case "fork":
		event = gogithub.ForkEvent{}
	case "gollum":
		event = gogithub.GollumEvent{}
	case "issue_comment":
		event = gogithub.IssueCommentEvent{}
	case "issues":
		event = gogithub.IssuesEvent{}
	case "member":
		event = gogithub.MemberEvent{}
	case "membership":
		event = gogithub.MembershipEvent{}
	case "page_build":
		event = gogithub.PageBuildEvent{}
	case "public":
		event = gogithub.PublicEvent{}
	case "pull_request_review_comment":
		event = gogithub.PullRequestReviewCommentEvent{}
	// case "pull_request_review":
	// 	event = gogithub.PullRequestReviewEvent{}
	case "pull_request":
		event = gogithub.PullRequestEvent{}
	case "push":
		event = gogithub.PushEvent{}
	case "repository":
		event = gogithub.RepositoryEvent{}
	case "release":
		event = gogithub.ReleaseEvent{}
	case "status":
		event = gogithub.StatusEvent{}
	case "team_add":
		event = gogithub.TeamAddEvent{}
	case "watch":
		event = gogithub.WatchEvent{}
	default:
		ok = false
	}

	if !ok {
		return nil, errors.New("Unknown event type")
	}

	err := json.Unmarshal(body, &event)

	if err != nil {
		return nil, err
	}

	return event, nil
}
