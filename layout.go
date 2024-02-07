// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/layout.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"

	layout_gio "github.com/utopiagio/gio/layout"
	//widget_gio "github.com/utopiagio/gio/widget"
	unit_gio "github.com/utopiagio/gio/unit"
	"github.com/utopiagio/utopia/metrics"
)

type GoLayoutStyle int

const (
	NoLayout 	GoLayoutStyle = iota
	HBoxLayout 							// gio.List{Axis: layout_gio.Horizontal}	
	VBoxLayout 							// gio.List{Axis: layout_gio.Vertical}
	HVBoxLayout							// Not Implemented *******************
	HFlexBoxLayout
	// gio.Flex{Axis: layout_gio.Horizontal, Spacing: 0, Alignment: Baseline, WeightSum: 0}
	VFlexBoxLayout						
	// gio.Flex{Axis: layout_gio.Vertical, Spacing: 0, Alignment: Baseline, WeightSum: 0}
	PopupMenuLayout

)

type GoLayoutDirection int

const (
	Horizontal = 0
	Vertical = 1
)

type GoLayoutSpacing uint8

const (
	// SpaceEnd leaves space at the end.
	SpaceEnd GoLayoutSpacing = iota
	// SpaceStart leaves space at the start.
	SpaceStart
	// SpaceSides shares space between the start and end.
	SpaceSides
	// SpaceAround distributes space evenly between children,
	// with half as much space at the start and end.
	SpaceAround
	// SpaceBetween distributes space evenly between children,
	// leaving no space at the start and end.
	SpaceBetween
	// SpaceEvenly distributes space evenly between children and
	// at the start and end.
	SpaceEvenly
)

type GoLayoutAlignment uint8

const (
	LayoutStart GoLayoutAlignment = iota
	LayoutEnd
	LayoutMiddle
	LayoutBaseline
)

func GoLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	switch style {
		case NoLayout:
			return GoBoxLayout(parent, NoLayout)
		case HBoxLayout:
			return GoHBoxLayout(parent)
		case VBoxLayout:
			return GoVBoxLayout(parent)
		case HFlexBoxLayout:
			return GoHFlexBoxLayout(parent)
		case VFlexBoxLayout:
			return GoVFlexBoxLayout(parent)
		default:
			return GoVBoxLayout(parent)
	}
}

func GoBoxLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	var axis layout_gio.Axis
	switch style {
		case NoLayout:
			axis = layout_gio.Horizontal
		case HBoxLayout:
			axis = layout_gio.Horizontal
		case VBoxLayout:
			axis = layout_gio.Vertical
	}
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		list_gio: &layout_gio.List{Axis: axis},
		style: style,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoHBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		list_gio: &layout_gio.List{Axis: layout_gio.Horizontal},
		style: HBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoVBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		list_gio: &layout_gio.List{Axis: layout_gio.Vertical},
		style: VBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoFlexBoxLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	var axis layout_gio.Axis
	switch style {
		case NoLayout:
			axis = layout_gio.Horizontal
		case HFlexBoxLayout:
			axis = layout_gio.Horizontal
		case VFlexBoxLayout:
			axis = layout_gio.Vertical
	}
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: axis},
		style: style,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoHFlexBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: layout_gio.Horizontal},
		style: HFlexBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoVFlexBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: layout_gio.Vertical},
		style: VFlexBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoPopupMenuLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: layout_gio.Vertical},
		style: PopupMenuLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

type GoLayoutObj struct {
	GioObject
	GioWidget
	list_gio 	*layout_gio.List
	flex_gio 	*layout_gio.Flex
	style 		GoLayoutStyle
	flexControls 	[]layout_gio.FlexChild
}

func (ob *GoLayoutObj) SetAlignment(alignment GoLayoutAlignment) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flex_gio.Alignment = layout_gio.Alignment(uint8(alignment))	// layout_gio.Alignment
	} else if ob.style == HBoxLayout || ob.style == VBoxLayout || ob.style == PopupMenuLayout {
		ob.list_gio.Alignment = layout_gio.Alignment(uint8(alignment))	// layout_gio.Alignment
	}
}


func (ob *GoLayoutObj) SetSpacing(spacing GoLayoutSpacing) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout || ob.style == PopupMenuLayout {
		ob.flex_gio.Spacing = layout_gio.Spacing(uint8(spacing))	// layout_gio.Spacing
	}
}

func (ob *GoLayoutObj) Style() (GoLayoutStyle) {
	return ob.style
}

func (ob *GoLayoutObj) addFlexControl(control GoObject) {
	if ob.style == HFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.FlexControl(control.SizePolicy().HFlex, control.SizePolicy().VFlex, 1, control.Draw))
	} else if ob.style == VFlexBoxLayout || ob.style == PopupMenuLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.FlexControl(control.SizePolicy().VFlex, control.SizePolicy().HFlex, 1, control.Draw))
	}
}

func (ob *GoLayoutObj) addFlexedControl(control GoObject) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout || ob.style == PopupMenuLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Flexed(1, control.Draw))
	}
}

func (ob *GoLayoutObj) addRigidControl(control GoObject) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout || ob.style == PopupMenuLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Rigid(control.Draw))
	}
}

func (ob *GoLayoutObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		if ob.style == HBoxLayout || ob.style == VBoxLayout {
			//log.Println("BoxLayout style:", ob.style)
			dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
				borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
					paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
						layoutDims := ob.list_gio.Layout(gtx, len(ob.Controls), func(gtx C, i int) D {
							//log.Println("Layout Control idx: ", i)
							return ob.Controls[i].Draw(gtx)
						})
						//log.Println("Layout LayoutDims: ", layoutDims)
						return layoutDims
					})
					//log.Println("Layout PaddingDims: ", paddingDims)
					return paddingDims
				})
				//log.Println("Layout BorderDims: ", borderDims)
				return borderDims
			})
			//log.Println("Layout MarginDims: ", dims)
		} else if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
			//log.Println("FlexBoxLayout style:", ob.style)
			ob.repack(gtx)
			dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
				borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
					paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
						layoutDims := ob.flex_gio.Layout(gtx, ob.flexControls... )
						return layoutDims
					})
					return paddingDims
				})
				return borderDims
			})
		}  else if ob.style == PopupMenuLayout {
			//log.Println("PopupMenuLayout style:", ob.style)
			ob.repack(gtx)
			dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
				borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
					paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
						layoutDims := ob.flex_gio.Layout(gtx, ob.flexControls... )
						layoutDims.Size.X = gtx.Dp(unit_gio.Dp(ob.MinWidth))
						layoutDims.Size.Y = gtx.Dp(unit_gio.Dp(ob.MinHeight))
						return layoutDims
					})
					//log.Println("Layout PaddingDims: ", paddingDims)
					return paddingDims
				})
				//log.Println("Layout BorderDims: ", borderDims)
				return borderDims
			})
		}
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoLayoutObj) ObjectType() (string) {
	return "GoLayoutObj"
}

func (ob *GoLayoutObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoLayoutObj) repack(gtx layout_gio.Context) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = []layout_gio.FlexChild{}
		for i := 0; i < len(ob.Controls); i++ {
			ob.addFlexControl(ob.Controls[i])
			if ob.Controls[i].ObjectType() == "GoLayoutObj" {
				ob.Controls[i].(*GoLayoutObj).repack(gtx)
			}
		}
	} else if ob.style == PopupMenuLayout {
		ob.flexControls = []layout_gio.FlexChild{}
		ob.Width = 0
		ob.Height = 0
		ob.MinWidth = 0
		ob.MinHeight = 0
		for i := 0; i < len(ob.Controls); i++ {
			ob.addFlexControl(ob.Controls[i])
			dims := ob.Controls[i].(*GoMenuItemObj).CalcSize(gtx)
			//ob.Height += metrics.PxToDp(GoDpr, dims.Size.Y)
			//log.Println("dims.Size.Y :", dims.Size.Y)
			ob.Height += dims.Size.Y - 1 		// ******* why (-1) *******
			menuItemWidth := metrics.PxToDp(GoDpr, dims.Size.X)
			if menuItemWidth > ob.Width {
				ob.Width = menuItemWidth
				ob.MinWidth = menuItemWidth
			}
		}
		ob.Height = metrics.PxToDp(GoDpr, ob.Height)
		ob.MinHeight = ob.Height
		
		for i := 0; i < len(ob.Controls); i++ {
			ob.Controls[i].(*GoMenuItemObj).MaxWidth = ob.Width
		}
	}
}
