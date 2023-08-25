/* canvaspath.go */

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

func GoCanvasPath(parent *GoCanvasObj) (hCanvasPath *GoCanvasPathObj) {
	item := goCanvasItem{
		active: false,
		animated: false,
		canvas: parent,
		enabled: false,
		selected: false,
		typeId: CanvasPath,
		visible: true,
		xVelocity: 0,
		yVelocity: 0,
		z: 0,
	}
	hCanvasPath = &GoCanvasPathObj{
		item,
		false,			// closed
		[]f32.Point{},	// controlPoints
		Color_White,	// fillColor
		Color_Black,	// lineColor
		2,				// lineWidth
		0,				// x
		0,				// xL
		0,				// xR
		0,				// y
		0,				// yL
		0,				// yR
	}	
	parent.AddItem(hCanvasPath)

	return
}

type GoCanvasPathObj struct {
	goCanvasItem
	closed bool
	controlPoints []f32.Point
	fillColor GoColor
	lineColor GoColor
	lineWidth float32
	x float32
	xL float32
	xR float32
	y float32
	yL float32
	yR float32
}

func (ob *GoCanvasPathObj) AddPoint(x float32, y float32)  {
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

func (ob *GoCanvasPathObj) Advance()  {
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

func (ob *GoCanvasPathObj) BoundingRect() (x float32, y float32, width float32, height float32) {
	return ob.xL, ob.yL, ob.xR, ob.yR
}

func (ob *GoCanvasPathObj) Draw(ops *op_gio.Ops) {
	var path clip_gio.Path
	
	path.Begin(ops)
	path.MoveTo(ob.controlPoints[0])
	for idx := 1; idx < len(ob.controlPoints); idx++ {
		path.LineTo(ob.controlPoints[idx])
	}
	if ob.closed {
		path.Close()
	}
	outline := path.End()
	if ob.closed {
		paint_gio.FillShape(ops, ob.fillColor.NRGBA(),
		clip_gio.Outline{
			Path:  outline,
		}.Op())
	}
	paint_gio.FillShape(ops, ob.lineColor.NRGBA(),
		clip_gio.Stroke{
			Path:  outline,
			Width: ob.lineWidth,
		}.Op())
}

func (ob *GoCanvasPathObj) Hide()  {
	ob.visible = false
}

func (ob *GoCanvasPathObj) Move(x float32, y float32) {
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

func (ob *GoCanvasPathObj) Centre(x float32, y float32) {
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

func (ob *GoCanvasPathObj) SetFillColor(color GoColor) {
	ob.fillColor = color
}

func (ob *GoCanvasPathObj) SetLineColor(color GoColor) {
	ob.lineColor = color
}

func (ob *GoCanvasPathObj) SetLineWidth(width float32) {
	ob.lineWidth = width
}

func (ob *GoCanvasPathObj) SetPos(x float32, y float32) {
	ob.xR = x + ob.xR - ob.xL
	ob.yR = y + ob.yR - ob.yL
	ob.xL = x 
	ob.yL = y 
}

func (ob *GoCanvasPathObj) updateBoundingRect(x float32, y float32)  {
	if len(ob.controlPoints) == 0 {
		ob.xL = x
		ob.xR = x
		ob.x = x
		ob.yL = y
		ob.yR = y
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