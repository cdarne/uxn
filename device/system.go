package device

type System struct {
	Halt   bool
	Status bool
}

func (d *System) Name() string {
	return "System"
}

func (d *System) Output(port byte, value byte) {
	switch port {
	case 0xe:
		d.Status = true
	case 0xf:
		d.Halt = true
	}
}
