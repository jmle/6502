package main

import (
	"testing"
	"log"
)

type Memory struct {
	memory	[1<<8]int
}

func (m *Memory) Read(addr int) int {
	return m.memory[addr]
}

func (m *Memory) Write(addr, value int) {
	m.memory[addr] = value
}

// interprets a word as bcd
func bcd(n int) int {
	return (n & 0xF) + (n&0xF0>>4 * 10)
}

func TestAdcHappyPath(t *testing.T) {
	log.Println("Test adc happy path")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 2
	mem.Write(0, 3)

	cpu.adc(0)

	if expect := 5; cpu.ac != expect {
		t.Errorf("Invalid ac value: %q != %q", cpu.ac, expect)
	}
	if cpu.p.v != 0 || cpu.p.n != 0 || cpu.p.z != 0 || cpu.p.c != 0 {
		t.Errorf("Wrong processor status")
	}
}

func TestAdcWithOverflow(t *testing.T) {
	log.Println("Test adc with overflow")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 20
	mem.Write(0, 120)

	cpu.adc(0)

	if cpu.p.v != 1 {
		t.Errorf("Overflow flag clear")
	}
}

func TestAdcWithNegative(t *testing.T) {
	log.Println("Test adc with negative")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.mem.Write(0, 200)

	cpu.adc(0)

	if cpu.p.n != 1 {
		t.Errorf("Negative flag clear")
	}
}

func TestAdcWithZero(t *testing.T) {
	log.Println("Test adc with zero")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	mem.Write(0, 0)

	cpu.adc(0)

	if cpu.p.z != 1 {
		t.Errorf("Zero flag clear")
	}
}

func TestAdcWithDecimalModeWithCarry(t *testing.T) {
	log.Println("Test adc with decimal")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.d = 1			// set bcd mode
	cpu.ac = 1
	mem.Write(0, 153)   // 99 in bcd is 153

	cpu.adc(0)

	if cpu.p.c != 1 {
		t.Errorf("Carry flag clear")
	}
	if expect := 0; cpu.ac != expect {
		t.Errorf("Invalid result: %b != %d", cpu.ac, expect)
	}
}

func TestAdcWithDecimalWithoutCarry(t *testing.T) {
	log.Println("Test adc with decimal")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.d = 1			// set bcd mode
	cpu.ac = 64			// bcd: 40
	mem.Write(0, 8)		// bcd: 8

	cpu.adc(0)

	if cpu.p.c != 0 {
		t.Errorf("Carry flag set")
	}
	if expect := 72; cpu.ac != expect {	// bcd: 48
		t.Errorf("Invalid result: %b != %d", cpu.ac, expect)
	}
}

func TestAdcWithCarry(t *testing.T) {
	log.Println("Test adc with decimal")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 255
	mem.Write(0, 1)

	cpu.adc(0)

	if cpu.p.c != 1 {
		t.Errorf("Carry flag clear")
	}
	if expect := 0; cpu.ac != expect {
		t.Errorf("Invalid result: %b != %d", cpu.ac, expect)
	}
}

func TestAnd(t *testing.T) {
	log.Println("Test and")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 240
	mem.Write(0, 255)

	cpu.and(0)

	if expect := 240; cpu.ac != expect {
		t.Errorf("Invalid result: %b != %b", cpu.ac, expect)
	}
	if cpu.p.n != 1 || cpu.p.z != 0 {
		t.Errorf("Invalid processor status")
	}
}

func TestAslaWithoutCarry(t *testing.T) {
	log.Println("Test asla without carry")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 16

	cpu.asla()

	if expect := 32; cpu.ac != expect {
		t.Errorf("Invalid result: %b != %b", cpu.ac, expect)
	}
	if cpu.p.n != 0 || cpu.p.z != 0 || cpu.p.c != 0 {
		t.Errorf("Invalid processor status")
	}
}

func TestAslaWithCarry(t *testing.T) {
	log.Println("Test asla with carry")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 128

	cpu.asla()

	if expect := 0; cpu.ac != 0 {
		t.Errorf("Invalid result: %b != %b", cpu.ac, expect)
	}
	if cpu.p.n != 0 || cpu.p.z != 1 || cpu.p.c != 1 {
		t.Errorf("Invalid processor status")
	}
}

func TestAslWithoutCarry(t *testing.T) {
	log.Println("Test asl without carry")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.mem.Write(0, 16)

	cpu.asl(0)

	actual := cpu.mem.Read(0)
	if expect := 32; actual != expect {
		t.Errorf("Invalid result: %b != %b", actual, expect)
	}
	if cpu.p.c != 0 || cpu.p.n != 0 || cpu.p.z != 0 {
		t.Errorf("Invalid processor status")
	}
}

func TestAslWithCarry(t *testing.T) {
	log.Println("Test asl with carry")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.mem.Write(0, 128)

	cpu.asl(0)

	actual := cpu.mem.Read(0)
	if expect := 0; actual != expect {
		t.Errorf("Invalid result: %b != %b", actual, expect)
	}
	if cpu.p.c != 1 || cpu.p.n != 0 || cpu.p.z != 1 {
		t.Errorf("Invalid processor status")
	}
}

func TestBccWithCarryClear(t *testing.T) {
	log.Println("Test bcc with carry clear")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.c = 0

	actual := cpu.bcc(16)

	if expect := true; actual != expect {
		t.Errorf("Invalid return value: %v != %v", actual, expect)
	}
	if expect := 16; cpu.pc != expect {
		t.Errorf("Wrong PC")
	}
}

func TestBccWithCarrySet(t *testing.T) {
	log.Println("Test bcc with carry set")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.c = 1

	actual := cpu.bcc(16)

	if expect := false; actual != expect {
		t.Errorf("Invalid return value: %v != %v", actual, expect)
	}
}
