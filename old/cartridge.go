/*
@author: sk
@date: 2023/10/12
*/
package main

import (
	"fmt"
	"os"
)

type Cartridge struct {
	PrgRom       []uint8
	ChrRom       []uint8
	MapperID     uint8
	PrgBankCount uint8
	ChrBankCount uint8
	Mapper       *Mapper
}

func (c *Cartridge) CpuRead(addr uint16) (uint8, bool) {
	if addr, ok := c.Mapper.CpuRead(addr); ok {
		return c.PrgRom[addr], true
	}
	return 0, false
}

func (c *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	if addr, ok := c.Mapper.CpuWrite(addr); ok {
		c.PrgRom[addr] = data
		return true
	}
	return false
}

func (c *Cartridge) PpuRead(addr uint16) (uint8, bool) {
	if addr, ok := c.Mapper.PpuRead(addr); ok {
		return c.ChrRom[addr], true
	}
	return 0, false
}

func (c *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	if addr, ok := c.Mapper.PpuWrite(addr); ok {
		c.ChrRom[addr] = data
		return true
	}
	return false
}

func NewCartridge(nesFile string) *Cartridge {
	file, err := os.Open(nesFile)
	HandleErr(err)
	defer file.Close()
	//name := ReadStr(file, 4)
	prgBankCount := ReadUint8(file)
	chrBankCount := ReadUint8(file)
	mapperID1 := ReadUint8(file)
	mapperID2 := ReadUint8(file)
	//ramSize := ReadUint8(file)  // 还是要跳过的，否则会出问题
	//tv1 := ReadUint8(file)
	//tv2 := ReadUint8(file)
	ReadBytes(file, 5)
	if mapperID1&0x4 > 0 {
		ReadBytes(file, 512)
	}
	mapperID := (mapperID2>>4)<<4 | mapperID1>>4
	prgRom := ReadBytes(file, int(prgBankCount)*16*1024)
	chrRom := ReadBytes(file, int(chrBankCount)*8*1024)
	mapper := createMapper(mapperID, prgBankCount, chrBankCount)
	return &Cartridge{
		PrgRom:       prgRom,
		ChrRom:       chrRom,
		MapperID:     mapperID,
		PrgBankCount: prgBankCount,
		ChrBankCount: chrBankCount,
		Mapper:       mapper,
	}
}

// 应该使用接口
func createMapper(mapperID uint8, prgBankCount, chrBankCount uint8) *Mapper {
	switch mapperID {
	case 0:
		return NewMapper(prgBankCount, chrBankCount)
	default:
		panic(fmt.Sprintf("Can't Handle MapperID %v", mapperID))
	}
}
