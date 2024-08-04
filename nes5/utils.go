/*
@author: sk
@date: 2024/2/24
*/
package nes5

import (
	"image"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func AddrEq(obj1, obj2 any) bool {
	return reflect.ValueOf(obj1).Pointer() == reflect.ValueOf(obj2).Pointer()
}

func Format[T uint16 | uint8](value T, count int) string {
	temp := strconv.FormatUint(uint64(value), 16)
	temp = strings.Repeat("0", count-len(temp)) + temp
	return strings.ToUpper(temp)
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Min[T int](val1, val2 T) T {
	if val1 < val2 {
		return val1
	}
	return val2
}

func OpenFile(path string) *os.File {
	open, err := os.Open(path)
	HandleErr(err)
	return open
}

func ReadStr(reader *os.File, count int) string {
	bs := ReadBytes(reader, count)
	return string(bs)
}

func ReadUint8(reader *os.File) uint8 {
	return ReadBytes(reader, 1)[0]
}

func ReadBytes(reader *os.File, count int) []byte {
	bs := make([]byte, count)
	_, err := reader.Read(bs)
	HandleErr(err)
	return bs
}

var (
	imgOption = &ebiten.DrawImageOptions{}
)

func DrawImage(screen *ebiten.Image, img *ebiten.Image, x, y float64) {
	imgOption.GeoM.Reset()
	imgOption.GeoM.Translate(x, y)
	screen.DrawImage(img, imgOption)
}

func DrawSubImage(screen *ebiten.Image, img *ebiten.Image, x, y float64, ix, iy, iw, ih int) {
	img = img.SubImage(image.Rect(ix, iy, ix+iw, iy+ih)).(*ebiten.Image)
	imgOption.GeoM.Reset()
	imgOption.GeoM.Translate(x, y)
	screen.DrawImage(img, imgOption)
}

func GetRangeBit8(data uint8, from, to int) uint8 {
	data <<= 8 - to
	data >>= 8 - to + from
	return data
}

func SetRangeBit8(data uint8, from, to int) uint8 {
	empty := 8 - (to - from) // 保证其他空间都是empty
	data <<= empty
	data >>= 8 - to
	return data
}

func GetRangeBit16(data uint16, from, to int) uint16 {
	data <<= 16 - to
	data >>= 16 - to + from
	return data
}

func SetRangeBit16(data uint16, from, to int) uint16 {
	empty := 16 - (to - from) // 保证其他空间都是empty
	data <<= empty
	data >>= 16 - to
	return data
}

func GetBit[T uint8 | uint16](data T, index int) bool {
	return (data & (1 << index)) > 0
}

func SetBit[T uint8 | uint16](data bool, index int) T {
	if !data {
		return 0
	}
	return 1 << index
}

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
