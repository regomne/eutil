package peHelper

import (
	"debug/pe"
	"os"
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

func ceilAlign(val uint32, align uint32) uint32 {
	l := val % align
	if l == 0 {
		return val
	}
	return (val/align + 1) * align
}

// LoadPEImage load a PE file to memory
func LoadPEImage(fname string) (image []byte, err error) {
	var peFile *pe.File
	peFile, err = pe.Open(fname)
	if err != nil {
		return
	}
	defer peFile.Close()
	var fs *os.File
	fs, err = os.Open(fname)
	if err != nil {
		return
	}
	defer fs.Close()
	lastSec := peFile.Sections[peFile.NumberOfSections-1]
	imageSize := lastSec.VirtualAddress + ceilAlign(lastSec.VirtualSize, 0x1000)
	image = make([]byte, imageSize)
	_, err = fs.Read(image[0:peFile.Sections[0].Offset])
	if err != nil {
		return
	}
	for _, sec := range peFile.Sections {
		fs.Seek(int64(sec.Offset), 0)
		_, err = fs.Read(image[sec.VirtualAddress : sec.VirtualAddress+sec.Size])
		if err != nil {
			return
		}
	}
	return
}
