/* spacer.go */

package utopia

import (
	"log"
	"image"
	

	//"github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	//semantic_gio "github.com/utopiagio/gio/io/semantic"
	unit_gio "github.com/utopiagio/gio/unit"
)

func GoSpacer(parent GoObject, space int) (hObj *GoSpacerObj) {
	if parent.ObjectType() != "GoLayoutObj" {
		log.Println("Cannot create GoSpacerObj for", parent.ObjectType())
		return nil
	}
	
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hSpacer := &GoSpacerObj{
		GioObject: object,
		GioWidget: widget,
		
		color: theme.ContrastFg,
		background: theme.ContrastBg,
		space: space,
	}
	if parent.(*GoLayoutObj).Style() == HFlexBoxLayout {
		hSpacer.width = unit_gio.Dp(space)
		hSpacer.height = unit_gio.Dp(0)
	} else if parent.(*GoLayoutObj).Style() == VFlexBoxLayout {
		hSpacer.width = unit_gio.Dp(0)
		hSpacer.height = unit_gio.Dp(space)
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
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}

	//log.Println("gtx.Constraints.Max: ", dims)
	if ob.Visible {
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
				//log.Println("PaddingDims: ", paddingDims)
				return paddingDims
			})
			//log.Println("BorderDims: ", borderDims)
			return borderDims
		})
		ob.dims = dims
		//log.Println("SpacerDims: ", dims)
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
		//log.Println("SpacerSize:", ob.Width, ob.Height)
	}
	return dims
}

func (ob *GoSpacerObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	return layout_gio.Dimensions {
		Size: image.Point{
			X: gtx.Dp(ob.width),
			Y: gtx.Dp(ob.height),
		},
	}
}

func (ob *GoSpacerObj) ObjectType() (string) {
	return "GoSpacerObj"
}

func (ob *GoSpacerObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}