// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/application.go */

package utopia

import (
	"errors"
	"log"
	"time"
	"os"

	app_gio "github.com/utopiagio/gio/app"
	"github.com/utopiagio/gio/font/gofont"
	event_gio "github.com/utopiagio/gio/io/event"
	key_gio "github.com/utopiagio/gio/io/key"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	"github.com/utopiagio/gio/io/system"
	layout_gio "github.com/utopiagio/gio/layout"
	"github.com/utopiagio/gio/op"

	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	_ "github.com/utopiagio/gio/widget"

	"github.com/utopiagio/utopia/desktop"
	"github.com/utopiagio/utopia/metrics"
	screen "github.com/utopiagio/utopia/internal/sysmetrics"
)

type (
	D = layout_gio.Dimensions
	C = layout_gio.Context
)

type GoApplicationMode int

const (
	WindowedMode GoApplicationMode = iota 	// enables all GoWindows.
	ModalMode  	// sets modal window on top.
)

var tagCounter int
var GoDpr float32
var GoSpr float32
var GoApp *GoApplicationObj = nil

type GoApplicationObj struct {
	name string
	windows	[]*GoWindowObj
	clipboard *GoClipBoardObj
	//desktop *GoDeskTopObj
	keyboard *GoKeyboardObj
	// Theme contains semantic style data. Extends `material.Theme`.
	theme *GoThemeObj
	//theme *material_gio.Theme
	// Shaper cache of registered fonts.
	shaper *text_gio.Shaper
	//fontCollection []text_gio.FontFace
	//dpr float32
	mode GoApplicationMode
	modalWindow *GoWindowObj
}

//- <a name=\"goApplication\"></a> [**GoApplication**](api.GoApplication#goApplication)( appName **string** )  ( app [***GoApplicationObj**](#goApplicationObj) )\n
//- Initialise the application. Instantiate the GoApp global reference.\n\n
func GoApplication(appName string) (a *GoApplicationObj) {
	clipboard := GoClipBoard()
	if clipboard.init() != nil {
		log.Println("ClipBoard Not Available!")
	}
	desktop.Init()
	keyboard := GoKeyboard()
	/*if keyboard.init() != nil {
		log.Println("Keyboard Not Available!")
	}*/
	theme := GoTheme(gofont.Collection())
	GoApp = &GoApplicationObj{
		name: appName,
		clipboard: clipboard,
		//desktop: desktop,
		keyboard: keyboard,
		theme: theme,
		mode: WindowedMode,
	}
	return GoApp
}

//- <a name=\"addWindow\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**AddWindow(** win [***GoWindowObj**](api.GoWindow#) **)**\n
//- - Add a new window to the application.\n\n
func (a *GoApplicationObj) AddWindow(w *GoWindowObj) {
	a.windows = append(a.windows, w)
}

//- <a name=\"removeWindow\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**RemoveWindow(** win [***GoWindowObj**](api.GoWindow#) **)**\n
//- - Remove a window from the application.\nIf the window is the main window then the application will be shut down.\n\n
func (a *GoApplicationObj) RemoveWindow(w *GoWindowObj) {
	if w.IsMainWindow() {
		os.Exit(0)
	}
	k := 0
	for _, v := range a.windows {
	    if v != w {
	        a.windows[k] = v
	        k++
	    }
	}
	a.windows = a.windows[:k] // set slice len to remaining elements
}

//- <a name=\"run\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**Run()**\n
//- - Run the application main loop.\n\n
func (a *GoApplicationObj) Run() {
	var gio *app_gio.Window = nil
	if len(a.windows) == 0 {
		err := errors.New("****************\n\nApplication has no main window!\n" +
											"Use GoMainWindow()) method to create new windows.\n\n")
		log.Fatal(err)
	}
	gio = a.windows[0].gio
	if gio == nil {
		err := errors.New("****************\n\nApplication has no active main window!\n" +
											"Use GoWindow.Show() method to activate windows.\n\n")
		log.Fatal(err)
	}
	a.windows[0].mainwindow = true
	// on windows, linux, darwin the gio.Main functions blocks the main thread forever
	app_gio.Main()	
}

//- <a name=\"clipBoard\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**ClipBoard() (** clipboard [***GoClipBoardObj**](api.GoClipBoard#) **)**\n
//- - Return the application clipboard.\n\n
func (a *GoApplicationObj) ClipBoard() (clipboard *GoClipBoardObj) {
	return a.clipboard
}

//- <a name=\"keyboard\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**Keyboard() (** keyboard [***GoKeyboardObj**](api.GoKeyboard#) **)**\n
//- - Return the application keyboard.\n\n
func (a *GoApplicationObj) Keyboard() (keyboard *GoKeyboardObj) {
	return a.keyboard
}

//- <a name=\"setModal\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**SetModal(** modalWin [***GoWindowObj**](api.GoWindow#) **)**\n
//- - Set the window to run as a modal window.\n All other windows running under the application will be disabled.\n\n
func (a *GoApplicationObj) SetModal(modalWin *GoWindowObj) {
	a.modalWindow = modalWin
	if modalWin == nil {
		a.mode = WindowedMode
		for _, w := range a.windows {
			if !w.IsModal() {
				w.eventmask.Hide()
				w.Refresh()
			}
		}
	} else {
		a.mode = ModalMode
		for _, w := range a.windows {
			if !w.IsModal() {
				w.eventmask.Show()
				w.Refresh()
			}
		}
	}
}

//- <a name=\"theme\"></a> **(ob** [***GoApplicationObj**](api.GoApplication#)**)**.**Theme() (** theme [***GoThemeObj**](api.GoTheme#) **)**\n
//- - Return the application main theme.\n\n
func (a *GoApplicationObj) Theme() (theme *GoThemeObj) {
	return a.theme
}

/*func (a *GoApplicationObj) Theme() (theme *material_gio.Theme) {
	t := &material_gio.Theme{
		Shaper: text_gio.NewShaper(a.fontCollection),
	}
	t.Palette = material_gio.Palette{
		Fg:         a.theme.Palette.Fg,
		Bg:         a.theme.Palette.Bg,
		ContrastBg: a.theme.Palette.ContrastBg,
		ContrastFg: a.theme.Palette.ContrastFg,
	}
	t.TextSize = a.theme.TextSize

	t.Icon = a.theme.Icon

	// 38dp is on the lower end of possible finger size.
	t.FingerSize = a.theme.FingerSize

	return t
}*/

type GoWindowObj struct {
	GioObject
	//goWidget
	GoSize 			// Current, Min and Max sizes
	GoPos
	gio *app_gio.Window
	title string
	frame *GoLayoutObj
	menubar *GoMenuBarObj
	statusbar *GoLayoutObj
	layout *GoLayoutObj
	eventmask *GoEventMaskObj
	mainwindow bool
	modalwindow bool
	modalstyle string
	ModalAction int
	ModalInfo string
	popupmenus []*GoPopupMenuObj
	popupwindow *GoPopupWindowObj
	onClose func() 
	onConfig func()
}

//- <a name=\"goMainWindow\"></a> [**GoMainWindow**](api.GoWindow#goMainWindow)( windowTitle **string** )  ( hWin [***GoWindowObj**](#goWindowObj) )\n
//- - Create a new main window\n\n
func GoMainWindow(windowTitle string) (hWin *GoWindowObj) {
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{320, 480, 640, 480, 1500, 1200, 640, 480}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, nil, false, false, "", -1, "", nil, nil, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	hWin.eventmask = GoEventMask(hWin)
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

//- <a name=\"goWindow\"></a> [**GoWindow**](api.GoWindow#goWindow)( windowTitle **string** )  ( hWin [***GoWindowObj**](#goWindowObj) )\n
//- - Create a new window\n\n
func GoWindow(windowTitle string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{0, 0, 640, 480, 1500, 1200, 640, 480}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, nil, false, false, "", -1, "", nil, nil, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	hWin.eventmask = GoEventMask(hWin)
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

//- <a name=\"goModalWindow\"></a> [**GoModalWindow**](api.GoWindow#goModalWindow)( windowTitle **string** )  ( hWin [***GoWindowObj**](#goWindowObj) )\n
//- - Create a new modal window\n\n
func GoModalWindow(modalStyle string, windowTitle string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{0, 0, 640, 450, 1500, 1000, 640, 450}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, nil, false, true, modalStyle, -1, "", nil, nil, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	hWin.eventmask = GoEventMask(hWin)
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}
//- <a name=\"addPopupMenu\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**AddPopupMenu(** popupMenu [***GoPopupMenuObj**](api.GoPopupMenu#) **)**\n
//- - Add a new popup menu.\n\n
func (ob *GoWindowObj) AddPopupMenu() (popupMenu *GoPopupMenuObj) {
	popupMenu = GoPopupMenu(ob)
	ob.popupmenus = append(ob.popupmenus, popupMenu)
	return
}

//- <a name=\"centre\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Centre()**\n
//- - Centre the window on the client screen.\n\n
func (ob *GoWindowObj) Centre() {
	screen.Width()
	ob.SetPos((screen.ClientWidth() - ob.Width) / 2, (screen.ClientHeight() - ob.Height) / 2)
}

//- <a name=\"clearPopupMenus\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ClearPopupMenus()**\n
//- - Clear the popup menus.\n\n
func (ob *GoWindowObj) ClearPopupMenus() {
	ob.popupmenus = []*GoPopupMenuObj{}
	ob.Refresh()
}

//- <a name=\"clientPos\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ClientPos(** x **int,** y **int )**\n
//- - Returns the client postion of the window.\nThis is usually set to (0,0)\n\n
func (ob *GoWindowObj) ClientPos() (x int, y int) {
	wx, wy := ob.gio.GetClientPos()
	return metrics.PxToDp(GoDpr, wx), metrics.PxToDp(GoDpr, wy)
}

//- <a name=\"clientSize\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ClientSize() (** width **int,** height **int )**\n
//- - Returns the inner client size of the window.\n\n
func (ob *GoWindowObj) ClientSize() (width int, height int) {
	ww, wh := ob.gio.GetClientSize()
	return metrics.PxToDp(GoDpr, ww), metrics.PxToDp(GoDpr, wh)
}

//- <a name=\"close\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Close()**\n
//- - Close the window.\n\n
func (ob *GoWindowObj) Close() {
	//if ob.IsMainWindow() {
		ob.gio.Perform(system.ActionClose)
	//}
}

//- <a name=\"escFullScreen\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**EscFullScreen()**\n
//- - Escape fullscreen.\n\n
func (ob *GoWindowObj) EscFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Windowed.Option())
	}
}

//- <a name=\"goFullScreen\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**GoFullScreen()**\n
//- - Switch to fullscreen.\n\n
func (ob *GoWindowObj) GoFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Fullscreen.Option())
	}
}

//- <a name=\"isMainWindow\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**IsMainWindow() ( bool )**\n
//- - Returns **true** if the window is the main window.\n\n
func (ob *GoWindowObj) IsMainWindow() (isMain bool) {
	return ob.mainwindow
}

//- <a name=\"isModal\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**IsModal() ( bool )**\n
//- - Returns **true** if the window is a modal window.\n\n
func (ob *GoWindowObj) IsModal() (isModal bool) {
	return ob.modalwindow
}

//- <a name=\"layout\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Layout() ( layout** [***GoLayoutObj**](api.GoLayout#) **)**\n
//- - Returns a pointer to the window central layout.\n\n
func (ob *GoWindowObj) Layout() (layout *GoLayoutObj) {
	return ob.layout
}

//- <a name=\"maximize\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Maximize()**\n
//- - Maximize the window.\n\n
func (ob *GoWindowObj) Maximize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Maximized.Option())
	}
}

//- <a name=\"menubar\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**MenuBar() (** menubar [***GoMenuBarObj**](api.GoMenuBar#) **)**\n
//- - Installs and returns a pointer to the window main menu bar.\n\n
func (ob *GoWindowObj) MenuBar() *GoMenuBarObj {
	ob.menubar.Show()
	return ob.menubar
}

//- <a name=\"menuPopup\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**MenuPopup() (** popupMenu [***GoPopupMenuObj**](api.GoPopupMenu#) **)**\n
//- - Returns a pointer to the popup menu at index idx.\n\n
func (ob *GoWindowObj) MenuPopup(idx int) (popupMenu *GoPopupMenuObj) {
	if len(ob.popupmenus) > idx {
		return ob.popupmenus[idx]
	} else {
		return ob.popupmenus[len(ob.popupmenus) - 1]
	}
	return nil
}

//- <a name=\"minimize\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Minimize()**\n
//- - Minimize the window.\n\n
func (ob *GoWindowObj) Minimize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Minimized.Option())
	}
}

//- <a name=\"modalStyle\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ModalStyle() (** style **string )**\n
//- - Returns the modal style of a modal window.\n\n
func (ob *GoWindowObj) ModalStyle() (string) {
	return ob.modalstyle
}

//- <a name=\"popupWindow\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**PopupWindow() (** popupWindow [***GoPopupWindowObj**](api.GoPopupWindow#) **)**\n
//- - Returns a pointer to the windows popup window.\n\n
func (ob *GoWindowObj) PopupWindow() *GoPopupWindowObj {
	return ob.popupwindow
}

//- <a name=\"objectType\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ObjectType() (** type **string )**\n
//- - Returns the object type as a string definition \"GoWindowObj\".\n\n
func (ob *GoWindowObj) ObjectType() (string) {
	return "GoWindowObj"
}

//- <a name=\"pos\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Pos() (** x **int, y int )**\n
//- - Returns the screen postion of the window.\n\n
func (ob *GoWindowObj) Pos() (x int, y int) {
	wx, wy := ob.gio.GetWindowPos()
	return metrics.PxToDp(GoDpr, wx), metrics.PxToDp(GoDpr, wy)
}

/*func (ob *GoWindowObj) GioWindow() *app_gio.Window {
	return ob.gio
}*/

//- <a name=\"raise\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Raise()**\n
//- - Raises the window to be top most window.\n\n
func (ob *GoWindowObj) Raise() {
	ob.gio.Perform(system.ActionRaise)
}

//- <a name=\"refresh\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Refresh()**\n
//- - Refresh the window.\n\n
func (ob *GoWindowObj) Refresh() {
	if ob.gio != nil {
		ob.gio.Invalidate()
	}
}

//- <a name=\"setBorder\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetBorder(** style [**GoBorderStyle**](api.GoBorderStyle#), width **int**, radius **int**, color [**GoColor**](api.GoColor#) **)**\n
//- - Add a border to the window.\n\n
func (ob *GoWindowObj) SetBorder(style GoBorderStyle, width int, radius int, color GoColor) {
	ob.layout.SetBorder(style, width, radius, color)
}

//- <a name=\"setBorderColor\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetBorderColor(** color [**GoColor**](api.GoColor#) **)**\n
//- - Change the border color of the window border.\n\n
func (ob *GoWindowObj) SetBorderColor(color GoColor) {
	ob.layout.SetBorderColor(color)
}

//- <a name=\"setBorderRadius\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetBorderRadius(** radius **int )**\n
//- - Change the border radius of the window border.\n\n
func (ob *GoWindowObj) SetBorderRadius(radius int) {
	ob.layout.SetBorderRadius(radius)
}

//- <a name=\"setBorderStyle\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetBorderStyle(** style [**GoBorderStyle**](api.GoBorderStyle#) **)**\n
//- - Change the border style of the window border.\n\n
func (ob *GoWindowObj) SetBorderStyle(style GoBorderStyle) {
	ob.layout.SetBorderStyle(style)
}

//- <a name=\"setBorderWidth\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetBorderWidth(** width **int )**\n
//- - Change the border width of the window border.\n\n
func (ob *GoWindowObj) SetBorderWidth(width int) {
	ob.layout.SetBorderWidth(width)
}

//- <a name=\"setLayoutStyle\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetLayoutStyle(** style [**GoLayoutStyle**](api.GoLayoutStyle#) **)**\n
//- - Changes the window central layout style.\n\n
func (ob *GoWindowObj)SetLayoutStyle(style GoLayoutStyle) {
	ob.frame.DeleteControl(ob.layout)
	switch style {
	case NoLayout:
		ob.layout = nil
	case HBoxLayout:
		ob.layout = GoHBoxLayout(ob.frame)
	case VBoxLayout:
		ob.layout = GoVBoxLayout(ob.frame)	
	case HVBoxLayout:
		// Not Implemented *******************
	case HFlexBoxLayout:
		ob.layout = GoHFlexBoxLayout(ob.frame)	
	case VFlexBoxLayout:						
		ob.layout = GoVFlexBoxLayout(ob.frame)	
	case PopupMenuLayout:
		// Not Implemented *******************
	}
}

//- <a name=\"setMargin\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetMargin(** left **int,** top **int,** bottom **int,** right **int )**\n
//- - Sets the window margin sizes.\n\n
func (ob *GoWindowObj) SetMargin(left int, top int, right int, bottom int) {
	ob.layout.SetMargin(left, top, right, bottom)
}

//- <a name=\"setOnClose\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetOnClose(** f **func() )**\n
//- - Adds a function to be called when the window is closed.\n\n
func (ob *GoWindowObj) SetOnClose(f func()) {
	ob.onClose = f
}

//- <a name=\"setOnConfig\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetOnConfig(** f **func() )**\n
//- - Adds a function to be called when the window is reconfigured.\n\n
func (ob *GoWindowObj) SetOnConfig(f func()) {
	ob.onConfig = f
}

//- <a name=\"setPadding\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetPadding(** left **int,** top **int,** bottom **int,** right **int )**\n
//- - Sets the window padding sizes.\n\n
func (ob *GoWindowObj) SetPadding(left int, top int, right int, bottom int) {
	ob.layout.SetPadding(left, top, right, bottom)
}

//- <a name=\"setPos\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetPos(** x **int, y int )**\n
//- - Moves the position of the window on screen.\n\n
func (ob *GoWindowObj) SetPos(x int, y int) {
	ob.X = x
	ob.Y = y
	if ob.gio != nil {
		ob.gio.Option(app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)))
	}
	ob.Refresh()
}

//- <a name=\"setSize\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetSize(** width **int, height int )**\n
//- - Resizes the window on screen.\n\n
func (ob *GoWindowObj) SetSize(width int, height int) {
	ob.Width = width
	ob.Height = height
	if ob.gio != nil {
		ob.gio.Option(app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)))
	}
	ob.Refresh()
}

//- <a name=\"setSpacing\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetSpacing(** style [**GoLayoutSpacing**](api.GoLayoutSpacing#) **)**\n
//- - Changes the window central layout widget spacing.\n\n
func (ob *GoWindowObj) SetSpacing(spacing GoLayoutSpacing) {
	ob.layout.SetSpacing(spacing)
}

//- <a name=\"setTitle\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**SetTitle(** title **string )**\n
//- - Change the text displayed in the window title bar.\n\n
func (ob *GoWindowObj) SetTitle(title string) {
	ob.title = title
	if ob.gio != nil {
		ob.gio.Option(app_gio.Title(title))
	}
}

//- <a name=\"size\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Size() (** width **int,** height **int )**\n
//- - Returns the outer size of the window.\n\n
func (ob *GoWindowObj) Size() (width int, height int) {
	ww, wh := ob.gio.GetWindowSize()
	return metrics.PxToDp(GoDpr, ww), metrics.PxToDp(GoDpr, wh)
}

//- <a name=\"show\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Show()**\n
//- - Activate the window loop and set as top window.\n\n
func (ob *GoWindowObj) Show() {
	ob.run()
}

//- <a name=\"showModal\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**ShowModal()**\n
//- - Activate the window as a modal application window and set as topmost window.\n\n
func (ob *GoWindowObj) ShowModal() (action int, info string) {
	action, info = ob.runModal()
	return
}

//- <a name=\"title\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**Title() (** title **string )**\n
//- - Return the text displayed by the window title bar.\n\n
func (ob *GoWindowObj) Title() (title string) {
	return ob.title
}

//- <a name=\"widget\"></a> **(ob** [***GoWindowObj)**](api.GoWindow#).**Widget() (** widget [**GioWidget**](api.GioWidget#) **)**\n
//- - Returns a pointer to the window widget properties.\n\n
func (ob *GoWindowObj) Widget() (*GioWidget) {
	return nil
}

//- <a name=\"run\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**run()**\n
//- - Creates the OS window and runs the window messaging loop, until the window is closed.\n\n
func (ob *GoWindowObj) run() {
	go func() {
	    // create new window
	    ob.gio = app_gio.NewWindow(
	      app_gio.Title(ob.title),
	      app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)),
	      app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)),
	    )
	    // draw on screen
	    if err := ob.loop(); err != nil {
	      log.Fatal(err)
		}
		GoApp.RemoveWindow(ob)
	}()
	time.Sleep(200 * time.Millisecond)
}

//- <a name=\"run\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**runModal()**\n
//- - Creates the OS window, sets the application to run in ModalMode and runs the modal window messaging loop, until the window is closed and returns its action and info.\n\n

func (ob *GoWindowObj) runModal() (action int, info string) {
    // create new modalwindow
    ob.gio = app_gio.NewWindow(
      app_gio.Title(ob.title),
      app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)),
      app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)),
    )
    GoApp.SetModal(ob)
    
    // draw on screen
    log.Println("Modal Dialog ob.loop()")
    if err := ob.loopModal(); err != nil {
      log.Fatal(err)
		}
		switch ob.ModalStyle() {
		case "GoFileDialog":
			log.Println("Modal Dialog Style: GoFileDialog")
			action = ob.ModalAction
			info = ob.ModalInfo
		case "GoMsgDialog":
			//log.Println("Modal Dialog Style: GoMsgDialog")
			action = ob.ModalAction
			info = ob.ModalInfo
		case "GoPrintDialog":
			//log.Println("Modal Dialog Style: GoPrintDialog")
			action = ob.ModalAction
			info = ob.ModalInfo
		default:
			action = 0
			info = ""
		}
		log.Println("ob.IsMainWindow() :", ob.IsMainWindow())
		GoApp.SetModal(nil)
		GoApp.RemoveWindow(ob)
	//time.Sleep(200 * time.Millisecond)		// Is this required?
	return action, info
}

//- <a name=\"run\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**loop()**\n
//- - Runs the window messaging loop, listening for Gio events, until the window is destroyed.\n\n
func (ob *GoWindowObj) loop() (err error) {
	// ops are the operations from the UI
    var ops op.Ops

    // listen for events in the window.
    for {
    	ev := ob.gio.NextEvent()
    	//if GoApp.mode == ModeWindowed {
    		//log.Println("ModeWindowed")
					// detect what type of event
					switch e := ev.(type) {
						//switch  e := ob.gio.NextEvent().(type) {
							case app_gio.DestroyEvent:
								if ob.onClose != nil {
			      			ob.onClose()
			      		}
			      		//log.Println("system.DestroyEvent.....")
			      		return e.Err
			      	// this is sent when the application should re-render.
			      	case app_gio.FrameEvent:
			      		//if !modal || ob.IsModal() {
				      		// Open an new context
				      		gtx := app_gio.NewContext(&ops, e)
				      		ob.update(gtx)		// receiveEvents
				      		ob.render(gtx)		// draw layout and signalEvents
				      		ob.paint(e, gtx)	// window paint
				      		
				      	//}
			      		//}
			      	case app_gio.ConfigEvent:
			      		if ob.onConfig != nil {
			      			ob.onConfig()
			      		}
			    }
			    if GoApp.mode == ModalMode {
						GoApp.modalWindow.Raise()
					}
	    /*for _, v := range GoApp.windows {
	    		if v.IsModal() {
	    				modal = true
	    		}
	    }*/
	    /*case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}
			ob.gio.Invalidate()			// redraw window*/
		//}
			}
    //}
	return nil
}

//- <a name=\"run\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**loopModal()**\n
//- - Runs the modal window messaging loop, listening for Gio events, until the modal window is destroyed.\n\n
func (ob *GoWindowObj) loopModal() (err error) {
	// ops are the operations from the UI
    var ops op.Ops

    // listen for events in the window.
    for {
    			// detect what type of event
					switch  e := ob.gio.NextEvent().(type) {

							case app_gio.DestroyEvent:
								if ob.onClose != nil {
			      			ob.onClose()
			      		}
			      		//log.Println("system.DestroyEvent.....")
			      		return e.Err
			      	// this is sent when the application should re-render.
			      	case app_gio.FrameEvent:
			      		//if !modal || ob.IsModal() {
				      		// Open an new context
				      		gtx := app_gio.NewContext(&ops, e)
				      		ob.update(gtx)		// receiveEvents
				      		ob.render(gtx)		// draw layout and signalEvents
				      		ob.paint(e, gtx)	// window paint
				      		/*if !ob.IsModal() {
											if GoApp.mode == ModeModal {
												GoApp.modalWindow.Raise()
											}
									}*/
				      	//}
			      		//}
			      	case app_gio.ConfigEvent:
			      		if ob.onConfig != nil {
			      			ob.onConfig()
			      		}
			    }
	    /*for _, v := range GoApp.windows {
	    		if v.IsModal() {
	    				modal = true
	    		}
	    }*/
	    /*case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}*/
			//ob.gio.Invalidate()			// redraw window
		//}
    }
	return nil
}

//- <a name=\"paint\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**paint(** e **app_gio.FrameEvent,** gtx **layout_gio.Context )**\n
//- - calls the Gio window Frame to repaint the window.\n\n
func (ob *GoWindowObj) paint(e app_gio.FrameEvent, gtx layout_gio.Context) {
	e.Frame(gtx.Ops)
}

//- <a name=\"render\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**render(** gtx **layout_gio.Context ) ( layout_gio.Dimensions )**\n
//- - writes all the rendering to repaint the window.\n\n
func (ob *GoWindowObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	
	// signal for window events
	ob.signalEvents(gtx)
		
	// draw ops for window frame layout
	//log.Println("(ob *GoWindowObj) frame.............")
	dims := ob.frame.Draw(gtx)

	// draw menupopup modal layout
	//log.Println("(ob *GoWindowObj) modal.............")
	if len(ob.popupmenus) > 0 {
		ob.popupmenus[0].Draw(gtx)
		for idx := 0; idx < len(ob.popupmenus); idx++ {
			if ob.popupmenus[idx].Visible {
				//ob.popupmenus[idx].Draw(gtx)
				ob.popupmenus[idx].layout.Draw(gtx)
			}
		}
	}
	if ob.popupwindow.Visible {
		ob.popupwindow.Draw(gtx)
		ob.popupwindow.layout.Draw(gtx)
	}
	if ob.eventmask.Visible {
			ob.eventmask.Draw(gtx)
	}
	return dims
}

//- <a name=\"signalEvents\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**signalEvents(** gtx **layout_gio.Context )**\n
//- - sets the application wide tag to enable the window to  receive all events.\n\n
func (ob *GoWindowObj) signalEvents(gtx layout_gio.Context) {
		// Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
		event_gio.Op(gtx.Ops, 0)
}

//- <a name=\"update\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**update(** gtx **layout_gio.Context )**\n
//- - sets the global screen scaling, calls updateLayout to update all the window layouts and controls, receives events, updates the keyboard focus all the rendering to repaint the window.\n\n
func (ob *GoWindowObj) update(gtx layout_gio.Context) {
	// set global screen pixel size
	GoDpr = gtx.Metric.PxPerDp
	GoSpr = gtx.Metric.PxPerSp
	for _, obj := range ob.frame.Controls {
		if obj.ObjectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		}
	}
	for {
			event, ok := gtx.Event(
				key_gio.FocusFilter{},
				//key_gio.Filter{Name: "A"},
				//key_gio.Filter{Name: key_gio.NameSpace},
				// list of filters in the form of key_gio.Filter{ Name: key_gio.NameEnter}...
				// also pointer_gio.Filter{Target: tag, Kinds: pointer_gio.KindPress...}
				pointer_gio.Filter{Target: 0, Kinds: pointer_gio.Press},
			)

			if !ok { break }
			if ev, ok := event.(key_gio.Event); ok {
				   	log.Println("ApplicationKey::Event -", "Name -", ev.Name, "Modifiers -", ev.Modifiers, "State -", ev.State)
			} else if ev, ok := event.(key_gio.EditEvent); ok {
		    		log.Println("ApplicationKey::EditEvent -", "Range -", ev.Range, "Text -", ev.Text)
			} else if ev, ok := event.(pointer_gio.Event); ok {
				 	switch ev.Kind {
						case pointer_gio.Press:
							if ev.Priority == pointer_gio.Grabbed {
								if ob.popupwindow.Visible {
									ob.popupwindow.Hide()
								}
								log.Println("ApplicationKey::keyPress -")
								gtx.Execute(key_gio.FocusCmd{Tag: nil})
								GoApp.Keyboard().SetFocusControl(nil)
							}
					}
			}
	}
	GoApp.Keyboard().Update()
}

//- <a name=\"updateLayout\"></a> **(ob** [***GoWindowObj**](api.GoWindow#)**)**.**updateLayout(** layout **GoObject,** gtx **layout_gio.Context )**\n
//- - updates all the window layouts and controls\n\n
func (ob *GoWindowObj) updateLayout(layout GoObject, gtx layout_gio.Context) {
	for _, obj := range layout.Objects() {
		if obj.ObjectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		}
	}
}