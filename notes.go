func FindWindow() (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		
		hwnd = h
		
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	
	return hwnd, nil
}


//******************************************************************************** круче

func FindWindow() {
	var hwnd syscall.Handle

	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {  //callback function

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
		return 1 // continue enumeration	
	})

	
	EnumWindows(cb, 0)
	return 
}