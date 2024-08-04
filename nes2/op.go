/*
@author: sk
@date: 2023/10/14
*/
package main

type OpFunc func(cpu *Cpu) uint8

func loadOps() map[uint8]OpFunc {
	res := make(map[uint8]OpFunc)
	return res
}
