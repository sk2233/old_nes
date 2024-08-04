/*
@author: sk
@date: 2023/11/11
*/
package main

import (
	"strconv"
	"strings"
)

//===============Cpu==================

type Cpu struct {
	Bus                 *Bus
	State               *CpuState
	A, X, Y             uint8          // 寄存器
	StackP              uint8          // 栈指针
	Pc                  uint16         // Pc寄存器
	Instructions        []*Instruction // 256 个操作指令
	Cycles              uint8          // 用于模拟指令执行消耗的时钟数
	InstructionComplete bool
}

func NewCpu() *Cpu {
	return &Cpu{State: NewCpuState(), Instructions: []*Instruction{
		{"BRK", &BRK{}, &IMM{}, 7}, {"ORA", &ORA{}, &IZX{}, 6}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 3}, {"ORA", &ORA{}, &ZP0{}, 3}, {"ASL", &ASL{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"PHP", &PHP{}, &IMP{}, 3}, {"ORA", &ORA{}, &IMM{}, 2}, {"ASL", &ASL{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"???", &NOP{}, &IMP{}, 4}, {"ORA", &ORA{}, &ABS{}, 4}, {"ASL", &ASL{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BPL", &BPL{}, &REL{}, 2}, {"ORA", &ORA{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"ORA", &ORA{}, &ZPX{}, 4}, {"ASL", &ASL{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"CLC", &CLC{}, &IMP{}, 2}, {"ORA", &ORA{}, &ABY{}, 4}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"ORA", &ORA{}, &ABX{}, 4}, {"ASL", &ASL{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
		{"JSR", &JSR{}, &ABS{}, 6}, {"AND", &AND{}, &IZX{}, 6}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"BIT", &BIT{}, &ZP0{}, 3}, {"AND", &AND{}, &ZP0{}, 3}, {"ROL", &ROL{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"PLP", &PLP{}, &IMP{}, 4}, {"AND", &AND{}, &IMM{}, 2}, {"ROL", &ROL{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"BIT", &BIT{}, &ABS{}, 4}, {"AND", &AND{}, &ABS{}, 4}, {"ROL", &ROL{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BMI", &BMI{}, &REL{}, 2}, {"AND", &AND{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"AND", &AND{}, &ZPX{}, 4}, {"ROL", &ROL{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"SEC", &SEC{}, &IMP{}, 2}, {"AND", &AND{}, &ABY{}, 4}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"AND", &AND{}, &ABX{}, 4}, {"ROL", &ROL{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
		{"RTI", &RTI{}, &IMP{}, 6}, {"EOR", &EOR{}, &IZX{}, 6}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 3}, {"EOR", &EOR{}, &ZP0{}, 3}, {"LSR", &LSR{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"PHA", &PHA{}, &IMP{}, 3}, {"EOR", &EOR{}, &IMM{}, 2}, {"LSR", &LSR{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"JMP", &JMP{}, &ABS{}, 3}, {"EOR", &EOR{}, &ABS{}, 4}, {"LSR", &LSR{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BVC", &BVC{}, &REL{}, 2}, {"EOR", &EOR{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"EOR", &EOR{}, &ZPX{}, 4}, {"LSR", &LSR{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"CLI", &CLI{}, &IMP{}, 2}, {"EOR", &EOR{}, &ABY{}, 4}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"EOR", &EOR{}, &ABX{}, 4}, {"LSR", &LSR{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
		{"RTS", &RTS{}, &IMP{}, 6}, {"ADC", &ADC{}, &IZX{}, 6}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 3}, {"ADC", &ADC{}, &ZP0{}, 3}, {"ROR", &ROR{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"PLA", &PLA{}, &IMP{}, 4}, {"ADC", &ADC{}, &IMM{}, 2}, {"ROR", &ROR{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"JMP", &JMP{}, &IND{}, 5}, {"ADC", &ADC{}, &ABS{}, 4}, {"ROR", &ROR{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BVS", &BVS{}, &REL{}, 2}, {"ADC", &ADC{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"ADC", &ADC{}, &ZPX{}, 4}, {"ROR", &ROR{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"SEI", &SEI{}, &IMP{}, 2}, {"ADC", &ADC{}, &ABY{}, 4}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"ADC", &ADC{}, &ABX{}, 4}, {"ROR", &ROR{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
		{"???", &NOP{}, &IMP{}, 2}, {"STA", &STA{}, &IZX{}, 6}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 6}, {"STY", &STY{}, &ZP0{}, 3}, {"STA", &STA{}, &ZP0{}, 3}, {"STX", &STX{}, &ZP0{}, 3}, {"???", &XXX{}, &IMP{}, 3}, {"DEY", &DEY{}, &IMP{}, 2}, {"???", &NOP{}, &IMP{}, 2}, {"TXA", &TXA{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"STY", &STY{}, &ABS{}, 4}, {"STA", &STA{}, &ABS{}, 4}, {"STX", &STX{}, &ABS{}, 4}, {"???", &XXX{}, &IMP{}, 4},
		{"BCC", &BCC{}, &REL{}, 2}, {"STA", &STA{}, &IZY{}, 6}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 6}, {"STY", &STY{}, &ZPX{}, 4}, {"STA", &STA{}, &ZPX{}, 4}, {"STX", &STX{}, &ZPY{}, 4}, {"???", &XXX{}, &IMP{}, 4}, {"TYA", &TYA{}, &IMP{}, 2}, {"STA", &STA{}, &ABY{}, 5}, {"TXS", &TXS{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 5}, {"???", &NOP{}, &IMP{}, 5}, {"STA", &STA{}, &ABX{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"???", &XXX{}, &IMP{}, 5},
		{"LDY", &LDY{}, &IMM{}, 2}, {"LDA", &LDA{}, &IZX{}, 6}, {"LDX", &LDX{}, &IMM{}, 2}, {"???", &XXX{}, &IMP{}, 6}, {"LDY", &LDY{}, &ZP0{}, 3}, {"LDA", &LDA{}, &ZP0{}, 3}, {"LDX", &LDX{}, &ZP0{}, 3}, {"???", &XXX{}, &IMP{}, 3}, {"TAY", &TAY{}, &IMP{}, 2}, {"LDA", &LDA{}, &IMM{}, 2}, {"TAX", &TAX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"LDY", &LDY{}, &ABS{}, 4}, {"LDA", &LDA{}, &ABS{}, 4}, {"LDX", &LDX{}, &ABS{}, 4}, {"???", &XXX{}, &IMP{}, 4},
		{"BCS", &BCS{}, &REL{}, 2}, {"LDA", &LDA{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 5}, {"LDY", &LDY{}, &ZPX{}, 4}, {"LDA", &LDA{}, &ZPX{}, 4}, {"LDX", &LDX{}, &ZPY{}, 4}, {"???", &XXX{}, &IMP{}, 4}, {"CLV", &CLV{}, &IMP{}, 2}, {"LDA", &LDA{}, &ABY{}, 4}, {"TSX", &TSX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 4}, {"LDY", &LDY{}, &ABX{}, 4}, {"LDA", &LDA{}, &ABX{}, 4}, {"LDX", &LDX{}, &ABY{}, 4}, {"???", &XXX{}, &IMP{}, 4},
		{"CPY", &CPY{}, &IMM{}, 2}, {"CMP", &CMP{}, &IZX{}, 6}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"CPY", &CPY{}, &ZP0{}, 3}, {"CMP", &CMP{}, &ZP0{}, 3}, {"DEC", &DEC{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"INY", &INY{}, &IMP{}, 2}, {"CMP", &CMP{}, &IMM{}, 2}, {"DEX", &DEX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 2}, {"CPY", &CPY{}, &ABS{}, 4}, {"CMP", &CMP{}, &ABS{}, 4}, {"DEC", &DEC{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BNE", &BNE{}, &REL{}, 2}, {"CMP", &CMP{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"CMP", &CMP{}, &ZPX{}, 4}, {"DEC", &DEC{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"CLD", &CLD{}, &IMP{}, 2}, {"CMP", &CMP{}, &ABY{}, 4}, {"NOP", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"CMP", &CMP{}, &ABX{}, 4}, {"DEC", &DEC{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
		{"CPX", &CPX{}, &IMM{}, 2}, {"SBC", &SBC{}, &IZX{}, 6}, {"???", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"CPX", &CPX{}, &ZP0{}, 3}, {"SBC", &SBC{}, &ZP0{}, 3}, {"INC", &INC{}, &ZP0{}, 5}, {"???", &XXX{}, &IMP{}, 5}, {"INX", &INX{}, &IMP{}, 2}, {"SBC", &SBC{}, &IMM{}, 2}, {"NOP", &NOP{}, &IMP{}, 2}, {"???", &SBC{}, &IMP{}, 2}, {"CPX", &CPX{}, &ABS{}, 4}, {"SBC", &SBC{}, &ABS{}, 4}, {"INC", &INC{}, &ABS{}, 6}, {"???", &XXX{}, &IMP{}, 6},
		{"BEQ", &BEQ{}, &REL{}, 2}, {"SBC", &SBC{}, &IZY{}, 5}, {"???", &XXX{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 8}, {"???", &NOP{}, &IMP{}, 4}, {"SBC", &SBC{}, &ZPX{}, 4}, {"INC", &INC{}, &ZPX{}, 6}, {"???", &XXX{}, &IMP{}, 6}, {"SED", &SED{}, &IMP{}, 2}, {"SBC", &SBC{}, &ABY{}, 4}, {"NOP", &NOP{}, &IMP{}, 2}, {"???", &XXX{}, &IMP{}, 7}, {"???", &NOP{}, &IMP{}, 4}, {"SBC", &SBC{}, &ABX{}, 4}, {"INC", &INC{}, &ABX{}, 7}, {"???", &XXX{}, &IMP{}, 7},
	}}
}

func (c *Cpu) ConnectBus(bus *Bus) {
	c.Bus = bus
}

func (c *Cpu) Read(addr uint16) byte {
	return c.Bus.CpuRead(addr)
}

func (c *Cpu) ReadUint16(addr uint16) uint16 {
	return c.Bus.ReadUint16(addr, addr+1)
}

func (c *Cpu) ReadUint16Bug(addr uint16) uint16 { // 模拟 低位+1不会进位的bug
	lo := (addr & 0x00FF) + 1                 // 计算低位
	hiAddr := (addr & 0xFF00) | (lo & 0x00FF) // 低位只要低位，抛弃进位
	return c.Bus.ReadUint16(addr, hiAddr)
}

func (c *Cpu) Write(addr uint16, data byte) {
	c.Bus.CpuWrite(addr, data)
}

func (c *Cpu) Clock() {
	if c.Cycles > 0 {
		c.Cycles--
		return
	}
	c.InstructionComplete = true // 完成开始准备下一个指令
	opcode := c.Read(c.Pc)
	c.Pc++
	instruction := c.Instructions[opcode]
	c.Cycles = instruction.Cycles
	addr, flipPage1 := instruction.Address.GetAddress(c)
	flipPage2 := instruction.Operation.Operation(c, addr, opcode)
	if flipPage1 && flipPage2 { // 必须读取地址与 操作都进行了翻页才累加耗时
		c.Cycles++
	}
}

func (c *Cpu) Reset() {
	c.Pc = c.ReadUint16(0xFFFC) // 设置到程序入口
	c.StackP = 0xFD             // 恢复栈指针
	c.State.Set(0x20)
	c.Cycles = 8 // 固定停止8个时钟
}

func (c *Cpu) Irq() { // 中断请求
	if c.State.UnInterrupt {
		return
	}
	c.PushUint16(c.Pc)
	c.Push(c.State.Get())
	c.Pc = c.ReadUint16(0xFFFE)
	c.State.UnInterrupt = true
	c.Cycles = 7
}

func (c *Cpu) Nmi() { // 不可屏蔽的中断请求
	c.PushUint16(c.Pc)
	c.Push(c.State.Get())
	c.Pc = c.ReadUint16(0xFFFA) // 跳转的程序入口不一致
	c.State.UnInterrupt = true
	c.Cycles = 8
}

func (c *Cpu) Push(data uint8) {
	c.Write(0x0100+uint16(c.StackP), data) // 0x0100 就是栈的基地址
	c.StackP--                             // 下压栈
}

func (c *Cpu) PushUint16(data uint16) {
	c.Push(uint8(data >> 8)) // 先高后低
	c.Push(uint8(data))
}

func (c *Cpu) Pop() uint8 {
	c.StackP++
	return c.Read(0x0100 + uint16(c.StackP))
}

func (c *Cpu) PopUint16() uint16 {
	lo := c.Pop()
	hi := c.Pop()
	return uint16(lo) | (uint16(hi) << 8)
}

func (c *Cpu) Disassemble(start, end uint32) (map[uint16]string, []uint16) {
	res := make(map[uint16]string)
	addrs := make([]uint16, 0)
	addr := start
	for addr < end {
		tempAddr := uint16(addr)
		opcode := c.Read(uint16(addr))
		addr++
		instruction := c.Instructions[opcode]
		name := instruction.Address.GetName()
		line := Format(tempAddr, 4) + ": {" + instruction.Name + "} : [" + name + "] :"
		switch name {
		case "IMP":
		case "IMM", "ZP0", "ZPX", "ZPY", "IZX", "IZY", "REL":
			line += " " + Format(c.Read(uint16(addr)), 2)
			addr++
		case "ABS", "ABX", "ABY", "IND":
			line += " " + Format(c.ReadUint16(uint16(addr)), 4)
			addr += 2
		}
		res[tempAddr] = line
		addrs = append(addrs, tempAddr)
	}
	return res, addrs
}

func Format[T uint16 | uint8](value T, count int) string {
	temp := strconv.FormatUint(uint64(value), 16)
	temp = strings.Repeat("0", count-len(temp)) + temp
	return "0x" + strings.ToUpper(temp)
}

//==================Instruction===================

type Instruction struct {
	Name      string
	Operation IOperation
	Address   IAddress
	Cycles    uint8
}

//=================Address==================

type IAddress interface {
	GetAddress(cpu *Cpu) (uint16, bool)
	GetName() string // 获取地址，并返回是否内存翻页
}

func FlipPage(old, new uint16) bool { // 判断新旧地址是否发生翻页
	return old&0xFF00 != new&0xFF00
}

type IMP struct { // 以自身(A寄存器)为操作数，无需取值
}

func (I *IMP) GetName() string {
	return "IMP"
}

func (I *IMP) GetAddress(_ *Cpu) (uint16, bool) {
	return 0, false // 地址随意返回。不会被使用
}

type IMM struct {
}

func (I *IMM) GetName() string {
	return "IMM"
}

func (I *IMM) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.Pc
	cpu.Pc++
	return addr, false
}

type ZP0 struct {
}

func (Z *ZP0) GetName() string {
	return "ZP0"
}

func (Z *ZP0) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.Read(cpu.Pc)
	cpu.Pc++
	return uint16(addr), false
}

type ZPX struct {
}

func (Z *ZPX) GetName() string {
	return "ZPX"
}

func (Z *ZPX) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.Read(cpu.Pc) + cpu.X // 依旧是0页寻址 若是溢出 256 只取页内偏移 可以直接溢出
	cpu.Pc++
	return uint16(addr), false
}

type ZPY struct {
}

func (Z *ZPY) GetName() string {
	return "ZPY"
}

func (Z *ZPY) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.Read(cpu.Pc) + cpu.Y // 依旧是0页寻址 若是溢出 256 只取页内偏移 可以直接溢出
	cpu.Pc++
	return uint16(addr), false
}

type REL struct {
}

func (R *REL) GetName() string {
	return "REL"
}

func (R *REL) GetAddress(cpu *Cpu) (uint16, bool) { // 只适用于分支跳转，返回的也是绝对地址
	offset := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	if offset < 0x80 { // 128 位置存在就是负数 相对跳转位置位于 -127 ~ 127
		return cpu.Pc + offset, false
	} else {
		return cpu.Pc + offset - 0x100, false
	}
}

type ABS struct {
}

func (A *ABS) GetName() string {
	return "ABS"
}

func (A *ABS) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.ReadUint16(cpu.Pc)
	cpu.Pc += 2
	return addr, false
}

type ABX struct {
}

func (A *ABX) GetName() string {
	return "ABX"
}

func (A *ABX) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.ReadUint16(cpu.Pc)
	cpu.Pc += 2
	oldAddr := addr
	addr += uint16(cpu.X)
	return addr, FlipPage(oldAddr, addr) // 可能发生翻页
}

type ABY struct {
}

func (A *ABY) GetName() string {
	return "ABY"
}

func (A *ABY) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.ReadUint16(cpu.Pc)
	cpu.Pc += 2
	oldAddr := addr
	addr += uint16(cpu.Y)
	return addr, FlipPage(oldAddr, addr) // 可能发生翻页
}

type IND struct {
}

func (I *IND) GetName() string {
	return "IND"
}

func (I *IND) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.ReadUint16(cpu.Pc)
	cpu.Pc += 2
	return cpu.ReadUint16Bug(addr), false
}

type IZX struct {
}

func (I *IZX) GetName() string {
	return "IZX"
}

func (I *IZX) GetAddress(cpu *Cpu) (uint16, bool) {
	addr := cpu.Read(cpu.Pc) + cpu.X // 从0地址间接寻址 0xFF 溢出的会舍弃
	cpu.Pc++
	return cpu.ReadUint16(uint16(addr)), false
}

type IZY struct {
}

func (I *IZY) GetName() string {
	return "IZY"
}

func (I *IZY) GetAddress(cpu *Cpu) (uint16, bool) { // 与IZX类似宅0页操作，但是偏移Y是最后进行的，也需要判断是否发生页面跳转
	addr := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	addr = cpu.ReadUint16(addr)
	oldAddr := addr
	addr += uint16(cpu.Y)
	return addr, FlipPage(oldAddr, addr)
}

//====================Operation====================

type IOperation interface {
	Operation(cpu *Cpu, addr uint16, opcode byte) bool // 也需要返回是否内存翻页
}

func JumpTo(cpu *Cpu, addr uint16) {
	cpu.Pc = addr // 跳转到指定地址
	cpu.Cycles++  // 跳转本身就需要 一个周期 如果翻页就再需要一个周期
	if FlipPage(addr, cpu.Pc) {
		cpu.Cycles++
	}
}

func Compare(cpu *Cpu, a, b uint8) {
	cpu.State.UpdateZN(a - b)
	cpu.State.Carry = a >= b
}

type ADC struct {
}

func (A *ADC) Operation(cpu *Cpu, addr uint16, _ byte) bool { // 加法操作
	a := cpu.A
	b := cpu.Read(addr)
	c := cpu.State.GetCarryNum() // 进位数
	cpu.A = a + b + c
	cpu.State.UpdateZN(cpu.A)
	cpu.State.Carry = int(a)+int(b)+int(c) > 0xFF
	// 两位加数同号,但是最终结果却不与他们同号 结果溢出
	cpu.State.Overflow = (a^b)&0x80 == 0 && (a^cpu.A)&0x80 != 0
	return true
}

type AND struct {
}

func (A *AND) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	data := cpu.Read(addr)
	cpu.A &= data
	cpu.State.UpdateZN(cpu.A)
	return true
}

type ASL struct {
}

func (A *ASL) Operation(cpu *Cpu, addr uint16, opcode byte) bool {
	if cpu.Instructions[opcode].Address.GetName() == "IMP" { // IMP 是以自身为操作对象的操作 类似 i++ i<<=1
		val := uint16(cpu.A) << 1
		cpu.A = uint8(val)
		cpu.State.Carry = val > 0xFF
		cpu.State.UpdateZN(cpu.A)
	} else { // 非自身为操作对象
		val := uint16(cpu.Read(addr)) << 1
		cpu.Write(addr, byte(val))
		cpu.State.Carry = val > 0xFF
		cpu.State.UpdateZN(uint8(val))
	}
	return false
}

type BCC struct {
}

func (B *BCC) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if !cpu.State.Carry {
		JumpTo(cpu, addr)
	}
	return false
}

type BCS struct {
}

func (B *BCS) Operation(cpu *Cpu, addr uint16, _ byte) bool { // 分支判断
	if cpu.State.Carry {
		JumpTo(cpu, addr)
	}
	return false
}

type BEQ struct {
}

func (B *BEQ) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if cpu.State.Zero {
		JumpTo(cpu, addr)
	}
	return false
}

type BIT struct {
}

func (B *BIT) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	temp := cpu.Read(addr)
	cpu.State.Overflow = (temp>>6)&1 > 0
	cpu.State.UpdateZ(temp & cpu.A)
	cpu.State.UpdateN(temp)
	return false
}

type BMI struct {
}

func (B *BMI) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if cpu.State.Negative {
		JumpTo(cpu, addr)
	}
	return false
}

type BNE struct {
}

func (B *BNE) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if !cpu.State.Zero {
		JumpTo(cpu, addr)
	}
	return false
}

type BPL struct {
}

func (B *BPL) Operation(
	cpu *Cpu,
	addr uint16,
	_ byte,
) bool {
	if !cpu.State.Negative {
		JumpTo(cpu, addr)
	}
	return false
}

type BRK struct {
}

func (B *BRK) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.PushUint16(cpu.Pc)
	cpu.Push(cpu.State.Get() | 0x10) // 设置中断状态
	cpu.State.UnInterrupt = true
	cpu.Pc = cpu.ReadUint16(0xFFFE)
	return false
}

type BVC struct {
}

func (B *BVC) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if !cpu.State.Overflow {
		JumpTo(cpu, addr)
	}
	return false
}

type BVS struct {
}

func (B *BVS) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	if cpu.State.Overflow {
		JumpTo(cpu, addr)
	}
	return false
}

type CLC struct {
}

func (C *CLC) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Carry = false
	return false
}

type CLD struct {
}

func (C *CLD) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Decimal = false
	return false
}

type CLI struct {
}

func (C *CLI) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.UnInterrupt = false
	return false
}

type CLV struct {
}

func (C *CLV) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Overflow = false
	return false
}

type CMP struct {
}

func (C *CMP) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	val := cpu.Read(addr)
	Compare(cpu, cpu.A, val)
	return true
}

type CPX struct {
}

func (C *CPX) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	val := cpu.Read(addr)
	Compare(cpu, cpu.X, val)
	return false
}

type CPY struct {
}

func (C *CPY) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	val := cpu.Read(addr)
	Compare(cpu, cpu.Y, val)
	return false
}

type DEC struct {
}

func (D *DEC) Operation(cpu *Cpu, addr uint16, _ byte) bool { // i--
	val := cpu.Read(addr) - 1
	cpu.Write(addr, val)
	cpu.State.UpdateZN(val)
	return false
}

type DEX struct { // X--

}

func (D *DEX) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.X--
	cpu.State.UpdateZN(cpu.X)
	return false
}

type DEY struct { // y--

}

func (D *DEY) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Y--
	cpu.State.UpdateZN(cpu.Y)
	return false
}

type EOR struct { // 异或操作

}

func (R *EOR) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.A ^= cpu.Read(addr)
	cpu.State.UpdateZN(cpu.A)
	return false
}

type INC struct { // 内存自增
}

func (I *INC) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	val := cpu.Read(addr) + 1
	cpu.Write(addr, val)
	cpu.State.UpdateZN(val)
	return false
}

type INX struct { // 自增 X

}

func (I *INX) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.X++
	cpu.State.UpdateZN(cpu.X)
	return false
}

type INY struct {
}

func (I *INY) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Y++
	cpu.State.UpdateZN(cpu.Y)
	return false
}

type JMP struct {
}

func (J *JMP) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.Pc = addr
	return false
}

type JSR struct {
}

func (J *JSR) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.PushUint16(cpu.Pc - 1)
	cpu.Pc = addr
	return false
}

type LDA struct {
}

func (L *LDA) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.A = cpu.Read(addr)
	cpu.State.UpdateZN(cpu.A)
	return true
}

type LDX struct {
}

func (L *LDX) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.X = cpu.Read(addr)
	cpu.State.UpdateZN(cpu.X)
	return true
}

type LDY struct {
}

func (L *LDY) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.Y = cpu.Read(addr)
	cpu.State.UpdateZN(cpu.Y)
	return true
}

type LSR struct {
}

func (A *LSR) Operation(cpu *Cpu, addr uint16, opcode byte) bool {
	if cpu.Instructions[opcode].Address.GetName() == "IMP" { // IMP 是以自身为操作对象的操作 类似 i++ i<<=1
		cpu.State.Carry = cpu.A&0x01 > 0
		cpu.A >>= 1
		cpu.State.UpdateZN(cpu.A)
	} else { // 非自身为操作对象
		val := cpu.Read(addr)
		cpu.State.Carry = val&0x01 > 0
		val >>= 1
		cpu.Write(addr, val)
		cpu.State.UpdateZN(val)
	}
	return false
}

type NOP struct {
}

func (N *NOP) Operation(_ *Cpu, _ uint16, opcode byte) bool {
	switch opcode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return true
	default:
		return false
	}
}

type ORA struct {
}

func (O *ORA) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.A |= cpu.Read(addr)
	cpu.State.UpdateZN(cpu.A)
	return true
}

type PHA struct {
}

func (P *PHA) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Push(cpu.A)
	return false
}

type PHP struct {
}

func (P *PHP) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Push(cpu.State.Get() | 0x10)
	return false
}

type PLA struct {
}

func (P *PLA) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.A = cpu.Pop()
	cpu.State.UpdateZN(cpu.A)
	return false
}

type PLP struct {
}

func (P *PLP) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Set(cpu.Pop())
	return false
}

type ROL struct {
}

func (A *ROL) Operation(cpu *Cpu, addr uint16, opcode byte) bool {
	if cpu.Instructions[opcode].Address.GetName() == "IMP" { // IMP 是以自身为操作对象的操作 类似 i++ i<<=1
		c := cpu.State.GetCarryNum()
		cpu.State.Carry = (uint16(cpu.A) << 1) > 0xFF
		cpu.A = (cpu.A << 1) | c
		cpu.State.UpdateZN(cpu.A)
	} else { // 非自身为操作对象
		c := cpu.State.GetCarryNum()
		val := cpu.Read(addr)
		cpu.State.Carry = (uint16(val) << 1) > 0xFF
		val = (val << 1) | c
		cpu.Write(addr, val)
		cpu.State.UpdateZN(val)
	}
	return false
}

type ROR struct {
}

func (R *ROR) Operation(cpu *Cpu, addr uint16, opcode byte) bool {
	if cpu.Instructions[opcode].Address.GetName() == "IMP" { // IMP 是以自身为操作对象的操作 类似 i++ i<<=1
		c := cpu.State.GetCarryNum()
		cpu.State.Carry = (cpu.A & 1) > 0
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.State.UpdateZN(cpu.A)
	} else { // 非自身为操作对象
		c := cpu.State.GetCarryNum()
		val := cpu.Read(addr)
		cpu.State.Carry = (cpu.A & 1) > 0
		val = (val >> 1) | (c << 7)
		cpu.Write(addr, val)
		cpu.State.UpdateZN(val)
	}
	return false
}

type RTI struct {
}

func (R *RTI) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Set(cpu.Pop())
	cpu.Pc = cpu.PopUint16()
	return false
}

type RTS struct {
}

func (R *RTS) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Pc = cpu.PopUint16() + 1
	return false
}

type SBC struct {
}

func (S *SBC) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	a := cpu.A
	b := cpu.Read(addr)
	c := cpu.State.GetCarryNum()
	cpu.A = a - b - (1 - c) // 1-c 是下面进位条件导致的
	cpu.State.UpdateZN(cpu.A)
	cpu.State.Carry = int(a)-int(b)-int(1-c) >= 0
	// 减数与被减数符号不一致，且减数与结果符号也不一致
	cpu.State.Overflow = (a^b)&0x80 != 0 && (a^cpu.A)&0x80 != 0
	return true
}

type SEC struct {
}

func (S *SEC) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Carry = true
	return false
}

type SED struct {
}

func (S *SED) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.Decimal = true
	return false
}

type SEI struct {
}

func (S *SEI) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.State.UnInterrupt = true
	return false
}

type STA struct {
}

func (S *STA) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.Write(addr, cpu.A)
	return false
}

type STX struct {
}

func (S *STX) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.Write(addr, cpu.X)
	return false
}

type STY struct {
}

func (S *STY) Operation(cpu *Cpu, addr uint16, _ byte) bool {
	cpu.Write(addr, cpu.Y)
	return false
}

type TAX struct {
}

func (T *TAX) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.X = cpu.A
	cpu.State.UpdateZN(cpu.X)
	return false
}

type TAY struct {
}

func (T *TAY) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.Y = cpu.A
	cpu.State.UpdateZN(cpu.Y)
	return false
}

type TSX struct {
}

func (T *TSX) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.X = cpu.StackP
	cpu.State.UpdateZN(cpu.X)
	return false
}

type TXA struct {
}

func (T *TXA) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.A = cpu.X
	cpu.State.UpdateZN(cpu.A)
	return false
}

type TXS struct {
}

func (T *TXS) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.StackP = cpu.X
	return false
}

type TYA struct {
}

func (T *TYA) Operation(cpu *Cpu, _ uint16, _ byte) bool {
	cpu.A = cpu.Y
	cpu.State.UpdateZN(cpu.A)
	return false
}

type XXX struct { // 非官方操作码

}

func (X *XXX) Operation(*Cpu, uint16, byte) bool {
	return false
}

//==============CpuState================

type CpuState struct {
	Carry       bool // 是否有进位
	Zero        bool // 是否为0
	UnInterrupt bool // 是否屏蔽中断
	Decimal     bool // 十进制模式 暂时无用 没有实现
	Break       bool // Break 标记
	Unused      bool // 无用标记
	Overflow    bool // 溢出标记
	Negative    bool // 是否为负数标记
}

func (s *CpuState) UpdateZN(data uint8) {
	s.UpdateZ(data)
	s.UpdateN(data)
}

func (s *CpuState) UpdateZ(data uint8) {
	s.Zero = data == 0x00
}

func (s *CpuState) UpdateN(data uint8) {
	s.Negative = data >= 0x80
}

func (s *CpuState) GetCarryNum() uint8 { // 以数字形态获取进位方便使用
	if s.Carry {
		return 1
	}
	return 0
}

func (s *CpuState) Set(data uint8) {
	s.Carry = data&0x01 > 0
	s.Zero = data&0x02 > 0
	s.UnInterrupt = data&0x04 > 0
	s.Decimal = data&0x08 > 0
	s.Break = data&0x10 > 0
	s.Unused = data&0x20 > 0
	s.Overflow = data&0x40 > 0
	s.Negative = data&0x80 > 0
}

func (s *CpuState) Get() uint8 {
	data := uint8(0)
	if s.Carry {
		data |= 0x01
	}
	if s.Zero {
		data |= 0x02
	}
	if s.UnInterrupt {
		data |= 0x04
	}
	if s.Decimal {
		data |= 0x08
	}
	if s.Break {
		data |= 0x10
	}
	if s.Unused {
		data |= 0x20
	}
	if s.Overflow {
		data |= 0x40
	}
	if s.Negative {
		data |= 0x80
	}
	return data
}

func NewCpuState() *CpuState {
	return &CpuState{}
}
