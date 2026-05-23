package main

import (
	"log"
	"os"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/westleaf/workflow-tracker/assets"
)

func runSystray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Workflow notifier")
	systray.SetTooltip("Workflow notifier")
	if icon := assets.TrayIcon(); icon != nil {
		systray.SetIcon(icon)
	}

	// Show startup notification
	beeep.AppName = "Workflow notifier"
	if err := beeep.Notify("Running!", "Watching workflows", ""); err != nil {
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
	log.Println("Exiting...")
	os.Exit(0)
}
