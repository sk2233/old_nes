/*
@author: sk
@date: 2023/10/9
*/
package main

type Bus struct {
	Cpu          *Cpu
	Ppu          *Ppu
	CpuRam       [2 * 1024]uint8 // 使用固定长度的数组
	PpuRam       [2 * 1024]uint8
	ClockCounter uint64
	Cartridge    *Cartridge
}

func NewBus() *Bus {
	return &Bus{
		CpuRam: [2 * 1024]uint8{},
		PpuRam: [2 * 1024]uint8{},
	}
}

// CPU 可以寻址的空间是  8k 但是 实际只有2k 其他部分都是前面2k的镜像
func (b *Bus) Read(addr uint16) uint8 {
	if addr < 0x2000 { // CPU 地址 0 ~ 8k  前2k循环
		addr %= 2 * 1024
		return b.CpuRam[addr]
	} else if addr < 0x4000 { // PPU 地址 8k ~ 16k 前2k循环
		addr -= 0x2000
		addr %= 2 * 1024
		return b.PpuRam[addr]
	}
	return 0
}

func (b *Bus) Write(addr uint16, data uint8) {
	if addr < 0x2000 {
		addr %= 2 * 1024
		b.CpuRam[addr] = data
	} else if addr < 0x4000 {
		addr -= 0x2000
		addr %= 2 * 1024
		b.PpuRam[addr] = data
	}
}

func (b *Bus) SetCartridge(cartridge *Cartridge) {
	b.Cartridge = cartridge
	b.Ppu.SetCartridge(cartridge)

}

func (b *Bus) Reset() {
	b.Cpu.Reset()
	b.ClockCounter = 0
}

func (b *Bus) Clock() {
	b.Ppu.Clock() // ppu 频率是 cpu的3倍
	if b.ClockCounter%3 == 0 {
		b.Cpu.Clock()
	}
	b.ClockCounter++
}
