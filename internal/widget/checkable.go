// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/internal/widget/checkable.go */

package widget

import (
	"image"
	"image/color"
	//"log"

	"github.com/utopiagio/utopia/internal/f32color"
	layout_gio "github.com/utopiagio/gio/layout"
	clip_gio "github.com/utopiagio/gio/op/clip"
	font_gio "github.com/utopiagio/gio/font"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
	//widget_int "github.com/utopiagio/utopia/internal/widget"
)

type GioCheckable struct {
	Label              string
	Color              color.NRGBA
	Font               font_gio.Font
	TextSize           unit_gio.Sp
	IconColor          color.NRGBA
	Size               unit_gio.Dp
	Shaper             *text_gio.Shaper
	CheckedStateIcon   *widget_gio.Icon
	UncheckedStateIcon *widget_gio.Icon
}

func (c *GioCheckable) Layout(gtx layout_gio.Context, checked, hovered bool) layout_gio.Dimensions {
	var icon *widget_gio.Icon
	if checked {
		icon = c.CheckedStateIcon
	} else {
		icon = c.UncheckedStateIcon
	}

	dims := layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					size := gtx.Dp(c.Size) * 4 / 3
					dims := layout_gio.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !hovered {
						return dims
					}

					background := f32color.MulAlpha(c.IconColor, 70)

					b := image.Rectangle{Max: image.Pt(size, size)}
					paint_gio.FillShape(gtx.Ops, background, clip_gio.Ellipse(b).Op(gtx.Ops))

					return dims
				}),
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
						size := gtx.Dp(c.Size)
						col := c.IconColor
						if !gtx.Enabled() {
							col = f32color.Disabled(col)
						}
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, col)
						return layout_gio.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
			)
		}),

		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				colMacro := op_gio.Record(gtx.Ops)
				paint_gio.ColorOp{Color: c.Color}.Add(gtx.Ops)
				return widget_gio.Label{}.Layout(gtx, c.Shaper, c.Font, c.TextSize, c.Label, colMacro.Stop())
			})
		}),
	)
	return dims
	//log.Println("(c *GioCheckable) Layout dims: (", dims.Size.X, lbldims.Size.Y, ")")
	/*return layout_gio.Dimensions{
							Size: image.Point{X: dims.Size.X, Y: dims.Size.Y},
						}*/
}
