package frame

import (
	"syscall"
	"unsafe"
)

type (
	short int16
	word  uint16
	coord struct {
		x short
		y short
	}
	smallRect struct {
		left   short
		top    short
		right  short
		bottom short
	}
	consoleScreenBufferInfo struct {
		size              coord
		cursorPosition    coord
		attributes        word
		window            smallRect
		maximumWindowSize coord
	}
)

func (this coord) uintptr() uintptr {
	return uintptr(*(*int32)(unsafe.Pointer(&this)))
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var procSetConsoleCursorPosition = kernel32.NewProc("SetConsoleCursorPosition")
var procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
var tmpInfo consoleScreenBufferInfo

func setCursorRow(row int) (err error) {
	pos := coord{0, short(row)}
	r0, _, e1 := syscall.Syscall(procSetConsoleCursorPosition.Addr(), 2, uintptr(outHandle), pos.uintptr(), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func getConsoleScreenBufferInfo(h syscall.Handle, info *consoleScreenBufferInfo) (err error) {
	r0, _, e1 := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func GetCursorRow() (int, error) {
	err := getConsoleScreenBufferInfo(outHandle, &tmpInfo)
	if err != nil {
		return -1, err
	}
	return int(tmpInfo.cursorPosition.y), nil
}
