/*
@author: sk
@date: 2024/2/25
*/
package nes5

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CpuApp struct {
	Nes       *Bus
	codeMap   map[uint16]string
	codeLines []uint16
	codeIndex map[uint16]int
}

func NewCpuApp() *CpuApp {
	ebiten.SetWindowSize(488*2, 360*2)
	nes := NewBus()
	// 设置代码
	code := []uint8{0xA2, 0x0A, 0x8E, 0x00, 0x00, 0xA2, 0x03, 0x8E, 0x01, 0x00, 0xAC, 0x00, 0x00, 0xA9, 0x00, 0x18, 0x6D, 0x01, 0x00, 0x88, 0xD0, 0xFA, 0x8D, 0x02, 0x00, 0xEA, 0xEA, 0xEA}
	addr := 0x8000
	for i := 0; i < len(code); i++ {
		nes.CpuRam[addr+i] = code[i]
	}
	// 设置入口位置 若要测试需要更改CpuRam的大小
	//nes.CpuRam[0xFFFC] = 0x00
	//nes.CpuRam[0xFFFD] = 0x80
	// 反编译+cpu 复位
	codeMap := nes.Cpu.Disassemble(uint16(addr), uint16(addr+len(code)))
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
	nes.Cpu.Reset()
	return &CpuApp{
		Nes:       nes,
		codeMap:   codeMap,
		codeLines: codeLines,
		codeIndex: codeIndex,
	}
}

func (c *CpuApp) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.Nes.Cpu.Clock()
		for !c.Nes.Cpu.Complete() {
			c.Nes.Cpu.Clock()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		c.Nes.Cpu.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		c.Nes.Cpu.Irq()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		c.Nes.Cpu.Nmi()
	}
	return nil
}

func (c *CpuApp) Draw(screen *ebiten.Image) {
	c.DrawRam(screen, 2, 2, 0x0000, 16, 16)
	c.DrawRam(screen, 2, 182, 0x8000, 16, 16)
	c.DrawCpu(screen, 344, 2)
	c.DrawCode(screen, 344, 72, 26)
}

func (c *CpuApp) Layout(w, h int) (int, int) {
	return w / 2, h / 2
}

func (c *CpuApp) DrawRam(screen *ebiten.Image, x int, y int, addr uint16, xNum int, yNum int) {
	buff := strings.Builder{}
	for ty := 0; ty < yNum; ty++ {
		buff.Reset()
		buff.WriteString(fmt.Sprintf("$%s:", Format(addr, 4)))
		for tx := 0; tx < xNum; tx++ {
			buff.WriteString(fmt.Sprintf(" %s", Format(c.Nes.CpuRead(addr, true), 2)))
			addr++
		}
		ebitenutil.DebugPrintAt(screen, buff.String(), x, y+ty*10)
	}
}

func (c *CpuApp) DrawCpu(screen *ebiten.Image, x int, y int) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STATUS:%08bCZIDBUVN", c.Nes.Cpu.Status), x, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PC: $%s", Format(c.Nes.Cpu.Pc, 4)), x, y+10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A: $%s", Format(c.Nes.Cpu.A, 2)), x, y+20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: $%s", Format(c.Nes.Cpu.X, 2)), x, y+30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: $%s", Format(c.Nes.Cpu.Y, 2)), x, y+40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("StackP: $%s", Format(c.Nes.Cpu.StackP, 4)), x, y+50)
}

func (c *CpuApp) DrawCode(screen *ebiten.Image, x int, y int, num int) {
	index := c.codeIndex[c.Nes.Cpu.Pc]
	num = Min(num, len(c.codeLines)-index)
	for i := 0; i < num; i++ {
		line := c.codeLines[index+i]
		ebitenutil.DebugPrintAt(screen, c.codeMap[line], x, y+i*10)
	}
}
