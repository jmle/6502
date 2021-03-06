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

func TestGetAsWord(t *testing.T) {
	procStat := ProcStat{c:1, z:1, i:1, d:1, b:1, v:1, n:1}
	pstatus := procStat.getAsWord()

	if expProcStat := 223; pstatus != expProcStat {
		t.Errorf("Expected %+v, got %+v\n", expProcStat, pstatus)
	}
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
		{name: "With overflow 2",
			ac: 0xFF, val: 1, adc: 0,
			expAc:  0,
			expProc: ProcStat{n: 0, v: 0, c: 1},
		},
		{name: "With negative",
			val:     200,
			expAc:   200,
			expProc: ProcStat{n: 1},
		},
		{name: "With zero",
			expProc: ProcStat{z: 1},
		},
		{name: "With carry",
			ac: 255, val: 1,
			expProc: ProcStat{c: 1},
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

func TestBit(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr	int
		value	int
		ac		int
		// Expected
		expProc	ProcStat
	}{
		{name: "Sets overflow",
			value: 64,
			expProc: ProcStat{z: 0, n: 0, v: 1}},
		{name: "Overflow not set",
			value: 4,
			expProc: ProcStat{z:0, n:0, v:0},
		},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{},
		}
		cpu.mem.Write(tt.addr, tt.value)

		cpu.bit(tt.addr)
		t.Log(tt.name)

		if reflect.DeepEqual(cpu.p, tt.expProc) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestBmi(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr		int
		negative	int
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "With Negative",
			addr: 16,
			negative: 1,
			expPc: 16,
			expBranch: true},
		{name: "Without Negative",
			negative: 0,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{n: tt.negative},
		}

		branch := cpu.bmi(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}

func TestBne(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr		int
		zero		int
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "With Equals",
			addr: 16,
			zero: 0,
			expPc: 16,
			expBranch: true},
		{name: "Without Equals",
			zero: 1,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{z: tt.zero},
		}

		branch := cpu.bne(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}

func TestBpl(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr		int
		negative	int
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "With Positive",
			addr: 16,
			negative: 0,
			expPc: 16,
			expBranch: true},
		{name: "Without Equals",
			negative: 1,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{n: tt.negative},
		}

		branch := cpu.bpl(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}

func TestBvc(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr		int
		overflow	int
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "Without Overflow",
			addr: 16,
			overflow: 0,
			expPc: 16,
			expBranch: true},
		{name: "With Overflow",
			overflow: 1,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{v: tt.overflow},
		}

		branch := cpu.bvc(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}

func TestBvs(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		addr		int
		overflow	int
		// Expected
		expPc	int
		expBranch bool
	}{
		{name: "With Overflow",
			addr: 16,
			overflow: 1,
			expPc: 16,
			expBranch: true},
		{name: "Without Overflow",
			overflow: 0,
			expBranch: false},
	}{
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			p: ProcStat{v: tt.overflow},
		}

		branch := cpu.bvs(tt.addr)
		t.Log(tt.name)

		if cpu.pc != tt.expPc {
			t.Errorf("Expected %+v, got %+v\n", tt.expPc, cpu.pc)
		}
		if branch != tt.expBranch {
			t.Errorf("Expected %+v, got %+v\n", tt.expBranch, branch)
		}
	}
}

func TestClc(t *testing.T) {
	cpu := Cpu{}
	cpu.p.c = 1

	cpu.clc()

	if cpu.p.c != 0 {
		t.Errorf("Expected %+v, got %+v\n", 0, cpu.p.c)
	}
}

func TestCld(t *testing.T) {
	cpu := Cpu{}
	cpu.p.d = 1

	cpu.cld()

	if cpu.p.d != 0 {
		t.Errorf("Expected %+v, got %+v\n", 0, cpu.p.d)
	}
}

func TestCli(t *testing.T) {
	cpu := Cpu{}
	cpu.p.i = 1

	cpu.cli()

	if cpu.p.i != 0 {
		t.Errorf("Expected %+v, got %+v\n", 0, cpu.p.i)
	}
}

func TestClv(t *testing.T) {
	cpu := Cpu{}
	cpu.p.v = 1

	cpu.clv()

	if cpu.p.v != 0 {
		t.Errorf("Expected %+v, got %+v\n", 0, cpu.p.v)
	}
}

func TestCmp(t *testing.T) {
	for _, tt := range []struct {
		name		string
		reg, ac		int
		x, y		int
		data		int
		// Expected
		expProc		ProcStat
	}{
		{name: "accumulator - set carry",
			reg: A, ac: 20, data: 10,
			expProc: ProcStat{c:1, n:0, z:0},
		},
		{name: "accumulator - carry clear",
			reg: A, ac: 5, data: 10,
			expProc: ProcStat{c:0, n:1, z:0},
		},
		{name: "X - set carry",
			reg: X, x: 20, data: 10,
			expProc: ProcStat{c:1, n:0, z:0},
		},
		{name: "X - carry clear",
			reg: X, x: 5, data: 10,
			expProc: ProcStat{c:0, n:1, z:0},
		},
		{name: "Y - set carry",
			reg: Y, y: 20, data: 10,
			expProc: ProcStat{c:1, n:0, z:0},
		},
		{name: "Y - carry clear",
			reg: Y, y: 5, data: 10,
			expProc: ProcStat{c:0, n:1, z:0},
		},
	}{
		var mem Memory
		cpu := Cpu{ac: tt.ac, x: tt.x, y: tt.y}
		cpu.mem = &mem
		mem.Write(0, tt.data)

		t.Log(tt.name)

		cpu.cmp(0, tt.reg)

		if !reflect.DeepEqual(cpu.p, tt.expProc) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestDec(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem: &mem}

	mem.Write(0, 1)
	cpu.dec(0)

	if actual := cpu.mem.Read(0); actual != 0 {
		t.Errorf("Expected %+v, got %+v\n", 0, actual)
	}
	if exp := cpu.p.z; exp != 1 {
		t.Errorf("Expected %+v, got %+v\n", 1, exp)
	}
}

func TestDecxyRegX(t *testing.T) {
	cpu := Cpu{}
	cpu.x = 1

	cpu.decxy(X)

	if exp := 0; cpu.x != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.x)
	}
	if exp := (ProcStat{z:1, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestDecxyRegY(t *testing.T) {
	cpu := Cpu{}
	cpu.y = 1

	cpu.decxy(Y)

	if exp := 0; cpu.y != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.y)
	}
	if exp := (ProcStat{z:1, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestEor(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem: &mem}

	cpu.mem.Write(0, 15)
	cpu.ac = 12

	cpu.eor(0)

	if exp := 3; cpu.ac != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.ac)
	}
	if exp := (ProcStat{z:0, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestInc(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem: &mem}

	cpu.mem.Write(0, 255)

	cpu.inc(0)
	actual := cpu.mem.Read(0)

	if exp := 0; exp != actual {
		t.Errorf("Expected %+v, got %+v\n", exp, actual)
	}
	if exp := (ProcStat{z:1, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestIncxyRegX(t *testing.T) {
	cpu := Cpu{}
	cpu.x = 255

	cpu.incxy(X)

	if exp := 0; cpu.x != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.x)
	}
	if exp := (ProcStat{z:1, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestIncxyRegY(t *testing.T) {
	cpu := Cpu{}
	cpu.y = 255

	cpu.incxy(Y)

	if exp := 0; cpu.y != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.y)
	}
	if exp := (ProcStat{z:1, n:0}); !reflect.DeepEqual(exp, cpu.p) {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
	}
}

func TestJmp(t *testing.T) {
	cpu := Cpu{}

	exp := 24
	cpu.jmp(exp)

	if cpu.pc != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.pc)
	}
}

// TODO: test jsr

func TestLdrWithAc(t *testing.T) {
	for _, tt := range []struct {
		name				string
		reg, data, addr		int
		//exp
		expProc				ProcStat
	} {
		{name: "With ac",
			reg:A, data:100, addr:40,
			expProc:ProcStat{},
		},
		{name: "With X",
			reg:X, data:100, addr:40,
			expProc:ProcStat{},
		},
		{name: "With Y",
			reg:Y, data:100, addr:40,
			expProc:ProcStat{},
		},
	} {
		var mem Memory
		cpu := Cpu{mem: &mem}
		cpu.mem.Write(tt.addr, tt.data)

		cpu.ldr(tt.addr, tt.reg)

		var regValue int
		switch (tt.reg) {
		case A:
			regValue = cpu.ac
		case X:
			regValue = cpu.x
		case Y:
			regValue = cpu.y
		}

		if regValue != tt.data {
			t.Errorf("Expected %+v, got %+v\n", tt.data, regValue)
		}
		if exp := (ProcStat{z:0, n:0}); !reflect.DeepEqual(exp, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", exp, cpu.p)
		}
	}
}

func TestLsra(t *testing.T) {
	for _, tt := range []struct {
		name		string
		ac			int
		// exp
		expAc		int
		expProc		ProcStat
	} {
		{name: "With carry set",
			ac: 15, expAc: 7,
			expProc: ProcStat{c: 1},
		},
		{name: "With carry clear",
			ac: 14, expAc: 7,
			expProc: ProcStat{c: 0},
		},
	} {
		cpu := Cpu{}
		cpu.ac = tt.ac

		cpu.lsra()
		t.Log(tt.name)

		if tt.expAc != cpu.ac {
			t.Errorf("Expected %+v, got %+v\n", tt.expAc, cpu.ac)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}

	}
}

func TestLsrm(t *testing.T) {
	for _, tt := range []struct {
		name		string
		val			int
		// exp
		expVal		int
		expProc		ProcStat
	} {
		{name: "With carry set",
			val: 15, expVal: 7,
			expProc: ProcStat{c: 1},
		},
		{name: "With carry clear",
			val: 14, expVal: 7,
			expProc: ProcStat{c: 0},
		},
	} {
		var mem Memory
		cpu := Cpu{mem: &mem}
		cpu.mem.Write(0, tt.val)

		cpu.lsrm(0)
		t.Log(tt.name)

		actualVal := cpu.mem.Read(0)
		if tt.expVal != actualVal {
			t.Errorf("Expected %+v, got %+v\n", tt.expVal, actualVal)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}

	}
}

func TestNop(t *testing.T) {
	cpu := Cpu{}

	cpu.nop()

	if !reflect.DeepEqual(cpu, Cpu{}) {
		t.Errorf("Expected %+v, got %+v\n", cpu, Cpu{})
	}
}

func TestOra(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem:&mem}
	cpu.mem.Write(0, 1)
	cpu.ac = 14

	cpu.ora(0)

	if exp := 15; cpu.ac != exp {
		t.Errorf("Expected %+v, got %+v\n", exp, cpu.ac)
	}
}

func TestPha(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem:&mem, ac:14, sp:100}

	cpu.pha()

	if expSp := 99; cpu.sp != expSp {
		t.Errorf("Expected %+v, got %+v\n", expSp, cpu.sp)
	}
	if expMemSp := cpu.mem.Read(100); expMemSp != cpu.ac {
		t.Errorf("Expected %+v, got %+v\n", expMemSp, cpu.ac)
	}
}

func TestPhp(t *testing.T) {
	var mem Memory
	procStat := ProcStat{}
	procStat.setAsWord(255)
	cpu := Cpu{p:procStat, sp:40, mem:&mem}

	cpu.php()

	if expProc := 223; cpu.p.getAsWord() != expProc {
		t.Errorf("Expected %+v, got %+v\n", expProc, cpu.p.getAsWord())
	}
}

func TestPla(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem:&mem, sp:40}
	expAc := 5; mem.Write(41, expAc)

	cpu.pla()

	if cpu.ac != expAc {
		t.Errorf("Expected %+v, got %+v\n", expAc, cpu.ac)
	}
	if cpu.ac != expAc {
		t.Errorf("Expected %+v, got %+v\n", expAc, cpu.ac)
	}
}

func TestPlp(t *testing.T) {
	var mem Memory
	procStat := ProcStat{}
	procStat.setAsWord(223)
	cpu := Cpu{p:procStat, sp:40, mem:&mem}
	expSp := 41; expProc := 223;
	cpu.mem.Write(expSp, procStat.getAsWord())

	cpu.plp()

	if cpu.sp != expSp {
		t.Errorf("Expected %+v, got %+v\n", expSp, cpu.sp)
	}
	if cpu.p.getAsWord() != expProc {
		t.Errorf("Expected %+v, got %+v\n", expProc, cpu.p.getAsWord())
	}
}

func TestRola(t *testing.T) {
	for _, tt := range []struct {
		name		string
		ac			int
		proc		ProcStat
		// exp
		expAc		int
		expProc		ProcStat
	} {
		{name: "With carry bit",
			ac: 160, expAc: 64,
			proc: ProcStat{},
			expProc: ProcStat{c:1},
		},
		{name: "Without carry bit",
			ac: 64, expAc: 128,
			proc: ProcStat{},
			expProc: ProcStat{c:0, n:1},
		},
		{name: "With carry bit already set",
			ac: 64, expAc: 129,
			proc: ProcStat{c:1},
			expProc: ProcStat{c:0, n:1},
		},
	} {
		cpu := Cpu{}
		cpu.ac = tt.ac
		cpu.p = tt.proc

		cpu.rola()
		t.Log(tt.name)

		if tt.expAc != cpu.ac {
			t.Errorf("Expected %+v, got %+v\n", tt.expAc, cpu.ac)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestRolm(t *testing.T) {
	for _, tt := range []struct {
		name		string
		val			int
		proc		ProcStat
		// exp
		expVal		int
		expProc		ProcStat
	} {
		{name: "With carry bit",
			val: 160, expVal: 64,
			proc: ProcStat{},
			expProc: ProcStat{c:1},
		},
		{name: "Without carry bit",
			val: 64, expVal: 128,
			proc: ProcStat{},
			expProc: ProcStat{c:0, n:1},
		},
		{name: "With carry bit already set",
			val: 64, expVal: 129,
			proc: ProcStat{c:1},
			expProc: ProcStat{c:0, n:1},
		},
	} {
		var mem Memory
		cpu := Cpu{mem:&mem}
		cpu.p = tt.proc
		cpu.mem.Write(0, tt.val)

		cpu.rolm(0)
		t.Log(tt.name)

		if actVal := cpu.mem.Read(0); tt.expVal != actVal {
			t.Errorf("Expected %+v, got %+v\n", tt.expVal, actVal)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestRora(t *testing.T) {
	for _, tt := range []struct {
		name		string
		ac			int
		proc		ProcStat
		// exp
		expAc		int
		expProc		ProcStat
	} {
		{name: "With carry bit",
			ac: 65, expAc: 32,
			proc: ProcStat{},
			expProc: ProcStat{c:1},
		},
		{name: "Without carry bit",
			ac: 64, expAc: 32,
			proc: ProcStat{},
			expProc: ProcStat{},
		},
		{name: "With carry bit already set",
			ac: 16, expAc: 0x88,
			proc: ProcStat{c:1},
			expProc: ProcStat{c:0, n:1},
		},
	} {
		cpu := Cpu{}
		cpu.ac = tt.ac
		cpu.p = tt.proc

		cpu.rora()
		t.Log(tt.name)

		if tt.expAc != cpu.ac {
			t.Errorf("Expected %+v, got %+v\n", tt.expAc, cpu.ac)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestRorm(t *testing.T) {
	for _, tt := range []struct {
		name		string
		val			int
		proc		ProcStat
		// exp
		expVal		int
		expProc		ProcStat
	} {
		{name: "With carry bit",
			val: 65, expVal: 32,
			proc: ProcStat{},
			expProc: ProcStat{c:1},
		},
		{name: "Without carry bit",
			val: 64, expVal: 32,
			proc: ProcStat{},
			expProc: ProcStat{},
		},
		{name: "With carry bit already set",
			val: 16, expVal: 0x88,
			proc: ProcStat{c:1},
			expProc: ProcStat{c:0, n:1},
		},
	} {
		var mem Memory
		cpu := Cpu{mem:&mem}
		cpu.p = tt.proc
		cpu.mem.Write(0, tt.val)

		cpu.rorm(0)
		t.Log(tt.name)

		if actVal := cpu.mem.Read(0); tt.expVal != actVal {
			t.Errorf("Expected %+v, got %+v\n", tt.expVal, actVal)
		}
		if !reflect.DeepEqual(tt.expProc, cpu.p) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
	}
}

func TestRti(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem:&mem}
	sp := 100; cpu.sp = sp

	cpu.mem.Write(cpu.sp+1, 0xFF)
	cpu.mem.Write(cpu.sp+2, 0x10)
	cpu.mem.Write(cpu.sp+3, 0x10)

	cpu.rti()

	if expSp := sp+3; cpu.sp != expSp {
		t.Errorf("Expected %+v, got %+v\n", expSp, cpu.sp)
	}
	if expProc := 223; cpu.p.getAsWord() != expProc {
		t.Errorf("Expected %+v, got %+v\n", expProc, cpu.p.getAsWord())
	}
	if expPc := 0x1010; cpu.pc != expPc {
		t.Errorf("Expected %+v, got %+v\n", expPc, cpu.pc)
	}
}

func TestRts(t *testing.T) {
	var mem Memory
	cpu := Cpu{mem:&mem}
	sp := 100; cpu.sp = sp

	cpu.mem.Write(cpu.sp+1, 0x10)
	cpu.mem.Write(cpu.sp+2, 0x10)

	cpu.rts()

	if expPc := 0x1011; cpu.pc != expPc {
		t.Errorf("Expected %+v, got %+v\n", expPc, cpu.pc)
	}
	if expSp := sp+2; cpu.sp != expSp {
		t.Errorf("Expected %+v, got %+v\n", expSp, cpu.sp)
	}
}

func TestSbc(t *testing.T) {
	for _, tt := range []struct {
		name string
		// Set-up
		ac, val int
		proc    ProcStat
		// Expected
		expAc   int
		expProc ProcStat
	}{
		{name: "Happy Path",
			ac: 5, val: 3,
			proc: ProcStat{c: 1},
			expProc: ProcStat{c: 1},
			expAc: 2,
		},
		{name: "With overflow",
			ac: 0x50, val: 0xB0,
			proc: ProcStat{c:1},
			expAc:   0xA0,
			expProc: ProcStat{n:1, v:1},
		},
		{name: "With zero",
			proc:	 ProcStat{c:1},
			expProc: ProcStat{c:1, z:1},
		},
		// TODO: Review decimal mode
/**
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
*/
	} {
		var mem Memory
		cpu := Cpu{
			mem: &mem,
			ac:  tt.ac,
			p:   tt.proc,
		}
		mem.Write(0, tt.val)
		cpu.sbc(0)
		t.Log(tt.name)
		if !reflect.DeepEqual(cpu.p, tt.expProc) {
			t.Errorf("Expected %+v, got %+v\n", tt.expProc, cpu.p)
		}
		if cpu.ac != tt.expAc {
			t.Errorf("Expected ac %d, got %d\n", tt.expAc, cpu.ac)
		}
	}
}

