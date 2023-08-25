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

func (a *GoApplicationObj) addWindow(w *GoWindowObj) {
	a.windows = append(a.windows, w)
}

func (a *GoApplicationObj) Run() {
	var gio *app_gio.Window = nil
	if len(a.windows) == 0 {
		err := errors.New("****************\n\nApplication has no main windows!\n" +
											"Use GoWindow()) method to create new windows.\n\n")
		log.Fatal(err)
	}
	for _, window := range a.windows {
		gio = window.gio
	}
	if gio == nil {
		err := errors.New("****************\n\nApplication has no active main windows!\n" +
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
	gio *app_gio.Window
	name string
	layout *GoLayoutObj

}

func GoWindow(windowName string) (hWin *GoWindowObj) {
	//object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	object := GioObject{nil, nil, []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	hWin = &GoWindowObj{object, nil, windowName, nil}
	hWin.Window = hWin
	hWin.layout = GoVFlexBoxLayout(hWin)
	GoApp.addWindow(hWin)
	return
}

func (ob *GoWindowObj) Layout() *GoLayoutObj {
	return ob.layout
}

func (ob *GoWindowObj) ObjectType() (string) {
	return "GoWindowObj"
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
	if style == HBoxLayout || style == VBoxLayout {
		ob.layout = GoBoxLayout(ob, style)
	} else if style == HFlexBoxLayout || style == VFlexBoxLayout {
		ob.layout = GoFlexBoxLayout(ob, style)
	}
}

func (ob *GoWindowObj) SetMargin(left int, top int, right int, bottom int) {
	ob.layout.SetMargin(left, top, right, bottom)
}

func (ob *GoWindowObj) SetPadding(left int, top int, right int, bottom int) {
	ob.layout.SetPadding(left, top, right, bottom)
}

func (ob *GoWindowObj) SetSpacing(spacing GoLayoutSpacing) {
	ob.layout.SetSpacing(spacing)
}

func (ob *GoWindowObj) Show() {
	ob.run()
}

func (ob *GoWindowObj) run() {
	go func() {
	    // create new window
	    ob.gio = app_gio.NewWindow(
	      app_gio.Title(ob.name),
	      app_gio.Size(unit_gio.Dp(650), unit_gio.Dp(600)),
	    )
	    // draw on screen
	    if err := ob.loop(); err != nil {
	      log.Fatal(err)
		}
		os.Exit(0)
	}()
	time.Sleep(200 * time.Millisecond)
}

func (ob *GoWindowObj) loop() (err error) {
	var count int
	// ops are the operations from the UI
    var ops op.Ops

    // th defines the material design style
    //material_gio.NewTheme(gofont.Collection())
    //th := material.NewTheme(gofont.Collection())

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
	e.Frame(gtx.Ops)
}

func (ob *GoWindowObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	//cs := gtx.Constraints
	//log.Println("Layout Constraints:", cs)
	//log.Println("control length:", len(ob.controls))
	//return layout_gio.Inset{10, 10, 10, 10}.Layout(gtx, func(gtx C) D {
		ob.signalEvents(gtx)

		/*top := gtx.Dp(unit_gio.Dp(ob.layout.GoMargin.Top))
		right := gtx.Dp(unit_gio.Dp(ob.layout.GoMargin.Right))
		bottom := gtx.Dp(unit_gio.Dp(ob.layout.GoMargin.Bottom))
		left := gtx.Dp(unit_gio.Dp(ob.layout.GoMargin.Left))

		log.Println("left: ", left)
		log.Println("top: ", top)
		log.Println("right: ", right)
		log.Println("bottom: ", bottom)
		

		mcs := gtx.Constraints
		mcs.Max.X -= left + right
		if mcs.Max.X < 0 {
			left = 0
			right = 0
			mcs.Max.X = 0
		}
		if mcs.Min.X > mcs.Max.X {
			mcs.Min.X = mcs.Max.X
		}
		mcs.Max.Y -= top + bottom
		if mcs.Max.Y < 0 {
			bottom = 0
			top = 0
			mcs.Max.Y = 0
		}
		if mcs.Min.Y > mcs.Max.Y {
			mcs.Min.Y = mcs.Max.Y
		}
		gtx.Constraints = mcs*/
		log.Println("gtx.Constraints: ", gtx.Constraints)
		dims := ob.layout.Draw(gtx)
		log.Println("dims: ", dims)
		// add the events handler to receive widget pointer events
		
		return dims
	//})
	/*return ob.margin.Layout(gtx, func(gtx C) D {
		return ob.border.Layout(gtx, func(gtx C) D {
			return ob.padding.Layout(gtx, func(gtx C) D {
				return ob.gio.Layout(gtx, len(ob.controls), func(gtx layout_gio.Context, i int) layout_gio.Dimensions {
					//return list.Layout(gtx, 2, func(gtx layout_gio.Context, i int) layout_gio.Dimensions {
					//cs = gtx.Constraints
					log.Println("Object[", i, "].draw")
					return ob.controls[i].draw(gtx)
				})
			})
		})
	})*/
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
	for _, obj := range ob.layout.Controls {
		if obj.ObjectType() == "GoLayoutObj" {
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
		    	log.Println("ApplicationPointer::Event -", event.Type)
		    	switch event.Type {
					case pointer_gio.Press:
						log.Println("GoApp.Keyboard().SetFocusControl(nil)")
						GoApp.Keyboard().SetFocusControl(nil)
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