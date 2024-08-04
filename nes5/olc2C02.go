/*
@author: sk
@date: 2024/2/25
*/
package nes5

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Status struct {
	Unused         uint8 // 5
	SpriteOverflow bool
	SpriteZeroHit  bool
	VerticalBlank  bool
}

func (s *Status) Get() uint8 {
	res := uint8(0)
	res |= SetRangeBit8(s.Unused, 0, 5)
	res |= SetBit[uint8](s.SpriteOverflow, 5)
	res |= SetBit[uint8](s.SpriteZeroHit, 6)
	res |= SetBit[uint8](s.VerticalBlank, 7)
	return res
}

func (s *Status) Set(data uint8) {
	s.Unused = GetRangeBit8(data, 0, 5)
	s.SpriteOverflow = GetBit[uint8](data, 5)
	s.SpriteZeroHit = GetBit[uint8](data, 6)
	s.VerticalBlank = GetBit[uint8](data, 7)
}

type Mask struct { // 都是从低位到高位
	GrayScale            bool
	RenderBackgroundLeft bool
	RenderSpriteLeft     bool
	RenderBackground     bool
	RenderSprite         bool
	EnhanceRed           bool
	EnhanceGreen         bool
	EnhanceBlue          bool
}

func (m *Mask) Set(data uint8) {
	m.GrayScale = GetBit(data, 0)
	m.RenderBackgroundLeft = GetBit(data, 1)
	m.RenderSpriteLeft = GetBit(data, 2)
	m.RenderBackground = GetBit(data, 3)
	m.RenderSprite = GetBit(data, 4)
	m.EnhanceRed = GetBit(data, 5)
	m.EnhanceGreen = GetBit(data, 6)
	m.EnhanceBlue = GetBit(data, 7)
}

func (m *Mask) Get() uint8 {
	res := uint8(0)
	res |= SetBit[uint8](m.GrayScale, 0)
	res |= SetBit[uint8](m.RenderBackgroundLeft, 1)
	res |= SetBit[uint8](m.RenderSpriteLeft, 2)
	res |= SetBit[uint8](m.RenderBackground, 3)
	res |= SetBit[uint8](m.RenderSprite, 4)
	res |= SetBit[uint8](m.EnhanceRed, 5)
	res |= SetBit[uint8](m.EnhanceGreen, 6)
	res |= SetBit[uint8](m.EnhanceBlue, 7)
	return res
}

type Ctrl struct {
	NameTableX        bool
	NameTableY        bool
	IncrementMode     bool
	PatternSprite     bool
	PatternBackground bool
	SpriteSize        bool
	Unused            bool
	EnableNmi         bool
}

func (c *Ctrl) Set(data uint8) {
	c.NameTableX = GetBit(data, 0)
	c.NameTableY = GetBit(data, 1)
	c.IncrementMode = GetBit(data, 2)
	c.PatternSprite = GetBit(data, 3)
	c.PatternBackground = GetBit(data, 4)
	c.SpriteSize = GetBit(data, 5)
	c.Unused = GetBit(data, 6)
	c.EnableNmi = GetBit(data, 7)
}

func (c *Ctrl) Get() uint8 {
	res := uint8(0)
	res |= SetBit[uint8](c.NameTableX, 0)
	res |= SetBit[uint8](c.NameTableY, 1)
	res |= SetBit[uint8](c.IncrementMode, 2)
	res |= SetBit[uint8](c.PatternSprite, 3)
	res |= SetBit[uint8](c.PatternBackground, 4)
	res |= SetBit[uint8](c.SpriteSize, 5)
	res |= SetBit[uint8](c.Unused, 6)
	res |= SetBit[uint8](c.EnableNmi, 7)
	return res
}

type LoopReg struct {
	CoarseX    uint16 // 5
	CoarseY    uint16 // 5
	NameTableX bool
	NameTableY bool
	FineY      uint16 // 3
	Unused     bool
}

func (r *LoopReg) Get() uint16 {
	res := uint16(0)
	res |= SetRangeBit16(r.CoarseX, 0, 5)
	res |= SetRangeBit16(r.CoarseY, 5, 10)
	res |= SetBit[uint16](r.NameTableX, 10)
	res |= SetBit[uint16](r.NameTableY, 11)
	res |= SetRangeBit16(r.FineY, 12, 15)
	res |= SetBit[uint16](r.Unused, 15)
	return res
}

func (r *LoopReg) Set(data uint16) {
	r.CoarseX = GetRangeBit16(data, 0, 5)
	r.CoarseY = GetRangeBit16(data, 5, 10)
	r.NameTableX = GetBit(data, 10)
	r.NameTableY = GetBit(data, 11)
	r.FineY = GetRangeBit16(data, 12, 15)
	r.Unused = GetBit(data, 15)
}

type Olc2C02 struct {
	TabName          [][]uint8 // 2 1024
	TabPalette       []uint8   // 32
	Cartridge        *Cartridge
	PalCols          []color.Color   // 0x40
	Screen           *ebiten.Image   // 256 * 240
	NameTable        []*ebiten.Image // 2 * 256 * 240
	PatternTable     []*ebiten.Image // 2 * 128 * 128
	FrameComplete    bool
	Scanline         int16
	Cycle            int16
	Status           *Status
	Mask             *Mask
	Ctrl             *Ctrl
	AddrLatch        bool
	PpuDataBuffer    uint8
	PpuAddr          uint16
	Nmi              bool
	VRamAddr         *LoopReg
	TRamAddr         *LoopReg
	FineX            uint8
	BgNextTileID     uint8
	BgNextTileAttr   uint8
	BgNextTileLo     uint8
	BgNextTileHi     uint8
	BgShiftPatternLo uint16
	BgShiftPatternHi uint16
	BgShiftAttrLo    uint16
	BgShiftAttrHi    uint16
}

func NewOlc2C02() *Olc2C02 {
	palCols := make([]color.Color, 0x40)
	palCols[0x00] = color.RGBA{R: 84, G: 84, B: 84, A: 255}
	palCols[0x01] = color.RGBA{G: 30, B: 116, A: 255}
	palCols[0x02] = color.RGBA{R: 8, G: 16, B: 144, A: 255}
	palCols[0x03] = color.RGBA{R: 48, B: 136, A: 255}
	palCols[0x04] = color.RGBA{R: 68, B: 100, A: 255}
	palCols[0x05] = color.RGBA{R: 92, B: 48, A: 255}
	palCols[0x06] = color.RGBA{R: 84, G: 4, A: 255}
	palCols[0x07] = color.RGBA{R: 60, G: 24, A: 255}
	palCols[0x08] = color.RGBA{R: 32, G: 42, A: 255}
	palCols[0x09] = color.RGBA{R: 8, G: 58, A: 255}
	palCols[0x0A] = color.RGBA{G: 64, A: 255}
	palCols[0x0B] = color.RGBA{G: 60, A: 255}
	palCols[0x0C] = color.RGBA{G: 50, B: 60, A: 255}
	palCols[0x0D] = color.RGBA{A: 255}
	palCols[0x0E] = color.RGBA{A: 255}
	palCols[0x0F] = color.RGBA{A: 255}

	palCols[0x10] = color.RGBA{R: 152, G: 150, B: 152, A: 255}
	palCols[0x11] = color.RGBA{R: 8, G: 76, B: 196, A: 255}
	palCols[0x12] = color.RGBA{R: 48, G: 50, B: 236, A: 255}
	palCols[0x13] = color.RGBA{R: 92, G: 30, B: 228, A: 255}
	palCols[0x14] = color.RGBA{R: 136, G: 20, B: 176, A: 255}
	palCols[0x15] = color.RGBA{R: 160, G: 20, B: 100, A: 255}
	palCols[0x16] = color.RGBA{R: 152, G: 34, B: 32, A: 255}
	palCols[0x17] = color.RGBA{R: 120, G: 60, A: 255}
	palCols[0x18] = color.RGBA{R: 84, G: 90, A: 255}
	palCols[0x19] = color.RGBA{R: 40, G: 114, A: 255}
	palCols[0x1A] = color.RGBA{R: 8, G: 124, A: 255}
	palCols[0x1B] = color.RGBA{G: 118, B: 40, A: 255}
	palCols[0x1C] = color.RGBA{G: 102, B: 120, A: 255}
	palCols[0x1D] = color.RGBA{A: 255}
	palCols[0x1E] = color.RGBA{A: 255}
	palCols[0x1F] = color.RGBA{A: 255}

	palCols[0x20] = color.RGBA{R: 236, G: 238, B: 236, A: 255}
	palCols[0x21] = color.RGBA{R: 76, G: 154, B: 236, A: 255}
	palCols[0x22] = color.RGBA{R: 120, G: 124, B: 236, A: 255}
	palCols[0x23] = color.RGBA{R: 176, G: 98, B: 236, A: 255}
	palCols[0x24] = color.RGBA{R: 228, G: 84, B: 236, A: 255}
	palCols[0x25] = color.RGBA{R: 236, G: 88, B: 180, A: 255}
	palCols[0x26] = color.RGBA{R: 236, G: 106, B: 100, A: 255}
	palCols[0x27] = color.RGBA{R: 212, G: 136, B: 32, A: 255}
	palCols[0x28] = color.RGBA{R: 160, G: 170, A: 255}
	palCols[0x29] = color.RGBA{R: 116, G: 196, A: 255}
	palCols[0x2A] = color.RGBA{R: 76, G: 208, B: 32, A: 255}
	palCols[0x2B] = color.RGBA{R: 56, G: 204, B: 108, A: 255}
	palCols[0x2C] = color.RGBA{R: 56, G: 180, B: 204, A: 255}
	palCols[0x2D] = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	palCols[0x2E] = color.RGBA{A: 255}
	palCols[0x2F] = color.RGBA{A: 255}

	palCols[0x30] = color.RGBA{R: 236, G: 238, B: 236, A: 255}
	palCols[0x31] = color.RGBA{R: 168, G: 204, B: 236, A: 255}
	palCols[0x32] = color.RGBA{R: 188, G: 188, B: 236, A: 255}
	palCols[0x33] = color.RGBA{R: 212, G: 178, B: 236, A: 255}
	palCols[0x34] = color.RGBA{R: 236, G: 174, B: 236, A: 255}
	palCols[0x35] = color.RGBA{R: 236, G: 174, B: 212, A: 255}
	palCols[0x36] = color.RGBA{R: 236, G: 180, B: 176, A: 255}
	palCols[0x37] = color.RGBA{R: 228, G: 196, B: 144, A: 255}
	palCols[0x38] = color.RGBA{R: 204, G: 210, B: 120, A: 255}
	palCols[0x39] = color.RGBA{R: 180, G: 222, B: 120, A: 255}
	palCols[0x3A] = color.RGBA{R: 168, G: 226, B: 144, A: 255}
	palCols[0x3B] = color.RGBA{R: 152, G: 226, B: 180, A: 255}
	palCols[0x3C] = color.RGBA{R: 160, G: 214, B: 228, A: 255}
	palCols[0x3D] = color.RGBA{R: 160, G: 162, B: 160, A: 255}
	palCols[0x3E] = color.RGBA{A: 255}
	palCols[0x3F] = color.RGBA{A: 255}
	return &Olc2C02{
		TabName:      [][]uint8{make([]uint8, 1024), make([]uint8, 1024)},
		PalCols:      palCols,
		Screen:       ebiten.NewImage(256, 240),
		NameTable:    []*ebiten.Image{ebiten.NewImage(256, 240), ebiten.NewImage(256, 240)},
		PatternTable: []*ebiten.Image{ebiten.NewImage(128, 128), ebiten.NewImage(128, 128)},
		TabPalette:   make([]uint8, 32),
		Status:       &Status{},
		Mask:         &Mask{},
		Ctrl:         &Ctrl{},
		VRamAddr:     &LoopReg{},
		TRamAddr:     &LoopReg{},
	}
}

func (o *Olc2C02) UpdatePatternTable(index, palette uint8) {
	for tileY := uint16(0); tileY < 16; tileY++ {
		for tileX := uint16(0); tileX < 16; tileX++ {
			offset := (tileY*16 + tileX) * 8 * 2 // 一共过了tileY*16 + tileX 个Tile，每个Tile 8*8 需要 8*2 byte
			for row := uint16(0); row < 8; row++ {
				// 每行8个像素由2byte组成 i 是第几个tile表(一共2个) 一共8行，所以另外一个byte需要偏移8
				tileHi := o.PpuRead(uint16(index)*0x1000+offset+row+0x0000, false)
				tileLo := o.PpuRead(uint16(index)*0x1000+offset+row+0x0008, false)
				for col := uint16(0); col < 8; col++ {
					// 拼接高位与地位获取索引 获取颜色
					i := (tileHi&0x01)<<1 | (tileLo & 0x01)
					tileHi >>= 1
					tileLo >>= 1 // 之所以 7-col是因为 这里是从低位开始计算的，每次位移抹除的也是低位
					o.PatternTable[index].Set(int(tileX*8+(7-col)), int(tileY*8+row), o.GetColorFromPalette(palette, i))
				}
			}
		}
	}
}

func (o *Olc2C02) GetColorFromPalette(palette uint8, index uint8) color.Color {
	// 基本偏移  0x3F00    索引0~3正好占用低两位   调色板索引直接位移即可(*4)
	return o.PalCols[o.PpuRead(0x3F00+(uint16(palette)<<2)+uint16(index), false)&0x3F]
}

// read=true主要是给反编译使用的，因为读取是有损的，要防止其破坏系统
func (o *Olc2C02) CpuRead(addr uint16, read bool) uint8 {
	if read {
		switch addr {
		case 0x0000: // Ctrl
			return o.Ctrl.Get()
		case 0x0001: // Mask
			return o.Mask.Get()
		case 0x0002: // Status
			return o.Status.Get()
		}
	} else {
		switch addr {
		case 0x0002: // Status
			res := (o.Status.Get() & 0xE0) | (o.PpuDataBuffer & 0x1F)
			o.Status.VerticalBlank = false
			o.AddrLatch = false
			return res
		case 0x0007: // Data 延迟读取的
			res := o.PpuDataBuffer
			oldAddr := o.VRamAddr.Get()
			o.PpuDataBuffer = o.PpuRead(oldAddr, false)
			// 调色板是不需要延迟的
			if oldAddr >= 0x3F00 {
				res = o.PpuDataBuffer
			}
			o.VRamAddr.Set(oldAddr + If(o.Ctrl.IncrementMode, uint16(32), uint16(1))) // 是连续读取的
			return res
		}
	}
	return 0
}

func (o *Olc2C02) CpuWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Ctrl
		o.Ctrl.Set(data)
		o.TRamAddr.NameTableX = o.Ctrl.NameTableX
		o.TRamAddr.NameTableY = o.Ctrl.NameTableY
	case 0x0001: // Mask
		o.Mask.Set(data)
	case 0x0005: // Scroll
		if o.AddrLatch { // 连续写入 x,y偏移
			o.TRamAddr.FineY = uint16(data & 0x07)
			o.TRamAddr.CoarseY = uint16(data >> 3)
		} else {
			o.FineX = data & 0x07
			o.TRamAddr.CoarseX = uint16(data >> 3)
		}
		o.AddrLatch = !o.AddrLatch
	case 0x0006: // Addr 是分两次写入的
		if o.AddrLatch {
			o.TRamAddr.Set((o.TRamAddr.Get() & 0xFF00) | uint16(data))
			o.VRamAddr.Set(o.TRamAddr.Get())
		} else {
			o.TRamAddr.Set((o.TRamAddr.Get() & 0x00FF) | (uint16(data) << 8))
		}
		o.AddrLatch = !o.AddrLatch
	case 0x0007:
		o.PpuWrite(o.PpuAddr, data)
		o.PpuAddr += If(o.Ctrl.IncrementMode, uint16(32), uint16(1))
	}
}

func (o *Olc2C02) PpuRead(addr uint16, read bool) uint8 {
	addr &= 0x3FFF
	if res, ok := o.Cartridge.PpuRead(addr); ok {
		return res
	} else if addr < 0x2000 {
		panic("不支持，应该被Cartridge处理")
		// 这里不应该进来的
		// 获取指定地址的样式  前面正好是第一或2个图案表  后面是图案表中的位置
		//return o.TabName[(addr&0x1000)>>12][addr&0x0FFF]
	} else if addr < 0x3F00 {
		addr &= 0x0FFF
		if o.Cartridge.Mirror == VERTICAL { // 垂直镜像分别为4个屏幕设置画面
			if addr < 0x0400 { // 用于操作镜像数据
				return o.TabName[0][addr&0x03FF]
			} else if addr < 0x0800 {
				return o.TabName[1][addr&0x03FF]
			} else if addr < 0x0C00 {
				return o.TabName[0][addr&0x03FF]
			} else if addr < 0x1000 {
				return o.TabName[1][addr&0x03FF]
			}
		} else if o.Cartridge.Mirror == HORIZONTAL {
			if addr < 0x0400 {
				return o.TabName[0][addr&0x03FF]
			} else if addr < 0x0800 {
				return o.TabName[0][addr&0x03FF]
			} else if addr < 0x0C00 {
				return o.TabName[1][addr&0x03FF]
			} else if addr < 0x1000 {
				return o.TabName[1][addr&0x03FF]
			}
		}
	} else if addr < 0x4000 { // 调色板
		addr &= 0x1F
		switch addr { // 背景透明色处理
		case 0x10, 0x14, 0x18, 0x1C:
			addr -= 0x10
		}
		return o.TabPalette[addr]
	}
	return 0
}

func (o *Olc2C02) PpuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF
	if o.Cartridge.PpuWrite(addr, data) {

	} else if addr < 0x2000 {
		// 不支持写入，没有读取消费的地方
		//panic("不支持，应该被Cartridge处理")
	} else if addr < 0x3F00 {
		addr &= 0x0FFF
		if o.Cartridge.Mirror == VERTICAL { // 垂直镜像分别为4个屏幕设置画面
			if addr < 0x0400 { // 0X03FF = 1K -1
				o.TabName[0][addr&0x03FF] = data
			} else if addr < 0x0800 {
				o.TabName[1][addr&0x03FF] = data
			} else if addr < 0x0C00 {
				o.TabName[0][addr&0x03FF] = data
			} else if addr < 0x1000 {
				o.TabName[1][addr&0x03FF] = data
			}
		} else if o.Cartridge.Mirror == HORIZONTAL {
			if addr < 0x0400 {
				o.TabName[0][addr&0x03FF] = data
			} else if addr < 0x0800 {
				o.TabName[0][addr&0x03FF] = data
			} else if addr < 0x0C00 {
				o.TabName[1][addr&0x03FF] = data
			} else if addr < 0x1000 {
				o.TabName[1][addr&0x03FF] = data
			}
		}
	} else if addr < 0x4000 { // 调色板处理
		addr &= 0x1F
		switch addr { // 背景透明色处理
		case 0x10, 0x14, 0x18, 0x1C:
			addr -= 0x10
		}
		o.TabPalette[addr] = data
	}
}

func (o *Olc2C02) ConnectCartridge(cartridge *Cartridge) {
	o.Cartridge = cartridge
}

func (o *Olc2C02) IncrementScrollX() {
	if o.Mask.RenderBackground || o.Mask.RenderSprite {
		if o.VRamAddr.CoarseX == 31 {
			o.VRamAddr.CoarseX = 0
			o.VRamAddr.NameTableX = !o.VRamAddr.NameTableX
		} else {
			o.VRamAddr.CoarseX++
		}
	}
}

func (o *Olc2C02) IncrementScrollY() {
	if o.Mask.RenderBackground || o.Mask.RenderSprite {
		if o.VRamAddr.FineY < 7 {
			o.VRamAddr.FineY++
		} else {
			o.VRamAddr.FineY = 0
			if o.VRamAddr.CoarseY == 29 {
				o.VRamAddr.CoarseY = 0
				o.VRamAddr.NameTableY = !o.VRamAddr.NameTableY
			} else if o.VRamAddr.CoarseY == 31 { // ?
				o.VRamAddr.CoarseY = 0
			} else {
				o.VRamAddr.CoarseY++
			}
		}
	}
}

func (o *Olc2C02) TransferAddressX() {
	if o.Mask.RenderBackground || o.Mask.RenderSprite {
		o.VRamAddr.NameTableX = o.TRamAddr.NameTableX
		o.VRamAddr.CoarseX = o.TRamAddr.CoarseX
	}
}

func (o *Olc2C02) TransferAddressY() {
	if o.Mask.RenderBackground || o.Mask.RenderSprite {
		o.VRamAddr.FineY = o.TRamAddr.FineY
		o.VRamAddr.NameTableY = o.TRamAddr.NameTableY
		o.VRamAddr.CoarseY = o.TRamAddr.CoarseY
	}
}

func (o *Olc2C02) LoadBackgroundShift() {
	o.BgShiftPatternLo = (o.BgShiftPatternLo & 0xFF00) | uint16(o.BgNextTileLo)
	o.BgShiftPatternHi = (o.BgShiftPatternHi & 0xFF00) | uint16(o.BgNextTileHi)
	o.BgShiftAttrLo = (o.BgShiftAttrLo & 0xFF00) | If[uint16](o.BgNextTileAttr&0x01 > 0, 0xFF, 0x00)
	o.BgShiftAttrHi = (o.BgShiftAttrHi & 0xFF00) | If[uint16](o.BgNextTileAttr&0x02 > 0, 0xFF, 0x00)
}

func (o *Olc2C02) UpdateBackgroundShift() {
	if o.Mask.RenderBackground {
		o.BgShiftPatternLo <<= 1
		o.BgShiftPatternHi <<= 1
		o.BgShiftAttrLo <<= 1
		o.BgShiftAttrHi <<= 1
	}
}

func (o *Olc2C02) Clock() {
	if o.Scanline >= -1 && o.Scanline < 240 {
		if o.Scanline == 0 && o.Cycle == 0 {
			o.Cycle = 1
		}
		if o.Scanline == -1 && o.Cycle == 1 {
			o.Status.VerticalBlank = false
		}
		if (o.Cycle >= 2 && o.Cycle < 258) || (o.Cycle >= 321 && o.Cycle < 338) {
			o.UpdateBackgroundShift()
			switch (o.Cycle - 1) % 8 {
			case 0:
				o.LoadBackgroundShift()
				o.BgNextTileID = o.PpuRead(0x2000|(o.VRamAddr.Get()&0x0FFF), false)
			case 2:
				// << 偏移是为了放到对应的位置  >> 偏移是为了共用属性， 4*4个格子公用一个属性
				o.BgNextTileAttr = o.PpuRead(0x23C0|SetBit[uint16](o.VRamAddr.NameTableY, 11)|
					SetBit[uint16](o.VRamAddr.NameTableX, 10)|
					(o.VRamAddr.CoarseY>>2<<3)| // 原来5位因为公用格子抹除2位所以偏移3位即可
					(o.VRamAddr.CoarseX>>2), false)
				// 调色盘使用 4*4 但是属性使用 2*2的共享,还需进一步分割
				if o.VRamAddr.CoarseY&0x02 > 0 {
					o.BgNextTileAttr >>= 4
				}
				if o.VRamAddr.CoarseX&0x02 > 0 {
					o.BgNextTileAttr >>= 2
				}
				o.BgNextTileAttr &= 0x03
			case 4:
				o.BgNextTileLo = o.PpuRead(SetBit[uint16](o.Ctrl.PatternBackground, 12)+
					uint16(o.BgNextTileID)<<4+o.VRamAddr.FineY+0, false)
			case 6:
				o.BgNextTileHi = o.PpuRead(SetBit[uint16](o.Ctrl.PatternBackground, 12)+
					uint16(o.BgNextTileID)<<4+o.VRamAddr.FineY+8, false)
			case 7:
				o.IncrementScrollX()
			}
		}
		if o.Cycle == 256 {
			o.IncrementScrollY()
		}
		if o.Cycle == 257 {
			o.LoadBackgroundShift()
			o.TransferAddressX()
		}
		if o.Cycle == 338 || o.Cycle == 340 {
			o.BgNextTileID = o.PpuRead(0x2000|(o.VRamAddr.Get()&0x0FFF), false)
		}
		if o.Scanline == -1 && o.Cycle >= 280 && o.Cycle < 305 {
			o.TransferAddressY()
		}
	}

	if o.Scanline == 241 && o.Cycle == 1 {
		o.Status.VerticalBlank = true
		if o.Ctrl.EnableNmi {
			o.Nmi = true
		}
	}
	bgPixel := uint8(0x00)
	bgPalette := uint8(0x00)
	if o.Mask.RenderBackground {
		bitMask := uint16(0x8000 >> o.FineX) // 需要偏移器对应位置下的偏移
		pixelLo := If[uint8]((o.BgShiftPatternLo&bitMask) > 0, 1, 0)
		pixelHi := If[uint8]((o.BgShiftPatternHi&bitMask) > 0, 1, 0)
		bgPixel = (pixelHi << 1) | pixelLo
		bgPalLo := If[uint8]((o.BgShiftAttrLo&bitMask) > 0, 1, 0)
		bgPalHi := If[uint8]((o.BgShiftAttrLo&bitMask) > 0, 1, 0)
		bgPalette = (bgPalHi << 1) | bgPalLo
	}
	o.Screen.Set(int(o.Cycle-1), int(o.Scanline), o.GetColorFromPalette(bgPalette, bgPixel))
	o.Cycle++
	if o.Cycle >= 341 {
		o.Cycle = 0
		o.Scanline++
		if o.Scanline >= 261 {
			o.Scanline = -1
			o.FrameComplete = true
		}
	}
}

func (o *Olc2C02) Reset() {
	o.FineX = 0
	o.AddrLatch = false
	o.PpuDataBuffer = 0
	o.Scanline = 0
	o.Cycle = 0
	o.BgNextTileID = 0
	o.BgNextTileAttr = 0
	o.BgNextTileHi = 0
	o.BgNextTileLo = 0
	o.BgShiftAttrHi = 0
	o.BgShiftAttrLo = 0
	o.BgShiftPatternHi = 0
	o.BgShiftPatternLo = 0
	o.Status.Set(0)
	o.Mask.Set(0)
	o.Ctrl.Set(0)
	o.VRamAddr.Set(0)
	o.TRamAddr.Set(0)
}
