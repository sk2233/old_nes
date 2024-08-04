/*
@author: sk
@date: 2024/2/25
*/
package nes5

import "os"

type MirrorType int

const (
	HORIZONTAL MirrorType = 1
	VERTICAL   MirrorType = 2
)

type Cartridge struct {
	Mirror   MirrorType
	PRGMem   []uint8
	CHRMem   []uint8
	MapperID uint8
	PRGBank  uint8
	CHRBank  uint8
	Mapper   *Mapper0
}

type header struct {
	Name                     string
	PrgRomChunk, ChrRomChunk uint8
	Mapper1, Mapper2         uint8
	PrgRamSize               uint8
	TvSys1, TvSys2           uint8
	Unused                   string
}

func parseHeader(reader *os.File) *header {
	return &header{
		Name:        ReadStr(reader, 4),
		PrgRomChunk: ReadUint8(reader),
		ChrRomChunk: ReadUint8(reader),
		Mapper1:     ReadUint8(reader),
		Mapper2:     ReadUint8(reader),
		PrgRamSize:  ReadUint8(reader),
		TvSys1:      ReadUint8(reader),
		TvSys2:      ReadUint8(reader),
		Unused:      ReadStr(reader, 5),
	}
}

func NewCartridge(file string) *Cartridge {
	reader := OpenFile(file)
	defer reader.Close()
	temp := parseHeader(reader)
	if (temp.Mapper1 & 0x04) > 0 {
		ReadBytes(reader, 512)
	}
	//return &Cartridge{
	//	Mirror:    If((header.Mapper1&0x01) > 0, VERTICAL, HORIZONTAL),
	//	MapperID:  (header.Mapper1 >> 4) | (header.Mapper2 & 0xF0),
	//	PRGBanks:  header.PrgRomChunks,
	//	PRGMemory: ReadBytes(reader, int(header.PrgRomChunks)*16*1024),
	//	CHRBanks:  header.ChrRomChunks,
	//	CHRMemory: ReadBytes(reader, int(header.ChrRomChunks)*8*1024),
	//	Mapper:    NewMapper(header.PrgRomChunks, header.ChrRomChunks),
	//}
	return &Cartridge{
		Mirror:   If((temp.Mapper1&0x01) > 0, VERTICAL, HORIZONTAL),
		MapperID: (temp.Mapper1 >> 4) | ((temp.Mapper2 >> 4) << 4), // 必须是0 当前使用的是0的Mapper
		PRGMem:   ReadBytes(reader, int(temp.PrgRomChunk)*16*1024),
		CHRMem:   ReadBytes(reader, int(temp.ChrRomChunk)*8*1024),
		Mapper:   NewMapper0(temp.PrgRomChunk, temp.ChrRomChunk),
	}
}

func (c *Cartridge) CpuRead(addr uint16) (uint8, bool) {
	if mapAddr, ok := c.Mapper.CpuMapRead(addr); ok {
		return c.PRGMem[mapAddr], true
	}
	return 0, false
}

func (c *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	if mapAddr, ok := c.Mapper.CpuMapWrite(addr); ok {
		c.PRGMem[mapAddr] = data
		return true
	}
	return false
}

func (c *Cartridge) PpuRead(addr uint16) (uint8, bool) {
	if mapAddr, ok := c.Mapper.PpuMapRead(addr); ok {
		return c.CHRMem[mapAddr], true
	}
	return 0, false
}

func (c *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	if mapAddr, ok := c.Mapper.PpuMapWrite(addr); ok {
		c.CHRMem[mapAddr] = data
		return true
	}
	return false
}

func (c *Cartridge) Reset() {
	c.Mapper.Reset()
}
