package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cdarne/uxn/cpu"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cpu *cpu.CPU
}

func initialModel() model {
	rom, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	c := new(cpu.CPU)
	c.Reset()
	c.Load(rom, 0x100)

	return model{cpu: c}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// The "down" and "j" keys move the cursor down
		case " ":
			m.cpu.Step()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := fmt.Sprintf("CPU - PC: $%#x\n\n", m.cpu.PC)

	s += "Devices:\n"
	for i, d := range m.cpu.Devices {
		s += fmt.Sprintf("[%x] %s", i<<8, d.Name())
		s += "\n"
	}

	/*
		s += "\nMemory:\n"
		addr := 0
		for i := 0; i < 256; i++ {
			s += fmt.Sprintf("$%04x ", addr)
			for j := 0; j < 8; j++ {
				s += fmt.Sprintf(" %02x", m.cpu.RAM[addr+j])
			}
			s += " "
			for j := 8; j < 16; j++ {
				s += fmt.Sprintf(" %02x", m.cpu.RAM[addr+j])
			}
			s += "\n"
			addr += 0x10
		}
	*/

	s += "\nWorking Stack:\n"
	s += printStack(m.cpu.WS)
	/*
		s += "\nCurrent state:\n"
		start := max(0, m.cpu.PC-5)
		end := min(0xffff, m.cpu.PC+5)
		for p := start; start < end; p++ {
			s += fmt.Sprintf("$%02x:", p)
			b := m.cpu.RAM[m.cpu.PC]
			inst := instructions[b&0x1f]
			s += fmt.Sprintf(" %s (%02x)\n", inst.name, b)
		}
	*/
	// The footer
	s += "\nPress space to step, r to reset, q to quit.\n"

	// Send the UI for rendering
	return s
}

type instruction struct {
	name  string
	arity uint
}

var instructions = []instruction{
	{"LIT", 1},
	{"INC", 0},
	{"POP", 0},
	{"DUP", 0},
	{"NIP", 0},
	{"SWP", 0},
	{"OVR", 0},
	{"ROT", 0},

	{"EQU", 0},
	{"NEQ", 0},
	{"GTH", 0},
	{"LTH", 0},
	{"JMP", 1},
	{"JCN", 1},
	{"JSR", 1},
	{"STH", 1},

	{"LDZ", 1},
	{"STZ", 1},
	{"LDR", 1},
	{"STR", 1},
	{"LDA", 1},
	{"STA", 1},
	{"DEI", 1},
	{"DEO", 1},

	{"ADD", 1},
	{"SUB", 1},
	{"MUL", 1},
	{"DIV", 1},
	{"AND", 1},
	{"ORA", 1},
	{"EOR", 1},
	{"SFT", 1},
}

func printStack(s cpu.Stack) (ret string) {
	for p := byte(0); p < s.Ptr; p++ {
		ret += fmt.Sprintf(" - %02x\n", s.Data[p])
	}
	return ret
}

func min(a cpu.Word, b cpu.Word) cpu.Word {
	if a > b {
		return b
	}
	return a
}

func max(a cpu.Word, b cpu.Word) cpu.Word {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
