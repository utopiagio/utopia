// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/radiobutton.go */

package utopia

import (
	"log"
	"image"

	"github.com/utopiagio/utopia/metrics"

	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	widget_int "github.com/utopiagio/utopia/internal/widget"
)

type GoRadioButtonObj struct {
	GioObject
	GioWidget
	Checkable *widget_int.GioCheckable
	Key   string
	Group *GoButtonGroupObj
	selected bool
	onChange func(bool)
	onFocus func(string)
	onHover func(string)
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func GoRadioButton(parent GoObject, group *GoButtonGroupObj, key, label string) *GoRadioButtonObj {
	if group == nil {
		group = GoButtonGroup()
	}
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 200, 26, 16777215, 16777215, 200, 26},
		Visible: true,
	}
	checkable := &widget_int.GioCheckable{
		Label: 				label,
		Color:              theme.ColorFg.NRGBA(),
		IconColor:          theme.ContrastBg.NRGBA(),
		TextSize:           theme.TextSize, // * 14.0 / 16.0,
		Size:               16,
		Shaper:             theme.Shaper,
		CheckedStateIcon:   theme.Icon.RadioChecked,
		UncheckedStateIcon: theme.Icon.RadioUnchecked,
	}
	hRadioButton := &GoRadioButtonObj{
		GioObject: object,
		GioWidget: widget,
		Checkable: checkable,
		Key: key,
		Group: group,
	}
	parent.AddControl(hRadioButton)
	return hRadioButton
}

func (ob *GoRadioButtonObj) Selected() (bool) {
	return ob.selected
}

func (ob *GoRadioButtonObj) Focused() (bool) {
	key, focused := ob.Group.Focused()
	if focused {
		if key == ob.Key {
			return true
		}
	}
	return false
}

func (ob *GoRadioButtonObj) Hovered() (bool) {
	key, hovered := ob.Group.Hovered()
	if hovered {
		if key == ob.Key {
			return true
		}
	}
	return false
}

func (ob *GoRadioButtonObj) SetOnChange(f func(bool)) {
	ob.onChange = f
}

func (ob *GoRadioButtonObj) SetOnFocus(f func(string)) {
	ob.onFocus = f
}

func (ob *GoRadioButtonObj) SetOnHover(f func(string)) {
	ob.onHover = f
}

func (ob *GoRadioButtonObj) State() (bool) {
	return ob.Group.Value() == ob.Key
}

func (ob *GoRadioButtonObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	log.Println("GoRadioButtonObj::Draw()")
	cs := gtx.Constraints
	log.Println("gtx.Constraints Min = (", cs.Min.X, cs.Min.Y, ") Max = (", cs.Max.X, cs.Max.Y, ")")
	
	width := metrics.DpToPx(GoDpr, ob.Width)
	height := metrics.DpToPx(GoDpr, ob.Height)
	minWidth := metrics.DpToPx(GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		cs.Min.X = min(cs.Max.X, width)			// constrain to ob.Width
		cs.Max.X = min(cs.Max.X, width)			// constrain to ob.Width
	case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = minWidth						// set to ob.MinWidth
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = max(cs.Min.X, minWidth)		// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = maxWidth						// set to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		cs.Min.Y = min(cs.Max.Y, height)		// constrain to ob.Height
		cs.Max.Y = min(cs.Max.Y, height)		// constrain to ob.Height
	case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = minHeight					// set to ob.MinHeight
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = min(cs.Min.Y, minHeight)		// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = maxHeight					// set to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	case ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}
	gtx.Constraints = cs
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
		log.Println("GoRadioButton::Height: ", dims.Size.Y)
	}
	return dims
}

// Layout updates enum and displays the radio button.
func (ob *GoRadioButtonObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.Group.Update(gtx)
	if !ob.selected {
		if ob.Group.Value() == ob.Key {
			ob.selected = true
			if ob.onChange != nil {
				ob.onChange(true)
			}
		}
	} else {
		if ob.Group.Value() != ob.Key {
			ob.selected = false
			if ob.onChange != nil {
				ob.onChange(false)
			}
		}
	}
	hovered, hovering := ob.Group.Hovered()
	focus, focused := ob.Group.Focused()
	return ob.Group.Layout(gtx, ob.Key, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.RadioButton.Add(gtx.Ops)
		highlight := hovering && hovered == ob.Key || focused && focus == ob.Key
		return ob.Checkable.Layout(gtx, ob.Group.Value() == ob.Key, highlight)
	})
}

func (ob *GoRadioButtonObj) ObjectType() (string) {
	return "GoRadioButtonObj"
}

func (ob *GoRadioButtonObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}