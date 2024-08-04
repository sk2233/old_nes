/*
@author: sk
@date: 2023/10/14
*/
package main

type Bus struct {
	// 0x0000 zeroPage
	// 0x0100 stack(下压栈)
	// 0x0200 ram 主要cpu操作的内存
	// 0x0800 0x0000~0x07ff mirrors
	// 0x2000 ppuIO寄存器 cpu与ppu交互使用
	// PPUCTRL(0x2000)   0~1位:使用哪个nameTable  2位:0偏移1水平移动,1偏移32垂直移动  3位:精灵使用哪个tileSet(0:0x0000,1:0x1000)
	// 4位:背景使用哪个tileSet  5位:精灵大小(0:8*8,1:8*16)  7位:V_Blank开始时产生中断写数据到ppuRam
	// 0x2008 0x2000~0x2007 mirrors
	// 0x4000 otherIO寄存器(例如手柄)
	// 0x4020 exRom 卡带扩展内存
	// 0x6000 wRam 存档使用
	// 0x8000 prgRomLoBank 与下面的一样都用于存储代码
	// 0xC000 prgRomHiBank
	// 0xFFFA reset,irq,nmi触发后跳转的地址
	cpuRam [64 * K]byte
	// 0x0000 patternTab1 tileSet 与背景一致都是先存储低位再存储高位
	// 0x1000 patternTab2 tileSet
	// 0x2000 4组nameTab与attrs，每组1k，其中nameTab 960b(32P*30P,16*16个tile一个tile需要uin8存储)
	// attrs 64b 每16(4*4)个tile使用1b(uint8)控制，其中每4个使用2位控制选择哪个调色板
	// 0x3000 0x2000~0x2EFF mirrors 这里连前面的一次都没有镜像，纯粹填空间
	// 0x3F00 imgPalette 背景调色板 索引为0的颜色都是一样的背景色
	// 0x3F10 sprPalette 精灵调色板 索引为0的颜色都是一样的透明色
	// 0x3F20 0x3F00~0x3F1F mirrors
	// 0x4000 0x0000~0x3FFF mirrors
	ppuRam [64 * K]byte
}

func (b *Bus) PpuRead(addr uint16) byte {
	addr = b.mirrorPpuAddr(addr)
	return b.ppuRam[addr]
}

func (b *Bus) mirrorPpuAddr(addr uint16) uint16 {
	if addr >= 0x4000 { // 可能存在多重镜像先处理大的
		addr = addr & 0x3FFF
	}
	if addr >= 0x3F20 && addr < 0x4000 {
		addr = 0x3F00 + ((addr - 0x3F00) & 0x01F)
	}
	if addr >= 0x3000 && addr < 0x3F00 {
		addr = addr - 0x1000
	}
	return addr
}

func (b *Bus) CpuRead(addr uint16) byte {
	addr = b.mirrorCpuAddr(addr)
	return b.cpuRam[addr]
}

func (b *Bus) mirrorCpuAddr(addr uint16) uint16 { // 处理镜像地址
	if addr >= 0x2008 && addr < 0x4000 { // 可能存在多重镜像先处理大的
		addr = 0x2000 + ((addr - 0x2000) & 0x0007)
	}
	if addr >= 0x0800 && addr < 0x2000 {
		addr = addr & 0x7FF
	}
	return addr
}

func NewBus() *Bus {
	return &Bus{}
}
