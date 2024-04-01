// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/slider.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"
	"github.com/utopiagio/utopia/internal/f32color"
	"github.com/utopiagio/utopia/metrics"
	widget_int "github.com/utopiagio/utopia/internal/widget"

	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	
)

// Slider is for selecting a value in a range.
func GoSlider(parent GoObject, min int, max int) *GoSliderObj {
	var theme *GoThemeObj = GoApp.Theme()
	var gioFloat *widget_int.GioFloat = new(widget_int.GioFloat)
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, FixedHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{4,4,4,4},
		GoSize: GoSize{0, 0, 200, 20, 16777215, 16777215, 200, 20},
		Visible: true,
		tag: tagCounter,
	}
	hSlider := &GoSliderObj{
		GioObject: object,
		GioWidget: widget,
		Color:      theme.ContrastBg,
		GioFloat:   gioFloat,
		FingerSize: theme.FingerSize,
		ThumbRadius: theme.ThumbRadius,
		TrackWidth: theme.TrackWidth,
		min:        float32(min),
		max:        float32(max),
	}
	parent.AddControl(hSlider)
	return hSlider
}

type GoSliderObj struct {
	GioObject
	GioWidget
	Axis    layout_gio.Axis
	//invert  bool
	Color   GoColor
	GioFloat *widget_int.GioFloat
	FingerSize unit_gio.Dp
	ThumbRadius unit_gio.Dp
	TrackWidth unit_gio.Dp
	min 	float32
	max 	float32
	changed bool
	onChange func(int)
	onDrag 	func(float32)
}

func (ob *GoSliderObj) SetMaxValue(max int) {
	ob.max = float32(max)
}

func (ob *GoSliderObj) SetMinValue(min int) {
	ob.min = float32(min)
}

func (ob *GoSliderObj) Changed() bool {
	return ob.changed
}

func (ob *GoSliderObj) Dragging() bool {
	return ob.GioFloat.Dragging()
}

func (ob *GoSliderObj) ObjectType() (string) {
	return "GoSliderObj"
}

func (ob *GoSliderObj) SetOnChange(f func(int)) {
	ob.onChange = f
}

func (ob *GoSliderObj) SetOnDrag(f func(float32)) {
	ob.onDrag = f
}

func (ob *GoSliderObj) Value() int {
	 return int((ob.GioFloat.Value * (ob.max - ob.min)) + ob.min)
}

func (ob *GoSliderObj) SetValue(value int) {
	ob.GioFloat.Value = (float32(value) - ob.min) / (ob.max - ob.min)
}

func (ob *GoSliderObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoSliderObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
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

func (ob *GoSliderObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.changed = ob.GioFloat.Update(gtx)
	if ob.changed {
		if ob.onChange != nil {
			ob.onChange(int((ob.GioFloat.Value * (ob.max - ob.min)) + ob.min))
		}
	}

	//const thumbRadius unit_gio.Dp = 6

	tr := gtx.Dp(ob.ThumbRadius)
	trackWidth := gtx.Dp(ob.TrackWidth)

	axis := ob.Axis

	/*if ob.SizePolicy().Horiz == FixedWidth {
		if axis == layout_gio.Horizontal{
			gtx.Constraints.Min.X = metrics.DpToPx(GoDpr, ob.Width)
		} else {
			gtx.Constraints.Min.X = metrics.DpToPx(GoDpr, ob.Height)
		}
	}*/

	// Keep a minimum length so that the track is always visible.
	minLength := tr + 3*tr + tr
	// Try to expand to finger size, but only if the constraints
	// allow for it.
	touchSizePx := ob.minValue(gtx.Dp(ob.FingerSize), axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := ob.maxValue(axis.Convert(gtx.Constraints.Max).X, minLength)
	sizeCross := ob.maxValue(2*tr, touchSizePx)
	size := axis.Convert(image.Pt(sizeMain, sizeCross + 1))

	o := axis.Convert(image.Pt(tr, 0))
	trans := op_gio.Offset(o).Push(gtx.Ops)
	gtx.Constraints.Min = axis.Convert(image.Pt(sizeMain-2*tr, sizeCross))
	dims := ob.GioFloat.Layout(gtx, axis, ob.ThumbRadius)
	gtx.Constraints.Min = gtx.Constraints.Min.Add(axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := tr + int(ob.GioFloat.Value*float32(axis.Convert(dims.Size).X))
	trans.Pop()

	color := ob.Color.NRGBA()
	if !gtx.Enabled() {
		color = f32color.Disabled(color)
	}

	rect := func(minx, miny, maxx, maxy int) image.Rectangle {
		r := image.Rect(minx, miny, maxx, maxy)
		if axis == layout_gio.Vertical {
			r.Max.X, r.Min.X = sizeMain-r.Min.X, sizeMain-r.Max.X
		}
		r.Min = axis.Convert(r.Min)
		r.Max = axis.Convert(r.Max)
		return r
	}

	// Draw track before thumb.
	track := rect(
		tr, sizeCross/2-trackWidth/2,
		thumbPos, sizeCross/2+trackWidth/2,
	)
	paint_gio.FillShape(gtx.Ops, color, clip_gio.Rect(track).Op())

	// Draw track after thumb.
	track = rect(
		thumbPos, axis.Convert(track.Min).Y,
		sizeMain-tr, axis.Convert(track.Max).Y,
	)
	paint_gio.FillShape(gtx.Ops, f32color.MulAlpha(color, 96), clip_gio.Rect(track).Op())

	// Draw thumb.
	pt := image.Pt(thumbPos, sizeCross/2)
	thumb := rect(
		pt.X-tr, pt.Y-tr,
		pt.X+tr, pt.Y+tr,
	)
	paint_gio.FillShape(gtx.Ops, color, clip_gio.Ellipse(thumb).Op(gtx.Ops))
	if ob.SizePolicy().Horiz == ExpandingWidth {
		size = image.Pt(size.X, gtx.Constraints.Max.X)
	}
	if ob.SizePolicy().Vert == ExpandingHeight {
		size = image.Pt(size.X, gtx.Constraints.Max.Y)
	}
	return layout_gio.Dimensions{Size: size}
}

func (ob *GoSliderObj) maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (ob *GoSliderObj) minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}