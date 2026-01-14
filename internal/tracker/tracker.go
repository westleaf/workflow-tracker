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

	"github.com/gen2brain/beeep"
	"github.com/google/go-github/v81/github"
	"github.com/westleaf/workflow-tracker/assets"
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
			log.Printf("PR %s updated. New head SHA: %s\n", update.key, update.state.HeadSHA)
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
	key := key(owner, repo, number)
	oldState := t.state.PRState.PRs[key]
	prState := oldState

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, number)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return prUpdate{}, err
	}

	if oldState.Etag != "" {
		req.Header.Set("If-None-Match", oldState.Etag)
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

	if resp != nil && resp.StatusCode == http.StatusNotModified {
		if oldState.HeadSHA != "" {
			workflowRun, err := t.getWorkflowStatus(ctx, owner, repo, oldState.HeadSHA)
			if err != nil {
				return prUpdate{}, err
			}

			if workflowRun.GetID() != 0 {
				prState.WorkflowStatus = workflowRun.GetStatus()
				prState.WorkflowConclusion = workflowRun.GetConclusion()
				prState.WorkflowRunID = int(workflowRun.GetID())
			}

			checkShouldNotify(oldState, prState)

			if oldState.WorkflowStatus != prState.WorkflowStatus {
				return prUpdate{key: key, state: prState, modified: true}, nil
			}
		}
		return prUpdate{modified: false}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return prUpdate{}, err
	}

	var pr github.PullRequest
	if err := json.Unmarshal(body, &pr); err != nil {
		return prUpdate{}, err
	}

	newHeadSha := pr.GetHead().GetSHA()

	if oldState.HeadSHA != "" && oldState.HeadSHA != newHeadSha {
		prState.WorkflowStatus = "pending"
		prState.WorkflowRunID = 0
		oldState.WorkflowStatus = "pending"
		log.Printf("new commit on PR %s: %s -> %s", key, oldState.HeadSHA, newHeadSha)
	}

	workflowRun, err := t.getWorkflowStatus(ctx, owner, repo, newHeadSha)
	if err != nil {
		return prUpdate{}, err
	}

	if workflowRun.GetID() != 0 {
		prState.WorkflowStatus = workflowRun.GetStatus()
		prState.WorkflowRunID = int(workflowRun.GetID())
		prState.WorkflowConclusion = workflowRun.GetConclusion()
	} else {
		prState.WorkflowStatus = "not_found"
	}

	checkShouldNotify(oldState, prState)

	prState.Number = number
	prState.Repo = repo
	prState.HeadSHA = pr.GetHead().GetSHA()
	prState.Etag = resp.Header.Get("ETag")
	prState.UpdatedAt = time.Now()

	return prUpdate{
		key:      key,
		state:    prState,
		modified: true,
		err:      nil,
	}, nil
}

func checkShouldNotify(oldState, newState config.PRState) {
	if oldState.WorkflowStatus != "completed" && newState.WorkflowStatus == "completed" {
		if newState.WorkflowConclusion == "success" || newState.WorkflowConclusion == "failure" {
			title := fmt.Sprintf("Pr #%d Workflow Complete", newState.Number)
			message := fmt.Sprintf("%s - %s", newState.Repo, newState.WorkflowConclusion)
			if newState.WorkflowConclusion == "success" {
				_ = beeep.Notify(title, message, assets.SuccessIcon())
			} else {
				_ = beeep.Notify(title, message, assets.FailIcon())
			}
		}
	}
}

func (t *Tracker) getWorkflowStatus(ctx context.Context, owner, repo, sha string) (workflowRun github.WorkflowRun, err error) {
	opts := &github.ListWorkflowRunsOptions{
		HeadSHA: sha,
	}

	runs, _, err := t.state.Client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	if err != nil {
		return github.WorkflowRun{}, err
	}

	if len(runs.WorkflowRuns) == 0 {
		return github.WorkflowRun{}, nil
	}

	return *runs.WorkflowRuns[0], nil
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
