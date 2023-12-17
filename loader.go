// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/loader.go */

package utopia

import (
	"image"
	"math"
	"time"

	"github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
)

type GoLoaderObj struct {
	GioObject
	GioWidget
	color GoColor
}

func GoLoader(parent GoObject) *GoLoaderObj {
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLoader := &GoLoaderObj{
		GioObject: object,
		GioWidget: widget,
		color: theme.ContrastBg,
	}
	parent.AddControl(hLoader)
	return hLoader
}

func (ob *GoLoaderObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}

func (ob *GoLoaderObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	diam := gtx.Constraints.Min.X
	if minY := gtx.Constraints.Min.Y; minY > diam {
		diam = minY
	}
	if diam == 0 {
		diam = gtx.Dp(24)
	}
	sz := gtx.Constraints.Constrain(image.Pt(diam, diam))
	radius := sz.X / 2
	defer op_gio.Offset(image.Pt(radius, radius)).Push(gtx.Ops).Pop()

	dt := float32((time.Duration(gtx.Now.UnixNano()) % (time.Second)).Seconds())
	startAngle := dt * math.Pi * 2
	endAngle := startAngle + math.Pi*1.5

	defer clipLoader(gtx.Ops, startAngle, endAngle, float32(radius)).Push(gtx.Ops).Pop()
	paint_gio.ColorOp{
		Color: ob.color.NRGBA(),
	}.Add(gtx.Ops)
	defer op_gio.Offset(image.Pt(-radius, -radius)).Push(gtx.Ops).Pop()
	paint_gio.PaintOp{}.Add(gtx.Ops)
	op_gio.InvalidateOp{}.Add(gtx.Ops)
	return layout_gio.Dimensions{
		Size: sz,
	}
}

func (ob *GoLoaderObj) ObjectType() (string) {
	return "GoLoaderObj"
}

func (ob *GoLoaderObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func clipLoader(ops *op_gio.Ops, startAngle, endAngle, radius float32) clip_gio.Op {
	const thickness = .25

	var (
		width = radius * thickness
		delta = endAngle - startAngle

		vy, vx = math.Sincos(float64(startAngle))

		inner  = radius * (1. - thickness*.5)
		pen    = f32.Pt(float32(vx), float32(vy)).Mul(inner)
		center = f32.Pt(0, 0).Sub(pen)

		p clip_gio.Path
	)

	p.Begin(ops)
	p.Move(pen)
	p.Arc(center, center, delta)
	return clip_gio.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op()
}