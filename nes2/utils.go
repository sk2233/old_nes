/*
@author: sk
@date: 2023/10/14
*/
package main

type Read func(addr uint16) uint8

func ReadUnit16(read Read, addrLo, addrHi uint16) uint16 {
	lo := read(addrLo)
	hi := read(addrHi)
	return uint16(lo) | (uint16(hi) << 8)
}
