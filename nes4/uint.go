/*
@author: sk
@date: 2023/10/28
*/
package nes4

type bitRange struct {
	Start, End int
}

var (
	StatusUnused         = &bitRange{0, 5}
	StatusSpriteOverflow = &bitRange{5, 6}
	StatusSpriteZeroHit  = &bitRange{6, 7}
	StatusVerticalBlank  = &bitRange{7, 8}
)

var (
	MaskGrayScale            = &bitRange{0, 1}
	MaskRenderBackgroundLeft = &bitRange{1, 2}
	MaskRenderSpriteLeft     = &bitRange{2, 3}
	MaskRenderBackground     = &bitRange{3, 4}
	MaskRenderSprite         = &bitRange{4, 5}
	MaskEnhanceRed           = &bitRange{5, 6}
	MaskEnhanceGreen         = &bitRange{6, 7}
	MaskEnhanceBlue          = &bitRange{7, 8}
)

var (
	ControlNameTableX        = &bitRange{0, 1}
	ControlNameTableY        = &bitRange{1, 2}
	ControlIncrementMode     = &bitRange{2, 3}
	ControlPatternSprite     = &bitRange{3, 4}
	ControlPatternBackground = &bitRange{4, 5}
	ControlSpriteSize        = &bitRange{5, 6}
	ControlUnused            = &bitRange{6, 7}
	ControlEnableNmi         = &bitRange{7, 8}
)

var (
	LoopCoarseX    = &bitRange{0, 5}   // 5
	LoopCoarseY    = &bitRange{5, 10}  // 5
	LoopNameTableX = &bitRange{10, 11} // 1
	LoopNameTableY = &bitRange{11, 12} // 1
	LoopFineY      = &bitRange{12, 15} // 3
	LoopUnused     = &bitRange{15, 16} // 1
)

type Uint8 struct {
	Data uint8
}

func NewUint8() *Uint8 {
	return &Uint8{}
}

// TODO 是从低位开始的

func (u *Uint8) Get(key *bitRange) uint8 {
	return u.Data << key.Start >> (8 - key.End + key.Start)
}

func (u *Uint8) Set(key *bitRange, value uint8) {
	value = value << (8 - (key.End - key.Start)) >> key.Start
	u.Data |= value
}

type Uint16 struct {
	Data uint16
}

func NewUint16() *Uint16 {
	return &Uint16{}
}

func (u *Uint16) Get(key *bitRange) uint16 {
	return u.Data << key.Start >> (16 - key.End + key.Start)
}

func (u *Uint16) Set(key *bitRange, value uint16) {
	value = value << (16 - (key.End - key.Start)) >> key.Start
	u.Data |= value
}
