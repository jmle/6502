package main

import (
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

func TestAdc(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		ac, val int
		adc     int
		proc    ProcStat
		// Expected
		expAc   int
		expProc ProcStat
	}{
		{name: "Happy Path",
			ac: 2, val: 3, adc: 0,
			expAc: 5,
		},
		{name: "With overflow",
			ac: 20, val: 120, adc: 0,
			expAc:   140,
			expProc: ProcStat{n: 1, v: 1},
		},
		{name: "With negative",
			val:     200,
			expAc:   200,
			expProc: ProcStat{n: 1, v: 1},
		},
		{name: "With zero",
			expProc: ProcStat{z: 1},
		},
		{name: "With carry",
			ac: 255, val: 1,
			expProc: ProcStat{c: 1, v: 1},
		},
		{name: "Decimal mode, without Carry",
			ac: 64, val: 8,
			proc:    ProcStat{d: 1},
			expProc: ProcStat{d: 1},
			expAc:   72,
		},
		{name: "Decimal mode, with carry",
			ac: 1, val: 153,
			proc:    ProcStat{d: 1},
			expProc: ProcStat{c: 1, d: 1},
		},
	} {
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			ac:  tt.ac,
			p:   tt.proc,
		}
		mem.Write(0, tt.val)
		cpu.adc(tt.adc)
		t.Log(tt.name)
		if !reflect.DeepEqual(cpu.p, tt.expProc) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
		if cpu.ac != tt.expAc {
			t.Errorf("Expected ac %d, got %d\n", tt.expAc, cpu.ac)
		}
	}
}

func TestAdcWithCarry(t *testing.T) {
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
	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.c = 1

	actual := cpu.bcc(16)

	if expect := false; actual != expect {
		t.Errorf("Invalid return value: %v != %v", actual, expect)
	}
}

func TestBcsWithCarryClear(t *testing.T) {
	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.c = 0

	actual := cpu.bcs(16)

	if expect := false; actual != expect {
		t.Errorf("Invalid return value: %v != %v", actual, expect)
	}
}

func TestBcsWithCarrySet(t *testing.T) {
	cpu := Cpu{}
	mem := &Memory{}
	cpu.mem = mem
	cpu.p.c = 1

	actual := cpu.bcs(16)

	if expect := true; actual != expect {
		t.Errorf("Invalid return value: %v != %v", actual, expect)
	}
	if expect := 16; cpu.pc != expect {
		t.Errorf("Wrong PC")
	}
}

func TestBeq(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr	int
		zero	int
		proc	ProcStat
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "With Zero",
			addr: 16,
			zero: 1,
			expPc: 16,
			expBranch: true},
		{name: "Without Zero",
			zero: 0,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{z: tt.zero},
		}

		branch := cpu.beq(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}
