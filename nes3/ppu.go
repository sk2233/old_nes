/*
@author: sk
@date: 2023/11/11
*/
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ppu struct {
	Cartridge     *Cartridge
	TabName       [][]byte
	TabPalette    []byte          // 调色盘
	PalColors     []color.Color   // 全部系统调色盘
	Scanline      int             // 当前处理的是第几行
	Cycle         int             // 某一行的周期
	Screen        *ebiten.Image   // 显示屏幕
	Tiles         []*ebiten.Image // tile图片
	FrameComplete bool
	// 几个主要的 CPU寄存器
	Ctrl  *PpuCtrl  // 0x2000
	Mask  *PpuMask  // 0x2001
	State *PpuState // 0x2002
	// 通过地址数据传送
	FirstAddr bool // 地址需要两次传输，用来区分第一次与第二次
	LastData  byte // 读数据会延迟1需要保留上次的结果
	Nmi       bool // 用来触发CPU进入中断，由bus处理
	TRamAddr  *PpuLoopy
	VRamAddr  *PpuLoopy
	FineX     uint8 // 3bit
	// 用来渲染下一个Tile进行准备
	NextBgTileID   uint8
	NextBgTileAttr uint8
	NextBgTileHi   uint8
	NextBgTileLo   uint8
	// 当前在渲染的Tile信息 下面的注释是指在加载新值时的情况
	BgShifterTileAttr uint32 // 前16b 2个一组是当前的8个调色板编号，后16b是NextTile的8个的，实际没8个都是同一个调色板，这样做只是为了方便位移操作
	BgShifterTileHi   uint16 // 前8b 当前tile 8个像素的高位 后8b是NextTile 8个像素的高位
	BgShifterTileLo   uint16 // 与BgTileHi类似不过是低位
}

func NewPpu() *Ppu {
	tabName := make([][]byte, 2)
	tabName[0] = make([]byte, KB)
	tabName[1] = make([]byte, KB)

	palColors := make([]color.Color, 16*4)
	palColors[0x00] = NewColor(84, 84, 84)
	palColors[0x01] = NewColor(0, 30, 116)
	palColors[0x02] = NewColor(8, 16, 144)
	palColors[0x03] = NewColor(48, 0, 136)
	palColors[0x04] = NewColor(68, 0, 100)
	palColors[0x05] = NewColor(92, 0, 48)
	palColors[0x06] = NewColor(84, 4, 0)
	palColors[0x07] = NewColor(60, 24, 0)
	palColors[0x08] = NewColor(32, 42, 0)
	palColors[0x09] = NewColor(8, 58, 0)
	palColors[0x0A] = NewColor(0, 64, 0)
	palColors[0x0B] = NewColor(0, 60, 0)
	palColors[0x0C] = NewColor(0, 50, 60)
	palColors[0x0D] = NewColor(0, 0, 0)
	palColors[0x0E] = NewColor(0, 0, 0)
	palColors[0x0F] = NewColor(0, 0, 0)

	palColors[0x10] = NewColor(152, 150, 152)
	palColors[0x11] = NewColor(8, 76, 196)
	palColors[0x12] = NewColor(48, 50, 236)
	palColors[0x13] = NewColor(92, 30, 228)
	palColors[0x14] = NewColor(136, 20, 176)
	palColors[0x15] = NewColor(160, 20, 100)
	palColors[0x16] = NewColor(152, 34, 32)
	palColors[0x17] = NewColor(120, 60, 0)
	palColors[0x18] = NewColor(84, 90, 0)
	palColors[0x19] = NewColor(40, 114, 0)
	palColors[0x1A] = NewColor(8, 124, 0)
	palColors[0x1B] = NewColor(0, 118, 40)
	palColors[0x1C] = NewColor(0, 102, 120)
	palColors[0x1D] = NewColor(0, 0, 0)
	palColors[0x1E] = NewColor(0, 0, 0)
	palColors[0x1F] = NewColor(0, 0, 0)

	palColors[0x20] = NewColor(236, 238, 236)
	palColors[0x21] = NewColor(76, 154, 236)
	palColors[0x22] = NewColor(120, 124, 236)
	palColors[0x23] = NewColor(176, 98, 236)
	palColors[0x24] = NewColor(228, 84, 236)
	palColors[0x25] = NewColor(236, 88, 180)
	palColors[0x26] = NewColor(236, 106, 100)
	palColors[0x27] = NewColor(212, 136, 32)
	palColors[0x28] = NewColor(160, 170, 0)
	palColors[0x29] = NewColor(116, 196, 0)
	palColors[0x2A] = NewColor(76, 208, 32)
	palColors[0x2B] = NewColor(56, 204, 108)
	palColors[0x2C] = NewColor(56, 180, 204)
	palColors[0x2D] = NewColor(60, 60, 60)
	palColors[0x2E] = NewColor(0, 0, 0)
	palColors[0x2F] = NewColor(0, 0, 0)

	palColors[0x30] = NewColor(236, 238, 236)
	palColors[0x31] = NewColor(168, 204, 236)
	palColors[0x32] = NewColor(188, 188, 236)
	palColors[0x33] = NewColor(212, 178, 236)
	palColors[0x34] = NewColor(236, 174, 236)
	palColors[0x35] = NewColor(236, 174, 212)
	palColors[0x36] = NewColor(236, 180, 176)
	palColors[0x37] = NewColor(228, 196, 144)
	palColors[0x38] = NewColor(204, 210, 120)
	palColors[0x39] = NewColor(180, 222, 120)
	palColors[0x3A] = NewColor(168, 226, 144)
	palColors[0x3B] = NewColor(152, 226, 180)
	palColors[0x3C] = NewColor(160, 214, 228)
	palColors[0x3D] = NewColor(160, 162, 160)
	palColors[0x3E] = NewColor(0, 0, 0)
	palColors[0x3F] = NewColor(0, 0, 0)

	tiles := make([]*ebiten.Image, 2)
	tiles[0] = ebiten.NewImage(128, 128)
	tiles[1] = ebiten.NewImage(128, 128)
	return &Ppu{TabName: tabName, TabPalette: make([]byte, 32), FirstAddr: true, Tiles: tiles,
		Screen: ebiten.NewImage(256, 240), PalColors: palColors, Ctrl: &PpuCtrl{},
		Mask: &PpuMask{}, State: &PpuState{}, VRamAddr: &PpuLoopy{}, TRamAddr: &PpuLoopy{}}
}

func NewColor(r uint8, g uint8, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}
}

func (p *Ppu) ConnectCartridge(cartridge *Cartridge) {
	p.Cartridge = cartridge
}

func (p *Ppu) Clock() {
	// 按照那张流程图进行，特定时间区间干特定的事
	if p.Scanline >= -1 && p.Scanline < 240 {
		if p.Scanline == 0 && p.Cycle == 0 {
			p.Cycle = 1
		}
		if p.Scanline == -1 && p.Cycle == 1 { // 重新开始渲染
			p.State.VBlank = false
		}
		if (p.Cycle > 1 && p.Cycle < 258) || (p.Cycle > 320 && p.Cycle < 338) {
			p.UpdateShifter()    // 该绘制横向的下一个像素了，偏移一位
			switch p.Cycle % 8 { // 一个Tile横向8个像素
			case 1:
				p.LoadShifter() // 已经偏移8次，原来的高位全部偏移掉了，原来的高位变成了现在的高位，补充低位
				// 取的是 NameTabY NameTabX CoarseY CoarseX 的组合，取NameTab的相对位置实现滚屏操作
				p.NextBgTileID = p.PpuRead(0x2000 | p.VRamAddr.Get()&0x0FFF)
			case 3:
				// 0x23C0 是属性字段的基地址，CoarseY与CoarseX >>2 是因为4*4个格子使用一个byte,先取出这个byte
				p.NextBgTileAttr = p.PpuRead(0x23C0 | uint16(p.VRamAddr.NameTab)<<10 |
					uint16(p.VRamAddr.CoarseY>>2)<<3 | uint16(p.VRamAddr.CoarseX)>>2)
				// 每2*2个格子公用2b，作为调色板选择
				if p.VRamAddr.CoarseY&0x02 > 0 { // 偏移到正确的2b后只取2b使用
					p.NextBgTileAttr >>= 4
				}
				if p.VRamAddr.CoarseX&0x02 > 0 {
					p.NextBgTileAttr >>= 2
				}
				p.NextBgTileAttr &= 0x03
			case 5:
				// FineY 是向下偏移的行数每偏移1行 +1 每8个像素由2byte高低位组成
				// NextBgTileID 每个TileID 需要 2*8 个byte渲染,所以偏移 4
				// BackgroundTile 每个tile集有 16*16个tile，所以偏移 4+8=12
				p.NextBgTileLo = p.PpuRead(uint16(p.Ctrl.BackgroundTile)<<12 + uint16(p.NextBgTileID)<<4 +
					uint16(p.VRamAddr.FineY))
			case 7:
				// 与上面类似，不过这个取高位，要加8
				p.NextBgTileHi = p.PpuRead(uint16(p.Ctrl.BackgroundTile)<<12 + uint16(p.NextBgTileID)<<4 +
					uint16(p.VRamAddr.FineY) + 8)
			case 0:
				p.IncrementScrollX() // 渲染完一行的一个tile,
			}
		}
		if p.Cycle == 256 { // 渲染完一行
			p.IncrementScrollY()
		}
		if p.Cycle == 257 { // 重置X的位置到开始绘制这帧画面的时候
			p.LoadShifter()
			p.TransferAddressX()
		}
		if p.Cycle == 338 || p.Cycle == 340 {
			p.NextBgTileID = p.PpuRead(0x2000 | p.VRamAddr.Get()&0x0FFF)
		}
		if p.Scanline == -1 && p.Cycle >= 280 && p.Cycle < 305 { // TODO 为什么重置多次？
			p.TransferAddressY() // 重置Y具体到像素位置
		}
	}

	if p.Scanline == 241 && p.Cycle == 1 {
		p.State.VBlank = true // 开始渲染屏幕外,通知cpu可以填充数据了
		if p.Ctrl.EnableNmi {
			p.Nmi = true
		}
	}

	if p.Mask.ShowBackground {
		// 需要应用 FineX 的基础偏移
		lo := (p.BgShifterTileLo >> (15 - p.FineX)) & 0x01
		hi := (p.BgShifterTileHi >> (15 - p.FineX)) & 0x01
		index := uint8(lo | (hi << 1))
		// 也需要应用 FineX 的基础偏移，且要注意这个是2bit
		pal := uint8(p.BgShifterTileAttr >> ((15 - p.FineX) * 2) & 0x03)
		p.Screen.Set(p.Cycle-1, p.Scanline, p.GetColorFromPalette(pal, index)) // 256 * 240
	}

	p.Cycle++
	if p.Cycle > 340 { // 扫描完一行了
		p.Cycle = 0
		p.Scanline++
		if p.Scanline > 260 { // 扫描完一屏幕了
			p.Scanline = -1
			p.FrameComplete = true
		}
	}
}

func (p *Ppu) CpuRead(addr uint16) byte { // cpu与ppu的交互通过这8个地址实现
	switch addr {
	case 0x0000: // Control
		return p.Ctrl.Get()
	case 0x0001: // Mask
		return p.Mask.Get()
	case 0x0002: // Status
		res := p.State.Get()   // 实际只有高3位有用
		p.State.VBlank = false // 每次获取都是渲染一屏幕了吗?
		p.FirstAddr = true     // 进行复位
		return res
	case 0x0003: // OAM Addr
	case 0x0004: // OAM Data
	case 0x0005: // Scroll
	case 0x0006: // Addr
	case 0x0007: // Data
		res := p.LastData // 真实的电子结构决定这里需要延迟 代码也需要模拟
		p.LastData = p.PpuRead(p.VRamAddr.Get())
		if p.VRamAddr.Get() >= 0x3F00 { // 读取调色板数据的读取不需要这个延迟
			res = p.LastData
		}
		if p.Ctrl.Increment32 { // 保证连续读取，一般都是设置一个地址，然后连续读取 注意垂直与水平连续
			p.VRamAddr.CoarseY++
		} else {
			p.VRamAddr.CoarseX++
		}
		return res
	}
	return 0
}

func (p *Ppu) CpuWrite(addr uint16, data byte) {
	switch addr {
	case 0x0000: // Control
		p.Ctrl.Set(data)
		p.TRamAddr.NameTab = p.Ctrl.NameTab
	case 0x0001: // Mask
		p.Mask.Set(data)
	case 0x0002: // Status
	case 0x0003: // OAM Addr
	case 0x0004: // OAM Data
	case 0x0005: // Scroll 写入滚动信息
		if p.FirstAddr {
			p.FineX = data & 0x07
			p.TRamAddr.CoarseX = data >> 3
			p.FirstAddr = false
		} else {
			p.TRamAddr.FineY = data & 0x07
			p.TRamAddr.CoarseY = data >> 3
			p.FirstAddr = true
		}
	case 0x0006: // Addr
		if p.FirstAddr { // 先写高位，再写低位  这里注意
			p.TRamAddr.Set((p.TRamAddr.Get() & 0x00FF) | (uint16(data&0x3F) << 8))
			p.FirstAddr = false
		} else {
			p.TRamAddr.Set((p.TRamAddr.Get() & 0xFF00) | uint16(data))
			p.FirstAddr = true
			p.VRamAddr = p.TRamAddr
		}
	case 0x0007: // Data
		p.PpuWrite(p.VRamAddr.Get(), data)
		if p.Ctrl.Increment32 { // 一般都是确定一个地址后连续写入，不过要注意是水平连续，还是垂直连续
			p.VRamAddr.CoarseY++ // 直接写CoarseY 也相当于+32，CoarseY在高位
		} else {
			p.VRamAddr.CoarseX++
		}
	}
}

func (p *Ppu) GetTile(i, palette uint8) *ebiten.Image {
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			offset := (tileY*16 + tileX) * 8 * 2 // 一共过了tileY*16 + tileX 个Tile，每个Tile 8*8 需要 8*2 byte
			for row := 0; row < 8; row++ {
				// 每行8个像素由2byte组成 i 是第几个tile表(一共2个) 一共8行，所以另外一个byte需要偏移8
				tileHi := p.PpuRead(uint16(int(i)*0x1000 + offset + row + 0x0000))
				tileLo := p.PpuRead(uint16(int(i)*0x1000 + offset + row + 0x0008))
				for col := 0; col < 8; col++ {
					// 拼接高位与地位获取索引 获取颜色
					index := (tileHi&0x01)<<1 | (tileLo & 0x01)
					tileHi >>= 1
					tileLo >>= 1
					clr := p.GetColorFromPalette(palette, index)
					// 之所以 7-col是因为 这里是从低位开始计算的，每次位移抹除的也是低位
					p.Tiles[i].Set(tileX*8+(7-col), tileY*8+row, clr)
				}
			}
		}
	}
	return p.Tiles[i]
}

func (p *Ppu) PpuRead(addr uint16) byte {
	addr &= 0x3FFF
	if data, ok := p.Cartridge.PpuRead(addr); ok { // 0 ~ 0x2000
		return data
	} else if addr < 0x3F00 { // nameTable
		addr &= 0x0FFF
		// 0~1KB 1~2KB
		// 2~3KB 3~4KB
		if p.Cartridge.VMirror { // 垂直镜像 上下镜像
			if addr < KB {
				return p.TabName[0][addr&0x03FF]
			} else if addr < 2*KB {
				return p.TabName[1][addr&0x03FF]
			} else if addr < 3*KB {
				return p.TabName[0][addr&0x03FF]
			} else if addr < 4*KB {
				return p.TabName[1][addr&0x03FF]
			}
		} else {
			if addr < KB {
				return p.TabName[0][addr&0x03FF]
			} else if addr < 2*KB {
				return p.TabName[0][addr&0x03FF]
			} else if addr < 3*KB {
				return p.TabName[1][addr&0x03FF]
			} else if addr < 4*KB {
				return p.TabName[1][addr&0x03FF]
			}
		}
	} else if addr < 0x4000 { // 操作调色板
		addr &= 0x1F // 只有32个 8 * 4
		switch addr {
		// 前4个是背景调色板，后4个是精灵调色板，一一对应，精灵调色板的背景色使用背景调色板的背景色实现透明效果
		case 0x10, 0x14, 0x18, 0x1C:
			addr -= 0x10
		}
		return p.TabPalette[addr]
	}
	return 0
}

func (p *Ppu) PpuWrite(addr uint16, data byte) {
	addr &= 0x3FFF
	if p.Cartridge.PpuWrite(addr, data) {
		// 卡带拦截处理
		return
	} else if addr < 0x2000 { // 样式表 因为不允许写入 这里要拦截一下防止混入nameTable的操作中

	} else if addr < 0x3F00 { // nameTable
		addr &= 0x0FFF
		// 0~1KB 1~2KB
		// 2~3KB 3~4KB
		if p.Cartridge.VMirror { // 垂直镜像 上下镜像
			if addr < KB {
				p.TabName[0][addr&0x03FF] = data
			} else if addr < 2*KB {
				p.TabName[1][addr&0x03FF] = data
			} else if addr < 3*KB {
				p.TabName[0][addr&0x03FF] = data
			} else if addr < 4*KB {
				p.TabName[1][addr&0x03FF] = data
			}
		} else {
			if addr < KB {
				p.TabName[0][addr&0x03FF] = data
			} else if addr < 2*KB {
				p.TabName[0][addr&0x03FF] = data
			} else if addr < 3*KB {
				p.TabName[1][addr&0x03FF] = data
			} else if addr < 4*KB {
				p.TabName[1][addr&0x03FF] = data
			}
		}
	} else if addr < 0x4000 { // 调色板
		addr &= 0x1F  // 只有32个 8 * 4
		switch addr { // 也可以写入与读取同步都不做这个特化，只需要镜像 32 即可
		// 前4个是背景调色板，后4个是精灵调色板，一一对应，精灵调色板的背景色使用背景调色板的背景色实现透明效果
		case 0x10, 0x14, 0x18, 0x1C:
			addr -= 0x10
		}
		p.TabPalette[addr] = data
	}
}

func (p *Ppu) GetColorFromPalette(palette uint8, index uint8) color.Color {
	// 0x3F00 是调色板基地址 每个调色板4个颜色，获取全局调色板的索引
	index = p.PpuRead(0x3F00 + uint16(palette)<<2 + uint16(index))
	return p.PalColors[index&0x3F]
}

func (p *Ppu) IncrementScrollX() {
	if p.Mask.ShowBackground || p.Mask.ShowSprite {
		p.VRamAddr.CoarseX++
		if p.VRamAddr.CoarseX == 32 { // 到32了
			p.VRamAddr.CoarseX = 0
			p.VRamAddr.NameTab ^= 0x01 // 切换到新的NameTab
		}
	}
}

func (p *Ppu) IncrementScrollY() {
	if p.Mask.ShowBackground || p.Mask.ShowSprite {
		p.VRamAddr.FineY++
		if p.VRamAddr.FineY == 8 { // 到8了,下一格
			p.VRamAddr.FineY = 0
			p.VRamAddr.CoarseY++
			if p.VRamAddr.CoarseY == 30 { // 到30了，下一个NameTab,竖向只有 30*8=240个
				p.VRamAddr.CoarseY = 0
				p.VRamAddr.NameTab ^= 0x02
			}
		}
	}
}

func (p *Ppu) TransferAddressX() {
	if p.Mask.ShowBackground || p.Mask.ShowSprite { // 恢复之前X轴方向
		p.VRamAddr.NameTab = (p.VRamAddr.NameTab & 0x02) | (p.TRamAddr.NameTab & 0x01)
		p.VRamAddr.CoarseX = p.TRamAddr.CoarseX
	}
}

func (p *Ppu) TransferAddressY() {
	if p.Mask.ShowBackground || p.Mask.ShowSprite { // 恢复之前的Y轴具体到像素
		p.VRamAddr.FineY = p.TRamAddr.FineY
		p.VRamAddr.NameTab = (p.VRamAddr.NameTab & 0x01) | (p.TRamAddr.NameTab & 0x02)
		p.VRamAddr.CoarseY = p.TRamAddr.CoarseY
	}
}

func (p *Ppu) UpdateShifter() {
	if p.Mask.ShowBackground {
		p.BgShifterTileLo <<= 1
		p.BgShifterTileHi <<= 1
		p.BgShifterTileAttr <<= 2 // 2位一组
	}
}

func (p *Ppu) LoadShifter() {
	p.BgShifterTileLo = (p.BgShifterTileLo & 0xFF00) | uint16(p.NextBgTileLo)
	p.BgShifterTileHi = (p.BgShifterTileHi & 0xFF00) | uint16(p.NextBgTileHi)
	p.BgShifterTileAttr = p.BgShifterTileAttr & 0xFFFF0000
	for i := 0; i < 8; i++ { // 8个像素是一个Tile选择的调色板都一样，方便统一位移操作，复制8次
		p.BgShifterTileAttr |= (uint32(p.NextBgTileAttr) & 0x03) << (i * 2)
	}
}

//==============PpuState================

type PpuState struct { // 8bit
	Unused   uint8 // 5bit
	Overflow bool  // 精灵渲染溢出，默认一行最多8个精灵
	ZeroHit  bool  // 精灵非背景色与背景非背景色第一次重叠标记
	VBlank   bool  // 垂直渲染是否超出屏幕
}

func (s *PpuState) Get() byte {
	res := s.Unused & 0x1F
	if s.Overflow {
		res |= 0x20
	}
	if s.ZeroHit {
		res |= 0x40
	}
	if s.VBlank {
		res |= 0x80
	}
	return res
}

//================PpuMask===============

type PpuMask struct { // 8bit
	GrayScale          bool
	ShowLeftBackground bool
	ShowLeftSprite     bool
	ShowBackground     bool
	ShowSprite         bool
	RedTint            bool
	GreenTint          bool
	BlueTint           bool
}

func (m *PpuMask) Set(data byte) {
	m.GrayScale = (data & 0x01) > 0
	m.ShowLeftBackground = (data & 0x02) > 0
	m.ShowLeftSprite = (data & 0x04) > 0
	m.ShowBackground = (data & 0x08) > 0
	m.ShowSprite = (data & 0x10) > 0
	m.RedTint = (data & 0x20) > 0
	m.GreenTint = (data & 0x40) > 0
	m.BlueTint = (data & 0x80) > 0
}

func (m *PpuMask) Get() byte {
	res := uint8(0)
	if m.GrayScale {
		res |= 0x01
	}
	if m.ShowLeftBackground {
		res |= 0x02
	}
	if m.ShowLeftSprite {
		res |= 0x04
	}
	if m.ShowBackground {
		res |= 0x08
	}
	if m.ShowSprite {
		res |= 0x10
	}
	if m.RedTint {
		res |= 0x20
	}
	if m.GreenTint {
		res |= 0x40
	}
	if m.BlueTint {
		res |= 0x80
	}
	return res
}

//==============PpuCtrl================

type PpuCtrl struct { // 8bit
	NameTab        uint8 // 2b 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	Increment32    bool  // false 1, true 32
	SpriteTile     uint8 // 1b
	BackgroundTile uint8 // 1b
	SpriteSize     bool  // false 8*8 , true 8*16
	Unused         bool
	EnableNmi      bool
}

func (c *PpuCtrl) Set(data byte) {
	c.NameTab = data & 0x03
	c.Increment32 = (data & 0x04) > 0
	c.SpriteTile = (data >> 3) & 0x01
	c.BackgroundTile = (data >> 4) & 0x01
	c.SpriteSize = (data & 0x20) > 0
	c.Unused = (data & 0x40) > 0
	c.EnableNmi = (data & 0x80) > 0
}

func (c *PpuCtrl) Get() byte {
	res := c.NameTab & 0x03
	if c.Increment32 {
		res |= 0x04
	}
	res |= (c.SpriteTile & 0x01) << 3
	res |= (c.BackgroundTile & 0x01) << 4
	if c.SpriteSize {
		res |= 0x20
	}
	if c.Unused {
		res |= 0x40
	}
	if c.EnableNmi {
		res |= 0x80
	}
	return res
}

//===============PpuLoopy=================

type PpuLoopy struct { // 16bit
	CoarseX uint8 // 5
	CoarseY uint8 // 5
	NameTab uint8 // 2b 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	FineY   uint8 // 3
	Unused  bool  // 1
}

func (l *PpuLoopy) Get() uint16 {
	res := uint16(l.CoarseX & 0x1F)
	res |= uint16(l.CoarseY&0x1F) << 5
	res |= uint16(l.NameTab&0x03) << 10
	res |= uint16(l.FineY&0x07) << 12
	if l.Unused {
		res |= 0x80
	}
	return res
}

func (l *PpuLoopy) Set(data uint16) {
	l.CoarseX = uint8(data) & 0x1F
	l.CoarseY = uint8(data>>5) & 0x1F
	l.NameTab = uint8(data>>10) & 0x03
	l.FineY = uint8(data>>12) & 0x07
	l.Unused = (data >> 15) > 0
}
