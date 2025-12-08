package config

import (
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	_ "github.com/bartventer/httpcache/store/memcache"
)

type Config struct {
	Client         *github.Client
	Token          *oauth2.Token
	UpdateInterval time.Duration
}

var oauth2Config = &oauth2.Config{
	ClientID:     "Iv1.0123456789abcdef",
	ClientSecret: "0123456789abcdef0123456789abcdef01234567",
	Scopes:       []string{"read:packages", "repo"},
}

func NewConfig(*Config, error) {
	c := &Config{}
	// token, err := loadTokenFromFile("token.json")
}

func loadTokenFromFile(path string) (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}
