package notify

import (
	"os/exec"
)

type Notifier interface {
	Notify(title, message, url string) error
	// NotifySuccess(title, message, url string) error
	// NotifyFailure(title, message, url string) error
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
