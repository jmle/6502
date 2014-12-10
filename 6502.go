package main

// the main struct, containing the main registers,
// the processor's status and the memory
type Cpu struct {
	pc, sp, ac, x, y int
	pbCrossed        bool
	p                ProcStat
	mem              Mem
}

// the status flags of the processor
// carry, zero, interrupt, decimal, break,
// negative, overflow
type ProcStat struct {
	c, z, i, d, b, n, v int
}

// returns the status as an int word
func (p *ProcStat) getAsWord() (pstatus int) {
	pstatus = p.c | p.z<<1 | p.i<<2 | p.d<<3 | p.b<<4 | p.v<<6 | p.n<<7
	return
}

func (p *ProcStat) setAsWord(pstatus int) {
	if pstatus&BIT_0 == 0 {
		p.c = 0
	} else {
		p.c = 1
	}
	if pstatus&BIT_1 == 0 {
		p.z = 0
	} else {
		p.z = 1
	}
	if pstatus&BIT_2 == 0 {
		p.i = 0
	} else {
		p.i = 1
	}
	if pstatus&BIT_3 == 0 {
		p.d = 0
	} else {
		p.d = 1
	}
	if pstatus&BIT_4 == 0 {
		p.b = 0
	} else {
		p.b = 1
	}
	if pstatus&BIT_6 == 0 {
		p.n = 0
	} else {
		p.n = 1
	}
	if pstatus&BIT_7 == 0 {
		p.v = 0
	} else {
		p.v = 1
	}
}

// sets the negative flag (n) from the
// data given
func (p *ProcStat) setN(data int) {
	if data&BIT_7 == BIT_7 {
		p.n = 1
	} else {
		p.n = 0
	}
}

// set the zero flag (z) from the data
// given
func (p *ProcStat) setZ(data int) {
	if data == 0 {
		p.z = 1
	} else {
		p.z = 0
	}
}

// the memory interface
// a memory must be provided for the cpu to work
type Mem interface {
	Read(addr int) int
	Write(addr, value int)
}

// interprets a word as bcd
func bcd(n int) int {
	return (n & 0xF) + (n & 0xF0 >> 4 * 10)
}

// Register constants
const (
	A = iota
	X
	Y
)

// Bit constants
const (
	BIT_0 = 1 << iota
	BIT_1
	BIT_2
	BIT_3
	BIT_4
	BIT_5
	BIT_6
	BIT_7
	BIT_8
)

func (cpu *Cpu) execute() (resCycles int) {
	// grab current instruction and increment pc
	inst := cpu.mem.Read(cpu.pc)
	cpu.pc++

	switch inst {
	case 0x69:
		cpu.adc(cpu.imm())
		resCycles = 2

	case 0x65:
		cpu.adc(cpu.zp())
		resCycles = 3

	case 0x75:
		cpu.adc(cpu.zpx())
		resCycles = 4

	case 0x6D:
		cpu.adc(cpu.abs())
		resCycles = 4

	case 0x7D:
		cpu.adc(cpu.abx())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x79:
		cpu.adc(cpu.aby())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x61:
		cpu.adc(cpu.indx())
		resCycles = 6

	case 0x71:
		cpu.adc(cpu.indy())
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// AND's
	case 0x29:
		cpu.and(cpu.imm())
		resCycles = 2

	case 0x25:
		cpu.and(cpu.zp())
		resCycles = 2

	case 0x35:
		cpu.and(cpu.zpx())
		resCycles = 3

	case 0x2D:
		cpu.and(cpu.abs())
		resCycles = 4

	case 0x3D:
		cpu.and(cpu.abx())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x39:
		cpu.and(cpu.aby())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x21:
		cpu.and(cpu.indx())
		resCycles = 6

	case 0x31:
		cpu.and(cpu.indy())
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// ASL's
	case 0x0A:
		cpu.asla()
		resCycles = 2

	case 0x06:
		cpu.asl(cpu.zp())
		resCycles = 5

	case 0x16:
		cpu.asl(cpu.zpx())
		resCycles = 6

	case 0x0E:
		cpu.asl(cpu.abs())
		resCycles = 6

	case 0x1E:
		cpu.asl(cpu.abx())
		resCycles = 7

	// BCC
	case 0x90:
		if cpu.bcc(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BCS
	case 0xB0:
		if cpu.bcs(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BEQ
	case 0xF0:
		if cpu.beq(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BIT
	case 0x24:
		cpu.bit(cpu.zp())
		resCycles = 3

	case 0x2C:
		cpu.bit(cpu.abs())
		resCycles = 4

	// BMI
	case 0x30:
		if cpu.bmi(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BNE
	case 0xD0:
		if cpu.bne(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BPL
	case 0x10:
		if cpu.bpl(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BRK
	case 0x00:
		cpu.brk()
		resCycles = 7

	// BVC
	case 0x50:
		if cpu.bvc(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	// BVS
	case 0x70:
		if cpu.bvs(cpu.rel()) {
			if cpu.pbCrossed {
				resCycles = 4
			} else {
				resCycles = 3
			}
		} else {
			resCycles = 2
		}

	case 0x18:
		cpu.clc()
		resCycles = 2

	case 0xD8:
		cpu.cld()
		resCycles = 2

	case 0x58:
		cpu.cli()
		resCycles = 2

	case 0xB8:
		cpu.clv()

	// CMP
	case 0xC9:
		cpu.cmp(cpu.imm(), A)
		resCycles = 2

	case 0xC5:
		cpu.cmp(cpu.zp(), A)
		resCycles = 3

	case 0xD5:
		cpu.cmp(cpu.zpx(), A)
		resCycles = 4

	case 0xCD:
		cpu.cmp(cpu.abs(), A)
		resCycles = 4

	case 0xDD:
		cpu.cmp(cpu.abx(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xD9:
		cpu.cmp(cpu.aby(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xC1:
		cpu.cmp(cpu.indx(), A)
		resCycles = 6

	case 0xD1:
		cpu.cmp(cpu.indy(), A)
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// CPX
	case 0xE0:
		cpu.cmp(cpu.imm(), X)
		resCycles = 2

	case 0xE4:
		cpu.cmp(cpu.zp(), X)
		resCycles = 3

	case 0xEC:
		cpu.cmp(cpu.abs(), X)
		resCycles = 4

	// CPY
	case 0xC0:
		cpu.cmp(cpu.imm(), Y)
		resCycles = 2

	case 0xC4:
		cpu.cmp(cpu.zp(), Y)
		resCycles = 3

	case 0xCC:
		cpu.cmp(cpu.abs(), Y)
		resCycles = 4

	// DEC
	case 0xC6:
		cpu.dec(cpu.zp())
		resCycles = 5

	case 0xD6:
		cpu.dec(cpu.zpx())
		resCycles = 6

	case 0xCE:
		cpu.dec(cpu.abs())
		resCycles = 6

	case 0xDE:
		cpu.dec(cpu.abx())
		resCycles = 7

	// DEX
	case 0xCA:
		cpu.decxy(X)
		resCycles = 2

	// DEY
	case 0x88:
		cpu.decxy(Y)
		resCycles = 2

	// EOR
	case 0x49:
		cpu.eor(cpu.imm())
		resCycles = 2

	case 0x45:
		cpu.eor(cpu.zp())
		resCycles = 3

	case 0x55:
		cpu.eor(cpu.zpx())
		resCycles = 4

	case 0x4D:
		cpu.eor(cpu.abs())
		resCycles = 4

	case 0x5D:
		cpu.eor(cpu.abx())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x59:
		cpu.eor(cpu.aby())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x41:
		cpu.eor(cpu.indx())
		resCycles = 6

	case 0x51:
		cpu.eor(cpu.indy())
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// INC
	case 0xE6:
		cpu.inc(cpu.zp())
		resCycles = 5

	case 0xF6:
		cpu.inc(cpu.zpx())
		resCycles = 6

	case 0xEE:
		cpu.inc(cpu.abs())
		resCycles = 6

	case 0xFE:
		cpu.inc(cpu.abx())
		resCycles = 7

	// INX
	case 0xE8:
		cpu.incxy(X)
		resCycles = 2

	// INY
	case 0xC8:
		cpu.incxy(Y)
		resCycles = 2

	// JMP
	case 0x4C:
		cpu.jmp(cpu.abs())
		resCycles = 3

	case 0x6C:
		cpu.jmp(cpu.ind())
		resCycles = 5

	// JSR
	case 0x20:
		cpu.jsr(cpu.abs())
		resCycles = 6

	// LDA
	case 0xA9:
		cpu.ldr(cpu.imm(), A)
		resCycles = 2

	case 0xA5:
		cpu.ldr(cpu.zp(), A)
		resCycles = 3

	case 0xB5:
		cpu.ldr(cpu.zpx(), A)
		resCycles = 4

	case 0xAD:
		cpu.ldr(cpu.abs(), A)
		resCycles = 4

	case 0xBD:
		cpu.ldr(cpu.abx(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xB9:
		cpu.ldr(cpu.aby(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xA1:
		cpu.ldr(cpu.indx(), A)
		resCycles = 6

	case 0xB1:
		cpu.ldr(cpu.indy(), A)
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// LDX
	case 0xA2:
		cpu.ldr(cpu.imm(), X)
		resCycles = 2

	case 0xA6:
		cpu.ldr(cpu.zp(), X)
		resCycles = 3

	case 0xB6:
		cpu.ldr(cpu.zpy(), X)
		resCycles = 4

	case 0xAE:
		cpu.ldr(cpu.abs(), X)
		resCycles = 4

	case 0xBE:
		cpu.ldr(cpu.aby(), X)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	// LDY
	case 0xA0:
		cpu.ldr(cpu.imm(), Y)
		resCycles = 2

	case 0xA4:
		cpu.ldr(cpu.zp(), Y)
		resCycles = 3

	case 0xB4:
		cpu.ldr(cpu.zpx(), Y)
		resCycles = 4

	case 0xAC:
		cpu.ldr(cpu.abs(), Y)
		resCycles = 4

	case 0xBC:
		cpu.ldr(cpu.abx(), Y)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	// LSR
	case 0x4A:
		cpu.lsra()
		resCycles = 2

	case 0x46:
		cpu.lsrm(cpu.zp())
		resCycles = 5

	case 0x56:
		cpu.lsrm(cpu.zpx())
		resCycles = 6

	case 0x4E:
		cpu.lsrm(cpu.abs())
		resCycles = 6

	case 0x5E:
		cpu.lsrm(cpu.abx())
		resCycles = 7

	// NOP
	case 0xEA:
		cpu.nop()
		resCycles = 2

	// ORA
	case 0x09:
		cpu.ora(cpu.imm())
		resCycles = 2

	case 0x05:
		cpu.ora(cpu.imm())
		resCycles = 2

	case 0x15:
		cpu.ora(cpu.imm())
		resCycles = 3

	case 0x1D:
		cpu.ora(cpu.imm())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x19:
		cpu.ora(cpu.imm())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x01:
		cpu.ora(cpu.imm())
		resCycles = 6

	case 0x11:
		cpu.ora(cpu.imm())
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// PHA
	case 0x48:
		cpu.pha()
		resCycles = 3

	// PHP
	case 0x08:
		cpu.php()
		resCycles = 3

	// PLA
	case 0x68:
		cpu.pla()
		resCycles = 4

	// PLP
	case 0x28:
		cpu.plp()
		resCycles = 4

	// ROL
	case 0x2A:
		cpu.rola()
		resCycles = 2

	case 0x26:
		cpu.rolm(cpu.zp())
		resCycles = 5

	case 0x36:
		cpu.rolm(cpu.zpx())
		resCycles = 6

	case 0x2E:
		cpu.rolm(cpu.abs())
		resCycles = 6

	case 0x3E:
		cpu.rolm(cpu.abx())
		resCycles = 7

	// ROR
	case 0x6A:
		cpu.rora()
		resCycles = 2

	case 0x66:
		cpu.rorm(cpu.zp())
		resCycles = 5

	case 0x76:
		cpu.rorm(cpu.zpx())
		resCycles = 6

	case 0x6E:
		cpu.rorm(cpu.abs())
		resCycles = 6

	case 0x7E:
		cpu.rorm(cpu.abx())
		resCycles = 7

	// RTI
	case 0x40:
		cpu.rti()
		resCycles = 6

	// RTS
	case 0x60:
		cpu.rts()
		resCycles = 6

	// SBC
	case 0xE9:
		cpu.sbc(cpu.imm())
		resCycles = 2

	case 0xE5:
		cpu.sbc(cpu.zp())
		resCycles = 3

	case 0xF5:
		cpu.sbc(cpu.zpx())
		resCycles = 4

	case 0xED:
		cpu.sbc(cpu.abs())
		resCycles = 4

	case 0xFD:
		cpu.sbc(cpu.abx())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xF9:
		cpu.sbc(cpu.aby())
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0xE1:
		cpu.sbc(cpu.indx())
		resCycles = 6

	case 0xF1:
		cpu.sbc(cpu.indy())
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// SEC
	case 0x38:
		cpu.sec()
		resCycles = 2

	// SED
	case 0xF8:
		cpu.sed()
		resCycles = 2

	// SEI
	case 0x78:
		cpu.sei()
		resCycles = 2

	// STA
	case 0x85:
		cpu.st(cpu.zp(), A)
		resCycles = 3

	case 0x95:
		cpu.st(cpu.zpx(), A)
		resCycles = 4

	case 0x8D:
		cpu.st(cpu.abs(), A)
		resCycles = 4

	case 0x9D:
		cpu.st(cpu.abx(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x99:
		cpu.st(cpu.aby(), A)
		if cpu.pbCrossed {
			resCycles = 5
		} else {
			resCycles = 4
		}

	case 0x81:
		cpu.st(cpu.indx(), A)
		resCycles = 6

	case 0x91:
		cpu.st(cpu.indy(), A)
		if cpu.pbCrossed {
			resCycles = 6
		} else {
			resCycles = 5
		}

	// STX
	case 0x86:
		cpu.st(cpu.zp(), X)
		resCycles = 3

	case 0x96:
		cpu.st(cpu.zpy(), X)
		resCycles = 4

	case 0x8E:
		cpu.st(cpu.abs(), X)
		resCycles = 4

	// STY
	case 0x84:
		cpu.st(cpu.zp(), Y)
		resCycles = 3

	case 0x94:
		cpu.st(cpu.zpx(), Y)
		resCycles = 4

	case 0x8C:
		cpu.st(cpu.abs(), Y)
		resCycles = 4

	// TAX
	case 0xAA:
		cpu.taxy(X)
		resCycles = 2

	// TAY
	case 0xA8:
		cpu.taxy(Y)
		resCycles = 2

	// TSX
	case 0xBA:
		cpu.tsx()
		resCycles = 2

	// TXA
	case 0x8A:
		cpu.txya(X)
		resCycles = 2

	// TXS
	case 0x9A:
		cpu.txs()
		resCycles = 2

	// TYA
	case 0x98:
		cpu.txya(Y)
		resCycles = 2
	}

	return
}

// instruction implementations
// add with carry
func (cpu *Cpu) adc(addr int) {
	data := cpu.mem.Read(addr)

	if cpu.p.d == 1 {
		// Calculate auxiliary value
		aux := bcd2bin(cpu.ac) + bcd2bin(data) + cpu.p.c

		if aux > 99 {
			aux -= 100
			cpu.p.c = 1
		} else {
			cpu.p.c = 0
		}

		cpu.ac = bin2bcd(aux)
	} else {
		// Calculate auxiliary value
		aux := cpu.ac + data + cpu.p.c

		// Set flags: overflow, sign, zero, and carry
		overflow := ((cpu.ac & BIT_7) != (aux & BIT_7))
		if overflow {
			cpu.p.v = 1
		} else {
			cpu.p.v = 0
		}
		cpu.p.setN(aux)
		cpu.p.setZ(aux)

		if aux > 255 {
			cpu.p.c = 1
		} else {
			cpu.p.c = 0
		}

		// take the possible carry out
		cpu.ac = aux & 0xFF
	}
}

// and accumulator with memory
func (cpu *Cpu) and(addr int) {
	data := cpu.mem.Read(addr)
	cpu.ac &= data

	// flags: sign, zero.
	cpu.p.setN(cpu.ac)
	cpu.p.setZ(cpu.ac)
}

// asymetric shift left accumulator
func (cpu *Cpu) asla() {
	carry := (cpu.ac & BIT_7) == BIT_7
	if carry {
		cpu.p.c = 1
	} else {
		cpu.p.c = 0
	}
	cpu.ac <<= 1
	cpu.ac &= 0xFF

	cpu.p.setN(cpu.ac)
	cpu.p.setZ(cpu.ac)
}

// asymetric shift left memory
func (cpu *Cpu) asl(addr int) {
	data := cpu.mem.Read(addr)

	carry := (data & BIT_7) == BIT_7
	if carry {
		cpu.p.c = 1
	} else {
		cpu.p.c = 0
	}
	data = (data << 1) & 0xFE

	cpu.p.setN(data)
	cpu.p.setZ(data)

	cpu.mem.Write(addr, data)
}

// branch if carry clear
func (cpu *Cpu) bcc(addr int) bool {
	if cpu.p.c == 0 {
		cpu.pc = addr
		return true
	}
	return false
}

// branch if carry set
func (cpu *Cpu) bcs(addr int) bool {
	if cpu.p.c == 1 {
		cpu.pc = addr
		return true
	}
	return false
}

// branch if equals (checks zero)
func (cpu *Cpu) beq(addr int) bool {
	if cpu.p.z == 1 {
		cpu.pc = addr
		return true
	}
	return false
}

// bit test
func (cpu *Cpu) bit(addr int) {
	data := cpu.mem.Read(addr) & cpu.ac

	if data&BIT_6 != 0 {
		cpu.p.v = 1
	} else {
		cpu.p.v = 0
	}
	cpu.p.setN(data)
	cpu.p.setZ(data)
}

// branch if negative
func (cpu *Cpu) bmi(addr int) bool {
	if cpu.p.n == 1 {
		cpu.pc = addr
		return true
	}
	return false
}

// branch if not equal (checks zero)
func (cpu *Cpu) bne(addr int) bool {
	if cpu.p.z == 0 {
		cpu.pc = addr
		return true
	}
	return false
}

// branch if positive
func (cpu *Cpu) bpl(addr int) bool {
	if cpu.p.n == 0 {
		cpu.pc = addr
		return true
	}
	return false
}

// break
// TODO
func (cpu *Cpu) brk() {
	var l, h int

	// Even though the brk instruction is just one byte long, the pc is
	// incremented, meaning that the instruction after brk is ignored.
	cpu.pc++
	/*
		cpu.mem.Write(cpu.sp, cpu.pc&0xF0)
		cpu.mem.commit()
		cpu.sp--
		cpu.mem.Write(cpu.sp, cpu.pc&0xF)
		cpu.mem.commit()
		cpu.sp--
		cpu.mem.Write(cpu.sp, cpu.p.b)
		cpu.sp--
	*/
	l = cpu.mem.Read(0xFFFE)
	h = cpu.mem.Read(0xFFFF) << 8

	cpu.pc = h | l
}

// branch if bit clear
func (cpu *Cpu) bvc(addr int) bool {
	if cpu.p.v == 0 {
		cpu.pc = addr
		return true
	}
	return false
}

// branch if bit set
func (cpu *Cpu) bvs(addr int) bool {
	if cpu.p.v == 1 {
		cpu.pc = addr
		return true
	}
	return false
}

// clear carry flag
func (cpu *Cpu) clc() {
	cpu.p.c = 0
}

// clear decimal flag
func (cpu *Cpu) cld() {
	cpu.p.d = 0
}

// clear interrupt flag
func (cpu *Cpu) cli() {
	cpu.p.i = 0
}

// clear bit bit
func (cpu *Cpu) clv() {
	cpu.p.v = 0
}

// compare accumulator with memory
func (cpu *Cpu) cmp(addr, r int) {
	data := cpu.mem.Read(addr)

	// Calculate auxiliary value
	t := 0
	switch r {
	case A:
		t = cpu.ac - data
		if cpu.ac >= data {
			cpu.p.c = 1
		} else {
			cpu.p.c = 0
		}

	case X:
		t = cpu.x - data
		if cpu.x >= data {
			cpu.p.c = 1
		} else {
			cpu.p.c = 0
		}

	case Y:
		t = cpu.y - data
		if cpu.y >= data {
			cpu.p.c = 1
		} else {
			cpu.p.c = 0
		}
	}

	// Set flags
	cpu.p.setN(t)
	cpu.p.setZ(t)
}

// decrement memory
func (cpu *Cpu) dec(addr int) {
	data := cpu.mem.Read(addr)

	// Decrement & AND 0xFF
	data = (data - 1) & 0xFF
	cpu.mem.Write(addr, data)

	// Set flags
	cpu.p.setN(data)
	cpu.p.setZ(data)
}

// decrement register
func (cpu *Cpu) decxy(r int) {
	switch r {
	case X:
		cpu.x = (cpu.x - 1) & 0xFF
		cpu.p.setN(cpu.x)
		cpu.p.setZ(cpu.x)

	case Y:
		cpu.y = (cpu.y - 1) & 0xFF
		cpu.p.setN(cpu.y)
		cpu.p.setZ(cpu.y)
	}
}

// exclusive or accumulator and memory
func (cpu *Cpu) eor(addr int) {
	data := cpu.mem.Read(addr)

	cpu.ac ^= data
	cpu.p.setN(cpu.ac)
	cpu.p.setZ(cpu.ac)
}

// increment memory
func (cpu *Cpu) inc(addr int) {
	data := cpu.mem.Read(addr)

	data++
	data &= 0xFF
	cpu.mem.Write(addr, data)

	cpu.p.setN(data)
	cpu.p.setZ(data)
}

// increment register
func (cpu *Cpu) incxy(r int) {
	switch r {
	case X:
		cpu.x = (cpu.x + 1) & 0xFF
		cpu.p.setN(cpu.x)
		cpu.p.setZ(cpu.x)

	case Y:
		cpu.y = (cpu.y + 1) & 0xFF
		cpu.p.setN(cpu.y)
		cpu.p.setZ(cpu.y)
	}
}

// jump to address
func (cpu *Cpu) jmp(addr int) {
	cpu.pc = addr
}

// jump to subrutine
func (cpu *Cpu) jsr(addr int) {
	t := cpu.pc - 1

	// Push PC onto the stack
	cpu.mem.Write(cpu.sp, (t&0xFF00)>>8)
	// TODO: what, why?
	//cpu.mem.commit()
	cpu.sp--
	cpu.mem.Write(cpu.sp, t&0xFF)
	cpu.sp--

	// Jump
	cpu.pc = addr
}

// load memory to register
func (cpu *Cpu) ldr(addr, r int) {
	data := cpu.mem.Read(addr)

	// One function for three different opcodes. Have to switch the register
	switch r {
	case A:
		cpu.ac = data
		cpu.p.setN(cpu.ac)
		cpu.p.setZ(cpu.ac)

	case X:
		cpu.x = data
		cpu.p.setN(cpu.x)
		cpu.p.setZ(cpu.x)

	case Y:
		cpu.y = data
		cpu.p.setN(cpu.y)
		cpu.p.setZ(cpu.y)
	}
}

// shift right accumulator
func (cpu *Cpu) lsra() {
	cpu.p.n = 0
	if cpu.ac&BIT_0 == 0 {
		cpu.p.c = 0
	} else {
		cpu.p.c = 1
	}

	cpu.ac = (cpu.ac >> 1) & 0x7F
	cpu.p.setZ(cpu.ac)
}

// right shift memory
func (cpu *Cpu) lsrm(addr int) {
	data := cpu.mem.Read(addr)

	cpu.p.n = 0
	if data&BIT_0 == 0 {
		cpu.p.c = 0
	} else {
		cpu.p.c = 1
	}
	data = (data >> 1) & 0x7F
	cpu.p.setZ(data)

	cpu.mem.Write(addr, data)
}

// no operation
func (cpu *Cpu) nop() {

}

// or with accumulator
func (cpu *Cpu) ora(addr int) {
	data := cpu.mem.Read(addr)

	cpu.ac |= data
	cpu.p.setZ(cpu.ac)
	cpu.p.setN(cpu.ac)
}

// push accumulator to stack
func (cpu *Cpu) pha() {
	cpu.mem.Write(cpu.sp, cpu.ac)
	cpu.sp--
}

// push processor status to stack
func (cpu *Cpu) php() {
	cpu.mem.Write(cpu.sp, cpu.p.getAsWord())
	cpu.sp--
}

// put stack in accumulator
func (cpu *Cpu) pla() {
	cpu.sp++
	cpu.ac = cpu.mem.Read(cpu.sp)

	cpu.p.setN(cpu.ac)
	cpu.p.setZ(cpu.ac)
}

// set push stack to processor status
func (cpu *Cpu) plp() {
	cpu.sp++
	cpu.p.setAsWord(cpu.mem.Read(cpu.sp))
}

// rotate accumulator left
func (cpu *Cpu) rola() {
	// This opcode uses the carry to fill the LSB, and then sets the carry
	// according to the MSB of the rolled byte

	// Take from the byte what will be the future carry
	var t int
	if cpu.ac&BIT_7 != 0 {
		t = 1
	} else {
		t = 0
	}

	// Rotate left and &
	cpu.ac = (cpu.ac << 1) & 0xFE
	// Set LSB with the carry value from before the operation
	cpu.ac |= cpu.p.c
	// Set the next carry
	cpu.p.c = t
	// Set flags
	cpu.p.setZ(cpu.ac)
	cpu.p.setN(cpu.ac)
}

// rotate memory left
func (cpu *Cpu) rolm(addr int) {
	data := cpu.mem.Read(addr)
	var t int
	if data&BIT_7 != 0 {
		t = 1
	} else {
		t = 0
	}

	// Rotate left and &
	data = (data << 1) & 0xFE
	// Set LSB with the carry value from before the operation
	data |= cpu.p.c
	// Set the next carry
	cpu.p.c = t
	// Set flags
	cpu.p.setZ(data)
	cpu.p.setN(data)

	// Write to memory
	cpu.mem.Write(addr, data)
}

// rorate accumulator right
func (cpu *Cpu) rora() {
	// This opcode uses the carry to fill the MSB, and then sets the carry
	// according to the LSB of the rolled byte

	// Take from the byte what will be the future carry
	var t int
	if cpu.ac&BIT_0 != 0 {
		t = 1
	} else {
		t = 0
	}

	// Rotate right and &
	cpu.ac = (cpu.ac >> 1) & 0x7F

	// Set MSB with the carry value from before the operation
	if cpu.p.c == 1 {
		cpu.ac |= 0x80
	} else {
		cpu.ac |= 0x00
	}

	// Set the next carry
	cpu.p.c = t
	// Set flags
	cpu.p.setZ(cpu.ac)
	cpu.p.setN(cpu.ac)
}

// rotate memory right
func (cpu *Cpu) rorm(addr int) {
	data := cpu.mem.Read(addr)
	var t int
	if data&BIT_0 != 0 {
		t = 1
	} else {
		t = 0
	}

	// Rotate right and &
	data = (data >> 1) & 0x7F

	// Set LSB with the carry value from before the operation
	if cpu.p.c == 1 {
		data |= 0x80
	} else {
		data |= 0x00
	}

	// Set the next carry
	cpu.p.c = t
	// Set flags
	cpu.p.setZ(data)
	cpu.p.setN(data)

	// Write to memory
	cpu.mem.Write(addr, data)
}

// return from interrupt
func (cpu *Cpu) rti() {
	var l, h int
	// TODO: increment or decrement sp?
	cpu.sp++
	cpu.p.setAsWord(cpu.mem.Read(cpu.sp))
	cpu.sp++
	l = cpu.mem.Read(cpu.sp)
	cpu.sp++
	h = cpu.mem.Read(cpu.sp)

	cpu.pc = (h << 8) | l
}

// return from subrutine
func (cpu *Cpu) rts() {
	var l, h int

	cpu.sp++
	l = cpu.mem.Read(cpu.sp)
	cpu.sp++
	h = cpu.mem.Read(cpu.sp)

	cpu.pc = ((h << 8) | l) + 1
}

// substract with carry
func (cpu *Cpu) sbc(addr int) {
	data := cpu.mem.Read(addr)

	var t int
	// If decimal mode is on...
	if cpu.p.d == 1 {
		// When using SBC, the code should have used SEC to set the carry
		// before. This is to make sure that, if we need to borrow, there is
		// something to borrow.
		var negcarry int
		if cpu.p.c&BIT_0 != 0 {
			negcarry = 0
		} else {
			negcarry = 1
		}
		t = bcd2bin(cpu.ac) - bcd2bin(data) - negcarry

		if t > 99 || t < 0 {
			cpu.p.v = 1
		} else {
			cpu.p.v = 0
		}
	} else {
		var negcarry int
		if cpu.p.c&BIT_0 != 0 {
			negcarry = 0
		} else {
			negcarry = 1
		}
		t = cpu.ac - data - negcarry

		if t > 127 || t < -128 {
			cpu.p.v = 1
		} else {
			cpu.p.v = 0
		}
	}

	// Set the flags
	if t >= 0 {
		cpu.p.c = 1
	} else {
		cpu.p.c = 0
	}
	cpu.p.n = t
	cpu.p.z = t

	// Write the result (ANDed, just in case it overflowed)
	cpu.ac = t & 0xFF
}

// set carry flag
func (cpu *Cpu) sec() {
	cpu.p.c = 1
}

// set decimal flag
func (cpu *Cpu) sed() {
	cpu.p.d = 1
}

// set interrupt
func (cpu *Cpu) sei() {
	cpu.p.i = 1
}

// store register in memory
func (cpu *Cpu) st(addr, r int) {
	switch r {
	case A:
		cpu.mem.Write(addr, cpu.ac)

	case X:
		cpu.mem.Write(addr, cpu.x)

	case Y:
		cpu.mem.Write(addr, cpu.y)
	}
}

// copy accumulator in register
func (cpu *Cpu) taxy(r int) {
	switch r {
	case X:
		cpu.x = cpu.ac
		cpu.p.n = cpu.x
		cpu.p.z = cpu.x

	case Y:
		cpu.y = cpu.ac
		cpu.p.n = cpu.y
		cpu.p.z = cpu.y
	}
}

// load stack in register
func (cpu *Cpu) tsx() {
	cpu.x = cpu.sp
	cpu.p.n = cpu.x
	cpu.p.z = cpu.x
}

// load accumulator with register
func (cpu *Cpu) txya(r int) {
	switch r {
	case X:
		cpu.ac = cpu.x

	case Y:
		cpu.ac = cpu.y
	}

	cpu.p.n = cpu.ac
	cpu.p.z = cpu.ac
}

// set stack to register x
func (cpu *Cpu) txs() {
	cpu.sp = cpu.x
}

// -----------------------------------
// Addressing modes
// - Page crossing is checked
// - The operand is retrieved and stored for debugging purposes
// -----------------------------------

/**
 * Immediate: The operand is used directly to perform the computation.
 */
func (cpu *Cpu) imm() int {
	addr := cpu.pc
	cpu.pc++
	return addr
}

// Zero page: A single byte specifies an address in the first page of mem
// ($00xx), also known as the zero page, and the byte at that address is
// used to perform the computation.
func (cpu *Cpu) zp() int {
	addr := cpu.mem.Read(cpu.pc) & 0xFF
	cpu.pc++
	return addr
}

// Zero page,X: The value in X is added to the specified zero page address
// for a sum address. The value at the sum address is used to perform the
// computation.
func (cpu *Cpu) zpx() int {
	addr := cpu.mem.Read(cpu.pc)
	cpu.pc++
	return (addr + cpu.x) & 0xFF
}

// Zero page,Y: The value in Y is added to the specified zero page address
// for a sum address. The value at the sum address is used to perform the
// computation.
func (cpu *Cpu) zpy() int {
	addr := cpu.mem.Read(cpu.pc)
	cpu.pc++
	return (addr + cpu.y) & 0xFF
}

// The offset specified is added to the current address stored in the
// Program Counter (PC). Offsets can range from -128 to +127.
func (cpu *Cpu) rel() int {
	addr := cpu.mem.Read(cpu.pc)
	cpu.pc++
	offset := int((byte(addr)))
	addr = cpu.pc + offset

	cpu.pageBoundaryCrossed(cpu.pc, addr)

	return addr
}

// Absolute: A full 16-bit address is specified and the byte at that address
// is used to perform the computation.
func (cpu *Cpu) abs() int {
	op1 := cpu.mem.Read(cpu.pc)
	cpu.pc++
	op2 := cpu.mem.Read(cpu.pc)
	cpu.pc++
	addr := cpu.mem.Read(op1) | (cpu.mem.Read(op2) << 8)

	return addr
}

// Absolute indexed with X: The value in X is added to the specified address
// for a sum address. The value at the sum address is used to perform the
// computation.
func (cpu *Cpu) abx() int {
	op1 := cpu.mem.Read(cpu.pc)
	cpu.pc++
	op2 := cpu.mem.Read(cpu.pc)
	cpu.pc++
	addr := (cpu.mem.Read(op1) | (cpu.mem.Read(op2) << 8))

	before := addr
	after := (before + cpu.x)

	cpu.pageBoundaryCrossed(before, after)

	return after & 0xFFFF
}

// Absolute indexed with Y: The value in Y is added to the specified address
// for a sum address. The value at the sum address is used to perform the
// computation.
func (cpu *Cpu) aby() int {
	op1 := cpu.pc
	cpu.pc++
	op2 := cpu.pc
	cpu.pc++
	addr := (cpu.mem.Read(op1) | (cpu.mem.Read(op2) << 8))
	before := addr
	after := (before + cpu.y)

	cpu.pageBoundaryCrossed(before, after)

	return after & 0xFFFF
}

// Indirect addressing. With this instruction, the 8-bit address (location)
// supplied by the programmer is considered to be a Zero-Page address, that
// is, an address in the first 256 (0..255) bytes of memory. The content of
// this Zero-Page address must contain the low 8-bits of a memory address.
// The following byte (the contents of address+1) must contain the upper
// 8-bits of a memory address
func (cpu *Cpu) ind() int {
	addr := cpu.mem.Read(cpu.pc) & 0xFF
	cpu.pc++

	return cpu.mem.Read(addr) | (cpu.mem.Read(addr+1) << 8)
}

// Zero Page Indexed Indirect: Much like Indirect Addressing, but the
// content of the index register is added to the Zero-Page address
// (location)
func (cpu *Cpu) indx() int {
	addr := cpu.mem.Read(cpu.pc) & 0xFF
	cpu.pc++

	return (cpu.mem.Read(addr+cpu.x) | (cpu.mem.Read(addr+1+cpu.x) << 8))
}

// Indirect Indexed Addressing: Much like Indexed Addressing, but the
// contents of the index register is added to the Base_Location after it is
// read from Zero-Page memory.
func (cpu *Cpu) indy() int {
	addr := cpu.mem.Read(cpu.pc) & 0xFF
	cpu.pc++

	before := cpu.mem.Read(cpu.mem.Read(addr) | (cpu.mem.Read(addr+1) << 8))
	after := before + cpu.y

	cpu.pageBoundaryCrossed(before, after)

	return after
}

// helper functions

// Checks if a page boundary was crossed between two addresses.
//
// "For example, in the instruction LDA 1234,X, where the value in the X
// register is added to address 1234 to get the effective address to load
// the accumulator from, the operand's low byte is fetched before the high
// byte, so the processor can start adding the X register's value before it
// has the high byte. If there is no carry operation, the entire indexed
// operation takes only four clocks, which is one microsecond at 4MHz. If
// there is a carry requiring the high byte to be incremented, it takes one
// additional clock." (Taken from the AtariAge forums)
func (cpu *Cpu) pageBoundaryCrossed(addr1, addr2 int) {
	cpu.pbCrossed = ((addr1 ^ addr2) & BIT_8) != 0
}

// interprets a word as bcd
func bcd2bin(n int) int {
	return (n & 0xF) + ((n & 0xF0) >> 4 * 10)
}

// turns a binary word into bcd
func bin2bcd(n int) int {
	units := n % 10
	tens := (n - units) / 10

	return (tens << 4) | units
}
