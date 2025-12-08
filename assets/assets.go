package assets

import (
	_ "embed"
)

//go:embed success.png
var success []byte

//go:embed fail.png
var fail []byte

func SuccessIcon() []byte {
	return success
}

func FailIcon() []byte {
	return fail
}
