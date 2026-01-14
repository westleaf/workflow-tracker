//go:build linux

package notify

import ()

type linuxNotifier struct{}

func New() Notifier { return &linuxNotifier{} }

func (n *linuxNotifier) Notify(title, message, url string) error {
	return nil
}
