package cpu

import (
	"log"

	"github.com/cdarne/uxn/device"
)

type Word uint16
type Stack struct {
	Ptr  byte
	Data []byte
}

func newStack() Stack {
	return Stack{Data: make([]byte, 256)}
}

func (s *Stack) Push(b byte) {
	if s.Ptr == 0xff {
		panic("stack overflow")
	}
	s.Data[s.Ptr] = b
	s.Ptr++
}

func (s *Stack) Pop() byte {
	if s.Ptr == 0x00 {
		panic("stack underflow")
	}
	s.Ptr--
	return s.Data[s.Ptr]
}

func (s *Stack) PopWord() Word {
	if s.Ptr <= 0x00 {
		panic("stack underflow")
	}
	l := Word(s.Data[s.Ptr-1])
	h := Word(s.Data[s.Ptr-2])
	s.Ptr -= 2
	return l + (h << 8)
}

func (s *Stack) PushWord(w Word) {
	if s.Ptr >= 0xfe {
		panic("stack overflow")
	}
	s.Data[s.Ptr] = byte((w & 0xff00) >> 8)
	s.Data[s.Ptr+1] = byte(w & 0xff)
	s.Ptr += 2
}

type CPU struct {
	PC Word // program counter

	WS Stack // working stack
	RS Stack // return stack

	RAM []byte // random access memory

	Devices []device.Device
}

func (u *CPU) Reset() {
	u.PC = 0x100
	u.WS = newStack()
	u.RS = newStack()
	u.RAM = make([]byte, 64*1024) // 64KB
	u.Devices = make([]device.Device, 16)
	u.Devices[0] = new(device.System)
	u.Devices[1] = new(device.Console)
	for i := 2; i < 16; i++ {
		u.Devices[i] = new(device.Null)
	}
}

func (u *CPU) FetchByte() (b byte) {
	b = u.RAM[u.PC]
	u.PC++
	return b
}

func (u *CPU) Execute() {
	system := u.Devices[0].(*device.System)
	for !system.Halt && u.Step() {
		if system.Status {
			system.Status = false
			log.Printf("CPU: %+v\n", u)
		}
	}
}

func (u *CPU) Step() bool {
	op := u.FetchByte()
	if op == 0 {
		return false
	}
	shortMode := op&0x20 == 0x20
	switch op & 0x1f { // first 5 bits describe the opcode

	case 0x00: // LIT - Pushes the next value seen in the program onto the stack.
		if shortMode {
			a := u.FetchByte()
			u.WS.Push(a)
			a = u.FetchByte()
			u.WS.Push(a)
		} else {
			a := u.FetchByte()
			u.WS.Push(a)
		}
	case 0x01: // INC
		if shortMode {
			w := u.WS.PopWord()
			u.WS.PushWord(w + 1)
		} else {
			a := u.WS.Pop()
			u.WS.Push(a + 1)
		}
	case 0x17: // DEO
		b := u.WS.Pop()
		a := u.WS.Pop()
		deviceId := (b & 0xf0) >> 4 // mask High nibble of byte then shift right
		port := b & 0x0f            // mask Low nibble of byte
		u.Devices[deviceId].Output(port, a)
	case ADD: // Pushes the sum of the two values at the top of the stack.
		b := u.WS.Pop()
		a := u.WS.Pop()
		u.WS.Push(a + b)
	case SUB: // Pushes the difference of the first value minus the second, to the top of the stack.
		b := u.WS.Pop()
		a := u.WS.Pop()
		u.WS.Push(a - b)
	default:
		log.Fatalf("unsupported opcode: %x", op)
	}
	return true
}

func (u *CPU) Load(data []byte, addr Word) {
	copy(u.RAM[addr:], data)
}

const (
	// Opcodes

	// stack
	BRK byte = iota
	INC      // increment
	POP      // pop stack
	DUP      // duplicate
	NIP      // nip
	SWP      // swap
	OVR      // over
	ROT      // rotate

	// logic
	EQU
	NEQ
	GTH
	LTH
	JMP
	JCN
	JSR
	STH

	// memory
	LDZ
	STZ
	LDR
	STR
	LDA
	STA
	DEI
	DEO

	// arithmetic
	ADD
	SUB
	MUL
	DIV
	AND
	ORA
	EOR
	SFT

	LIT  byte = 0x80
	LIT2 byte = 0xA0
)
