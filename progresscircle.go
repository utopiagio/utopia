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
)

type GoProgressCircleObj struct {
	goObject
	goWidget
	color GoColor 	//color.NRGBA
	progress int
	totalSteps int
}

func GoProgressCircle(parent GoObject, totalSteps int) *GoProgressCircleObj {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
		//target: nil,
	}
	hProgress :=  &GoProgressCircleObj{
		goObject: object,
		goWidget: widget,
		color:    theme.GoPalette.ContrastBg,
		totalSteps: totalSteps,
		progress: 0,
	}
	parent.addControl(hProgress)
	return hProgress
}

func (ob *GoProgressCircleObj) SetProgress(progress int) {
	if progress > ob.totalSteps {
		progress = 0
	}
	ob.progress = progress
}

func (ob *GoProgressCircleObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}

func (ob *GoProgressCircleObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.visible {
		dims = ob.goMargin.layout(gtx, func(gtx C) D {
			return ob.goBorder.layout(gtx, func(gtx C) D {
				return ob.goPadding.layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
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