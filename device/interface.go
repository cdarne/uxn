package device

type Device interface {
	Name() string
	Output(port byte, value byte)
}

type Null struct {
}

func (d *Null) Name() string {
	return "Null"
}

func (d *Null) Output(byte, byte) {
}
