/*
@author: sk
@date: 2023/10/26
*/
package nes4

type Ppu struct {
	TableNames         [2][1024]uint8     // 两个屏幕
	TablePatterns      [2][4 * 1024]uint8 // 2个tile
	TablePalettes      [32]uint8          // 调色板
	PalScreen          [16 * 4]*Color     // 总调色板
	Status             *Uint8
	Mask               *Uint8
	Control            *Uint8
	VRamAddr, TRamAddr *Uint16
	FineX              uint8
	AddrLatch          uint8
	PpuDataBuffer      uint8
	// 逐点扫描
	Scanline int16
	Cycle    uint16
	// 背景渲染
	BgNextTileID     uint8
	BgNextTileAttrib uint8
	BgNextTileLsb    uint8
	BgNextTileMsb    uint8
	// 用于偏移的信息,滚动屏幕
	BgShifterPatternLo     uint16
	BgShifterPatternHi     uint16
	BgShifterAttribLo      uint16
	BgShifterAttribHi      uint16
	OAM                    [64]*ObjectAttr // 64 * 4   256 正好一页数据，一次DMA拷贝
	OamAddr                uint8
	SpriteScanline         [8]*ObjectAttr
	SpriteCount            uint8
	SpriteShifterPatternLo [8]uint8
	SpriteShifterPatternHi [8]uint8
	SpriteZeroHitPossible  bool
	SpriteZeroBeingRender  bool
	Cartridge              *Cartridge
	Nmi                    bool
	SprScreen              [256][240]*Color    // 显示屏幕
	SprNameTables          [2][256][240]*Color // 两个屏幕
	SprPatternTables       [2][128][128]*Color // 两个精灵表
}

type ObjectAttr struct {
	Y    uint8
	ID   uint8
	Attr uint8
	X    uint8
}

func (a *ObjectAttr) SetData(index uint8, data uint8) {
	switch index {
	case 0:
		a.Y = data
	case 1:
		a.ID = data
	case 2:
		a.Attr = data
	case 3:
		a.X = data
	}
}

type Color struct {
	R, G, B uint8
}

func NewColor(r uint8, g uint8, b uint8) *Color {
	return &Color{R: r, G: g, B: b}
}

func NewPpu() *Ppu {
	palScreen := [16 * 4]*Color{}
	palScreen[0x00] = NewColor(84, 84, 84)
	palScreen[0x01] = NewColor(0, 30, 116)
	palScreen[0x02] = NewColor(8, 16, 144)
	palScreen[0x03] = NewColor(48, 0, 136)
	palScreen[0x04] = NewColor(68, 0, 100)
	palScreen[0x05] = NewColor(92, 0, 48)
	palScreen[0x06] = NewColor(84, 4, 0)
	palScreen[0x07] = NewColor(60, 24, 0)
	palScreen[0x08] = NewColor(32, 42, 0)
	palScreen[0x09] = NewColor(8, 58, 0)
	palScreen[0x0A] = NewColor(0, 64, 0)
	palScreen[0x0B] = NewColor(0, 60, 0)
	palScreen[0x0C] = NewColor(0, 50, 60)
	palScreen[0x0D] = NewColor(0, 0, 0)
	palScreen[0x0E] = NewColor(0, 0, 0)
	palScreen[0x0F] = NewColor(0, 0, 0)

	palScreen[0x10] = NewColor(152, 150, 152)
	palScreen[0x11] = NewColor(8, 76, 196)
	palScreen[0x12] = NewColor(48, 50, 236)
	palScreen[0x13] = NewColor(92, 30, 228)
	palScreen[0x14] = NewColor(136, 20, 176)
	palScreen[0x15] = NewColor(160, 20, 100)
	palScreen[0x16] = NewColor(152, 34, 32)
	palScreen[0x17] = NewColor(120, 60, 0)
	palScreen[0x18] = NewColor(84, 90, 0)
	palScreen[0x19] = NewColor(40, 114, 0)
	palScreen[0x1A] = NewColor(8, 124, 0)
	palScreen[0x1B] = NewColor(0, 118, 40)
	palScreen[0x1C] = NewColor(0, 102, 120)
	palScreen[0x1D] = NewColor(0, 0, 0)
	palScreen[0x1E] = NewColor(0, 0, 0)
	palScreen[0x1F] = NewColor(0, 0, 0)

	palScreen[0x20] = NewColor(236, 238, 236)
	palScreen[0x21] = NewColor(76, 154, 236)
	palScreen[0x22] = NewColor(120, 124, 236)
	palScreen[0x23] = NewColor(176, 98, 236)
	palScreen[0x24] = NewColor(228, 84, 236)
	palScreen[0x25] = NewColor(236, 88, 180)
	palScreen[0x26] = NewColor(236, 106, 100)
	palScreen[0x27] = NewColor(212, 136, 32)
	palScreen[0x28] = NewColor(160, 170, 0)
	palScreen[0x29] = NewColor(116, 196, 0)
	palScreen[0x2A] = NewColor(76, 208, 32)
	palScreen[0x2B] = NewColor(56, 204, 108)
	palScreen[0x2C] = NewColor(56, 180, 204)
	palScreen[0x2D] = NewColor(60, 60, 60)
	palScreen[0x2E] = NewColor(0, 0, 0)
	palScreen[0x2F] = NewColor(0, 0, 0)

	palScreen[0x30] = NewColor(236, 238, 236)
	palScreen[0x31] = NewColor(168, 204, 236)
	palScreen[0x32] = NewColor(188, 188, 236)
	palScreen[0x33] = NewColor(212, 178, 236)
	palScreen[0x34] = NewColor(236, 174, 236)
	palScreen[0x35] = NewColor(236, 174, 212)
	palScreen[0x36] = NewColor(236, 180, 176)
	palScreen[0x37] = NewColor(228, 196, 144)
	palScreen[0x38] = NewColor(204, 210, 120)
	palScreen[0x39] = NewColor(180, 222, 120)
	palScreen[0x3A] = NewColor(168, 226, 144)
	palScreen[0x3B] = NewColor(152, 226, 180)
	palScreen[0x3C] = NewColor(160, 214, 228)
	palScreen[0x3D] = NewColor(160, 162, 160)
	palScreen[0x3E] = NewColor(0, 0, 0)
	palScreen[0x3F] = NewColor(0, 0, 0)
	return &Ppu{PalScreen: palScreen, Status: NewUint8(), Mask: NewUint8(), Control: NewUint8(),
		VRamAddr: NewUint16(), TRamAddr: NewUint16()}
}

func (p *Ppu) GetScreen() [256][240]*Color {
	return p.SprScreen
}

func (p *Ppu) GetPatternTable(i, palette uint8) [128][128]*Color {
	for tileY := uint16(0); tileY < 16; tileY++ {
		for tileX := uint16(0); tileX < 16; tileX++ {
			offset := (tileY*16 + tileX) * 8 * 2 // 一共过了tileY*16 + tileX 个Tile，每个Tile 8*8 需要 8*2 byte
			for row := uint16(0); row < 8; row++ {
				// 每行8个像素由2byte组成 i 是第几个tile表(一共2个) 一共8行，所以另外一个byte需要偏移8
				tileHi := p.PpuRead(uint16(i)*0x1000 + offset + row + 0x0000)
				tileLo := p.PpuRead(uint16(i)*0x1000 + offset + row + 0x0008)
				for col := uint16(0); col < 8; col++ {
					// 拼接高位与地位获取索引 获取颜色
					index := (tileHi&0x01)<<1 | (tileLo & 0x01)
					tileHi >>= 1
					tileLo >>= 1 // 之所以 7-col是因为 这里是从低位开始计算的，每次位移抹除的也是低位
					p.SprPatternTables[i][tileX*8+(7-col)][tileY*8+row] = p.GetColorFromPalette(palette, index)
				}
			}
		}
	}
	return p.SprPatternTables[i]
}

func (p *Ppu) GetColorFromPalette(palette uint8, index uint8) *Color {
	// 基本偏移  0x3F00    索引0~3正好占用低两位   调色板索引直接位移即可(*4)
	return p.PalScreen[p.PpuRead(0x3F00+(uint16(palette)<<2)+uint16(index))&0x3F]
}

func (p *Ppu) CpuRead(addr uint16) uint8 { // TODO 没有readOnly参数
	switch addr {
	case 0x0000:
		return p.Control.Data
	case 0x0001:
		return p.Mask.Data
	case 0x0002:
		return p.Status.Data
	}
	return 0
}

func (p *Ppu) CpuWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // 写入控制状态
		p.Control.Data = data
		p.TRamAddr.Set(LoopNameTableX, uint16(p.Control.Get(ControlNameTableX)))
		p.TRamAddr.Set(LoopNameTableY, uint16(p.Control.Get(ControlNameTableY)))
	case 0x0001: // 写入Mask
		p.Mask.Data = data
	case 0x0003: // 与下面共同作用写入OAM 不过一般使用 DMA
		p.OamAddr = data
	case 0x0004:
		p.OAM[p.OamAddr/4].SetData(p.OamAddr%4, data)
	case 0x0005: // 滚屏操作
		if p.AddrLatch == 0 { // 第一次写X
			p.FineX = data & 0x07                        // 低3位
			p.TRamAddr.Set(LoopCoarseX, uint16(data>>3)) // 高5位
			p.AddrLatch = 1
		} else { // 第二次写Y
			p.TRamAddr.Set(LoopFineY, uint16(data&0x07))
			p.TRamAddr.Set(LoopCoarseY, uint16(data>>3))
			p.AddrLatch = 0
		}
	case 0x0006:
		if p.AddrLatch == 0 {
			// 占用TRamAddr存储地址 高位只取 6位
			p.TRamAddr.Data = uint16(data&0x3F)<<8 | p.TRamAddr.Data&0x00FF
			p.AddrLatch = 1
		} else { // 写入低位并复制给VRamAddr
			p.TRamAddr.Data = p.TRamAddr.Data&0xFF00 | uint16(data)
			p.VRamAddr.Data = p.TRamAddr.Data
			p.AddrLatch = 0
		}
	case 0x0007:
		p.PpuWrite(p.VRamAddr.Data, data)
		// 水平模式还是垂直模式 用于自增 LoopCoarseX 与 LoopCoarseY
		p.VRamAddr.Data += If(p.Control.Get(ControlIncrementMode) == 1, uint16(32), uint16(1))
	}
}

func (p *Ppu) PpuRead(addr uint16) uint8 {
	addr &= 0x3FFF // 满足镜像内存
	if data, ok := p.Cartridge.PpuRead(addr); ok {
		return data
	} else if addr <= 0x1FFF {
		// 获取指定地址的样式  前面正好是第一或2个图案表  后面是图案表中的位置
		return p.TablePatterns[(addr&0x1000)>>12][addr&0x0FFF]
	} else if addr <= 0x3EFF {
		addr &= 0x0FFF
		if p.Cartridge.Mirror == VERTICAL { // 垂直镜像分别为4个屏幕设置画面
			if addr <= 0x03FF { // 用于操作镜像数据
				return p.TableNames[0][addr&0x03FF]
			} else if addr <= 0x07FF {
				return p.TableNames[1][addr&0x03FF]
			} else if addr <= 0x0BFF {
				return p.TableNames[0][addr&0x03FF]
			} else if addr <= 0x0FFF {
				return p.TableNames[1][addr&0x03FF]
			}
		} else if p.Cartridge.Mirror == HORIZONTAL {
			if addr <= 0x03FF {
				return p.TableNames[0][addr&0x03FF]
			} else if addr <= 0x07FF {
				return p.TableNames[0][addr&0x03FF]
			} else if addr <= 0x0BFF {
				return p.TableNames[1][addr&0x03FF]
			} else if addr <= 0x0FFF {
				return p.TableNames[1][addr&0x03FF]
			}
		}
	} else if addr <= 0x3FFF { // 获取调色板
		addr &= 0x001F
		switch addr { // 一共32个  前4组是背景调色板  后4组是精灵调色板  每组第一个是背景色
		case 0x0010, 0x0014, 0x0018, 0x001C: // 若是精灵背景色取对应背景调色板的背景色
			addr -= 0x0010
		}
		return p.TablePalettes[addr] & If(p.Mask.Get(MaskGrayScale) == 1, uint8(0x30), uint8(0x3F))
	}
	return 0
}

func (p *Ppu) PpuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF
	if p.Cartridge.PpuWrite(addr, data) {

	} else if addr < 0x1FFF {
		p.TablePatterns[(addr&0x1000)>>12][addr&0x0FFF] = data
	} else if addr < 0x3EFF {
		addr &= 0x0FFF
		if p.Cartridge.Mirror == VERTICAL { // 垂直镜像分别为4个屏幕设置画面
			if addr <= 0x03FF { // 0X03FF = 1K -1
				p.TableNames[0][addr&0x03FF] = data
			} else if addr <= 0x07FF {
				p.TableNames[1][addr&0x03FF] = data
			} else if addr <= 0x0BFF {
				p.TableNames[0][addr&0x03FF] = data
			} else if addr <= 0x0FFF {
				p.TableNames[1][addr&0x03FF] = data
			}
		} else if p.Cartridge.Mirror == HORIZONTAL {
			if addr <= 0x03FF {
				p.TableNames[0][addr&0x03FF] = data
			} else if addr <= 0x07FF {
				p.TableNames[0][addr&0x03FF] = data
			} else if addr <= 0x0BFF {
				p.TableNames[1][addr&0x03FF] = data
			} else if addr <= 0x0FFF {
				p.TableNames[1][addr&0x03FF] = data
			}
		}
	} else if addr < 0x3FFF {
		addr &= 0x001F
		switch addr {
		case 0x0010, 0x0014, 0x0018, 0x001C:
			addr -= 0x0010
		}
		p.TablePalettes[addr] = data
	}
}

func (p *Ppu) ConnectCartridge(cartridge *Cartridge) {
	p.Cartridge = cartridge
}

func (p *Ppu) Reset() {
	p.FineX = 0
	p.AddrLatch = 0
	p.PpuDataBuffer = 0
	p.Scanline = 0
	p.Cycle = 0
	p.BgNextTileID = 0
	p.BgNextTileAttrib = 0
	p.BgNextTileLsb = 0
	p.BgNextTileMsb = 0
	p.BgShifterAttribLo = 0
	p.BgShifterAttribHi = 0
	p.BgShifterPatternLo = 0
	p.BgShifterPatternHi = 0
	p.Status.Data = 0
	p.Mask.Data = 0
	p.Control.Data = 0
	p.VRamAddr.Data = 0
	p.TRamAddr.Data = 0
}

func (p *Ppu) Clock() {
	// 逐条渲染 -1是准备阶段 Scanline 纵向，Cycle横向
	if p.Scanline >= -1 && p.Scanline < 240 {
		// 跳过第一次
		if p.Scanline == 0 && p.Cycle == 0 {
			p.Cycle = 1
		}
		// 初始化数据
		if p.Scanline == -1 && p.Cycle == 1 {
			p.Status.Set(StatusVerticalBlank, 0)
			p.Status.Set(StatusSpriteOverflow, 0)
			p.Status.Set(StatusSpriteZeroHit, 0)
			for i := 0; i < 8; i++ {
				p.SpriteShifterPatternLo[i] = 0
				p.SpriteShifterPatternHi[i] = 0
			}
		}
		if (p.Cycle >= 2 && p.Cycle < 258) || (p.Cycle >= 321 && p.Cycle < 338) {
			p.UpdateShifter()
			switch (p.Cycle - 1) % 8 {
			case 0:
				p.LoadBgShifter()
				// p.VRamAddr.Data & 0x0FFF 获取相对地址
				// 偏移到对应的NameTable区域
				p.BgNextTileID = p.PpuRead(0x2000 | (p.VRamAddr.Data & 0x0FFF))
			case 2:
				// << 偏移是为了放到对应的位置  >> 偏移是为了共用属性， 4*4个格子公用一个属性
				p.BgNextTileAttrib = p.PpuRead(0x23C0 | p.VRamAddr.Get(LoopNameTableY)<<11 |
					p.VRamAddr.Get(LoopNameTableX)<<10 |
					p.VRamAddr.Get(LoopCoarseY)>>2<<3 | // 原来5位因为公用格子抹除2位所以偏移3位即可
					p.VRamAddr.Get(LoopCoarseX)>>2)
				// 调色盘使用 4*4 但是属性使用 2*2的共享,还需进一步分割
				if p.VRamAddr.Get(LoopCoarseY)&0x02 > 0 {
					p.BgNextTileAttrib >>= 4
				}
				if p.VRamAddr.Get(LoopCoarseX)&0x02 > 0 {
					p.BgNextTileAttrib >>= 2
				}
				p.BgNextTileAttrib &= 0x03
			case 4:
				p.BgNextTileLsb = p.PpuRead(uint16(p.Control.Get(ControlPatternBackground))<<12 +
					uint16(p.BgNextTileID)<<4 + p.VRamAddr.Get(LoopFineY) + 0)
			case 6:
				p.BgNextTileMsb = p.PpuRead(uint16(p.Control.Get(ControlPatternBackground))<<12 +
					uint16(p.BgNextTileID)<<4 + p.VRamAddr.Get(LoopFineY) + 8)
			case 7:
				p.IncrementScrollX()
			}
		}
		if p.Cycle == 256 {
			p.IncrementScrollY()
		}
		// 横向扫描结束，重新设置X位置
		if p.Cycle == 257 {
			p.LoadBackgroundShifters()
			p.TransferAddressX()
		}
		if p.Cycle == 338 || p.Cycle == 340 {
			p.BgNextTileID = p.PpuRead(0x2000 | (p.VRamAddr.Data & 0x0FFF))
		}
		if p.Scanline == -1 && p.Cycle >= 280 && p.Cycle < 305 {
			p.TransferAddressY()
		}
		// 绘制前景
		if p.Cycle == 257 && p.Scanline >= 0 {
			//p.SpriteScanline
			p.SpriteCount = 0
			for i := 0; i < 8; i++ {
				p.SpriteShifterPatternLo[i] = 0
				p.SpriteShifterPatternHi[i] = 0
			}
			oamCount := 0
			p.SpriteZeroHitPossible = false // 寻找要渲染的8个对象 顺便填写溢出标记
			for oamCount < 64 && p.SpriteCount < 9 {
				diff := p.Scanline - int16(p.OAM[oamCount].Y)
				if diff >= 0 && diff < If(p.Control.Get(ControlSpriteSize) == 1, int16(16), int16(8)) {
					if p.SpriteCount < 8 {
						if oamCount == 0 {
							p.SpriteZeroHitPossible = true
						}
						p.SpriteScanline[p.SpriteCount] = p.OAM[oamCount]
						p.SpriteCount++
					}
				}
				oamCount++
			}
			p.Status.Set(StatusSpriteOverflow, If(p.SpriteCount >= 8, uint8(1), uint8(0)))
		}
		if p.Cycle == 340 {
			for i := uint8(0); i < p.SpriteCount; i++ {

			}
		}
	}
}

func (p *Ppu) UpdateShifter() {

}

func (p *Ppu) LoadBgShifter() {

}

func (p *Ppu) IncrementScrollX() {

}

func (p *Ppu) IncrementScrollY() {

}

func (p *Ppu) LoadBackgroundShifters() {

}

func (p *Ppu) TransferAddressX() {

}

func (p *Ppu) TransferAddressY() {

}
