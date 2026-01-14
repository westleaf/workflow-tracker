//go:build darwin

package notify

import (
	"os/exec"
)

type darwinNotifier struct{}

func New() Notifier { return &darwinNotifier{} }

func (n *darwinNotifier) Notify(title, message, url string) error {
	cmd := exec.Command("terminal-notifier",
		"-title", title,
		"-message", message,
		"-open", url,
	)
	return cmd.Run()
}
