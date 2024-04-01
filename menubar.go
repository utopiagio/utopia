// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/menubar.go */

package utopia

import (
	//"log"
	"image"
	
	layout_gio "github.com/utopiagio/gio/layout"

	"github.com/utopiagio/utopia/metrics"
)

func GoMenuBar(parent GoObject) (hObj *GoMenuBarObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, FixedHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 24, 100, 24, 1000, 24, 100, 24},
		FocusPolicy: NoFocus,
		Visible: false,
		tag: tagCounter,
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
	layout *GoLayoutObj
	menus []*GoMenuObj
}

func (ob *GoMenuBarObj) AddMenu(text string) (*GoMenuObj){
	menu := GoMenu(ob, text, len(ob.menus))
	ob.menus = append(ob.menus, menu)
	return menu
}


func (ob *GoMenuBarObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		ob.repack(gtx)
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoMenuBarObj) MenuOffset(id int) (int) {
	var offset int
	if id > len(ob.menus) {return 0}
	for idx := 0; idx < id; idx++ {
		offset += ob.menus[idx].AbsWidth
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

func (ob *GoMenuBarObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	layoutDims := ob.layout.flex_gio.Layout(gtx, ob.layout.flexControls... )
	dims := image.Point{X: layoutDims.Size.X, Y: layoutDims.Size.Y}
	return layout_gio.Dimensions{Size: dims}
}

func (ob *GoMenuBarObj) repack(gtx layout_gio.Context) {
	ob.layout.flexControls = []layout_gio.FlexChild{}
	for i := 0; i < len(ob.Controls); i++ {
		ob.layout.addFlexControl(ob.Controls[i])
		if ob.Controls[i].ObjectType() == "GoLayoutObj" {
			ob.Controls[i].(*GoLayoutObj).repack(gtx)
		}
	}
}