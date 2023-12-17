// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/iconvg.go */

package utopia

import (
	"image"
	//"image/color"
	"image/draw"
	"log"
	"reflect"

	"github.com/utopiagio/utopia/internal/f32color"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	
	"golang.org/x/exp/shiny/iconvg"

	//archive "golang.org/x/exp/shiny/materialdesign/icons"	// eg: archive.FileFolder
)

const defaultIconColor = Color_Black
const defaultIconSize = 24


// example : folderIcon := GoIcon(parent, archive.FileFolder)

// Icon returns a new Icon from IconVG data.
func GoIconVG(data []byte, args ...interface{}) (*GoIconVGObj) {
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
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil
	}

	//object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	
	hIcon := &GoIconVGObj{
		//GioObject: object,
		//GioWidget: widget,
		src: data,
		color: color,
		description: text,
		size: size,
	}
	//parent.AddControl(hIcon)

	return hIcon
}

type GoIconVGObj struct {
	//GioObject
	//GioWidget

	src []byte
	color GoColor
	description string
	size int
	// Cached values.
	op       paint_gio.ImageOp
	imgSize  int
	imgColor GoColor	//color.NRGBA
}

/*func (ob *GoIconVGObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

func (ob *GoIconVGObj) ObjectType() (string) {
	return "GoIconVGObj"
}

/*func (ob *GoIconVGObj) Widget() (*GioWidget) {
	return nil
}*/

func (ob *GoIconVGObj) Size() (int) {
	return ob.size
}

// Layout displays the icon with its size set to the X minimum constraint.
func (ob *GoIconVGObj) Layout(gtx layout_gio.Context, color GoColor) layout_gio.Dimensions {
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
}

func (ob *GoIconVGObj) image(sz int, color GoColor) paint_gio.ImageOp {
	if sz == ob.imgSize && color == ob.imgColor {
		return ob.op
	}
	m, _ := iconvg.DecodeMetadata(ob.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = f32color.NRGBAToLinearRGBA(color.NRGBA())
	iconvg.Decode(&ico, ob.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ob.op = paint_gio.NewImageOp(img)
	ob.imgSize = sz
	ob.imgColor = color
	return ob.op
}