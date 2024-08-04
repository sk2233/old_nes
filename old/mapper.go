/*
@author: sk
@date: 2023/10/14
*/
package main

type Mapper struct {
	PrgBankCount uint8
	ChrBankCount uint8
}

func NewMapper(prgBankCount uint8, chrBankCount uint8) *Mapper {
	return &Mapper{PrgBankCount: prgBankCount, ChrBankCount: chrBankCount}
}

func (m *Mapper) CpuRead(addr uint16) (uint16, bool) {
	if addr >= 0x8000 { // 只对 64k 的后32k感兴趣
		// 如果多于1个程序页，就是使用 32k rom (其他寻址地址都是镜像) 取余 32k   否则就是用 16k rom (其他寻址地址都是镜像) 取余  16k
		mask := If(m.PrgBankCount > 1, 0x7fff, 0x3fff)
		return addr & uint16(mask), true
	}
	return 0, false
}

func (m *Mapper) CpuWrite(addr uint16) (uint16, bool) {
	if addr > 0x8000 {
		mask := If(m.PrgBankCount > 1, 0x7fff, 0x3fff)
		return addr & uint16(mask), true
	}
	return 0, false
}

func (m *Mapper) PpuRead(addr uint16) (uint16, bool) {
	if addr < 0x1fff { // ppu只读 64k的 前面8k
		return addr, true
	}
	return 0, false
}

func (m *Mapper) PpuWrite(addr uint16) (uint16, bool) { // PPU NOT WRITE
	return 0, false
}
