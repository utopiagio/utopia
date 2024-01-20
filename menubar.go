// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/menubar.go */

package utopia

import (
	"log"
	"image"
	
	layout_gio "github.com/utopiagio/gio/layout"

	"github.com/utopiagio/utopia/metrics"
)

func GoMenuBar(parent GoObject) (hObj *GoMenuBarObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 24, 100, 24, 1000, 24, 100, 24},
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
	layout *GoLayoutObj
	menus []*GoMenuObj
}

func (ob *GoMenuBarObj) AddMenu(text string) (*GoMenuObj){
	menu := GoMenu(ob, text, len(ob.menus))
	ob.menus = append(ob.menus, menu)
	return menu
}


func (ob *GoMenuBarObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	log.Println("GoMenuBarObj::Draw()")
	cs := gtx.Constraints
	log.Println("gtx.Constraints Min = (", cs.Min.X, cs.Min.Y, ") Max = (", cs.Max.X, cs.Max.Y, ")")
	
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		cs.Min.X = min(cs.Max.X, width)			// constrain to ob.Width
		cs.Max.X = min(cs.Max.X, width)			// constrain to ob.Width
	case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = minWidth					// set to ob.MinWidth
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = max(cs.Min.X, minWidth)		// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = maxWidth						// set to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		cs.Min.Y = min(cs.Max.Y, height)		// constrain to ob.Height
		cs.Max.Y = min(cs.Max.Y, height)		// constrain to ob.Height
	case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = minHeight				// set to ob.MinHeight
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = min(cs.Min.Y, minHeight)		// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = maxHeight					// set to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	case ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}
	gtx.Constraints = cs
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
		log.Println("ob.MenuBar height :", ob.AbsHeight)
	}
	log.Println("return dims = (", dims.Size.X, ",", dims.Size.Y, ")")
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