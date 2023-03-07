// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoRadioButtonObj struct {
	goObject
	goWidget
	checkable
	key   string
	group *widget_gio.Enum
	onChange func()
	onFocus func(string)
	onHover func(string)
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func GoRadioButton(parent GoObject, key, label string) *GoRadioButtonObj {
	var theme *GoThemeObj = goApp.Theme()
	var group *widget_gio.Enum = new(widget_gio.Enum)
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hRadioButton := &GoRadioButtonObj{
		goObject: object,
		goWidget: widget,
		checkable: checkable{
			label: label,
			color:              theme.ColorFg,
			iconColor:          theme.ContrastBg,
			fontSize:           theme.TextSize, // * 14.0 / 16.0,
			size:               26,
			shaper:             theme.Shaper,
			checkedStateIcon:   theme.Icon.RadioChecked,
			uncheckedStateIcon: theme.Icon.RadioUnchecked,
		},
		key: key,
		group: group,
	}
	parent.addControl(hRadioButton)
	return hRadioButton
}

func (ob *GoRadioButtonObj) Changed() bool {
	return ob.group.Changed()
}

func (ob *GoRadioButtonObj) Focused() (string, bool) {
	return ob.group.Focused()
}

func (ob *GoRadioButtonObj) Hovered() (string, bool) {
	return ob.group.Hovered()
}

func (ob *GoRadioButtonObj) SetOnChange(f func()) {
	ob.onChange = f
}

func (ob *GoRadioButtonObj) SetOnFocus(f func(string)) {	// 
	ob.onFocus = f
}

func (ob *GoRadioButtonObj) SetOnHover(f func(string)) {
	ob.onHover = f
}

func (ob *GoRadioButtonObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

// Layout updates enum and displays the radio button.
func (ob *GoRadioButtonObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	hovered, hovering := ob.group.Hovered()
	focus, focused := ob.group.Focused()
	return ob.group.Layout(gtx, ob.key, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.RadioButton.Add(gtx.Ops)
		highlight := hovering && hovered == ob.key || focused && focus == ob.key
		return ob.checkable.layout(gtx, ob.group.Value == ob.key, highlight)
	})
}

func (ob *GoRadioButtonObj) objectType() (string) {
	return "GoRadioButtonObj"
}