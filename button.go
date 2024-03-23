// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/button.go */

package utopia

import (
	_ "log"
	"image"
	"image/color"
	"math"

	//"github.com/utopiagio/gio/font/gofont"
	
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	//material_gio "github.com/utopiagio/gio/widget/material"
	//semantic_gio "github.com/utopiagio/gio/io/semantic"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	text_gio "github.com/utopiagio/gio/text"
	font_gio "github.com/utopiagio/gio/font"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
	widget_int "github.com/utopiagio/utopia/internal/widget"

	"github.com/utopiagio/utopia/metrics"
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

func GoButton(parent GoObject, text string) (hObj *GoButtonObj) {
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 200, 36, 16777215, 16777215, 200, 36},
		FocusPolicy: StrongFocus,
		Visible: true,
	}
	hButton := &GoButtonObj{
		GioObject: object,
		GioWidget: widget,
		
		fontSize: theme.TextSize,
		text: text,
		textColor: theme.ContrastFg,
		faceColor: theme.ContrastBg,
		cornerRadius: 4,
		inset: layout_gio.Inset{
			Top: 8, Bottom: 8,
			Left: 10, Right: 10,
		},
		shaper: theme.Shaper,
	}
	hButton.SetOnKeyPress(nil)
	hButton.SetOnKeyRelease(nil)
	hButton.SetOnKeyEdit(nil)
	hButton.SetOnPointerPress(nil)
	hButton.SetOnPointerRelease(hButton.click)
	//hButton.SetOnPointerMove(nil)
	//hButton.SetOnPointerClick(nil)
	hButton.SetOnPointerEnter(nil)
	hButton.SetOnPointerLeave(nil)
	parent.AddControl(hButton)
	return hButton
}

type GoButtonObj struct {
	GioObject
	GioWidget
	//theme *GoThemeObj
	font font_gio.Font
	fontSize unit_gio.Sp
	text string
	textColor GoColor
	faceColor GoColor
	cornerRadius unit_gio.Dp
	inset layout_gio.Inset
	shaper *text_gio.Shaper
	onClick func()
	//textAlign text.Alignment
}

func (ob *GoButtonObj) click(e pointer_gio.Event) {
	if ob.onClick != nil {
		ob.onClick()
	}
}

func (ob *GoButtonObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
				return paddingDims
			})
			return borderDims
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoButtonObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx, nil)
	textColorMacro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.textColor.NRGBA()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()
	return ob.layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		insetDims := ob.inset.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
			//log.Println("Button label color:", ob.color.NRGBA())
			//paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
			dims := widget_int.GioLabel{Alignment: text_gio.Middle, MaxLines: 2,}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.text, textColor)
			return dims
		})
		return insetDims
	})
}

func (ob *GoButtonObj) ObjectType() (string) {
	return "GoButtonObj"
}

func (ob *GoButtonObj) SetFaceColor(color GoColor) {
	ob.faceColor = color
}

func (ob *GoButtonObj) SetOnClick(f func()) {
	ob.onClick = f
}

func (ob *GoButtonObj) SetText(text string) {
	ob.text = text
}

func (ob *GoButtonObj) SetTextColor(color GoColor) {
	ob.textColor = color
}

func (ob *GoButtonObj) Text() (text string) {
	return ob.text
}

func (ob *GoButtonObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoButtonObj) layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	min := gtx.Constraints.Min
	//semantic_gio.Button.Add(gtx.Ops)
	
	return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
		layout_gio.Expanded(func(gtx layout_gio.Context) layout_gio.Dimensions {
			rr := gtx.Dp(ob.cornerRadius)
			defer clip_gio.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
			background := ob.faceColor.NRGBA()
			switch {
			case !gtx.Enabled():
				background = DisabledBlend(background)
			case ob.IsHovered() || ob.HasFocus():
				background = HoveredBlend(background)
			}
			paint_gio.Fill(gtx.Ops, background)
			/*for _, c := range ob.clickable.History() {
				drawInk(gtx, c)
			}*/
			ob.SignalEvents(gtx)

			return layout_gio.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
			gtx.Constraints.Min = min
			return layout_gio.Center.Layout(gtx, w)
		}),
	)
}

type GoIconButtonObj struct {
	GioObject
	GioWidget
	// Color is the icon color.
	textColor GoColor
	faceColor GoColor
	cornerRadius unit_gio.Dp
	iconVG  *GoIconVGObj
	iconPNG *GoIconPNGObj
	// Size is the icon size.
	size        unit_gio.Dp
	inset       layout_gio.Inset
	onClick func()
	//clickable   *widget_gio.Clickable
	//description string
}

//func GoIconButton(parent GoObject, icon *GoIconObj), description string) *GoIconButtonObj {
func GoIconVGButton(parent GoObject, icon *GoIconVGObj) *GoIconButtonObj {
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 30, 30, 16777215, 16777215, 30, 30},
		Visible: true,
	}
	hIconButton := &GoIconButtonObj{
		GioObject: object,
		GioWidget: widget,
		faceColor: theme.ColorBg,
		textColor: theme.ColorFg,
		cornerRadius: 4,
		iconVG: icon,
		size: 20,
		inset: layout_gio.UniformInset(4),
	}
	hIconButton.SetOnPointerRelease(hIconButton.Click)
	hIconButton.SetOnPointerEnter(nil)
	hIconButton.SetOnPointerLeave(nil)
	parent.AddControl(hIconButton)
	return hIconButton
}

//func GoIconButton(parent GoObject, icon *GoIconObj), description string) *GoIconButtonObj {
func GoIconPNGButton(parent GoObject, icon *GoIconPNGObj) *GoIconButtonObj {
	var theme *GoThemeObj = GoApp.Theme()
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{0, 0, 24, 24, 16777215, 16777215, 24, 24},
		Visible: true,
	}
	hIconButton := &GoIconButtonObj{
		GioObject: object,
		GioWidget: widget,
		faceColor: theme.ColorBg,
		textColor: theme.ColorFg,
		cornerRadius: 4,
		iconPNG: icon,
		size: 20,
		inset: layout_gio.UniformInset(4),
	}
	hIconButton.SetOnPointerRelease(hIconButton.Click)
	hIconButton.SetOnPointerEnter(nil)
	hIconButton.SetOnPointerLeave(nil)
	parent.AddControl(hIconButton)
	return hIconButton
}

func (ob *GoIconButtonObj) Click(e pointer_gio.Event) {
	if ob.onClick != nil {
		ob.onClick()
	}
}

func (ob *GoIconButtonObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx C) D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
				return paddingDims
			})
			return borderDims
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

func (ob *GoIconButtonObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx, nil)
	return ob.layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		insetDims := ob.inset.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
			size := gtx.Dp(ob.size)
			if ob.iconVG != nil || ob.iconPNG != nil {
				gtx.Constraints.Min = image.Point{X: size}
				if ob.iconVG != nil {
					ob.iconVG.Layout(gtx, ob.textColor)
				} else if ob.iconPNG != nil {
					ob.iconPNG.Layout(gtx)
				}
			}
			return layout_gio.Dimensions{
				Size: image.Point{X: size, Y: size},
			}
		})
		return insetDims
	})
}

func (ob *GoIconButtonObj) layout(gtx layout_gio.Context, w layout_gio.Widget) layout_gio.Dimensions {
	
	min := gtx.Constraints.Min

	//semantic_gio.Button.Add(gtx.Ops)
	
	return layout_gio.Stack{Alignment: layout_gio.Center}.Layout(gtx,
		layout_gio.Expanded(func(gtx layout_gio.Context) layout_gio.Dimensions {
			rr := gtx.Dp(ob.cornerRadius)
			defer clip_gio.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, rr).Push(gtx.Ops).Pop()
			background := ob.faceColor.NRGBA()
			switch {
			case !gtx.Enabled():
				background = DisabledBlend(background)
			case ob.IsHovered() || ob.HasFocus():
				background = HoveredBlend(background)
			}
			paint_gio.Fill(gtx.Ops, background)
			/*for _, c := range ob.clickable.History() {
				drawInk(gtx, c)
			}*/
			ob.SignalEvents(gtx)
			return layout_gio.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout_gio.Stacked(func(gtx layout_gio.Context) layout_gio.Dimensions {
			gtx.Constraints.Min = min
			return layout_gio.Center.Layout(gtx, w)
		}),
	)
}
/*func (ob *GoIconButtonObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
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
				case !gtx.Enabled():
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
}*/

func (ob *GoIconButtonObj) ObjectType() (string) {
	return "GoIconButtonObj"
}

func (ob *GoIconButtonObj) SetFaceColor(color GoColor) {
	ob.faceColor = color
}

func (ob *GoIconButtonObj) SetOnClick(f func()) {
	ob.onClick = f
}

func (ob *GoIconButtonObj) SetSize(size int) {
	ob.size = unit_gio.Dp(size)
}

func (ob *GoIconButtonObj) SetTextColor(color GoColor) {
	ob.textColor = color
}

func (ob *GoIconButtonObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

// Layout and update the button state
/*func clickableLayout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	b.update(gtx)
	m := op.Record(gtx.Ops)
	dims := w(gtx)
	c := m.Stop()
	defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
	disabled := !gtx.Enabled()
	semantic.DisabledOp(disabled).Add(gtx.Ops)
	b.click.Add(gtx.Ops)
	if !disabled {
		keys := key.Set("⏎|Space")
		if !b.focused {
			keys = ""
		}
		key.InputOp{Tag: &b.keyTag, Keys: keys}.Add(gtx.Ops)
		if b.requestFocus {
			key.FocusOp{Tag: &b.keyTag}.Add(gtx.Ops)
			b.requestFocus = false
		}
	} else {
		b.focused = false
	}
	c.Add(gtx.Ops)
	for len(b.history) > 0 {
		c := b.history[0]
		if c.End.IsZero() || gtx.Now.Sub(c.End) < 1*time.Second {
			break
		}
		n := copy(b.history, b.history[1:])
		b.history = b.history[:n]
	}
	return dims
}*/

func (ob *GoButtonObj) drawInk(gtx layout_gio.Context, c widget_gio.Press) {
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
		gtx.Execute(op_gio.InvalidateCmd{})
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
