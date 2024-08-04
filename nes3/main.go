/*
@author: sk
@date: 2023/11/11
*/
package main

import (
	"fmt"
	"image"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func main() {
	ebiten.SetWindowSize(780*2, 480*2)
	err := ebiten.RunGame(NewMainGame())
	HandleErr(err)
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

type MainGame struct {
	bus             *Bus
	codes           map[uint16]string
	addrs           []uint16
	option          *ebiten.DrawImageOptions
	selectedPalette uint8 // 当前选择的调色板
}

func NewMainGame() *MainGame {
	//// 注入测试代码
	//codes := strings.Split("A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA", " ")
	//offset := 0x8000
	//bus := NewBus()
	//for _, code := range codes {
	//	temp, _ := strconv.ParseUint(code, 16, 8)
	//	bus.CpuRam[offset] = uint8(temp)
	//	offset++
	//}
	//// 指定代码位置
	//bus.CpuRam[0xFFFC] = 0x00
	//bus.CpuRam[0xFFFD] = 0x80
	//bus.Cpu.Reset()
	//// 反编译代码
	//res, addrs := bus.Cpu.Disassemble(0x8000, uint16(offset))
	bus := NewBus()
	bus.InsertCartridge(NewCartridge("/Users/bytedance/Documents/go/nes/res/小蜜蜂.NES"))
	codes, addrs := bus.Cpu.Disassemble(0x0000, 0xFFFF)
	bus.Reset()
	return &MainGame{bus: bus, codes: codes, addrs: addrs, option: &ebiten.DrawImageOptions{}}
}

func (m *MainGame) Update() error {
	//if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	//	for true {
	//		m.bus.Clock()
	//		if m.bus.Cpu.Complete() {
	//			break
	//		}
	//	}
	//}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) { // 单步执行
		for true {
			m.bus.Clock()
			if m.bus.Ppu.FrameComplete {
				break
			}
		}
		m.bus.Ppu.FrameComplete = false
		for true {
			m.bus.Clock()
			if m.bus.Cpu.InstructionComplete {
				break
			}
		}
		m.bus.Cpu.InstructionComplete = false
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) { // 每秒 60帧执行
		for true {
			m.bus.Clock()
			if m.bus.Ppu.FrameComplete {
				break
			}
		}
		m.bus.Ppu.FrameComplete = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		m.bus.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) { // 切换调色板
		m.selectedPalette = (m.selectedPalette + 1) & 0x07
	}
	return nil
}

func (m *MainGame) Draw(screen *ebiten.Image) {
	//m.DrawRam(screen, 2, 2, 0x0000, 16, 16)
	//m.DrawRam(screen, 2, 182, 0x8000, 16, 16)
	m.DrawCpu(screen, 516, 2)
	for p := uint8(0); p < 8; p++ { // 绘制调色板
		for i := uint8(0); i < 4; i++ {
			clr := m.bus.Ppu.GetColorFromPalette(p, i)
			vector.DrawFilledRect(screen, 516+float32((p*5+i)*5), 340, 5, 5, clr, false)
		}
	}
	vector.StrokeRect(screen, 516+float32(m.selectedPalette)*5*5, 340, 4*5, 5, 1, colornames.White, false)
	// 绘制tile
	tile0 := m.bus.Ppu.GetTile(0, m.selectedPalette)
	m.DrawImage(screen, tile0, 516, 348, 1)
	tile1 := m.bus.Ppu.GetTile(1, m.selectedPalette)
	m.DrawImage(screen, tile1, 648, 348, 1)
	//m.DrawImage(screen, m.bus.Ppu.Screen, 0, 0, 2) // 绘制主屏幕
	for y := 0; y < 30; y++ {
		for x := 0; x < 32; x++ {
			//str := Format(m.bus.Ppu.TabName[0][y*32+x], 2)// 绘制背景编号
			//ebitenutil.DebugPrintAt(screen, str[2:], x*16, y*16)
			index := m.bus.Ppu.TabName[0][y*32+x]
			m.DrawSubImage(screen, tile0, float64(x*8), float64(y*8), int(index&0x0F)*8, int(index>>4)*8, 8, 8, 1)
		}
	}
	for y := 0; y < 30; y++ {
		for x := 0; x < 32; x++ {
			//str := Format(m.bus.Ppu.TabName[0][y*32+x], 2)// 绘制背景编号
			//ebitenutil.DebugPrintAt(screen, str[2:], x*16, y*16)
			index := m.bus.Ppu.TabName[1][y*32+x]
			m.DrawSubImage(screen, tile0, float64(x*8)+32*8, float64(y*8), int(index&0x0F)*8, int(index>>4)*8, 8, 8, 1)
		}
	}
	m.DrawCode(screen, 516+6*3, 82+12*10, 12)
}

func (m *MainGame) Layout(w, h int) (int, int) {
	return w / 2, h / 2
}

func (m *MainGame) DrawRam(screen *ebiten.Image, posX int, posY int, addr uint16, xNum int, yNum int) {
	for y := 0; y < yNum; y++ {
		for x := 0; x < xNum; x++ {
			ebitenutil.DebugPrintAt(screen, Format(m.bus.CpuRead(addr), 2)[2:], posX+x*20, posY+y*10)
			addr++
		}
	}
}

func (m *MainGame) DrawCpu(screen *ebiten.Image, x int, y int) {
	cpu := m.bus.Cpu
	state := cpu.State
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STATUS:N(%v),V(%v),U(%v),B(%v)\nSTATUS:D(%v),I(%v),Z(%v),C(%v)", state.Negative, state.Overflow, state.Unused, state.Break, state.Decimal, state.UnInterrupt, state.Zero, state.Carry), x, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PC:%v", Format(cpu.Pc, 4)), x, y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A:%v", Format(cpu.A, 2)), x, y+40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X:%v", Format(cpu.X, 2)), x, y+50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y:%v", Format(cpu.Y, 2)), x, y+60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("StackP:%v", Format(cpu.StackP, 2)), x, y+70)
}

func (m *MainGame) DrawCode(screen *ebiten.Image, x int, y int, lines int) {
	index := -1
	for i, addr := range m.addrs {
		if addr == m.bus.Cpu.Pc {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	ebitenutil.DebugPrintAt(screen, "-->"+m.codes[m.addrs[index]], x-6*3, y)
	for i := 1; i < lines; i++ {
		if index+i < len(m.addrs) { // 向下扩展
			ebitenutil.DebugPrintAt(screen, m.codes[m.addrs[index+i]], x, y+i*10)
		}
		if index-i >= 0 { // 向上扩展
			ebitenutil.DebugPrintAt(screen, m.codes[m.addrs[index-i]], x, y-i*10)
		}
	}
}

func (m *MainGame) DrawImage(screen *ebiten.Image, img *ebiten.Image, x float64, y float64, scale float64) {
	m.option.GeoM.Reset()
	m.option.GeoM.Scale(scale, scale)
	m.option.GeoM.Translate(x, y)
	screen.DrawImage(img, m.option)
}

func (m *MainGame) DrawSubImage(screen *ebiten.Image, img *ebiten.Image, x float64, y float64,
	px int, py int, w int, h int, scale float64) {
	img = img.SubImage(image.Rect(px, py, px+w, py+h)).(*ebiten.Image)
	m.DrawImage(screen, img, x, y, scale)
}
