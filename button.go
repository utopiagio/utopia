/* button.go */

package utopia

import (
	//"log"
	"image"
	"image/color"
	"math"

	//"github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	//material_gio "github.com/utopiagio/gio/widget/material"
	semantic_gio "github.com/utopiagio/gio/io/semantic"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
)

/*type ButtonStyle struct {
	Text string
	// Color is the text color.
	Color        color.NRGBA
	Font         text.Font
	TextSize     unit.Sp
	Background   color.NRGBA
	CornerRadius unit.Dp
	Inset        layout.Inset
	Button       *widget.Clickable
	shaper       *text.Shaper
}*/



//func Button(th *Theme, button *widget.Clickable, txt string) ButtonStyle {

func GoButton(parent GoObject, text string) (hObj *GoButtonObj) {
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hButton := &GoButtonObj{
		goObject: object,
		goWidget: widget,
		
		fontSize: theme.TextSize,
		text: text,
		color: theme.ContrastFg,
		background: theme.ContrastBg,
		cornerRadius: 4,
		inset: layout_gio.Inset{
			Top: 10, Bottom: 10,
			Left: 12, Right: 12,
		},
		clickable: new(widget_gio.Clickable),
		shaper: theme.Shaper,
		onClick: nil,
	}
	parent.addControl(hButton)
	return hButton
}

type GoButtonObj struct {
	goObject
	goWidget
	//gio material_gio.ButtonStyle
	//theme *GoThemeObj
	font text_gio.Font
	fontSize unit_gio.Sp
	text string
	color GoColor
	background GoColor
	cornerRadius unit_gio.Dp
	inset layout_gio.Inset
	clickable *widget_gio.Clickable
	shaper *text_gio.Shaper
	onClick func()
	onFocus func()
	onHover func()
	onPress func()
	//textAlign text.Alignment
}

func (ob *GoButtonObj) Clicked() bool {
	return ob.clickable.Clicked()
}

func (ob *GoButtonObj) Focused() bool {
	return ob.clickable.Focused()
}

func (ob *GoButtonObj) Hovered() bool {
	return ob.clickable.Hovered()
}

func (ob *GoButtonObj) Pressed() bool {
	return ob.clickable.Pressed()
}

func (ob *GoButtonObj) SetOnClick(f func()) {
	ob.onClick = f
}

func (ob *GoButtonObj) SetOnFocus(f func()) {
	ob.onFocus = f
}

func (ob *GoButtonObj) SetOnHover(f func()) {
	ob.onHover = f
}

func (ob *GoButtonObj) SetOnPress(f func()) {
	ob.onPress = f
}

func (ob *GoButtonObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}

func (ob *GoButtonObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.visible {
	//margin := layout_gio.Inset(ob.margin.Left)
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

func (ob *GoButtonObj) objectType() (string) {
	return "GoButtonObj"
}

func (ob *GoButtonObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	return ButtonLayoutStyle{
		background:   ob.background.NRGBA(),
		cornerRadius: ob.cornerRadius,
		clickable:    ob.clickable,
	}.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		return ob.inset.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
			//log.Println("Button label color:", ob.color.NRGBA())
			paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			return widget_gio.Label{Alignment: text_gio.Middle}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.text)
		})
	})
}

type ButtonLayoutStyle struct {
	background   color.NRGBA
	cornerRadius unit_gio.Dp
	clickable    *widget_gio.Clickable
}

func (ob ButtonLayoutStyle) Layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	min := gtx.Constraints.Min
	return ob.clickable.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.Button.Add(gtx.Ops)
		return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
			layout_gio.Expanded(func(gtx layout_gio.Context) layout_gio.Dimensions {
				rr := gtx.Dp(ob.cornerRadius)
				defer clip_gio.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				background := ob.background
				switch {
				case gtx.Queue == nil:
					background = DisabledBlend(ob.background)
				case ob.clickable.Hovered() || ob.clickable.Focused():
					background = HoveredBlend(ob.background)
				}
				paint_gio.Fill(gtx.Ops, background)
				for _, c := range ob.clickable.History() {
					drawInk(gtx, c)
				}
				return layout_gio.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
				gtx.Constraints.Min = min
				return layout_gio.Center.Layout(gtx, w)
			}),
		)
	})
}

/*func ButtonLayout(th *GoThemeObj, button *widget_gio.Clickable) ButtonLayoutStyle {
	return ButtonLayoutStyle{
		Button:       button,
		Background:   th.Palette.ContrastBg,
		CornerRadius: 4,
	}
}*/

type GoIconButtonObj struct {
	goObject
	goWidget

	background GoColor
	// Color is the icon color.
	color 	GoColor
	icon  *widget_gio.Icon
	// Size is the icon size.
	size        unit_gio.Dp
	inset       layout_gio.Inset
	clickable   *widget_gio.Clickable
	description string
}

func GoIconButton(parent GoObject, icon *widget_gio.Icon, description string) *GoIconButtonObj {
	var theme *GoThemeObj = goApp.Theme()
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hIconButton := &GoIconButtonObj{
		goObject: object,
		goWidget: widget,
		background:  theme.ContrastBg,
		color:       theme.ContrastFg,
		icon:        icon,
		size:        24,
		inset:       layout_gio.UniformInset(12),
		clickable: 	 new(widget_gio.Clickable),
		description: description,
	}
	parent.addControl(hIconButton)
	return hIconButton
}

func (ob *GoIconButtonObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	m := op_gio.Record(gtx.Ops)
	dims := ob.clickable.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		semantic_gio.Button.Add(gtx.Ops)
		if d := ob.description; d != "" {
			semantic_gio.DescriptionOp(ob.description).Add(gtx.Ops)
		}
		return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
			layout_gio.Expanded(func(gtx layout_gio.Context) layout_gio.Dimensions {
				rr := (gtx.Constraints.Min.X + gtx.Constraints.Min.Y) / 4
				defer clip_gio.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
				background := ob.background.NRGBA()
				switch {
				case gtx.Queue == nil:
					background = DisabledBlend(ob.background.NRGBA())
				case ob.clickable.Hovered() || ob.clickable.Focused():
					background = HoveredBlend(ob.background.NRGBA())
				}
				paint_gio.Fill(gtx.Ops, background)
				for _, c := range ob.clickable.History() {
					drawInk(gtx, c)
				}
				return layout_gio.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
				return ob.inset.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
					size := gtx.Dp(ob.size)
					if ob.icon != nil {
						gtx.Constraints.Min = image.Point{X: size}
						ob.icon.Layout(gtx, ob.color.NRGBA())
					}
					return layout_gio.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			}),
		)
	})
	c := m.Stop()
	bounds := image.Rectangle{Max: dims.Size}
	defer clip_gio.Ellipse(bounds).Push(gtx.Ops).Pop()
	c.Add(gtx.Ops)
	return dims
}

func (ob *GoIconButtonObj) objectType() (string) {
	return "GoIconButtonObj"
}

func drawInk(gtx layout_gio.Context, c widget_gio.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op_gio.InvalidateOp{}.Add(gtx.Ops)
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	// Cover the entire constraints min rectangle and
	// apply curve values to size and color.
	size = int(float32(size) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint_gio.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op_gio.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip_gio.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint_gio.PaintOp{}.Add(gtx.Ops)
}
