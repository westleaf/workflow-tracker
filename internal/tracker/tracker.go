package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/westleaf/workflow-tracker/internal/config"
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
	if t.state.PRState.PRs == nil {
		t.state.PRState.PRs = make(map[string]config.PRState)
	}

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

type prUpdate struct {
	key      string
	state    config.PRState
	modified bool
	err      error
}

func (t *Tracker) checkPRs(ctx context.Context) error {
	q := fmt.Sprintf("is:pr author:%s state:open", t.state.Config.CurrentUser)
	res, _, err := t.state.Client.Search.Issues(ctx, q, &github.SearchOptions{})
	if err != nil {
		return err
	}

	updatesChan := make(chan prUpdate, len(res.Issues))

	sem := make(chan struct{}, 5)

	for _, pr := range res.Issues {
		sem <- struct{}{}

		go func(pr *github.Issue) {
			defer func() { <-sem }()

			owner, repo := ExtractRepo(pr)
			number := pr.GetNumber()

			update, err := t.fetchPRDetails(ctx, owner, repo, number)
			if err != nil {
				updatesChan <- prUpdate{err: err}
				return
			}

			updatesChan <- update
		}(pr)
	}

	for i := 0; i < len(res.Issues); i++ {
		update := <-updatesChan
		if update.err != nil {
			return update.err
		}
		if update.modified {
			fmt.Printf("PR %s updated. New head SHA: %s\n", update.key, update.state.HeadSHA)
			t.state.PRState.PRs[update.key] = update.state
		}
	}

	err = config.WriteState(*t.state.PRState)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tracker) fetchPRDetails(ctx context.Context, owner, repo string, number int) (prUpdate, error) {
	prState := t.state.PRState.PRs[key(owner, repo, number)]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return prUpdate{}, err
	}

	if prState.Etag != "" {
		req.Header.Set("If-None-Match", prState.Etag)
	}

	httpClient := t.state.Client.Client()
	resp, err := httpClient.Do(req)
	if err != nil {
		return prUpdate{}, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// pr, resp, err := t.state.Client.PullRequests.Get(ctx, owner, repo, number)
	// if err != nil {
	if resp != nil && resp.StatusCode == http.StatusNotModified {
		log.Printf("no change on pr %s", url)
		return prUpdate{modified: false}, nil
	}
	// return prUpdate{}, err
	// }

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return prUpdate{}, err
	}

	var pr github.PullRequest
	if err := json.Unmarshal(body, &pr); err != nil {
		return prUpdate{}, err
	}

	prState.Number = number
	prState.Repo = repo
	prState.HeadSHA = pr.GetHead().GetSHA()
	prState.Etag = resp.Header.Get("ETag")
	prState.UpdatedAt = time.Now()

	t.state.PRState.PRs[key(owner, repo, number)] = prState

	status, conclusion, err := t.getWorkflowStatus(ctx, owner, repo, prState.HeadSHA)
	if err != nil {
		return prUpdate{}, err
	}

	log.Printf("%s, %s", status, conclusion)

	return prUpdate{
		key:      key(owner, repo, number),
		state:    prState,
		modified: true,
		err:      nil,
	}, nil
}

func (t *Tracker) getWorkflowStatus(ctx context.Context, owner, repo, headSHA string) (status, conclusion string, err error) {
	opts := &github.ListWorkflowRunsOptions{
		HeadSHA: headSHA,
	}

	runs, _, err := t.state.Client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	if err != nil {
		return "", "", err
	}

	if len(runs.WorkflowRuns) == 0 {
		log.Printf("no workflow runs found for sha %s", headSHA)
		return
	}
	return *runs.WorkflowRuns[0].Status, "", nil
}

func key(owner, repo string, number int) string {
	return fmt.Sprintf("%s/%s+#%d", owner, repo, number)
}

func ExtractRepo(issue *github.Issue) (owner, repo string) {
	repoURL := issue.GetRepositoryURL()
	parts := strings.Split(repoURL, "/")
	owner = parts[len(parts)-2]
	repo = parts[len(parts)-1]

	return owner, repo
}
