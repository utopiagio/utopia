// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoCheckBoxObj struct {
	GioObject
	GioWidget
	checkable
	checkBox *widget_gio.Bool
}

func GoCheckBox(parent GoObject, label string) *GoCheckBoxObj {
	var theme *GoThemeObj = GoApp.Theme()
	var checkbox *widget_gio.Bool = new(widget_gio.Bool)
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		Visible: true,
	}
	hCheckBox := &GoCheckBoxObj{
		GioObject: object,
		GioWidget: widget,
		checkBox: checkbox,
		checkable: checkable{
			label:              label,
			color:              theme.ColorFg,
			iconColor:          theme.ContrastBg,
			fontSize:           theme.TextSize, // * 14.0 / 16.0,
			size:               26,
			shaper:             theme.Shaper,
			checkedStateIcon:   theme.Icon.CheckBoxChecked,
			uncheckedStateIcon: theme.Icon.CheckBoxUnchecked,
		},
	}
	parent.AddControl(hCheckBox)
	return hCheckBox
}

func (ob *GoCheckBoxObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
	}
	return dims
}

// Layout updates the checkBox and displays it.
func (ob *GoCheckBoxObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	return ob.checkBox.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.CheckBox.Add(gtx.Ops)
		return ob.checkable.layout(gtx, ob.checkBox.Value, ob.checkBox.Hovered() || ob.checkBox.Focused())
	})
}

func (ob *GoCheckBoxObj) ObjectType() (string) {
	return "GoCheckBoxObj"
}