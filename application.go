/* application.go */

package utopia

import (
	"errors"
	"log"
	"time"
	"os"
	app_gio "github.com/utopiagio/gio/app"
	"github.com/utopiagio/gio/font/gofont"
	"github.com/utopiagio/gio/io/key"
	"github.com/utopiagio/gio/io/pointer"
	"github.com/utopiagio/gio/io/system"
	layout_gio "github.com/utopiagio/gio/layout"
	"github.com/utopiagio/gio/op"
	text_gio "github.com/utopiagio/gio/text"
	"github.com/utopiagio/gio/unit"
	_ "github.com/utopiagio/gio/widget"
)

type (
	D = layout_gio.Dimensions
	C = layout_gio.Context
)

var goApp *GoApplicationObj = nil

type GoApplicationObj struct {
	name string
	windows	[]*GoWindowObj
	// Theme contains semantic style data. Extends `material.Theme`.
	theme *GoThemeObj
	//theme *material_gio.Theme
	// Shaper cache of registered fonts.
	shaper *text_gio.Shaper
	//fontCollection []text_gio.FontFace
}

func GoApplication(appName string) (a *GoApplicationObj) {
	theme := GoTheme(gofont.Collection())
	//theme := material_gio.NewTheme(gofont.Collection())
	//sh := text_gio.NewShaper(gofont.Collection())
	goApp = &GoApplicationObj{
		name: appName,
		theme: theme,
		//fontCollection: gofont.Collection(),
		
	}
	return goApp
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
	goObject
	//goWidget
	gio *app_gio.Window
	name string
	layout *GoLayoutObj

}

func GoWindow(windowName string) (hWin *GoWindowObj) {
	object := goObject{nil, nil, []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	hWin = &GoWindowObj{object, nil, windowName, nil}
	hWin.goObject.window = hWin
	hWin.layout = GoVFlexBoxLayout(hWin)
	goApp.addWindow(hWin)
	return
}

func (ob *GoWindowObj) Layout() *GoLayoutObj {
	return ob.layout
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

func (ob *GoWindowObj) Show() {
	ob.run()
}

func (ob *GoWindowObj) run() {
	go func() {
	    // create new window
	    ob.gio = app_gio.NewWindow(
	      app_gio.Title(ob.name),
	      app_gio.Size(unit.Dp(650), unit.Dp(600)),
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
	      		log.Println("Window.update(gtx).....", count)
	      		count++
	      		ob.update(gtx)
	      		//log.Println("Window.render(gtx).....")
	      		ob.render(gtx)
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

func (ob *GoWindowObj) objectType() (string) {
	return "GoWindowObj"
}

func (ob *GoWindowObj) paint(e system.FrameEvent, gtx layout_gio.Context) {
	e.Frame(gtx.Ops)
}

func (ob *GoWindowObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	//cs := gtx.Constraints
	//log.Println("Layout Constraints:", cs)
	//log.Println("control length:", len(ob.controls))
	//return layout_gio.Inset{10, 10, 10, 10}.Layout(gtx, func(gtx C) D {
		return ob.layout.draw(gtx)
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

func (ob *GoWindowObj) update(gtx layout_gio.Context) {
	//log.Println("(ob *GoWindowObj) update.............")
	for _, obj := range ob.layout.controls {
		if obj.objectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		} else {
			if obj.objectType() == "GoButtonObj"{
				button := obj.(*GoButtonObj)
				if button.Clicked() {
					log.Println("GoButtonObj:Clicked()")
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
					//log.Println("GoButtonObj:Pressed()")
					if button.onPress != nil {
						button.onPress()
					}
				}
			} else if obj.objectType() == "GoRadioButtonObj"{
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

			} else if obj.objectType() == "GoSwitchObj" {
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
			} else if obj.objectType() == "GoSliderObj" {
				slider := obj.(*GoSliderObj)
				if slider.Changed() {
					log.Println("GoSliderObj:Changed()")
					if slider.onChange != nil {
						slider.onChange(slider.gioSlider.Value)
					}
				}
				if slider.Dragging() {
					log.Println("GoSliderObj:Dragging()")
					if slider.onDrag != nil {
						slider.onDrag(slider.gioSlider.Value)
					}
				}
			} //else if obj.objectType() == "GoTextEditObj" {
				//textedit := obj.(*GoTextEditObj)
		}
	}
	for _, gtxEvent := range gtx.Events(0) {
	    switch gtxEvent.(type) { //gtxE := gtxEvent.(type) {

		    case key.EditEvent:

		    case key.Event:
		      
		    case pointer.Event:

	    }
	}
}

func (ob *GoWindowObj) updateLayout(layout GoObject, gtx layout_gio.Context) {
	//log.Println("updateLayout() = ENTRY")
	for _, obj := range layout.objects() {
		if obj.objectType() == "GoLayoutObj" {
			ob.updateLayout(obj, gtx)
		} else {
			if obj.objectType() == "GoButtonObj"{
				button := obj.(*GoButtonObj)
				if button.Clicked() {
					log.Println("GoButtonObj:Clicked()")
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
					//log.Println("GoButtonObj:Pressed()")
					if button.onPress != nil {
						button.onPress()
					}
				}
			} else if obj.objectType() == "GoRadioButtonObj"{
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

			} else if obj.objectType() == "GoSwitchObj" {
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
			} else if obj.objectType() == "GoSliderObj" {
				slider := obj.(*GoSliderObj)
				if slider.Changed() {
					log.Println("GoSliderObj:Changed()")
					if slider.onChange != nil {
						slider.onChange(slider.gioSlider.Value)
					}
				}
				if slider.Dragging() {
					log.Println("GoSliderObj:Dragging()")
					if slider.onDrag != nil {
						slider.onDrag(slider.gioSlider.Value)
					}
				}
			}
		}

	}
}