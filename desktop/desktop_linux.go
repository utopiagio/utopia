// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/dekstop_linux.go */

package desktop

// DeviceCaps are obtained for the primary screen only. See monitor package for other options.

/*
#cgo LDFLAGS: -lX11

#include <X11/Xlib.h>
*/
import "C"

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
	deviceCaps.getPrimaryDeviceCaps()
}

func (ob goDeviceCaps) getPrimaryDeviceCaps() {
	dy := C.XOpenDisplay(nil)
	if dy == nil {
		return
		log.Println("return nil")
	}
	defer C.XCloseDisplay(dy)

	screen := C.XDefaultScreenOfDisplay(dy)
	width := int(C.XWidthOfScreen(screen))
	height := int(C.XHeightOfScreen(screen))
	deviceCaps.height = height
	deviceCaps.width = width
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

/*func (ob goDeviceCaps) getDeviceCaps() {
	 Display *display;
    Screen *screen;

    // open a display
    display = XOpenDisplay(NULL);

    // return the number of available screens
    int count_screens = ScreenCount(display);

    printf("Total count screens: %d\n", count_screens);


    for (int i = 0; i < count_screens; ++i) {
        screen = ScreenOfDisplay(display, i);
        printf("\tScreen %d: %dX%d\n", i + 1, screen->width, screen->height);
    }

    // close the display
    XCloseDisplay(display);
	//hWnd := ob.goWidget.hWnd
	//hDC := sysmetrics.GetDC(hWnd)*/
	/*deviceCaps.aspectX = sysmetrics.AspectX()
	deviceCaps.aspectY = sysmetrics.AspectY()
	deviceCaps.aspectXY = sysmetrics.AspectXY()

	deviceCaps.clientHeight = sysmetrics.ClientHeight()
	deviceCaps.clientWidth = sysmetrics.ClientWidth()
	deviceCaps.height = sysmetrics.Height()
	deviceCaps.width = sysmetrics.Width()
	deviceCaps.vertSize = sysmetrics.VerticalSize() //hDC)		//GetDeviceCaps(hDC, w32.VERTSIZE)
	deviceCaps.horizSize = sysmetrics.HorizontalSize() //hDC) 	//GetDeviceCaps(hDC, w32.HORZSIZE)
	deviceCaps.vertRes = sysmetrics.VerticalRes() //hDC) 		//GetDeviceCaps(hDC, w32.LOGPIXELSY)
	deviceCaps.horizRes = sysmetrics.HorizontalRes() //hDC) 	//GetDeviceCaps(hDC, w32.LOGPIXELSX)*/

//}

func AspectX() int {
	return 0
}

func AspectY() int {
	return 0
}

func AspectXY() int {
	return 0
}

func ClientHeight() int {
	return 0
}

func ClientWidth() int {
	return 0
}

func Height() int {
	return deviceCaps.height
}

func HorizontalRes() int {
	return 0
}

func HorizontalSize() int {
	return 0
}

func TaskBarHeight() int {
	return 0
}

func VerticalRes() int {
	return 0
}

func VerticalSize() int {
	return 0
}

func Width() int {
	return deviceCaps.width
}
