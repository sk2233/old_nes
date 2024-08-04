/*
@author: sk
@date: 2023/11/11
*/
package main

import (
	"encoding/binary"
	"os"
)

const NESFileMagic = 0x1a53454e

type Cartridge struct {
	PRGRam   []byte
	CHRRam   []byte
	MapperID uint8
	PRGBanks uint8
	CHRBanks uint8
	VMirror  bool
	Mapper   *Mapper
}

func NewCartridge(path string) *Cartridge {
	file, err := os.Open(path)
	HandleErr(err)
	defer file.Close()
	header := &FileHeader{}
	err = binary.Read(file, binary.LittleEndian, header)
	HandleErr(err)

	if header.Magic != NESFileMagic {
		panic("invalid nes file")
	}
	mapperID := (header.Mapper1 >> 4) | (header.Mapper2 >> 4 << 4)
	if header.Mapper1&0x04 > 0 {
		_, err = file.Read(make([]byte, 512)) // 抛弃无用信息
		HandleErr(err)
	}
	prgRam := make([]byte, int(header.PRGBanks)*16*KB)
	_, err = file.Read(prgRam)
	HandleErr(err)
	chrRam := make([]byte, int(header.CHRBanks)*8*KB)
	_, err = file.Read(chrRam)
	HandleErr(err)

	return &Cartridge{
		PRGRam:   prgRam,
		CHRRam:   chrRam,
		MapperID: mapperID,
		PRGBanks: header.PRGBanks,
		CHRBanks: header.CHRBanks,
		Mapper:   NewMapper(header.PRGBanks, header.CHRBanks),
		VMirror:  header.Mapper1&0x01 > 0,
	}
}

func (p *Cartridge) CpuRead(addr uint16) (byte, bool) { // 读取程序代码
	if addr, ok := p.Mapper.CpuReadMap(addr); ok {
		return p.PRGRam[addr], true
	}
	return 0, false
}

func (p *Cartridge) CpuWrite(addr uint16, data byte) bool { // 写入存档
	if addr, ok := p.Mapper.CpuWriteMap(addr); ok {
		p.PRGRam[addr] = data
		return true
	}
	return false
}

func (p *Cartridge) PpuRead(addr uint16) (byte, bool) { // 读取tile
	if addr, ok := p.Mapper.PpuReadMap(addr); ok {
		return p.CHRRam[addr], true
	}
	return 0, false
}

func (p *Cartridge) PpuWrite(addr uint16, data byte) bool { // 写入tile 一般不会支持
	if addr, ok := p.Mapper.PpuWriteMap(addr); ok {
		p.CHRRam[addr] = data // 一般来说这里是不会被写入的
		return true
	}
	return false
}

//===============FileHeader================

type FileHeader struct { // 文件头
	Magic      uint32 // 'NES '
	PRGBanks   uint8
	CHRBanks   uint8
	Mapper1    byte
	Mapper2    byte
	PRGRAMSize uint8
	_          [7]byte // 无用字段
}
