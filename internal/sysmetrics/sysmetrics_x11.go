// SPDX-License-Identifier: Unlicense OR MIT

//go:build ((linux && !android) || freebsd || openbsd) && !nox11
// +build linux,!android freebsd openbsd
// +build !nox11

package sysmetrics

/*
#cgo freebsd openbsd CFLAGS: -I/usr/X11R6/include -I/usr/local/include
#cgo freebsd openbsd LDFLAGS: -L/usr/X11R6/lib -L/usr/local/lib
#cgo freebsd openbsd LDFLAGS: -lX11 -lxkbcommon -lxkbcommon-x11 -lX11-xcb -lXcursor -lXfixes
#cgo linux pkg-config: x11 xkbcommon xkbcommon-x11 x11-xcb xcursor xfixes

#include <stdlib.h>
#include <locale.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <X11/Xresource.h>
#include <X11/XKBlib.h>
#include <X11/Xlib-xcb.h>
#include <X11/extensions/Xfixes.h>
#include <X11/Xcursor/Xcursor.h>
#include <xkbcommon/xkbcommon-x11.h>

*/
import "C"
import "log"

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
		log.Println("return nil")
	}
	defer C.XCloseDisplay(dy)

	screen := C.XDefaultScreenOfDisplay(dy)
	width := int(C.XWidthOfScreen(screen))
	height := int(C.XHeightOfScreen(screen))
	deviceCaps.height = height
	deviceCaps.width = width
	deviceCaps.clientHeight = height
	deviceCaps.clientWidth = width

}


/*func (ob goDeviceCaps) getDeviceCaps(win *app_gio.Window) {

	display *C.Display
  screen *C.Screen

	// open a display
  display = C.XOpenDisplay(nil)

	screen = C.XDefaultScreen(win)

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

}*/

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