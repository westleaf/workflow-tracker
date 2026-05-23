package notify

import (
	"os"
	"path/filepath"

	"github.com/go-toast/toast"
)

type windowsNotifier struct{}

func New() Notifier { return &windowsNotifier{} }

func (n *windowsNotifier) Notify(title, message string, icon []byte, url string) error {
	notification := toast.Notification{
		AppID:   "wft",
		Title:   title,
		Message: message,
		Actions: []toast.Action{
			{Type: "protocol", Label: "Open", Arguments: url},
		},
	}
	if len(icon) > 0 {
		iconPath := filepath.Join(os.TempDir(), "wft-icon.png")
		os.WriteFile(iconPath, icon, 0644)
		notification.Icon = iconPath
	}
	return notification.Push()
}
