// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvascircle.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"
	//"math"
	//"github.com/utopiagio/gio/f32"
	//"github.com/utopiagio/gio/font/gofont"
	//layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	//text_gio "github.com/utopiagio/gio/text"
	//unit_gio "github.com/utopiagio/gio/unit"
)

func GoCanvasCircle(parent *GoCanvasObj) (hCanvasCircle *GoCanvasCircleObj) {
	item := goCanvasItem{
		active: false,
		animated: false,
		canvas: parent,
		enabled: false,
		selected: false,
		typeId: CanvasCircle,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		z: 0,
	}
	hCanvasCircle = &GoCanvasCircleObj{
		item,
		Color_White,	// fillColor
		Color_Black,	// lineColor
		2,				// lineWidth
		0,				// radius
		0,				// x
		0,				// xL
		0,				// xR
		0,				// y
		0,				// yL
		0,				// yR
	}	
	parent.AddItem(hCanvasCircle)

	return
}

type GoCanvasCircleObj struct {
	goCanvasItem
	fillColor GoColor
	lineColor GoColor
	lineWidth float32
	radius float32
	x float32
	xL float32
	xR float32
	y float32
	yL float32
	yR float32
}

func (ob *GoCanvasCircleObj) Advance()  {
	if ob.animated == true {
		ob.x += float32(ob.xVelocity)
		ob.y += float32(ob.yVelocity)
		ob.xL += float32(ob.xVelocity)
		ob.yL += float32(ob.yVelocity)
		ob.xR += float32(ob.xVelocity)
		ob.yR += float32(ob.yVelocity)
	}
}

func (ob *GoCanvasCircleObj) BoundingRect() (x float32, y float32, width float32, height float32) {
	return ob.xL, ob.yL, ob.xR, ob.yR
}

func (ob *GoCanvasCircleObj) Draw(ops *op_gio.Ops) {
	ellipse := clip_gio.Ellipse{
		Min: image.Pt(int(ob.xL), int(ob.yL)),
		Max: image.Pt(int(ob.xR), int(ob.yR)),
	}
	paint_gio.FillShape(ops, ob.fillColor.NRGBA(),
		clip_gio.Outline{
			Path:  ellipse.Path(ops),
		}.Op())
	paint_gio.FillShape(ops, ob.lineColor.NRGBA(),
		clip_gio.Stroke{
			Path:  ellipse.Path(ops),
			Width: ob.lineWidth,
		}.Op())
}

func (ob *GoCanvasCircleObj) Hide()  {
	ob.visible = false
}

func (ob *GoCanvasCircleObj) Move(x float32, y float32) {
	ob.x += x
	ob.y += y
	ob.xL += x
	ob.yL += y
	ob.xR += x
	ob.yR += y
}

func (ob *GoCanvasCircleObj) Centre(x float32, y float32) {
	ob.x = x
	ob.y = y
	ob.xL = ob.x - ob.radius
	ob.yL = ob.y - ob.radius
	ob.xR = ob.x + ob.radius
	ob.yR = ob.y + ob.radius
}

func (ob *GoCanvasCircleObj) SetFillColor(color GoColor) {
	ob.fillColor = color
}

func (ob *GoCanvasCircleObj) SetLineColor(color GoColor) {
	ob.lineColor = color
}

func (ob *GoCanvasCircleObj) SetLineWidth(width float32) {
	ob.lineWidth = width
}

func (ob *GoCanvasCircleObj) SetPos(x float32, y float32) {
	ob.xR = x + ob.xR - ob.xL
	ob.yR = y + ob.yR - ob.yL
	ob.xL = x 
	ob.yL = y 
}

func (ob *GoCanvasCircleObj) SetRadius(radius float32) {
	ob.radius = radius
	ob.xL = ob.x - radius
	ob.yL = ob.y - radius
	ob.xR = ob.x + radius
	ob.yR = ob.y + radius
}

/*func (ob *GoCanvasCircleObj) updateBoundingRect(x float32, y float32)  {
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
}*/