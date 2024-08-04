/*
@author: sk
@date: 2023/10/27
*/
package nes4

type Mapper struct {
	PRGBanks, CHRBanks uint8
}

func (m *Mapper) CpuMapRead(addr uint16) (uint16, bool) {
	if addr >= 0x8000 { // 0x8000~0xFFFF
		if m.PRGBanks > 1 { // 对于 32k的 只有前 32k是内容，其他都是镜像
			return addr & 0x7FFF, true
		} else { // 对于 16k的 只有前 16k是内容，其他都是镜像
			return addr & 0x3FFF, true
		}
	}
	return 0, false
}

func (m *Mapper) CpuMapWrite(addr uint16, data uint8) (uint16, bool) {
	if addr >= 0x8000 { // 暂时只做地址转换 与 CpuMapRead 一致
		if m.PRGBanks > 1 { // 对于 32k的 只有前 32k是内容，其他都是镜像
			return addr & 0x7FFF, true
		} else { // 对于 16k的 只有前 16k是内容，其他都是镜像
			return addr & 0x3FFF, true
		}
	}
	return 0, false
}

func (m *Mapper) PpuMapRead(addr uint16) (uint16, bool) {
	if addr <= 0x1FFF { // ppu占8k且无需映射  范围  0~0x1FFF
		return addr, true
	}
	return 0, false
}

func (m *Mapper) PpuMapWrite(addr uint16, data uint8) (uint16, bool) {
	if addr <= 0x1FFF { // ppu占8k且无需映射
		return addr, true
	}
	return 0, false
}

func (m *Mapper) Reset() {
}

func NewMapper(prgBanks, chrBanks uint8) *Mapper { // 调用 reset了
	return &Mapper{PRGBanks: prgBanks, CHRBanks: chrBanks}
}
