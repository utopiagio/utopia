// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/layout.go */

package utopia

import (
	//"log"
	"image"
	"math"

	layout_gio "github.com/utopiagio/gio/layout"
	widget_gio "github.com/utopiagio/gio/widget"
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
	theme := GoApp.Theme()
	lightFg := theme.ColorFg.NRGBA()
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	state := &widget_gio.List{}
	state.Axis = axis
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
	}
	scrollbar := GoScrollbar{
		Scrollbar: &state.Scrollbar,
		Track: ScrollTrackStyle{
			MajorPadding: 2,
			MinorPadding: 2,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  8,
			MinorWidth:   6,
			CornerRadius: 3,
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		AnchorStrategy: Occupy,
		list_gio: state,
		style: style,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoHBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	theme := GoApp.Theme()
	lightFg := theme.ColorFg.NRGBA()
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	state := &widget_gio.List{}
	state.Axis = layout_gio.Horizontal
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
	}
	scrollbar := GoScrollbar{
		Scrollbar: &state.Scrollbar,
		Track: ScrollTrackStyle{
			MajorPadding: 2,
			MinorPadding: 2,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  8,
			MinorWidth:   6,
			CornerRadius: 3,
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		AnchorStrategy: Occupy,
		list_gio: state,
		//list_gio: &layout_gio.List{Axis: layout_gio.Horizontal},
		style: HBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.AddControl(hLayout)
	return hLayout
}

func GoVBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	theme := GoApp.Theme()
	lightFg := theme.ColorFg.NRGBA()
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	state := &widget_gio.List{}
	state.Axis = layout_gio.Vertical
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
	}
	scrollbar := GoScrollbar{
		Scrollbar: &state.Scrollbar,
		Track: ScrollTrackStyle{
			MajorPadding: 2,
			MinorPadding: 2,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  8,
			MinorWidth:   6,
			CornerRadius: 3,
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
	hLayout := &GoLayoutObj{
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		AnchorStrategy: Occupy,
		list_gio: state,
		//list_gio: &layout_gio.List{Axis: layout_gio.Vertical},
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
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
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
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
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
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
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
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
		FocusPolicy: NoFocus,
		Visible: true,
		tag: tagCounter,
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
	Scrollbar GoScrollbar
	AnchorStrategy
	list_gio 	*widget_gio.List
	flex_gio 	*layout_gio.Flex
	style 		GoLayoutStyle
	flexControls 	[]layout_gio.FlexChild
}

// ScrollBy scrolls the list by a relative amount of items.
// Fractional scrolling may be inaccurate for items of differing
// dimensions. This includes scrolling by integer amounts if the current
// l.Position.Offset is non-zero.
func (ob *GoLayoutObj) ScrollBy(num float32) {	// num listItem.offset:dx
	if ob.style == HBoxLayout || ob.style == VBoxLayout {
		ob.list_gio.ScrollBy(num)
	}
}

func (ob *GoLayoutObj) ScrollToOffset(dx int) {		// dx offset - pixels
	if ob.style == HBoxLayout || ob.style == VBoxLayout {
		ob.list_gio.ScrollToOffset(dx)
	}
}

func (ob *GoLayoutObj) ScrollTo(num int) {			// num - listItem
	if ob.style == HBoxLayout || ob.style == VBoxLayout {
		ob.list_gio.ScrollTo(num)
	}
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
	//log.Println("GoLayoutObj::Draw()")
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		if ob.style == HBoxLayout || ob.style == VBoxLayout {
			//log.Println("BoxLayout style:", ob.style)
			dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
				borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
					paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
						layoutDims := ob.Layout(gtx, len(ob.Controls), func(gtx C, i int) D {	
							return ob.Controls[i].Draw(gtx)
						})
						//log.Println("Layout BeforeEnd", ob.list_gio.Position.BeforeEnd)
						//log.Println("Layout First", ob.list_gio.Position.First)
						//log.Println("Layout Offset", ob.list_gio.Position.Offset)
						//log.Println("Layout OffsetLast", ob.list_gio.Position.OffsetLast)
						//log.Println("Layout Count", ob.list_gio.Position.Count)
						//log.Println("Layout Length", ob.list_gio.Position.Length)
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

// layout the list and its scrollbar.
func (ob *GoLayoutObj) Layout(gtx layout_gio.Context, length int, w layout_gio.ListElement) layout_gio.Dimensions {
	//log.Println("GoLayoutObj::Layout()")
	originalConstraints := gtx.Constraints

	// Determine how much space the scrollbar occupies.
	barWidth := gtx.Dp(ob.Scrollbar.Width())

	if ob.AnchorStrategy == Occupy {

		// Reserve space for the scrollbar using the gtx constraints.
		max := ob.list_gio.Axis.Convert(gtx.Constraints.Max)
		min := ob.list_gio.Axis.Convert(gtx.Constraints.Min)
		max.Y -= barWidth
		if max.Y < 0 {
			max.Y = 0
		}
		min.Y -= barWidth
		if min.Y < 0 {
			min.Y = 0
		}
		gtx.Constraints.Max = ob.list_gio.Axis.Convert(max)
		gtx.Constraints.Min = ob.list_gio.Axis.Convert(min)
	}

	listDims := ob.list_gio.List.Layout(gtx, length, w)
	gtx.Constraints = originalConstraints

	// Draw the scrollbar.
	anchoring := layout_gio.E
	if ob.list_gio.Axis == layout_gio.Horizontal {
		anchoring = layout_gio.S
	}
	majorAxisSize := ob.list_gio.Axis.Convert(listDims.Size).X
	start, end := ob.fromListPosition(ob.list_gio.Position, length, majorAxisSize)
	// layout.Direction respects the minimum, so ensure that the
	// scrollbar will be drawn on the correct edge even if the provided
	// layout.Context had a zero minimum constraint.
	gtx.Constraints.Min = listDims.Size
	if ob.AnchorStrategy == Occupy {
		min := ob.list_gio.Axis.Convert(gtx.Constraints.Min)
		min.Y += barWidth
		gtx.Constraints.Min = ob.list_gio.Axis.Convert(min)
	}
	anchoring.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		return ob.Scrollbar.Layout(gtx, ob.list_gio.Axis, start, end)
	})

	if delta := ob.list_gio.ScrollDistance(); delta != 0 {
		// Handle any changes to the list position as a result of user interaction
		// with the scrollbar.
		ob.list_gio.List.Position.Offset += int(math.Round(float64(float32(ob.list_gio.Position.Length) * delta)))

		// Ensure that the list pays attention to the Offset field when the scrollbar drag
		// is started while the bar is at the end of the list. Without this, the scrollbar
		// cannot be dragged away from the end.
		ob.list_gio.List.Position.BeforeEnd = true
	}

	if ob.AnchorStrategy == Occupy {
		// Increase the width to account for the space occupied by the scrollbar.
		cross := ob.list_gio.Axis.Convert(listDims.Size)
		cross.Y += barWidth
		listDims.Size = ob.list_gio.Axis.Convert(cross)
	}
	//log.Println("listDims :", listDims)
	if ob.MinWidth > listDims.Size.X {
		listDims.Size.X = ob.MinWidth
	}
	if ob.MinHeight > listDims.Size.Y {
		listDims.Size.Y = ob.MinHeight
	}
	return listDims
}


func (ob *GoLayoutObj) ObjectType() (string) {
	return "GoLayoutObj"
}

func (ob *GoLayoutObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

// fromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list
// in order to do this. The returned values will be in the range [0,1], and
// start will be less than or equal to end.
/*func (ob *GoLayoutObj) fromListPosition(lp layout_gio.Position, elements int, majorAxisSize int) (start, end float32) {
	// Approximate the size of the scrollable content.
	lengthPx := float32(lp.Length)
	meanElementHeight := lengthPx / float32(elements)

	// Determine how much of the content is visible.
	listOffsetF := float32(lp.Offset)
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthPx

	// Compute the location of the beginning of the viewport.
	viewportStart := (float32(lp.First)*meanElementHeight + listOffsetF) / lengthPx

	return viewportStart, clamp1(viewportStart + visibleFraction)
}*/

// fromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list
// in order to do this. The returned values will be in the range [0,1], and
// start will be less than or equal to end.
func (ob *GoLayoutObj) fromListPosition(lp layout_gio.Position, elements int, majorAxisSize int) (start, end float32) {
	// Approximate the size of the scrollable content.
	lengthEstPx := float32(lp.Length)
	elementLenEstPx := lengthEstPx / float32(elements)

	// Determine how much of the content is visible.
	listOffsetF := float32(lp.Offset)
	listOffsetL := float32(lp.OffsetLast)

	// Compute the location of the beginning of the viewport using estimated element size and known
	// pixel offsets.
	viewportStart := clamp1((float32(lp.First)*elementLenEstPx + listOffsetF) / lengthEstPx)
	viewportEnd := clamp1((float32(lp.First+lp.Count)*elementLenEstPx + listOffsetL) / lengthEstPx)
	viewportFraction := viewportEnd - viewportStart

	// Compute the expected visible proportion of the list content based solely on the ratio
	// of the visible size and the estimated total size.
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthEstPx

	// Compute the error between the two methods of determining the viewport and diffuse the
	// error on either end of the viewport based on how close we are to each end.
	err := visibleFraction - viewportFraction
	adjStart := viewportStart
	adjEnd := viewportEnd
	if viewportFraction < 1 {
		startShare := viewportStart / (1 - viewportFraction)
		endShare := (1 - viewportEnd) / (1 - viewportFraction)
		startErr := startShare * err
		endErr := endShare * err

		adjStart -= startErr
		adjEnd += endErr
	}
	return adjStart, adjEnd
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
