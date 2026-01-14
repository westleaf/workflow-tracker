package auth

import (
	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
)

func AuthHandler() (*api.AccessToken, error) {
	host, err := oauth.NewGitHubHost("https://github.com")
	if err != nil {
		return nil, err
	}

	flow := &oauth.Flow{
		Host:     host,
		ClientID: "Ov23lizLQqTQJR22bGlf",
		Scopes:   []string{"repo", "read:user", "workflow"},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return nil, err
	}

	return accessToken, err
}
