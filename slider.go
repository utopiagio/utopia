// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	//"log"
	"image"
	//"image/color"

	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

// Slider is for selecting a value in a range.
func GoSlider(parent GoObject, min float32, max float32) *GoSliderObj {
	var theme *GoThemeObj = GoApp.Theme()
	var slider *widget_gio.Float = new(widget_gio.Float)
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hSlider := &GoSliderObj{
		GioObject: object,
		GioWidget: widget,
		min:        min,
		max:        max,
		color:      theme.ContrastBg,
		gioSlider:   slider,
		fingerSize: theme.FingerSize,
	}
	parent.AddControl(hSlider)
	return hSlider
}

type GoSliderObj struct {
	GioObject
	GioWidget
	axis    layout_gio.Axis
	min 	float32
	max 	float32
	invert  bool
	color   GoColor
	gioSlider *widget_gio.Float
	fingerSize unit_gio.Dp

	onChange func(float32)
	onDrag 	func(float32)
}

func (ob *GoSliderObj) SetMaxValue(max float32) {
	ob.max = max
}

func (ob *GoSliderObj) SetMinValue(min float32) {
	ob.min = min
}

func (ob *GoSliderObj) Changed() bool {
	return ob.gioSlider.Changed()
}

func (ob *GoSliderObj) Dragging() bool {
	return ob.gioSlider.Dragging()
}

func (ob *GoSliderObj) ObjectType() (string) {
	return "GoSliderObj"
}

func (ob *GoSliderObj) SetOnChange(f func(float32)) {
	ob.onChange = f
}

func (ob *GoSliderObj) SetOnDrag(f func(float32)) {
	ob.onDrag = f
}

/*func (ob *GoSliderObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoSliderObj) Value() float32 {
	return ob.gioSlider.Value
}

func (ob *GoSliderObj) SetValue(value float32) {
	ob.gioSlider.Value = value
}

func (ob *GoSliderObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
	}
	return dims
}

func (ob *GoSliderObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	thumbRadius := gtx.Dp(6)
	trackWidth := gtx.Dp(2)

	axis := ob.axis
	// Keep a minimum length so that the track is always visible.
	minLength := thumbRadius + 3*thumbRadius + thumbRadius
	// Try to expand to finger size, but only if the constraints
	// allow for it.
	touchSizePx := ob.minValue(gtx.Dp(ob.fingerSize), axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := ob.maxValue(axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := ob.maxValue(2*thumbRadius, touchSizePx)
	size := axis.Convert(image.Pt(sizeMain, sizeCross))

	offset := axis.Convert(image.Pt(thumbRadius, 0))
	trans := op_gio.Offset(offset).Push(gtx.Ops)
	gtx.Constraints.Min = axis.Convert(image.Pt(sizeMain-2*thumbRadius, sizeCross))
	ob.gioSlider.Layout(gtx, axis, ob.min, ob.max, ob.invert, thumbRadius)
	gtx.Constraints.Min = gtx.Constraints.Min.Add(axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := thumbRadius + int(ob.gioSlider.Pos())
	trans.Pop()

	color := ob.color.NRGBA()
	if gtx.Queue == nil {
		color = DisabledBlend(color)
	}

	rect := func(minx, miny, maxx, maxy int) image.Rectangle {
		r := image.Rect(minx, miny, maxx, maxy)
		if ob.invert != (axis == layout_gio.Vertical) {
			r.Max.X, r.Min.X = sizeMain-r.Min.X, sizeMain-r.Max.X
		}
		r.Min = axis.Convert(r.Min)
		r.Max = axis.Convert(r.Max)
		return r
	}

	// Draw track before thumb.
	track := rect(
		thumbRadius, sizeCross/2-trackWidth/2,
		thumbPos, sizeCross/2+trackWidth/2,
	)
	paint_gio.FillShape(gtx.Ops, color, clip_gio.Rect(track).Op())

	// Draw track after thumb.
	track = rect(
		thumbPos, axis.Convert(track.Min).Y,
		sizeMain-thumbRadius, axis.Convert(track.Max).Y,
	)
	paint_gio.FillShape(gtx.Ops, MulAlpha(color, 96), clip_gio.Rect(track).Op())

	// Draw thumb.
	pt := image.Pt(thumbPos, sizeCross/2)
	thumb := rect(
		pt.X-thumbRadius, pt.Y-thumbRadius,
		pt.X+thumbRadius, pt.Y+thumbRadius,
	)
	paint_gio.FillShape(gtx.Ops, color, clip_gio.Ellipse(thumb).Op(gtx.Ops))

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