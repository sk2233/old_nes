/*
@author: sk
@date: 2024/2/26
*/
package nes5

import (
	"fmt"
	"sort"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BackgroundApp struct {
	Nes       *Bus
	CodeMap   map[uint16]string
	CodeLines []uint16
	CodeIndex map[uint16]int
	Palette   uint8
}

func (b *BackgroundApp) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		b.Palette = (b.Palette + 1) % 8
		b.Nes.Ppu.UpdatePatternTable(0, b.Palette)
		b.Nes.Ppu.UpdatePatternTable(1, b.Palette)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		for !b.Nes.Cpu.Complete() { // 走一个时钟完结 一个code
			b.Nes.Clock()
		}
		for b.Nes.Cpu.Complete() { // 读取下一个指令初始化下一个时钟
			b.Nes.Clock()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		for !b.Nes.Ppu.FrameComplete { // 运行整个一帧
			b.Nes.Clock()
		}
		for !b.Nes.Cpu.Complete() { // 并保证运算补齐
			b.Nes.Clock()
		}
		b.Nes.Ppu.FrameComplete = false
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		for !b.Nes.Ppu.FrameComplete { // 运行整个一帧
			b.Nes.Clock()
		}
		for !b.Nes.Cpu.Complete() { // 并保证运算补齐
			b.Nes.Clock()
		}
		b.Nes.Ppu.FrameComplete = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		b.Nes.Reset()
	}
	return nil
}

func (b *BackgroundApp) DrawCpu(screen *ebiten.Image, x int, y int) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STATUS: %08b [CZIDBUVN]", b.Nes.Cpu.Status), x, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PC: $%s", Format(b.Nes.Cpu.Pc, 4)), x, y+10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A: $%s", Format(b.Nes.Cpu.A, 2)), x, y+20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: $%s", Format(b.Nes.Cpu.X, 2)), x, y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: $%s", Format(b.Nes.Cpu.Y, 2)), x, y+40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("StackP: $%s", Format(b.Nes.Cpu.StackP, 4)), x, y+50)
}

func (b *BackgroundApp) Draw(screen *ebiten.Image) {
	for y := 0; y < 30; y++ {
		for x := 0; x < 32; x++ {
			id := b.Nes.Ppu.TabName[0][y*32+x]
			DrawSubImage(screen, b.Nes.Ppu.PatternTable[0], float64(x*8), float64(y*8),
				int(id&0xF)*8, int(id>>4)*8, 8, 8)
		}
	}
	b.DrawCpu(screen, 514, 2)
	b.DrawCode(screen, 514, 72, 26)
	DrawImage(screen, b.Nes.Ppu.PatternTable[0], 512, 352)
	DrawImage(screen, b.Nes.Ppu.PatternTable[1], 640, 352)
	// 32 2 7 7 7 7 2
	// 8 * 4
	for palette := 0; palette < 8; palette++ {
		base := 512 + 32*palette + 2
		for index := 0; index < 4; index++ {
			clr := b.Nes.Ppu.GetColorFromPalette(uint8(palette), uint8(index))
			vector.DrawFilledRect(screen, float32(base+index*7), 343, 7, 7, clr, false)
		}
	}
	vector.StrokeRect(screen, 513+float32(b.Palette)*32, 342, 30, 9, 2, colornames.White, false)
}

func (b *BackgroundApp) Layout(w, h int) (int, int) {
	return w / 2, h / 2
}

func (b *BackgroundApp) DrawCode(screen *ebiten.Image, x int, y int, num int) {
	index := b.CodeIndex[b.Nes.Cpu.Pc]
	num = Min(num, len(b.CodeLines)-index)
	for i := 0; i < num; i++ {
		line := b.CodeLines[index+i]
		ebitenutil.DebugPrintAt(screen, b.CodeMap[line], x, y+i*10)
	}
}

func NewBackgroundApp() *BackgroundApp {
	ebiten.SetWindowSize(768*2, 480*2)
	cartridge := NewCartridge("/Users/bytedance/Documents/go/nes/res/nestest.nes")
	nes := NewBus()
	nes.InsertCartridge(cartridge)
	codeMap := nes.Cpu.Disassemble(0x0000, 0xFFFF)
	codeLines := make([]uint16, 0, len(codeMap))
	for key := range codeMap {
		codeLines = append(codeLines, key)
	}
	sort.Slice(codeLines, func(i, j int) bool {
		return codeLines[i] < codeLines[j]
	})
	codeIndex := make(map[uint16]int, len(codeLines))
	for i := 0; i < len(codeLines); i++ {
		codeIndex[codeLines[i]] = i
	}
	nes.Reset()
	return &BackgroundApp{
		Nes:       nes,
		CodeMap:   codeMap,
		CodeLines: codeLines,
		CodeIndex: codeIndex,
	}
}
