/*
@author: sk
@date: 2023/10/14
*/
package test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"
	"strconv"
	"testing"

	"golang.org/x/image/colornames"
)

func TestReadFile(t *testing.T) {
	file, err := os.Open("/Users/bytedance/Documents/go/nes/res/小蜜蜂.NES")
	HandleErr(err)
	defer file.Close()
	header := parseHeader(file)
	if header.HasTrainer() { // 没用的部分
		ReadBytes(file, 512)
	}
	prgRam := ReadBytes(file, 16*1024*int(header.PrgCount))
	chrRam := ReadBytes(file, 8*1024*int(header.ChrCount))
	fmt.Println(len(prgRam))
	for i := 0; i < int(header.ChrCount)*2; i++ {
		out(chrRam[i*4*1024 : (i+1)*4*1024])
	}
}

var (
	imgIndex = 0
)

func out(bs []byte) {
	img := image.NewRGBA(image.Rect(0, 0, 16*8, 16*8))
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 16; tileX++ {
			offset := (tileY*16 + tileX) * 8 * 2
			for y := 0; y < 8; y++ {
				lo := bs[offset+y*2]
				hi := bs[offset+y*2+1]
				for x := 0; x < 8; x++ {
					index := (lo & 0x01) | ((hi & 0x01) << 1)
					lo >>= 1
					hi >>= 1
					img.Set(tileX*8+x, tileY*8+y, getColor(index))
				}
			}
		}
	}
	temp, err := os.Create(fmt.Sprintf("/Users/bytedance/Documents/go/nes/res/%v.png", imgIndex))
	HandleErr(err)
	imgIndex++
	defer temp.Close()
	err = png.Encode(temp, img)
	HandleErr(err)
}

func getColor(index byte) color.Color {
	switch index {
	case 0b00:
		return colornames.Red
	case 0b01:
		return colornames.Green
	case 0b10:
		return colornames.Blue
	case 0b11:
		return colornames.Black
	default:
		panic("index out of range")
	}
}

func parseHeader(file *os.File) *Header {
	name := ReadStr(file, 4)
	prgCount := ReadUint8(file)
	chrCount := ReadUint8(file)
	flag1 := ReadUint8(file)
	flag2 := ReadUint8(file)
	return &Header{
		Name:     name,
		PrgCount: prgCount,
		ChrCount: chrCount,
		Flag1:    flag1,
		Flag2:    flag2,
		UnUse:    ReadBytes(file, 8),
	}
}

type Header struct {
	Name         string // 0 ~ 3
	PrgCount     uint8  // 4
	ChrCount     uint8  // 5
	Flag1, Flag2 uint8  // 6 ~ 7
	UnUse        []byte // 8 ~ 15
}

func (h *Header) HasTrainer() bool {
	return h.Flag1&0x04 == 1
}

func (h *Header) GetMapperID() uint8 {
	return (h.Flag2 & 0xf0) | (h.Flag1 >> 4)
}

func (h *Header) Test() bool {
	return false
}

func TestAt(t *testing.T) {
	//num := uint64(0b00001)
	//fmt.Println(strconv.FormatUint(num, 2))
	//h := &Header{}
	//fmt.Println(reflect.ValueOf(h.Test).UnsafePointer() == reflect.ValueOf(h.HasTrainer).UnsafePointer())
	//fmt.Println(12 ^ 12)
	fmt.Println(strconv.FormatUint(uint64(2233), 16))
}

type TestIt struct {
	Names [10][10]string
}

func TestUint(t *testing.T) {
	//for i := 0; i < 100; i++ {
	//	num := rand.Uint32()
	//	fmt.Println(uint16(num) == uint16(num&0x0000FFFF))
	//}
	ti := &TestIt{}
	ti.Names[4][4] = "sfsdfsd"
	fmt.Println(ti.Names)
}

func TestUint2(t *testing.T) {
	num := uint8(0b001100)
	num = num << 2 >> 4
	fmt.Println(num)
}

type bitRange struct {
	Start, End int
}

type Uint8 struct {
	Data uint8
}

func (u *Uint8) Get(key *bitRange) uint8 {
	return u.Data << key.Start >> (8 - key.End + key.Start)
}

func (u *Uint8) Set(key *bitRange, value uint8) {
	value = value << (8 - (key.End - key.Start)) >> key.Start
	u.Data |= value
}

func TestReadUint8(t *testing.T) {
	u := &Uint8{}
	u.Set(&bitRange{1, 2}, 1)
	fmt.Println(u.Data)
	fmt.Println(u.Get(&bitRange{2, 4}))
}

func TestOver(t *testing.T) {
	var a uint8
	a = 255
	fmt.Println(a + 22)
	fmt.Println((uint16(a) + 22) & 0x00FF)
}

type Temp struct {
	FuncA func()
	FuncB func()
}

func (t *Temp) Test() {

}

func TestFunc(t *testing.T) {
	ts := &Temp{}
	ts.FuncA = func() {

	}
	ts.FuncB = ts.Test
	fmt.Println(reflect.ValueOf(ts.FuncA).Pointer() == reflect.ValueOf(ts.FuncB).Pointer())
}
