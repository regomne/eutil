package peHelper

import (
	"debug/pe"
)

// RvaToOffset convert pe rva to file offset
func RvaToOffset(peFile *pe.File, rva uint32) (off uint32, secIdx int) {
	off = 0xffffffff
	secIdx = 1
	secs := peFile.Sections
	if rva < secs[0].VirtualAddress {
		if rva < secs[0].Offset {
			return rva, -1
		}
		return
	}
	for idx, sec := range secs {
		if rva >= sec.VirtualAddress &&
			rva < sec.VirtualAddress+sec.VirtualSize {
			secOff := rva - sec.VirtualAddress
			if secOff <= sec.Size {
				return sec.Offset + (rva - sec.VirtualAddress), idx
			} else {
				return
			}
		}
	}
	return
}

//OffsetToRva conver pe file offset to rva
func OffsetToRva(peFile *pe.File, off uint32) (rva uint32, secIdx int) {
	secs := peFile.Sections
	if off < secs[0].Offset {
		return off, -1
	}
	for idx, sec := range secs {
		if off >= sec.Offset &&
			off < sec.Offset+sec.Size {
			return sec.VirtualAddress + (off - sec.Offset), idx
		}
	}
	return 0xffffffff, -1
}
