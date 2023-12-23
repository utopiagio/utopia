package sysmetrics

import (
	"fmt"
	"log"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

// GetSystemMetrics constants
const (
	SM_CXSCREEN             = 0
	SM_CYSCREEN             = 1
	SM_CXVSCROLL            = 2
	SM_CYHSCROLL            = 3
	SM_CYCAPTION            = 4
	SM_CXBORDER             = 5
	SM_CYBORDER             = 6
	SM_CXDLGFRAME           = 7
	SM_CYDLGFRAME           = 8
	SM_CYVTHUMB             = 9
	SM_CXHTHUMB             = 10
	SM_CXICON               = 11
	SM_CYICON               = 12
	SM_CXCURSOR             = 13
	SM_CYCURSOR             = 14
	SM_CYMENU               = 15
	SM_CXFULLSCREEN         = 16
	SM_CYFULLSCREEN         = 17
	SM_CYKANJIWINDOW        = 18
	SM_MOUSEPRESENT         = 19
	SM_CYVSCROLL            = 20
	SM_CXHSCROLL            = 21
	SM_DEBUG                = 22
	SM_SWAPBUTTON           = 23
	SM_RESERVED1            = 24
	SM_RESERVED2            = 25
	SM_RESERVED3            = 26
	SM_RESERVED4            = 27
	SM_CXMIN                = 28
	SM_CYMIN                = 29
	SM_CXSIZE               = 30
	SM_CYSIZE               = 31
	SM_CXFRAME              = 32
	SM_CYFRAME              = 33
	SM_CXMINTRACK           = 34
	SM_CYMINTRACK           = 35
	SM_CXDOUBLECLK          = 36
	SM_CYDOUBLECLK          = 37
	SM_CXICONSPACING        = 38
	SM_CYICONSPACING        = 39
	SM_MENUDROPALIGNMENT    = 40
	SM_PENWINDOWS           = 41
	SM_DBCSENABLED          = 42
	SM_CMOUSEBUTTONS        = 43
	SM_CXFIXEDFRAME         = SM_CXDLGFRAME
	SM_CYFIXEDFRAME         = SM_CYDLGFRAME
	SM_CXSIZEFRAME          = SM_CXFRAME
	SM_CYSIZEFRAME          = SM_CYFRAME
	SM_SECURE               = 44
	SM_CXEDGE               = 45
	SM_CYEDGE               = 46
	SM_CXMINSPACING         = 47
	SM_CYMINSPACING         = 48
	SM_CXSMICON             = 49
	SM_CYSMICON             = 50
	SM_CYSMCAPTION          = 51
	SM_CXSMSIZE             = 52
	SM_CYSMSIZE             = 53
	SM_CXMENUSIZE           = 54
	SM_CYMENUSIZE           = 55
	SM_ARRANGE              = 56
	SM_CXMINIMIZED          = 57
	SM_CYMINIMIZED          = 58
	SM_CXMAXTRACK           = 59
	SM_CYMAXTRACK           = 60
	SM_CXMAXIMIZED          = 61
	SM_CYMAXIMIZED          = 62
	SM_NETWORK              = 63
	SM_CLEANBOOT            = 67
	SM_CXDRAG               = 68
	SM_CYDRAG               = 69
	SM_SHOWSOUNDS           = 70
	SM_CXMENUCHECK          = 71
	SM_CYMENUCHECK          = 72
	SM_SLOWMACHINE          = 73
	SM_MIDEASTENABLED       = 74
	SM_MOUSEWHEELPRESENT    = 75
	SM_XVIRTUALSCREEN       = 76
	SM_YVIRTUALSCREEN       = 77
	SM_CXVIRTUALSCREEN      = 78
	SM_CYVIRTUALSCREEN      = 79
	SM_CMONITORS            = 80
	SM_SAMEDISPLAYFORMAT    = 81
	SM_IMMENABLED           = 82
	SM_CXFOCUSBORDER        = 83
	SM_CYFOCUSBORDER        = 84
	SM_TABLETPC             = 86
	SM_MEDIACENTER          = 87
	SM_STARTER              = 88
	SM_SERVERR2             = 89
	SM_CMETRICS             = 91
	SM_REMOTESESSION        = 0x1000
	SM_SHUTTINGDOWN         = 0x2000
	SM_REMOTECONTROL        = 0x2001
	SM_CARETBLINKINGENABLED = 0x2002
)

// GetDeviceCaps index constants
const (
	DRIVERVERSION   = 0
	TECHNOLOGY      = 2
	HORZSIZE        = 4
	VERTSIZE        = 6
	HORZRES         = 8
	VERTRES         = 10
	LOGPIXELSX      = 88
	LOGPIXELSY      = 90
	BITSPIXEL       = 12
	PLANES          = 14
	NUMBRUSHES      = 16
	NUMPENS         = 18
	NUMFONTS        = 22
	NUMCOLORS       = 24
	NUMMARKERS      = 20
	ASPECTX         = 40
	ASPECTY         = 42
	ASPECTXY        = 44
	PDEVICESIZE     = 26
	CLIPCAPS        = 36
	SIZEPALETTE     = 104
	NUMRESERVED     = 106
	COLORRES        = 108
	PHYSICALWIDTH   = 110
	PHYSICALHEIGHT  = 111
	PHYSICALOFFSETX = 112
	PHYSICALOFFSETY = 113
	SCALINGFACTORX  = 114
	SCALINGFACTORY  = 115
	VREFRESH        = 116
	DESKTOPHORZRES  = 118
	DESKTOPVERTRES  = 117
	BLTALIGNMENT    = 119
	SHADEBLENDCAPS  = 120
	COLORMGMTCAPS   = 121
	RASTERCAPS      = 38
	CURVECAPS       = 28
	LINECAPS        = 30
	POLYGONALCAPS   = 32
	TEXTCAPS        = 34
)

const MONITOR_DEFAULTTOPRIMARY = 1
const MDT_EFFECTIVE_DPI = 0


var (
	user32 = syscall.NewLazySystemDLL("user32.dll")
	_GetDC            = user32.NewProc("GetDC")
	_GetMonitorInfo   = user32.NewProc("GetMonitorInfoW")
	_GetSystemMetrics = user32.NewProc("GetSystemMetrics")
	_MonitorFromPoint = user32.NewProc("MonitorFromPoint")
	_MonitorFromWindow = user32.NewProc("MonitorFromWindow")
	_ReleaseDC         = user32.NewProc("ReleaseDC")

	shcore = syscall.NewLazySystemDLL("shcore")
	_GetDpiForMonitor = shcore.NewProc("GetDpiForMonitor")

	gdi32 = syscall.NewLazySystemDLL("gdi32")
	_GetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
)

type Point struct {
	X, Y int32
}

type Rect struct {
	Left, Top, Right, Bottom int32
}

func GetDC(hwnd syscall.Handle) (syscall.Handle, error) {
	hdc, _, err := _GetDC.Call(uintptr(hwnd))
	if hdc == 0 {
		return 0, fmt.Errorf("GetDC failed: %v", err)
	}
	return syscall.Handle(hdc), nil
}


func GetSystemMetrics(nIndex int) int {
	r, _, _ := _GetSystemMetrics.Call(uintptr(nIndex))
	return int(r)
}

func GetMonitorInfo(hwnd syscall.Handle) GoMonitorInfo {
	var mi GoMonitorInfo
	mi.cbSize = uint32(unsafe.Sizeof(mi))
	v, _, _ := _MonitorFromWindow.Call(uintptr(hwnd), MONITOR_DEFAULTTOPRIMARY)
	_GetMonitorInfo.Call(v, uintptr(unsafe.Pointer(&mi)))
	return mi
}

// GetSystemDPI returns the effective DPI of the system.
func GetSystemDPI() int {
	// Check for GetDpiForMonitor, introduced in Windows 8.1.
	if _GetDpiForMonitor.Find() == nil {
		hmon := monitorFromPoint(Point{}, MONITOR_DEFAULTTOPRIMARY)
		return getDpiForMonitor(hmon, MDT_EFFECTIVE_DPI)
	} else {
		// Fall back to the physical device DPI.
		screenDC, err := GetDC(0)
		if err != nil {
			log.Println("GetSystemDPI() GetDC(0):", err)
			return 96
		}
		defer ReleaseDC(screenDC)
		return getDeviceCaps(screenDC, LOGPIXELSX)
	}
}

func ReleaseDC(hdc syscall.Handle) {
	_ReleaseDC.Call(uintptr(hdc))
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

func monitorFromPoint(pt Point, flags uint32) syscall.Handle {
	r, _, _ := _MonitorFromPoint.Call(uintptr(pt.X), uintptr(pt.Y), uintptr(flags))
	return syscall.Handle(r)
}

/*func GetSystemMetrics() *GoSystemMetrics {
	metrics := &GoSystemMetrics{}
	metrics.setSystemMetrics()
	return metrics
}*/

type GoMonitorInfo struct {
	cbSize   uint32
	Monitor  Rect
	Screen Rect
	Flags    uint32
}

/*func GetMonitorInfo() {
	mi := windows.GetMonitorInfo()
		x, y = mi.Monitor.Left, mi.Monitor.Top
		width = mi.Monitor.Right - mi.Monitor.Left
		height = mi.Monitor.Bottom - mi.Monitor.Top
}*/

type GoSystemMetrics struct {
	AspectX 	int
	AspectY 	int
	AspectXY 	int
	ClientHeight int
	ClientWidth int
	Height 		int
	Width 		int
	HorizRes 	int
	VertRes 	int
	HorizSize int
	VertSize	int

	//BorderWidth int
	//BorderHeight int
	//FixedBorderWidth int
	//FixedBorderHeight int
	//ResizeBorderWidth int
	//ResizeBorderHeight int

	//CaptionHeight int
	//CaptionButtonHeight int
	//CaptionButtonWidth int
	//SmallIconHeight int
	//SmallIconWidth int
	//IconHeight int
	//IconWidth int
	//ScrollbarWidth int
	//ScrollbarHeight int
	//MenuHeight int
	//GetDC(hWnd syscall.Handle) (syscall.Handle)
}

func (sm *GoSystemMetrics) init() {
	hScreen, _ := GetDC(0)

	sm.AspectX = getDeviceCaps(hScreen, ASPECTX)
	sm.AspectY = getDeviceCaps(hScreen, ASPECTY)

	sm.ClientHeight = GetSystemMetrics(SM_CYFULLSCREEN)
	sm.ClientWidth = GetSystemMetrics(SM_CXFULLSCREEN)
	sm.Height = GetSystemMetrics(SM_CYSCREEN)
	sm.Width = GetSystemMetrics(SM_CXSCREEN)

	sm.VertSize = getDeviceCaps(hScreen, VERTSIZE)
	sm.HorizSize = getDeviceCaps(hScreen, HORZSIZE)
	sm.VertRes = getDeviceCaps(hScreen, LOGPIXELSY)
	sm.HorizRes = getDeviceCaps(hScreen, LOGPIXELSX)

	//sm.borderWidth = w32.GetSystemMetrics(w32.SM_CXBORDER)
	//sm.borderHeight = w32.GetSystemMetrics(w32.SM_CYBORDER)
	//sm.fixedBorderWidth = w32.GetSystemMetrics(w32.SM_CXFIXEDFRAME)
	//sm.fixedBorderHeight = w32.GetSystemMetrics(w32.SM_CYFIXEDFRAME)
	//sm.resizeBorderWidth = w32.GetSystemMetrics(w32.SM_CXSIZEFRAME)
	//sm.resizeBorderHeight = w32.GetSystemMetrics(w32.SM_CYSIZEFRAME)
	//sm.captionHeight = w32.GetSystemMetrics(w32.SM_CYCAPTION)
	//sm.captionButtonHeight = w32.GetSystemMetrics(w32.SM_CYSIZE)
	//sm.captionButtonWidth = w32.GetSystemMetrics(w32.SM_CXSIZE)

	//sm.smallIconHeight = w32.GetSystemMetrics(w32.SM_CYSMICON)
	//sm.smallIconWidth = w32.GetSystemMetrics(w32.SM_CXSMICON)
	//sm.iconHeight = w32.GetSystemMetrics(w32.SM_CYICON)
	//sm.iconWidth = w32.GetSystemMetrics(w32.SM_CXICON)

	//sm.scrollbarWidth = w32.GetSystemMetrics(w32.SM_CXHSCROLL)
	//sm.scrollbarHeight = w32.GetSystemMetrics(w32.SM_CYHSCROLL)
	//sm.menuHeight = w32.GetSystemMetrics(w32.SM_CYMENU)
	ReleaseDC(hScreen)
	
}

func AspectX() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, ASPECTX)
	ReleaseDC(hScreen)
	return ret
}

func AspectY() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, ASPECTY)
	ReleaseDC(hScreen)
	return ret
}

func AspectXY() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, ASPECTXY)
	ReleaseDC(hScreen)
	return ret
}

func ClientHeight() int {
	return GetSystemMetrics(SM_CYFULLSCREEN)
}

func ClientWidth() int {
	return GetSystemMetrics(SM_CXFULLSCREEN)
}

func Height() int {
	return GetSystemMetrics(SM_CYSCREEN)
}

func HorizontalRes() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, LOGPIXELSX)
	ReleaseDC(hScreen)
	return ret
}

func HorizontalSize() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, HORZSIZE)
	ReleaseDC(hScreen)
	return ret
}

func VerticalRes() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, LOGPIXELSY)
	ReleaseDC(hScreen)
	return ret
}

func VerticalSize() int {
	hScreen, _ := GetDC(0)
	ret := getDeviceCaps(hScreen, VERTSIZE)
	ReleaseDC(hScreen)
	return ret
}

func Width() int {
	return GetSystemMetrics(SM_CXSCREEN)
}

/*func (sm *GoSystemMetrics) FixedBorderHeight() int {
	return sm.fixedBorderHeight
}

func (sm *GoSystemMetrics) FixedBorderWidth() int {
	return sm.fixedBorderWidth
}

func (sm *GoSystemMetrics) ResizeBorderHeight() int {
	return sm.resizeBorderHeight
}

func (sm *GoSystemMetrics) ResizeBorderWidth() int {
	return sm.resizeBorderWidth
}

func (sm *GoSystemMetrics) BorderHeight() int {
	return sm.borderHeight
}

func (sm *GoSystemMetrics) BorderWidth() int {
	return sm.borderWidth
}

func (sm *GoSystemMetrics) SmallIconWidth() int {
	return sm.smallIconWidth
}

func (sm *GoSystemMetrics) SmallIconHeight() int {
	return sm.smallIconHeight
}

func (sm *GoSystemMetrics) IconWidth() int {
	return sm.iconWidth
}

func (sm *GoSystemMetrics) IconHeight() int {
	return sm.iconHeight
}

func (sm *GoSystemMetrics) CaptionHeight() int {
	return sm.captionHeight
}

func (sm *GoSystemMetrics) CaptionButtonHeight() int {
	return sm.captionButtonHeight
}

func (sm *GoSystemMetrics) CaptionButtonWidth() int {
	return sm.captionButtonWidth
}

func (sm *GoSystemMetrics) MenuHeight() int {
	return sm.menuHeight
}

func (sm *GoSystemMetrics) ScrollbarHeight() int {
	return sm.scrollbarHeight
}

func (sm *GoSystemMetrics) ScrollbarWidth() int {
	return sm.scrollbarWidth
}

func (sm *GoSystemMetrics) WindowSizeX(clientWidth int, scrollbar bool) int {
	width := clientWidth + (sm.borderWidth * 2)
	if scrollbar == true {
		width += sm.scrollbarWidth * 2
	}
	return width
}

func (sm *GoSystemMetrics) WindowSizeY(clientHeight int, menu bool, scrollbar bool) int {
	height := clientHeight + (sm.borderWidth * 2)
	height += sm.captionHeight
	if menu {
		height += sm.menuHeight
	}
	if scrollbar {
		height += sm.scrollbarHeight * 2
	}
	return height
}*/