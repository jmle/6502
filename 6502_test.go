package main

import (
	"log"
	"reflect"
	"testing"
)

type Memory struct {
	memory [1 << 8]int
}

func (m *Memory) Read(addr int) int {
	return m.memory[addr]
}

func (m *Memory) Write(addr, value int) {
	m.memory[addr] = value
}

// interprets a word as bcd
func bcd(n int) int {
	return (n & 0xF) + (n & 0xF0 >> 4 * 10)
}

func TestAdc(t *testing.T) {
	for _, tt := range []struct {
		name    string
		ac, val int
		adcAddr int
		expProc ProcStat
	}{
		{name: "Happy Path",
			ac: 2, val: 3, adcAddr: 0,
			expProc: ProcStat{0, 0, 0, 0, 0, 0, 0},
		},
		{name: "With Overflow",
			ac: 20, val: 120, adcAddr: 0,
			expProc: ProcStat{0, 0, 0, 0, 0, 1, 1},
		},
	} {
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			ac:  tt.ac,
		}
		log.Println(tt.name)
		mem.Write(0, tt.val)
		cpu.adc(tt.adcAddr)
		if !reflect.DeepEqual(cpu.p, tt.expProc) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestAdcHappyPath(t *testing.T) {
	log.Println("Test adc happy path")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.ac = 2
	mem.Write(0, 3)

	cpu.adc(0)

	if cpu.ac != 5 {
		t.Errorf("Invalid ac value: %q != %q", cpu.ac, 5)
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
	cpu.p.d = 1 // set bcd mode
	cpu.ac = 1
	mem.Write(0, 153) // 99 in bcd is 153

	cpu.adc(0)

	if cpu.p.c != 1 {
		t.Errorf("Carry flag clear")
	}
	if cpu.ac != 0 {
		t.Errorf("Invalid result: %b != %d", cpu.ac, 0)
	}
}

func TestAdcWithDecimalWithoutCarry(t *testing.T) {
	log.Println("Test adc with decimal")

	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.d = 1     // set bcd mode
	cpu.ac = 64     // bcd: 40
	mem.Write(0, 8) // bcd: 8

	cpu.adc(0)

	if cpu.p.c != 0 {
		t.Errorf("Carry flag set")
	}
	if cpu.ac != 72 { // bcd: 48
		t.Errorf("Invalid result: %b != %d", cpu.ac, 72)
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
	if cpu.ac != 0 {
		t.Errorf("Invalid result: %b != %d", cpu.ac, 0)
	}
}
