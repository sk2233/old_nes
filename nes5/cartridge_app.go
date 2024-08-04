/*
@author: sk
@date: 2024/2/25
*/
package nes5

import (
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CartridgeApp struct {
	Nes       *Bus
	CodeMap   map[uint16]string
	CodeLines []uint16
	CodeIndex map[uint16]int
}

func (c *CartridgeApp) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		for !c.Nes.Cpu.Complete() { // 走一个时钟完结 一个code
			c.Nes.Clock()
		}
		for c.Nes.Cpu.Complete() { // 读取下一个指令初始化下一个时钟
			c.Nes.Clock()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		for !c.Nes.Ppu.FrameComplete { // 运行整个一帧
			c.Nes.Clock()
		}
		for !c.Nes.Cpu.Complete() { // 并保证运算补齐
			c.Nes.Clock()
		}
		c.Nes.Ppu.FrameComplete = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		c.Nes.Reset()
	}
	return nil
}

func (c *CartridgeApp) DrawCpu(screen *ebiten.Image, x int, y int) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STATUS:%08bCZIDBUVN", c.Nes.Cpu.Status), x, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PC: $%s", Format(c.Nes.Cpu.Pc, 4)), x, y+10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A: $%s", Format(c.Nes.Cpu.A, 2)), x, y+20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: $%s", Format(c.Nes.Cpu.X, 2)), x, y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: $%s", Format(c.Nes.Cpu.Y, 2)), x, y+40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("StackP: $%s", Format(c.Nes.Cpu.StackP, 4)), x, y+50)
}

func (c *CartridgeApp) DrawCode(screen *ebiten.Image, x int, y int, num int) {
	index := c.CodeIndex[c.Nes.Cpu.Pc]
	num = Min(num, len(c.CodeLines)-index)
	for i := 0; i < num; i++ {
		line := c.CodeLines[index+i]
		ebitenutil.DebugPrintAt(screen, c.CodeMap[line], x, y+i*10)
	}
}

func (c *CartridgeApp) Draw(screen *ebiten.Image) {
	c.DrawCpu(screen, 258, 2)
	c.DrawCode(screen, 258, 72, 16)
	DrawImage(screen, c.Nes.Ppu.Screen, 0, 0)
}

func (c *CartridgeApp) Layout(w, h int) (int, int) {
	return w / 2, h / 2
}

func NewCartridgeApp() *CartridgeApp {
	ebiten.SetWindowSize(400*2, 240*2)
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
	return &CartridgeApp{
		Nes:       nes,
		CodeMap:   codeMap,
		CodeIndex: codeIndex,
		CodeLines: codeLines,
	}
}
