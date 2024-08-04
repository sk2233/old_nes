/*
@author: sk
@date: 2023/10/14
*/
package main

import "fmt"

func main() {
	bus := NewBus()
	cpu := NewCpu(bus)
	ppu := NewPpu(bus)
	fmt.Println(cpu)
	fmt.Println(ppu)
}
