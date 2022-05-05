//go:build (386 || amd64 || amd64p32)

#include "textflag.h"

// func cpuid(eax, ecx uint32) [4]uint32
TEXT ·cpuid(SB),NOSPLIT,$0-24
	MOVL eax+0(FP), AX
	MOVL ecx+4(FP), CX
	CPUID
	MOVL AX, eax+8(FP)
	MOVL BX, ebx+12(FP)
	MOVL CX, ecx+16(FP)
	MOVL DX, edx+20(FP)
	RET

// func xgetbv(arg0 uint32) [2]uint32
TEXT ·xgetbv(SB),NOSPLIT,$0-16
	MOVL arg0+0(FP), CX
	XGETBV
	MOVL AX, eax+8(FP)
	MOVL DX, edx+12(FP)
	RET
