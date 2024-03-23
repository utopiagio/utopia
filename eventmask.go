// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/eventmask.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"

	clip_gio "github.com/utopiagio/gio/op/clip"
	layout_gio "github.com/utopiagio/gio/layout"
	//pointer_gio "github.com/utopiagio/gio/io/pointer"

	"github.com/utopiagio/utopia/metrics"
)

func GoEventMask(parent GoObject) (hEventMask *GoEventMaskObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 100, 100, 1000, 1000, 100, 100},
		FocusPolicy: StrongFocus,
		Visible: false,
	}
	hEventMask = &GoEventMaskObj{GioObject: object, GioWidget: widget, alpha: 90}
	hEventMask.SetOnPointerClick(nil)
	return hEventMask
}

type GoEventMaskObj struct {
	GioObject
	GioWidget
	// FinalAlpha is the final opacity of the scrim on a scale from 0 to 255.
	alpha uint8
}

func (ob *GoEventMaskObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	//gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			/*return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {*/
					return ob.Layout(gtx)
				})
			/*})
		})*/
	}
	ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
	ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	return dims
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
func (ob *GoEventMaskObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	dims := gtx.Constraints.Max
		//if !s.Visible() {
			//return layout.Dimensions{}
		//}
		ob.ReceiveEvents(gtx, nil)
		//gtx.Constraints.Min = gtx.Constraints.Max
		//currentAlpha := s.FinalAlpha
		/*if anim.Animating() {
			revealed := anim.Revealed(gtx)
			currentAlpha = uint8(float32(s.FinalAlpha) * revealed)
		}*/
		//color := th.Fg
		//color.A = currentAlpha
		defer clip_gio.Rect(image.Rectangle{Max: dims}).Push(gtx.Ops).Pop()
		fill := Color_WhiteSmoke.MulAlpha(ob.alpha)
		PaintRect(gtx, dims, fill)
		ob.SignalEvents(gtx)
		return layout_gio.Dimensions{Size: dims}
}

func (ob *GoEventMaskObj) ObjectType() (string) {
	return "GoEventMaskObj"
}

func (ob *GoEventMaskObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}