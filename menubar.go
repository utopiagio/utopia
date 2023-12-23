// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/menubar.go */

package utopia

import (
	//"log"
	"image"
	//"math"
	//"time"

	//"github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	//unit_gio "github.com/utopiagio/gio/unit"

	"github.com/utopiagio/utopia/metrics"
)

func GoMenuBar(parent GoObject) (hObj *GoMenuBarObj) {
	//var fontSize unit_gio.Sp = 14
	//var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 24, 1000, 24, 1000, 24},
		
		FocusPolicy: NoFocus,
		Visible: false,
	}
	hMenuBar := &GoMenuBarObj{
		GioObject: object,
		GioWidget: widget,
		//background: Color_WhiteSmoke,
		//cornerRadius: 4,
	}
	hMenuBar.layout = GoHFlexBoxLayout(hMenuBar)
	parent.AddControl(hMenuBar)

	return hMenuBar
}

type GoMenuBarObj struct {
	GioObject
	GioWidget
	//background GoColor
	//cornerRadius unit_gio.Dp
	layout *GoLayoutObj
	//menuOffset []int
	menus []*GoMenuObj
}

func (ob *GoMenuBarObj) AddMenu(text string) (*GoMenuObj){
	menu := GoMenu(ob, text, len(ob.menus))
	ob.menus = append(ob.menus, menu)
	return menu
}


func (ob *GoMenuBarObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	//var paddingDims layout_gio.Dimensions
	//log.Println("gtx.Constraints.Max: ", gtx.Constraints.Max)
	//dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Min,}

	//log.Println("gtx.Constraints.Max: ", dims)
	if ob.Visible {
		ob.repack(gtx)
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
					barDims := ob.render(gtx)
					//log.Println("BarDims: ", barDims)
					return barDims
				})
				//log.Println("PaddingDims: ", paddingDims)
				return paddingDims
			})
			//log.Println("BorderDims: ", borderDims)
			return borderDims
		})
		//log.Println("GoMenuBar dims: ", dims)
		ob.dims = dims
		ob.Width = metrics.PxToDp(GoDpr, dims.Size.X)	//(int(float32(dims.Size.X) / GoDpr))
		ob.Height = metrics.PxToDp(GoDpr, dims.Size.Y)	//(int(float32(dims.Size.Y) / GoDpr))
		//log.Println("GoMenuBar Size:", ob.Width, ",", ob.Height)
	}
	return dims
}

func (ob *GoMenuBarObj) MenuOffset(id int) (int) {
	var offset int
	if id > len(ob.menus) {return 0}
	for idx := 0; idx < id; idx++ {
		offset += ob.menus[idx].Width
	}
	return offset
}

func (ob *GoMenuBarObj) ObjectType() (string) {
	return "GoMenuBarObj"
}

func (ob *GoMenuBarObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoMenuBarObj) SetFixedHeight(height int) {
	ob.Height = height
	ob.MinHeight = height
	ob.MaxHeight = height
}

func (ob *GoMenuBarObj) render(gtx layout_gio.Context) layout_gio.Dimensions {
	//log.Println("MenuBar Height: ", ob.Height())
	//log.Println("MenuBar Width: ", ob.Width())
	/*height := gtx.Dp(unit_gio.Dp(ob.Height))
	//height := gtx.Dp(ob.Height)

	//log.Println("MenuBar height: ", height)
	width := gtx.Dp(unit_gio.Dp(ob.Width))
	//width := gtx.Dp(ob.Width)
	//log.Println("MenuBar width: ", width)
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
		//log.Println("MenuBar width: ", width)
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
		//log.Println("MenuBar height: ", height)
	}*/
	// Save the operations in an independent ops value (the cache).
	//macro := op_gio.Record(gtx.Ops)
	layoutDims := ob.layout.flex_gio.Layout(gtx, ob.layout.flexControls... )
	//call := macro.Stop()
	dims := image.Point{X: layoutDims.Size.X, Y: layoutDims.Size.Y}
	//r := dims
	
	//r.X += ob.GoPadding.Left + ob.GoPadding.Right
	//r.Y += ob.GoPadding.Top + ob.GoPadding.Bottom
	
	//rr := gtx.Dp(ob.cornerRadius)
	//defer clip_gio.UniformRRect(image.Rectangle{Max: r}, rr).Push(gtx.Ops).Pop()

	// paint background
	//background := ob.background.NRGBA()
	//paint_gio.Fill(gtx.Ops, background)
	// Draw the operations from the cache.
	//call.Add(gtx.Ops)
	//layoutDims := ob.layout.flex_gio.Layout(gtx, ob.layout.flexControls... )
	return layout_gio.Dimensions{Size: dims}
}

func (ob *GoMenuBarObj) repack(gtx layout_gio.Context) {
	ob.layout.flexControls = []layout_gio.FlexChild{}
	for i := 0; i < len(ob.Controls); i++ {
		
		ob.layout.addFlexControl(ob.Controls[i])
		/*if ob.controls[i].sizePolicy().HFixed {
			ob.addRigidControl(ob.controls[i])
		} else {
			ob.addFlexedControl(ob.controls[i])
		}*/
		if ob.Controls[i].ObjectType() == "GoLayoutObj" {
			ob.Controls[i].(*GoLayoutObj).repack(gtx)
		}
	}
}