/*
@author: sk
@date: 2023/10/26
*/
package nes4

import (
	"reflect"
)

const (
	C = 1 << iota // 进位标记
	Z             // 是否为0标记
	I             // 屏蔽普通中断标记
	D             // Unused
	B             // 中断标记
	U             // Unused
	V             // 计算溢出标记
	N             // 负数标记
)

type Instruction struct { // 一条指令的抽象
	Name     string
	Operate  func() uint8
	AddrMode func() uint8
	Cycle    uint8
}

type Cpu struct { // Olc6502 Cpu
	A, X, Y      uint8  // 寄存器
	StackP       uint8  // 栈指针
	Pc           uint16 // 程序计数器
	State        uint8  // 状态
	Bus          *Bus
	Fetched      uint8          // 指令获取到的结果
	Temp         uint16         //临时变量
	AddrAbs      uint16         //相对地址
	AddrRel      uint16         // 绝对地址
	OpCode       uint8          // 当前执行的指令
	Cycles       uint8          // 当前指令还有多少周期
	ClockCount   uint32         // 全局Clock计数
	Instructions []*Instruction // 所有可用指令
}

func (c *Cpu) ConnectBus(bus *Bus) {
	c.Bus = bus
}

func (c *Cpu) Reset() { // 重置
	// 重置程序计数器
	c.AddrAbs = 0xFFFC // 默认reset程序跳转的地址
	lo := c.Read(c.AddrAbs + 0)
	hi := c.Read(c.AddrAbs + 1)
	c.Pc = uint16(lo) | (uint16(hi) << 8)
	// 重置寄存器
	c.A = 0
	c.X = 0
	c.Y = 0
	c.StackP = 0xFD    // 栈基地址
	c.State = 0x00 | U // TODO 应该设置为 0x00也行
	// 重置地址
	c.AddrRel = 0x0000
	c.AddrAbs = 0x0000
	c.Fetched = 0x00
	// 重置循环次数
	c.Cycles = 8
}

func (c *Cpu) Nmi() { // 不可屏蔽的中断
	// 存储程序计数器
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc>>8))
	c.StackP--
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc))
	c.StackP--
	// 存储 cpu状态
	c.SetFlag(B, false)
	c.SetFlag(U, true)
	c.SetFlag(I, true)
	c.Write(0x0100+uint16(c.StackP), c.State)
	c.StackP--
	// 跳转到中断处理程序,注意于Irq相比处理中断程序的位置不同，且周期不同
	c.AddrAbs = 0xFFFA
	lo := c.Read(c.AddrAbs + 0)
	hi := c.Read(c.AddrAbs + 1)
	c.Pc = (uint16(hi) << 8) | uint16(lo)
	c.Cycles = 8
}

func (c *Cpu) Irq() { // 可以屏蔽的中断
	if c.GetFlag(I) {
		return
	}
	// 存储程序计数器
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc>>8))
	c.StackP--
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc))
	c.StackP--
	// 存储 cpu状态
	c.SetFlag(B, false)
	c.SetFlag(U, true)
	c.SetFlag(I, true)
	c.Write(0x0100+uint16(c.StackP), c.State)
	c.StackP--
	// 跳转到中断处理程序
	c.AddrAbs = 0xFFFE
	lo := c.Read(c.AddrAbs + 0)
	hi := c.Read(c.AddrAbs + 1)
	c.Pc = (uint16(hi) << 8) | uint16(lo)
	c.Cycles = 7
}

func (c *Cpu) Clock() {
	// 上一个指令执行完了，读取下一个指令执行，并等待对应时间
	if c.Cycles == 0 {
		c.OpCode = c.Read(c.Pc)
		c.SetFlag(U, true)
		c.Pc++
		instruction := c.Instructions[c.OpCode]
		c.Cycles = instruction.Cycle
		cycle1 := instruction.AddrMode()
		cycle2 := instruction.Operate()
		c.Cycles += cycle1 & cycle2 // TODO 为什么使用 &
		c.SetFlag(U, true)
	}
	// 模拟时钟行走
	c.ClockCount++
	c.Cycles--
}

func (c *Cpu) Complete() bool { // 当前指令是否执行完成
	return c.Cycles == 0
}

// 反编译指定范围内的代码
func (c *Cpu) Disassemble(start, end uint16) map[uint16]string {
	res := make(map[uint16]string, end-start)
	for addr := start; addr < end; addr++ {
		currentAddr := addr
		opCode := c.Read(addr)
		addr++
		current := "0x" + Format(currentAddr, 4) + ": "
		switch reflect.ValueOf(c.Instructions[opCode].AddrMode).UnsafePointer() {
		case reflect.ValueOf(c.IMP).UnsafePointer():
			current += "{IMP}"
		case reflect.ValueOf(c.IMM).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{IMM} 0x" + Format(value, 2)
		case reflect.ValueOf(c.ZP0).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{ZP0} 0x" + Format(value, 2)
		case reflect.ValueOf(c.ZPX).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{ZPX} 0x" + Format(value, 2)
		case reflect.ValueOf(c.ZPY).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{ZPY} 0x" + Format(value, 2)
		case reflect.ValueOf(c.IZX).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{IZX} 0x" + Format(value, 2)
		case reflect.ValueOf(c.IZY).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{IZY} 0x" + Format(value, 2)
		case reflect.ValueOf(c.ABS).UnsafePointer():
			lo := c.Read(addr)
			addr++
			hi := c.Read(addr)
			addr++
			current += "{ABS} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(c.ABX).UnsafePointer():
			lo := c.Read(addr)
			addr++
			hi := c.Read(addr)
			addr++
			current += "{ABX} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(c.ABY).UnsafePointer():
			lo := c.Read(addr)
			addr++
			hi := c.Read(addr)
			addr++
			current += "{ABY} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(c.IND).UnsafePointer():
			lo := c.Read(addr)
			addr++
			hi := c.Read(addr)
			addr++
			current += "{IND} 0x" + Format(uint16(lo)|(uint16(hi)<<8), 4)
		case reflect.ValueOf(c.REL).UnsafePointer():
			value := c.Read(addr)
			addr++
			current += "{REL} 0x" + Format(value, 2)
		}
		res[currentAddr] = current
	}
	return res
}

func (c *Cpu) GetFlag(key uint8) bool {
	return c.State&key > 0
}

func (c *Cpu) SetFlag(key uint8, flag bool) {
	if flag {
		c.State |= key
	} else {
		c.State &= ^key
	}
}

func (c *Cpu) Read(addr uint16) uint8 {
	return c.Bus.CpuRead(addr)
}

func (c *Cpu) Write(addr uint16, data uint8) {
	c.Bus.CpuWrite(addr, data)
}

func (c *Cpu) Fetch() uint8 {
	// IMP 取值方法会直接给 Fetched 赋值，其他都需要使用地址再获取最终值
	if reflect.ValueOf(c.Instructions[c.OpCode].AddrMode).UnsafePointer() != reflect.ValueOf(c.IMP).UnsafePointer() {
		c.Fetched = c.Read(c.AddrAbs)
	}
	return c.Fetched
}

//=============AddrMode=============
// 一共 256页，每页256b 组成  64k  若是操作码与取值操作都需要跨页则需要额外的时钟周期，这也是上面为什么使用 & 的原因

func (c *Cpu) IMP() uint8 {
	c.Fetched = c.A
	return 0
}

func (c *Cpu) IMM() uint8 {
	c.AddrAbs = c.Pc
	c.Pc++
	return 0
}

func (c *Cpu) ZP0() uint8 {
	c.AddrAbs = uint16(c.Read(c.Pc))
	c.Pc++
	return 0
}

func (c *Cpu) ZPX() uint8 {
	c.AddrAbs = uint16(c.Read(c.Pc)) + uint16(c.X)
	c.AddrAbs &= 0x00FF
	c.Pc++
	return 0
}

func (c *Cpu) ZPY() uint8 {
	c.AddrAbs = uint16(c.Read(c.Pc)) + uint16(c.Y)
	c.AddrAbs &= 0x00FF
	c.Pc++
	return 0
}

func (c *Cpu) REL() uint8 {
	c.AddrRel = uint16(c.Read(c.Pc))
	if c.AddrRel&0x80 > 0 { // 128
		c.AddrRel |= 0xFF00
	}
	c.Pc++
	return 0
}

func (c *Cpu) ABS() uint8 {
	lo := c.Read(c.Pc)
	c.Pc++
	hi := c.Read(c.Pc)
	c.Pc++
	c.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	return 0
}

func (c *Cpu) ABX() uint8 {
	lo := c.Read(c.Pc)
	c.Pc++
	hi := c.Read(c.Pc)
	c.Pc++
	c.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	c.AddrAbs += uint16(c.X)
	// 因为需要 额外加 x，可能高 8位发生变化 产生换页 耗时加1
	if c.AddrAbs&0xFF00 != uint16(hi)<<8 {
		return 1
	}
	return 0
}

func (c *Cpu) ABY() uint8 {
	lo := c.Read(c.Pc)
	c.Pc++
	hi := c.Read(c.Pc)
	c.Pc++
	c.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	c.AddrAbs += uint16(c.Y)
	// 因为需要 额外加 x，可能高 8位发生变化 产生换页 耗时加1
	if c.AddrAbs&0xFF00 != uint16(hi)<<8 {
		return 1
	}
	return 0
}

func (c *Cpu) IND() uint8 {
	// 先获取地址
	ptrLo := c.Read(c.Pc)
	c.Pc++
	ptrHi := c.Read(c.Pc)
	c.Pc++
	ptr := (uint16(ptrHi) << 8) | uint16(ptrLo)
	// 再获取地址对应的值 作为地址
	if ptrLo == 0xFF { // 模拟硬件上的一个bug相当于没有翻页，还是原页，页内地址变成了0
		c.AddrAbs = (uint16(c.Read(ptr&0xFF00)) << 8) | uint16(c.Read(ptr+0))
	} else {
		c.AddrAbs = (uint16(c.Read(ptr+1)) << 8) | uint16(c.Read(ptr+0))
	}
	return 0
}

func (c *Cpu) IZX() uint8 {
	// 在第0页进行寻址,只需要先找页内地址即可
	ptr := c.Read(c.Pc)
	c.Pc++
	// 需要根据地址找对应的值作为地址，一切都要在第0页进行
	lo := c.Read(uint16(ptr+c.X) & 0x00FF)
	hi := c.Read(uint16(ptr+c.X+1) & 0x00FF)
	c.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	return 0
}

func (c *Cpu) IZY() uint8 {
	// 与IZX类似宅0页操作，但是偏移Y是最后进行的，也需要判断是否发生页面跳转
	ptr := c.Read(c.Pc)
	c.Pc++
	lo := c.Read(uint16(ptr))
	hi := c.Read(uint16(ptr+1) & 0x00FF)
	c.AddrAbs = (uint16(hi) << 8) | uint16(lo)
	c.AddrAbs += uint16(c.Y)
	if c.AddrAbs&0xFF00 != uint16(hi)<<8 {
		return 1
	}
	return 0
}

//============Opcodes=============
//***********官方操作码***********

func (c *Cpu) ADC() uint8 {
	// 累加值到 Temp 并考虑进位 最终操作始终在 0~0x00FF
	c.Fetch()
	c.Temp = uint16(c.A) + uint16(c.Fetched)
	if c.GetFlag(C) {
		c.Temp++
	}
	// 设置计算后的各种标记
	c.SetFlag(C, c.Temp >= 0x00FF)
	c.SetFlag(Z, c.Temp&0x00FF == 0)
	// 计算溢出标记，先计算相关数据的负数标记
	tA := c.A&0x80 > 0
	tF := c.Fetched&0x80 > 0
	tT := c.Temp&0x80 > 0
	if tA && tF && !tT { // 负数+负数=正数 溢出
		c.SetFlag(V, true)
	} else if !tA && !tF && tT { // 正数+正数=负数 溢出
		c.SetFlag(V, true)
	} else { // 不会溢出
		c.SetFlag(V, false)
	}
	c.SetFlag(N, c.Temp&0x80 > 0)
	c.A = uint8(c.Temp) // 最终还是应用到累加器上
	return 1
}

func (c *Cpu) AND() uint8 {
	c.Fetch()
	c.A &= c.Fetched
	c.SetFlag(Z, c.A == 0x00)
	c.SetFlag(N, c.A&0x80 > 0)
	return 1
}

func (c *Cpu) ASL() uint8 {
	c.Fetch()
	c.Temp = uint16(c.Fetched) << 1
	c.SetFlag(C, c.Temp > 0x00FF)
	c.SetFlag(Z, c.Temp&0x00FF == 0x0000)
	c.SetFlag(N, c.Temp&0x80 > 0)
	if reflect.ValueOf(c.Instructions[c.OpCode].AddrMode).UnsafePointer() == reflect.ValueOf(c.IMP).UnsafePointer() {
		c.A = uint8(c.Temp)
	} else {
		c.Write(c.AddrAbs, uint8(c.Temp))
	}
	return 0
}

func (c *Cpu) BCC() uint8 {
	if !c.GetFlag(C) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BCS() uint8 {
	if c.GetFlag(C) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BEQ() uint8 {
	if c.GetFlag(Z) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BIT() uint8 {
	c.Fetch()
	c.Temp = uint16(c.A & c.Fetched)
	c.SetFlag(Z, c.Temp&0x00FF == 0)
	// 为什么突然使用 Fetched 了?
	c.SetFlag(N, c.Fetched&0x80 > 0)
	c.SetFlag(V, c.Fetched&0x40 > 0)
	return 0
}

func (c *Cpu) BMI() uint8 {
	if c.GetFlag(N) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BNE() uint8 {
	if !c.GetFlag(Z) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BPL() uint8 {
	if !c.GetFlag(N) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BRK() uint8 {
	// 保存指针与状态
	c.Pc++
	c.SetFlag(I, true)
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc>>8))
	c.StackP--
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc))
	c.StackP--
	c.SetFlag(B, true)
	c.Write(0x0100+uint16(c.StackP), c.State)
	c.StackP--
	c.SetFlag(B, false)
	// 跳转到中断处理程序
	c.Pc = uint16(c.Read(0xFFFF))<<8 | uint16(c.Read(0xFFFE))
	return 0
}

func (c *Cpu) BVC() uint8 {
	if !c.GetFlag(V) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) BVS() uint8 {
	if c.GetFlag(V) { // 分支判断 跳转到指定位置
		c.Cycles++
		c.AddrAbs = c.Pc + c.AddrRel // 并判断是否发生切页
		if c.AddrAbs&0xFF00 != c.Pc&0xFF00 {
			c.Cycles++
		}
		c.Pc = c.AddrAbs
	}
	return 0
}

func (c *Cpu) CLC() uint8 {
	c.SetFlag(C, false)
	return 0
}

func (c *Cpu) CLD() uint8 {
	c.SetFlag(D, false)
	return 0
}

func (c *Cpu) CLI() uint8 {
	c.SetFlag(I, false)
	return 0
}

func (c *Cpu) CLV() uint8 {
	c.SetFlag(V, false)
	return 0
}

func (c *Cpu) CMP() uint8 {
	c.Fetch()
	c.Temp = uint16(c.A) - uint16(c.Fetched)
	c.SetFlag(C, c.A >= c.Fetched)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	return 1
}

func (c *Cpu) CPX() uint8 {
	c.Fetch()
	c.Temp = uint16(c.X) - uint16(c.Fetched)
	c.SetFlag(C, c.X >= c.Fetched)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	return 0
}

func (c *Cpu) CPY() uint8 {
	c.Fetch()
	c.Temp = uint16(c.Y) - uint16(c.Fetched)
	c.SetFlag(C, c.Y >= c.Fetched)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	return 0
}

func (c *Cpu) DEC() uint8 {
	c.Fetch()
	c.Temp = uint16(c.Fetched - 1)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	return 0
}

func (c *Cpu) DEX() uint8 {
	c.X--
	c.SetFlag(Z, c.X == 0)
	c.SetFlag(N, c.X&0x80 > 0)
	return 0
}

func (c *Cpu) DEY() uint8 {
	c.Y--
	c.SetFlag(Z, c.Y == 0)
	c.SetFlag(N, c.Y&0x80 > 0)
	return 0
}

func (c *Cpu) EOR() uint8 {
	c.Fetch()
	c.A ^= c.Fetched
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 1
}

func (c *Cpu) INC() uint8 {
	c.Fetch()
	c.Temp = uint16(c.Fetched + 1)
	c.Write(c.AddrAbs, uint8(c.Temp))
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	return 0
}

func (c *Cpu) INX() uint8 {
	c.X++
	c.SetFlag(Z, c.X&0xFF == 0)
	c.SetFlag(N, c.X&0x80 > 0)
	return 0
}

func (c *Cpu) INY() uint8 {
	c.Y++
	c.SetFlag(Z, c.Y&0xFF == 0)
	c.SetFlag(N, c.Y&0x80 > 0)
	return 0
}

func (c *Cpu) JMP() uint8 {
	c.Pc = c.AddrAbs
	return 0
}

func (c *Cpu) JSR() uint8 {
	// 先保存pc指针再跳转
	c.Pc--
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc>>8))
	c.StackP--
	c.Write(0x0100+uint16(c.StackP), uint8(c.Pc))
	c.StackP--
	c.Pc = c.AddrAbs
	return 0
}

func (c *Cpu) LDA() uint8 {
	c.Fetch()
	c.A = c.Fetched
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 1
}

func (c *Cpu) LDX() uint8 {
	c.Fetch()
	c.X = c.Fetched
	c.SetFlag(Z, c.X == 0)
	c.SetFlag(N, c.X&0x80 > 0)
	return 1
}

func (c *Cpu) LDY() uint8 {
	c.Fetch()
	c.Y = c.Fetched
	c.SetFlag(Z, c.Y == 0)
	c.SetFlag(N, c.Y&0x80 > 0)
	return 1
}

func (c *Cpu) LSR() uint8 { // 移位运算符
	c.Fetch()
	c.SetFlag(C, c.Fetched&0x1 > 0)
	c.Temp = uint16(c.Fetched >> 1)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	if reflect.ValueOf(c.Instructions[c.OpCode].AddrMode).UnsafePointer() == reflect.ValueOf(c.IMP).UnsafePointer() {
		c.A = uint8(c.Temp)
	} else {
		c.Write(c.AddrAbs, uint8(c.Temp))
	}
	return 0
}

func (c *Cpu) NOP() uint8 {
	switch c.OpCode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return 1
	}
	return 0
}

func (c *Cpu) ORA() uint8 {
	c.Fetch()
	c.A |= c.Fetched
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 1
}

func (c *Cpu) PHA() uint8 {
	c.Write(0x0100+uint16(c.StackP), c.A)
	c.StackP--
	return 0
}

func (c *Cpu) PHP() uint8 {
	c.Write(0x0100+uint16(c.StackP), c.State|B|U)
	c.StackP--
	c.SetFlag(B, false)
	c.SetFlag(U, false)
	return 0
}

func (c *Cpu) PLA() uint8 { // 弹出栈
	c.StackP++
	c.A = c.Read(0x0100 + uint16(c.StackP))
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 0
}

func (c *Cpu) PLP() uint8 {
	c.StackP++
	c.State = c.Read(0x0100 + uint16(c.StackP))
	c.SetFlag(U, true)
	return 0
}

func (c *Cpu) ROL() uint8 {
	c.Fetch()
	c.Temp = uint16(c.Fetched << 1)
	if c.GetFlag(C) {
		c.Temp |= 0x0001
	}
	c.SetFlag(C, c.Temp > 0xFF)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	if reflect.ValueOf(c.Instructions[c.OpCode].AddrMode).UnsafePointer() == reflect.ValueOf(c.IMP).UnsafePointer() {
		c.A = uint8(c.Temp)
	} else {
		c.Write(c.AddrAbs, uint8(c.Temp))
	}
	return 0
}

func (c *Cpu) ROR() uint8 { // 位移操作
	c.Fetch()
	c.Temp = uint16(c.Fetched) >> 1
	if c.GetFlag(C) {
		c.Temp |= 0x80
	}
	c.SetFlag(C, c.Temp&0x01 > 0)
	c.SetFlag(Z, c.Temp&0xFF == 0)
	c.SetFlag(N, c.Temp&0x80 > 0)
	if reflect.ValueOf(c.Instructions[c.OpCode].AddrMode).UnsafePointer() == reflect.ValueOf(c.IMP).UnsafePointer() {
		c.A = uint8(c.Temp)
	} else {
		c.Write(c.AddrAbs, uint8(c.Temp))
	}
	return 0
}

func (c *Cpu) RTI() uint8 {
	// 恢复到原来的状态与地址
	c.StackP++
	c.State = c.Read(0x0100 + uint16(c.StackP))
	c.SetFlag(B, false)
	c.SetFlag(U, false)
	c.StackP++
	lo := c.Read(0x0100 + uint16(c.StackP))
	c.StackP++
	hi := c.Read(0x0100 + uint16(c.StackP))
	c.Pc = uint16(lo) | (uint16(hi) << 8)
	return 0
}

func (c *Cpu) RTS() uint8 { // 只恢复指针
	c.StackP++
	lo := c.Read(0x0100 + uint16(c.StackP))
	c.StackP++
	hi := c.Read(0x0100 + uint16(c.StackP))
	c.Pc = uint16(lo) | (uint16(hi) << 8)
	c.Pc++
	return 0
}

func (c *Cpu) SBC() uint8 {
	// 累减值到 Temp 并考虑进位 最终操作始终在 0~0x00FF
	c.Fetch()
	c.Temp = uint16(c.A) - uint16(c.Fetched)
	if c.GetFlag(C) {
		c.Temp++
	}
	// 设置计算后的各种标记
	c.SetFlag(C, c.Temp >= 0x00FF)
	c.SetFlag(Z, c.Temp&0x00FF == 0)
	// 计算溢出标记，先计算相关数据的负数标记
	tA := c.A&0x80 > 0
	tF := c.Fetched&0x80 > 0
	tT := c.Temp&0x80 > 0
	if tA && !tF && !tT { // 负数-正数=正数 溢出
		c.SetFlag(V, true)
	} else if !tA && tF && tT { // 正数-负数=负数 溢出
		c.SetFlag(V, true)
	} else { // 不会溢出
		c.SetFlag(V, false)
	}
	c.SetFlag(N, c.Temp&0x80 > 0)
	c.A = uint8(c.Temp) // 最终还是应用到累加器上
	return 1
}

func (c *Cpu) SEC() uint8 {
	c.SetFlag(C, true)
	return 0
}

func (c *Cpu) SED() uint8 {
	c.SetFlag(D, true)
	return 0
}

func (c *Cpu) SEI() uint8 {
	c.SetFlag(I, true)
	return 0
}

func (c *Cpu) STA() uint8 {
	c.Write(c.AddrAbs, c.A)
	return 0
}

func (c *Cpu) STX() uint8 {
	c.Write(c.AddrAbs, c.X)
	return 0
}

func (c *Cpu) STY() uint8 {
	c.Write(c.AddrAbs, c.Y)
	return 0
}

func (c *Cpu) TAX() uint8 {
	c.X = c.A
	c.SetFlag(Z, c.X == 0)
	c.SetFlag(N, c.X&0x80 > 0)
	return 0
}

func (c *Cpu) TAY() uint8 {
	c.Y = c.A
	c.SetFlag(Z, c.Y == 0)
	c.SetFlag(N, c.Y&0x80 > 0)
	return 0
}

func (c *Cpu) TSX() uint8 {
	c.X = c.StackP
	c.SetFlag(Z, c.X == 0)
	c.SetFlag(N, c.X&0x80 > 0)
	return 0
}

func (c *Cpu) TXA() uint8 {
	c.A = c.X
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 0
}

func (c *Cpu) TXS() uint8 {
	c.StackP = c.X
	return 0
}

func (c *Cpu) TYA() uint8 {
	c.A = c.Y
	c.SetFlag(Z, c.A == 0)
	c.SetFlag(N, c.A&0x80 > 0)
	return 0
}

//*********非官方操作码统一使用这个***********

func (c *Cpu) XXX() uint8 {
	return 0
}

func NewCpu() *Cpu {
	res := &Cpu{}
	// 各种指令的映射关系
	res.Instructions = []*Instruction{
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
