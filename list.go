// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/list.go */

package utopia

import (
	"log"
	"image"
	"image/color"
	"math"

	pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"

	"github.com/utopiagio/utopia/metrics"
	//widget_int "github.com/utopiagio/utopia/internal/widget"
)

// fromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list
// in order to do this. The returned values will be in the range [0,1], and
// start will be less than or equal to end.
func fromListPosition(lp layout_gio.Position, elements int, majorAxisSize int) (start, end float32) {
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
}

// rangeIsScrollable returns whether the viewport described by start and end
// is smaller than the underlying content (such that it can be scrolled).
// start and end are expected to each be in the range [0,1], and start
// must be less than or equal to end.
func rangeIsScrollable(start, end float32) bool {
	return end-start < 1
}

// ScrollTrackStyle configures the presentation of a track for a scroll area.
type ScrollTrackStyle struct {
	// MajorPadding and MinorPadding along the major and minor axis of the
	// scrollbar's track. This is used to keep the scrollbar from touching
	// the edges of the content area.
	MajorPadding, MinorPadding unit_gio.Dp
	// Color of the track background.
	Color color.NRGBA
}

// ScrollIndicatorStyle configures the presentation of a scroll indicator.
type ScrollIndicatorStyle struct {
	// MajorMinLen is the smallest that the scroll indicator is allowed to
	// be along the major axis.
	MajorMinLen unit_gio.Dp
	// MinorWidth is the width of the scroll indicator across the minor axis.
	MinorWidth unit_gio.Dp
	// Color and HoverColor are the normal and hovered colors of the scroll
	// indicator.
	Color, HoverColor color.NRGBA
	// CornerRadius is the corner radius of the rectangular indicator. 0
	// will produce square corners. 0.5*MinorWidth will produce perfectly
	// round corners.
	CornerRadius unit_gio.Dp
}

// ScrollbarStyle configures the presentation of a scrollbar.
type goScrollbar struct {
	Scrollbar *widget_gio.Scrollbar
	Track     ScrollTrackStyle
	Indicator ScrollIndicatorStyle
}

// Scrollbar configures the presentation of a scrollbar using the provided
// theme and state.
/*func goScrollbar(theme *Theme, state *widget_gio.Scrollbar) *GoScrollbarObj {
	lightFg := theme.ColorFg
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200

	return GoScrollbarObj{
		Scrollbar: state,
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
}*/

// Width returns the minor axis width of the scrollbar in its current
// configuration (taking padding for the scroll track into account).
func (s *goScrollbar) Width() unit_gio.Dp {
	return s.Indicator.MinorWidth + s.Track.MinorPadding + s.Track.MinorPadding
}

// Layout the scrollbar.
func (s *goScrollbar) Layout(gtx layout_gio.Context, axis layout_gio.Axis, viewportStart, viewportEnd float32) layout_gio.Dimensions {
	if !rangeIsScrollable(viewportStart, viewportEnd) {
		return layout_gio.Dimensions{}
	}

	// Set minimum constraints in an axis-independent way, then convert to
	// the correct representation for the current axis.
	convert := axis.Convert
	maxMajorAxis := convert(gtx.Constraints.Max).X
	gtx.Constraints.Min.X = maxMajorAxis
	gtx.Constraints.Min.Y = gtx.Dp(s.Width())
	gtx.Constraints.Min = convert(gtx.Constraints.Min)
	gtx.Constraints.Max = gtx.Constraints.Min

	s.Scrollbar.Update(gtx, axis, viewportStart, viewportEnd)

	// Darken indicator if hovered.
	if s.Scrollbar.IndicatorHovered() {
		s.Indicator.Color = s.Indicator.HoverColor
	}

	return s.layout(gtx, axis, viewportStart, viewportEnd)
}

// layout the scroll track and indicator.
func (s *goScrollbar) layout(gtx layout_gio.Context, axis layout_gio.Axis, viewportStart, viewportEnd float32) layout_gio.Dimensions {
	inset := layout_gio.Inset{
		Top:    s.Track.MajorPadding,
		Bottom: s.Track.MajorPadding,
		Left:   s.Track.MinorPadding,
		Right:  s.Track.MinorPadding,
	}
	if axis == layout_gio.Horizontal {
		inset.Top, inset.Bottom, inset.Left, inset.Right = inset.Left, inset.Right, inset.Top, inset.Bottom
	}
	// Capture the outer constraints because layout_gio.Stack will reset
	// the minimum to zero.
	outerConstraints := gtx.Constraints

	return layout_gio.Stack{}.Layout(gtx,
		layout_gio.Expanded(func(gtx layout_gio.Context) layout_gio.Dimensions {
			// Lay out the draggable track underneath the scroll indicator.
			area := image.Rectangle{
				Max: gtx.Constraints.Min,
			}
			pointerArea := clip_gio.Rect(area)
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddDrag(gtx.Ops)

			// Stack a normal clickable area on top of the draggable area
			// to capture non-dragging clicks.
			defer pointer_gio.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddTrack(gtx.Ops)

			paint_gio.FillShape(gtx.Ops, s.Track.Color, clip_gio.Rect(area).Op())
			return layout_gio.Dimensions{}
		}),
		layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
			gtx.Constraints = outerConstraints
			return inset.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				// Use axis-independent constraints.
				gtx.Constraints.Min = axis.Convert(gtx.Constraints.Min)
				gtx.Constraints.Max = axis.Convert(gtx.Constraints.Max)

				// Compute the pixel size and position of the scroll indicator within
				// the track.
				trackLen := gtx.Constraints.Min.X
				viewStart := int(math.Round(float64(viewportStart) * float64(trackLen)))
				viewEnd := int(math.Round(float64(viewportEnd) * float64(trackLen)))
				indicatorLen := s.max(viewEnd-viewStart, gtx.Dp(s.Indicator.MajorMinLen))
				if viewStart+indicatorLen > trackLen {
					viewStart = trackLen - indicatorLen
				}
				indicatorDims := axis.Convert(image.Point{
					X: indicatorLen,
					Y: gtx.Dp(s.Indicator.MinorWidth),
				})
				radius := gtx.Dp(s.Indicator.CornerRadius)

				// Lay out the indicator.
				offset := axis.Convert(image.Pt(viewStart, 0))
				defer op_gio.Offset(offset).Push(gtx.Ops).Pop()
				paint_gio.FillShape(gtx.Ops, s.Indicator.Color, clip_gio.RRect{
					Rect: image.Rectangle{
						Max: indicatorDims,
					},
					SW: radius,
					NW: radius,
					NE: radius,
					SE: radius,
				}.Op(gtx.Ops))

				// Add the indicator pointer hit area.
				area := clip_gio.Rect(image.Rectangle{Max: indicatorDims})
				defer pointer_gio.PassOp{}.Push(gtx.Ops).Pop()
				defer area.Push(gtx.Ops).Pop()
				s.Scrollbar.AddIndicator(gtx.Ops)

				return layout_gio.Dimensions{Size: axis.Convert(gtx.Constraints.Min)}
			})
		}),
	)
}

func (s *goScrollbar) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *goScrollbar) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AnchorStrategy defines a means of attaching a scrollbar to content.
type AnchorStrategy uint8

const (
	// Occupy reserves space for the scrollbar, making the underlying
	// content region smaller on one axis.
	Occupy AnchorStrategy = iota
	// Overlay causes the scrollbar to float atop the content without
	// occupying any space. Content in the underlying area can be occluded
	// by the scrollbar.
	Overlay
)

// ListStyle configures the presentation of a layout_gio.List with a scrollbar.
type GoListBoxObj struct {
	GioObject
	GioWidget
	Scrollbar goScrollbar
	AnchorStrategy
	state *widget_gio.List
}

// List constructs a ListStyle using the provided theme and state.
func GoListBox(parent GoObject) *GoListBoxObj {
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
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 300, 300, 16777215, 16777215, 300, 300},
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
	hListBox := &GoListBoxObj {
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		//ScrollbarStyle: Scrollbar(theme, state.Scrollbar),
		AnchorStrategy: Occupy,
		state: 	state,
	}
	parent.AddControl(hListBox)
	return hListBox
}

func (ob *GoListBoxObj) ObjectType() (string) {
	return "GoListBoxObj"
}

func (ob *GoListBoxObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoListBoxObj) SetLayoutMode(mode GoLayoutDirection) {
	ob.state.Axis = layout_gio.Axis(uint8(mode))
}

/*func (ob *GoListBoxObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.SetSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoListBoxObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	log.Println("GoListBoxObj::Draw()")
	cs := gtx.Constraints
	//clipper := gtx.Constraints
	log.Println("gtx.Constraints Min = (", cs.Min.X, cs.Min.Y, ") Max = (", cs.Max.X, cs.Max.Y, ")")
	
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		log.Println("FixedWidth............")
		//log.Println("object Width = (", width, " )")
		cs.Min.X = min(cs.Max.X, width)
		log.Println("cs.Min.X = (", cs.Min.X, " )")
		cs.Max.X = min(cs.Max.X, width)
		log.Println("cs.Max.X = (", cs.Max.X, " )")
	/*case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = min(cs.Min.X, minWidth)
		cs.Max.X = min(cs.Max.X, maxWidth)*/
	case PreferredWidth:		// SizeHint is Preferred
		log.Println("PreferredWidth............")
		log.Println("object MinWidth = (", minWidth, " )")
		log.Println("object MaxWidth = (", maxWidth, " )")
		cs.Min.X = min(cs.Max.X, width)
		cs.Max.X = min(cs.Max.X, width)
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
		log.Println("PreferredWidth............")
		log.Println("object MinHeight = (", minHeight, " )")
		log.Println("object MaxHeight = (", maxHeight, " )")
		cs.Min.Y = min(cs.Max.Y, height)
		cs.Max.Y = min(cs.Max.Y, height)
	/*case MaximumHeight:			// SizeHint is Maximum
		cs.Min.Y = min(cs.Min.Y, ob.MinHeight) 	// No change to gtx.Constraints.Y
		cs.Max.Y = min(cs.Max.Y, ob.MaxHeight)*/
	case ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}

	gtx.Constraints = cs
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Min,}
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
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

// layout the list and its scrollbar.
func (ob *GoListBoxObj) Layout(gtx layout_gio.Context, length int, w layout_gio.ListElement) layout_gio.Dimensions {
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


