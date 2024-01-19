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
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{space, space, space, space, 16777215, 16777215, space, space},
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
	//log.Println("GoSpacerObj::Draw()")
	cs := gtx.Constraints
	//clipper := gtx.Constraints
	//log.Println("gtx.Constraints Min = (", cs.Min.X, cs.Min.Y, ") Max = (", cs.Max.X, cs.Max.Y, ")")
	
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		//log.Println("FixedWidth............")
		//log.Println("object Width = (", width, " )")
		cs.Min.X = min(cs.Max.X, width)
		//log.Println("cs.Min.X = (", cs.Min.X, " )")
		cs.Max.X = min(cs.Max.X, width)
		//log.Println("cs.Max.X = (", cs.Max.X, " )")
	/*case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = min(cs.Min.X, minWidth)
		cs.Max.X = min(cs.Max.X, maxWidth)*/
	case PreferredWidth:		// SizeHint is Preferred
		//log.Println("PreferredWidth............")
		//log.Println("object MinWidth = (", minWidth, " )")
		//log.Println("object MaxWidth = (", maxWidth, " )")
		cs.Min.X = max(cs.Min.X, minWidth)
		cs.Max.X = min(cs.Max.X, maxWidth)
	/*case MaximumWidth:			// SizeHint is Maximum
		cs.Min.X = max(cs.Min.X, minWidth) 	// No change to gtx.Constraints.X
		cs.Max.X = min(cs.Max.X, maxWidth)*/
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		cs.Min.Y = min(cs.Max.Y, height)
		cs.Max.Y = min(cs.Max.Y, height)
	/*case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = min(cs.Min.Y, ob.MinHeight)
		cs.Max.Y = min(cs.Max.Y, ob.MaxHeight)*/
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = min(cs.Min.Y, minHeight)
		cs.Max.Y = min(cs.Max.Y, maxHeight)
	/*case MaximumHeight:			// SizeHint is Maximum
		cs.Min.Y = min(cs.Min.Y, ob.MinHeight) 	// No change to gtx.Constraints.Y
		cs.Max.Y = min(cs.Max.Y, ob.MaxHeight)*/
	case ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}

	gtx.Constraints = cs
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Min,}

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
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
		//log.Println("SpacerSize:", ob.Width, ob.Height)
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