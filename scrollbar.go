// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/list.go */

package utopia

import (
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
)



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
type GoScrollbar struct {
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
func (ob *GoScrollbar) Width() unit_gio.Dp {
	return ob.Indicator.MinorWidth + ob.Track.MinorPadding + ob.Track.MinorPadding
}

// Layout the scrollbar.
func (ob *GoScrollbar) Layout(gtx layout_gio.Context, axis layout_gio.Axis, viewportStart, viewportEnd float32) layout_gio.Dimensions {
	if !ob.rangeIsScrollable(viewportStart, viewportEnd) {
		return layout_gio.Dimensions{}
	}

	// Set minimum constraints in an axis-independent way, then convert to
	// the correct representation for the current axis.
	convert := axis.Convert
	maxMajorAxis := convert(gtx.Constraints.Max).X
	gtx.Constraints.Min.X = maxMajorAxis
	gtx.Constraints.Min.Y = gtx.Dp(ob.Width())
	gtx.Constraints.Min = convert(gtx.Constraints.Min)
	gtx.Constraints.Max = gtx.Constraints.Min

	ob.Scrollbar.Update(gtx, axis, viewportStart, viewportEnd)

	// Darken indicator if hovered.
	if ob.Scrollbar.IndicatorHovered() {
		ob.Indicator.Color = ob.Indicator.HoverColor
	}

	return ob.layout(gtx, axis, viewportStart, viewportEnd)
}

// layout the scroll track and indicator.
func (ob *GoScrollbar) layout(gtx layout_gio.Context, axis layout_gio.Axis, viewportStart, viewportEnd float32) layout_gio.Dimensions {
	inset := layout_gio.Inset{
		Top:    ob.Track.MajorPadding,
		Bottom: ob.Track.MajorPadding,
		Left:   ob.Track.MinorPadding,
		Right:  ob.Track.MinorPadding,
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
			ob.Scrollbar.AddDrag(gtx.Ops)

			// Stack a normal clickable area on top of the draggable area
			// to capture non-dragging clicks.
			defer pointer_gio.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			ob.Scrollbar.AddTrack(gtx.Ops)

			paint_gio.FillShape(gtx.Ops, ob.Track.Color, clip_gio.Rect(area).Op())
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
				indicatorLen := ob.max(viewEnd-viewStart, gtx.Dp(ob.Indicator.MajorMinLen))
				if viewStart+indicatorLen > trackLen {
					viewStart = trackLen - indicatorLen
				}
				indicatorDims := axis.Convert(image.Point{
					X: indicatorLen,
					Y: gtx.Dp(ob.Indicator.MinorWidth),
				})
				radius := gtx.Dp(ob.Indicator.CornerRadius)

				// Lay out the indicator.
				offset := axis.Convert(image.Pt(viewStart, 0))
				defer op_gio.Offset(offset).Push(gtx.Ops).Pop()
				paint_gio.FillShape(gtx.Ops, ob.Indicator.Color, clip_gio.RRect{
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
				ob.Scrollbar.AddIndicator(gtx.Ops)

				return layout_gio.Dimensions{Size: axis.Convert(gtx.Constraints.Min)}
			})
		}),
	)
}

// rangeIsScrollable returns whether the viewport described by start and end
// is smaller than the underlying content (such that it can be scrolled).
// start and end are expected to each be in the range [0,1], and start
// must be less than or equal to end.
func (ob *GoScrollbar) rangeIsScrollable(start, end float32) bool {
	return end - start < 1
}

func (ob *GoScrollbar) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (ob *GoScrollbar) min(a, b int) int {
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