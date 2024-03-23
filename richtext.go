// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/richtext.go */

package utopia

import (
	"log"
	"image"
	
	font_gio "github.com/utopiagio/gio/font"
	layout_gio "github.com/utopiagio/gio/layout"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	"github.com/utopiagio/utopia/metrics"

	"github.com/utopiagio/gio-x/markdown"
	"github.com/utopiagio/gio-x/richtext"
)

type GoRichTextObj struct {
	GioObject
	GioWidget
	name string
	title string
	anchorTable map[string]int
	font     font_gio.Font
	fontSize unit_gio.Sp
	
	spans []richtext.SpanStyle
	
	// SelectionColor is the color of the background for selected text.
	selectionColor GoColor
	selectionColorIndex int
	
	state richtext.InteractiveText
	shaper *text_gio.Shaper

	onLinkClick func(string)
}

func GoRichText(parent GoObject, name string) (hObj *GoRichTextObj) {
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
		name: name,
		anchorTable: make(map[string]int),
		selectionColor:	NRGBAColor(MulAlpha(theme.ContrastBg.NRGBA(), 0x60)),
		selectionColorIndex: 0,
		shaper: 	theme.Shaper,
	}
	parent.AddControl(hRichText)
	return hRichText
}

func (ob *GoRichTextObj) AnchorTable(ref string) (offset int) {
	return ob.anchorTable[ref]
}

func (ob *GoRichTextObj) AddContent(spans []richtext.SpanStyle) {
	ob.spans = append(ob.spans, spans...)
}

func (ob *GoRichTextObj) Clear() {
	ob.spans = []richtext.SpanStyle{}
}

/*func (ob *GoRichTextObj) Click(e pointer_gio.Event) {
	log.Println("GoRichTextObj::Click(e pointer_gio.Event)")
	if ob.onLinkClick != nil {
		ob.onLinkClick
	}
}*/

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
	/*for x := 0; x < len(spans_richtext); x++ {
		log.Printf("1 %f", spans_richtext[x].Content)
	}*/
	ob.AddContent(spans_richtext)
}

func (ob *GoRichTextObj) Layout(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	for {
		span, event, ok := ob.state.Update(gtx)
		if !ok {
			break
		}
		//content, _ := span.Content()
		switch event.Type {
		case richtext.Click:
			//log.Println(event.ClickData.Kind)
			if event.ClickData.Kind == 1 {	//gesture.KindClick {
				if url, ok := span.Get(markdown.MetadataURL).(string); ok && url != "" {
					if ob.onLinkClick != nil {
						ob.onLinkClick(url)
					}
				}
				ob.ParentWindow().Refresh()
			}
		}
	}
	ob.anchorTable, dims = richtext.Text(&ob.state, ob.shaper, ob.spans...).Layout(gtx)
	return dims
}

func (ob *GoRichTextObj) Name() (string) {
	return ob.name
}

func (ob *GoRichTextObj) Title() (string) {
	return ob.title
}

func (ob *GoRichTextObj) ObjectType() (string) {
	return "GoRichTextObj"
}

func (ob *GoRichTextObj) SetOnLinkClick(f func(string)) {
	ob.onLinkClick = f
}

func (ob *GoRichTextObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}