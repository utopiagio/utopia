package utopia

import (
	"image"
	//"image/color"
	//"image/draw"
	"log"
	"os"
	"reflect"

	//f32_ui "github.com/utopiagio/utopia/colorf32"
	f32_gio "github.com/utopiagio/gio/f32"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	op_gio "github.com/utopiagio/gio/op"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	
	

	
)
// declared in iconvg.go
//const defaultIconColor = Color_Black
//const defaultIconSize = 24


// example : folderIcon := GoIcon(parent, archive.FileFolder)

// Icon returns a new Icon from IconVG data.
func GoIconPNG(filePath string, args ...interface{}) (*GoIconPNGObj) {
	var color GoColor
	var size int
	var text string
	color = defaultIconColor
	size = defaultIconSize
	text = ""
	for i, v := range args {
		//log.Println("GoIcon() - arg:", i, "value:", v)
		switch v := reflect.ValueOf(v); v.Kind() {
			case reflect.String:
				//log.Println("GoIcon() - v.String():", v.String())
				text = v.String()
			case reflect.Int:
				//log.Println("GoIcon() - v.Int():", v.Int())
				size = args[i].(int)
			case reflect.Uint32:
				//log.Println("GoIcon() - v.Uint32():", args[i].(GoColor))
				color = args[i].(GoColor)
			default:
				log.Println("GoIcon() - Not String or GoColor")
		}
	}
	reader, err := os.Open(filePath)
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
	
	//object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	
	hIcon := &GoIconPNGObj{
		//GioObject: object,
		//GioWidget: widget,
		src: filePath,
		color: color,
		description: text,
		fit: Contain,
		positioning: layout_gio.Center,
		size: size,
		imgOp: imageOp,
		imgSize: size,
	}
	//parent.AddControl(hIcon)

	return hIcon
}

type GoIconPNGObj struct {
	//GioObject
	//GioWidget

	src string
	color GoColor
	description string
	fit GoFit
	positioning layout_gio.Direction
	size int
	// Cached values.
	imgOp       paint_gio.ImageOp
	imgSize  int
	imgColor GoColor	//color.NRGBA
}

/*func (ob *GoIconPNGObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	
	if ob.Visible {
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx C, ) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx, ob.color)
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}*/

func (ob *GoIconPNGObj) ObjectType() (string) {
	return "GoIconPNGObj"
}

/*func (ob *GoIconPNGObj) Widget() (*GioWidget) {
	return nil
}*/

func (ob *GoIconPNGObj) Size() (int) {
	return ob.size
}

// Layout displays the icon with its size set to the X minimum constraint.
/*func (ob *GoIconPNGObj) Layout(gtx layout_gio.Context, color GoColor) layout_gio.Dimensions {
	//ob.ReceiveEvents(gtx)
	
	sz := gtx.Constraints.Min.X
	//log.Println("Icon ob.size =", sz)
	if sz == 0 {
		sz = gtx.Dp(unit_gio.Dp(ob.size))
	}
	//log.Println("Icon sz =", sz)
	size := gtx.Constraints.Constrain(image.Pt(sz, sz))
	defer clip_gio.Rect{Max: size}.Push(gtx.Ops).Pop()

	ico := ob.image(size.X, color)
	ico.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	// add the events handler to receive widget pointer events
	//ob.SignalEvents(gtx)

	return layout_gio.Dimensions{
		Size: ico.Size(),
	}
}*/

const defaultScale = float32(160.0 / 72.0)

func (ob *GoIconPNGObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	scale := defaultScale
	/*scale := ob.scale
	if scale == 0 {
		scale = defaultScale
	}*/
	//sz := gtx.Constraints.Min.X
	//log.Println("Icon ob.size =", sz)
	//if sz == 0 {
		//sz = gtx.Dp(unit_gio.Dp(ob.size))
	//}
	size := ob.imgOp.Size()

	//size := gtx.Constraints.Constrain(image.Pt(sz, sz))
	//defer clip_gio.Rect{Max: size}.Push(gtx.Ops).Pop()

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

	dims, trans := ob.scale(gtx, ob.positioning, layout_gio.Dimensions{Size: image.Pt(w, h)})
	
	defer clip_gio.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	//defer clip_gio.Ellipse{Max: dims.Size}.Push(gtx.Ops).Pop()

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32_gio.Affine2D{}.Scale(f32_gio.Point{}, f32_gio.Pt(pixelScale, pixelScale)))
	defer op_gio.Affine(trans).Push(gtx.Ops).Pop()

	ob.imgOp.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	return dims
}

/*func (ob *GoIconPNGObj) image(sz int, color GoColor) paint_gio.ImageOp {
	if sz == ob.imgSize && color == ob.imgColor {
		return ob.op
	}
	m, _ := iconvg.DecodeMetadata(ob.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = f32_gio.NRGBAToLinearRGBA(color.NRGBA())
	iconvg.Decode(&ico, ob.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ob.op = paint_gio.NewImageOp(img)
	ob.imgSize = sz
	ob.imgColor = color
	return ob.op
}*/

// scale computes the new dimensions and transformation required to fit dims to cs, given the position.
func (ob *GoIconPNGObj) scale(gtx layout_gio.Context, positioning layout_gio.Direction, dims layout_gio.Dimensions) (layout_gio.Dimensions, f32_gio.Affine2D) {
	widgetSize := dims.Size

	if ob.fit == Unscaled || dims.Size.X == 0 || dims.Size.Y == 0 {
		dims.Size = gtx.Constraints.Constrain(dims.Size)

		offset := positioning.Position(widgetSize, dims.Size)
		dims.Baseline += offset.Y
		return dims, f32_gio.Affine2D{}.Offset(layout_gio.FPt(offset))
	}

	//sz = gtx.Dp(unit_gio.Dp(ob.size))
	width := gtx.Dp(unit_gio.Dp(ob.size))
	height := gtx.Dp(unit_gio.Dp(ob.size))
	
	scale := f32_gio.Point{
		X: float32(width) / float32(dims.Size.X),
		Y: float32(height) / float32(dims.Size.Y),
	}
	
	switch ob.fit {
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

			offset := positioning.Position(widgetSize, dims.Size)
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

	offset := positioning.Position(scaledSize, dims.Size)
	trans := f32_gio.Affine2D{}.
		Scale(f32_gio.Point{}, scale).
		Offset(layout_gio.FPt(offset))

	dims.Baseline += offset.Y

	return dims, trans
}