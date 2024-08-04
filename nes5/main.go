/*
@author: sk
@date: 2024/2/24
*/
package nes5

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	//err := ebiten.RunGame(NewCpuApp())
	err := ebiten.RunGame(NewCartridgeApp())
	//err := ebiten.RunGame(NewBackgroundApp())
	HandleErr(err)
}
