// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvasitem.go */

package utopia

import (
	//"log"
	//"image"
	//"image/color"
	//"math"
	//"github.com/utopiagio/gio/f32"
	//"github.com/utopiagio/gio/font/gofont"
	//layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	//text_gio "github.com/utopiagio/gio/text"
	//unit_gio "github.com/utopiagio/gio/unit"
)

const (
	CanvasLine int = iota
	CanvasRectangle
	CanvasCircle
	CanvasEllipse
	CanvasPath
	CanvasSpline
)

type GoCanvasItem interface {
	Advance()
	BoundingRect() (float32, float32, float32, float32)
	Draw(ops *op_gio.Ops)
	Hide()
	Move(x float32, y float32)
	Centre(x float32, y float32)
	Show()
	Type() (int)
}


/*func GoCanvasItem(parent *GoCanvasObj) (hCanvasItem *CanvasItem) {
	hCanvasItem = &GoCanvasItemObj{
		active: false,
		animated: false,
		canvas: parent,
		color: Color_Red,
		enabled: false,
		selected: false,
		typeId: 0,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		x: 0,
		y: 0,
		z: 0,
	}	
	parent.AddItem(hCanvasItem)

	return hCanvasItem
}*/

type goCanvasItem struct {
	active bool
	animated bool
	canvas *GoCanvasObj
	enabled bool
	selected bool
	typeId int
	visible bool
	xVelocity int
	yVelocity int
	z int
}

func (ob *goCanvasItem) Hide()  {
	ob.visible = false
	if ob.canvas != nil {
		ob.canvas.Update()
	}
}

func (ob *goCanvasItem) IsActive() (active bool) {
	return ob.active
}

func (ob *goCanvasItem) IsAnimated() (animated bool) {
	return ob.animated
}

func (ob *goCanvasItem) IsEnabled() (enabled bool) {
	return ob.enabled
}

func (ob *goCanvasItem) IsSelected() (selected bool) {
	return ob.selected
}

func (ob *goCanvasItem) IsVisible() (visible bool) {
	return ob.visible
}

func (ob *goCanvasItem) SetActive(active bool) {
	ob.active = active
}

func (ob *goCanvasItem) SetAnimated(animated bool) {
	ob.animated = animated
}

func (ob *goCanvasItem) SetEnabled(enabled bool) {
	ob.enabled = enabled
}

func (ob *goCanvasItem) SetSelected(selected bool) {
	ob.selected = selected
}

func (ob *goCanvasItem) SetVelocity(x int, y int) {
	ob.xVelocity = x
	ob.yVelocity = y
}

func (ob *goCanvasItem) SetVisible(visible bool) {
	ob.visible = visible
}

func (ob *goCanvasItem) SetXVelocity(x int) {
	ob.xVelocity = x
}

func (ob *goCanvasItem) SetYVelocity(y int) {
	ob.yVelocity = y
}

func (ob *goCanvasItem) SetZ(z int) {
	ob.z = z
}

func (ob *goCanvasItem) Show()  {
	ob.visible = true
	if ob.canvas != nil {
		ob.canvas.Update()
	}
}

func (ob *goCanvasItem) Type() (typeId int) {
	return ob.typeId
}

func (ob *goCanvasItem) XVelocity() (x int) {
	return ob.xVelocity
}

func (ob *goCanvasItem) YVelocity() (y int) {
	return ob.yVelocity
}

func (ob *goCanvasItem) Z() (z int) {
	return ob.z
}

