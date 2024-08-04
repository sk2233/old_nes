/*
@author: sk
@date: 2023/10/26
*/
package nes4

import (
	"os"
)

type MirrorType int

const (
	HORIZONTAL MirrorType = 1
	VERTICAL   MirrorType = 2
)

type header struct {
	Name                       string
	PrgRomChunks, ChrRomChunks uint8
	Mapper1, Mapper2           uint8
	PrgRamSize                 uint8
	TvSys1, TvSys2             uint8
	Other                      string
}

type Cartridge struct {
	Mirror               MirrorType
	MapperID             uint8
	PRGBanks, CHRBanks   uint8
	PRGMemory, CHRMemory []uint8
	Mapper               *Mapper
}

func NewCartridge(nesFile string) *Cartridge {
	reader := OpenFile(nesFile)
	defer reader.Close()
	header := ReadHeader(reader)
	if (header.Mapper1 & 0x04) > 0 {
		ReadBytes(reader, 512)
	}
	return &Cartridge{
		Mirror:    If((header.Mapper1&0x01) > 0, VERTICAL, HORIZONTAL),
		MapperID:  (header.Mapper1 >> 4) | (header.Mapper2 & 0xF0),
		PRGBanks:  header.PrgRomChunks,
		PRGMemory: ReadBytes(reader, int(header.PrgRomChunks)*16*1024),
		CHRBanks:  header.ChrRomChunks,
		CHRMemory: ReadBytes(reader, int(header.ChrRomChunks)*8*1024),
		Mapper:    NewMapper(header.PrgRomChunks, header.ChrRomChunks),
	}
}

func ReadHeader(reader *os.File) *header {
	return &header{
		Name:         ReadStr(reader, 4),
		PrgRomChunks: ReadUint8(reader),
		ChrRomChunks: ReadUint8(reader),
		Mapper1:      ReadUint8(reader),
		Mapper2:      ReadUint8(reader),
		PrgRamSize:   ReadUint8(reader),
		TvSys1:       ReadUint8(reader),
		TvSys2:       ReadUint8(reader),
		Other:        ReadStr(reader, 5),
	}
}

func (c *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	if addr, ok := c.Mapper.CpuMapWrite(addr, data); ok {
		c.PRGMemory[addr] = data
		return true
	}
	return false
}

func (c *Cartridge) CpuRead(addr uint16) (uint8, bool) {
	if addr, ok := c.Mapper.CpuMapRead(addr); ok {
		return c.PRGMemory[addr], true
	}
	return 0, false
}

func (c *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	if addr, ok := c.Mapper.PpuMapWrite(addr, data); ok {
		c.CHRMemory[addr] = data
		return true
	}
	return false
}

func (c *Cartridge) PpuRead(addr uint16) (uint8, bool) {
	if addr, ok := c.Mapper.PpuMapRead(addr); ok {
		return c.CHRMemory[addr], true
	}
	return 0, false
}

func (c *Cartridge) Reset() {
	c.Mapper.Reset()
}
