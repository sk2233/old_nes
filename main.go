/*
@author: sk
@date: 2024/5/30
*/
package main

import (
	"nes/nes5"

	"github.com/hajimehoshi/ebiten/v2"
)

// https://bugzmanov.github.io/nes_ebook/chapter_1.html

func main() {
	err := ebiten.RunGame(nes5.NewCartridgeApp())
	nes5.HandleErr(err)
}
