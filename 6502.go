package main

// the main struct, containing the main registers,
// the processor's status and the memory
type Cpu struct {
	pc, sp, ac, x, y 	int
	pbCrossed 			bool
	p 					ProcStat
	mem 				Mem
}

// the status flags of the processor
// carry, zero, interrupt, decimal, break,
// negative, overflow
type ProcStat struct {
	c, z, i, d, b, n, v int
}

// the memory is basically an array with
// fixed size (64K) 
// TODO: reader/writer
type Mem struct {
	m 		[1 << 16]int
}

func (mem *Mem) read(i int) (value int) {
	value = mem.m[i]
	return
}

func (mem *Mem) write(addr, i int) {
	mem.m[addr] = i
}

func (cpu *Cpu) execute() (resCycles int) {
	// grab current instruction and increment pc
	inst := cpu.mem.read(cpu.pc)
	cpu.pc++

	switch inst {
	case 0x69:
		// TODO: these are methods now, not functions, like this:
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
        cpu.aslacpu.()
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
        cpu.brk(cpu.)
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
        cpu.clc(cpu.)
        resCycles = 2

    case 0xD8:
        cpu.cld(cpu.)
        resCycles = 2

    case 0x58:
        cpu.cli(cpu.)
        resCycles = 2

    case 0xB8:
        cpu.clv(cpu.)
        
    // CMP
    case 0xC9:
        cpu.cmp(cpu.imm(), R.A)
        resCycles = 2

    case 0xC5:
        cpu.cmp(cpu.zp(), R.A)
        resCycles = 3

    case 0xD5:
        cpu.cmp(cpu.zpx(), R.A)
        resCycles = 4

    case 0xCD:
        cpu.cmp(cpu.abs(), R.A)
        resCycles = 4

    case 0xDD:
        cpu.cmp(cpu.abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xD9:
        cpu.cmp(cpu.aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xC1:
        cpu.cmp(cpu.indx(), R.A)
        resCycles = 6

    case 0xD1:
        cpu.cmp(cpu.indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // CPX
    case 0xE0:
        cpu.cmp(cpu.imm(), R.X)
        resCycles = 2

    case 0xE4:
        cpu.cmp(cpu.zp(), R.X)
        resCycles = 3

    case 0xEC:
        cpu.cmp(cpu.abs(), R.X)
        resCycles = 4

    // CPY
    case 0xC0:
        cpu.cmp(cpu.imm(), R.Y)
        resCycles = 2

    case 0xC4:
        cpu.cmp(cpu.zp(), R.Y)
        resCycles = 3

    case 0xCC:
        cpu.cmp(cpu.abs(), R.Y)
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
        cpu.decxy(R.X)
        resCycles = 2

    // DEY
    case 0x88:
        cpu.decxy(R.Y)
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
        cpu.incxy(R.X)
        resCycles = 2

    // INY
    case 0xC8:
        cpu.incxy(R.Y)
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
        cpu.ldr(cpu.imm(), R.A)
        resCycles = 2

    case 0xA5:
        cpu.ldr(cpu.zp(), R.A)
        resCycles = 3

    case 0xB5:
        cpu.ldr(cpu.zpx(), R.A)
        resCycles = 4

    case 0xAD:
        cpu.ldr(cpu.abs(), R.A)
        resCycles = 4

    case 0xBD:
        cpu.ldr(cpu.abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xB9:
        cpu.ldr(cpu.aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xA1:
        cpu.ldr(cpu.indx(), R.A)
        resCycles = 6

    case 0xB1:
        cpu.ldr(cpu.indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // LDX
    case 0xA2:
        cpu.ldr(cpu.imm(), R.X)
        resCycles = 2

    case 0xA6:
        cpu.ldr(cpu.zp(), R.X)
        resCycles = 3

    case 0xB6:
        cpu.ldr(cpu.zpy(), R.X)
        resCycles = 4

    case 0xAE:
        cpu.ldr(cpu.abs(), R.X)
        resCycles = 4

    case 0xBE:
        cpu.ldr(cpu.aby(), R.X)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    // LDY
    case 0xA0:
        cpu.ldr(cpu.imm(), R.Y)
        resCycles = 2

    case 0xA4:
        cpu.ldr(cpu.zp(), R.Y)
        resCycles = 3

    case 0xB4:
        cpu.ldr(cpu.zpx(), R.Y)
        resCycles = 4

    case 0xAC:
        cpu.ldr(cpu.abs(), R.Y)
        resCycles = 4

    case 0xBC:
        cpu.ldr(cpu.abx(), R.Y)
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
        cpu.lsrmcpu.(zp())
        resCycles = 5

    case 0x56:
        cpu.lsrmcpu.(zpx())
        resCycles = 6

    case 0x4E:
        cpu.lsrmcpu.(abs())
        resCycles = 6

    case 0x5E:
        cpu.lsrmcpu.(abx())
        resCycles = 7

    // NOP
    case 0xEA:
        cpu.nop()
        resCycles = 2

    // ORA
    case 0x09:
        cpu.ora(imm())
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
        cpu.rolacpu.()
        resCycles = 2

    case 0x6E:
        cpu.rormcpu.(abs())
        resCycles = 6

    case 0x7E:
        cpu.rormcpu.(abx())
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
        cpu.st(cpu.zp(), R.A)
        resCycles = 3

    case 0x95:
        cpu.st(cpu.zpx(), R.A)
        resCycles = 4

    case 0x8D:
        cpu.st(cpu.abs(), R.A)
        resCycles = 4

    case 0x9D:
        cpu.st(cpu.abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x99:
        cpu.st(cpu.aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x81:
        cpu.st(cpu.indx(), R.A)
        resCycles = 6

    case 0x91:
        cpu.st(indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }
        
    // STX
    case 0x86:
        cpu.st(cpu.zp(), R.X)
        resCycles = 3

    case 0x96:
        cpu.st(cpu.zpy(), R.X)
        resCycles = 4

    case 0x8E:
        cpu.st(cpu.abs(), R.X)
        resCycles = 4

    // STY
    case 0x84:
        cpu.st(cpu.zp(), R.Y)
        resCycles = 3

    case 0x94:
        cpu.st(cpu.zpx(), R.Y)
        resCycles = 4

    case 0x8C:
        cpu.st(cpu.abs(), R.Y)
        resCycles = 4

    // TAX
    case 0xAA:
        cpu.taxy(R.X)
        resCycles = 2

    // TAY
    case 0xA8:
        cpu.taxy(R.Y)
        resCycles = 2

    // TSX
    case 0xBA:
        cpu.tsx()
        resCycles = 2

    // TXA
    case 0x8A:
        cpu.txya(R.X)
        resCycles = 2

    // TXS
    case 0x9A:
        cpu.txs()
        resCycles = 2

    // TYA
    case 0x98:
        cpu.txya(R.Y)
        resCycles = 2
    }		

    return
}

// instruction implementations
func (cpu *Cpu) adc(addr int) {
    data := cpu.mem.read(addr)

    // Calculate auxiliary value
    t := cpu.ac + data + cpu.p.c

    // Set flags: overflow, sign, zero, and carry
    overflow := ((cpu.ac & M.BIT_7) != (t & M.BIT_7))
    if overflow {
    	cpu.p.v = 1
    } else {
    	cpu.p.v = 0
    }
    cpu.p.n = t
    cpu.p.z = t

    if (cpu.p.d == 1) {
        t = bcd(cpu.ac) + bcd(data) + cpu.p.c
        if (t > 99) {
        	cpu.p.c = 1
        } else {
        	cpu.p.c = 0
        }
    } else {
    	if (t > 255) {
        	cpu.p.c = 1
        } else {
        	cpu.p.c = 0
        }
    }

    // take the possible carry out
    cpu.ac = t & 0xFF
}

func (cpu *Cpu) and(int addr) {
    data := cpu.mem.read(addr)
    cpu.ac &= data

    // flags: sign, zero.
    cpu.p.n = cpu.ac
    cpu.p.z = cpu.ac
}

func (cpu *Cpu) asla() {
	carry := (cpu.ac & M.BIT_7) == M.BIT_7
	if carry {
		cpu.p.c = 1
	} else {
		cpu.p.c = 0
	}
    cpu.ac <<= 1

    cpu.p.n = cpu.ac
    cpu.p.z = cpu.ac
}

func (cpu *Cpu) asl(int addr) {
    data := cpu.mem.read(addr)

    carry := (data & M.BIT_7) == M.BIT_7
    if carry {
    	cpu.p.c = 1
    } else {
    	cpu.p.c = 0
    }
    data = (data << 1) & 0xFE

    cpu.p.n = data
    cpu.p.z = data

    cpu.mem.write(addr, data)
}

func (cpu *Cpu) bcc(int addr) bool {
    if cpu.p.c == 0 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) bcs(addr int) bool {
    if cpu.p.c == 1 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) beq(addr int) bool {
    if cpu.p.z == 1 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) bit(addr int) {
    data := cpu.mem.read(addr) & cpu.ac

	if data & M.BIT_6 != 0 {
		cpu.p.v = 1
	} else {
		cpu.p.v = 0
	}
    cpu.p.n = data
    cpu.p.z = data
}

func (cpu *Cpu) bmi(addr int) bool {
    if cpu.p.n == 1 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) bne(addr int) bool {
    if p.z == 0 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) bpl(addr int) bool {
    if p.n == 0 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) brk() {
    var l, h int

    // Even though the brk instruction is just one byte long, the pc is
    // incremented, meaning that the instruction after brk is ignored.
    cpu.pc++

    cpu.mem.write(cpu.sp, cpu.pc & 0xF0)
    cpu.mem.commit()
    cpu.sp--
    cpu.mem.write(cpu.sp, cpu.pc & 0xF)
    cpu.mem.commit()
    cpu.sp--
    cpu.mem.write(cpu.sp, cpu.p.b)
    cpu.sp--

    l = cpu.mem.read(0xFFFE)
    h = cpu.mem.read(0xFFFF) << 8

    cpu.pc = h | l
}

func (cpu *Cpu) bvc(addr int) bool {
    if cpu.p.v == 0 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) bvs(addr int) bool {
    if p.v == 1 {
        cpu.pc = addr
        return true
    }
    return false
}

func (cpu *Cpu) clc() {
    cpu.p.c = 0
}

func (cpu *Cpu) cld() {
    cpu.p.d = 0
}

func (cpu *Cpu) cli() {
    cpu.p.i = 0
}

func (cpu *Cpu) clv() {
    cpu.p.v = 0
}

// TODO: registers
func (cpu *Cpu) cmp(addr int, R r) {
    data := cpu.mem.read(addr)

    // Calculate auxiliary value
    t := 0
    switch (r) {
        case A:
            t = ac - data
			if cpu.ac >= data {
				cpu.p.c = 1
			} else {
				cpu.p.c = 0
			}

        case X:
            t = cpu.x - data
			if cpu.x >= data {
				p.c = 1
			} else {
				p.c =0
			}

        case Y:
            t = cpu.y - data
			if cpu.y >= data {
				p.c = 1
			} else {
				p.c = 0
			}
    }

    // Set flags
    cpu.p.n = t
    cpu.p.z = t
}

func (cpu *Cpu) dec(addr int) {
    data = cpu.mem.read(addr)

    // Decrement & AND 0xFF
    data = (data-1) & 0xFF
    cpu.mem.write(addr, data)

    cpu.p.n = t
    cpu.p.z = t
}

func (cpu *Cpu) decxy(R r) {
    switch (r) {
        case X:
            cpu.x = (cpu.x - 1) & 0xFF
            cpu.p.n = cpu.x
            cpu.p.z = cpu.x

        case Y:
            cpu.y = (cpu.y - 1) & 0xFF
            cpu.p.setN(cpu.y)
            cpu.p.setZ(cpu.y)
    }
}

func (cpu *Cpu) eor(addr int) {
    data := cpu.mem.read(addr)

    cpu.ac ^= data
    cpu.p.n = ac
    cpu.p.z = ac
}

func (cpu *Cpu) inc(addr int) {
    data := cpu.mem.read(addr)

	data++
    data &= 0xFF
    cpu.mem.write(addr, data)

    cpu.p.n = data
    cpu.p.z = data
}

func (cpu *Cpu) incxy(R r) {
    switch (r) {
        case X:
            cpu.x = (cpu.x + 1) & 0xFF
            cpu.p.n = x
            cpu.p.z = x

        case Y:
            cpu.y = (cpu.y + 1) & 0xFF
            cpu.p.n = y
            cpu.p.z = y
    }
}

func (cpu *Cpu) jmp(addr int) {
    cpu.pc = addr
}

func (cpu *Cpu) jsr(addr int) {
    t := cpu.pc - 1

    // Push PC onto the stack
    cpu.mem.write(cpu.sp, (t & 0xFF00) >> 8)
    cpu.mem.commit()
    cpu.sp--
    cpu.mem.write(cpu.sp, t & 0xFF)
    cpu.sp--

    // Jump
    cpu.pc = addr
}

func (cpu *Cpu) ldr(int addr, R r) {
    data := cpu.mem.read(addr)

    // One func (cpu *Cpu)tion for three different opcodes. Have to switch the register
    switch (r) {
        case A:
            cpu.ac = data
            cpu.p.n = ac
            cpu.p.z = ac
            
        case X:
            x = data
            cpu.p.n = x
            cpu.p.z = x

        case Y:
            y = data
            cpu.p.n = y
            cpu.p.z = y
    }
}

func (cpu *Cpu) lsra() {
    cpu.p.n = 0
    cpu.p.c = ((cpu.ac & M.BIT_0) == 0) ? 0 : 1
    cpu.ac = (cpu.ac >> 1) & 0x7F
    cpu.p.z = cpu.ac
}

func (cpu *Cpu) lsrm(int addr) {
    data := cpu.mem.read(addr)

    cpu.p.n = 0
	if data & M.BIT_0 == 0 {
		cpu.p.c = 0
	} else {
		cpu.p.c = 1
	}
    data = (data >> 1) & 0x7F
    cpu.p.z = data

    cpu.mem.write(addr, data)
}

func (cpu *Cpu) nop() {

}

func (cpu *Cpu) ora(addr int) {
    data = cpu.mem.read(addr)

    cpu.ac |= data
    cpu.p.n = data
    cpu.p.z = data
}

func (cpu *Cpu) pha() {
    cpu.mem.write(cpu.sp, cpu.ac)
    cpu.sp--
}

func (cpu *Cpu) php() {
    cpu.mem.write(cpu.sp, cpu.p.getProcessorStatus())
    cpu.sp--
}

func (cpu *Cpu) pla() {
    cpu.sp++
    cpu.ac = cpu.mem.read(cpu.sp)

    cpu.p.n = ac
    cpu.p.z = ac
}

func (cpu *Cpu) plp() {
    cpu.sp++
    cpu.p = cpu.mem.read(sp)
}

func (cpu *Cpu) rola() {
    // This opcode uses the carry to fill the LSB, and then sets the carry
    // according to the MSB of the rolled byte

    // Take from the byte what will be the future carry
	var t
	if cpu.ac & M.BIT_7 != 0 {
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
    cpu.p.z = cpu.ac
    cpu.p.n = cpu.ac
}

func (cpu *Cpu) rolm(addr int) {
    data := cpu.mem.read(addr)
    var t
	if data & M.BIT_7 != 0 {
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
    cpu.p.z = data
    cpu.p.n = data

    // Write to memory
    cpu.mem.write(addr, data)
}

func (cpu *Cpu) rora() {
    // This opcode uses the carry to fill the MSB, and then sets the carry
    // according to the LSB of the rolled byte

    // Take from the byte what will be the future carry
	var t
	if ac & M.BIT_0 != 0 {
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
    cpu.p.z = cpu.ac
    cpu.p.n = cpu.ac
}

func (cpu *Cpu) rorm(addr int) {
    data := cpu.mem.read(addr)
    var t
	if data & M.BIT_0 != 0 {
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
    cpu.p.z = data
    cpu.p.n = data

    // Write to memory
    cpu.mem.write(addr, data)
}

func (cpu *Cpu) rti() {
    var l, h

    cpu.sp--
    cpu.p = mem.read(cpu.sp)
    cpu.sp--
    l = cpu.mem.read(cpu.sp)
    cpu.sp--
    h = cpu.mem.read(cpu.sp)

    cpu.pc = (h << 8) | l
}

func (cpu *Cpu) rts() {
    var l, h

    cpu.sp++
    l = cpu.mem.read(cpu.sp)
    cpu.sp++
    h = cpu.mem.read(cpu.sp)

    cpu.pc = ((h << 8) | l) + 1
}

func (cpu *Cpu) sbc(int addr) {
    data := cpu.mem.read(addr)
    var t

    // If decimal mode is on...
    if (cpu.p.d == 1) {
        // When using SBC, the code should have used SEC to set the carry
        // before. This is to make sure that, if we need to borrow, there is
        // something to borrow.
		var negcarry
		if cpu.p.c & M.BIT_0 != 0 {
			negcarry = 0
		else {
			negcarry = 1
		}
        t = bcd(cpu.ac) - bcd(data) - negcarry

		if t > 99 || t < 0 {
			cpu.p.v = 1
		} else {
			cpu.p.v = 0
		}
    } else {
        t = cpu.ac - data - (((p.c & M.BIT_0) != 0) ? 0 : 1)
		var negcarry
		if cpu.p.c & M.BIT_0 != 0 {
			negcarry = 0
		} else {
			negcarry = 1
		}
		
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

func (cpu *Cpu) sec() {
    cpu.p.c = 1
}

func (cpu *Cpu) sed() {
    cpu.p.d = 1
}

func (cpu *Cpu) sei() {
    cpu.p.i = 1
}

func (cpu *Cpu) st(addr int, R r) {
    switch (r) {
        case A:
            cpu.mem.write(addr, cpu.ac)

        case X:
            cpu.mem.write(addr, cpu.x)

        case Y:
            cpu.mem.write(addr, cpu.y)
    }
}

func (cpu *Cpu) taxy(R r) {
    switch (r) {
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

func (cpu *Cpu) tsx() {
    cpu.x = cpu.sp
    cpu.p.n = cpu.x
    cpu.p.z = cpu.x
}

func (cpu *Cpu) txya(R r) {
    switch (r) {
        case X:
            cpu.ac = cpu.x
            
        case Y:
            cpu.ac = cpu.y
    }

    cpu.p.n = cpu.ac
    cpu.p.z = cpu.ac
}

func (cpu *Cpu) txs() {
    cpu.sp = cpu.x
}

// helper functions

// returns the bcd equivalent of the given number
func bcd(n int) int {
    return (n & 0xF) + ((n & 0xF0) * 10);
}

