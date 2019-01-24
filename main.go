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
	procIsWindowVisible = user32.MustFindProc("IsWindowVisible")
	procEnumChildWindows = user32.MustFindProc("EnumChildWindows")
	procGetClassName = user32.MustFindProc("GetClassNameW")
	procRegisterWindowMessage = user32.MustFindProc("RegisterWindowMessageW")
	procSendMessageTimeout = user32.MustFindProc("SendMessageTimeoutW")
	mshtml             = syscall.MustLoadDLL("mshtml.dll")
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

func IsWindowVisible(hwnd syscall.Handle) bool {
	ret, _, _ := syscall.Syscall(procIsWindowVisible.Addr(), 1, uintptr(hwnd), 0, 0)

	return ret != 0
}

func EnumChildWindows(hWndParent syscall.Handle, lpEnumFunc, lParam uintptr) bool {
	ret, _, _ := syscall.Syscall(procEnumChildWindows.Addr(), 3, uintptr(hWndParent), lpEnumFunc, lParam)

	return ret != 0
}

func GetClassName(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r1, _, e1 := syscall.Syscall(procGetClassName.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
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

func RegisterWindowMessage(lpString *uint16) uint32 {
	ret, _, _ := syscall.Syscall(procRegisterWindowMessage.Addr(), 1, uintptr(unsafe.Pointer(lpString)), 0, 0)

	return uint32(ret)
}


func SendMessageTimeout(
		hwnd syscall.Handle,
		msg uint32,
		w *uint16,
		l *uint16,
		flags, 
		timeout uint32,
result *uint32) uint32 {
	ret, _, _ := syscall.Syscall(procSendMessageTimeout.Addr(), 7,
	uintptr(hwnd),
	uintptr(msg),
	uintptr(unsafe.Pointer(w)), 
	uintptr(unsafe.Pointer(l)), 
	flags,
	uintptr(timeout),
	uintptr(unsafe.Pointer(result)))

	return uint32(ret)

}

func get_ihtml(hwnd syscall.Handle, ppHTMLDoc2 **IHTMLDocument2){
	lRes :=	0
	
	
	var MSG uint32 = RegisterWindowMessage(_T("WM_HTML_GETOBJECT"))
	SendMessageTimeout(hWnd, MSG, 0, 0, SMTO_ABORTIFHUNG, 1000, &lRes)
	
	hr	:=	(*pfObjectFromLresult)(lRes,IID_IHTMLDocument2,0, **ppHTMLDoc2)

	return hr
}

func get_ihtml_url(pHTMLDoc2 *mshtml.IHTMLDocument2, pUrl *string) string {
	if(pHTMLDoc2==nil){
		return syscall.E_INVALIDARG
	}
	hr := pHTMLDoc2.get_URL(pUrl)
	return hr	
}

var ihtml IHTMLDocument2



func FindWindow() {
	var hwnd syscall.Handle

	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {  //callback function


	buf := make([]char16, 32)
	var url string;
	if(!h){
		return 0;
	}
	GetClassName(h,buf,32)

	if(strstr(buf,"Explorer_Server"))
	{
		get_ihtml(hwnd,&ihtml)
		get_ihtml_url(ihtml,&url)

		if(IsWindowVisible(h)){
			b := make([]uint16, 200)
			_, err := GetWindowText(h, &b[0], int32(len(b)))


			if err != nil {
				// ignore the error
				return 1 // continue enumeration
			}
			
			hwnd = h
			
		
			title := syscall.UTF16ToString(b)
			fmt.Printf("Found '%s' window: handle=0x%x\n", title, hwnd)
		}

	}
		return 1 // continue enumeration	
	})

	
	EnumWindows(cb, 0)
	return 
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
