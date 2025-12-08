package main

import (
	_ "embed"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/westleaf/workflow-tracker/internal/config"
)

func main() {
	config := newConfig()

	systray.Run(onReady, onExit)
	watchWorkflows(config)
}

func onReady() {
	systray.SetTitle("Workflow notifier")
	systray.SetTooltip("din din dan")

	// Show startup notification
	beeep.AppName = "Workflow notifier"
	if err := beeep.Notify("Setup!", "Watching workflows", ""); err != nil {
		log.Println("Notification error:", err)
	}

	// Add menu items
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

}

func onExit() {
	// Cleanup code goes here (close connections, save state, etc.)
	log.Println("Exiting...")
}

func checkWorkflows(c *config.Config) error {
	return nil
}
