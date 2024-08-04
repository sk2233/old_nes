/*
@author: sk
@date: 2023/10/28
*/
package nes4

import (
	"os"
	"strconv"
	"strings"
)

func OpenFile(path string) *os.File {
	open, err := os.Open(path)
	HandleErr(err)
	return open
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
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

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func Max[T int](value1, value2 T) T {
	if value1 > value2 {
		return value1
	}
	return value2
}

func Format[T uint16 | uint8](value T, count int) string {
	temp := strconv.FormatUint(uint64(value), 16)
	temp = strings.Repeat("0", count-len(temp)) + temp
	return strings.ToUpper(temp)
}
