// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	//"log"
	//"image"
	//"image/color"
	"math"

	//pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	//unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

// ListStyle configures the presentation of a layout_gio.List with a scrollbar.
type GoListViewObj struct {
	GioObject
	GioWidget
	Scrollbar goScrollbar
	AnchorStrategy
	state *widget_gio.List
	columns int
}

// List constructs a ListStyle using the provided theme and state.
func GoListView(parent GoObject) *GoListViewObj {
	theme := GoApp.Theme()
	lightFg := theme.ColorFg.NRGBA()
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	state := &widget_gio.List{}
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{4,4,4,4},
		Visible: true,
		//target: nil,
	}
	scrollbar := goScrollbar{
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
	hListView := &GoListViewObj {
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		//ScrollbarStyle: Scrollbar(theme, state.Scrollbar),
		AnchorStrategy: Occupy,
		state: 	state,
		columns: 1,
	}
	parent.AddControl(hListView)
	return hListView
}



func (ob *GoListViewObj) ObjectType() (string) {
	return "GoListViewObj"
}

func (ob *GoListViewObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoListViewObj) SetLayoutMode(mode GoLayoutDirection) {
	ob.state.Axis = layout_gio.Axis(uint8(mode))
}

/*func (ob *GoListBoxObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.SetSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoListViewObj) Draw(gtx layout_gio.Context) layout_gio.Dimensions {
	dims := layout_gio.Dimensions{Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx, len(ob.Controls), func(gtx C, i int) D {
						return ob.Controls[i].Draw(gtx)
					})
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}

// layout the list and its scrollbar.
func (ob *GoListViewObj) Layout(gtx layout_gio.Context, length int, w layout_gio.ListElement) layout_gio.Dimensions {
	originalConstraints := gtx.Constraints

	// Determine how much space the scrollbar occupies.
	barWidth := gtx.Dp(ob.Scrollbar.Width())

	if ob.AnchorStrategy == Occupy {

		// Reserve space for the scrollbar using the gtx constraints.
		max := ob.state.Axis.Convert(gtx.Constraints.Max)
		min := ob.state.Axis.Convert(gtx.Constraints.Min)
		max.Y -= barWidth
		if max.Y < 0 {
			max.Y = 0
		}
		min.Y -= barWidth
		if min.Y < 0 {
			min.Y = 0
		}
		gtx.Constraints.Max = ob.state.Axis.Convert(max)
		gtx.Constraints.Min = ob.state.Axis.Convert(min)
	}

	listDims := ob.state.List.Layout(gtx, length, w)
	gtx.Constraints = originalConstraints

	// Draw the scrollbar.
	anchoring := layout_gio.E
	if ob.state.Axis == layout_gio.Horizontal {
		anchoring = layout_gio.S
	}
	majorAxisSize := ob.state.Axis.Convert(listDims.Size).X
	start, end := fromListPosition(ob.state.Position, length, majorAxisSize)
	// layout.Direction respects the minimum, so ensure that the
	// scrollbar will be drawn on the correct edge even if the provided
	// layout.Context had a zero minimum constraint.
	gtx.Constraints.Min = listDims.Size
	if ob.AnchorStrategy == Occupy {
		min := ob.state.Axis.Convert(gtx.Constraints.Min)
		min.Y += barWidth
		gtx.Constraints.Min = ob.state.Axis.Convert(min)
	}
	anchoring.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		return ob.Scrollbar.Layout(gtx, ob.state.Axis, start, end)
	})

	if delta := ob.state.ScrollDistance(); delta != 0 {
		// Handle any changes to the list position as a result of user interaction
		// with the scrollbar.
		ob.state.List.Position.Offset += int(math.Round(float64(float32(ob.state.Position.Length) * delta)))

		// Ensure that the list pays attention to the Offset field when the scrollbar drag
		// is started while the bar is at the end of the list. Without this, the scrollbar
		// cannot be dragged away from the end.
		ob.state.List.Position.BeforeEnd = true
	}

	if ob.AnchorStrategy == Occupy {
		// Increase the width to account for the space occupied by the scrollbar.
		cross := ob.state.Axis.Convert(listDims.Size)
		cross.Y += barWidth
		listDims.Size = ob.state.Axis.Convert(cross)
	}
	//log.Println("listDims :", listDims)
	if ob.MinWidth > listDims.Size.X {
		listDims.Size.X = ob.MinWidth
	}
	return listDims
}

func (ob *GoListViewObj) insertItem(item *GoListViewItemObj) {

}
