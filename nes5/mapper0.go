/*
@author: sk
@date: 2024/2/25
*/
package nes5

type Mapper0 struct {
	PRGBank uint8
	CHRBank uint8
}

func NewMapper0(prgBank uint8, chrBank uint8) *Mapper0 {
	return &Mapper0{PRGBank: prgBank, CHRBank: chrBank}
}

func (m *Mapper0) CpuMapRead(addr uint16) (uint16, bool) {
	if addr >= 0x8000 {
		if m.PRGBank > 1 {
			return addr & 0x7FFF, true
		}
		return addr & 0x3FFF, true
	}
	return 0, false
}

func (m *Mapper0) CpuMapWrite(addr uint16) (uint16, bool) {
	if addr >= 0x8000 {
		if m.PRGBank > 1 {
			return addr & 0x7FFF, true
		}
		return addr & 0x3FFF, true
	}
	return 0, false
}

func (m *Mapper0) PpuMapRead(addr uint16) (uint16, bool) {
	if addr < 0x2000 {
		return addr, true
	}
	return 0, false
}

func (m *Mapper0) PpuMapWrite(addr uint16) (uint32, bool) {
	//if addr < 0x2000 {
	//	return 0, true
	//}
	return 0, false
}

func (m *Mapper0) Reset() {

}
