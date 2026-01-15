package main

import (
	"context"
	_ "embed"
	"log"

	"github.com/google/go-github/v81/github"
	"github.com/joho/godotenv"

	"github.com/westleaf/workflow-tracker/internal/auth"
	"github.com/westleaf/workflow-tracker/internal/config"
	"github.com/westleaf/workflow-tracker/internal/runtime"
	"github.com/westleaf/workflow-tracker/internal/tracker"
)

func main() {
	godotenv.Load()

	if err := config.EnsureConfigExists(); err != nil {
		log.Fatal(err)
	}

	if err := config.EnsureStateExists(); err != nil {
		log.Fatal(err)
	}

	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if conf.Token == "" {
		token, err := auth.AuthHandler()
		if err != nil {
			log.Fatal(err)
		}

		conf.SetToken(token.Token)
	}

	ghclient := github.NewClient(nil).WithAuthToken(conf.Token)

	user, _, err := ghclient.Users.Get(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	if conf.CurrentUser == "" {
		conf.SetUser(user.GetName())
	}

	prState, err := config.ReadState()
	if err != nil {
		log.Fatal(err)
	}

	st := runtime.State{
		Config:  &conf,
		Client:  ghclient,
		PRState: &prState,
	}

	tracker := tracker.NewTracker(&st)

	go func() {
		err = tracker.Start("2m")
		if err != nil {
			log.Printf("tracker error: %v", err)
		}
	}()

	runSystray()
}
