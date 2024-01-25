// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/canvas.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"
	//"math"

	//"github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	clip_gio "github.com/utopiagio/gio/op/clip"
	font_gio "github.com/utopiagio/gio/font"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"

	"github.com/utopiagio/utopia/metrics"
)

func GoCanvas(parent GoObject) (hObj *GoCanvasObj) {
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 100, 100, 100, 1000, 1000, 100, 100},
		
		FocusPolicy: StrongFocus,
		Visible: true,
	}
	hCanvas := &GoCanvasObj{
		GioObject: object,
		GioWidget: widget,
		
		fontSize: theme.TextSize,
		color: theme.ContrastFg,
		background: theme.ContrastBg,
		cornerRadius: 4,
		shaper: theme.Shaper,
	}	
	parent.AddControl(hCanvas)

	return hCanvas
}

type GoCanvasObj struct {
	GioObject
	GioWidget
	//theme *GoThemeObj
	font font_gio.Font
	fontSize unit_gio.Sp
	color GoColor
	background GoColor
	cornerRadius unit_gio.Dp

	items []GoCanvasItem
	//inset layout_gio.Inset
	shaper *text_gio.Shaper
}

func (ob *GoCanvasObj) AddItem(item GoCanvasItem) {
	ob.items = append(ob.items, item)
}

func (ob *GoCanvasObj) AddCircle(radius, centreX, centreY float32) (hCanvasCircle *GoCanvasCircleObj) {
	hCanvasCircle = GoCanvasCircle(ob)
	hCanvasCircle.SetRadius(radius)
	hCanvasCircle.Centre(centreX, centreY)
	ob.items = append(ob.items, hCanvasCircle)
	return hCanvasCircle
}

func (ob *GoCanvasObj) AddEllipse(height, width, centreX, centreY float32) (hCanvasEllipse *GoCanvasEllipseObj) {
	hCanvasEllipse = GoCanvasEllipse(ob)
	hCanvasEllipse.SetSize(height, width)
	hCanvasEllipse.Centre(centreX, centreY)
	ob.items = append(ob.items, hCanvasEllipse)
	return hCanvasEllipse
}

func (ob *GoCanvasObj) AddLine(startX, startY, endX, endY float32) (hCanvasLine *GoCanvasLineObj) {
	hCanvasLine = GoCanvasLine(ob)
	hCanvasLine.SetPoints(startX, startY, endX, endY)
	ob.items = append(ob.items, hCanvasLine)
	return hCanvasLine
}

func (ob *GoCanvasObj) RemoveItem(item GoCanvasItem) {
	k := 0
	for _, v := range ob.items {
	    if v != item {
	        ob.items[k] = v
	        k++
	    }
	}
	ob.items = ob.items[:k] // set slice len to remaining elements
}

func (ob *GoCanvasObj) SetBackgroundColor(color GoColor) {
	ob.background = color
}

/*func (ob *GoCanvasObj) PointerDragged(e pointer_gio.Event) {
	log.Println("Type:", e.Type, "Pos:", e.Position, "Buttons:", e.Buttons, "Scroll:", e.Scroll)
}

func (ob *GoCanvasObj) PointerEntered(e pointer_gio.Event) {
	log.Println("Type:", e.Type, "Pos:", e.Position, "Scroll:", e.Scroll)
}

func (ob *GoCanvasObj) PointerLeft(e pointer_gio.Event) {
	log.Println("Type:", e.Type, "Pos:", e.Position, "Scroll:", e.Scroll)
}

func (ob *GoCanvasObj) PointerMoved(e pointer_gio.Event) {
	log.Println("Type:", e.Type, "Pos:", e.Position, "Scroll:", e.Scroll)
}

func (ob *GoCanvasObj) PointerPressed(e pointer_gio.Event) {
	log.Println("GoCanvasObj::PointerPressed")
	log.Println("Type:", e.Type)
	log.Println("Source:", e.Source)
	log.Println("PointerID:", e.PointerID)
	log.Println("Priority:", e.Priority)
	log.Println("Time:", e.Time)
	log.Println("Buttons:", e.Buttons)
	log.Println("Position:", e.Position)
	log.Println("Scroll:", e.Scroll)
	log.Println("Modifiers:", e.Modifiers)
}

func (ob *GoCanvasObj) PointerReleased(e pointer_gio.Event) {
	log.Println("Type:", e.Type)
	log.Println("Source:", e.Source)
	log.Println("PointerID:", e.PointerID)
	log.Println("Priority:", e.Priority)
	log.Println("Time:", e.Time)
	log.Println("Buttons:", e.Buttons)
	log.Println("Position:", e.Position)
	log.Println("Scroll:", e.Scroll)
	log.Println("Modifiers:", e.Modifiers)
}*/

func (ob *GoCanvasObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	cs := gtx.Constraints
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		w := min(maxWidth, width)			// constrain to ob.MaxWidth
		cs.Min.X = max(minWidth, w)				// constrain to ob.MinWidth 
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = minWidth						// set to ob.MinWidth
		cs.Max.X = minWidth						// set to ob.MinWidth
	case PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = minWidth						// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = maxWidth						// set to ob.MaxWidth
		cs.Min.X = maxWidth						// set to ob.MaxWidth
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		w := min(maxHeight, height)				// constrain to ob.MaxHeight
		cs.Min.Y = max(minHeight, w)			// constrain to ob.MinHeight 
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = minHeight					// set to ob.MinHeight
		cs.Max.Y = minHeight					// set to ob.MinHeight
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = max(0, minHeight)			// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = maxHeight					// set to ob.MaxHeight
		cs.Min.Y = maxHeight					// set to ob.MaxHeight
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
					return ob.Layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoCanvasObj) ObjectType() (string) {
	return "GoCanvasObj"
}

func (ob *GoCanvasObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoCanvasObj) Update() {
	
}

func (ob *GoCanvasObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx)

	// paint object
	width := gtx.Dp(unit_gio.Dp(ob.Width))
	height := gtx.Dp(unit_gio.Dp(ob.Height))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	dims := image.Point{X: width, Y: height}
	rr := gtx.Dp(ob.cornerRadius)
	defer clip_gio.UniformRRect(image.Rectangle{Max: dims}, rr).Push(gtx.Ops).Pop()
	// paint background
	background := ob.background.NRGBA()
	paint_gio.Fill(gtx.Ops, background)

	// paint foreground
	for idx := 0; idx < len(ob.items); idx++ {
		ob.items[idx].Draw(gtx.Ops)
	}
	// add the events handler to receive widget pointer events
	ob.SignalEvents(gtx)

	return layout_gio.Dimensions{Size: dims}
}

