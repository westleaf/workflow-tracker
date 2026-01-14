package notify

import (
	"github.com/go-toast/toast"
)

type windowsNotifier struct{}

func New() Notifier { return &windowsNotifier{} }

func (n *windowsNotifier) Notify(title, message, url string) error {
	notification := toast.Notification{
		AppID:   "wft",
		Title:   title,
		Message: message,
		Actions: []toast.Action{
			{Type: "protocol", Label: "Open", Arguments: url},
		},
	}
	return notification.Push()
}
