// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/popupwindow.go */

package utopia

import (
	//"log"
	//"image"
	//"image/color"

	//clip_gio "github.com/utopiagio/gio/op/clip"
	layout_gio "github.com/utopiagio/gio/layout"
	//pointer_gio "github.com/utopiagio/gio/io/pointer"
)

func GoPopupWindow(parent GoObject) (hPopupWindow *GoPopupWindowObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		FocusPolicy: StrongFocus,
		Visible: false,
		tag: tagCounter,
	}
	hPopupWindow = &GoPopupWindowObj{GioObject: object, GioWidget: widget, alpha: 90}
	hPopupWindow.layout = GoVFlexBoxLayout(hPopupWindow)
	hPopupWindow.layout.SetMargin(3,3,3,3)
	hPopupWindow.layout.SetBorder(BorderSingleLine, 2, 3, Color_Blue)
	hPopupWindow.layout.SetPadding(3,3,3,3)
	
	hPopupWindow.layout.SetOnPointerPress(hPopupWindow.Click)
	hPopupWindow.layout.SetOnPointerEnter(nil)
	hPopupWindow.layout.SetOnPointerLeave(nil)
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

func (ob *GoPopupWindowObj) Click(e GoPointerEvent) {
	ob.Hide()
}

func (ob *GoPopupWindowObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		return ob.layout.Draw(gtx)
		/*dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
			})
		})*/
	}
	return dims
}

func (ob *GoPopupWindowObj) Hide() {
	ob.GioWidget.Hide()
	ob.layout.Clear()
	//ob.layout = GoVFlexBoxLayout(ob)
	//ob.layout.SetMargin(3,3,3,3)
	//ob.layout.SetBorder(BorderSingleLine, 1, 3, Color_Blue)
	//ob.layout.SetPadding(3,3,3,3)
	
}

// Layout draws the scrim using the provided animation. If the animation indicates
// that the scrim is not visible, this is a no-op.
/*func (ob *GoPopupWindowObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	dims := gtx.Constraints.Max
		//if !s.Visible() {
			//return layout.Dimensions{}
		//}
		ob.ReceiveEvents(gtx, nil)
		//gtx.Constraints.Min = gtx.Constraints.Max
		//currentAlpha := s.FinalAlpha
		if anim.Animating() {
			revealed := anim.Revealed(gtx)
			currentAlpha = uint8(float32(s.FinalAlpha) * revealed)
		}
		//color := th.Fg
		//color.A = currentAlpha
		defer clip_gio.Rect(image.Rectangle{Max: dims}).Push(gtx.Ops).Pop()
		fill := Color_WhiteSmoke.MulAlpha(ob.alpha)
		PaintRect(gtx, dims, fill)
		ob.SignalEvents(gtx)
		return layout_gio.Dimensions{Size: dims}
}*/

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

func (ob *GoPopupWindowObj)SetLayoutStyle(style GoLayoutStyle) {
	ob.DeleteControl(ob.layout)
	switch style {
	case NoLayout:
		ob.layout = nil
	case HBoxLayout:
		ob.layout = GoHBoxLayout(ob)
	case VBoxLayout:
		ob.layout = GoVBoxLayout(ob)	
	case HVBoxLayout:
		// Not Implemented *******************
	case HFlexBoxLayout:
		ob.layout = GoHFlexBoxLayout(ob)	
	case VFlexBoxLayout:						
		ob.layout = GoVFlexBoxLayout(ob)	
	case PopupMenuLayout:
		// Not Implemented *******************
	}
}

func (ob *GoPopupWindowObj) SetSize(width int, height int) {
	ob.Width = width
	ob.Height = height
	if ob.Visible {
		ob.ParentWindow().Refresh()
	}
}

func (ob *GoPopupWindowObj) Show() {
	ob.GioWidget.Show()
	ob.ParentWindow().Refresh()
	//ob.layout = GoVFlexBoxLayout(ob)
	//ob.layout.SetMargin(3,3,3,3)
	//ob.layout.SetBorder(BorderSingleLine, 1, 3, Color_Blue)
	//ob.layout.SetPadding(3,3,3,3)
	
}

func (ob *GoPopupWindowObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoPopupWindowObj) Layout() (*GoLayoutObj) {
	return ob.layout
}