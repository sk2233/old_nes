/*
@author: sk
@date: 2024/2/24
*/
package nes5

type Bus struct {
	Cpu           *Olc6502
	Ppu           *Olc2C02
	Cartridge     *Cartridge
	CpuRam        [2 * 1024]uint8
	SysClockCount uint32
}

func NewBus() *Bus {
	res := &Bus{
		Cpu: NewOlc6502(),
		Ppu: NewOlc2C02(),
	}
	res.Cpu.ConnectBus(res)
	return res
}

func (b *Bus) CpuWrite(addr uint16, data uint8) {
	if b.Cartridge.CpuWrite(addr, data) {

	} else if addr < 0x2000 {
		b.CpuRam[addr&0x07FF] = data
	} else if addr < 0x4000 {
		b.Ppu.CpuWrite(addr&0x0007, data)
	}
}

// read = true 主要是给反编译使用的，因为read是有损的
func (b *Bus) CpuRead(addr uint16, read bool) uint8 {
	if res, ok := b.Cartridge.CpuRead(addr); ok {
		return res
	} else if addr < 0x2000 {
		return b.CpuRam[addr&0x07FF]
	} else if addr < 0x4000 {
		return b.Ppu.CpuRead(addr&0x0007, read)
	}
	return 0
}

func (b *Bus) InsertCartridge(cartridge *Cartridge) {
	b.Cartridge = cartridge
	b.Ppu.ConnectCartridge(cartridge)
}

func (b *Bus) Reset() {
	b.Cartridge.Reset()
	b.Cpu.Reset()
	b.Ppu.Reset()
	b.SysClockCount = 0
}

func (b *Bus) Clock() {
	b.Ppu.Clock()
	if b.SysClockCount%3 == 0 {
		b.Cpu.Clock()
	}
	if b.Ppu.Nmi {
		b.Ppu.Nmi = false
		b.Cpu.Nmi()
	}
	b.SysClockCount++
}
