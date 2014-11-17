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
type Mem struct {
	m 		[1 << 16]int
}

func (mem *Mem) read(i int) (value int) {
	value = mem.m[i]
	return
}

func (cpu *Cpu) execute() (resCycles int) {
	// grab current instruction and increment pc
	inst := cpu.mem.mem[cpu.pc]
	cpu.pc++

	switch inst {
	case 0x69:
        adc(imm())
        resCycles = 2

    case 0x65:
        adc(zp())
        resCycles = 3

    case 0x75:
        adc(zpx())
        resCycles = 4

    case 0x6D:
        adc(abs())
        resCycles = 4

    case 0x7D:
        adc(abx())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x79:
        adc(aby())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x61:
        adc(indx())
        resCycles = 6

    case 0x71:
        adc(indy())
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // AND's
    case 0x29:
        and(imm())
        resCycles = 2

    case 0x25:
        and(zp())
        resCycles = 2

    case 0x35:
        and(zpx())
        resCycles = 3

    case 0x2D:
        and(abs())
        resCycles = 4

    case 0x3D:
        and(abx())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x39:
        and(aby())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x21:
        and(indx())
        resCycles = 6

    case 0x31:
        and(indy())
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // ASL's
    case 0x0A:
        asla()
        resCycles = 2

    case 0x06:
        asl(zp())
        resCycles = 5

    case 0x16:
        asl(zpx())
        resCycles = 6

    case 0x0E:
        asl(abs())
        resCycles = 6

    case 0x1E:
        asl(abx())
        resCycles = 7

    // BCC
    case 0x90:
    	if bcc(rel()) {
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
    	if bcs(rel()) {
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
    	if beq(rel()) {
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
        bit(zp())
        resCycles = 3

    case 0x2C:
        bit(abs())
        resCycles = 4

    // BMI
    case 0x30:
    	if bmi(rel()) {
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
    	if bne(rel()) {
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
    	if bpl(rel()) {
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
        brk()
        resCycles = 7

    // BVC
    case 0x50:
    	if bvc(rel()) {
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
    	if bvs(rel()) {
    		if cpu.pbCrossed {
    			resCycles = 4
    		} else {
    			resCycles = 3
    		}
    	} else {
    		resCycles = 2
    	}

    case 0x18:
        clc()
        resCycles = 2

    case 0xD8:
        cld()
        resCycles = 2

    case 0x58:
        cli()
        resCycles = 2

    case 0xB8:
        clv()
        
    // CMP
    case 0xC9:
        cmp(imm(), R.A)
        resCycles = 2

    case 0xC5:
        cmp(zp(), R.A)
        resCycles = 3

    case 0xD5:
        cmp(zpx(), R.A)
        resCycles = 4

    case 0xCD:
        cmp(abs(), R.A)
        resCycles = 4

    case 0xDD:
        cmp(abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xD9:
        cmp(aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xC1:
        cmp(indx(), R.A)
        resCycles = 6

    case 0xD1:
        cmp(indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // CPX
    case 0xE0:
        cmp(imm(), R.X)
        resCycles = 2

    case 0xE4:
        cmp(zp(), R.X)
        resCycles = 3

    case 0xEC:
        cmp(abs(), R.X)
        resCycles = 4

    // CPY
    case 0xC0:
        cmp(imm(), R.Y)
        resCycles = 2

    case 0xC4:
        cmp(zp(), R.Y)
        resCycles = 3

    case 0xCC:
        cmp(abs(), R.Y)
        resCycles = 4

    // DEC
    case 0xC6:
        dec(zp())
        resCycles = 5

    case 0xD6:
        dec(zpx())
        resCycles = 6

    case 0xCE:
        dec(abs())
        resCycles = 6

    case 0xDE:
        dec(abx())
        resCycles = 7

    // DEX
    case 0xCA:
        decxy(R.X)
        resCycles = 2

    // DEY
    case 0x88:
        decxy(R.Y)
        resCycles = 2

    // EOR
    case 0x49:
        eor(imm())
        resCycles = 2

    case 0x45:
        eor(zp())
        resCycles = 3

    case 0x55:
        eor(zpx())
        resCycles = 4

    case 0x4D:
        eor(abs())
        resCycles = 4

    case 0x5D:
        eor(abx())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x59:
        eor(aby())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x41:
        eor(indx())
        resCycles = 6

    case 0x51:
        eor(indy())
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // INC
    case 0xE6:
        inc(zp())
        resCycles = 5

    case 0xF6:
        inc(zpx())
        resCycles = 6

    case 0xEE:
        inc(abs())
        resCycles = 6

    case 0xFE:
        inc(abx())
        resCycles = 7

    // INX
    case 0xE8:
        incxy(R.X)
        resCycles = 2

    // INY
    case 0xC8:
        incxy(R.Y)
        resCycles = 2

    // JMP
    case 0x4C:
        jmp(abs())
        resCycles = 3

    case 0x6C:
        jmp(ind())
        resCycles = 5

    // JSR
    case 0x20:
        jsr(abs())
        resCycles = 6

    // LDA
    case 0xA9:
        ldr(imm(), R.A)
        resCycles = 2

    case 0xA5:
        ldr(zp(), R.A)
        resCycles = 3

    case 0xB5:
        ldr(zpx(), R.A)
        resCycles = 4

    case 0xAD:
        ldr(abs(), R.A)
        resCycles = 4

    case 0xBD:
        ldr(abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xB9:
        ldr(aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xA1:
        ldr(indx(), R.A)
        resCycles = 6

    case 0xB1:
        ldr(indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // LDX
    case 0xA2:
        ldr(imm(), R.X)
        resCycles = 2

    case 0xA6:
        ldr(zp(), R.X)
        resCycles = 3

    case 0xB6:
        ldr(zpy(), R.X)
        resCycles = 4

    case 0xAE:
        ldr(abs(), R.X)
        resCycles = 4

    case 0xBE:
        ldr(aby(), R.X)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    // LDY
    case 0xA0:
        ldr(imm(), R.Y)
        resCycles = 2

    case 0xA4:
        ldr(zp(), R.Y)
        resCycles = 3

    case 0xB4:
        ldr(zpx(), R.Y)
        resCycles = 4

    case 0xAC:
        ldr(abs(), R.Y)
        resCycles = 4

    case 0xBC:
        ldr(abx(), R.Y)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    // LSR
    case 0x4A:
        lsra()
        resCycles = 2

    case 0x46:
        lsrm(zp())
        resCycles = 5

    case 0x56:
        lsrm(zpx())
        resCycles = 6

    case 0x4E:
        lsrm(abs())
        resCycles = 6

    case 0x5E:
        lsrm(abx())
        resCycles = 7

    // NOP
    case 0xEA:
        nop()
        resCycles = 2

    // ORA
    case 0x09:
        ora(imm())
        resCycles = 2

    case 0x05:
        ora(imm())
        resCycles = 2

    case 0x15:
        ora(imm())
        resCycles = 3

    case 0x1D:
        ora(imm())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x19:
        ora(imm())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x01:
        ora(imm())
        resCycles = 6

    case 0x11:
        ora(imm())
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // PHA
    case 0x48:
        pha()
        resCycles = 3

    // PHP
    case 0x08:
        php()
        resCycles = 3

    // PLA
    case 0x68:
        pla()
        resCycles = 4

    // PLP
    case 0x28:
        plp()
        resCycles = 4

    // ROL
    case 0x2A:
        rola()
        resCycles = 2

    case 0x6E:
        rorm(abs())
        resCycles = 6

    case 0x7E:
        rorm(abx())
        resCycles = 7

    // RTI
    case 0x40:
        rti()
        resCycles = 6

    // RTS
    case 0x60:
        rts()
        resCycles = 6

    // SBC
    case 0xE9:
        sbc(imm())
        resCycles = 2

    case 0xE5:
        sbc(zp())
        resCycles = 3

    case 0xF5:
        sbc(zpx())
        resCycles = 4

    case 0xED:
        sbc(abs())
        resCycles = 4

    case 0xFD:
        sbc(abx())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xF9:
        sbc(aby())
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0xE1:
        sbc(indx())
        resCycles = 6

    case 0xF1:
        sbc(indy())
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }

    // SEC
    case 0x38:
        sec()
        resCycles = 2

    // SED
    case 0xF8:
        sed()
        resCycles = 2

    // SEI
    case 0x78:
        sei()
        resCycles = 2

    // STA
    case 0x85:
        st(zp(), R.A)
        resCycles = 3

    case 0x95:
        st(zpx(), R.A)
        resCycles = 4

    case 0x8D:
        st(abs(), R.A)
        resCycles = 4

    case 0x9D:
        st(abx(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x99:
        st(aby(), R.A)
        if cpu.pbCrossed {
        	resCycles = 5
        } else {
        	resCycles = 4
        }

    case 0x81:
        st(indx(), R.A)
        resCycles = 6

    case 0x91:
        st(indy(), R.A)
        if cpu.pbCrossed {
        	resCycles = 6
        } else {
        	resCycles = 5
        }
        
    // STX
    case 0x86:
        st(zp(), R.X)
        resCycles = 3

    case 0x96:
        st(zpy(), R.X)
        resCycles = 4

    case 0x8E:
        st(abs(), R.X)
        resCycles = 4

    // STY
    case 0x84:
        st(zp(), R.Y)
        resCycles = 3

    case 0x94:
        st(zpx(), R.Y)
        resCycles = 4

    case 0x8C:
        st(abs(), R.Y)
        resCycles = 4

    // TAX
    case 0xAA:
        taxy(R.X)
        resCycles = 2

    // TAY
    case 0xA8:
        taxy(R.Y)
        resCycles = 2

    // TSX
    case 0xBA:
        tsx()
        resCycles = 2

    // TXA
    case 0x8A:
        txya(R.X)
        resCycles = 2

    // TXS
    case 0x9A:
        txs()
        resCycles = 2

    // TYA
    case 0x98:
        txya(R.Y)
        resCycles = 2
    }		

    return
}

// instruction implementations
func (cpu *Cpu) adc(int addr) {
    data := cpu.mem.read(addr)

    // Calculate auxiliary value
    t := ac + data + cpu.p.c

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

    carry := ((data & M.BIT_7) == M.BIT_7)
    if carry {
    	cpu.p.c = 1
    } else {
    	cpu.p.c = 0
    }
    data = (data << 1) & 0xFE

    p.setN(data)
    p.setZ(data)

    mem.write(addr, data)
}

func (cpu *Cpu) bcc(int addr) bool {
    if (p.c == 0) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) bcs(int addr) bool {
    if (p.c == 1) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) beq(int addr) bool {
    if (p.z == 1) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) bit(int addr) {
    int data = mem.read(addr) & ac

    p.setN(data)
    p.v = ((data & M.BIT_6) != 0) ? 1 : 0
    p.setZ(data)
}

func (cpu *Cpu) bmi(int addr) bool {
    if (p.n == 1) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) bne(int addr) bool {
    if (p.z == 0) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) bpl(int addr) bool {
    if (p.n == 0) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) brk() {
    int l, h

    // Even though the brk instruction is just one byte long, the pc is
    // incremented, meaning that the instruction after brk is ignored.
    pc++

    mem.write(sp, pc & 0xF0)
    mem.commit()
    sp--
    mem.write(sp, pc & 0xF)
    mem.commit()
    sp--
    mem.write(sp, p.b)
    sp--

    l = mem.read(0xFFFE)
    h = mem.read(0xFFFF) << 8

    pc = h | l
}

func (cpu *Cpu) bvc(int addr) bool {
    if (p.v == 0) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) bvs(int addr) bool {
    if (p.v == 1) {
        pc = addr
        return true
    }

    return false
}

func (cpu *Cpu) clc() {
    p.c = 0
}

func (cpu *Cpu) cld() {
    p.d = 0
}

func (cpu *Cpu) cli() {
    p.i = 0
}

func (cpu *Cpu) clv() {
    p.v = 0
}

func (cpu *Cpu) cmp(int addr, R r) {
    int data = mem.read(addr)

    // Calculate auxiliary value
    int t = 0
    switch (r) {
        case A:
            t = ac - data
            p.c = (ac >= data) ? 1 : 0
            

        case X:
            t = x - data
            p.c = (x >= data) ? 1 : 0
            

        case Y:
            t = y - data
            p.c = (y >= data) ? 1 : 0
            

        default:
            
    }

    // Set flags
    p.setN(t)
    p.setZ(t)
}

func (cpu *Cpu) dec(int addr) {
    int data = mem.read(addr)

    // Decrement & AND 0xFF
    data = --data & 0xFF
    mem.write(addr, data)

    p.setN(data)
    p.setZ(data)
}

func (cpu *Cpu) decxy(R r) {
    switch (r) {
        case X:
            x = (x - 1) & 0xFF
            p.setN(x)
            p.setZ(x)
            

        case Y:
            y = (y - 1) & 0xFF
            p.setN(y)
            p.setZ(y)
            

        default:
            
    }
}

func (cpu *Cpu) eor(int addr) {
    int data = mem.read(addr)

    ac ^= data
    p.setN(ac)
    p.setZ(ac)
}

func (cpu *Cpu) inc(int addr) {
    int data = mem.read(addr)

    data = ++data & 0xFF
    mem.write(addr, data)

    p.setN(data)
    p.setZ(data)
}

func (cpu *Cpu) incxy(R r) {
    switch (r) {
        case X:
            x = (x + 1) & 0xFF
            p.setN(x)
            p.setZ(x)
            

        case Y:
            y = (y + 1) & 0xFF
            p.setN(y)
            p.setZ(y)
            

        default:
            
    }
}

func (cpu *Cpu) jmp(int addr) {
    pc = addr
}

func (cpu *Cpu) jsr(int addr) {
    int t = pc - 1

    // Push PC onto the stack
    mem.write(sp, (t & 0xFF00) >> 8)
    mem.commit()
    sp--
    mem.write(sp, t & 0xFF)
    sp--

    // Jump
    pc = addr
}

func (cpu *Cpu) ldr(int addr, R r) {
    int data = mem.read(addr)

    // One func (cpu *Cpu)tion for three different opcodes. Have to switch the register
    switch (r) {
        case A:
            ac = data
            p.setN(ac)
            p.setZ(ac)
            

        case X:
            x = data
            p.setN(x)
            p.setZ(x)
            

        case Y:
            y = data
            p.setN(y)
            p.setZ(y)
            

        default:
            
    }
}

func (cpu *Cpu) lsra() {
    p.n = 0
    p.c = ((ac & M.BIT_0) == 0) ? 0 : 1
    ac = (ac >> 1) & 0x7F
    p.setZ(ac)
}

func (cpu *Cpu) lsrm(int addr) {
    int data = mem.read(addr)

    p.n = 0
    p.c = ((data & M.BIT_0) == 0) ? 0 : 1
    data = (data >> 1) & 0x7F
    p.setZ(data)

    mem.write(addr, data)
}

func (cpu *Cpu) nop() {

}

func (cpu *Cpu) ora(int addr) {
    int data = mem.read(addr)

    ac |= data
    p.setN(data)
    p.setZ(data)
}

func (cpu *Cpu) pha() {
    mem.write(sp, ac)
    sp--
}

func (cpu *Cpu) php() {
    mem.write(sp, p.getProcessorStatus())
    sp--
}

func (cpu *Cpu) pla() {
    sp++
    ac = mem.read(sp)

    p.setN(ac)
    p.setZ(ac)
}

func (cpu *Cpu) plp() {
    sp++
    p.setProcessorStatus(mem.read(sp))
}

func (cpu *Cpu) rola() {
    // This opcode uses the carry to fill the LSB, and then sets the carry
    // according to the MSB of the rolled byte

    // Take from the byte what will be the future carry
    int t = ((ac & M.BIT_7) != 0) ? 1 : 0

    // Rotate left and &
    ac = (ac << 1) & 0xFE
    // Set LSB with the carry value from before the operation
    ac |= p.c
    // Set the next carry
    p.c = t
    // Set flags
    p.setZ(ac)
    p.setN(ac)
}

func (cpu *Cpu) rolm(int addr) {
    int data = mem.read(addr)
    int t = ((data & M.BIT_7) != 0) ? 1 : 0

    // Rotate left and &
    data = (data << 1) & 0xFE
    // Set LSB with the carry value from before the operation
    data |= p.c
    // Set the next carry
    p.c = t
    // Set flags
    p.setZ(data)
    p.setN(data)

    // Write to memory
    mem.write(addr, data)
}

func (cpu *Cpu) rora() {
    // This opcode uses the carry to fill the MSB, and then sets the carry
    // according to the LSB of the rolled byte

    // Take from the byte what will be the future carry
    int t = ((ac & M.BIT_0) != 0) ? 1 : 0

    // Rotate right and &
    ac = (ac >> 1) & 0x7F
    // Set MSB with the carry value from before the operation
    ac |= (((p.c == 1) ? 0x80 : 0x00))
    // Set the next carry
    p.c = t
    // Set flags
    p.setZ(ac)
    p.setN(ac)
}

func (cpu *Cpu) rorm(int addr) {
    int data = mem.read(addr)
    int t = ((data & M.BIT_0) != 0) ? 1 : 0

    // Rotate right and &
    data = (data >> 1) & 0x7F
    // Set LSB with the carry value from before the operation
    data |= (((p.c == 1) ? 0x80 : 0x00))
    // Set the next carry
    p.c = t
    // Set flags
    p.setZ(data)
    p.setN(data)

    // Write to memory
    mem.write(addr, data)
}

func (cpu *Cpu) rti() {
    int l, h

    sp--
    p.setProcessorStatus(mem.read(sp))
    sp--
    l = mem.read(sp)
    sp--
    h = mem.read(sp)

    pc = (h << 8) | l
}

func (cpu *Cpu) rts() {
    int l, h

    sp++
    l = mem.read(sp)
    sp++
    h = mem.read(sp)

    pc = ((h << 8) | l) + 1
}

func (cpu *Cpu) sbc(int addr) {
    int data = mem.read(addr)
    int t

    // If decimal mode is on...
    if (p.d == 1) {
        // When using SBC, the code should have used SEC to set the carry
        // before. This is to make sure that, if we need to borrow, there is
        // something to borrow.
        t = bcd(ac) - bcd(data) - (((p.c & M.BIT_0) != 0) ? 0 : 1)
        p.v = (t > 99 || t < 0) ? 1 : 0
    } else {
        t = ac - data - (((p.c & M.BIT_0) != 0) ? 0 : 1)
        p.v = (t > 127 || t < -128) ? 1 : 0
    }

    // Set the flags
    p.c = (t >= 0) ? 1 : 0
    p.setN(t)
    p.setZ(t)

    // Write the result (ANDed, just in case it overflowed)
    ac = t & 0xFF
}

func (cpu *Cpu) sec() {
    p.c = 1
}

func (cpu *Cpu) sed() {
    p.d = 1
}

func (cpu *Cpu) sei() {
    p.i = 1
}

func (cpu *Cpu) st(int addr, R r) {
    switch (r) {
        case A:
            mem.write(addr, ac)
            

        case X:
            mem.write(addr, x)
            

        case Y:
            mem.write(addr, y)
            

        default:
            
    }
}

func (cpu *Cpu) taxy(R r) {
    switch (r) {
        case X:
            x = ac
            p.setN(x)
            p.setZ(x)
            

        case Y:
            y = ac
            p.setN(y)
            p.setZ(y)
            

        default:
            
    }
}

func (cpu *Cpu) tsx() {
    x = sp
    p.setN(x)
    p.setZ(x)
}

func (cpu *Cpu) txya(R r) {
    switch (r) {
        case X:
            ac = x
            

        case Y:
            ac = y
            

        default:
            
    }

    p.setN(ac)
    p.setZ(ac)
}

func (cpu *Cpu) txs() {
    sp = x
}

// helper functions

// returns the bcd equivalent of the given number
func bcd(n int) int {
    return (n & 0xF) + ((n & 0xF0) * 10);
}

