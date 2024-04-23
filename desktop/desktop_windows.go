// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/dekstop_windows.go */

package desktop

// DeviceCaps are obtained for the primary screen only. See monitor package for other options.
import (
	//"log"
	"fmt"
	"unsafe"
	syscall "golang.org/x/sys/windows"
	"github.com/utopiagio/utopia/internal/sysmetrics"
)

const (

	LOGPIXELSX = 88

	MDT_EFFECTIVE_DPI = 0

	MONITOR_DEFAULTTOPRIMARY = 1

	SM_CXSIZEFRAME = 32
	SM_CYSIZEFRAME = 33

	SWP_FRAMECHANGED  = 0x0020
	SWP_NOMOVE        = 0x0002
	SWP_NOOWNERZORDER = 0x0200
	SWP_NOSIZE        = 0x0001
	SWP_NOZORDER      = 0x0004
	SWP_SHOWWINDOW    = 0x0040
)

type PROCESS_DPI_AWARENESS int
type DPI_AWARENESS int

const (
    PROCESS_DPI_UNAWARE = 0
    PROCESS_SYSTEM_DPI_AWARE = 1
    PROCESS_PER_MONITOR_DPI_AWARE = 2
)

const  (
    DPI_AWARENESS_INVALID DPI_AWARENESS = -1
    DPI_AWARENESS_UNAWARE DPI_AWARENESS = 0
    DPI_AWARENESS_SYSTEM_AWARE DPI_AWARENESS = 1
    DPI_AWARENESS_PER_MONITOR_AWARE DPI_AWARENESS = 2
)

const (
	DPI_AWARENESS_CONTEXT_UNAWARE              = -1
	DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = -2
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = -3
	DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = -4
	DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = -5
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	_GetClientRect = user32.NewProc("GetClientRect")
	_GetDC = user32.NewProc("GetDC")
	_GetDesktopWindow = user32.NewProc("GetDesktopWindow")
	_GetDpiForWindow = user32.NewProc("GetDpiForWindow")
	_GetSystemMetrics = user32.NewProc("GetSystemMetrics")
	_GetWindowRect = user32.NewProc("GetWindowRect")
	_MonitorFromPoint = user32.NewProc("MonitorFromPoint")
	_ReleaseDC = user32.NewProc("ReleaseDC")
	_SetProcessDPIAware = user32.NewProc("SetProcessDPIAware")
	_SetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	_SetWindowPos = user32.NewProc("SetWindowPos")
	
	shcore = syscall.NewLazyDLL("shcore.dll")
	_GetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")
	_GetScaleFactorForMonitor = shcore.NewProc("GetScaleFactorForMonitor")
	
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
	_GetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
)

type Point struct {
	X, Y int32
}

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

func Init() (dpr float32, spr float32) {
	deviceCaps.getDeviceCaps()
	//ok := SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_UNAWARE)
	//ok := SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_SYSTEM_AWARE)
	SetProcessDPIAware()
	//log.Println("OK =", ok)
	dpi := GetSystemDPI()
	dpr, spr = configureForDPI(dpi)
	return
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

func TaskBarHeight() int {
	return deviceCaps.height - deviceCaps.clientHeight
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
	//return GetSystemMetrics(nIndex int)
	return deviceCaps.vertRes
}

func VerticalSize() int {
	return deviceCaps.vertSize
}

func Width() int {
	return deviceCaps.width
}

func GetDC(hwnd syscall.Handle) (syscall.Handle, error) {
	hdc, _, err := _GetDC.Call(uintptr(hwnd))
	if hdc == 0 {
		return 0, fmt.Errorf("GetDC failed: %v", err)
	}
	return syscall.Handle(hdc), nil
}

func ReleaseDC(hdc syscall.Handle) {
	_ReleaseDC.Call(uintptr(hdc))
}

func GetSystemScale() (float32, float32) {
	// Check for GetScaleFactorForMonitor, introduced in Windows 8.1.
	hmon := monitorFromPoint(Point{}, MONITOR_DEFAULTTOPRIMARY)
	scaler := getScaleForMonitor(hmon)
	return scaler, scaler
}

// GetSystemDPI returns the effective DPI of the system.
func GetSystemDPI() int {
	// Check for GetDpiForMonitor, introduced in Windows 8.1.
	if _GetDpiForMonitor.Find() == nil {
		hmon := monitorFromPoint(Point{}, MONITOR_DEFAULTTOPRIMARY)
		dpi := getDpiForMonitor(hmon, MDT_EFFECTIVE_DPI)
		return dpi
	} else {
		// Fall back to the physical device DPI.
		screenDC, err := GetDC(0)
		if err != nil {
			return 96
		}
		defer ReleaseDC(screenDC)
		return getDeviceCaps(screenDC, LOGPIXELSX)
	}
}

func GetSystemMetrics(nIndex int) int {
	r, _, _ := _GetSystemMetrics.Call(uintptr(nIndex))
	return int(r)
}

func SetProcessDPIAware() int {
	r, _, _ := _SetProcessDPIAware.Call()
	return int(r)
}

func SetProcessDpiAwarenessContext(aware int) int {
	r, _, _ := _SetProcessDpiAwarenessContext.Call(uintptr(aware))
	return int(r)
}

func configureForDPI(dpi int) (float32, float32) {
	const inchPrDp = 1.0 / 96.0
	ppdp := float32(dpi) * inchPrDp
	return ppdp, ppdp
}

func getDeviceCaps(hdc syscall.Handle, index int32) int {
	c, _, _ := _GetDeviceCaps.Call(uintptr(hdc), uintptr(index))
	return int(c)
}

func getDpiForMonitor(hmonitor syscall.Handle, dpiType uint32) int {
	var dpiX, dpiY uintptr
	_GetDpiForMonitor.Call(uintptr(hmonitor), uintptr(dpiType), uintptr(unsafe.Pointer(&dpiX)), uintptr(unsafe.Pointer(&dpiY)))
	return int(dpiX)
}

func getScaleForMonitor(hmonitor syscall.Handle) float32 {
	var scalef uintptr
	_GetScaleFactorForMonitor.Call(uintptr(hmonitor), uintptr(unsafe.Pointer(&scalef)))
	return float32(scalef) / 100
}

func monitorFromPoint(pt Point, flags uint32) syscall.Handle {
	r, _, _ := _MonitorFromPoint.Call(uintptr(pt.X), uintptr(pt.Y), uintptr(flags))
	return syscall.Handle(r)
}