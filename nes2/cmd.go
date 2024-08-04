/*
@author: sk
@date: 2023/10/14
*/
package main

type Cmd struct {
	AddrCode uint8
	OpCode   uint8
}

func loadCmds() map[uint8]*Cmd {
	res := make(map[uint8]*Cmd)

	return res
}
