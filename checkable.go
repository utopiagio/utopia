// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"image"
	//"image/color"

	//"github.com/utopiagio/gio/internal/f32color"
	layout_gio "github.com/utopiagio/gio/layout"
	clip_gio "github.com/utopiagio/gio/op/clip"
	font_gio "github.com/utopiagio/gio/font"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
	widget_int "github.com/utopiagio/utopia/internal/widget"
)

type checkable struct {
	label              string
	color              GoColor
	font               font_gio.Font
	fontSize           unit_gio.Sp
	iconColor          GoColor
	size               unit_gio.Dp
	shaper             *text_gio.Shaper
	checkedStateIcon   *widget_gio.Icon
	uncheckedStateIcon *widget_gio.Icon
}

func (ob *checkable) layout(gtx layout_gio.Context, checked, hovered bool) layout_gio.Dimensions {
	var icon *widget_gio.Icon
	textColorMacro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()

	if checked {
		icon = ob.checkedStateIcon
	} else {
		icon = ob.uncheckedStateIcon
	}

	dims := layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					size := gtx.Dp(ob.size) * 4 / 3
					dims := layout_gio.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !hovered {
						return dims
					}

					background := MulAlpha(ob.iconColor.NRGBA(), 70)

					b := image.Rectangle{Max: image.Pt(size, size)}
					paint_gio.FillShape(gtx.Ops, background, clip_gio.Ellipse(b).Op(gtx.Ops))

					return dims
				}),
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
						size := gtx.Dp(ob.size)
						col := ob.iconColor.NRGBA()
						if gtx.Queue == nil {
							col = DisabledBlend(col)
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
				paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
				return widget_int.GioLabel{}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.label, textColor)
			})
		}),
	)
	return dims
}