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
)

type (
	D = layout_gio.Dimensions
	C = layout_gio.Context
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
		//fontCollection: gofont.Collection(),
		
	}
	return GoApp
}

func (a *GoApplicationObj) AddWindow(w *GoWindowObj) {
	a.windows = append(a.windows, w)
}

func (a *GoApplicationObj) RemoveWindow(w *GoWindowObj) {
	k := 0
	for _, v := range a.windows {
	    if v != w {
	        a.windows[k] = v
	        k++
	    }
	}
	a.windows = a.windows[:k] // set slice len to remaining elements
	if len(a.windows) == 0 {
		os.Exit(0)
	}
}

func (a *GoApplicationObj) Run() {
	var gio *app_gio.Window = nil
	if len(a.windows) == 0 {
		err := errors.New("****************\n\nApplication has no main window!\n" +
											"Use GoWindow()) method to create new windows.\n\n")
		log.Fatal(err)
	}
	gio = a.windows[0].gio
	if gio == nil {
		err := errors.New("****************\n\nApplication has no active main window!\n" +
											"Use GoWindow.Show() method to activate windows.\n\n")
		log.Fatal(err)
	}
	app_gio.Main()
}

func (a *GoApplicationObj) ClipBoard() (clipboard *GoClipBoardObj) {
	return a.clipboard
}

func (a *GoApplicationObj) Keyboard() (keyboard *GoKeyboardObj) {
	return a.keyboard
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
	//mainwindow bool
	modalwindow bool
	modalstyle string
	ModalAction int
	ModalInfo string
	popupmenus []*GoPopupMenuObj
	popupwindow *GoPopupWindowObj

	onConfig func()
}

func GoMainWindow(windowTitle string) (hWin *GoWindowObj) {
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{320, 480, 640, 480, 1500, 1200, 640, 480}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, false, "", -1, "", nil, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	
	hWin.menubar = GoMenuBar(hWin.frame)
	//hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	//hWin.AddPopupMenu(GoPopupMenu(hWin))
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

func GoWindow(windowTitle string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{0, 0, 640, 480, 1500, 1200, 640, 480}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, false, "", -1, "", nil, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	//hWin.popupmenus = GoPopupMenu(hWin)
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

func GoModalWindow(modalStyle string, windowTitle string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{0, 0, 640, 450, 1500, 1000, 640, 450}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowTitle, nil, nil, nil, nil, true, modalStyle, -1, "", nil, nil, nil}
	hWin.Window = hWin

	hWin.frame = GoVFlexBoxLayout(hWin)
	
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	hWin.menubar.SetBackgroundColor(Color_WhiteSmoke)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	//hWin.popupmenus = GoPopupMenu(hWin)
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

func (ob *GoWindowObj) AddPopupMenu() (popupMenu *GoPopupMenuObj) {
	popupMenu = GoPopupMenu(ob)
	//menu := GoMenuItem(ob, len(ob.popupmenus), nil)
	ob.popupmenus = append(ob.popupmenus, popupMenu)
	return
}

func (ob *GoWindowObj) ClearPopupMenus() {
	ob.popupmenus = []*GoPopupMenuObj{}
	ob.Refresh()
}

func (ob *GoWindowObj) Close() {
	ob.gio.Perform(system.ActionClose)
}

func (ob *GoWindowObj) EscFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Windowed.Option())
	}
}

func (ob *GoWindowObj) GoFullScreen() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Fullscreen.Option())
	}
}

func (ob *GoWindowObj) IsMainWindow() bool {
	//return ob.frame
	return !ob.modalwindow
}

func (ob *GoWindowObj) IsModal() bool {
	//return ob.frame
	return ob.modalwindow
}

func (ob *GoWindowObj) Layout() *GoLayoutObj {
	//return ob.frame
	return ob.layout
}

func (ob *GoWindowObj) Maximize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Maximized.Option())
	}
}

func (ob *GoWindowObj) Minimize() {
	if ob.gio != nil {
		ob.gio.Option(app_gio.Minimized.Option())
	}
}

func (ob *GoWindowObj) MenuBar() *GoMenuBarObj {
	ob.menubar.Show()
	return ob.menubar
}

func (ob *GoWindowObj) MenuPopup(idx int) *GoPopupMenuObj {
	if len(ob.popupmenus) > idx {
		return ob.popupmenus[idx]
	} else {
		return ob.popupmenus[len(ob.popupmenus) - 1]
	}
	return nil
}

/*func (ob *GoWindowObj) ModalAction() (int) {
	return ob.ModalAction
}

func (ob *GoWindowObj) ModalInfo() (string) {
	return ob.ModalInfo
}*/

func (ob *GoWindowObj) ModalStyle() (string) {
	return ob.modalstyle
}

func (ob *GoWindowObj) PopupWindow() *GoPopupWindowObj {
	//return ob.frame
	return ob.popupwindow
}

/*func (ob *GoWindowObj) ModalLayout() *GoLayoutObj {
	//return ob.frame
	return ob.modal.layout
}*/

func (ob *GoWindowObj) ObjectType() (string) {
	return "GoWindowObj"
}

func (ob *GoWindowObj) ClientPos() (x int, y int) {
	//log.Println("ob.gio.HWND() = ", ob.gio.HWND())
	//log.Println("GoWindowObj::ClientPos:()")
	//log.Println("ob.Pos = ", ob.X, ob.Y)
	x, y, _, _ = desktop.GetClientRect(ob.gio.HWND())
	x = metrics.PxToDp(GoDpr, x)
	y = metrics.PxToDp(GoDpr, y)
	//log.Println("GoWindowObj::ClientPos:", x, y)
	return
}

func (ob *GoWindowObj) Pos() (x int, y int) {
	//log.Println("ob.gio.HWND() = ", ob.gio.HWND())
	//log.Println("GoWindowObj::Pos:()")
	//log.Println("ob.Pos = ", ob.X, ob.Y)
	x, y, _, _ = desktop.GetWindowRect(ob.gio.HWND())
	x = metrics.PxToDp(GoDpr, x)
	y = metrics.PxToDp(GoDpr, y)
	//log.Println("GoWindowObj::Pos:", x, y)
	//ob.X = x
	//ob.Y = y
	return
}


func (ob *GoWindowObj) Widget() (*GioWidget) {
	return nil
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
	/*if style == HBoxLayout || style == VBoxLayout {
		ob.layout = GoBoxLayout(ob.frame, style)
	} else if style == HFlexBoxLayout || style == VFlexBoxLayout {
		ob.layout.style = HFlexBoxLayout	GoFlexBoxLayout(ob.frame, style)
	}*/
	ob.layout.style = style
}

func (ob *GoWindowObj) SetMargin(left int, top int, right int, bottom int) {
	ob.layout.SetMargin(left, top, right, bottom)
}

func (ob *GoWindowObj) SetOnConfig(f func()) {
	ob.onConfig = f
}

func (ob *GoWindowObj) SetPadding(left int, top int, right int, bottom int) {
	ob.layout.SetPadding(left, top, right, bottom)
}

func (ob *GoWindowObj) SetPos(x int, y int) {
	ob.X = x
	ob.Y = y
	//log.Println("GoWindowObj::SetPos:", x, y)
	//log.Println("GoDpr =", GoDpr)
	//log.Println("GoWindowObj::SetPos unit_gio:", unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y))
	//log.Println("GoWindowObj::SetPos metrics:", metrics.PxToDp(GoDpr, ob.X), metrics.PxToDp(GoDpr, ob.Y))
	if ob.gio != nil {
		ob.gio.Option(app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)))
	}
}

func (ob *GoWindowObj) SetSize(width int, height int) {
	ob.Width = width
	ob.Height = height
	//log.Println("GoWindowObj::SetSize:", width, height)
	//log.Println("GoDpr =", GoDpr)
	//log.Println("GoWindowObj::SetSize: unit_gio", unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height))
	//log.Println("GoWindowObj::SetSize metrics:", metrics.DpToPx(GoDpr, ob.Width), metrics.PxToDp(GoDpr, ob.Height))
	if ob.gio != nil {
		ob.gio.Option(app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)))
		//ob.gio.Option(app_gio.Size(ob.Width, ob.Height))
	}
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
	//log.Println("GoWindowObj::ClientSize:()")
	//log.Println("ob.ClientSize = ", ob.Width, ob.Height)
	_, _, width, height = desktop.GetClientRect(ob.gio.HWND())
	width = metrics.PxToDp(GoDpr, width)
	height = metrics.PxToDp(GoDpr, height)
	//log.Println("GoWindowObj::ClientSize:", width, height)
	return
}

func (ob *GoWindowObj) Size() (width int, height int) {
	_, _, width, height = desktop.GetWindowRect(ob.gio.HWND())
	width = metrics.PxToDp(GoDpr, width)
	height = metrics.PxToDp(GoDpr, height)
	return
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
		if ob.IsMainWindow() {
			os.Exit(0)
		}
		GoApp.RemoveWindow(ob)
	}()
	time.Sleep(200 * time.Millisecond)
}

func (ob *GoWindowObj) runModal() (action int, info string) {
    // create new modalwindow
    ob.gio = app_gio.NewWindow(
      app_gio.Title(ob.title),
      app_gio.Pos(unit_gio.Dp(ob.X), unit_gio.Dp(ob.Y)),
      app_gio.Size(unit_gio.Dp(ob.Width), unit_gio.Dp(ob.Height)),
    )
    // draw on screen
    log.Println("Modal Dialog ob.loop()")
    if err := ob.loop(); err != nil {
      log.Fatal(err)
	}
	switch ob.ModalStyle() {
	case "GoFileDialog":
		log.Println("Modal Dialog Style: GoFileDialog")
		action = ob.ModalAction
		info = ob.ModalInfo
	case "GoMsgDialog":
		log.Println("Modal Dialog Style: GoMsgDialog")
		action = ob.ModalAction
		info = ob.ModalInfo
	case "GoPrintDialog":
		log.Println("Modal Dialog Style: GoPrintDialog")
		action = ob.ModalAction
		info = ob.ModalInfo
	}
	log.Println("ob.IsMainWindow() :", ob.IsMainWindow())
	if ob.IsMainWindow() {
		os.Exit(0)
	}
	GoApp.RemoveWindow(ob)
	
	time.Sleep(200 * time.Millisecond)
	return action, info
}

func (ob *GoWindowObj) loop() (err error) {
	//var count int
	// ops are the operations from the UI
    var ops op.Ops

    // listen for events in the window.
    for {
		//select {
    	//case e := <-ob.gio.Events():

			// detect what type of event
			switch  e := ob.gio.NextEvent().(type) {
					case system.DestroyEvent:
	      		log.Println("system.DestroyEvent.....")
	      		return e.Err
	      	// this is sent when the application should re-render.
	      	case system.FrameEvent:
	      		// Open an new context
	      		gtx := layout_gio.NewContext(&ops, e)
	      		//log.Println("Window.update(gtx).....", count)
	      		//count++
	      		ob.update(gtx)	// receiveEvents
	      		//log.Println("Window.render(gtx).....")
	      		ob.render(gtx)	// draw layout and signalEvents
	      		// window paint
	      		//e.Frame(gtx.Ops)
	      		//log.Println("Window.paint(e, gtx).....")
	      		ob.paint(e, gtx)
	      	case app_gio.ConfigEvent:
	      		if ob.onConfig != nil {
	      			ob.onConfig()
	      		}
	      	}
	      
	    /*case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}
			ob.gio.Invalidate()			// redraw window*/
		//}
    }
	return nil
}

func (ob *GoWindowObj) paint(e system.FrameEvent, gtx layout_gio.Context) {
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
	return dims
}

func (ob *GoWindowObj) signalEvents(gtx layout_gio.Context) {
	if GoApp.Keyboard().GetFocus() == nil {
		key_gio.FocusOp{
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
}

func (ob *GoWindowObj) update(gtx layout_gio.Context) {
	//log.Println("(ob *GoWindowObj) update.............")
	// set global screen pixel size
	GoDpr = gtx.Metric.PxPerDp
	GoSpr = gtx.Metric.PxPerSp
	//log.Println("GoDpr =", GoDpr)
	for _, obj := range ob.frame.Controls {
		if obj.ObjectType() == "GoLayoutObj" {
			//log.Println("(ob *GoLayoutObj) updateLayout.............")
			ob.updateLayout(obj, gtx)
		} /*else {
			if obj.ObjectType() == "GoButtonObj"{
				button := obj.(*GoButtonObj)
				if button.Clicked() {
					log.Println("GoWindow::GoButtonObj:Clicked()")
					//GoApp.Keyboard().SetFocus(nil)
					if button.onClick != nil{
						button.onClick()
					}
				} 
				if button.Focused() {
					log.Println("GoWindow::GoButtonObj:Focused()")
					if button.onFocus != nil {
						button.onFocus()
					}
				}
				if button.Hovered() {
					log.Println("GoWindow::GoButtonObj:Hovered()")
					if button.onHover != nil {
						button.onHover()
					}
				}
				if button.Pressed() {
					log.Println("GoWindow::GoButtonObj:Pressed()")
					//GoApp.Keyboard().SetFocus(nil)
					if button.onPress != nil {
						button.onPress()
					}
				}
			} else if obj.ObjectType() == "GoRadioButtonObj"{
				button := obj.(*GoRadioButtonObj)
				if button.Changed() {
					log.Println("GoWindow::GoRadioButtonObj:Changed()")
					if button.onChange != nil{
						button.onChange()
					}
				} 
				if tag, focus := button.Focused(); focus {
					log.Println("GoWindow::GoRadioButtonObj:Focused()")
					if button.onFocus != nil {
						button.onFocus(tag)
					}
				}
				if tag, hover := button.Hovered(); hover {
					log.Println("GoWindow::GoRadioButtonObj:Hovered()")
					if button.onHover != nil {
						button.onHover(tag)
					}
				}

			} else if obj.ObjectType() == "GoSwitchObj" {
				swtch := obj.(*GoSwitchObj)
				if swtch.Changed() {
					log.Println("GoWindow::GoSwitchObj:Changed()")
					if swtch.onChange != nil {
						swtch.onChange(swtch.goSwitch.Value)
					}
				}
				if swtch.Focused() {
					log.Println("GoWindow::GoSwitchObj:Focused()")
					if swtch.onFocus != nil {
						swtch.onFocus()
					}
				}
				if swtch.Hovered() {
					log.Println("GoWindow::GoSwitchObj:Hovered()")
					if swtch.onHover != nil {
						swtch.onHover()
					}
				}
				if swtch.Pressed() {
					log.Println("GoWindow::GoSwitchObj:Pressed()")
					if swtch.onPress != nil {
						swtch.onPress()
					}
				}
			} else if obj.ObjectType() == "GoSliderObj" {
				slider := obj.(*GoSliderObj)
				if slider.Changed() {
					//log.Println("GoSliderObj:Changed()")
					if slider.onChange != nil {
						slider.onChange(slider.gioSlider.Value)
					}
				}
				if slider.Dragging() {
					//log.Println("GoSliderObj:Dragging()")
					if slider.onDrag != nil {
						slider.onDrag(slider.gioSlider.Value)
					}
				}
			} //else if obj.objectType() == "GoTextEditObj" {
				//textedit := obj.(*GoTextEditObj)
		}*/
	}
	for _, gtxEvent := range gtx.Events(0) {
		//log.Println("gtxEvent -", gtxEvent.Type)
	    switch event := gtxEvent.(type) {
		    case key_gio.EditEvent:
				log.Println("ApplicationKey::EditEvent -", "Range -", event.Range, "Text -", event.Text)
		    case key_gio.Event:
		    	log.Println("ApplicationKey::Event -", "Name -", event.Name, "Modifiers -", event.Modifiers, "State -", event.State)
		    case pointer_gio.Event:
		    	//log.Println("ApplicationPointer::Event -", event.Type)

		    	switch event.Kind {
					case pointer_gio.Press:
						if event.Priority == pointer_gio.Grabbed {
							log.Println("GoApp.Keyboard().SetFocusControl(nil)")
							GoApp.Keyboard().SetFocusControl(nil)

						}

				}

	    }
	}
	GoApp.Keyboard().Update()
}

func (ob *GoWindowObj) updateLayout(layout GoObject, gtx layout_gio.Context) {
	//log.Println("updateLayout() = ENTRY")
	for _, obj := range layout.Objects() {
		if obj.ObjectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		}/* else {
			if obj.ObjectType() == "GoButtonObj"{
				button := obj.(*GoButtonObj)
				if button.Clicked() {
					log.Println("GoButtonObj:Clicked()")
					//GoApp.Keyboard().SetFocus(nil)
					if button.onRelease != nil{
						button.onRelease()
					}
					if button.onClick != nil{
						button.onClick()
					}
				} 
				if button.Focused() {
					//log.Println("GoButtonObj:Focused()")
					if button.onFocus != nil {
						button.onFocus()
					}
				}
				if button.Hovered() {
					//log.Println("GoButtonObj:Hovered()")
					if button.onHover != nil {
						button.onHover()
					}
				}
				if button.Pressed() {
					log.Println("GoLayout::GoButtonObj:Pressed()")
					log.Println("GoApp.Keyboard().SetFocusControl(button.GoWidget)")
					GoApp.Keyboard().SetFocusControl(&button.GioWidget)
					if button.onPress != nil {
						button.onPress()
					}
				}
			} else if obj.ObjectType() == "GoRadioButtonObj"{
				button := obj.(*GoRadioButtonObj)
				if button.Changed() {
					log.Println("GoRadioButtonObj:Changed()")
					if button.onChange != nil{
						button.onChange()
					}
				} 
				if tag, focus := button.Focused(); focus {
					log.Println("GoRadioButtonObj:Focused()")
					if button.onFocus != nil {
						button.onFocus(tag)
					}
				}
				if tag, hover := button.Hovered(); hover {
					log.Println("GoRadioButtonObj:Hovered()")
					if button.onHover != nil {
						button.onHover(tag)
					}
				}

			} else if obj.ObjectType() == "GoSwitchObj" {
				swtch := obj.(*GoSwitchObj)
				if swtch.Changed() {
					log.Println("GoSwitchObj:Changed()")
					if swtch.onChange != nil {
						swtch.onChange(swtch.goSwitch.Value)
					}
				}
				if swtch.Focused() {
					log.Println("GoSwitchObj:Focused()")
					if swtch.onFocus != nil {
						swtch.onFocus()
					}
				}
				if swtch.Hovered() {
					log.Println("GoSwitchObj:Hovered()")
					if swtch.onHover != nil {
						swtch.onHover()
					}
				}
				if swtch.Pressed() {
					log.Println("GoSwitchObj:Pressed()")
					if swtch.onPress != nil {
						swtch.onPress()
					}
				}
			} else if obj.ObjectType() == "GoSliderObj" {
				slider := obj.(*GoSliderObj)
				if slider.Changed() {
					//log.Println("GoSliderObj:Changed()")
					if slider.onChange != nil {
						slider.onChange(slider.gioSlider.Value)
					}
				}
				if slider.Dragging() {
					//log.Println("GoSliderObj:Dragging()")
					if slider.onDrag != nil {
						slider.onDrag(slider.gioSlider.Value)
					}
				}
			}
		}*/

	}
}