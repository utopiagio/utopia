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
	ModeWindowed GoApplicationMode = iota 	// enables all GoWindows.
	ModeModal  	// sets modal window on top.
)

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
		mode: ModeWindowed,
	}
	return GoApp
}

func (a *GoApplicationObj) AddWindow(w *GoWindowObj) {
	a.windows = append(a.windows, w)
}


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

func (a *GoApplicationObj) ClipBoard() (clipboard *GoClipBoardObj) {
	return a.clipboard
}

func (a *GoApplicationObj) Keyboard() (keyboard *GoKeyboardObj) {
	return a.keyboard
}

func (a *GoApplicationObj) SetModal(modalWin *GoWindowObj) {
	a.modalWindow = modalWin
	if modalWin == nil {
		a.mode = ModeWindowed
		for _, w := range a.windows {
			if !w.IsModal() {
				w.eventmask.Hide()
				w.Refresh()
			}
		}
	} else {
		a.mode = ModeModal
		for _, w := range a.windows {
			if !w.IsModal() {
				w.eventmask.Show()
				w.Refresh()
			}
		}
	}
}

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

//- (**[GoWindowObj]**) **AddPopupMenu()** (popupMenu [**\*GoPopupMenuObj**][1])  -//
//- ... Add a new popup menu.  -//
func (ob *GoWindowObj) AddPopupMenu() (popupMenu *GoPopupMenuObj) {
	popupMenu = GoPopupMenu(ob)
	ob.popupmenus = append(ob.popupmenus, popupMenu)
	return
}

//- (**[GoWindowObj]**) **Centre()**  -//
//- ... Centre the window on the client screen.  -//
func (ob *GoWindowObj) Centre() {
	screen.Width()
	ob.SetPos((screen.ClientWidth() - ob.Width) / 2, (screen.ClientHeight() - ob.Height) / 2)
}

//- (**[GoWindowObj]**) **ClearPopupMenus()**  -//
//- ... Clear the popup menus.  -//
func (ob *GoWindowObj) ClearPopupMenus() {
	ob.popupmenus = []*GoPopupMenuObj{}
	ob.Refresh()
}

//- (**[GoWindowObj]**) **Close()**  -//
//- ... Close the window.  -//
func (ob *GoWindowObj) Close() {
	//if ob.IsMainWindow() {
		ob.gio.Perform(system.ActionClose)
	//}
}

//- (**[GoWindowObj]**) **EscFullScreen()**  -//
//- ... Escape fullscreen.  -//
func (ob *GoWindowObj) EscFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Windowed.Option())
	}
}

//- (**[GoWindowObj]**) **GoFullScreen()**  -//
//- ... Switch to fullscreen.  -//
func (ob *GoWindowObj) GoFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Fullscreen.Option())
	}
}

//- (**[GoWindowObj]**) **IsMainWindow()** ([**bool**][2])  -//
//- ... Returns **true** if the window is the main window.  -//
func (ob *GoWindowObj) IsMainWindow() (isMain bool) {
	return ob.mainwindow
}

//- (**[GoWindowObj]**) **IsModal()** ([**bool**][2])  -//
//- Returns **true** if the window is a modal window.  -//
func (ob *GoWindowObj) IsModal() (isModal bool) {
	return ob.modalwindow
}

//- (**[GoWindowObj]**) **Layout()** (layout [**\*GoLayoutObj**][3])  -//
//- Returns a pointer to the window centre layout.  -//
func (ob *GoWindowObj) Layout() (layout *GoLayoutObj) {
	return ob.layout
}

//- (**[GoWindowObj]**) **Maximize()**  -//
//- Maximize the window.  -//
func (ob *GoWindowObj) Maximize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Maximized.Option())
	}
}

//- (**[GoWindowObj]**) **Minimize()**  -//
//- Minimize the window.  -//
func (ob *GoWindowObj) Minimize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Minimized.Option())
	}
}

//- (**[GoWindowObj]**) **MenuBar()** (layout [**\*GoMenuBarObj**][3])  -//
//- Installs and returns a pointer to the window main menu bar.  -//
func (ob *GoWindowObj) MenuBar() *GoMenuBarObj {
	ob.menubar.Show()
	return ob.menubar
}

//- (**[GoWindowObj]**) **MenuPopup(** idx **[int][int] )** ( popupMenu [**\*GoPopupMenuObj**][GoPopupMenuObj] )  -//
//- ... Returns a pointer to the popup menu at index idx.  -//
func (ob *GoWindowObj) MenuPopup(idx int) (popupMenu *GoPopupMenuObj) {
	if len(ob.popupmenus) > idx {
		return ob.popupmenus[idx]
	} else {
		return ob.popupmenus[len(ob.popupmenus) - 1]
	}
	return nil
}

func (ob *GoWindowObj) ModalStyle() (string) {
	return ob.modalstyle
}

func (ob *GoWindowObj) PopupWindow() *GoPopupWindowObj {
	return ob.popupwindow
}

func (ob *GoWindowObj) ObjectType() (string) {
	return "GoWindowObj"
}

func (ob *GoWindowObj) ClientPos() (x int, y int) {
	wx, wy := ob.gio.GetClientPos()
	return metrics.PxToDp(GoDpr, wx), metrics.PxToDp(GoDpr, wy)
}

func (ob *GoWindowObj) Pos() (x int, y int) {
	wx, wy := ob.gio.GetWindowPos()
	return metrics.PxToDp(GoDpr, wx), metrics.PxToDp(GoDpr, wy)
}

func (ob *GoWindowObj) Widget() (*GioWidget) {
	return nil
}

func (ob *GoWindowObj) GioWindow() *app_gio.Window {
	return ob.gio
}

func (ob *GoWindowObj) Refresh() {
	if ob.gio != nil {
		ob.gio.Invalidate()
	}
}

func (ob *GoWindowObj) SetBorder(style GoBorderStyle, width int, radius int, color GoColor) {
	ob.layout.SetBorder(style, width, radius, color)
}

func (ob *GoWindowObj) SetBorderColor(color GoColor) {
	ob.layout.SetBorderColor(color)
}

func (ob *GoWindowObj) SetBorderRadius(radius int) {
	ob.layout.SetBorderRadius(radius)
}

func (ob *GoWindowObj) SetBorderStyle(style GoBorderStyle) {
	ob.layout.SetBorderStyle(style)
}

func (ob *GoWindowObj) SetBorderWidth(width int) {
	ob.layout.SetBorderWidth(width)
}

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

func (ob *GoWindowObj) SetMargin(left int, top int, right int, bottom int) {
	ob.layout.SetMargin(left, top, right, bottom)
}

func (ob *GoWindowObj) SetOnClose(f func()) {
	ob.onClose = f
}

func (ob *GoWindowObj) SetOnConfig(f func()) {
	ob.onConfig = f
}

func (ob *GoWindowObj) Raise() {
	ob.gio.Perform(system.ActionRaise)
}

func (ob *GoWindowObj) SetPadding(left int, top int, right int, bottom int) {
	ob.layout.SetPadding(left, top, right, bottom)
}

func (ob *GoWindowObj) SetPos(x int, y int) {
	ob.X = x
	ob.Y = y
	if ob.gio != nil {
		ob.gio.Option(app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)))
	}
	ob.Refresh()
}

func (ob *GoWindowObj) SetSize(width int, height int) {
	ob.Width = width
	ob.Height = height
	if ob.gio != nil {
		ob.gio.Option(app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)))
	}
	ob.Refresh()
}

func (ob *GoWindowObj) SetSpacing(spacing GoLayoutSpacing) {
	ob.layout.SetSpacing(spacing)
}

func (ob *GoWindowObj) SetTitle(title string) {
	ob.title = title
	if ob.gio != nil {
		ob.gio.Option(app_gio.Title(title))
	}
}

func (ob *GoWindowObj) ClientSize() (width int, height int) {
	ww, wh := ob.gio.GetClientSize()
	return metrics.PxToDp(GoDpr, ww), metrics.PxToDp(GoDpr, wh)
}

func (ob *GoWindowObj) Size() (width int, height int) {
	ww, wh := ob.gio.GetWindowSize()
	return metrics.PxToDp(GoDpr, ww), metrics.PxToDp(GoDpr, wh)
}

func (ob *GoWindowObj) Show() {
	ob.run()
}

func (ob *GoWindowObj) ShowModal() (action int, info string) {
	action, info = ob.runModal()
	return
}

func (ob *GoWindowObj) Title() (title string) {
	return ob.title
}

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

func (ob *GoWindowObj) runModal() (action int, info string) {
  go func() {  
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
			//log.Println("Modal Dialog Style: GoFileDialog")
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
	}()
	//time.Sleep(200 * time.Millisecond)		// Is this required?
	return action, info
}

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
			    if GoApp.mode == ModeModal {
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


func (ob *GoWindowObj) paint(e app_gio.FrameEvent, gtx layout_gio.Context) {
	//log.Println("GoWindow.paint(e, gtx)")
	e.Frame(gtx.Ops)
}

func (ob *GoWindowObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	
	// signal for window events
	ob.signalEvents(gtx)
		
	// draw window frame layout
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

func (ob *GoWindowObj) signalEvents(gtx layout_gio.Context) {
	//if GoApp.Keyboard().GetFocus() == nil {

		event_gio.Op(gtx.Ops, 0)
		/*key_gio.FocusOp{
			Tag: 0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
		}.Add(gtx.Ops)

		// 3) Finally we add key.InputOp to catch specific keys
		// (Shift) means an optional Shift
		// These inputs are retrieved as key.Event
		key_gio.InputOp{
			//Keys: key_gio.Set("F|S|U|D|J|(K|(W|N|Space"),
			Tag:  0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
		}.Add(gtx.Ops)
	}
		pointer_gio.InputOp{
			Tag:   0,
			Grab:  false,
			Kinds: pointer_gio.Press,
		}.Add(gtx.Ops)
	}*/
}

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
								gtx.Execute(key_gio.FocusCmd{Tag: nil})
								GoApp.Keyboard().SetFocusControl(nil)
							}
					}
			}
	}
	GoApp.Keyboard().Update()
}

func (ob *GoWindowObj) updateLayout(layout GoObject, gtx layout_gio.Context) {
	for _, obj := range layout.Objects() {
		if obj.ObjectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		}
	}
}