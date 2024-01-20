// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/loader.go */

package utopia

import (
	"image"
	"log"
	"math"
	"time"

	"github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"

	"github.com/utopiagio/utopia/metrics"
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
	log.Println("GoLoaderObj::Draw()")
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
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
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