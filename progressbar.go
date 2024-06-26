// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/progressbar.go */

package utopia

import (
	"image"
	"image/color"
	//"log"
	
	layout_gio"github.com/utopiagio/gio/layout"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"

	"github.com/utopiagio/utopia/metrics"
)

type GoProgressBarObj struct {
	GioObject
	GioWidget
	color      GoColor
	trackColor GoColor
	progress   int
	totalSteps int
	thickness  int
}

func GoProgressBar(parent GoObject, totalSteps int) *GoProgressBarObj {
	theme := GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{4,4,4,4},
		GoSize: GoSize{0, 0, 100, 20, 16777215, 16777215, 100, 20},
		Visible: true,
		//target: nil,
		tag: tagCounter,
	}
	hProgress :=  &GoProgressBarObj{
		GioObject: object,
		GioWidget: widget,
		color:    theme.GoPalette.ContrastBg,
		trackColor: NRGBAColor(MulAlpha(theme.ColorFg.NRGBA(), 0x88)),
		totalSteps: totalSteps,
		progress: 0,
		thickness: 12,
	}
	parent.AddControl(hProgress)
	return hProgress
}

func (ob *GoProgressBarObj) LineThickness() (int) {
	return ob.thickness
}

func (ob *GoProgressBarObj) ObjectType() (string) {
	return "GoProgressBarObj"
}

func (ob *GoProgressBarObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoProgressBarObj) Progress() (int) {
	return ob.progress
}

func (ob *GoProgressBarObj) TotalSteps() (int) {
	return ob.totalSteps
}

func (ob *GoProgressBarObj) SetLineThickness(thickness int) {
	ob.thickness = thickness
}

func (ob *GoProgressBarObj) SetProgress(progress int) {
	if progress > ob.totalSteps {
		progress = ob.totalSteps
	} else if progress < 0 {
		progress = 0
	}
	ob.progress = progress
}

func (ob *GoProgressBarObj) SetTotalSteps(totalSteps int) {
	if totalSteps < ob.progress {
		totalSteps = ob.progress
	}
	ob.totalSteps = totalSteps
}

func (ob *GoProgressBarObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
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

func (ob *GoProgressBarObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	shader := func(width int, color color.NRGBA) layout_gio.Dimensions {
		maxHeight := unit_gio.Dp(ob.thickness)
		rr := gtx.Dp(unit_gio.Dp(ob.thickness / 2))

		d := image.Point{X: width, Y: gtx.Dp(maxHeight)}

		defer clip_gio.UniformRRect(image.Rectangle{Max: image.Pt(width, d.Y)}, rr).Push(gtx.Ops).Pop()
		paint_gio.ColorOp{Color: color}.Add(gtx.Ops)
		paint_gio.PaintOp{}.Add(gtx.Ops)

		return layout_gio.Dimensions{Size: d}
	}
	//log.Println("ob.Width: ", ob.Width, "...................")
	//log.Println("gtx.Constraints.Max.X: ", gtx.Constraints.Max.X, "...................")
	progressBarWidth := gtx.Constraints.Max.X
	//log.Println("progressBarWidth: ", progressBarWidth, "...................")
	dims := layout_gio.Flex{Alignment: layout_gio.Start}.Layout(gtx,
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.Stack{Alignment: layout_gio.W}.Layout(gtx,
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					return shader(progressBarWidth, ob.trackColor.NRGBA())
				}),
				layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
					fillWidth := int(float32(progressBarWidth) * clamp1(float32(ob.progress) / float32(ob.totalSteps)))
					fillColor := ob.color.NRGBA()
					if !gtx.Enabled() {
						fillColor = DisabledBlend(fillColor)
					}
					return shader(fillWidth, fillColor)
				}),
			)
		}),
	)
	return dims
}

// clamp1 limits v to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
	}
}
