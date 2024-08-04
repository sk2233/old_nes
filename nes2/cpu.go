/*
@author: sk
@date: 2023/10/14
*/
package main

type Cpu struct {
	bus     *Bus
	pc      uint16
	sp      uint8 // 0xff->0x00间循环 栈指针 这里只是存储偏移
	a, x, y uint8
	state   uint8
	addrs   map[uint8]AddrFunc
	ops     map[uint8]OpFunc
	cmds    map[uint8]*Cmd
}

func NewCpu(bus *Bus) *Cpu {
	return &Cpu{
		bus:   bus,
		addrs: loadAddrs(),
		ops:   loadOps(),
		cmds:  loadCmds(),
	}
}

func (c *Cpu) Irq() { // 普通中断
	c.pc = ReadUnit16(c.bus.CpuRead, 0xFFFE, 0xFFFF)
}

func (c *Cpu) Nmi() { // 不可屏蔽中断
	c.pc = ReadUnit16(c.bus.CpuRead, 0xFFFA, 0xFFFB)
}

func (c *Cpu) Reset() {
	c.pc = ReadUnit16(c.bus.CpuRead, 0xFFFC, 0xFFFD)
}

func (c *Cpu) Rti() int { // 从中断恢复
	// 恢复状态
	c.state = c.Pop()
	// 恢复PC
	c.pc = c.PopUint16()
	return 0
}

func (c *Cpu) Pop() uint8 {
	c.sp++
	c.sp &= 0xFF
	return c.bus.CpuRead(0x0100 + uint16(c.sp))
}

func (c *Cpu) PopUint16() uint16 {
	lo := c.Pop()
	hi := c.Pop()
	return uint16(lo) | (uint16(hi) << 8)
}
