//go:build !windows

package assets

func TrayIcon() []byte {
	return success
}
