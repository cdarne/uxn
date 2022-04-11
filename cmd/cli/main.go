package main

import (
	"log"
	"os"

	"github.com/cdarne/uxn/cpu"
)

func main() {
	rom, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	cpu := new(cpu.CPU)
	cpu.Reset()
	cpu.Load(rom, 0x100)
	cpu.Execute()
}
