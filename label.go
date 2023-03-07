/* label.go */

package utopia

import (
	//"log"
	//"image/color"

	//gofont_gio "github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"

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
	theme := goApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
		//target: nil,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	theme.TextSize,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := ColorFromRGB(127, 0, 0)
	//hLabel.color = maroon
	//hLabel.font.Weight = text_gio.Bold 			// Thin - Medium - Bold
	//hLabel.font.Style = text_gio.Italic 		// Regular - Italic
	//hLabel.font.Variant = "Mono"	// Mono - Smallcaps
	parent.addControl(hLabel)
	return hLabel
}

type GoLabelObj struct {
	goObject
	goWidget
	alignment text_gio.Alignment
	font text_gio.Font
	fontSize unit_gio.Sp
	maxLines int
	text string
	color GoColor
	selectionColor GoColor
	//textAlign text.Alignment
	
	shaper *text_gio.Shaper
	state *widget_gio.Selectable
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

func (ob *GoLabelObj) SetSelectable(selectable bool) {
	if selectable {
		ob.state = new(widget_gio.Selectable)
	} else {
		ob.state = nil
	}
}

func (ob *GoLabelObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}

func (ob *GoLabelObj) SetText(text string) {
	ob.text = text
}

func (ob *GoLabelObj) SetTextAlignment(alignment GoTextAlignment) {
	ob.alignment = text_gio.Alignment(uint8(alignment))
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

func (ob *GoLabelObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

func (ob *GoLabelObj) layout(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	lbl := widget_gio.Label{Alignment: ob.alignment, MaxLines: ob.maxLines, Selectable: ob.state}
	if ob.state == nil {
		return lbl.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.text)
	}
	return lbl.LayoutSelectable(gtx, ob.shaper, ob.font, ob.fontSize, ob.text, func(gtx layout_gio.Context) layout_gio.Dimensions {
		paint_gio.ColorOp{Color: ob.selectionColor.NRGBA()}.Add(gtx.Ops)
		ob.state.PaintSelection(gtx)
		paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
		ob.state.PaintText(gtx)
		return layout_gio.Dimensions{}
	})
}

func (ob *GoLabelObj) objectType() (string) {
	return "GoLabelObj"
}

func H1Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
		//target: nil,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	32,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := ColorFromRGB(127, 0, 0)
	//hLabel.color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.addControl(hLabel)
	return hLabel
}

func H2Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	26,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.addControl(hLabel)
	return hLabel
}

func H3Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	22,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.addControl(hLabel)
	return hLabel
}

func H4Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	20,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.addControl(hLabel)
	return hLabel
}

func H5Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	16,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	parent.addControl(hLabel)
	return hLabel
}

func H6Label(parent GoObject, text string) (hObj *GoLabelObj) {
	theme := goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLabel := &GoLabelObj{
		goObject: object,
		goWidget: widget,
		fontSize: 	12,
		maxLines:	1,
		text: 		text,
		color: 		theme.ColorFg,
		selectionColor:	theme.ContrastBg,
		shaper: 	theme.Shaper,
		state: 		nil,
	}
	//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	//hH2Label.gio.Color = maroon
	hLabel.font.Weight = text_gio.Medium
	parent.addControl(hLabel)
	return hLabel
}