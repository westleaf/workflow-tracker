//go:build linux

package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

type linuxNotifier struct{}

func New() Notifier { return &linuxNotifier{} }

func (n *linuxNotifier) Notify(title, message string, icon []byte, url string) error {
	msg := message
	if url != "" {
		msg = fmt.Sprintf("%s\n%s", message, url)
	}
	return beeep.Notify(title, msg, string(icon))
}
