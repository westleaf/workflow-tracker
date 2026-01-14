package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/google/go-github/v81/github"
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

	ghclient := github.NewClient(nil).WithAuthToken(os.Getenv("GH_TOKEN"))

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
