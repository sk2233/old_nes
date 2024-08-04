/*
@author: sk
@date: 2023/11/11
*/
package main

type Bus struct {
	Cpu        *Cpu
	Ppu        *Ppu
	CpuRam     []byte
	ClockCount uint32
	Cartridge  *Cartridge
}

func NewBus() *Bus {
	bus := &Bus{CpuRam: make([]byte, 2*KB), Cpu: NewCpu(), Ppu: NewPpu()}
	bus.Cpu.ConnectBus(bus)
	return bus
}

func (b *Bus) InsertCartridge(cartridge *Cartridge) {
	b.Cartridge = cartridge
	b.Ppu.ConnectCartridge(cartridge)
}

func (b *Bus) Reset() {
	b.Cpu.Reset()
	b.ClockCount = 0
}

func (b *Bus) Clock() {
	b.Ppu.Clock()
	if b.ClockCount%3 == 0 { // cpu频率比ppu慢3倍
		b.Cpu.Clock()
	}
	if b.Ppu.Nmi {
		b.Ppu.Nmi = false
		b.Cpu.Nmi()
	}
	b.ClockCount++
}

func (b *Bus) CpuRead(addr uint16) byte { // ReadOnly = false
	if data, ok := b.Cartridge.CpuRead(addr); ok { // 卡带拦截
		return data
	} else if addr < 0x2000 {
		return b.CpuRam[addr&0x07FF] // 进行镜像
	} else if addr < 0x4000 {
		return b.Ppu.CpuRead(addr & 0x0007)
	}
	return 0
}

func (b *Bus) CpuWrite(addr uint16, data byte) {
	if b.Cartridge.CpuWrite(addr, data) {
		// 卡带作为拦截器处理
		return
	} else if addr < 0x2000 {
		b.CpuRam[addr&0x07FF] = data // 同样镜像
	} else if addr < 0x4000 {
		b.Ppu.CpuWrite(addr&0x007, data) // 只可写入8个
	}
}

func (b *Bus) ReadUint16(loAddr, hiAddr uint16) uint16 { // 先低位，再高位
	lo := b.CpuRead(loAddr)
	hi := b.CpuRead(hiAddr)
	return uint16(lo) | (uint16(hi) << 8)
}
