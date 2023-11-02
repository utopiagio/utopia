package utopia

import (
	"image"
	"image/color"

	layout_gio "github.com/utopiagio/gio/layout"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
)

type GoRect struct {
	Color color.NRGBA
	Size  image.Point
	Radii int
}

func (r GoRect) Layout(gtx C) D {
	paint_gio.FillShape(
		gtx.Ops,
		r.Color,
		clip_gio.UniformRRect(
			image.Rectangle{
				Max: r.Size,
			},
			r.Radii,
		).Op(gtx.Ops))
	return layout_gio.Dimensions{Size: r.Size}
}

func PaintRect(gtx layout_gio.Context, size image.Point, fill GoColor) {
	GoRect{
		Color: fill.NRGBA(),
		Size:  size,
	}.Layout(gtx)
}