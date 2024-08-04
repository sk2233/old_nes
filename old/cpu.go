/*
@author: sk
@date: 2023/10/9
*/
package main

import (
	"math/rand"
	"strconv"
)

type Cpu struct {
	// 读地址
	Bus *Bus
	// 主要寄存器
	A, X, Y uint8
	StackP  uint8
	Pc      uint16
	State   uint8
	// 寻址方式
	AddrModes map[uint8]func() uint8
	// 操作处理
	OpCodes map[uint8]func() uint8
	// 消耗时间
	CostClocks map[uint8]uint8
	// 机器时钟
	Cycles uint8
	// 暂存寻址后的绝对地址
	Addr uint16
	// 从绝对地址中获取到的值
	Fetched uint8
}

func NewCpu(bus *Bus) *Cpu {
	return &Cpu{Bus: bus,
		AddrModes: map[uint8]func() uint8{}, OpCodes: map[uint8]func() uint8{}, CostClocks: map[uint8]uint8{}}
}

//=============AddrModes=================

func (o *Cpu) ZP0() uint8 {
	o.Addr = uint16(o.Read(o.Pc))
	o.Pc++
	return 0
}

func (o *Cpu) ZPX() uint8 {
	o.Addr = uint16(o.Read(o.Pc) + o.X)
	o.Pc++
	return 0
}

func (o *Cpu) ZPY() uint8 {
	o.Addr = uint16(o.Read(o.Pc) + o.Y)
	o.Pc++
	return 0
}

func (o *Cpu) ABS() uint8 { // 直接找出地址
	lo := o.Read(o.Pc)
	o.Pc++
	hi := o.Read(o.Pc)
	o.Pc++
	o.Addr = uint16(hi)<<8 | uint16(lo)
	return 0
}

func (o *Cpu) IND() uint8 {
	lo := o.Read(o.Pc)
	o.Pc++
	hi := o.Read(o.Pc)
	o.Pc++
	addr := uint16(hi)<<8 | uint16(lo)
	// 间接寻址，先找到对应的地址，再找对应的数据，可以封装出 Read8,Read16 自动偏移指针,CpuRead(uint16)的方法方便使用
	lo = o.Read(addr)
	hi = o.Read(addr + 1)
	o.Addr = uint16(hi)<<8 | uint16(lo)
	return 0
}

//==============OpCode==============

func (o *Cpu) AND() uint8 {
	o.Fetch()
	o.A &= o.Fetched
	o.SetState(CpuStateZ, o.A == 0)
	o.SetState(CpuStateN, (o.A&0x80) > 0)
	return 1
}

func (o *Cpu) ADC() int { // 减法，可以使用加法加 补码实现
	o.Fetch()
	temp := uint16(o.A) + uint16(o.Fetched)
	if o.HasState(CpuStateC) {
		temp++
	}
	o.SetState(CpuStateC, temp > 0xff)                                 // 设置是否需要进位
	o.SetState(CpuStateZ, temp&0xff == 0)                              // 是否为0值
	o.SetState(CpuStateN, temp&0x80 > 0)                               // 是否为负数
	o.SetState(CpuStateV, int16(o.A)*int16(o.Fetched)*int16(temp) < 0) // 是否溢出
	o.A = uint8(temp & 0xff)                                           // 赋值给加法器
	return 1
}

func (o *Cpu) PHA() int { // 压栈操作
	o.Write(0x0100+uint16(o.StackP), o.A) // 写入数据到栈空间，0x0100 是固定栈空间起始点
	o.StackP--                            // 与普通栈相比方向是反的
	return 0
}

func (o *Cpu) PLA() { // 弹栈操作
	o.StackP++
	o.A = o.Read(0x0100 + uint16(o.StackP)) // 把值再放回 寄存器 A
	o.SetState(CpuStateZ, o.A == 0)
	o.SetState(CpuStateN, o.A&0x80 > 0)
}

func (o *Cpu) Fetch() {
	opcode := strconv.Itoa(rand.Intn(100))
	if opcode != "IMP" {
		// IMP 指令是空指令 不会修改 Fetched
		o.Fetched = o.Read(o.Addr)
	}
}

func (o *Cpu) Read(addr uint16) uint8 {
	return o.Bus.Read(addr)
}

func (o *Cpu) Write(addr uint16, data uint8) {
	o.Bus.Write(addr, data)
}

func (o *Cpu) SetState(state CpuState, ok bool) {
	if ok {
		o.State |= uint8(state)
	} else {
		o.State &= ^uint8(state)
	}
}

func (o *Cpu) HasState(state CpuState) bool {
	return o.State&uint8(state) > 0
}

func (o *Cpu) Clock() {
	if o.Cycles == 0 {
		opCode := o.Read(o.Pc)
		o.Pc++
		cycles := o.CostClocks[opCode]
		cycles1 := o.AddrModes[opCode]()
		cycles2 := o.OpCodes[opCode]()
		cycles += cycles1 & cycles2
		o.Cycles = cycles
	}
	o.Cycles--
}

func (o *Cpu) Reset() { // 复位电源
	o.A = 0
	o.X = 0
	o.Y = 0
	o.StackP = 0xFD            // 固定栈空间位置
	o.State = uint8(CpuStateU) // 设置为 为使用的
	lo := o.Read(0xfffc)       // 默认程序入口地址通过一个固定地址获取
	hi := o.Read(0xfffd)
	o.Pc = uint16(hi)<<8 | uint16(lo)
	o.Addr = 0x0000
	o.Fetched = 0x00
	o.Cycles = 8
}

func (o *Cpu) Irq() { // 中断请求
	if o.HasState(CpuStateI) { // 可以屏蔽的中断请求
		return
	}
	o.Nmi()
}

func (o *Cpu) Nmi() { // 不可屏蔽的中断请求
	// 保存pc
	o.Write(0x0100+uint16(o.StackP), uint8((o.Pc>>8)&0x00ff))
	o.StackP--
	o.Write(0x0100+uint16(o.StackP), uint8(o.Pc&0x00ff))
	o.StackP--
	// 设置相关状态
	o.SetState(CpuStateB, false)
	o.SetState(CpuStateU, true)
	o.SetState(CpuStateI, true)
	// 保存状态
	o.Write(0x0100+uint16(o.StackP), o.State)
	o.StackP--
	// 调整到中断处理地址
	lo := o.Read(0xfffe)
	hi := o.Read(0xffff)
	o.Pc = uint16(hi)<<8 | uint16(lo)
	o.Cycles = 7
}

func (o *Cpu) RTI() int { // 从中断恢复
	// 恢复状态
	o.StackP++
	o.State = o.Read(0x0100 + uint16(o.StackP))
	o.SetState(CpuStateB, false)
	o.SetState(CpuStateU, false)
	// 恢复PC
	o.StackP++
	lo := o.Read(0x0100 + uint16(o.StackP))
	o.StackP++
	hi := o.Read(0x0100 + uint16(o.StackP))
	o.Pc = uint16(hi)<<8 | uint16(lo)
	return 0
}
