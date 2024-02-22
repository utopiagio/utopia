// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/richtext.go */

package utopia

import (
	"log"
	"image"
	
	"github.com/utopiagio/gio-x/markdown"
	"github.com/utopiagio/gio-x/richtext"
	//f32_gio "github.com/utopiagio/gio/f32"
	font_gio "github.com/utopiagio/gio/font"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	//pointer_gio "github.com/utopiagio/gio/io/pointer"
	//semantic_gio "github.com/utopiagio/gio/io/semantic"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_int "github.com/utopiagio/utopia/internal/widget"
	"github.com/utopiagio/utopia/metrics"
	
	//"golang.org/x/image/math/fixed"
)

type GoRichTextObj struct {
	GioObject
	GioWidget
	font     font_gio.Font
	fontSize unit_gio.Sp
	
	spans []richtext.SpanStyle
	
	// SelectionColor is the color of the background for selected text.
	selectionColor GoColor
	selectionColorIndex int
	
	state richtext.InteractiveText
	shaper *text_gio.Shaper

	//onFocus func()
}

func GoRichText(parent GoObject) (hObj *GoRichTextObj) {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{4,2,4,2},
		GoSize: GoSize{0, 0, 300, 26, 16777215, 16777215, 300, 26},
		FocusPolicy: StrongFocus,
		Visible: true,
	}
	hRichText := &GoRichTextObj{
		GioObject: 	object,
		GioWidget: 	widget,
		

		selectionColor:	NRGBAColor(MulAlpha(theme.ContrastBg.NRGBA(), 0x60)),
		selectionColorIndex: 0,
		shaper: 	theme.Shaper,
	}
	parent.AddControl(hRichText)
	return hRichText
}

func (ob *GoRichTextObj) AddContent(spans []richtext.SpanStyle) {
	ob.spans = append(ob.spans, spans...)
}

func (ob *GoRichTextObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
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

func (ob *GoRichTextObj) LoadMarkDown(src string) {
	buf := []byte(src)
	renderer := markdown.NewRenderer()
	renderer.Config.MonospaceFont.Typeface = "Go Mono"
	spans_richtext, err := renderer.Render(buf)
	if err != nil {
		log.Println("Render error..", err)
	}
	
	/*log.Println("len.6", len(spans_richtext[0].Content))
	for x := 0; x < len(spans_richtext[0].Content); x++ {
		log.Printf("1 %d %c", spans_richtext[0].Content[x], rune(spans_richtext[0].Content[x]))
	}*/
	ob.AddContent(spans_richtext)
}

func (ob *GoRichTextObj) Layout(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	return richtext.Text(&ob.state, ob.shaper, ob.spans...).Layout(gtx)
}

func (ob *GoRichTextObj) ObjectType() (string) {
	return "GoRichTextObj"
}

func (ob *GoRichTextObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}