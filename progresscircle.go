// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/progresscircle.go */

package utopia

import (
	"image"
	//"image/color"
	//"log"
	"math"
	"github.com/utopiagio/gio/op/clip"
	"github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"

	"github.com/utopiagio/utopia/metrics"
)

type GoProgressCircleObj struct {
	GioObject
	GioWidget
	color GoColor 	//color.NRGBA
	progress int
	totalSteps int
}

func GoProgressCircle(parent GoObject, totalSteps int) *GoProgressCircleObj {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
		//target: nil,
	}
	hProgress :=  &GoProgressCircleObj{
		GioObject: object,
		GioWidget: widget,
		color:    theme.GoPalette.ContrastBg,
		totalSteps: totalSteps,
		progress: 0,
	}
	parent.AddControl(hProgress)
	return hProgress
}

func (ob *GoProgressCircleObj) ObjectType() (string) {
	return "GoProgressCircleObj"
}

func (ob *GoProgressCircleObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoProgressCircleObj) SetProgress(progress int) {
	if progress > ob.totalSteps {
		progress = 0
	}
	ob.progress = progress
}

/*func (ob *GoProgressCircleObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.SetSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoProgressCircleObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
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

func (ob *GoProgressCircleObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
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
	
	defer ob.clipLoader(gtx.Ops, - math.Pi / 2, - math.Pi / 2 + math.Pi * 2 * (float32(ob.progress) / float32(ob.totalSteps)), float32(radius)).Push(gtx.Ops).Pop()
	paint_gio.ColorOp{
		Color: ob.color.NRGBA(),
	}.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)
	return layout_gio.Dimensions{
		Size: sz,
	}
}

func (ob *GoProgressCircleObj) start() {
	/*go func() {
		for {
			time.Sleep(time.Second)
			ob.progress += 0.1
			if ob.progress > 1.0 {
				ob.progress = 0
			}
			progressIncrementer <- 0.1
		}
	}()*/
}

func (ob *GoProgressCircleObj) stop() {
	
}

func (ob *GoProgressCircleObj) clipLoader(ops *op_gio.Ops, startAngle, endAngle, radius float32) clip.Op {
	const thickness = .25

	var (
		width = radius * thickness
		delta = endAngle - startAngle

		vy, vx = math.Sincos(float64(startAngle))

		inner  = radius * (1. - thickness*.5)
		pen    = f32.Pt(float32(vx), float32(vy)).Mul(inner)
		center = f32.Pt(0, 0).Sub(pen)

		p clip.Path
	)

	p.Begin(ops)
	p.Move(pen)
	p.Arc(center, center, delta)
	return clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op()
}