// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/statusbar.go */

package utopia

import (
	//"log"
	"image"
	
	layout_gio "github.com/utopiagio/gio/layout"

	"github.com/utopiagio/utopia/metrics"
)

func GoStatusBar(parent GoObject) (hObj *GoStatusBarObj) {
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
	hStatusBar := &GoStatusBarObj{
		GioObject: object,
		GioWidget: widget,
		//background: Color_WhiteSmoke,
		//cornerRadius: 4,
	}
	hStatusBar.layout = GoHFlexBoxLayout(hStatusBar)
	parent.AddControl(hStatusBar)

	return hStatusBar
}

type GoStatusBarObj struct {
	GioObject
	GioWidget
	layout *GoLayoutObj
}

func (ob *GoStatusBarObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

func (ob *GoStatusBarObj) ObjectType() (string) {
	return "GoStatusBarObj"
}

func (ob *GoStatusBarObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoStatusBarObj) SetFixedHeight(height int) {
	ob.Height = height
	ob.MinHeight = height
	ob.MaxHeight = height
}

func (ob *GoStatusBarObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	layoutDims := ob.layout.flex_gio.Layout(gtx, ob.layout.flexControls... )
	dims := image.Point{X: layoutDims.Size.X, Y: layoutDims.Size.Y}
	return layout_gio.Dimensions{Size: dims}
}

func (ob *GoStatusBarObj) repack(gtx layout_gio.Context) {
	ob.layout.flexControls = []layout_gio.FlexChild{}
	for i := 0; i < len(ob.Controls); i++ {
		ob.layout.addFlexControl(ob.Controls[i])
		if ob.Controls[i].ObjectType() == "GoLayoutObj" {
			ob.Controls[i].(*GoLayoutObj).repack(gtx)
		}
	}
}