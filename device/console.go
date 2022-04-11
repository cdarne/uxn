package device

import (
	"io"
	"os"
)

type Console struct {
}

func (d *Console) Name() string {
	return "Console"
}

func (d *Console) Output(port byte, value byte) {
	w := io.Discard
	if port == 0x08 {
		w = os.Stdout
	} else if port == 0x09 {
		w = os.Stdout
	}
	w.Write([]byte{value})
}
