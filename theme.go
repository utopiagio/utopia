/* theme */

package utopia

import (
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/icons"

	//"github.com/utopiagio/gio/font/gofont"
	//"github.com/utopiagio/gio/font/opentype"

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

func GoTheme(fontCollection []text_gio.FontFace) *GoThemeObj {
	//ltrFace, _ := opentype.Parse(goregular.TTF)
	//collection := text_gio[]FontFace{{Face: ltrFace}}
	th := &GoThemeObj{
		Shaper: text_gio.NewShaper(text_gio.NoSystemFonts(), text_gio.WithCollection(fontCollection)),

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
	th.FingerSize = 38

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
}

func (ob *GoThemeObj) WithPalette(p GoPalette) *GoThemeObj {
	ob.GoPalette = p
	return ob
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