/*
@author: sk
@date: 2023/10/14
*/
package main

const (
	K = 1024
)

const (
	StateC = 1 << iota // Carry  进位 在 加法运算时需要额外加上
	StateZ             // Zero 最近指令结果是否为0
	StateI             // Disable 是否忽略普通中断
	StateD             // Decimal  nes不使用
	StateB             // Break
	StateU             // Unused
	StateV             // Overflow 是否计算溢出
	StateN             // Negative 计算结果的最高位
)
