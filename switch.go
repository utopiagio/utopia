// SPDX-License-Identifier: Unlicense OR MIT

package utopia

import (
	"image"
	//"image/color"

	//"github.com/utopiagio/gio/internal/f32color"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	widget_gio "github.com/utopiagio/gio/widget"
)

type GoSwitchObj struct {
	goObject
	goWidget
	description string
	color       struct {
		enabled  GoColor
		disabled GoColor
		track    GoColor
	}
	goSwitch *widget_gio.Bool
	onChange func(bool)
	onFocus func()
	onHover func()
	onPress func()
}

// Switch is for selecting a boolean value.
func GoSwitch(parent GoObject, description string) *GoSwitchObj {
	var theme *GoThemeObj = goApp.Theme()
	var swtch *widget_gio.Bool = new(widget_gio.Bool)
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hSwitch := &GoSwitchObj{
		goObject: object,
		goWidget: widget,
		description: description,
		goSwitch: swtch,
	}
	hSwitch.color.enabled = theme.ContrastBg
	hSwitch.color.disabled = theme.ColorBg
	hSwitch.color.track = NRGBAColor(MulAlpha(theme.ColorFg.NRGBA(), 0x88))
	parent.addControl(hSwitch)
	return hSwitch
}

func (ob *GoSwitchObj) Changed() bool {
	return ob.goSwitch.Changed()
}

func (ob *GoSwitchObj) Focused() bool {
	return ob.goSwitch.Focused()
}

func (ob *GoSwitchObj) Hovered() bool {
	return ob.goSwitch.Hovered()
}

func (ob *GoSwitchObj) Pressed() bool {
	return ob.goSwitch.Pressed()
}

func (ob *GoSwitchObj) SetOnChange(f func(bool)) {
	ob.onChange = f
}

func (ob *GoSwitchObj) SetOnFocus(f func()) {
	ob.onFocus = f
}

func (ob *GoSwitchObj) SetOnHover(f func()) {
	ob.onHover = f
}

func (ob *GoSwitchObj) SetOnPress(f func()) {
	ob.onPress = f
}

func (ob *GoSwitchObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
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

// Layout updates the switch and displays it.
func (ob *GoSwitchObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	trackWidth := gtx.Dp(36)
	trackHeight := gtx.Dp(16)
	thumbSize := gtx.Dp(20)
	trackOff := (thumbSize - trackHeight) / 2

	// Draw track.
	trackCorner := trackHeight / 2
	trackRect := image.Rectangle{Max: image.Point{
		X: trackWidth,
		Y: trackHeight,
	}}
	col := ob.color.disabled.NRGBA()
	if ob.goSwitch.Value {
		col = ob.color.enabled.NRGBA()
	}
	if gtx.Queue == nil {
		col = DisabledBlend(col)
	}
	trackColor := ob.color.track.NRGBA()
	t := op_gio.Offset(image.Point{Y: trackOff}).Push(gtx.Ops)
	cl := clip_gio.UniformRRect(trackRect, trackCorner).Push(gtx.Ops)
	paint_gio.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)
	cl.Pop()
	t.Pop()

	// Draw thumb ink.
	inkSize := gtx.Dp(44)
	rr := inkSize / 2
	inkOff := image.Point{
		X: trackWidth/2 - rr,
		Y: -rr + trackHeight/2 + trackOff,
	}
	t = op_gio.Offset(inkOff).Push(gtx.Ops)
	gtx.Constraints.Min = image.Pt(inkSize, inkSize)
	cl = clip_gio.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops)
	for _, p := range ob.goSwitch.History() {
		drawInk(gtx, p)
	}
	cl.Pop()
	t.Pop()

	// Compute thumb offset.
	if ob.goSwitch.Value {
		xoff := trackWidth - thumbSize
		defer op_gio.Offset(image.Point{X: xoff}).Push(gtx.Ops).Pop()
	}

	thumbRadius := thumbSize / 2

	circle := func(x, y, r int) clip_gio.Op {
		b := image.Rectangle{
			Min: image.Pt(x-r, y-r),
			Max: image.Pt(x+r, y+r),
		}
		return clip_gio.Ellipse(b).Op(gtx.Ops)
	}
	// Draw hover.
	if ob.goSwitch.Hovered() || ob.goSwitch.Focused() {
		r := thumbRadius * 10 / 17
		background := MulAlpha(ob.color.enabled.NRGBA(), 70)
		paint_gio.FillShape(gtx.Ops, background, circle(thumbRadius, thumbRadius, r))
	}

	// Draw thumb shadow, a translucent disc slightly larger than the
	// thumb itself.
	// Center shadow horizontally and slightly adjust its Y.
	paint_gio.FillShape(gtx.Ops, argb(0x55000000), circle(thumbRadius, thumbRadius+gtx.Dp(.25), thumbRadius+1))

	// Draw thumb.
	paint_gio.FillShape(gtx.Ops, col, circle(thumbRadius, thumbRadius, thumbRadius))

	// Set up click area.
	clickSize := gtx.Dp(40)
	clickOff := image.Point{
		X: (thumbSize - clickSize) / 2,
		Y: (trackHeight-clickSize)/2 + trackOff,
	}
	defer op_gio.Offset(clickOff).Push(gtx.Ops).Pop()
	sz := image.Pt(clickSize, clickSize)
	defer clip_gio.Ellipse(image.Rectangle{Max: sz}).Push(gtx.Ops).Pop()
	ob.goSwitch.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		if d := ob.description; d != "" {
			semantic_gio.DescriptionOp(d).Add(gtx.Ops)
		}
		semantic_gio.Switch.Add(gtx.Ops)
		return layout_gio.Dimensions{Size: sz}
	})

	dims := image.Point{X: trackWidth, Y: thumbSize}
	return layout_gio.Dimensions{Size: dims}
}

func (ob *GoSwitchObj) objectType() (string) {
	return "GoSwitchObj"
}