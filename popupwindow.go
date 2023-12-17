// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/popupwindow.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"

	clip_gio "github.com/utopiagio/gio/op/clip"
	layout_gio "github.com/utopiagio/gio/layout"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
)

func GoPopupWindow(parent GoObject) (hPopupWindow *GoPopupWindowObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		FocusPolicy: StrongFocus,
		Visible: false,
	}
	hPopupWindow = &GoPopupWindowObj{GioObject: object, GioWidget: widget, alpha: 90}
	hPopupWindow.layout = GoPopupMenuLayout(hPopupWindow)
	//hPopupMenu.layout.SetPadding(3,3,3,3)
	hPopupWindow.layout.SetBorder(BorderSingleLine, 1, 3, Color_Blue)
	hPopupWindow.SetOnPointerRelease(hPopupWindow.Click)
	hPopupWindow.SetOnPointerEnter(nil)
	hPopupWindow.SetOnPointerLeave(nil)
	return hPopupWindow
}

type GoPopupWindowObj struct {
	GioObject
	GioWidget
	layout *GoLayoutObj
	// FinalAlpha is the final opacity of the scrim on a scale from 0 to 255.
	alpha uint8
	//visible bool
}

func (ob *GoPopupWindowObj) Click(e pointer_gio.Event) {
	ob.Hide()
}

func (ob *GoPopupWindowObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			/*return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {*/
					return ob.Layout(gtx)
				})
			/*})
		})*/
	}
	return dims
}

func (ob *GoPopupWindowObj) Hide() {
	ob.GioWidget.Hide()
	ob.layout = GoVFlexBoxLayout(ob)
		//ob.layout.SetPadding(3,3,3,3)
	ob.layout.SetBorder(BorderSingleLine, 1, 3, Color_Blue)
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
func (ob *GoPopupWindowObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	dims := gtx.Constraints.Max
		//if !s.Visible() {
			//return layout.Dimensions{}
		//}
		ob.ReceiveEvents(gtx)
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

func (ob *GoPopupWindowObj) ObjectType() (string) {
	return "GoPopupWindowObj"
}

/*func (ob *GoPopupMenuObj) Style() (style GoModalStyle) {
	return ob.style
}*/

/*func (ob *GoPopupMenuObj) SetStyle(style GoModalStyle) {
	ob.style = style
	if ob.style == GoPopupWindow {
		ob.layout = GoVFlexBoxLayout(ob)
	} else if ob.style == GoPopupMenu {
		ob.layout = GoPopupMenuLayout(ob)
		//ob.layout.SetPadding(3,3,3,3)
		ob.layout.SetBorder(BorderSingleLine, 1, 3, Color_Blue)
		//ob.SetPadding(0,25,0,0)

	}
}*/

func (ob *GoPopupWindowObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}