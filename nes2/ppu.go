/*
@author: sk
@date: 2023/10/14
*/
package main

type Ppu struct {
	bus *Bus
	// 64个精灵的属性，每个精灵使用4b
	// 0:Y 单位像素
	// 1:TileIndex 0x2000控制是8*8模式还是8*16模式,第一个Tile都是索引都是TileIndex,第二个Tile索引是TileIndex+1(只有8*16有两个,这要求特殊的tile图片设计格式)
	// 2:精灵属性   0~1位:使用的调色板  2~4位:未使用  5位:精灵是否在背景后面  6位:水平翻转  7位:垂直翻转
	// 3:X 单位像素
	oamRam [256]byte
	// 0x2000  0~1位:选择哪个nameTab  2位:0增长1水平移动,1增长32纵向移动  3位:精灵使用哪个tileSet  4位:背景使用哪个tileSet  5位:精灵大小0:8*8,1:8*16  7位:在V_Blank开始时产生NMI
	ctrl uint8
	// 0位:是否显示为黑白  3位:是否渲染背景  4位:是否渲染精灵
	mask uint8
	// 5位:精灵是否溢出  6位:  7位:是否处于V_Blank
	status uint8
	// 确定开始画面显示的位置,第一次写入x偏移,第二次写入y偏移
	scroll uint8
	// 0~4位:xTileIndex  5~9位:yTileIndex  10~11位:选择哪个nameTab  12~14位:y在tile内的像素偏移
	// 对t+1水平移动(xTileIndex),对t+32纵向移动(yTileIndex)
	t uint16
	// 0~2位:x在tile内的像素偏移
	x uint8
	// 写地址时需要写两次,使用这个开关控制,例如写入偏移时先写xTileIndex与x,并设置w=true,再写入yTileIndex与t的12~14位
	w bool
}

func NewPpu(bus *Bus) *Ppu {
	return &Ppu{bus: bus}
}
