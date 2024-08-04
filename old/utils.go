/*
@author: sk
@date: 2023/10/12
*/
package main

import "os"

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadBytes(file *os.File, count int) []byte {
	bs := make([]byte, count)
	_, err := file.Read(bs)
	HandleErr(err)
	return bs
}

func ReadStr(file *os.File, count int) string {
	bs := ReadBytes(file, count)
	return string(bs)
}

func ReadUint8(file *os.File) uint8 {
	bs := ReadBytes(file, 1)
	return bs[0]
}

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
