/*
@author: sk
@date: 2024/2/26
*/
package nes5

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetRangeBit8(t *testing.T) {
	res := GetRangeBit8(0b001011, 1, 4)
	fmt.Println(strconv.FormatInt(int64(res), 2))
}

func TestSetRangeBit8(t *testing.T) {
	res := SetRangeBit8(0b001110, 1, 4)
	fmt.Println(strconv.FormatInt(int64(res), 2))
}

func TestGetRangeBit16(t *testing.T) {
	res := GetRangeBit16(0b0010101001011, 3, 7)
	fmt.Println(strconv.FormatInt(int64(res), 2))
}

func TestSetRangeBit16(t *testing.T) {
	res := SetRangeBit16(0b011101110, 1, 4)
	fmt.Println(strconv.FormatInt(int64(res), 2))
}

func TestGetBit(t *testing.T) {
	res := GetBit[uint8](0b001011, 3)
	fmt.Println(res)
}

func TestSetBit(t *testing.T) {
	res := SetBit[uint8](true, 3)
	fmt.Println(strconv.FormatInt(int64(res), 2))
}

func TestStatus(t *testing.T) {
	s := &Status{}
	s.Unused = 22
	v := s.Get()
	s.Unused = 0
	s.SpriteOverflow = true
	v = s.Get()
	s.SpriteZeroHit = true
	v = s.Get()
	s.VerticalBlank = true
	v = s.Get()
	fmt.Println(v)
}

func TestLoop(t *testing.T) {
	l := &LoopReg{}
	for i := 0; i < 0xFFFF; i += i + 1 {
		l.Set(uint16(i))
		fmt.Println(l)
	}
	l.CoarseX = 22
	v := l.Get()
	l.CoarseY = 11
	v = l.Get()
	l.NameTableX = true
	v = l.Get()
	l.NameTableY = true
	v = l.Get()
	l.FineY = 2
	v = l.Get()
	fmt.Println(v)
}

func TestVal(t *testing.T) {
	//[]*Instruction{
	//	{"BRK", res.BRK, 5, 7}, {"ORA", res.ORA, 7, 6}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 3}, {"ORA", res.ORA, 11, 3}, {"ASL", res.ASL, 11, 5}, {"???", res.XXX, 6, 5}, {"PHP", res.PHP, 6, 3}, {"ORA", res.ORA, 5, 2}, {"ASL", res.ASL, 6, 2}, {"???", res.XXX, 6, 2}, {"???", res.NOP, 6, 4}, {"ORA", res.ORA, 1, 4}, {"ASL", res.ASL, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BPL", res.BPL, 10, 2}, {"ORA", res.ORA, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"ORA", res.ORA, 12, 4}, {"ASL", res.ASL, 12, 6}, {"???", res.XXX, 6, 6}, {"CLC", res.CLC, 6, 2}, {"ORA", res.ORA, 3, 4}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"ORA", res.ORA, 2, 4}, {"ASL", res.ASL, 2, 7}, {"???", res.XXX, 6, 7},
	//	{"JSR", res.JSR, 1, 6}, {"AND", res.AND, 7, 6}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"BIT", res.BIT, 11, 3}, {"AND", res.AND, 11, 3}, {"ROL", res.ROL, 11, 5}, {"???", res.XXX, 6, 5}, {"PLP", res.PLP, 6, 4}, {"AND", res.AND, 5, 2}, {"ROL", res.ROL, 6, 2}, {"???", res.XXX, 6, 2}, {"BIT", res.BIT, 1, 4}, {"AND", res.AND, 1, 4}, {"ROL", res.ROL, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BMI", res.BMI, 10, 2}, {"AND", res.AND, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"AND", res.AND, 12, 4}, {"ROL", res.ROL, 12, 6}, {"???", res.XXX, 6, 6}, {"SEC", res.SEC, 6, 2}, {"AND", res.AND, 3, 4}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"AND", res.AND, 2, 4}, {"ROL", res.ROL, 2, 7}, {"???", res.XXX, 6, 7},
	//	{"RTI", res.RTI, 6, 6}, {"EOR", res.EOR, 7, 6}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 3}, {"EOR", res.EOR, 11, 3}, {"LSR", res.LSR, 11, 5}, {"???", res.XXX, 6, 5}, {"PHA", res.PHA, 6, 3}, {"EOR", res.EOR, 5, 2}, {"LSR", res.LSR, 6, 2}, {"???", res.XXX, 6, 2}, {"JMP", res.JMP, 1, 3}, {"EOR", res.EOR, 1, 4}, {"LSR", res.LSR, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BVC", res.BVC, 10, 2}, {"EOR", res.EOR, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"EOR", res.EOR, 12, 4}, {"LSR", res.LSR, 12, 6}, {"???", res.XXX, 6, 6}, {"CLI", res.CLI, 6, 2}, {"EOR", res.EOR, 3, 4}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"EOR", res.EOR, 2, 4}, {"LSR", res.LSR, 2, 7}, {"???", res.XXX, 6, 7},
	//	{"RTS", res.RTS, 6, 6}, {"ADC", res.ADC, 7, 6}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 3}, {"ADC", res.ADC, 11, 3}, {"ROR", res.ROR, 11, 5}, {"???", res.XXX, 6, 5}, {"PLA", res.PLA, 6, 4}, {"ADC", res.ADC, 5, 2}, {"ROR", res.ROR, 6, 2}, {"???", res.XXX, 6, 2}, {"JMP", res.JMP, 8, 5}, {"ADC", res.ADC, 1, 4}, {"ROR", res.ROR, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BVS", res.BVS, 10, 2}, {"ADC", res.ADC, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"ADC", res.ADC, 12, 4}, {"ROR", res.ROR, 12, 6}, {"???", res.XXX, 6, 6}, {"SEI", res.SEI, 6, 2}, {"ADC", res.ADC, 3, 4}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"ADC", res.ADC, 2, 4}, {"ROR", res.ROR, 2, 7}, {"???", res.XXX, 6, 7},
	//	{"???", res.NOP, 6, 2}, {"STA", res.STA, 7, 6}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 6}, {"STY", res.STY, 11, 3}, {"STA", res.STA, 11, 3}, {"STX", res.STX, 11, 3}, {"???", res.XXX, 6, 3}, {"DEY", res.DEY, 6, 2}, {"???", res.NOP, 6, 2}, {"TXA", res.TXA, 6, 2}, {"???", res.XXX, 6, 2}, {"STY", res.STY, 1, 4}, {"STA", res.STA, 1, 4}, {"STX", res.STX, 1, 4}, {"???", res.XXX, 6, 4},
	//	{"BCC", res.BCC, 10, 2}, {"STA", res.STA, 9, 6}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 6}, {"STY", res.STY, 12, 4}, {"STA", res.STA, 12, 4}, {"STX", res.STX, 13, 4}, {"???", res.XXX, 6, 4}, {"TYA", res.TYA, 6, 2}, {"STA", res.STA, 3, 5}, {"TXS", res.TXS, 6, 2}, {"???", res.XXX, 6, 5}, {"???", res.NOP, 6, 5}, {"STA", res.STA, 2, 5}, {"???", res.XXX, 6, 5}, {"???", res.XXX, 6, 5},
	//	{"LDY", res.LDY, 5, 2}, {"LDA", res.LDA, 7, 6}, {"LDX", res.LDX, 5, 2}, {"???", res.XXX, 6, 6}, {"LDY", res.LDY, 11, 3}, {"LDA", res.LDA, 11, 3}, {"LDX", res.LDX, 11, 3}, {"???", res.XXX, 6, 3}, {"TAY", res.TAY, 6, 2}, {"LDA", res.LDA, 5, 2}, {"TAX", res.TAX, 6, 2}, {"???", res.XXX, 6, 2}, {"LDY", res.LDY, 1, 4}, {"LDA", res.LDA, 1, 4}, {"LDX", res.LDX, 1, 4}, {"???", res.XXX, 6, 4},
	//	{"BCS", res.BCS, 10, 2}, {"LDA", res.LDA, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 5}, {"LDY", res.LDY, 12, 4}, {"LDA", res.LDA, 12, 4}, {"LDX", res.LDX, 13, 4}, {"???", res.XXX, 6, 4}, {"CLV", res.CLV, 6, 2}, {"LDA", res.LDA, 3, 4}, {"TSX", res.TSX, 6, 2}, {"???", res.XXX, 6, 4}, {"LDY", res.LDY, 2, 4}, {"LDA", res.LDA, 2, 4}, {"LDX", res.LDX, 3, 4}, {"???", res.XXX, 6, 4},
	//	{"CPY", res.CPY, 5, 2}, {"CMP", res.CMP, 7, 6}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 8}, {"CPY", res.CPY, 11, 3}, {"CMP", res.CMP, 11, 3}, {"DEC", res.DEC, 11, 5}, {"???", res.XXX, 6, 5}, {"INY", res.INY, 6, 2}, {"CMP", res.CMP, 5, 2}, {"DEX", res.DEX, 6, 2}, {"???", res.XXX, 6, 2}, {"CPY", res.CPY, 1, 4}, {"CMP", res.CMP, 1, 4}, {"DEC", res.DEC, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BNE", res.BNE, 10, 2}, {"CMP", res.CMP, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"CMP", res.CMP, 12, 4}, {"DEC", res.DEC, 12, 6}, {"???", res.XXX, 6, 6}, {"CLD", res.CLD, 6, 2}, {"CMP", res.CMP, 3, 4}, {"NOP", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"CMP", res.CMP, 2, 4}, {"DEC", res.DEC, 2, 7}, {"???", res.XXX, 6, 7},
	//	{"CPX", res.CPX, 5, 2}, {"SBC", res.SBC, 7, 6}, {"???", res.NOP, 6, 2}, {"???", res.XXX, 6, 8}, {"CPX", res.CPX, 11, 3}, {"SBC", res.SBC, 11, 3}, {"INC", res.INC, 11, 5}, {"???", res.XXX, 6, 5}, {"INX", res.INX, 6, 2}, {"SBC", res.SBC, 5, 2}, {"NOP", res.NOP, 6, 2}, {"???", res.SBC, 6, 2}, {"CPX", res.CPX, 1, 4}, {"SBC", res.SBC, 1, 4}, {"INC", res.INC, 1, 6}, {"???", res.XXX, 6, 6},
	//	{"BEQ", res.BEQ, 10, 2}, {"SBC", res.SBC, 9, 5}, {"???", res.XXX, 6, 2}, {"???", res.XXX, 6, 8}, {"???", res.NOP, 6, 4}, {"SBC", res.SBC, 12, 4}, {"INC", res.INC, 12, 6}, {"???", res.XXX, 6, 6}, {"SED", res.SED, 6, 2}, {"SBC", res.SBC, 3, 4}, {"NOP", res.NOP, 6, 2}, {"???", res.XXX, 6, 7}, {"???", res.NOP, 6, 4}, {"SBC", res.SBC, 2, 4}, {"INC", res.INC, 2, 7}, {"???", res.XXX, 6, 7},
	//}
}
