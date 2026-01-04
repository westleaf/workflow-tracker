package runtime

import (
	"github.com/google/go-github/v62/github"
	"github.com/westleaf/workflow-tracker/internal/config"
)

type State struct {
	Config  *config.Config
	Client  *github.Client
	PRState *config.State
}
