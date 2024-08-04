/*
@author: sk
@date: 2023/11/12
*/
package main

type Mapper struct {
	PRGBanks uint8
	CHRBanks uint8
}

func NewMapper(prgBanks uint8, chrBanks uint8) *Mapper {
	return &Mapper{PRGBanks: prgBanks, CHRBanks: chrBanks}
}

func (m *Mapper) CpuReadMap(addr uint16) (uint16, bool) {
	if addr >= 0x8000 { // 转换到对应的合法区域
		mask := uint16(0x3FFF)
		if m.PRGBanks > 1 {
			mask = 0x7FFF
		}
		return addr & mask, true
	}
	return 0, false
}

func (m *Mapper) CpuWriteMap(addr uint16) (uint16, bool) {
	if addr >= 0x8000 {
		mask := uint16(0x3FFF)
		if m.PRGBanks > 1 {
			mask = 0x7FFF
		}
		return addr & mask, true
	}
	return 0, false
}

func (m *Mapper) PpuReadMap(addr uint16) (uint16, bool) { // 转换到对应的地址
	if addr < 0x2000 {
		return addr, true
	}
	return 0, false
}

func (m *Mapper) PpuWriteMap(addr uint16) (uint16, bool) { // ppu不应该写入
	return 0, false
}
