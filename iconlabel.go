// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/iconlabel.go */

package utopia

import (
	"image"
	//"image/color"
	"image/draw"
	"log"
	"reflect"

	"github.com/utopiagio/utopia/internal/f32color"
	"github.com/utopiagio/utopia/metrics"
	font_gio "github.com/utopiagio/gio/font"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
	widget_int "github.com/utopiagio/utopia/internal/widget"
	
	"golang.org/x/exp/shiny/iconvg"

	//archive "golang.org/x/exp/shiny/materialdesign/icons"	// eg: archive.FileFolder
)

//const defaultIconColor = Color_Black
//const defaultIconSize = 24

// example : folderIcon := GoIcon(parent, archive.FileFolder)

// GoIcon returns a new Icon from IconVG data.
func GoIconLabel(parent GoObject, data []byte, args ...interface{}) (*GoIconLabelObj) {
	var color GoColor
	var size int
	var text string
	var theme *GoThemeObj = GoApp.Theme()
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

	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 24, 100, 1000, 1000, 24, 100},
		FocusPolicy: StrongFocus,
		Visible: true,

		ForeColor: theme.ColorFg,
		
		//target: nil,
	}
	
	hIcon := &GoIconLabelObj{
		GioObject: object,
		GioWidget: widget,
		color: theme.ColorFg,
		fontSize: theme.TextSize,
		icon: data,
		iconColor: color,
		label: text,
		size: size,
		shaper: theme.Shaper,
	}
	parent.AddControl(hIcon)

	return hIcon
}

type GoIconLabelObj struct {
	GioObject
	GioWidget
	color GoColor
	font               font_gio.Font
	fontSize           unit_gio.Sp
	icon []byte
	iconColor GoColor
	label string
	size int
	shaper             *text_gio.Shaper
	// Cached values.
	op       paint_gio.ImageOp
	imgSize  int
	imgColor GoColor	//color.NRGBA
}

func (ob *GoIconLabelObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	cs := gtx.Constraints
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		w := min(maxWidth, width)			// constrain to ob.MaxWidth
		cs.Min.X = max(minWidth, w)				// constrain to ob.MinWidth 
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = minWidth						// set to ob.MinWidth
		cs.Max.X = minWidth						// set to ob.MinWidth
	case PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = minWidth						// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = maxWidth						// set to ob.MaxWidth
		cs.Min.X = maxWidth						// set to ob.MaxWidth
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		w := min(maxHeight, height)				// constrain to ob.MaxHeight
		cs.Min.Y = max(minHeight, w)			// constrain to ob.MinHeight 
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = minHeight					// set to ob.MinHeight
		cs.Max.Y = minHeight					// set to ob.MinHeight
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = max(0, minHeight)			// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = maxHeight					// set to ob.MaxHeight
		cs.Min.Y = maxHeight					// set to ob.MaxHeight
	case ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}
	
	gtx.Constraints = cs
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C, ) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoIconLabelObj) ObjectType() (string) {
	return "GoIconLabelObj"
}

func (ob *GoIconLabelObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoIconLabelObj) Size() (int) {
	return ob.size
}

// Layout displays the icon with its size set to the X minimum constraint.
func (ob *GoIconLabelObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx)
	textColorMacro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()
	dims := layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				gtx.Constraints.Min = image.Point{X: ob.size}
				return ob.layoutIcon(gtx)
			})
		}),
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				paint_gio.ColorOp{Color: ob.ForeColor.NRGBA()}.Add(gtx.Ops)
				return widget_int.GioLabel{}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.label, textColor)
			})
		}),
	)
	return dims
}

func (ob *GoIconLabelObj) layoutIcon(gtx layout_gio.Context) layout_gio.Dimensions {
	rect := image.Point{X: ob.size, Y: ob.size}
	defer clip_gio.Rect{Max: rect}.Push(gtx.Ops).Pop()

	icon := ob.image(ob.size, ob.iconColor)
	icon.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	// add the events handler to receive widget pointer events
	//ob.SignalEvents(gtx)

	return layout_gio.Dimensions{
		Size: icon.Size(),
	}
}

func (ob *GoIconLabelObj) image(sz int, color GoColor) paint_gio.ImageOp {
	if sz == ob.imgSize && color == ob.imgColor {
		return ob.op
	}
	m, _ := iconvg.DecodeMetadata(ob.icon)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = f32color.NRGBAToLinearRGBA(color.NRGBA())
	iconvg.Decode(&ico, ob.icon, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ob.op = paint_gio.NewImageOp(img)
	ob.imgSize = sz
	ob.imgColor = color
	return ob.op
}