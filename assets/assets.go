package assets

import (
	_ "embed"
)

//go:embed success.png
var success []byte

//go:embed fail.png
var fail []byte

//go:embed tray.ico
var trayIcon []byte

func SuccessIcon() []byte {
	return success
}

func FailIcon() []byte {
	return fail
}


