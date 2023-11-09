/* application.go */

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
)

type (
	D = layout_gio.Dimensions
	C = layout_gio.Context
)

var GoDpr float32
var GoApp *GoApplicationObj = nil

type GoApplicationObj struct {
	name string
	windows	[]*GoWindowObj
	clipboard *GoClipBoardObj
	keyboard *GoKeyboardObj
	// Theme contains semantic style data. Extends `material.Theme`.
	theme *GoThemeObj
	//theme *material_gio.Theme
	// Shaper cache of registered fonts.
	shaper *text_gio.Shaper
	//fontCollection []text_gio.FontFace
	dpr float32
}

func GoApplication(appName string) (a *GoApplicationObj) {
	clipboard := GoClipBoard()
	if clipboard.init() != nil {
		log.Println("ClipBoard Not Available!")
	}
	keyboard := GoKeyboard()
	/*if keyboard.init() != nil {
		log.Println("Keyboard Not Available!")
	}*/
	theme := GoTheme(gofont.Collection())
	GoApp = &GoApplicationObj{
		name: appName,
		clipboard: clipboard,
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
	name string
	frame *GoLayoutObj
	menubar *GoMenuBarObj
	statusbar *GoLayoutObj
	layout *GoLayoutObj
	//mainwindow bool
	modalwindow bool
	popupmenus []*GoPopupMenuObj
	popupwindow *GoPopupWindowObj

}

func GoMainWindow(windowName string) (hWin *GoWindowObj) {
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{640, 450, 0, 0, 1500, 1000}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowName, nil, nil, nil, nil, false, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	//hWin.menubar.SetBackgroundColor(Color_Gray)
	//hWin.menubar.SetBorder(BorderSingleLine, 5, 5, Color_Red)
	hWin.layout = GoVFlexBoxLayout(hWin.frame)
	//hWin.AddPopupMenu(GoPopupMenu(hWin))
	hWin.popupwindow = GoPopupWindow(hWin)
	GoApp.AddWindow(hWin)
	return hWin
}

func GoWindow(windowName string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	size := GoSize{640, 450, 0, 0, 1500, 1000}
	pos := GoPos{-1, -1}
	hWin = &GoWindowObj{object, size, pos, nil, windowName, nil, nil, nil, nil, false, nil, nil}
	hWin.Window = hWin
	hWin.frame = GoVFlexBoxLayout(hWin)
	
	hWin.menubar = GoMenuBar(hWin.frame)
	hWin.menubar.SetSizePolicy(ExpandingWidth, FixedHeight)
	//hWin.menubar.SetBackgroundColor(Color_Gray)
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
	ob.popupmenus = nil
}

func (ob *GoWindowObj) IsMainWindow() bool {
	//return ob.frame
	return !ob.modalwindow
}

func (ob *GoWindowObj) Layout() *GoLayoutObj {
	//return ob.frame
	return ob.layout
}

func (ob *GoWindowObj) MenuBar() *GoMenuBarObj {
	//return ob.frame
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

func (ob *GoWindowObj) SetPadding(left int, top int, right int, bottom int) {
	ob.layout.SetPadding(left, top, right, bottom)
}

func (ob *GoWindowObj) SetPos(x int, y int) {
	ob.X = x
	ob.Y = y
}

func (ob *GoWindowObj) SetSize(width int, height int) {
	ob.Width = width
	ob.Height = height
}

func (ob *GoWindowObj) SetSpacing(spacing GoLayoutSpacing) {
	ob.layout.SetSpacing(spacing)
}

func (ob *GoWindowObj) Show() {
	ob.run()
}

func (ob *GoWindowObj) ShowModal() {
	ob.modalwindow = true
	ob.run()
}



func (ob *GoWindowObj) run() {
	go func() {
	    // create new window
	    ob.gio = app_gio.NewWindow(
	      app_gio.Title(ob.name),
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

func (ob *GoWindowObj) loop() (err error) {
	var count int
	// ops are the operations from the UI
    var ops op.Ops

    // listen for events in the window.
    for {
		select {
    	case e := <-ob.gio.Events():

			// detect what type of event
			switch  e := e.(type) {
	      	case system.DestroyEvent:
	      		log.Println("system.DestroyEvent.....")

	      		return e.Err
	      	// this is sent when the application should re-render.
	      	case system.FrameEvent:
	      		// Open an new context
	      		gtx := layout_gio.NewContext(&ops, e)
	      		//log.Println("Window.update(gtx).....", count)
	      		count++
	      		ob.update(gtx)	// receiveEvents
	      		//log.Println("Window.render(gtx).....")
	      		ob.render(gtx)	// draw layout and signalEvents
	      		// window paint
	      		//e.Frame(gtx.Ops)
	      		//log.Println("Window.paint(e, gtx).....")
	      		ob.paint(e, gtx)
	      	}
	    /*case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}
			ob.gio.Invalidate()			// redraw window*/
		}
    }
	return nil
}

func (ob *GoWindowObj) paint(e system.FrameEvent, gtx layout_gio.Context) {
	//log.Println("GoWindow.paint(e, gtx)")
	e.Frame(gtx.Ops)
}

func (ob *GoWindowObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	// set global screen pixel size
	GoDpr = gtx.Metric.PxPerDp
	//log.Println("GoDpr =", GoDpr)
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
			Types: pointer_gio.Press,
		}.Add(gtx.Ops)
}

func (ob *GoWindowObj) update(gtx layout_gio.Context) {

	//log.Println("(ob *GoWindowObj) update.............")
	for _, obj := range ob.frame.Controls {
		if obj.ObjectType() == "GoLayoutObj" {
			//log.Println("(ob *GoLayoutObj) updateLayout.............")
			ob.updateLayout(obj, gtx)
		} else {
			if obj.ObjectType() == "GoButtonObj"{
				/*button := obj.(*GoButtonObj)
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
				}*/
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
				/*swtch := obj.(*GoSwitchObj)
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
				}*/
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
		}
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

		    	switch event.Type {
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
		} else {
			if obj.ObjectType() == "GoButtonObj"{
				/*button := obj.(*GoButtonObj)
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
				}*/
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
				/*swtch := obj.(*GoSwitchObj)
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
				}*/
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
		}

	}
}