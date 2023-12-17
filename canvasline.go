// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvasline.go */

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

func GoCanvasLine(parent *GoCanvasObj) (hCanvasLine *GoCanvasLineObj) {
	item := goCanvasItem{
		active: false,
		animated: false,
		canvas: parent,
		enabled: false,
		selected: false,
		typeId: CanvasLine,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		z: 0,
	}
	hCanvasLine = &GoCanvasLineObj{
		item,
		0,				// endX
		0,				// endY
		Color_Red,		// lineColor
		2,				// lineWidth
		0,				// startX
		0,				// startY
		0,				// x
		0,				// y
	}	
	parent.AddItem(hCanvasLine)
	return
}

type GoCanvasLineObj struct {
	goCanvasItem
	endX float32
	endY float32
	lineColor GoColor
	lineWidth float32
	startX float32
	startY float32
	x float32
	y float32
	
}

func (ob *GoCanvasLineObj) Advance()  {
	if ob.animated == true {
		ob.startX += float32(ob.xVelocity)
		ob.startY += float32(ob.yVelocity)
		ob.endX += float32(ob.xVelocity)
		ob.endY += float32(ob.yVelocity)
		ob.x += float32(ob.xVelocity)
		ob.y += float32(ob.yVelocity)
	}
}

func (ob *GoCanvasLineObj) BoundingRect() (x float32, y float32, width float32, height float32) {
	return ob.startX, ob.startY, ob.endX - ob.startX, ob.endY - ob.startY
}

func (ob *GoCanvasLineObj) Draw(ops *op_gio.Ops) {
		var path clip_gio.Path
		path.Begin(ops)
		path.MoveTo(f32.Pt(ob.startX, ob.startY))
		path.LineTo(f32.Pt(ob.endX, ob.endY))
		path.Close()
		paint_gio.FillShape(ops, ob.lineColor.NRGBA(),
			clip_gio.Stroke{
				Path:  path.End(),
				Width: ob.lineWidth,
			}.Op())
}

func (ob *GoCanvasLineObj) Move(x float32, y float32) {
	ob.startX += x
	ob.startY += y
	ob.endX += x
	ob.endY += y
	ob.x += x
	ob.y += y
}

func (ob *GoCanvasLineObj) Centre(x float32, y float32) {
	ob.startX += x - ob.x
	ob.startY += y - ob.y
	ob.endX += x - ob.x
	ob.endY += y - ob.y
	ob.x = x
	ob.y = y
}

func (ob *GoCanvasLineObj) SetPoints(startX, startY, endX, endY float32) {
	ob.startX = startX
	ob.startY = startY
	ob.endX = endX
	ob.endY = endY
}

func (ob *GoCanvasLineObj) SetLineColor(color GoColor) {
	ob.lineColor = color
}

func (ob *GoCanvasLineObj) SetLineWidth(width float32) {
	ob.lineWidth = width
}

