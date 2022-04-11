build:
	go build ./cmd/emu

rom:
	uxnasm programs/hello_char.tal test.rom

dump:
	hexdump -C test.rom