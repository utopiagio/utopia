// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/dekstop_windows.go */

package desktop

// DeviceCaps are obtained for the primary screen only. See monitor package for other options.
import (
	//"log"
	"unsafe"
	syscall "golang.org/x/sys/windows"
	//ui "github.com/utopiagio/utopia"
	"github.com/utopiagio/utopia/internal/sysmetrics"
)

const (
	SWP_FRAMECHANGED  = 0x0020
	SWP_NOMOVE        = 0x0002
	SWP_NOOWNERZORDER = 0x0200
	SWP_NOSIZE        = 0x0001
	SWP_NOZORDER      = 0x0004
	SWP_SHOWWINDOW    = 0x0040
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	_GetClientRect = user32.NewProc("GetClientRect")
	_GetDesktopWindow = user32.NewProc("GetDesktopWindow")
	_GetWindowRect = user32.NewProc("GetWindowRect")
	_SetWindowPos = user32.NewProc("SetWindowPos")
)

type Rect struct {
	Left, Top, Right, Bottom int32
}

// it will panic when the function fails
func GetHWnd() syscall.HWND {
	ret, _, _ := _GetDesktopWindow.Call()
	return syscall.HWND(ret)
}

func getClientRect(hwnd syscall.Handle) Rect {
	var r Rect
	_GetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r
}

func getWindowRect(hwnd syscall.Handle) Rect {
	var r Rect
	_GetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	//log.Println("GetWindowRect : ", r)
	return r
}

func setWindowPos(hwnd syscall.Handle, hwndInsertAfter uint32, x, y, dx, dy int32, style uintptr) {
	_SetWindowPos.Call(uintptr(hwnd), uintptr(hwndInsertAfter),
		uintptr(x), uintptr(y),
		uintptr(dx), uintptr(dy),
		style,
	)
}

type goDeviceCaps struct {
	aspectX 	int
	aspectY 	int
	aspectXY 	int
	clientHeight int
	clientWidth int
	height 		int
	width 		int
	horizRes 	int
	vertRes 	int
	horizSize 	int
	vertSize	int
}

var deviceCaps goDeviceCaps

func Init() {
	deviceCaps.getDeviceCaps()
}

/*func GoDeskTop() (hdesktop *GoDeskTopObj) {
	hWnd := GetDesktopWindow()
	screen := GioWidget{
		className: 	"",
		windowName: "",
		style: 		0,
		exStyle:    0,
		state:      0,
		x:      	0,
		y:      	0,
		width:  	0,
		height: 	0,
		hWndParent:	0,
		id:			0,
		instance:	0,
		param:		0,
		hWnd: 		0,
		parent:  	nil,
		disabled: 	false,
		visible:   	true,
		cursor:     nil,	// *goCursor
		font:       nil, 	// *goFont
		text: 	 	"",
		window:     false,
		widgets:    map[int]*goWidget{},
		alpha:		0, 		//uint8
		backcolor:		0,
		forecolor:		0,
		onClose:       nil, 	//func()
		onCanClose:    nil, 	//func() bool
		onMouseMove:   nil, 	//func(x, y int)
		onMouseWheel:  nil, 	//func(x, y int, delta float64)
		onMouseDown:   nil, 	//func(button MouseButton, x, y int)
		onMouseUp:     nil, 	//func(button MouseButton, x, y int)
		onKeyDown:     nil, 	//func(key int)
		onKeyUp:       nil, 	//func(key int)
		onResize:      nil, 	//func()
	}
	//screen := goWidget{"", "", 0, 0, 0, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 0, 0, 0, 0, 0, nil, false, false, nil, map[int]*goWidget{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
	screen := GioWidget{}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	devicecaps := GoDeviceCaps{}
	hdesktop = &GoDeskTopObj{object,screen, devicecaps}
	//hdesktop.hWnd = hWnd
	//hdesktop.instance = goApp.hInstance
	hdesktop.getDeviceCaps()
	return hdesktop
}*/

/*type GoDeskTopObj struct {
	GioObject
	GioWidget
	goDeviceCaps		
}*/

//func (d *GoDeskTopWindow) GetHWnd() (w32.HWND) {
//	return d.hWnd
//}

/*func (ob *GoDeskTopObj) wid() (*ui.GioWidget) {
	return &ob.GioWidget
}

func (ob *GoDeskTopObj) Screen() (*GoDeviceCaps) {
	return &ob.GoDeviceCaps
}*/

func (ob goDeviceCaps) getDeviceCaps() {
	//hWnd := ob.goWidget.hWnd
	//hDC := sysmetrics.GetDC(hWnd)
	deviceCaps.aspectX = sysmetrics.AspectX()
	deviceCaps.aspectY = sysmetrics.AspectY()
	deviceCaps.aspectXY = sysmetrics.AspectXY()

	deviceCaps.clientHeight = sysmetrics.ClientHeight()
	deviceCaps.clientWidth = sysmetrics.ClientWidth()
	deviceCaps.height = sysmetrics.Height()
	deviceCaps.width = sysmetrics.Width()
	deviceCaps.vertSize = sysmetrics.VerticalSize() //hDC)		//GetDeviceCaps(hDC, w32.VERTSIZE)
	deviceCaps.horizSize = sysmetrics.HorizontalSize() //hDC) 	//GetDeviceCaps(hDC, w32.HORZSIZE)
	deviceCaps.vertRes = sysmetrics.VerticalRes() //hDC) 		//GetDeviceCaps(hDC, w32.LOGPIXELSY)
	deviceCaps.horizRes = sysmetrics.HorizontalRes() //hDC) 	//GetDeviceCaps(hDC, w32.LOGPIXELSX)

}

func AspectX() int {
	return deviceCaps.aspectX
}

func AspectY() int {
	return deviceCaps.aspectY
}

func AspectXY() int {
	return deviceCaps.aspectXY
}

func ClientHeight() int {
	return deviceCaps.clientHeight
}

func ClientWidth() int {
	return deviceCaps.clientWidth
}

func GetClientRect(hWnd syscall.Handle) (x, y, width, height int) {
	rc := getClientRect(hWnd)
	return int(rc.Left), int(rc.Top), int(rc.Right - rc.Left), int(rc.Bottom - rc.Top)
}

func GetWindowRect(hWnd syscall.Handle) (x, y, width, height int) {
	rc := getWindowRect(hWnd)
	return int(rc.Left), int(rc.Top), int(rc.Right - rc.Left), int(rc.Bottom - rc.Top)
}

func Height() int {
	return deviceCaps.height
}

func HorizontalRes() int {
	return deviceCaps.horizRes
}

func HorizontalSize() int {
	return deviceCaps.horizSize
}

func SetWindowPos(hWnd syscall.Handle, hwndInsertAfter int, x int, y int, width int, height int, style uintptr) {
	var swpStyle uintptr
	if style == 0 {
		swpStyle = uintptr(SWP_NOZORDER | SWP_FRAMECHANGED)
	} else {
		swpStyle = style
	}
	setWindowPos(hWnd, 0, int32(x), int32(y), int32(width), int32(height), swpStyle)
}

func VerticalRes() int {
	return deviceCaps.vertRes
}

func VerticalSize() int {
	return deviceCaps.vertSize
}

func Width() int {
	return deviceCaps.width
}
