package main

import (
	"fmt"
	_"log"
	"syscall"
	"unsafe"
)


var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow")
)




func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r1, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r1)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}


func GetForegroundWindow() syscall.Handle {
	ret, _, _ := syscall.Syscall(procGetForegroundWindow.Addr(), 0, 0, 0, 0)
		
	return syscall.Handle(ret)
}





 
func main() {
	
	var hwnd syscall.Handle
	

	hwnd = GetForegroundWindow()
	

		b := make([]uint16, 200)
		_, err := GetWindowText(hwnd, &b[0], int32(len(b)))

		if err != nil {
			
		}
		
		fmt.Printf("Found '%s' window: handle=0x%x\n", syscall.UTF16ToString(b), hwnd)
	
	

	


}
