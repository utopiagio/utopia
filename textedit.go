// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"log"
	"image"
	"image/color"

	//"github.com/utopiagio/gio/internal/f32color"
	clip_gio "github.com/utopiagio/gio/op/clip"
	key_gio "github.com/utopiagio/gio/io/key"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoTextEditObj struct {
	GioObject
	GioWidget
	font     text_gio.Font
	fontSize unit_gio.Sp
	// Color is the text color.
	color GoColor
	// Hint contains the text displayed when the editor is empty.
	hint string
	// HintColor is the color of hint text.
	hintColor GoColor
	// SelectionColor is the color of the background for selected text.
	selectionColor GoColor
	editor         *Editor

	shaper *text_gio.Shaper

	onFocus func()
}

func GoTextEdit(parent GoObject, hintText string) *GoTextEditObj {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		FocusPolicy: StrongFocus,
		Visible: true,
		keys: "←|→|↑|↓|⏎|⌤|⎋|⇱|⇲|⌫|⌦|⇞|⇟",
		//target: nil,
	}
	hTextEdit := &GoTextEditObj{
		GioObject: object,
		GioWidget: widget,
		fontSize: 	theme.TextSize,
		color: 		theme.ColorFg,
		hint: 		hintText,
		hintColor: 		NRGBAColor(MulAlpha(theme.ColorFg.NRGBA(), 0xbb)),
		selectionColor:	NRGBAColor(MulAlpha(theme.ContrastBg.NRGBA(), 0x60)),
		editor: 	new(Editor),
		shaper: 	theme.Shaper,
	}
	hTextEdit.SetOnSetFocus(hTextEdit.GotFocus)
	hTextEdit.SetOnClearFocus(hTextEdit.LostFocus)
	hTextEdit.SetOnKeyEdit(hTextEdit.KeyEdit)
	hTextEdit.SetOnKeyPress(hTextEdit.KeyPressed)
	//hTextEdit.SetOnKeyRelease(hTextEdit.KeyReleased)
	hTextEdit.SetOnPointerDrag(hTextEdit.PointerDragged)
	hTextEdit.SetOnPointerPress(hTextEdit.PointerPressed)
	hTextEdit.SetOnPointerRelease(hTextEdit.PointerReleased)
	parent.AddControl(hTextEdit)
	return hTextEdit
}

func (ob *GoTextEditObj) ClearSelection() {
	ob.editor.ClearSelection()
}

func (ob *GoTextEditObj) Font() text_gio.Font {
	return ob.font
}

func (ob *GoTextEditObj) Focused() bool {
	return ob.editor.Focused()
}

func (ob *GoTextEditObj) GotFocus() {
	log.Println("GoTextEditObj::GotFocus()")
	ob.editor.focused = true
}

func (ob *GoTextEditObj) LostFocus() {
	log.Println("GoTextEditObj::LostFocus()")
	ob.editor.focused = false
	ob.focus = false
}

func (ob *GoTextEditObj) Insert(text string) {
	ob.editor.Insert(text)
}

func (ob *GoTextEditObj) KeyEdit(e key_gio.EditEvent) {
	log.Println("GoTextEditObj::KeyEdit()")
	ob.Insert(e.Text)
}

func (ob *GoTextEditObj) KeyPressed(e key_gio.Event) {
	log.Println("GoTextEditObj::KeyPressed()")
	ob.editor.processKey(e)
	//log.Println("REFRESH........")
	ob.ParentWindow().Refresh()
}

func (ob *GoTextEditObj) KeyReleased(e key_gio.Event) {
	log.Println("GoTextEditObj::KeyReleased()")
	//ob.editor.Insert(text)
}

func (ob *GoTextEditObj) Length() (length int) {
	return ob.editor.Len()
}

/*func (ob *GoTextEditObj) MoveCoord(pos image.Point) {
	x := fixed.I(pos.X + ob.scrollOff.X)
	y := pos.Y + ob.scrollOff.Y
	ob.caret.start = ob.closestToXY(x, y).runes
	ob.caret.xoff = 0
}

func (ob *GoTextEditObj) MoveCaret(startDelta, endDelta int) {
	ob.caret.xoff = 0
	ob.caret.start = ob.closestToRune(ob.caret.start + startDelta).runes
	ob.caret.end = ob.closestToRune(ob.caret.end + endDelta).runes
}

func (ob *GoTextEditObj) MoveStart(selAct selectionAction) {
	caret := ob.closestToRune(ob.caret.start)
	caret = ob.closestToLineCol(caret.lineCol.line, 0)
	ob.caret.start = caret.runes
	ob.caret.xoff = -caret.x
	ob.updateSelection(selAct)
}

func (ob *GoTextEditObj) MoveEnd(selAct selectionAction) {
	caret := ob.closestToRune(e.caret.start)
	caret = ob.closestToLineCol(caret.lineCol.line, math.MaxInt)
	ob.caret.start = caret.runes
	ob.caret.xoff = fixed.I(ob.maxWidth) - caret.x
	ob.updateSelection(selAct)
}

func (ob *GoTextEditObj) MoveLines(distance int, selAct selectionAction) {
	caretStart := ob.closestToRune(ob.caret.start)
	x := caretStart.x + ob.caret.xoff
	// Seek to line.
	pos := ob.closestToLineCol(caretStart.lineCol.line+distance, 0)
	pos = ob.closestToXY(x, pos.y)
	ob.caret.start = pos.runes
	ob.caret.xoff = x - pos.x
	ob.updateSelection(selAct)
}

func (ob *GoTextEditObj) MovePages(pages int, selAct selectionAction) {
	caret := ob.closestToRune(ob.caret.start)
	x := caret.x + ob.caret.xoff
	y := caret.y + pages*ob.viewSize.Y
	pos := ob.closestToXY(x, y)
	ob.caret.start = pos.runes
	ob.caret.xoff = x - pos.x
	ob.updateSelection(selAct)
}*/

func (ob *GoTextEditObj) PointerDragged(e pointer_gio.Event) {
	ob.editor.pointerDragged(e)
}

func (ob *GoTextEditObj) PointerPressed(e pointer_gio.Event) {
	ob.editor.pointerPressed(e)
}

func (ob *GoTextEditObj) PointerReleased(e pointer_gio.Event) {
	ob.editor.pointerReleased(e)
	//ob.editor.focused = true
}

func (ob *GoTextEditObj) SetFont(typeface string, variant string, style GoFontStyle, weight GoFontWeight) {
	ob.font = text_gio.Font{text_gio.Typeface(typeface), text_gio.Variant(variant), text_gio.Style(int(style)), text_gio.Weight(int(weight))}
}

func (ob *GoTextEditObj) SetFontBold(bold bool) {
	if bold {
		ob.font.Weight = text_gio.Bold
	} else {
		ob.font.Weight = text_gio.Light
	}
}

func (ob *GoTextEditObj) SetFontSize(size int) {
	ob.fontSize = unit_gio.Sp(size)
}

func (ob *GoTextEditObj) SetFontVariant(variant string) {
	ob.font.Variant = text_gio.Variant(variant)
}

func (ob *GoTextEditObj) SetFontWeight(weight GoFontWeight) {
	ob.font.Weight = text_gio.Weight(int(weight))
}

func (ob *GoTextEditObj) SelectedText() (text string) {
	return ob.editor.SelectedText()
}



/*func (ob *GoTextEditObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.SetSizePolicy(GetSizePolicy(horiz, vert))
}*/

func (ob *GoTextEditObj) SetText(text string) {
	ob.editor.SetText(text)
}

func (ob *GoTextEditObj) Text() (text string) {
	return ob.editor.Text()
}

func (ob *GoTextEditObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	//log.Println("GoTextEditObj::Draw()")
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

func (ob *GoTextEditObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	//log.Println("*GoTextEditObj::layout()")
	ob.ReceiveEvents(gtx)
	// *** create hint label macro
	macro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.hintColor.NRGBA()}.Add(gtx.Ops)
	var maxlines int
	if ob.editor.SingleLine {
		maxlines = 1
	}
	tl := widget_gio.Label{Alignment: ob.editor.Alignment, MaxLines: maxlines}
	dims := tl.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.hint)
	call := macro.Stop()
	// *** end
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = ob.editor.Layout(gtx, ob.shaper, ob.font, ob.fontSize, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.Editor.Add(gtx.Ops)
		//disabled := gtx.Queue == nil
		disabled := ob.HasFocus()
		//log.Println("disabled =", disabled)
		if ob.editor.Len() > 0 {
			paint_gio.ColorOp{Color: blendDisabledColor(disabled, ob.selectionColor.NRGBA())}.Add(gtx.Ops)
			//paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			ob.editor.PaintSelection(gtx)
			paint_gio.ColorOp{Color: blendDisabledColor(disabled, ob.color.NRGBA())}.Add(gtx.Ops)
			paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			ob.editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if ob.HasFocus() {
			paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			ob.editor.PaintCaret(gtx)
		}
		return dims //layout_gio.Dimensions{Size: gtx.Constraints.Min}
	})
	//defer clip_gio.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()
	defer clip_gio.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	// add the events handler to receive widget pointer events
	pointer_gio.CursorText.Add(gtx.Ops)
	ob.SignalEvents(gtx)
	return dims
}

func (ob *GoTextEditObj) ObjectType() (string) {
	return "GoTextEditObj"
}

func (ob *GoTextEditObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return DisabledBlend(c)
	}
	return c
}