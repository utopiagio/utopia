/* label.go */

package utopia

import (
	"log"
	"image"
	//"image/color"

	//gofont_gio "github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	
	"golang.org/x/image/math/fixed"

)

type GoTextAlignment uint8

const (
	TextStart GoTextAlignment = iota
	TextEnd
	TextMiddle
)

// Typeface identifies a particular typeface design. The empty
// string denotes the default typeface.
//type gio.Typeface string

// Variant denotes a typeface variant such as "Mono" or "Smallcaps".
//type gio.Variant string

// Style is the font style.
type GoFontStyle int

const (
	Regular GoFontStyle = iota
	Italic
)

// Weight is a font weight, in CSS units subtracted 400 so the zero value
// is normal text weight.
type GoFontWeight int

const (
	Thin       GoFontWeight = 100 - 400
	Hairline   GoFontWeight = Thin
	ExtraLight GoFontWeight = 200 - 400
	UltraLight GoFontWeight = ExtraLight
	Light      GoFontWeight = 300 - 400
	Normal     GoFontWeight = 400 - 400
	Medium     GoFontWeight = 500 - 400
	SemiBold   GoFontWeight = 600 - 400
	DemiBold   GoFontWeight = SemiBold
	Bold       GoFontWeight = 700 - 400
	ExtraBold  GoFontWeight = 800 - 400
	UltraBold  GoFontWeight = ExtraBold
	Black      GoFontWeight = 900 - 400
	Heavy      GoFontWeight = Black
	ExtraBlack GoFontWeight = 950 - 400
	UltraBlack GoFontWeight = ExtraBlack
)


func GoLabel(parent GoObject, text string) (hObj *GoLabelObj) {
	//var fontSize unit_gio.Sp = 12
	theme := GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		FocusPolicy: NoFocus,
		Visible: true,
		//target: nil,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	theme.TextSize,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	NRGBAColor(MulAlpha(theme.ContrastBg.NRGBA(), 0x60)),
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := ColorFromRGB(127, 0, 0)
	//hLabel.color = maroon
	//hLabel.font.Weight = text_gio.Bold 			// Thin - Medium - Bold
	//hLabel.font.Style = text_gio.Italic 		// Regular - Italic
	//hLabel.font.Variant = "Mono"	// Mono - Smallcaps
	
	parent.AddControl(hLabel)
	return hLabel
}

type GoLabelObj struct {
	GioObject
	GioWidget
	//GioSelectable
	alignment text_gio.Alignment
	font text_gio.Font
	fontSize unit_gio.Sp
	maxLines int
	text string
	color GoColor
	selectionColor GoColor
	//textAlign text.Alignment
	selectable bool
	shaper *text_gio.Shaper
	state *GioSelectable
}

/*
type Font struct {
	Typeface Typeface
	Variant  Variant
	Style    Style
	// Weight is the text weight. If zero, Normal is used instead.
	Weight Weight
}

	// Face defines the text style.
	Font text.Font
	// Color is the text color.
	Color color.NRGBA
	// SelectionColor is the color of the background for selected text.
	SelectionColor color.NRGBA
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Sp

	shaper *text.Shaper
	State  *widget.Selectable
*/

func (ob *GoLabelObj) GotFocus() {
	log.Println("GoLabelObj::GotFocus()")
	ob.state.focused = true
}

func (ob *GoLabelObj) LostFocus() {
	log.Println("GoLabelObj::LostFocus()")
	if ob.selectable == true {
		ob.state.focused = false
		ob.state.text.ClearSelection()
	}
	ob.focus = false
}

func (ob *GoLabelObj) SetFont(typeface string, variant string, style GoFontStyle, weight GoFontWeight) {
	ob.font = text_gio.Font{text_gio.Typeface(typeface), text_gio.Variant(variant), text_gio.Style(int(style)), text_gio.Weight(int(weight))}
}

func (ob *GoLabelObj) SetFontBold(bold bool) {
	if bold {
		ob.font.Weight = text_gio.Bold
	} else {
		ob.font.Weight = text_gio.Light
	}
}

func (ob *GoLabelObj) SetFontSize(size int) {
	ob.fontSize = unit_gio.Sp(size)
}

func (ob *GoLabelObj) SetFontVariant(variant string) {
	ob.font.Variant = text_gio.Variant(variant)
}

func (ob *GoLabelObj) SetFontWeight(weight GoFontWeight) {
	ob.font.Weight = text_gio.Weight(int(weight))
}

func (ob *GoLabelObj) SetHiliteColor(color GoColor) {
	ob.selectionColor = color
}

func (ob *GoLabelObj) SetFontItalic(italic bool) {
	if italic {
		ob.font.Style = text_gio.Italic
	} else {
		ob.font.Style = text_gio.Regular
	}
}

func (ob *GoLabelObj) SetMaxLines(size int) {
	ob.maxLines = size
}

func (ob *GoLabelObj) pointerDoubleClicked(e pointer_gio.Event) {
	log.Println("GoLabelObj::pointerDoubleClicked()")
	if ob.selectable == true {
		ob.state.pointerDoubleClicked(e)
		ob.ParentWindow().Refresh()
	}
}

func (ob *GoLabelObj) pointerDragged(e pointer_gio.Event) {
	if ob.selectable == true {
		ob.state.pointerDragged(e)
	}
}

func (ob *GoLabelObj) pointerPressed(e pointer_gio.Event) {
	log.Println("GoLabelObj::pointerPressed()")
	if ob.selectable == true {
		ob.state.pointerPressed(e)
	}
}

func (ob *GoLabelObj) pointerReleased(e pointer_gio.Event) {
	if ob.selectable == true {
		ob.state.pointerReleased(e)
	}
	//ob.editor.focused = true
}


func (ob *GoLabelObj) SetSelectable(selectable bool) {
	ob.selectable = selectable
	if selectable {
		ob.SetFocusPolicy(StrongFocus)
		ob.state = &GioSelectable{}
		ob.state.text.Alignment = ob.alignment
		ob.state.text.MaxLines = ob.maxLines
		ob.state.SetText(ob.text)
		ob.SetOnSetFocus(ob.GotFocus)
		ob.SetOnClearFocus(ob.LostFocus)
		ob.SetOnPointerPress(ob.pointerPressed)
		ob.SetOnPointerDrag(ob.pointerDragged)
		ob.SetOnPointerRelease(ob.pointerReleased)
		ob.SetOnPointerDoubleClick(ob.pointerDoubleClicked)
	} else {
		ob.SetFocusPolicy(NoFocus)
		ob.state = nil
		ob.SetOnSetFocus(nil)
		ob.SetOnClearFocus(nil)
		ob.SetOnPointerPress(nil)
		ob.SetOnPointerDrag(nil)
		ob.SetOnPointerRelease(nil)
		ob.SetOnPointerDoubleClick(nil)
	}
}

/*func (ob *GoLabelObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoLabelObj) SetText(text string) {
	ob.text = text
	if ob.selectable == true {
		ob.state.SetText(ob.text)
	}
}

func (ob *GoLabelObj) SetTextAlignment(alignment GoTextAlignment) {
	ob.alignment = text_gio.Alignment(uint8(alignment))
	if ob.selectable == true {
		ob.state.text.Alignment = ob.alignment
	}
}

func (ob *GoLabelObj) SetTextColor(color GoColor) {
	ob.color = color
}

func (ob *GoLabelObj) Text() (text string) {
	return ob.text
}

func (ob *GoLabelObj) TextColor() (color GoColor) {
	return ob.color
}

func (ob *GoLabelObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}

func (ob *GoLabelObj) layout(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	//lbl := widget_gio.Label{Alignment: ob.alignment, MaxLines: ob.maxLines} //, Selectable: ob.state}
	//if ob.state == nil {
	if ob.selectable == false {
		return ob.render(gtx, ob.shaper, ob.font, ob.fontSize, ob.text)
	} else {
		ob.ReceiveEvents(gtx)
		dims := ob.state.renderSelectable(gtx, ob.shaper, ob.font, ob.fontSize, func(gtx layout_gio.Context) layout_gio.Dimensions {
			paint_gio.ColorOp{Color: ob.selectionColor.NRGBA()}.Add(gtx.Ops)
			ob.state.PaintSelection(gtx)
			paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			ob.state.PaintText(gtx)
			
			return layout_gio.Dimensions{}
		})
		defer clip_gio.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
		//log.Println("dims.Size:", dims.Size)
		pointer_gio.CursorText.Add(gtx.Ops)
		// add the events handler to receive widget pointer events
		ob.SignalEvents(gtx)
		return dims
	}
}

func (ob *GoLabelObj) render(gtx layout_gio.Context, lt *text_gio.Shaper, font text_gio.Font, size unit_gio.Sp, txt string) layout_gio.Dimensions {
	cs := gtx.Constraints
	textSize := fixed.I(gtx.Sp(size))
	lt.LayoutString(text_gio.Parameters{
		Font:      font,
		PxPerEm:   textSize,
		MaxLines:  ob.maxLines,
		Alignment: ob.alignment,
	}, cs.Min.X, cs.Max.X, gtx.Locale, txt)
	m := op_gio.Record(gtx.Ops)
	viewport := image.Rectangle{Max: cs.Max}
	it := textIterator{viewport: viewport, maxLines: ob.maxLines}
	semantic_gio.LabelOp(txt).Add(gtx.Ops)
	var glyphs [32]text_gio.Glyph
	line := glyphs[:0]
	for g, ok := lt.NextGlyph(); ok; g, ok = lt.NextGlyph() {
		var ok bool
		if line, ok = it.paintGlyph(gtx, lt, g, line); !ok {
			break
		}
	}
	call := m.Stop()
	viewport.Min = viewport.Min.Add(it.padding.Min)
	viewport.Max = viewport.Max.Add(it.padding.Max)
	clipStack := clip_gio.Rect(viewport).Push(gtx.Ops)
	call.Add(gtx.Ops)
	dims := layout_gio.Dimensions{Size: it.bounds.Size()}
	dims.Size = cs.Constrain(dims.Size)
	dims.Baseline = dims.Size.Y - it.baseline
	clipStack.Pop()
	return dims
}

func (ob *GoLabelObj) ObjectType() (string) {
	return "GoLabelObj"
}

func (ob *GoLabelObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func H1Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
		//target: nil,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	32,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := ColorFromRGB(127, 0, 0)
	//hLabel.color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.AddControl(hLabel)
	return hLabel
}

func H2Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	26,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.AddControl(hLabel)
	return hLabel
}

func H3Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	22,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.AddControl(hLabel)
	return hLabel
}

func H4Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	20,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.AddControl(hLabel)
	return hLabel
}

func H5Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	16,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.AddControl(hLabel)
	return hLabel
}

func H6Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hLabel := &GoLabelObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	12,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		//state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.AddControl(hLabel)
	return hLabel
}

// textIterator computes the bounding box of and paints text.
type textIterator struct {
	// viewport is the rectangle of document coordinates that the iterator is
	// trying to fill with text.
	viewport image.Rectangle
	// maxLines is the maximum number of text lines that should be displayed.
	maxLines int

	// linesSeen tracks the quantity of line endings this iterator has seen.
	linesSeen int
	// lineOff tracks the origin for the glyphs in the current line.
	lineOff image.Point
	// padding is the space needed outside of the bounds of the text to ensure no
	// part of a glyph is clipped.
	padding image.Rectangle
	// bounds is the logical bounding box of the text.
	bounds image.Rectangle
	// visible tracks whether the most recently iterated glyph is visible within
	// the viewport.
	visible bool
	// first tracks whether the iterator has processed a glyph yet.
	first bool
	// baseline tracks the location of the first line of text's baseline.
	baseline int
}

// processGlyph checks whether the glyph is visible within the iterator's configured
// viewport and (if so) updates the iterator's text dimensions to include the glyph.
func (it *textIterator) processGlyph(g text_gio.Glyph, ok bool) (_ text_gio.Glyph, visibleOrBefore bool) {
	if it.maxLines > 0 {
		if g.Flags&text_gio.FlagLineBreak != 0 {
			it.linesSeen++
		}
		if it.linesSeen == it.maxLines && g.Flags&text_gio.FlagParagraphBreak != 0 {
			return g, false
		}
	}
	// Compute the maximum extent to which glyphs overhang on the horizontal
	// axis.
	if d := g.Bounds.Min.X.Floor(); d < it.padding.Min.X {
		it.padding.Min.X = d
	}
	if d := (g.Bounds.Max.X - g.Advance).Ceil(); d > it.padding.Max.X {
		it.padding.Max.X = d
	}
	logicalBounds := image.Rectangle{
		Min: image.Pt(g.X.Floor(), int(g.Y)-g.Ascent.Ceil()),
		Max: image.Pt((g.X + g.Advance).Ceil(), int(g.Y)+g.Descent.Ceil()),
	}
	if !it.first {
		it.first = true
		it.baseline = int(g.Y)
		it.bounds = logicalBounds
	}

	above := logicalBounds.Max.Y < it.viewport.Min.Y
	below := logicalBounds.Min.Y > it.viewport.Max.Y
	left := logicalBounds.Max.X < it.viewport.Min.X
	right := logicalBounds.Min.X > it.viewport.Max.X
	it.visible = !above && !below && !left && !right
	if it.visible {
		it.bounds.Min.X = it.minValue(it.bounds.Min.X, logicalBounds.Min.X)
		it.bounds.Min.Y = it.minValue(it.bounds.Min.Y, logicalBounds.Min.Y)
		it.bounds.Max.X = it.maxValue(it.bounds.Max.X, logicalBounds.Max.X)
		it.bounds.Max.Y = it.maxValue(it.bounds.Max.Y, logicalBounds.Max.Y)
	}
	return g, ok && !below
}

// paintGlyph buffers up and paints text glyphs. It should be invoked iteratively upon each glyph
// until it returns false. The line parameter should be a slice with
// a backing array of sufficient size to buffer multiple glyphs.
// A modified slice will be returned with each invocation, and is
// expected to be passed back in on the following invocation.
// This design is awkward, but prevents the line slice from escaping
// to the heap.
func (it *textIterator) paintGlyph(gtx layout_gio.Context, shaper *text_gio.Shaper, glyph text_gio.Glyph, line []text_gio.Glyph) ([]text_gio.Glyph, bool) {
	_, visibleOrBefore := it.processGlyph(glyph, true)
	if it.visible {
		if len(line) == 0 {
			it.lineOff = image.Point{X: glyph.X.Floor(), Y: int(glyph.Y)}.Sub(it.viewport.Min)
		}
		line = append(line, glyph)
	}
	if glyph.Flags&text_gio.FlagLineBreak != 0 || cap(line)-len(line) == 0 || !visibleOrBefore {
		t := op_gio.Offset(it.lineOff).Push(gtx.Ops)
		op := clip_gio.Outline{Path: shaper.Shape(line)}.Op().Push(gtx.Ops)
		paint_gio.PaintOp{}.Add(gtx.Ops)
		op.Pop()
		t.Pop()
		line = line[:0]
	}
	return line, visibleOrBefore
}

func (it *textIterator) maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (it *textIterator) minValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}