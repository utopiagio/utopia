// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvasellipse.go */

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

func GoCanvasEllipse(parent *GoCanvasObj) (hCanvasEllipse *GoCanvasEllipseObj) {
	item := goCanvasItem{
		active: false,
		animated: false,
		canvas: parent,
		enabled: false,
		selected: false,
		typeId: CanvasEllipse,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		z: 0,
	}
	hCanvasEllipse = &GoCanvasEllipseObj{
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
	parent.AddItem(hCanvasEllipse)

	return
}

type GoCanvasEllipseObj struct {
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

func (ob *GoCanvasEllipseObj) Advance()  {
	if ob.animated == true {
		ob.x += float32(ob.xVelocity)
		ob.y += float32(ob.yVelocity)
		ob.xL += float32(ob.xVelocity)
		ob.yL += float32(ob.yVelocity)
		ob.xR += float32(ob.xVelocity)
		ob.yR += float32(ob.yVelocity)
	}
}

func (ob *GoCanvasEllipseObj) BoundingRect() (x float32, y float32, width float32, height float32) {
	return ob.xL, ob.yL, ob.xR, ob.yR
}

func (ob *GoCanvasEllipseObj) Draw(ops *op_gio.Ops) {
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

func (ob *GoCanvasEllipseObj) Height() (height float32) {
	return ob.yR - ob.yL
}

func (ob *GoCanvasEllipseObj) Hide()  {
	ob.visible = false
}

func (ob *GoCanvasEllipseObj) Move(x float32, y float32) {
	ob.x += x
	ob.y += y
	ob.xL += x
	ob.yL += y
	ob.xR += x
	ob.yR += y
}

func (ob *GoCanvasEllipseObj) Centre(x float32, y float32) {
	ob.xL += x - ob.x
	ob.yL += y - ob.y
	ob.xR += x - ob.x
	ob.yR += y - ob.y
	ob.x = x
	ob.y = y
}

func (ob *GoCanvasEllipseObj) SetFillColor(color GoColor) {
	ob.fillColor = color
}

func (ob *GoCanvasEllipseObj) SetHeight(height float32) {
	ob.xL = ob.x - height / 2
	ob.xR = ob.x + height / 2
}

func (ob *GoCanvasEllipseObj) SetLineColor(color GoColor) {
	ob.lineColor = color
}

func (ob *GoCanvasEllipseObj) SetLineWidth(width float32) {
	ob.lineWidth = width
}

func (ob *GoCanvasEllipseObj) SetPos(x float32, y float32) {
	ob.xR = x + ob.xR - ob.xL
	ob.yR = y + ob.yR - ob.yL
	ob.xL = x 
	ob.yL = y 
}

func (ob *GoCanvasEllipseObj) SetSize(width float32, height float32) {
	ob.xL = ob.x - width / 2
	ob.yL = ob.y - height / 2
	ob.xR = ob.x + width / 2
	ob.yR = ob.y + height / 2
}

func (ob *GoCanvasEllipseObj) SetWidth(width float32) {
	ob.xL = ob.x - width / 2
	ob.xR = ob.x + width / 2
}

func (ob *GoCanvasEllipseObj) Width() (width float32) {
	return ob.xR - ob.xL
}