/*
@author: sk
@date: 2023/10/12
*/
package main

type Ppu struct {
	Cartridge     *Cartridge
	NameTab       [2][1024]uint8 // 屏幕中每个位置的tile
	PaletteTab    [32]uint8      // 调色盘
	Cycle         uint16
	ScanLine      int16
	FrameComplete bool
	// 其他字段
	dataBuffer  uint8
	dataAddress uint16
	status      uint8
	statusV     bool
	addrLatch   bool
}

func (p *Ppu) CpuRead(addr uint16) uint8 {
	var res uint8
	switch addr {
	case 0x0000: // Control
	case 0x0001: // Mask
	case 0x0002: // Status
		res = p.status // 除了读取状态，还会重置状态
		p.statusV = false
		p.addrLatch = false
	case 0x0003: // OAM Address
	case 0x0004: // OAM Data
	case 0x0005: // Scroll
	case 0x0006: // PPU Address

	case 0x0007: // PPU Data
		res = p.dataBuffer // 有一定的滞后
		p.dataBuffer = p.PpuRead(p.dataAddress)
		if p.dataAddress > 0x3f00 { // 这种情况下是没有滞后的
			res = p.dataBuffer
		}
	}
	return res
}

func (p *Ppu) CpuWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
	case 0x0001: // Mask
	case 0x0002: // Status
	case 0x0003: // OAM Address
	case 0x0004: // OAM Data
	case 0x0005: // Scroll
	case 0x0006: // PPU Address
		if p.addrLatch { // 是分两次设置值的
			p.dataAddress |= uint16(data)
		} else {
			p.dataAddress = uint16(data) << 8
		}
		p.addrLatch = !p.addrLatch
	case 0x0007: // PPU Data
	}
}

func (p *Ppu) SetCartridge(cartridge *Cartridge) {
	p.Cartridge = cartridge
}

func (p *Ppu) Clock() { // 绘制行为
	p.Cycle++
	if p.Cycle > 340 {
		p.Cycle = 0
		p.ScanLine++
		if p.ScanLine > 260 {
			p.ScanLine = -1
			p.FrameComplete = true
		}
	}
}

type Tile [8][8]uint8

func (p *Ppu) GetPatternTable(i uint8) { // 0 or 1  精灵 or  BG
	tiles := [16][16]Tile{}
	// 由 16 * 16 个Tile组成
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			offset := (tileY*16 + tileX) * 8 * 2
			tile := [8][8]uint8{}
			// 每个Tile 8 * 8 像素
			for y := 0; y < 8; y++ {
				// 正好 一个tile 图像占用  0x1000 = 16 * 16 * 8 * 2   两个int8 16b 每2b一个像素索引
				tile0 := p.PpuRead(uint16(i)*0x1000 + uint16(y+offset))
				tile1 := p.PpuRead(uint16(i)*0x1000 + uint16(y+offset) + 8) // 相邻存储
				for x := 0; x < 8; x++ {
					index := tile0&0x01 + tile1&0x01 // 得到当前位置的颜色索引  0 ~ 3
					tile0 >>= 1                      // 移动到下一个位置
					tile1 >>= 1
					tile[x][y] = index
				}
			}
			tiles[tileX][tileY] = tile
		}
	}
}

// 获取指定调色盘的指定索引颜色
func (p *Ppu) GetColor(palette, index uint8) uint8 {
	// 获取对应的颜色索引
	index = p.PpuRead(0x3f00 + uint16(palette*4+index))
	return index
}

func (p *Ppu) PpuRead(addr uint16) uint8 {
	addr &= 0x3fff // 限制读取的范围
	if data, ok := p.Cartridge.PpuRead(addr); ok {
		return data
	}
	return 0
}
