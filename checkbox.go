// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoCheckBoxObj struct {
	goObject
	goWidget
	checkable
	checkBox *widget_gio.Bool
}

func GoCheckBox(parent GoObject, label string) *GoCheckBoxObj {
	var theme *GoThemeObj = goApp.Theme()
	var checkbox *widget_gio.Bool = new(widget_gio.Bool)
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hCheckBox := &GoCheckBoxObj{
		goObject: object,
		goWidget: widget,
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
	parent.addControl(hCheckBox)
	return hCheckBox
}

func (ob *GoCheckBoxObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.visible {
		dims = ob.goMargin.layout(gtx, func(gtx C) D {
			return ob.goBorder.layout(gtx, func(gtx C) D {
				return ob.goPadding.layout(gtx, func(gtx C) D {
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

func (ob *GoCheckBoxObj) objectType() (string) {
	return "GoCheckBoxObj"
}