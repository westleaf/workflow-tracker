package tracker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/westleaf/workflow-tracker/internal/runtime"
)

type Tracker struct {
	state *runtime.State
}

func NewTracker(state *runtime.State) *Tracker {
	return &Tracker{
		state: state,
	}
}

func (t *Tracker) Start(interval string) error {
	timeBetweenReqs, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		err := t.checkPRs(context.Background())
		if err != nil {
			return err
		}
	}
}

func (t *Tracker) checkPRs(ctx context.Context) error {
	q := fmt.Sprintf("is:pr author:%s state:open", t.state.Config.CurrentUser)
	res, _, err := t.state.Client.Search.Issues(ctx, q, &github.SearchOptions{})
	if err != nil {
		return err
	}

	for _, pr := range res.Issues {
		owner, repo := ExtractRepo(pr)
		fmt.Printf("%s, %s, %s\n", *pr.Title, owner, repo)
	}

	return nil
}

func ExtractRepo(issue *github.Issue) (owner, repo string) {
	repoURL := issue.GetRepositoryURL()
	parts := strings.Split(repoURL, "/")
	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]

	return owner, repo
}
