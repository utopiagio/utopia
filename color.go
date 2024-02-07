// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/color.go */

package utopia

import (
	//"log"
	"image/color"

)

type GoColor uint32

const (
	Transparent GoColor = 0x00000000
	Color_AliceBlue GoColor = 0xFFF0F8FF
	Color_AntiqueWhite GoColor = 0xFFFAEBD7
	Color_Aqua GoColor = 0xFF00FFFF
	Color_Aquamarine GoColor = 0xFF7FFFD4
	Color_Azure GoColor = 0xFFF0FFFF
	Color_Beige GoColor = 0xFFF5F5DC
	Color_Bisque GoColor = 0xFFFFE4C4
	Color_Black GoColor = 0xFF000000
	Color_BlanchedAlmond GoColor = 0xFFFFEBCD
	Color_Blue GoColor = 0xFF0000FF
	Color_BlueViolet GoColor = 0xFF8A2BE2
	Color_Brown GoColor = 0xFFA52A2A
	Color_BurlyWood GoColor = 0xFFDEB887
	Color_CadetBlue GoColor = 0xFF5F9EA0
	Color_Chartreuse GoColor = 0xFF7FFF00
	Color_Chocolate GoColor = 0xFFD2691E
	Color_Coral GoColor = 0xFFFF7F50
	Color_CornflowerBlue = 0xFF6495ED
	Color_Cornsilk = 0xFFFFF8DC
	Color_Crimson GoColor = 0xFFDC143C
	Color_Cyan GoColor = 0xFF00FFFF
	Color_DarkBlue GoColor = 0xFF00008B
	Color_DarkCyan GoColor = 0xFF008B8B
	Color_DarkGoldenrod GoColor = 0xFFB8860B
	Color_DarkGray GoColor = 0xFFA9A9A9
	Color_DarkGrey GoColor = 0xFFA9A9A9
	Color_DarkGreen GoColor = 0xFF006400
	Color_DarkKhaki GoColor = 0xFFBDB76B
	Color_DarkMagenta GoColor = 0xFF8B008B
	Color_DarkOliveGreen GoColor = 0xFF556B2F
	Color_DarkOrange GoColor = 0xFFFF8C00
	Color_DarkOrchid GoColor = 0xFF9932CC
	Color_DarkRed GoColor = 0xFF8B0000
	Color_DarkSalmon GoColor = 0xFFE9967A
	Color_DarkSeaGreen GoColor = 0xFF8FBC8F
	Color_DarkSlateBlue GoColor = 0xFF483D8B
	Color_DarkSlateGray GoColor = 0xFF2F4F4F
	Color_DarkTurquoise GoColor = 0xFF00CED1
	Color_DarkViolet GoColor = 0xFF9400D3
	Color_DeepPink GoColor = 0xFFFF1493
	Color_DeepSkyBlue GoColor = 0xFF00BFFF
	Color_DimGray GoColor = 0xFF696969
	Color_DodgerBlue GoColor = 0xFF1E90FF
	Color_Firebrick GoColor = 0xFFB22222
	Color_FloralWhite GoColor = 0xFFFFFAF0
	Color_ForestGreen GoColor = 0xFF228B22
	Color_Fuchsia GoColor = 0xFFFF00FF
	Color_Gainsboro GoColor = 0xFFDCDCDC
	Color_GhostWhite GoColor = 0xFFF8F8FF
	Color_Gold GoColor = 0xFFFFD700
	Color_Goldenrod GoColor = 0xFFDAA520
	Color_Gray GoColor = 0xFF808080
	Color_Grey GoColor = 0xFF808080
	Color_Green GoColor = 0xFF008000
	Color_GreenYellow GoColor = 0xFFADFF2F
	Color_Honeydew GoColor = 0xFFF0FFF0
	Color_HotPink GoColor = 0xFFFF69B4
	Color_IndianRed GoColor = 0xFFCD5C5C
	Color_Indigo GoColor = 0xFF4B0082
	Color_Ivory GoColor = 0xFFFFFFF0
	Color_Khaki GoColor = 0xFFF0E68C
	Color_Lavender GoColor = 0xFFE6E6FA
	Color_LavenderBlush GoColor = 0xFFFFF0F5
	Color_LawnGreen GoColor = 0xFF7CFC00
	Color_LemonChiffon GoColor = 0xFFFFFACD
	Color_LightBlue GoColor = 0xFFADD8E6
	Color_LightCoral GoColor = 0xFFF08080
	Color_LightCyan GoColor = 0xFFE0FFFF
	Color_LightGoldenrodYellow GoColor = 0xFFFAFAD2
	Color_LightGray GoColor = 0xFFD3D3D3
	Color_LightGrey GoColor = 0xFFD3D3D3
	Color_LightGreen GoColor = 0xFF90EE90
	Color_LightPink GoColor = 0xFFFFB6C1
	Color_LightSalmon GoColor = 0xFFFFA07A
	Color_LightSeaGreen GoColor = 0xFF20B2AA
	Color_LightSkyBlue GoColor = 0xFF87CEFA
	Color_LightSlateGray GoColor = 0xFF778899
	Color_LightSteelBlue GoColor = 0xFFB0C4DE
	Color_Lime GoColor = 0xFF00FF00
	Color_LimeGreen GoColor = 0xFF32CD32
	Color_Linen GoColor = 0xFFFAF0E6
	Color_Magenta GoColor = 0xFFFF00FF
	Color_Maroon GoColor = 0xFF800000
	Color_MediumAquamarine GoColor = 0xFF66CDAA
	Color_MediumBlue GoColor = 0xFF0000CD
	Color_MediumOrchid GoColor = 0xFFBA55D3
	Color_MediumPurple GoColor = 0xFF9370DB
	Color_MediumSeaGreen GoColor = 0xFF3CB371
	Color_MediumSlateBlue GoColor = 0xFF7B68EE
	Color_MediumSpringGreen GoColor = 0xFF00FA9A
	Color_MediumTurquoise GoColor = 0xFF48D1CC
	Color_MediumVioletRed GoColor = 0xFFC71585
	Color_MidnightBlue GoColor = 0xFF191970
	Color_MintCream GoColor = 0xFFF5FFFA
	Color_MistyRose GoColor = 0xFFFFE4E1
	Color_Moccasin GoColor = 0xFFFFE4B5
	Color_NavajoWhite GoColor = 0xFFFFDEAD
	Color_Navy GoColor = 0xFF000080
	Color_OldLace GoColor = 0xFFFDF5E6
	Color_Olive GoColor = 0xFF808000
	Color_OliveDrab GoColor = 0xFF6B8E23
	Color_Orange GoColor = 0xFFFFA500
	Color_OrangeRed GoColor = 0xFFFF4500
	Color_Orchid GoColor = 0xFFDA70D6
	Color_PaleGoldenrod GoColor = 0xFFEEE8AA
	Color_PaleGreen GoColor = 0xFF98FB98
	Color_PaleTurquoise GoColor = 0xFFAFEEEE
	Color_PaleVioletRed GoColor = 0xFFDB7093
	Color_PapayaWhip GoColor = 0xFFFFEFD5
	Color_PeachPuff GoColor = 0xFFFFDAB9
	Color_Peru GoColor = 0xFFCD853F
	Color_Pink GoColor = 0xFFFFC0CB
	Color_Plum GoColor = 0xFFDDA0DD
	Color_PowderBlue GoColor = 0xFFB0E0E6
	Color_Purple GoColor = 0xFF800080
	Color_Red GoColor = 0xFFFF0000
	Color_RosyBrown GoColor = 0xFFBC8F8F
	Color_RoyalBlue GoColor = 0xFF4169E1
	Color_SaddleBrown GoColor = 0xFF8B4513
	Color_Salmon GoColor = 0xFFFA8072
	Color_SandyBrown GoColor = 0xFFF4A460
	Color_SeaGreen GoColor = 0xFF2E8B57
	Color_SeaShell GoColor = 0xFFFFF5EE
	Color_Sienna GoColor = 0xFFA0522D
	Color_Silver GoColor = 0xFFC0C0C0
	Color_SkyBlue GoColor = 0xFF87CEEB
	Color_SlateBlue GoColor = 0xFF6A5ACD
	Color_SlateGray GoColor = 0xFF708090
	Color_Snow GoColor = 0xFFFFFAFA
	Color_SpringGreen GoColor = 0xFF00FF7F
	Color_SteelBlue GoColor = 0xFF4682B4
	Color_Tan GoColor = 0xFFD2B48C
	Color_Teal GoColor = 0xFF008080
	Color_Thistle GoColor = 0xFFD8BFD8
	Color_Tomato GoColor = 0xFFFF6347
	Color_Transparent GoColor = 0x00FFFFFF
	Color_Turquoise GoColor = 0xFF40E0D0
	Color_Violet GoColor = 0xFFEE82EE
	Color_Wheat GoColor = 0xFFF5DEB3
	Color_White GoColor = 0xFFFFFFFF
	Color_WhiteSmoke GoColor = 0xFFF5F5F5
	Color_Yellow GoColor = 0xFFFFFF00
	Color_YellowGreen GoColor = 0xFF9ACD32
)

func ColorIndex(idx int) GoColor {
	switch idx {
		case 0	:
		 	return GoColor(0x00000000)
		case 1 	:
			return GoColor(0xFFF0F8FF)
		case 2 	:
			return GoColor(0xFFFAEBD7)
		case 3 	:
			return GoColor(0xFF00FFFF)
		case 4 	:
			return GoColor(0xFF7FFFD4)
		case 5 	:
			return GoColor(0xFFF0FFFF)
		case 6 	:
			return GoColor(0xFFF5F5DC)
		case 7 	:
			return GoColor(0xFFFFE4C4)

	}
	return GoColor(0xFF000000)
	/*8 	:	GoColor = 0xFF000000
	9 	:	GoColor = 0xFFFFEBCD
	10 	:	GoColor = 0xFF0000FF
	11 	:	GoColor = 0xFF8A2BE2
	12 	:	GoColor = 0xFFA52A2A
	13 	:	GoColor = 0xFFDEB887
	14 	:	GoColor = 0xFF5F9EA0
	15 	:	GoColor = 0xFF7FFF00
	16 	:	GoColor = 0xFFD2691E
	17 	:	GoColor = 0xFFFF7F50
	18 	:	GoColor = 0xFF6495ED
	19 	:	GoColor = 0xFFFFF8DC
	20 	:	GoColor = 0xFFDC143C
	21 	:	GoColor = 0xFF00FFFF
	22 	:	GoColor = 0xFF00008B
	23 	:	GoColor = 0xFF008B8B*/

}


func ColorFromRGB(r uint8, g uint8, b uint8, a uint8) (GoColor) {
	color := uint32(uint32(a) + uint32(r) << 16 + uint32(g) << 8 + uint32(b))
	return GoColor(color)
}

func ColorToRGB(color GoColor) (r uint8, g uint8, b uint8) {
	return uint8((color & 0xFF0000) >> 16), uint8((color & 0xFF00) >> 8), uint8(color & 0xFF)
}

/*func colorFromW32(ref w32.COLORREF) (GoColor) {
	color := uint32(0xFF000000 + (ref & 0xFF) << 16 + (ref & 0xFF00) + (ref & 0xFF0000) >> 16)
	return GoColor(color)
}

func colorToW32(col GoColor) (w32.COLORREF) {
	w32Color := uint32((col & 0xFF) << 16 + (col & 0xFF00) + (col & 0xFF0000) >> 16)
	return w32.COLORREF(w32Color)
}*/

/*func colorToImgColor(col GoColor) (color.NRGBA) {
	return color.NRGBA{R: uint8((col & 0xFF0000)) >> 16, G: uint8((col & 0xFF00)) >> 8, B: uint8(col & 0xFF), A: uint8((col & 0xFF000000) >> 24)}
}*/

func RGBAColor(r uint8, g uint8, b uint8, a uint8) (GoColor) {
	color := uint32(uint32(a) << 24 + uint32(r) << 16 + uint32(g) << 8 + uint32(b))
	return GoColor(color)
}

func (c GoColor) RGBA() (r uint8, g uint8, b uint8, a uint8) {
	return uint8((c & 0xFF0000) >> 16), uint8((c & 0xFF00) >> 8), uint8(c & 0xFF), uint8((c & 0xFF000000) >> 24)
}

/*func (c GoColor) LinearRGBA() (color.RGBA) {
	colr := color.NRGBA{R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c), A: uint8(c >> 24)}
	if colr.A == 0xFF {
		return color.RGBA(colr)
	}
	cl := LinearFromSRGB(col)
	return color.RGBA{
		R: uint8(cl.R*255 + .5),
		G: uint8(cl.G*255 + .5),
		B: uint8(cl.B*255 + .5),
		A: col.A,
	}
}*/

// LinearFromSRGB converts from col in the sRGB colorspace to RGBA.
/*func LinearFromSRGB(col color.NRGBA) RGBA {
	af := float32(col.A) / 0xFF
	return RGBA{
		R: srgb8ToLinear[col.R] * af, // sRGBToLinear(float32(col.R)/0xff) * af,
		G: srgb8ToLinear[col.G] * af, // sRGBToLinear(float32(col.G)/0xff) * af,
		B: srgb8ToLinear[col.B] * af, // sRGBToLinear(float32(col.B)/0xff) * af,
		A: af,
	}
}*/

func NRGBAColor(col color.NRGBA) (GoColor) {
	color := uint32(uint32(col.A) << 24 + uint32(col.R) << 16 + uint32(col.G) << 8 + uint32(col.B))
	return GoColor(color)
}

// func NRGBA returns NRGBA color from GoColor
func (c GoColor) NRGBA() (color.NRGBA) {
	return color.NRGBA{R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c), A: uint8(c >> 24)}
}

// MulAlpha applies the alpha to the color.
func (c GoColor) MulAlpha(alpha uint8) (color GoColor) {
	col := c.NRGBA()
	col.A = uint8(uint32(col.A) * uint32(alpha) / 0xFF)
	return NRGBAColor(col)
}

// MulAlpha applies the alpha to the color.
func MulAlpha(col color.NRGBA, alpha uint8) color.NRGBA {
	col.A = uint8(uint32(col.A) * uint32(alpha) / 0xFF)
	return col
}

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func DisabledBlend(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	d = mix(c, color.NRGBA{A: c.A, R: lum, G: lum, B: lum}, r)
	d = MulAlpha(d, 128+32)
	return
}

// Hovered blends dark colors towards white, and light colors towards
// black. It is approximate because it operates in non-linear sRGB space.
func HoveredBlend(c color.NRGBA) (h color.NRGBA) {
	if c.A == 0 {
		// Provide a reasonable default for transparent widgets.
		return color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
	}
	const ratio = 0x20
	m := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: c.A}
	if approxLuminance(c) > 128 {
		m = color.NRGBA{A: c.A}
	}
	return mix(m, c, ratio)
}

// mix mixes c1 and c2 weighted by (1 - a/256) and a/256 respectively.
func mix(c1, c2 color.NRGBA, a uint8) color.NRGBA {
	ai := int(a)
	return color.NRGBA{
		R: byte((int(c1.R)*ai + int(c2.R)*(256-ai)) / 256),
		G: byte((int(c1.G)*ai + int(c2.G)*(256-ai)) / 256),
		B: byte((int(c1.B)*ai + int(c2.B)*(256-ai)) / 256),
		A: byte((int(c1.A)*ai + int(c2.A)*(256-ai)) / 256),
	}
}

// approxLuminance is a fast approximate version of RGBA.Luminance.
func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}