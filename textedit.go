// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"image/color"

	//"github.com/utopiagio/gio/internal/f32color"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoTextEditObj struct {
	goObject
	goWidget
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
	editor         *widget_gio.Editor

	shaper *text_gio.Shaper

	onFocus func()
}

func GoTextEdit(parent GoObject, hintText string) *GoTextEditObj {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
		//target: nil,
	}
	hTextEdit := &GoTextEditObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	theme.TextSize,
		color: 		theme.ColorFg,
		hint: 		hintText,
		hintColor: 		NRGBAColor(MulAlpha(theme.ColorFg.NRGBA(), 0xbb)),
		selectionColor:	NRGBAColor(MulAlpha(theme.ContrastBg.NRGBA(), 0x60)),
		editor: 	new(widget_gio.Editor),
		shaper: 	theme.Shaper,
	}
	parent.addControl(hTextEdit)
	return hTextEdit
}

func (ob *GoTextEditObj) ClearSelection() {
	ob.editor.ClearSelection()
}

func (ob *GoTextEditObj) Focused() bool {
	return ob.editor.Focused()
}

func (ob *GoTextEditObj) Insert(text string) {
	ob.editor.Insert(text)
}

func (ob *GoTextEditObj) Length() (length int) {
	return ob.editor.Len()
}

func (ob *GoTextEditObj) SelectedText() (text string) {
	return ob.editor.SelectedText()
}

func (ob *GoTextEditObj) SetOnFocus(f func()) {
	ob.onFocus = f
}

func (ob *GoTextEditObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}

func (ob *GoTextEditObj) SetText(text string) {
	ob.editor.SetText(text)
}

func (ob *GoTextEditObj) Text() (text string) {
	return ob.editor.Text()
}

func (ob *GoTextEditObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.visible {
		dims = ob.goMargin.layout(gtx, func(gtx C) D {
			return ob.goBorder.layout(gtx, func(gtx C) D {
				return ob.goPadding.layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
	}
	return dims
}

func (ob *GoTextEditObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	macro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.hintColor.NRGBA()}.Add(gtx.Ops)
	var maxlines int
	if ob.editor.SingleLine {
		maxlines = 1
	}
	tl := widget_gio.Label{Alignment: ob.editor.Alignment, MaxLines: maxlines}
	dims := tl.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.hint)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = ob.editor.Layout(gtx, ob.shaper, ob.font, ob.fontSize, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.Editor.Add(gtx.Ops)
		disabled := gtx.Queue == nil
		if ob.editor.Len() > 0 {
			paint_gio.ColorOp{Color: blendDisabledColor(disabled, ob.selectionColor.NRGBA())}.Add(gtx.Ops)
			ob.editor.PaintSelection(gtx)
			paint_gio.ColorOp{Color: blendDisabledColor(disabled, ob.color.NRGBA())}.Add(gtx.Ops)
			ob.editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled {
			paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			ob.editor.PaintCaret(gtx)
		}
		return dims
	})
	return dims
}

func (ob *GoTextEditObj) objectType() (string) {
	return "GoTextEditObj"
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return DisabledBlend(c)
	}
	return c
}