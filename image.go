// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"image"

	f32_gio "github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

// Image is a widget that displays an image.
type GoImageObj struct {
	// Src is the image to display.
	Src paint_gio.ImageOp
	// Fit specifies how to scale the image to the constraints.
	// By default it does not do any scaling.
	fit widget_gio.Fit
	// Position specifies where to position the image within
	// the constraints.
	Position layout_gio.Direction
	// Scale is the ratio of image pixels to
	// dps. If Scale is zero Image falls back to
	// a scale that match a standard 72 DPI.
	Scale float32
}

const defaultScale = float32(160.0 / 72.0)

func (im *GoImageObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	scale := im.Scale
	if scale == 0 {
		scale = defaultScale
	}

	size := im.Src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Dp(unit_gio.Dp(wf*scale)), gtx.Dp(unit_gio.Dp(hf*scale))

	dims, trans := im.scale(gtx.Constraints, im.Position, layout_gio.Dimensions{Size: image.Pt(w, h)})
	defer clip_gio.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32_gio.Affine2D{}.Scale(f32_gio.Point{}, f32_gio.Pt(pixelScale, pixelScale)))
	defer op_gio.Affine(trans).Push(gtx.Ops).Pop()

	im.Src.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	return dims
}

// scale computes the new dimensions and transformation required to fit dims to cs, given the position.
func (im *GoImageObj) scale(cs layout_gio.Constraints, pos layout_gio.Direction, dims layout_gio.Dimensions) (layout_gio.Dimensions, f32_gio.Affine2D) {
	widgetSize := dims.Size

	if im.fit == widget_gio.Unscaled || dims.Size.X == 0 || dims.Size.Y == 0 {
		dims.Size = cs.Constrain(dims.Size)

		offset := pos.Position(widgetSize, dims.Size)
		dims.Baseline += offset.Y
		return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
	}

	scale := f32_gio.Point{
		X: float32(cs.Max.X) / float32(dims.Size.X),
		Y: float32(cs.Max.Y) / float32(dims.Size.Y),
	}

	switch im.fit {
	case widget_gio.Contain:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case widget_gio.Cover:
		if scale.Y > scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case widget_gio.ScaleDown:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}

		// The widget would need to be scaled up, no change needed.
		if scale.X >= 1 {
			dims.Size = cs.Constrain(dims.Size)

			offset := pos.Position(widgetSize, dims.Size)
			dims.Baseline += offset.Y
			return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
		}
	case widget_gio.Fill:
	}

	var scaledSize image.Point
	scaledSize.X = int(float32(widgetSize.X) * scale.X)
	scaledSize.Y = int(float32(widgetSize.Y) * scale.Y)
	dims.Size = cs.Constrain(scaledSize)
	dims.Baseline = int(float32(dims.Baseline) * scale.Y)

	offset := pos.Position(scaledSize, dims.Size)
	trans := f32_gio.Affine2D{}.
		Scale(f32_gio.Point{}, scale).
		Offset(layout_gio.FPt(offset))

	dims.Baseline += offset.Y

	return dims, trans
}