// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/checkbox.go */

package utopia

import (
	"log"
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
	log.Println("GoLabelObj::Draw()")
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
		log.Println("FixedWidth............")
		//log.Println("object Width = (", width, " )")
		cs.Min.X = min(cs.Max.X, width)
		log.Println("cs.Min.X = (", cs.Min.X, " )")
		cs.Max.X = min(cs.Max.X, width)
		log.Println("cs.Max.X = (", cs.Max.X, " )")
	/*case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = min(cs.Min.X, minWidth)
		cs.Max.X = min(cs.Max.X, maxWidth)*/
	case PreferredWidth:		// SizeHint is Preferred
		log.Println("PreferredWidth............")
		log.Println("object MinWidth = (", minWidth, " )")
		log.Println("object MaxWidth = (", maxWidth, " )")
		cs.Min.X = max(cs.Min.X, minWidth)
		cs.Max.X = min(cs.Max.X, maxWidth)
	/*case MaximumWidth:			// SizeHint is Maximum
		cs.Min.X = max(cs.Min.X, minWidth) 	// No change to gtx.Constraints.X
		cs.Max.X = min(cs.Max.X, maxWidth)*/
	case ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		cs.Min.Y = min(cs.Max.Y, height)
		cs.Max.Y = min(cs.Max.Y, height)
	/*case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = min(cs.Min.Y, ob.MinHeight)
		cs.Max.Y = min(cs.Max.Y, ob.MaxHeight)*/
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = min(cs.Min.Y, minHeight)
		cs.Max.Y = min(cs.Max.Y, maxHeight)
	/*case MaximumHeight:			// SizeHint is Maximum
		cs.Min.Y = min(cs.Min.Y, ob.MinHeight) 	// No change to gtx.Constraints.Y
		cs.Max.Y = min(cs.Max.Y, ob.MaxHeight)*/
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