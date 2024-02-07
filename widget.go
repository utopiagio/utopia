// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/widget.go */

package utopia

import (
	"image"
	"log"
	"time"

	f32_gio "github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	key_gio "github.com/utopiagio/gio/io/key"
	paint_gio "github.com/utopiagio/gio/op/paint"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	unit_gio "github.com/utopiagio/gio/unit"
)

// The duration is somewhat arbitrary.
const clickDuration = 200 * time.Millisecond
const doubleClickDuration = 400 * time.Millisecond

type GoFocusPolicy int

const (
	NoFocus GoFocusPolicy = iota 	// the widget does not accept focus.
	TabFocus   						// the widget accepts focus by tabbing.
	ClickFocus  					// the widget accepts focus by clicking.
	StrongFocus 					// the widget accepts focus by both tabbing and clicking. ***
	WheelFocus 						// like StrongFocus plus the widget accepts focus by using the mouse wheel.
	// *** On Mac OS X this will also be indicate that the widget accepts tab focus when in 'Text/List focus mode'.
)

type GoPos struct {
	X int
	Y int
}

func (p GoPos) ImPos() (image.Point) {
	return image.Point{X: p.X, Y: p.Y}
}

type GoSize struct {
	MinWidth int
	MinHeight int
	Width int
	Height int
	MaxWidth int
	MaxHeight int
	AbsWidth int
	AbsHeight int
}

func (s GoSize) ImMax() (image.Point) {
	return image.Point{X: s.MaxWidth, Y: s.MaxHeight}
}

func (s GoSize) ImMin() (image.Point) {
	return image.Point{X: s.MinWidth, Y: s.MinHeight}
}

func (s GoSize) ImSize() (image.Point) {
	return image.Point{X: s.Width, Y: s.Height}
}

// Fit scales a widget to fit and clip to the constraints.
type GoFit uint8

const (
	// Unscaled does not alter the scale of a widget.
	Unscaled GoFit = iota
	// Contain scales widget as large as possible without cropping
	// and it preserves aspect-ratio.
	Contain
	// Cover scales the widget to cover the constraint area and
	// preserves aspect-ratio.
	Cover
	// ScaleDown scales the widget smaller without cropping,
	// when it exceeds the constraint area.
	// It preserves aspect-ratio.
	ScaleDown
	// Fill stretches the widget to the constraints and does not
	// preserve aspect-ratio.
	Fill
)

// scale computes the new dimensions and transformation required to fit dims to cs, given the position.
func (fit GoFit) scale(cs layout_gio.Constraints, pos layout_gio.Direction, dims layout_gio.Dimensions) (layout_gio.Dimensions, f32_gio.Affine2D) {
	widgetSize := dims.Size

	if fit == Unscaled || dims.Size.X == 0 || dims.Size.Y == 0 {
		dims.Size = cs.Constrain(dims.Size)

		offset := pos.Position(widgetSize, dims.Size)
		dims.Baseline += offset.Y
		return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
	}

	scale := f32_gio.Point{
		X: float32(cs.Max.X) / float32(dims.Size.X),
		Y: float32(cs.Max.Y) / float32(dims.Size.Y),
	}

	switch fit {
	case Contain:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case Cover:
		if scale.Y > scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case ScaleDown:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}

		// The widget would need to be scaled up, no change needed.
		if scale.X >= 1 {
			dims.Size = cs.Constrain(dims.Size)

			offset := pos.Position(widgetSize, dims.Size)
			dims.Baseline += offset.Y
			return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
		}
	case Fill:
	}

	var scaledSize image.Point
	scaledSize.X = int(float32(widgetSize.X) * scale.X)
	scaledSize.Y = int(float32(widgetSize.Y) * scale.Y)
	dims.Size = cs.Constrain(scaledSize)
	dims.Baseline = int(float32(dims.Baseline) * scale.Y)

	offset := pos.Position(scaledSize, dims.Size)
	trans := f32_gio.Affine2D{}.
		Scale(f32_gio.Point{}, scale).
		Offset(layout_gio.FPt(offset))

	dims.Baseline += offset.Y

	return dims, trans
}

type GoBorderStyle int

const (
	BorderNone GoBorderStyle = iota
	BorderSingleLine
	BorderSunken
	BorderSunkenThick
	BorderRaised
)

type GoBorder struct {
	BStyle 	GoBorderStyle
	BColor 	GoColor
	BRadius 	int
	BWidth 	int
	FillColor 	GoColor
}

func (b GoBorder) Layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	rr := gtx.Dp(unit_gio.Dp(b.BRadius))
	width := gtx.Dp(unit_gio.Dp(b.BWidth))
	//log.Println("width =", width)

	mcs := gtx.Constraints
	mcs.Max.X -= width * 2
	if mcs.Max.X < 0 {
		width = 0
		mcs.Max.X = 0
	}
	if mcs.Min.X > mcs.Max.X {
		mcs.Min.X = mcs.Max.X
	}
	mcs.Max.Y -= width * 2
	if mcs.Max.Y < 0 {
		width = 0
		mcs.Max.Y = 0
	}
	if mcs.Min.Y > mcs.Max.Y {
		mcs.Min.Y = mcs.Max.Y
	}
	gtx.Constraints = mcs
	
	trans := op_gio.Offset(image.Pt(width, width)).Push(gtx.Ops)
	// Save the operations in an independent ops value (the cache).
	macro := op_gio.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()
	//trans := op_gio.Offset(image.Pt(width, width)).Push(gtx.Ops)
	//dims := w(gtx)

	r := image.Rectangle{Max: dims.Size}
	//log.Println("image.Rect=", r)
	r = r.Inset(-width / 2)

	if b.FillColor > 0x00 {
		paint_gio.FillShape(gtx.Ops,
			b.FillColor.NRGBA(),
			clip_gio.UniformRRect(r, rr).Op(gtx.Ops),
		)
	}
	// Draw the operations from the cache.
	call.Add(gtx.Ops)
	// Paint the Border
	paint_gio.FillShape(gtx.Ops,
		b.BColor.NRGBA(),
		clip_gio.Stroke{
			Path:  clip_gio.UniformRRect(r, rr).Path(gtx.Ops),
			Width: float32(width),
		}.Op(),
	)
	trans.Pop()

	return layout_gio.Dimensions{
		Size:     dims.Size.Add(image.Point{X: width * 2, Y: width * 2}),
		Baseline: dims.Baseline + width,
	}
}

type GoMargin struct {
	Left 	int
	Top 	int
	Right 	int
	Bottom 	int
}

func (m *GoMargin) Layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	top := gtx.Dp(unit_gio.Dp(m.Top))
	right := gtx.Dp(unit_gio.Dp(m.Right))
	bottom := gtx.Dp(unit_gio.Dp(m.Bottom))
	left := gtx.Dp(unit_gio.Dp(m.Left))
	mcs := gtx.Constraints
	mcs.Max.X -= left + right
	if mcs.Max.X < 0 {
		left = 0
		right = 0
		mcs.Max.X = 0
	}
	if mcs.Min.X > mcs.Max.X {
		mcs.Min.X = mcs.Max.X
	}
	mcs.Max.Y -= top + bottom
	if mcs.Max.Y < 0 {
		bottom = 0
		top = 0
		mcs.Max.Y = 0
	}
	if mcs.Min.Y > mcs.Max.Y {
		mcs.Min.Y = mcs.Max.Y
	}
	gtx.Constraints = mcs
	trans := op_gio.Offset(image.Pt(left, top)).Push(gtx.Ops)
	dims := w(gtx)
	trans.Pop()
	return layout_gio.Dimensions{
		Size:     dims.Size.Add(image.Point{X: right + left, Y: top + bottom}),
		Baseline: dims.Baseline + bottom,
	}
}

type GoPadding struct {
	Left 	int
	Top 	int
	Right 	int
	Bottom 	int
}

func (p *GoPadding) Layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	top := gtx.Dp(unit_gio.Dp(p.Top))
	right := gtx.Dp(unit_gio.Dp(p.Right))
	bottom := gtx.Dp(unit_gio.Dp(p.Bottom))
	left := gtx.Dp(unit_gio.Dp(p.Left))
	mcs := gtx.Constraints
	mcs.Max.X -= left + right
	if mcs.Max.X < 0 {
		left = 0
		right = 0
		mcs.Max.X = 0
	}
	if mcs.Min.X > mcs.Max.X {
		mcs.Min.X = mcs.Max.X
	}
	mcs.Max.Y -= top + bottom
	if mcs.Max.Y < 0 {
		bottom = 0
		top = 0
		mcs.Max.Y = 0
	}
	if mcs.Min.Y > mcs.Max.Y {
		mcs.Min.Y = mcs.Max.Y
	}
	gtx.Constraints = mcs
	trans := op_gio.Offset(image.Pt(left, top)).Push(gtx.Ops)
	dims := w(gtx)
	trans.Pop()
	return layout_gio.Dimensions{
		Size:     dims.Size.Add(image.Point{X: right + left, Y: top + bottom}),
		Baseline: dims.Baseline + bottom,
	}
}

type GioWidget struct {
	GoBorder 		// border drawn surrounding widget
	GoMargin		// clear margin surrounding widget
	GoPadding		// clear padding within widget
	GoSize 			// Fixed, Min and Max sizes
	
	dims layout_gio.Dimensions // Dimensions{Size image.Point{X: int, Y: int}, Baseline int}
	

	Visible bool
	// windows 10 System Colors
	BackColor 		GoColor 	// COLOR_WINDOW
	BackgroundColor GoColor 	// COLOR_WINDOW
	ButtonText 		GoColor 	// COLOR_BTNTEXT Foreground color for button text.
	FaceColor		GoColor  	// COLOR_BTNFACE Background color for button.
	ForeColor 		GoColor 	// COLOR_WINDOWTEXT
	GrayText 		GoColor 	// COLOR_GRAYTEXT Foreground color for disabled button text.
	Highlight 		GoColor 	// COLOR_HIGHLIGHT Background color for slected button.
	HighlightText 	GoColor 	// COLOR_HIGHLIGHTTEXT Foreground color for slected button text.
	Hotlight 		GoColor 	// COLOR_HOTLIGHT Hyperlink color.

	onClearFocus func()
	onSetFocus func()

	keys string
	onKeyEdit func(e key_gio.EditEvent)
	onKeyPress func(e key_gio.Event)
	onKeyRelease func(e key_gio.Event)


	onPointerClick func(e pointer_gio.Event)
	onPointerDoubleClick func(e pointer_gio.Event)
	onPointerDrag func(e pointer_gio.Event)
	onPointerEnter func(e pointer_gio.Event)
	onPointerLeave func(e pointer_gio.Event)
	onPointerMove func(e pointer_gio.Event)
	onPointerPress func(e pointer_gio.Event)
	onPointerRelease func(e pointer_gio.Event)

	events pointer_gio.Kind
	clicks int
	clickEvent pointer_gio.Event
	focus bool
	focusEnabled bool
	FocusPolicy GoFocusPolicy
	hovered bool
	selected bool
	//pressedAt time.Duration
	/*clickable *widget_gio.Clickable
	draggable *widget_gio.Draggable
	editor *widget_gio.Bool
	enum *widget_gio.Enum
	label *widget_gio.Label*/
}

func (w *GioWidget) Click(e pointer_gio.Event) {
	if w.onPointerClick != nil {
			w.onPointerClick(e)
		}
}

func (w *GioWidget) Clicked() (clicked bool) {
	log.Println("GioWidget::Clicked", w.clicks)
	if w.clicks == 1 {
		return true
	}
	return false
}

// func (*GioWidget) ClearFocus() (bool)
// Notifies the App to switch the keyboard focus to nil 
func (w *GioWidget) ClearFocus() bool {
	if GoApp.Keyboard().ClearFocus(w) {
		w.focus = false
		if w.onClearFocus != nil {
			w.onClearFocus()
		}
		return true
	}
	return false
}

func (w *GioWidget) HasFocus() bool {
	return w.focus
}

func (w *GioWidget) Hide() {
	w.Visible = false
}

func (w *GioWidget) IsFocusEnabled() bool {
	return w.focusEnabled
}

func (w *GioWidget) IsHovered() bool {
	return w.hovered
}

func (w *GioWidget) IsSelected() bool {
	return w.selected
}

func (w *GioWidget) IsVisible() bool {
	return w.Visible
}

// func (*GioWidget) Margin() (margin GoMargin)
// Returns the margin rect for the positioning dimensions for this goWidget within its parent.
// see goWidget positioning and geometry.
func (w *GioWidget) Margin() (margin GoMargin) {
	return w.GoMargin
}

// func (*GioWidget) Padding() (padding GoPadding)
// Returns the padding dimensions for the positioning of clients within their parent.
// see goWidget positioning and geometry.
func (w *GioWidget) Padding() (padding GoPadding) {
	return w.GoPadding
}

func (w *GioWidget) SetBackgroundColor(color GoColor) {
	//w.BackgroundColor = color
	w.GoBorder.FillColor = color
}

func (w *GioWidget) SetBorder(style GoBorderStyle, width int, radius int, color GoColor) {
	w.GoBorder = GoBorder{style, color, radius, width, w.GoBorder.FillColor}
}

func (w *GioWidget) SetBorderColor(color GoColor) {
	w.GoBorder.BColor = color
}

func (w *GioWidget) SetBorderRadius(radius int) {
	w.GoBorder.BRadius = radius
}

func (w *GioWidget) SetBorderStyle(style GoBorderStyle) {
	color := w.GoBorder.BColor
	switch style {
		case BorderNone:
			w.GoBorder = GoBorder{style, color, 0, 0, 0}
		case BorderSingleLine:
			w.GoBorder = GoBorder{style, color, 0, 1, 0}
		case BorderSunken:
			w.GoBorder = GoBorder{style, color, 0, 2, 0}
		case BorderSunkenThick:
			w.GoBorder = GoBorder{style, color, 0, 4, 0}
		case BorderRaised:
			w.GoBorder = GoBorder{style, color, 0, 4, 0}	
	}

}

func (w *GioWidget) SetBorderWidth(width int) {
	w.GoBorder.BWidth = width
}

// func (*GioWidget) LostFocus(f)
// Sets the keyboard focus on this widget
/*func (w *GioWidget) LostFocus() {
	GoApp.Keyboard().SetFocus(w)
	w.focus = focus
}*/

// func (*GioWidget) SetFocus() (bool)
// Notifies the App to switch the keyboard focus to this widget 
func (w *GioWidget) SetFocus() bool {
	//log.Println("GioWidget::SetFocus()")
	if w.FocusPolicy != NoFocus {
		if GoApp.Keyboard().SetFocus(w) {
			//w.focus = true
			//log.Println("w.onSetFocus =", w.onSetFocus)
			if w.onSetFocus != nil{
				//log.Println("w.onSetFocus()")
				w.onSetFocus()
			}
			return true
		}
	}
	return false
}

// func (*GioWidget) SetFocusPolicy(focusPolicy GoFocusPolicy)
// Sets the keyboard focus policy for this widget 
func (w *GioWidget) SetFocusPolicy(focusPolicy GoFocusPolicy) {
	w.FocusPolicy = focusPolicy
	if focusPolicy != NoFocus {
		w.focusEnabled = true
	}
}

// func (*GioWidget) ChangeFocus(focus bool) (bool)
// Changes focus on the widget if it can accept keyboard focus
func (w *GioWidget) ChangeFocus(focus bool) bool {
	//log.Println("GioWidget::ChangeFocus()")
	if w.FocusPolicy != NoFocus {
		w.focus = focus
		if w.onSetFocus != nil{
			//log.Println("w.onSetFocus()")
			w.onSetFocus()
		}
		return true
	}
	return false
}

func (w *GioWidget) SetHeight(height int) {
	w.GoSize.Height = height
}

func (w *GioWidget) SetMargin(left int, top int, right int, bottom int) {
	w.GoMargin = GoMargin{left, top, right, bottom}
}

func (w *GioWidget) SetMaxHeight(maxHeight int) {
	w.GoSize.MaxHeight = maxHeight
}

func (w *GioWidget) SetMinHeight(minHeight int) {
	w.GoSize.MinHeight = minHeight
}

func (w *GioWidget) SetMaxWidth(maxWidth int) {
	w.GoSize.MaxWidth = maxWidth
}

func (w *GioWidget) SetMinWidth(minWidth int) {
	w.GoSize.MinWidth = minWidth
}

func (w *GioWidget) SetOnClearFocus(f func()) {
	w.onClearFocus = f
}

func (w *GioWidget) SetOnSetFocus(f func()) {
	w.onSetFocus = f
}

func (w *GioWidget) SetOnKeyEdit(f func(e key_gio.EditEvent)) {
	w.onKeyEdit = f
}

func (w *GioWidget) SetOnKeyPress(f func(e key_gio.Event)) {
	w.onKeyPress = f
}

func (w *GioWidget) SetOnKeyRelease(f func(e key_gio.Event)) {
	w.onKeyRelease = f
}

func (w *GioWidget) SetOnPointerClick(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Press | pointer_gio.Release
	w.onPointerClick = f
}

func (w *GioWidget) SetOnPointerDoubleClick(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Press | pointer_gio.Release
	w.onPointerDoubleClick = f
}

func (w *GioWidget) SetOnPointerDrag(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Drag
	w.onPointerDrag = f
}

func (w *GioWidget) SetOnPointerEnter(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Enter
	w.onPointerEnter = f
}

func (w *GioWidget) SetOnPointerLeave(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Leave
	w.onPointerLeave = f
}

func (w *GioWidget) SetOnPointerMove(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Move
	w.onPointerMove = f
}

func (w *GioWidget) SetOnPointerPress(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Press
	w.onPointerPress = f
}

func (w *GioWidget) SetOnPointerRelease(f func(e pointer_gio.Event)) {
	w.events = w.events | pointer_gio.Release
	w.onPointerRelease = f
}


func (w *GioWidget) SetPadding(left int, top int, right int, bottom int) {
	w.GoPadding = GoPadding{left, top, right, bottom}
}

func (w *GioWidget) SetSelected(selected bool) {
	w.selected = selected
}

func (w *GioWidget) SetWidth(width int) {
	w.GoSize.Width = width
}

func (w *GioWidget) Show() {
	w.Visible = true
}

func (w *GioWidget) Size() (GoSize){
	return w.GoSize
}

func (w *GioWidget) SignalEvents(gtx layout_gio.Context) {
	//log.Println("GioWidget::SignalEvents", w.events)
	if w.events != 0 {
		pointer_gio.InputOp{
			Tag:   w,
			Grab:  false,
			Kinds: w.events,
		}.Add(gtx.Ops)
	}
	if w.focus {

		key_gio.FocusOp{
			Tag: w, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
		}.Add(gtx.Ops)

		// 3) Finally we add key.InputOp to catch specific keys
		// (Shift) means an optional Shift
		// These inputs are retrieved as key.Event
		key_gio.InputOp{
			Keys: key_gio.Set(w.keys),
			Tag:  w, // Use Tag: w as the event routing tag, and retrieve it through gtx.Events(w) in GioWidget::ReceiveEvents() routine.
		}.Add(gtx.Ops)
	}
}

func (w *GioWidget) ReceiveEvents(gtx layout_gio.Context) {
	for _, gtxEvent := range gtx.Events(w) {
		switch e := gtxEvent.(type) {
			case key_gio.EditEvent:
				//log.Println("WidgetKey::EditEvent -", "Range -", e.Range, "Text -", e.Text)
				if w.onKeyEdit != nil {
					//log.Println("WidgetKey::onKeyEdit() -")
					w.onKeyEdit(e)
				}
			case key_gio.Event:
				switch e.State {
					case key_gio.Press:
						//log.Println("WidgetKey::Event -", "Name -", e.Name, "Modifiers -", e.Modifiers, "State -", e.State)
						if w.onKeyPress != nil {
							//log.Println("WidgetKey::onKeyPress() -")
							w.onKeyPress(e)
						}
					case key_gio.Release:
						//log.Println("WidgetKey::Event -", "Name -", e.Name, "Modifiers -", e.Modifiers, "State -", e.State)
						if w.onKeyRelease != nil {
							w.onKeyRelease(e)
						}
				}
			case pointer_gio.Event:
				switch e.Kind {
					case pointer_gio.Press:
						//log.Println("MousePress:")
						//log.Println("e.Time: ", uint(e.Time))
						//log.Println("w.clickEvent.Time: ", uint(w.clickEvent.Time))
						//log.Println("doubleClickDuration: ", uint(doubleClickDuration))

						//log.Println("Duration: ", uint(e.Time - w.clickEvent.Time))

						if w.FocusPolicy >= ClickFocus && w.focus == false {
								//log.Println("GoApp.Keyboard().SetFocusControl(GoWidget)")
								//GoApp.Keyboard().SetFocusControl(w)
								log.Println("GioWidget.SetFocus() -")

								w.SetFocus()
						}
						if w.clickEvent.Time == 0 {
							w.clickEvent = e
							//log.Println("w.onPointerPress()")
							if w.onPointerPress != nil {
								//log.Println("w.onPointerPress() != nil")
								w.onPointerPress(e)
							}
						} else {
							if e.Time - w.clickEvent.Time < doubleClickDuration {
								//log.Println("MouseDoubleClick:")
								//log.Println("GoApp.Keyboard().SetFocusControl(GoWidget)")
								//GoApp.Keyboard().SetFocusControl(w)
								w.clicks = 2
							}
						}
					case pointer_gio.Release:
						//log.Println("MouseRelease:")
						//log.Println("e.Time: ", uint(e.Time))
						//log.Println("w.clickEvent.Time: ", uint(w.clickEvent.Time))
						//log.Println("clickDuration: ", uint(clickDuration))
						//log.Println("doubleClickDuration: ", uint(doubleClickDuration))

						//log.Println("Duration: ", uint(e.Time - w.clickEvent.Time))
						if e.Time - w.clickEvent.Time < clickDuration {
							//log.Println("MouseClick:")
							// call go routine
							go w.pointerClicked()
							w.clicks = 1
						} else {
							w.clickEvent.Time = 0
						}
						if w.onPointerRelease != nil {
							w.onPointerRelease(e)
						}
					case pointer_gio.Move:
						if w.onPointerMove != nil {
							w.onPointerMove(e)
						}
					case pointer_gio.Drag:
						if w.onPointerDrag != nil {
							w.onPointerDrag(e)
						}
					case pointer_gio.Enter:
						w.hovered = true
						if w.onPointerEnter != nil {
							w.onPointerEnter(e)
						}
					case pointer_gio.Leave:
						w.hovered = false
						if w.onPointerLeave != nil {
							w.onPointerLeave(e)
						}
				}
		}
	}
}

func (w *GioWidget) pointerClicked() {
	time.Sleep(doubleClickDuration)
	w.clickEvent.Time = 0
	if w.clicks == 1 {
		//log.Println("POINTER clicked:")
		if w.onPointerClick != nil {
			w.onPointerClick(w.clickEvent)
		}

	} else if w.clicks == 2 {
		//log.Println("POINTER doubleclicked:")
		if w.onPointerDoubleClick != nil {
			w.onPointerDoubleClick(w.clickEvent)
		}
	}
	
}