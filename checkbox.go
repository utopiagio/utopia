// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/checkbox.go */

package utopia

import (
	//"log"
	"image"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	widget_int "github.com/utopiagio/utopia/internal/widget"
	widget_gio "github.com/utopiagio/gio/widget"

	"github.com/utopiagio/utopia/metrics"
)

type GoCheckBoxObj struct {
	GioObject
	GioWidget
	checkable widget_int.GioCheckable
	checkBox *widget_gio.Bool
}

func GoCheckBox(parent GoObject, label string) *GoCheckBoxObj {
	var theme *GoThemeObj = GoApp.Theme()
	var GioCheckbox *widget_gio.Bool = new(widget_gio.Bool)
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 200, 26, 16777215, 16777215, 200, 26},
		Visible: true,
	}
	hCheckBox := &GoCheckBoxObj{
		GioObject: object,
		GioWidget: widget,
		checkBox: GioCheckbox,
		checkable: widget_int.GioCheckable{
			Label:              label,
			Color:              theme.ColorFg.NRGBA(),
			IconColor:          theme.ContrastBg.NRGBA(),
			TextSize:           theme.TextSize, // * 14.0 / 16.0,
			Size:               16,
			Shaper:             theme.Shaper,
			CheckedStateIcon:   theme.Icon.CheckBoxChecked,
			UncheckedStateIcon: theme.Icon.CheckBoxUnchecked,
		},
	}
	parent.AddControl(hCheckBox)
	return hCheckBox
}

func (ob *GoCheckBoxObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

// Layout updates the checkBox and displays it.
func (ob *GoCheckBoxObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	dims := ob.checkBox.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.CheckBox.Add(gtx.Ops)
		chdims := ob.checkable.Layout(gtx, ob.checkBox.Value, ob.checkBox.Hovered() || ob.checkBox.Focused())
		return chdims
	})
	return dims
}

func (ob *GoCheckBoxObj) ObjectType() (string) {
	return "GoCheckBoxObj"
}

func (ob *GoCheckBoxObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}