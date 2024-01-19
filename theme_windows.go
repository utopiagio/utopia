// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/theme_windows.go */

package utopia

import (
	"fmt"
	"image/color"
	"sync"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/font/gofont/gosmallcapsitalic"
	

	"github.com/utopiagio/gio/font/opentype"

	font_gio "github.com/utopiagio/gio/font"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"

	//"golang.org/x/image/font/gofont/goregular"
)

// Palette contains the minimal set of colors that a widget may need to
// draw itself.
type GoPalette struct {
	//BackColor 		GoColor 	// COLOR_WINDOW
	//BackgroundColor GoColor 	// COLOR_WINDOW
	//TextColor 		GoColor 	// COLOR_BTNTEXT Foreground color for button text.
	//FaceColor		GoColor  	// COLOR_BTNFACE Background color for button.
	//ForeColor 		GoColor 	// COLOR_WINDOWTEXT
	//GrayText 		GoColor 	// COLOR_GRAYTEXT Foreground color for disabled button text.
	//Highlight 		GoColor 	// COLOR_HIGHLIGHT Background color for slected button.
	//HighlightText 	GoColor 	// COLOR_HIGHLIGHTTEXT Foreground color for slected button text.
	//Hotlight 		GoColor 	// COLOR_HOTLIGHT Hyperlink color.

	// BackColor is the background color atop which content is currently being
	// drawn.
	ColorBg GoColor

	// ForeColor is a color suitable for drawing on top of Bg.
	ColorFg GoColor

	// ContrastBg is a color used to draw attention to active,
	// important, interactive widgets such as buttons.
	ContrastBg GoColor

	// ContrastFg is a color suitable for content drawn on top of
	// ContrastBg.
	ContrastFg GoColor

	// Provisional colors corresponding to Windows 10 Color Scheme
	TextColor GoColor

	BackColor GoColor

	FaceColor GoColor

	GrayText GoColor

	Highlight GoColor

	HighlightText GoColor

	Hotlight GoColor
}

var (
	regOnce    sync.Once
	reg        []font_gio.FontFace
	once       sync.Once
	collection []font_gio.FontFace
)

func GoTheme(fontCollection []text_gio.FontFace) *GoThemeObj {
	
	/*italicFace, _ := opentype.Parse(goitalic.TTF)
	boldFace, _ := opentype.Parse(gobold.TTF)
	bolditalicFace, _ := opentype.Parse(gobolditalic.TTF)
	mediumFace, _ := opentype.Parse(gomedium.TTF)
	mediumitalicFace, _ := opentype.Parse(gomediumitalic.TTF)
	ltrFace, _ := opentype.Parse(goregular.TTF)
	monoFace, _ := opentype.Parse(gomono.TTF)
	monoitalicFace, _ := opentype.Parse(gomonoitalic.TTF)
	monoboldFace, _ := opentype.Parse(gomonobold.TTF)
	monobolditalicFace, _ := opentype.Parse(gomonobolditalic.TTF)
	smallcapsFace, _ := opentype.Parse(gosmallcaps.TTF)
	smallcapsitalicFace, _ := opentype.Parse(gosmallcapsitalic.TTF)*/

	//collection := []text_gio.FontFace{{Face: italicFace}, {Face: boldFace}, {Face: bolditalicFace}, {Face: mediumFace}, {Face: mediumitalicFace}, {Face: ltrFace}}
	collection := Collection()
	th := &GoThemeObj{
		Shaper: text_gio.NewShaper(text_gio.NoSystemFonts(), text_gio.WithCollection(collection)),

	}
	th.GoPalette = GoPalette{
		ColorBg: 	GoColor(0xffffffff),
		ColorFg:	GoColor(0xff000000),
		ContrastBg: GoColor(0xff3f51b5),
		ContrastFg: GoColor(0xffffffff),
		TextColor:	GoColor(0xff000000),
		BackColor:	GoColor(0xffffffff),
		FaceColor:	GoColor(0xff000000),
		GrayText: 	GoColor(0xfff0f0f0),
		Highlight:	GoColor(0xff3f51b5),
		HighlightText: GoColor(0xffffffff),
		Hotlight: 	GoColor(0xff3f51b5),
	}
	th.TextSize = 14

	th.Icon.CheckBoxChecked = mustIcon(widget_gio.NewIcon(icons.ToggleCheckBox))
	th.Icon.CheckBoxUnchecked = mustIcon(widget_gio.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	th.Icon.RadioChecked = mustIcon(widget_gio.NewIcon(icons.ToggleRadioButtonChecked))
	th.Icon.RadioUnchecked = mustIcon(widget_gio.NewIcon(icons.ToggleRadioButtonUnchecked))

	// 38dp is on the lower end of possible finger size.
	th.FingerSize = 12
	// 6dp is the mimimum radius of slider finger
	th.ThumbRadius = 6
	// 2dp is the minimum size of slider track width
	th.TrackWidth = 4

	return th
}



type GoThemeObj struct {
	Shaper *text_gio.Shaper
	GoPalette
	TextSize unit_gio.Sp
	Icon     struct {
		CheckBoxChecked   *widget_gio.Icon
		CheckBoxUnchecked *widget_gio.Icon
		RadioChecked      *widget_gio.Icon
		RadioUnchecked    *widget_gio.Icon
	}

	// FingerSize is the minimum touch target size.
	FingerSize unit_gio.Dp
	// ThumbRadius is the mimimum radius of slider finger
	ThumbRadius unit_gio.Dp
	// TrackWidth is the minimum size of slider track width
	TrackWidth unit_gio.Dp 
}

// Regular returns a collection of all available Go font faces.
func Collection() []font_gio.FontFace {
	loadRegular()
	once.Do(func() {
		register(goitalic.TTF)
		register(gobold.TTF)
		register(gobolditalic.TTF)
		register(gomedium.TTF)
		register(gomediumitalic.TTF)
		register(gomono.TTF)
		register(gomonobold.TTF)
		register(gomonobolditalic.TTF)
		register(gomonoitalic.TTF)
		register(gosmallcaps.TTF)
		register(gosmallcapsitalic.TTF)
		// Ensure that any outside appends will not reuse the backing store.
		n := len(collection)
		collection = collection[:n:n]
	})
	return collection
}

func (ob *GoThemeObj) WithPalette(p GoPalette) *GoThemeObj {
	ob.GoPalette = p
	return ob
}

func loadRegular() {
	regOnce.Do(func() {
		faces, err := opentype.ParseCollection(goregular.TTF)
		if err != nil {
			panic(fmt.Errorf("failed to parse font: %v", err))
		}
		reg = faces
		collection = append(collection, reg[0])
	})
}

func register(ttf []byte) {
	faces, err := opentype.ParseCollection(ttf)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	collection = append(collection, faces[0])
}

func mustIcon(ic *widget_gio.Icon, err error) *widget_gio.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}