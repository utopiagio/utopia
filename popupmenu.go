// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/popupmenu.go */

package utopia

import (
	//"log"
	"image"
	//"image/color"

	clip_gio "github.com/utopiagio/gio/op/clip"
	layout_gio "github.com/utopiagio/gio/layout"
	pointer_gio "github.com/utopiagio/gio/io/pointer"

	"github.com/utopiagio/utopia/metrics"
)

func GoPopupMenu(parent GoObject) (hPopupMenu *GoPopupMenuObj) {
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 100, 100, 1000, 1000, 100, 100},
		FocusPolicy: StrongFocus,
		Visible: false,
	}
	hPopupMenu = &GoPopupMenuObj{GioObject: object, GioWidget: widget, alpha: 90}
	hPopupMenu.layout = GoPopupMenuLayout(hPopupMenu)
	//hPopupMenu.layout.SetPadding(3,3,3,3)
	hPopupMenu.layout.SetBorder(BorderSingleLine, 1, 2, Color_LightGray)
	hPopupMenu.SetOnPointerRelease(hPopupMenu.Click)
	hPopupMenu.SetOnPointerEnter(nil)
	hPopupMenu.SetOnPointerLeave(nil)
	return hPopupMenu
}

type GoPopupMenuObj struct {
	GioObject
	GioWidget
	layout *GoLayoutObj
	// FinalAlpha is the final opacity of the scrim on a scale from 0 to 255.
	alpha uint8
}

func (ob *GoPopupMenuObj) Clear() {
	ob.layout = GoPopupMenuLayout(ob)
	ob.layout.SetBorder(BorderSingleLine, 1, 2, Color_LightGray)
}

func (ob *GoPopupMenuObj) Click(e pointer_gio.Event) {
	ob.ParentWindow().ClearPopupMenus()
}

func (ob *GoPopupMenuObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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
func (ob *GoPopupMenuObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
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

func (ob *GoPopupMenuObj) ObjectType() (string) {
	return "GoPopupMenuObj"
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

func (ob *GoPopupMenuObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}