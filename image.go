// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"image"
	"log"
	"os"

	f32_gio "github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
)

func GoImage(parent GoObject, src string) (hObj *GoImageObj) {
	
	reader, err := os.Open(src)
	if err != nil {
	    log.Fatal(err)
	}
	defer reader.Close()
	//reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	imageOp := paint_gio.NewImageOp(m)

	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 100, 100, 100, 1000, 1000},
		FocusPolicy: NoFocus,
		Visible: true,
		//target: nil,
	}

	hImage := &GoImageObj{
		GioObject: object,
		GioWidget: widget,
		Src: imageOp,
		Fit: Contain,
		Position: layout_gio.Center,
		Scale: 0,	// defaults to 72 DPI
	}
	parent.AddControl(hImage)
	return hImage
}


// Image is a widget that displays an image.
type GoImageObj struct {
	GioObject
	GioWidget
	// Src is the image to display.
	Src paint_gio.ImageOp
	// Fit specifies how to scale the image to the constraints.
	// By default it does not do any scaling.
	Fit GoFit
	// Position specifies where to position the image within
	// the constraints.
	Position layout_gio.Direction
	// Scale is the ratio of image pixels to
	// dps. If Scale is zero Image falls back to
	// a scale that match a standard 72 DPI.
	Scale float32
}

func (ob *GoImageObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	/*X = gtx.Constraints.Max.X
	Y = gtx.Constraints.Max.Y
	if ob.SizePolicy.Horiz == FixedWidth {X = ob.X}
	if ob.SizePolicy.Vert == FixedHeight {Y = ob.Y}
	gtx.Constraints.Min = image.Point{X, Y}
	gtx.Constraints.Max = image.Point{X, Y}*/
	if ob.Visible {
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}

func (ob *GoImageObj) ObjectType() (string) {
	return "GoImageObj"
}

func (ob *GoImageObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

const defaultScale = float32(160.0 / 72.0)

func (ob *GoImageObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	scale := ob.Scale
	if scale == 0 {
		scale = defaultScale
	}

	size := ob.Src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Dp(unit_gio.Dp(wf*scale)), gtx.Dp(unit_gio.Dp(hf*scale))
	
	// paint object
	/*width := gtx.Dp(unit_gio.Dp(ob.Width))
	height := gtx.Dp(unit_gio.Dp(ob.Height))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	dims := image.Point{X: width, Y: height}*/

	dims, trans := ob.scale(gtx, ob.Position, layout_gio.Dimensions{Size: image.Pt(w, h)})
	
	defer clip_gio.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	//defer clip_gio.Ellipse{Max: dims.Size}.Push(gtx.Ops).Pop()

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32_gio.Affine2D{}.Scale(f32_gio.Point{}, f32_gio.Pt(pixelScale, pixelScale)))
	defer op_gio.Affine(trans).Push(gtx.Ops).Pop()

	ob.Src.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	return dims
}

// scale computes the new dimensions and transformation required to fit dims to cs, given the position.
func (ob *GoImageObj) scale(gtx layout_gio.Context, pos layout_gio.Direction, dims layout_gio.Dimensions) (layout_gio.Dimensions, f32_gio.Affine2D) {
	widgetSize := dims.Size

	if ob.Fit == Unscaled || dims.Size.X == 0 || dims.Size.Y == 0 {
		dims.Size = gtx.Constraints.Constrain(dims.Size)

		offset := pos.Position(widgetSize, dims.Size)
		dims.Baseline += offset.Y
		return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
	}


	width := gtx.Dp(unit_gio.Dp(ob.Width))
	height := gtx.Dp(unit_gio.Dp(ob.Height))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	scale := f32_gio.Point{
		X: float32(width) / float32(dims.Size.X),
		Y: float32(height) / float32(dims.Size.Y),
	}
	
	switch ob.Fit {
	case Contain:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case Cover:
		if scale.Y > scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case ScaleDown:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}

		// The widget would need to be scaled up, no change needed.
		if scale.X >= 1 {
			dims.Size = gtx.Constraints.Constrain(dims.Size)

			offset := pos.Position(widgetSize, dims.Size)
			dims.Baseline += offset.Y
			return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
		}
	case Fill:
	}

	var scaledSize image.Point
	scaledSize.X = int(float32(widgetSize.X) * scale.X)
	scaledSize.Y = int(float32(widgetSize.Y) * scale.Y)
	dims.Size = gtx.Constraints.Constrain(scaledSize)
	dims.Baseline = int(float32(dims.Baseline) * scale.Y)

	offset := pos.Position(scaledSize, dims.Size)
	trans := f32_gio.Affine2D{}.
		Scale(f32_gio.Point{}, scale).
		Offset(layout_gio.FPt(offset))

	dims.Baseline += offset.Y

	return dims, trans
}