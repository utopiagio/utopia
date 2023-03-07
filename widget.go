/* widget.go */

package utopia

import (
	"image"
	//"log"
	//"image/color"

	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	"github.com/utopiagio/gio/unit"
)

type GoBorderStyle int

const (
	BorderNone GoBorderStyle = iota
	BorderSingleLine
	BorderSunken
	BorderSunkenThick
	BorderRaised
)

type goBorder struct {
	style 	GoBorderStyle
	color 	GoColor
	radius 	int
	width 	int
	
}

func (b goBorder) layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	dims := w(gtx)
	sz := dims.Size

	rr := gtx.Dp(unit.Dp(b.radius))
	width := gtx.Dp(unit.Dp(b.width))
	whalf := (width + 1) / 2
	sz.X -= whalf * 2
	sz.Y -= whalf * 2

	r := image.Rectangle{Max: sz}
	r = r.Add(image.Point{X: whalf, Y: whalf})

	paint_gio.FillShape(gtx.Ops,
		b.color.NRGBA(),
		clip_gio.Stroke{
			Path:  clip_gio.UniformRRect(r, rr).Path(gtx.Ops),
			Width: float32(width),
		}.Op(),
	)

	return dims
}

/*func GoMargin(left int, top int, right int, bottom int) *GoMarginType {
	return &GoMarginType{left, top, right, bottom}
}*/

type goMargin struct {
	left 	int
	top 	int
	right 	int
	bottom 	int
}

func (m *goMargin) layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	top := gtx.Dp(unit.Dp(m.top))
	right := gtx.Dp(unit.Dp(m.right))
	bottom := gtx.Dp(unit.Dp(m.bottom))
	left := gtx.Dp(unit.Dp(m.left))
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

/*func GoPadding(left int, top int, right int, bottom int) *GoPaddingType {
	return &GoPaddingType{left, top, right, bottom}
}*/

type goPadding struct {
	left 	int
	top 	int
	right 	int
	bottom 	int
}

func (p *goPadding) layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	top := gtx.Dp(unit.Dp(p.top))
	right := gtx.Dp(unit.Dp(p.right))
	bottom := gtx.Dp(unit.Dp(p.bottom))
	left := gtx.Dp(unit.Dp(p.left))
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

type goWidget struct {
	goBorder 
	goMargin		// clear margin surrounding widget
	goPadding
	visible bool
	// windows 10 System Colors
	backcolor 		GoColor 	// COLOR_WINDOW
	backgroundcolor GoColor 	// COLOR_WINDOW
	buttonText 		GoColor 	// COLOR_BTNTEXT Foreground color for button text.
	facecolor		GoColor  	// COLOR_BTNFACE Background color for button.
	forecolor 		GoColor 	// COLOR_WINDOWTEXT
	grayText 		GoColor 	// COLOR_GRAYTEXT Foreground color for disabled button text.
	highlight 		GoColor 	// COLOR_HIGHLIGHT Background color for slected button.
	highlightText 	GoColor 	// COLOR_HIGHLIGHTTEXT Foreground color for slected button text.
	hotlight 		GoColor 	// COLOR_HOTLIGHT Hyperlink color.
	//onClick func()
	/*clickable *widget_gio.Clickable
	draggable *widget_gio.Draggable
	editor *widget_gio.Bool
	enum *widget_gio.Enum
	label *widget_gio.Label*/
}

/*func (w *goWidget) Click() {
	return w.clickable.Click()
}

func (w *goWidget) Clicked() bool {
	return w.clickable.Clicked()
}

func (w *goWidget) Clicks() []widget_gio.Click {
	return w.clickable.Clicks()
}

func (w *goWidget) Focus() {
	return w.clickable.Focus()
}

func (w *goWidget) Focused() bool {
	return w.clickable.Focused()
}

func (w *goWidget) Hovered() bool {
	return w.clickable.Hovered()
}

func (w *goWidget) Layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	return w.clickable.Layout(gtx layout_gio.Context, w layout_gio.Widget)
}*/

/*func (w *goWidget) Clickable() *widget_gio.Clickable {
	return w.clickable
}*/

/*func (w *goWidget) Pressed() bool {
	return w.clickable.Pressed()
}*/

func (w *goWidget) Hide() {
	w.visible = false
}

func (w *goWidget) SetBorder(style GoBorderStyle, width int, radius int, color GoColor) {
	w.goBorder = goBorder{style, color, radius, width}
}

func (w *goWidget) SetBorderColor(color GoColor) {
	w.goBorder.color = color
}

func (w *goWidget) SetBorderRadius(radius int) {
	w.goBorder.radius = radius
}

func (w *goWidget) SetBorderStyle(style GoBorderStyle) {
	color := w.goBorder.color
	switch style {
		case BorderNone:
			w.goBorder = goBorder{style, color, 0, 0}
		case BorderSingleLine:
			w.goBorder = goBorder{style, color, 0, 1}
		case BorderSunken:
			w.goBorder = goBorder{style, color, 0, 2}
		case BorderSunkenThick:
			w.goBorder = goBorder{style, color, 0, 4}
		case BorderRaised:
			w.goBorder = goBorder{style, color, 0, 4}	
	}

}

func (w *goWidget) SetBorderWidth(width int) {
	w.goBorder.width = width
}

// func (*goWidget) Margin()
// Returns the margin rect for the positioning dimensions for this goWidget within its parent.
// see goWidget positioning and geometry.
func (w *goWidget) Margin() (margin goMargin) {
	return w.goMargin
}

// func (*goWidget) Padding()
// Returns the padding dimensions for the positioning of clients within their parent.
// see goWidget positioning and geometry.
func (w *goWidget) Padding() (padding goPadding) {
	return w.goPadding
}

func (w *goWidget) SetMargin(left int, top int, right int, bottom int) {
	w.goMargin = goMargin{left, top, right, bottom}
}

func (w *goWidget) SetPadding(left int, top int, right int, bottom int) {
	w.goPadding = goPadding{left, top, right, bottom}
}

func (w *goWidget) Show() {
	w.visible = true
}