/*
@author: sk
@date: 2023/10/9
*/
package main

type CpuState uint8

const (
	CpuStateC CpuState = 1 << iota // Carry  进位 在 加法运算时需要额外加上
	CpuStateZ                      // Zero
	CpuStateI                      // Disable
	CpuStateD                      // Decimal
	CpuStateB                      // Break
	CpuStateU                      // Unused
	CpuStateV                      // Overflow
	CpuStateN                      // Negative
)
