/*
@author: sk
@date: 2024/2/24
*/
package nes5

import "reflect"

const (
	C = 1 << 0 // 进位标记
	Z = 1 << 1 // 是否为0标记
	I = 1 << 2 // 屏蔽普通中断标记
	D = 1 << 3 // Unused 十进制模式
	B = 1 << 4 // Break标记
	U = 1 << 5 // Unused
	V = 1 << 6 // 计算溢出标记
	N = 1 << 7 // 负数标记
)

type Instruction struct { // 一条指令的抽象
	Name     string
	Operate  func() bool
	AddrMode func() bool
	Cycles   uint8
}

type Olc6502 struct {
	Bus     *Bus
	A, X, Y uint8
	StackP  uint8
	Pc      uint16
	Status  uint8
	Fetched uint8
	AddrAbs uint16
	AddrRel uint16
	OpCode  uint8
	Cycles  uint8
	Lookup  []*Instruction
}

func NewOlc6502() *Olc6502 {
	res := &Olc6502{}
	// 各种指令的映射关系
	res.Lookup = []*Instruction{
		{"BRK", res.BRK, res.IMM, 7}, {"ORA", res.ORA, res.IZX, 6}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 3}, {"ORA", res.ORA, res.ZP0, 3}, {"ASL", res.ASL, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"PHP", res.PHP, res.IMP, 3}, {"ORA", res.ORA, res.IMM, 2}, {"ASL", res.ASL, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"???", res.NOP, res.IMP, 4}, {"ORA", res.ORA, res.ABS, 4}, {"ASL", res.ASL, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BPL", res.BPL, res.REL, 2}, {"ORA", res.ORA, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"ORA", res.ORA, res.ZPX, 4}, {"ASL", res.ASL, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"CLC", res.CLC, res.IMP, 2}, {"ORA", res.ORA, res.ABY, 4}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"ORA", res.ORA, res.ABX, 4}, {"ASL", res.ASL, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
		{"JSR", res.JSR, res.ABS, 6}, {"AND", res.AND, res.IZX, 6}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"BIT", res.BIT, res.ZP0, 3}, {"AND", res.AND, res.ZP0, 3}, {"ROL", res.ROL, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"PLP", res.PLP, res.IMP, 4}, {"AND", res.AND, res.IMM, 2}, {"ROL", res.ROL, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"BIT", res.BIT, res.ABS, 4}, {"AND", res.AND, res.ABS, 4}, {"ROL", res.ROL, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BMI", res.BMI, res.REL, 2}, {"AND", res.AND, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"AND", res.AND, res.ZPX, 4}, {"ROL", res.ROL, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"SEC", res.SEC, res.IMP, 2}, {"AND", res.AND, res.ABY, 4}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"AND", res.AND, res.ABX, 4}, {"ROL", res.ROL, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
		{"RTI", res.RTI, res.IMP, 6}, {"EOR", res.EOR, res.IZX, 6}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 3}, {"EOR", res.EOR, res.ZP0, 3}, {"LSR", res.LSR, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"PHA", res.PHA, res.IMP, 3}, {"EOR", res.EOR, res.IMM, 2}, {"LSR", res.LSR, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"JMP", res.JMP, res.ABS, 3}, {"EOR", res.EOR, res.ABS, 4}, {"LSR", res.LSR, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BVC", res.BVC, res.REL, 2}, {"EOR", res.EOR, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"EOR", res.EOR, res.ZPX, 4}, {"LSR", res.LSR, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"CLI", res.CLI, res.IMP, 2}, {"EOR", res.EOR, res.ABY, 4}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"EOR", res.EOR, res.ABX, 4}, {"LSR", res.LSR, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
		{"RTS", res.RTS, res.IMP, 6}, {"ADC", res.ADC, res.IZX, 6}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 3}, {"ADC", res.ADC, res.ZP0, 3}, {"ROR", res.ROR, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"PLA", res.PLA, res.IMP, 4}, {"ADC", res.ADC, res.IMM, 2}, {"ROR", res.ROR, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"JMP", res.JMP, res.IND, 5}, {"ADC", res.ADC, res.ABS, 4}, {"ROR", res.ROR, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BVS", res.BVS, res.REL, 2}, {"ADC", res.ADC, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"ADC", res.ADC, res.ZPX, 4}, {"ROR", res.ROR, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"SEI", res.SEI, res.IMP, 2}, {"ADC", res.ADC, res.ABY, 4}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"ADC", res.ADC, res.ABX, 4}, {"ROR", res.ROR, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
		{"???", res.NOP, res.IMP, 2}, {"STA", res.STA, res.IZX, 6}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 6}, {"STY", res.STY, res.ZP0, 3}, {"STA", res.STA, res.ZP0, 3}, {"STX", res.STX, res.ZP0, 3}, {"???", res.XXX, res.IMP, 3}, {"DEY", res.DEY, res.IMP, 2}, {"???", res.NOP, res.IMP, 2}, {"TXA", res.TXA, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"STY", res.STY, res.ABS, 4}, {"STA", res.STA, res.ABS, 4}, {"STX", res.STX, res.ABS, 4}, {"???", res.XXX, res.IMP, 4},
		{"BCC", res.BCC, res.REL, 2}, {"STA", res.STA, res.IZY, 6}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 6}, {"STY", res.STY, res.ZPX, 4}, {"STA", res.STA, res.ZPX, 4}, {"STX", res.STX, res.ZPY, 4}, {"???", res.XXX, res.IMP, 4}, {"TYA", res.TYA, res.IMP, 2}, {"STA", res.STA, res.ABY, 5}, {"TXS", res.TXS, res.IMP, 2}, {"???", res.XXX, res.IMP, 5}, {"???", res.NOP, res.IMP, 5}, {"STA", res.STA, res.ABX, 5}, {"???", res.XXX, res.IMP, 5}, {"???", res.XXX, res.IMP, 5},
		{"LDY", res.LDY, res.IMM, 2}, {"LDA", res.LDA, res.IZX, 6}, {"LDX", res.LDX, res.IMM, 2}, {"???", res.XXX, res.IMP, 6}, {"LDY", res.LDY, res.ZP0, 3}, {"LDA", res.LDA, res.ZP0, 3}, {"LDX", res.LDX, res.ZP0, 3}, {"???", res.XXX, res.IMP, 3}, {"TAY", res.TAY, res.IMP, 2}, {"LDA", res.LDA, res.IMM, 2}, {"TAX", res.TAX, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"LDY", res.LDY, res.ABS, 4}, {"LDA", res.LDA, res.ABS, 4}, {"LDX", res.LDX, res.ABS, 4}, {"???", res.XXX, res.IMP, 4},
		{"BCS", res.BCS, res.REL, 2}, {"LDA", res.LDA, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 5}, {"LDY", res.LDY, res.ZPX, 4}, {"LDA", res.LDA, res.ZPX, 4}, {"LDX", res.LDX, res.ZPY, 4}, {"???", res.XXX, res.IMP, 4}, {"CLV", res.CLV, res.IMP, 2}, {"LDA", res.LDA, res.ABY, 4}, {"TSX", res.TSX, res.IMP, 2}, {"???", res.XXX, res.IMP, 4}, {"LDY", res.LDY, res.ABX, 4}, {"LDA", res.LDA, res.ABX, 4}, {"LDX", res.LDX, res.ABY, 4}, {"???", res.XXX, res.IMP, 4},
		{"CPY", res.CPY, res.IMM, 2}, {"CMP", res.CMP, res.IZX, 6}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"CPY", res.CPY, res.ZP0, 3}, {"CMP", res.CMP, res.ZP0, 3}, {"DEC", res.DEC, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"INY", res.INY, res.IMP, 2}, {"CMP", res.CMP, res.IMM, 2}, {"DEX", res.DEX, res.IMP, 2}, {"???", res.XXX, res.IMP, 2}, {"CPY", res.CPY, res.ABS, 4}, {"CMP", res.CMP, res.ABS, 4}, {"DEC", res.DEC, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BNE", res.BNE, res.REL, 2}, {"CMP", res.CMP, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"CMP", res.CMP, res.ZPX, 4}, {"DEC", res.DEC, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"CLD", res.CLD, res.IMP, 2}, {"CMP", res.CMP, res.ABY, 4}, {"NOP", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"CMP", res.CMP, res.ABX, 4}, {"DEC", res.DEC, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
		{"CPX", res.CPX, res.IMM, 2}, {"SBC", res.SBC, res.IZX, 6}, {"???", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"CPX", res.CPX, res.ZP0, 3}, {"SBC", res.SBC, res.ZP0, 3}, {"INC", res.INC, res.ZP0, 5}, {"???", res.XXX, res.IMP, 5}, {"INX", res.INX, res.IMP, 2}, {"SBC", res.SBC, res.IMM, 2}, {"NOP", res.NOP, res.IMP, 2}, {"???", res.SBC, res.IMP, 2}, {"CPX", res.CPX, res.ABS, 4}, {"SBC", res.SBC, res.ABS, 4}, {"INC", res.INC, res.ABS, 6}, {"???", res.XXX, res.IMP, 6},
		{"BEQ", res.BEQ, res.REL, 2}, {"SBC", res.SBC, res.IZY, 5}, {"???", res.XXX, res.IMP, 2}, {"???", res.XXX, res.IMP, 8}, {"???", res.NOP, res.IMP, 4}, {"SBC", res.SBC, res.ZPX, 4}, {"INC", res.INC, res.ZPX, 6}, {"???", res.XXX, res.IMP, 6}, {"SED", res.SED, res.IMP, 2}, {"SBC", res.SBC, res.ABY, 4}, {"NOP", res.NOP, res.IMP, 2}, {"???", res.XXX, res.IMP, 7}, {"???", res.NOP, res.IMP, 4}, {"SBC", res.SBC, res.ABX, 4}, {"INC", res.INC, res.ABX, 7}, {"???", res.XXX, res.IMP, 7},
	}
	return res
}

func (o *Olc6502) ConnectBus(bus *Bus) {
	o.Bus = bus
}

func (o *Olc6502) Read(addr uint16) uint8 {
	return o.Bus.CpuRead(addr, false)
}

func (o *Olc6502) Write(addr uint16, data uint8) {
	o.Bus.CpuWrite(addr, data)
}

func (o *Olc6502) Clock() {
	if o.Cycles == 0 {
		o.OpCode = o.Read(o.Pc)
		o.Pc++
		instruction := o.Lookup[o.OpCode]
		o.Cycles = instruction.Cycles
		// 必须先拿出来确保运行防止短路
		add1 := instruction.AddrMode()
		add2 := instruction.Operate()
		if add1 && add2 {
			o.Cycles++
		}
	}
	o.Cycles--
}

func (o *Olc6502) Reset() {
	// 重置寄存器
	o.A = 0
	o.X = 0
	o.Y = 0
	o.StackP = 0xFD // 栈基地址
	o.Status = U
	// 默认reset程序跳转的地址
	o.AddrAbs = 0xFFFC
	lo := uint16(o.Read(o.AddrAbs + 0))
	hi := uint16(o.Read(o.AddrAbs + 1))
	o.Pc = lo | (hi << 8)
	o.AddrAbs = 0
	o.AddrRel = 0
	o.Fetched = 0
	o.Cycles = 8
}

func (o *Olc6502) Irq() {
	if o.GetFlag(I) {
		return
	}
	// 存储程序计数器
	o.Write(0x0100+uint16(o.StackP), uint8((o.Pc>>8)&0xFF))
	o.StackP--
	o.Write(0x0100+uint16(o.StackP), uint8(o.Pc&0xFF))
	o.StackP--
	// 存储 cpu状态
	o.SetFlag(B, false)
	o.SetFlag(U, true)
	o.SetFlag(I, true)
	o.Write(0x0100+uint16(o.StackP), o.Status)
	o.StackP--
	// 跳转到中断处理程序
	o.AddrAbs = 0xFFFE
	lo := uint16(o.Read(o.AddrAbs + 0))
	hi := uint16(o.Read(o.AddrAbs + 1))
	o.Pc = (hi << 8) | lo
	o.Cycles = 7
}

func (o *Olc6502) Nmi() {
	// 存储程序计数器
	o.Write(0x0100+uint16(o.StackP), uint8((o.Pc>>8)&0xFF))
	o.StackP--
	o.Write(0x0100+uint16(o.StackP), uint8(o.Pc&0xFF))
	o.StackP--
	// 存储 cpu状态
	o.SetFlag(B, false)
	o.SetFlag(U, true)
	o.SetFlag(I, true)
	o.Write(0x0100+uint16(o.StackP), o.Status)
	o.StackP--
	// 跳转到中断处理程序,注意于Irq相比处理中断程序的位置不同，且周期不同
	o.AddrAbs = 0xFFFA
	lo := uint16(o.Read(o.AddrAbs + 0))
	hi := uint16(o.Read(o.AddrAbs + 1))
	o.Pc = (hi << 8) | lo
	o.Cycles = 8
}

func (o *Olc6502) Fetch() uint8 {
	if !AddrEq(o.Lookup[o.OpCode].AddrMode, o.IMP) {
		o.Fetched = o.Read(o.AddrAbs)
	}
	return o.Fetched
}

func (o *Olc6502) SetFlag(key uint8, flag bool) {
	if flag {
		o.Status |= key
	} else {
		o.Status &= ^key
	}
}

func (o *Olc6502) GetFlag(key uint8) bool {
	return (o.Status & key) > 0
}

func (o *Olc6502) Complete() bool { // 当前指令是否执行完成
	return o.Cycles == 0
}

// 反编译指定范围内的代码
func (o *Olc6502) Disassemble(start, end uint16) map[uint16]string {
	res := make(map[uint16]string, end-start)
	for addr := start; addr < end; {
		currAddr := addr
		opCode := o.Bus.CpuRead(addr, true)
		addr++
		instruction := o.Lookup[opCode]
		curr := "$" + Format(currAddr, 4) + ": " + instruction.Name + " "
		switch reflect.ValueOf(instruction.AddrMode).Pointer() {
		case reflect.ValueOf(o.IMP).Pointer():
			curr += "{IMP}"
		case reflect.ValueOf(o.IMM).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{IMM} 0x" + Format(value, 2)
		case reflect.ValueOf(o.ZP0).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ZP0} 0x" + Format(value, 2)
		case reflect.ValueOf(o.ZPX).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ZPX} 0x" + Format(value, 2)
		case reflect.ValueOf(o.ZPY).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ZPY} 0x" + Format(value, 2)
		case reflect.ValueOf(o.IZX).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{IZX} 0x" + Format(value, 2)
		case reflect.ValueOf(o.IZY).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{IZY} 0x" + Format(value, 2)
		case reflect.ValueOf(o.ABS).Pointer():
			lo := o.Bus.CpuRead(addr, true)
			addr++
			hi := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ABS} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(o.ABX).Pointer():
			lo := o.Bus.CpuRead(addr, true)
			addr++
			hi := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ABX} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(o.ABY).Pointer():
			lo := o.Bus.CpuRead(addr, true)
			addr++
			hi := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{ABY} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(o.IND).Pointer():
			lo := o.Bus.CpuRead(addr, true)
			addr++
			hi := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{IND} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(o.REL).Pointer():
			value := o.Bus.CpuRead(addr, true)
			addr++
			curr += "{REL} 0x" + Format(value, 2)
		}
		res[currAddr] = curr
	}
	return res
}

//=============AddrMode=============
// 一共 256页，每页256b 组成  64k  若是操作码与取值操作都需要跨页则需要额外的时钟周期，这也是上面为什么使用 & 的原因

func (o *Olc6502) IMP() bool {
	o.Fetched = o.A
	return false
}

func (o *Olc6502) IMM() bool {
	o.AddrAbs = o.Pc
	o.Pc++
	return false
}

func (o *Olc6502) ZP0() bool {
	o.AddrAbs = uint16(o.Read(o.Pc)) & 0x00FF
	o.Pc++
	return false
}

func (o *Olc6502) ZPX() bool {
	o.AddrAbs = uint16(o.Read(o.Pc)) + uint16(o.X)
	o.Pc++
	o.AddrAbs &= 0x00FF
	return false
}

func (o *Olc6502) ZPY() bool {
	o.AddrAbs = uint16(o.Read(o.Pc)) + uint16(o.Y)
	o.Pc++
	o.AddrAbs &= 0x00FF
	return false
}

func (o *Olc6502) REL() bool {
	o.AddrRel = uint16(o.Read(o.Pc))
	o.Pc++
	if (o.AddrRel & 0x80) > 0 { // 128
		o.AddrRel |= 0xFF00
	}
	return false
}

func (o *Olc6502) ABS() bool {
	lo := uint16(o.Read(o.Pc))
	o.Pc++
	hi := uint16(o.Read(o.Pc))
	o.Pc++
	o.AddrAbs = (hi << 8) | lo
	return false
}

func (o *Olc6502) ABX() bool {
	lo := uint16(o.Read(o.Pc))
	o.Pc++
	hi := uint16(o.Read(o.Pc))
	o.Pc++
	o.AddrAbs = (hi << 8) | lo
	o.AddrAbs += uint16(o.X)
	// 因为需要 额外加 x，可能高 8位发生变化 产生换页 耗时加1
	if (o.AddrAbs & 0xFF00) != (hi << 8) {
		return true
	}
	return false
}

func (o *Olc6502) ABY() bool {
	lo := uint16(o.Read(o.Pc))
	o.Pc++
	hi := uint16(o.Read(o.Pc))
	o.Pc++
	o.AddrAbs = (hi << 8) | lo
	o.AddrAbs += uint16(o.Y)
	// 因为需要 额外加 x，可能高 8位发生变化 产生换页 耗时加1
	if (o.AddrAbs & 0xFF00) != (hi << 8) {
		return true
	}
	return false
}

func (o *Olc6502) IND() bool {
	// 先获取地址
	ptrLo := uint16(o.Read(o.Pc))
	o.Pc++
	ptrHi := uint16(o.Read(o.Pc))
	o.Pc++
	ptr := (ptrHi << 8) | ptrLo
	// 再获取地址对应的值 作为地址
	if ptrLo == 0x00FF { // 模拟硬件上的一个bug相当于没有翻页，还是原页，页内地址变成了0
		o.AddrAbs = (uint16(o.Read(ptr&0xFF00)) << 8) | uint16(o.Read(ptr+0))
	} else {
		o.AddrAbs = (uint16(o.Read(ptr+1)) << 8) | uint16(o.Read(ptr+0))
	}
	return false
}

func (o *Olc6502) IZX() bool {
	// 在第0页进行寻址,只需要先找页内地址即可
	ptr := uint16(o.Read(o.Pc))
	o.Pc++
	// 需要根据地址找对应的值作为地址，一切都要在第0页进行
	lo := uint16(o.Read((ptr + uint16(o.X)) & 0x00FF))
	hi := uint16(o.Read((ptr + uint16(o.X) + 1) & 0x00FF))
	o.AddrAbs = (hi << 8) | lo
	return false
}

func (o *Olc6502) IZY() bool {
	// 与IZX类似宅0页操作，但是偏移Y是最后进行的，也需要判断是否发生页面跳转
	ptr := uint16(o.Read(o.Pc))
	o.Pc++
	lo := uint16(o.Read(ptr & 0x00FF))
	hi := uint16(o.Read((ptr + 1) & 0x00FF))
	o.AddrAbs = (hi << 8) | lo
	o.AddrAbs += uint16(o.Y)
	if (o.AddrAbs & 0xFF00) != (hi << 8) {
		return true
	}
	return false
}

//============Opcodes=============
//***********官方操作码***********

func (o *Olc6502) ADC() bool {
	// 累加值到 Temp 并考虑进位 最终操作始终在 0~0x00FF
	o.Fetch()
	temp := uint16(o.A) + uint16(o.Fetched)
	if o.GetFlag(C) {
		temp++
	}
	// 设置计算后的各种标记
	o.SetFlag(C, temp > 0x00FF)
	o.SetFlag(Z, (temp&0x00FF) == 0)
	o.SetFlag(N, temp&0x80 > 0)
	// 计算溢出标记，先计算相关数据的负数标记
	tA := (o.A & 0x80) > 0
	tF := (o.Fetched & 0x80) > 0
	tT := (temp & 0x80) > 0
	if tA && tF && !tT { // 负数+负数=正数 溢出
		o.SetFlag(V, true)
	} else if !tA && !tF && tT { // 正数+正数=负数 溢出
		o.SetFlag(V, true)
	} else { // 不会溢出
		o.SetFlag(V, false)
	}
	o.A = uint8(temp & 0xFF) // 最终还是应用到累加器上
	return true
}

func (o *Olc6502) AND() bool {
	o.Fetch()
	o.A &= o.Fetched
	o.SetFlag(Z, o.A == 0x00)
	o.SetFlag(N, o.A&0x80 > 0)
	return true
}

func (o *Olc6502) ASL() bool {
	o.Fetch()
	temp := uint16(o.Fetched) << 1
	o.SetFlag(C, (temp&0xFF00) > 0)
	o.SetFlag(Z, (temp&0x00FF) == 0x0000)
	o.SetFlag(N, (temp&0x80) > 0)
	if AddrEq(o.Lookup[o.OpCode].AddrMode, o.IMP) {
		o.A = uint8(temp & 0xFF)
	} else {
		o.Write(o.AddrAbs, uint8(temp&0xFF))
	}
	return false
}

func (o *Olc6502) BCC() bool {
	if !o.GetFlag(C) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BCS() bool {
	if o.GetFlag(C) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BEQ() bool {
	if o.GetFlag(Z) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BIT() bool {
	o.Fetch()
	temp := uint16(o.A & o.Fetched)
	o.SetFlag(Z, (temp&0x00FF) == 0)
	// 为什么突然使用 Fetched 了?
	o.SetFlag(N, (o.Fetched&0x80) > 0)
	o.SetFlag(V, (o.Fetched&0x40) > 0)
	return false
}

func (o *Olc6502) BMI() bool {
	if o.GetFlag(N) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BNE() bool {
	if !o.GetFlag(Z) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BPL() bool {
	if !o.GetFlag(N) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BRK() bool {
	// 保存指针与状态
	o.Pc++
	o.SetFlag(I, true)
	o.Write(0x0100+uint16(o.StackP), uint8((o.Pc>>8)&0xFF))
	o.StackP--
	o.Write(0x0100+uint16(o.StackP), uint8(o.Pc&0xFF))
	o.StackP--
	o.SetFlag(B, true)
	o.Write(0x0100+uint16(o.StackP), o.Status)
	o.StackP--
	o.SetFlag(B, false)
	// 跳转到中断处理程序
	o.Pc = (uint16(o.Read(0xFFFF)) << 8) | uint16(o.Read(0xFFFE))
	return false
}

func (o *Olc6502) BVC() bool {
	if !o.GetFlag(V) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) BVS() bool {
	if o.GetFlag(V) { // 分支判断 跳转到指定位置
		o.Cycles++
		o.AddrAbs = o.Pc + o.AddrRel // 并判断是否发生切页
		if (o.AddrAbs & 0xFF00) != (o.Pc & 0xFF00) {
			o.Cycles++
		}
		o.Pc = o.AddrAbs
	}
	return false
}

func (o *Olc6502) CLC() bool {
	o.SetFlag(C, false)
	return false
}

func (o *Olc6502) CLD() bool {
	o.SetFlag(D, false)
	return false
}

func (o *Olc6502) CLI() bool {
	o.SetFlag(I, false)
	return false
}

func (o *Olc6502) CLV() bool {
	o.SetFlag(V, false)
	return false
}

func (o *Olc6502) CMP() bool {
	o.Fetch()
	temp := uint16(o.A) - uint16(o.Fetched)
	o.SetFlag(C, o.A >= o.Fetched)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	return true
}

func (o *Olc6502) CPX() bool {
	o.Fetch()
	temp := uint16(o.X) - uint16(o.Fetched)
	o.SetFlag(C, o.X >= o.Fetched)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	return false
}

func (o *Olc6502) CPY() bool {
	o.Fetch()
	temp := uint16(o.Y) - uint16(o.Fetched)
	o.SetFlag(C, o.Y >= o.Fetched)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	return false
}

func (o *Olc6502) DEC() bool {
	o.Fetch()
	temp := uint16(o.Fetched - 1)
	o.Write(o.AddrAbs, uint8(temp&0xFF))
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	return false
}

func (o *Olc6502) DEX() bool {
	o.X--
	o.SetFlag(Z, o.X == 0)
	o.SetFlag(N, (o.X&0x80) > 0)
	return false
}

func (o *Olc6502) DEY() bool {
	o.Y--
	o.SetFlag(Z, o.Y == 0)
	o.SetFlag(N, (o.Y&0x80) > 0)
	return false
}

func (o *Olc6502) EOR() bool {
	o.Fetch()
	o.A ^= o.Fetched
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return true
}

func (o *Olc6502) INC() bool {
	o.Fetch()
	temp := uint16(o.Fetched + 1)
	o.Write(o.AddrAbs, uint8(temp&0xFF))
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	return false
}

func (o *Olc6502) INX() bool {
	o.X++
	o.SetFlag(Z, o.X == 0)
	o.SetFlag(N, (o.X&0x80) > 0)
	return false
}

func (o *Olc6502) INY() bool {
	o.Y++
	o.SetFlag(Z, o.Y == 0)
	o.SetFlag(N, (o.Y&0x80) > 0)
	return false
}

func (o *Olc6502) JMP() bool {
	o.Pc = o.AddrAbs
	return false
}

func (o *Olc6502) JSR() bool {
	// 先保存pc指针再跳转
	o.Pc--
	o.Write(0x0100+uint16(o.StackP), uint8((o.Pc>>8)&0xFF))
	o.StackP--
	o.Write(0x0100+uint16(o.StackP), uint8(o.Pc&0xFF))
	o.StackP--
	o.Pc = o.AddrAbs
	return false
}

func (o *Olc6502) LDA() bool {
	o.Fetch()
	o.A = o.Fetched
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return true
}

func (o *Olc6502) LDX() bool {
	o.Fetch()
	o.X = o.Fetched
	o.SetFlag(Z, o.X == 0)
	o.SetFlag(N, (o.X&0x80) > 0)
	return true
}

func (o *Olc6502) LDY() bool {
	o.Fetch()
	o.Y = o.Fetched
	o.SetFlag(Z, o.Y == 0)
	o.SetFlag(N, (o.Y&0x80) > 0)
	return false
}

func (o *Olc6502) LSR() bool { // 移位运算符
	o.Fetch()
	o.SetFlag(C, (o.Fetched&0x1) > 0)
	temp := uint16(o.Fetched >> 1)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	if AddrEq(o.Lookup[o.OpCode].AddrMode, o.IMP) {
		o.A = uint8(temp & 0xFF)
	} else {
		o.Write(o.AddrAbs, uint8(temp&0xFF))
	}
	return false
}

func (o *Olc6502) NOP() bool {
	switch o.OpCode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return true
	}
	return false
}

func (o *Olc6502) ORA() bool {
	o.Fetch()
	o.A |= o.Fetched
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return true
}

func (o *Olc6502) PHA() bool {
	o.Write(0x0100+uint16(o.StackP), o.A)
	o.StackP--
	return false
}

func (o *Olc6502) PHP() bool {
	o.Write(0x0100+uint16(o.StackP), o.Status|B|U)
	o.StackP--
	o.SetFlag(B, false)
	o.SetFlag(U, false)
	return false
}

func (o *Olc6502) PLA() bool { // 弹出栈
	o.StackP++
	o.A = o.Read(0x0100 + uint16(o.StackP))
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return false
}

func (o *Olc6502) PLP() bool {
	o.StackP++
	o.Status = o.Read(0x0100 + uint16(o.StackP))
	o.SetFlag(U, true)
	return false
}

func (o *Olc6502) ROL() bool {
	o.Fetch()
	temp := uint16(o.Fetched << 1)
	if o.GetFlag(C) {
		temp |= 0x0001
	}
	o.SetFlag(C, temp > 0xFF)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	if AddrEq(o.Lookup[o.OpCode].AddrMode, o.IMP) {
		o.A = uint8(temp & 0xFF)
	} else {
		o.Write(o.AddrAbs, uint8(temp&0xFF))
	}
	return false
}

func (o *Olc6502) ROR() bool { // 位移操作
	o.Fetch()
	temp := uint16(o.Fetched) >> 1
	if o.GetFlag(C) {
		temp |= 0x80
	}
	o.SetFlag(C, (temp&0x01) > 0)
	o.SetFlag(Z, (temp&0xFF) == 0)
	o.SetFlag(N, (temp&0x80) > 0)
	if AddrEq(o.Lookup[o.OpCode].AddrMode, o.IMP) {
		o.A = uint8(temp & 0xFF)
	} else {
		o.Write(o.AddrAbs, uint8(temp&0xFF))
	}
	return false
}

func (o *Olc6502) RTI() bool {
	// 恢复到原来的状态与地址
	o.StackP++
	o.Status = o.Read(0x0100 + uint16(o.StackP))
	o.SetFlag(B, false)
	o.SetFlag(U, false)
	o.StackP++
	lo := uint16(o.Read(0x0100 + uint16(o.StackP)))
	o.StackP++
	hi := uint16(o.Read(0x0100 + uint16(o.StackP)))
	o.Pc = lo | (hi << 8)
	return false
}

func (o *Olc6502) RTS() bool { // 只恢复指针
	o.StackP++
	lo := uint16(o.Read(0x0100 + uint16(o.StackP)))
	o.StackP++
	hi := uint16(o.Read(0x0100 + uint16(o.StackP)))
	o.Pc = lo | (hi << 8)
	o.Pc++
	return false
}

func (o *Olc6502) SBC() bool {
	// 累减值到 Temp 并考虑进位 最终操作始终在 0~0x00FF
	o.Fetch()
	value := uint16(o.Fetched) ^ 0x00FF
	temp := uint16(o.A) + value
	if o.GetFlag(C) {
		temp++
	}
	// 设置计算后的各种标记
	o.SetFlag(C, temp > 0x00FF)
	o.SetFlag(Z, (temp&0x00FF) == 0)
	// 计算溢出标记，先计算相关数据的负数标记
	tA := (o.A & 0x80) > 0
	tF := (o.Fetched & 0x80) > 0
	tT := (temp & 0x80) > 0
	if tA && !tF && !tT { // 负数-正数=正数 溢出
		o.SetFlag(V, true)
	} else if !tA && tF && tT { // 正数-负数=负数 溢出
		o.SetFlag(V, true)
	} else { // 不会溢出
		o.SetFlag(V, false)
	}
	o.SetFlag(N, (temp&0x80) > 0)
	o.A = uint8(temp & 0xFF) // 最终还是应用到累加器上
	return true
}

func (o *Olc6502) SEC() bool {
	o.SetFlag(C, true)
	return false
}

func (o *Olc6502) SED() bool {
	o.SetFlag(D, true)
	return false
}

func (o *Olc6502) SEI() bool {
	o.SetFlag(I, true)
	return false
}

func (o *Olc6502) STA() bool {
	o.Write(o.AddrAbs, o.A)
	return false
}

func (o *Olc6502) STX() bool {
	o.Write(o.AddrAbs, o.X)
	return false
}

func (o *Olc6502) STY() bool {
	o.Write(o.AddrAbs, o.Y)
	return false
}

func (o *Olc6502) TAX() bool {
	o.X = o.A
	o.SetFlag(Z, o.X == 0)
	o.SetFlag(N, (o.X&0x80) > 0)
	return false
}

func (o *Olc6502) TAY() bool {
	o.Y = o.A
	o.SetFlag(Z, o.Y == 0)
	o.SetFlag(N, (o.Y&0x80) > 0)
	return false
}

func (o *Olc6502) TSX() bool {
	o.X = o.StackP
	o.SetFlag(Z, o.X == 0)
	o.SetFlag(N, (o.X&0x80) > 0)
	return false
}

func (o *Olc6502) TXA() bool {
	o.A = o.X
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return false
}

func (o *Olc6502) TXS() bool {
	o.StackP = o.X
	return false
}

func (o *Olc6502) TYA() bool {
	o.A = o.Y
	o.SetFlag(Z, o.A == 0)
	o.SetFlag(N, (o.A&0x80) > 0)
	return false
}

//*********非官方操作码统一使用这个***********

func (o *Olc6502) XXX() bool {
	return false
}
