/*
@author: sk
@date: 2023/10/14
*/
package main

type AddrFunc func(cpu *Cpu) uint8

func loadAddrs() map[uint8]AddrFunc {
	res := make(map[uint8]AddrFunc)
	// 寻址实现
	return res
}
