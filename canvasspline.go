// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvasspline.go */

package utopia

import (
	//"log"
	//"image"
	//"image/color"
	//"math"
	"github.com/utopiagio/gio/f32"
	//"github.com/utopiagio/gio/font/gofont"
	//layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	//text_gio "github.com/utopiagio/gio/text"
	//unit_gio "github.com/utopiagio/gio/unit"
)

func GoCanvasSpline(parent *GoCanvasObj) (hCanvasSpline *GoCanvasSplineObj) {
	hCanvasSpline = &GoCanvasSplineObj{
		active: false,
		animated: false,
		canvas: parent,
		closed: false,
		controlPoints: []f32.Point{},
		enabled: false,
		fillColor: Color_White,
		lineColor: Color_Black,
		lineWidth: 2,
		selected: false,
		typeId: 2,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		x: 0,
		xL:0,
		xR: 0,
		y: 0,
		yL: 0,
		yR: 0,
		z: 0,
	}	
	parent.AddItem(hCanvasSpline)

	return
}

type GoCanvasSplineObj struct {
	active bool
	animated bool
	canvas *GoCanvasObj
	closed bool
	
	controlPoints []f32.Point
	enabled bool
	fillColor GoColor
	lineColor GoColor
	lineWidth float32
	selected bool
	
	typeId int
	visible bool
	
	xVelocity int
	yVelocity int
	x float32
	xL float32
	xR float32
	y float32
	yL float32
	yR float32
	z int
}

func (ob *GoCanvasSplineObj) AddPoint(x float32, y float32)  {
	if !ob.closed {
		if len(ob.controlPoints) != 0 {
			if ob.controlPoints[0] == f32.Pt(x, y) {
				ob.closed = true
			} else {
				ob.controlPoints = append(ob.controlPoints, f32.Pt(x, y))
				ob.updateBoundingRect(x, y)
			}
		} else {
			ob.controlPoints = append(ob.controlPoints, f32.Pt(x, y))
			ob.updateBoundingRect(x, y)
		}
	}
}

func (ob *GoCanvasSplineObj) Advance()  {
	if ob.animated == true {
		for idx := 0; idx < len(ob.controlPoints); idx++ {
			ob.controlPoints[idx].X += float32(ob.xVelocity)
			ob.controlPoints[idx].Y += float32(ob.yVelocity)
		}
		ob.x += float32(ob.xVelocity)
		ob.y += float32(ob.yVelocity)
		ob.xL += float32(ob.xVelocity)
		ob.yL += float32(ob.yVelocity)
		ob.xR += float32(ob.xVelocity)
		ob.yR += float32(ob.yVelocity)
	}
}

func (ob *GoCanvasSplineObj) BoundingRect() (x float32, y float32, width float32, height float32) {
	return ob.xL, ob.yL, ob.xR, ob.yR
}

func (ob *GoCanvasSplineObj) Draw(ops *op_gio.Ops) {
	var path clip_gio.Path
	var fill clip_gio.Path
	fill.Begin(ops)
	fill.MoveTo(ob.controlPoints[0])
	for idx := 1; idx < len(ob.controlPoints); idx++ {
		fill.LineTo(ob.controlPoints[idx])
	}
	if ob.closed {
		fill.Close()
		paint_gio.FillShape(ops, ob.fillColor.NRGBA(),
		clip_gio.Outline{
			Path:  fill.End(),
		}.Op())
	}
	path.Begin(ops)
	path.MoveTo(ob.controlPoints[0])
	for idx := 1; idx < len(ob.controlPoints); idx++ {
		path.LineTo(ob.controlPoints[idx])
	}
	if ob.closed {
		path.Close()
	}
	paint_gio.FillShape(ops, ob.lineColor.NRGBA(),
		clip_gio.Stroke{
			Path:  path.End(),
			Width: ob.lineWidth,
		}.Op())
}

func (ob *GoCanvasSplineObj) Hide()  {
	ob.visible = false
}

func (ob *GoCanvasSplineObj) Move(x float32, y float32) {
	for idx := 0; idx < len(ob.controlPoints); idx++ {
			ob.controlPoints[idx].X += x
			ob.controlPoints[idx].Y += y
	}
	ob.x += x
	ob.y += y
	ob.xL += x
	ob.yL += y
	ob.xR += x
	ob.yR += y
}

func (ob *GoCanvasSplineObj) Centre(x float32, y float32) {
	for idx := 0; idx < len(ob.controlPoints); idx++ {
			ob.controlPoints[idx].X += x
			ob.controlPoints[idx].Y += y
	}
	ob.x = x
	ob.y = y
	ob.xL += x - ob.xL
	ob.yL += y - ob.yL
	ob.xR += x - ob.xR
	ob.yR += y - ob.yR
}



func (ob *GoCanvasSplineObj) Show()  {
	ob.visible = true
}

func (ob *GoCanvasSplineObj) Type() (typeId int) {
	return ob.typeId
}

func (ob *GoCanvasSplineObj) updateBoundingRect(x float32, y float32)  {
	if len(ob.controlPoints) == 0 {
		ob.xL = x
		ob.xR = x
		ob.x = x
		ob.yL = y
		ob.xR = y
		ob.y = y
	} else {
		if x > ob.xR {
			ob.xR = x
		} else if x < ob.xL {
			ob.xL = x
		}
		if y > ob.yR {
			ob.yR = y
		} else if y < ob.yL {
			ob.yL = y
		}
		ob.x = (ob.xL + ob.xR) / 2
		ob.y = (ob.yL + ob.yR) / 2
	}
}