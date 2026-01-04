package main

import (
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/joho/godotenv"

	"github.com/westleaf/workflow-tracker/internal/config"
	"github.com/westleaf/workflow-tracker/internal/runtime"
	"github.com/westleaf/workflow-tracker/internal/tracker"
)

func main() {
	godotenv.Load()

	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ghclient := github.NewTokenClient(context.Background(), os.Getenv("GH_TOKEN"))

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

	runSystray()

	err = tracker.Start("10s")
	if err != nil {
		log.Fatal(err)
	}
}
