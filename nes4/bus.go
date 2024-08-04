/*
@author: sk
@date: 2023/10/26
*/
package nes4

type Bus struct {
	Cpu                       *Cpu
	Ppu                       *Ppu
	Cartridge                 *Cartridge
	CpuRam                    [2048]uint8 // 2K 内存
	Controller                [2]uint8    // 控制位
	ControllerCache           [2]uint8    // Controller 的Cache
	SystemClockCounter        uint32      // 系统时钟
	DmaPage, DmaAddr, DmaData uint8       // Dma相关变量
	// DmaDummy 虚拟的Dma时钟周期， DmaTransfer DMA是否发生
	DmaDummy, DmaTransfer bool
}

func NewBus() *Bus {
	cpu := NewCpu()
	ppu := NewPpu()
	res := &Bus{Cpu: cpu, Ppu: ppu, DmaDummy: true}
	cpu.ConnectBus(res)
	return res
}

func (b *Bus) CpuWrite(addr uint16, data uint8) {
	if b.Cartridge.CpuWrite(addr, data) {
		// 卡带扩展能力，若是被卡带拦截，这里不做什么
	} else if addr < 0x1FFF {
		// 映射为了8k，但实际只有2k通过镜像分配
		b.CpuRam[addr%0x07FF] = data
	} else if addr < 0x3FFF {
		// 写入到ppu的地址空间，实际只有8b其他都是镜像
		b.Ppu.CpuWrite(addr%0x0007, data)
	} else if addr == 0x4014 {
		// 初始化DMA进行 DMA数据传输
		b.DmaPage = data
		b.DmaAddr = 0x00
		b.DmaTransfer = true
	} else if addr == 0x4016 { // 控制写入缓存
		b.ControllerCache[0] = b.Controller[0]
	} else if addr == 0x4017 {
		b.ControllerCache[1] = b.Controller[1]
	}
}

func (b *Bus) CpuRead(addr uint16) uint8 {
	if data, ok := b.Cartridge.CpuRead(addr); ok {
		// 还是先由卡带拦截
		return data
	} else if addr < 0x1FFF {
		return b.CpuRam[addr%0x07FF]
	} else if addr < 0x3FFF {
		return b.Ppu.CpuRead(addr % 0x0007)
	} else if addr == 0x4016 {
		data = (b.ControllerCache[0] & 0x80) >> 7 // 读取最高位 每次读取 1b
		b.ControllerCache[0] <<= 1                // 最高位读取完毕，移除最高位
		return data
	} else if addr == 0x4017 {
		data = (b.ControllerCache[1] & 0x80) >> 7 // 读取最高位 每次读取 1b
		b.ControllerCache[1] <<= 1                // 最高位读取完毕，移除最高位
		return data
	}
	return 0x00
}

func (b *Bus) SetCartridge(cartridge *Cartridge) { // 插入卡带
	b.Cartridge = cartridge
	b.Ppu.ConnectCartridge(cartridge)
}

func (b *Bus) Reset() {
	b.Cartridge.Reset()
	b.Cpu.Reset()
	b.Ppu.Reset()
	b.SystemClockCounter = 0
	b.DmaPage = 0x00
	b.DmaAddr = 0x00
	b.DmaData = 0x00
	b.DmaDummy = true
	b.DmaTransfer = false
}

func (b *Bus) Clock() {
	// 每次都运动
	b.Ppu.Clock()
	// CPU没3次执行一次
	if b.SystemClockCounter%3 == 0 {
		// DMA操作
		if b.DmaTransfer {
			// 模拟等待 因为模拟器太快?
			if b.DmaDummy {
				b.DmaDummy = b.SystemClockCounter%2 == 0
			} else {
				if b.SystemClockCounter%2 == 0 {
					// 读取数据，page与addr共同组成地址
					b.DmaData = b.CpuRead(uint16(b.DmaPage)<<8 | uint16(b.DmaAddr))
				} else {
					// 写入数据到ppu的 OAM内
					b.Ppu.OAM[b.DmaAddr/4].SetData(b.DmaAddr%4, b.DmaData)
					b.DmaAddr++ // 发生环绕证明 dam已经写完数据了，dam结束
					if b.DmaAddr == 0x00 {
						b.DmaTransfer = false
						b.DmaDummy = true
					}
				}
			}
		} else {
			// 正常执行CPU行为
			b.Cpu.Clock()
		}
	}
	// ppu发生中断 通知 Cpu
	if b.Ppu.Nmi {
		b.Ppu.Nmi = false
		b.Cpu.Nmi()
	}
	b.SystemClockCounter++
}
