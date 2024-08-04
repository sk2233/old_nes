/*
@author: sk
@date: 2023/10/9
*/
package main

func main() {
	//var num uint8
	//num = 0xff
	//fmt.Println(^num)
	cartridge := NewCartridge("")
	bus := NewBus()
	bus.SetCartridge(cartridge)
	bus.Reset()
	bus.Clock()
	for true {
		// 绘制一帧画面
		for !bus.Ppu.FrameComplete {
			bus.Clock()
		}
		bus.Ppu.FrameComplete = false
	}
}
