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
