// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/spacer.go */

package utopia

import (
	"log"
	"image"
	

	//"github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	//semantic_gio "github.com/utopiagio/gio/io/semantic"
	unit_gio "github.com/utopiagio/gio/unit"

	"github.com/utopiagio/utopia/metrics"
)

func GoSpacer(parent GoObject, space int) (hObj *GoSpacerObj) {
	if parent.ObjectType() != "GoLayoutObj" {
		log.Println("Cannot create GoSpacerObj for", parent.ObjectType())
		return nil
	}
	
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{space, space, space, space, 16777215, 16777215, space, space},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
	}
	hSpacer := &GoSpacerObj{
		GioObject: object,
		GioWidget: widget,
		
		color: theme.ContrastFg,
		background: theme.ContrastBg,
		space: space,
	}
	if parent.(*GoLayoutObj).Style() == HFlexBoxLayout {
		hSpacer.Width = space
		hSpacer.Height = 0
		//hSpacer.width = unit_gio.Dp(space)
		//hSpacer.height = unit_gio.Dp(0)
	} else if parent.(*GoLayoutObj).Style() == VFlexBoxLayout {
		hSpacer.Width = 0
		hSpacer.Height = space
		//hSpacer.width = unit_gio.Dp(0)
		//hSpacer.height = unit_gio.Dp(space)
	}
	
	parent.AddControl(hSpacer)
	return hSpacer
}

type GoSpacerObj struct {
	GioObject
	GioWidget
	//theme *GoThemeObj
	color GoColor
	background GoColor
	space int

	height unit_gio.Dp
	width unit_gio.Dp
}

func (ob *GoSpacerObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
				return paddingDims
			})
			return borderDims
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoSpacerObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	width := gtx.Dp(unit_gio.Dp(ob.MinWidth))
	height := gtx.Dp(unit_gio.Dp(ob.MinHeight))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	return layout_gio.Dimensions {
		Size: image.Point{
			X: width,
			Y: height,
		},
	}
}

func (ob *GoSpacerObj) ObjectType() (string) {
	return "GoSpacerObj"
}

func (ob *GoSpacerObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}