/* layout.go */

package utopia

import (
	//"log"
	//"image/color"

	layout_gio "github.com/utopiagio/gio/layout"
	//widget_gio "github.com/utopiagio/gio/widget"
	//"github.com/utopiagio/gio/unit"
)

type GoLayoutStyle int

const (
	NoLayout 	GoLayoutStyle = iota
	HBoxLayout 							// gio.List{Axis: layout_gio.Horizontal}	
	VBoxLayout 							// gio.List{Axis: layout_gio.Vertical}
	HVBoxLayout
	// gio.Flex{Axis: layout_gio.Horizontal, Spacing: 0, Alignment: Baseline, WeightSum: 0}
	HFlexBoxLayout							
	// gio.Flex{Axis: layout_gio.Vertical, Spacing: 0, Alignment: Baseline, WeightSum: 0}
	VFlexBoxLayout							
)

type GoLayoutDirection int

const (
	Horizontal = 0
	Vertical = 1
)

type GoLayoutSpacing uint8

const (
	// SpaceEnd leaves space at the end.
	SpaceEnd GoLayoutSpacing = iota
	// SpaceStart leaves space at the start.
	SpaceStart
	// SpaceSides shares space between the start and end.
	SpaceSides
	// SpaceAround distributes space evenly between children,
	// with half as much space at the start and end.
	SpaceAround
	// SpaceBetween distributes space evenly between children,
	// leaving no space at the start and end.
	SpaceBetween
	// SpaceEvenly distributes space evenly between children and
	// at the start and end.
	SpaceEvenly
)

type GoLayoutAlignment uint8

const (
	LayoutStart GoLayoutAlignment = iota
	LayoutEnd
	LayoutMiddle
	LayoutBaseline
)

func GoLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	switch style {
		case NoLayout:
			return GoBoxLayout(parent, NoLayout)
		case HBoxLayout:
			return GoHBoxLayout(parent)
		case VBoxLayout:
			return GoVBoxLayout(parent)
		case HFlexBoxLayout:
			return GoHFlexBoxLayout(parent)
		case VFlexBoxLayout:
			return GoVFlexBoxLayout(parent)
		default:
			return GoVBoxLayout(parent)
	}
}

func GoBoxLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	var axis layout_gio.Axis
	switch style {
		case NoLayout:
			axis = layout_gio.Horizontal
		case HBoxLayout:
			axis = layout_gio.Horizontal
		case VBoxLayout:
			axis = layout_gio.Vertical
	}
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		list_gio: &layout_gio.List{Axis: axis},
		style: style,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

func GoHBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		list_gio: &layout_gio.List{Axis: layout_gio.Horizontal},
		style: HBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

func GoVBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		list_gio: &layout_gio.List{Axis: layout_gio.Vertical},
		style: VBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

func GoFlexBoxLayout(parent GoObject, style GoLayoutStyle) (hObj *GoLayoutObj) {
	var axis layout_gio.Axis
	switch style {
		case NoLayout:
			axis = layout_gio.Horizontal
		case HFlexBoxLayout:
			axis = layout_gio.Horizontal
		case VFlexBoxLayout:
			axis = layout_gio.Vertical
	}
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: axis},
		style: style,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

func GoHFlexBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: layout_gio.Horizontal},
		style: HFlexBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

func GoVFlexBoxLayout(parent GoObject) (hObj *GoLayoutObj) {
	object := goObject{parent, parent.parentWindow(), []GoObject{}, GetSizePolicy(ExpandingWidth, ExpandingHeight)}
	//widget := goWidget{goMargin{0,0,0,0}, goPadding{0,0,0,0}, layout_gio.Inset{0,0,0,0}, layout_gio.Inset{0,0,0,0}}
	widget := goWidget{
		goBorder: goBorder{BorderNone, Color_Black, 0, 0},
		goMargin: goMargin{0,0,0,0},
		goPadding: goPadding{0,0,0,0},
		visible: true,
	}
	hLayout := &GoLayoutObj{
		goObject: object,
		goWidget: widget,
		flex_gio: &layout_gio.Flex{Axis: layout_gio.Vertical},
		style: VFlexBoxLayout,
		flexControls: []layout_gio.FlexChild{},
	}
	parent.addControl(hLayout)
	return hLayout
}

type GoLayoutObj struct {
	goObject
	goWidget
	list_gio 	*layout_gio.List
	flex_gio 	*layout_gio.Flex
	style 		GoLayoutStyle
	flexControls 	[]layout_gio.FlexChild
}

func (ob *GoLayoutObj) SetAlignment(alignment GoLayoutAlignment) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flex_gio.Alignment = layout_gio.Alignment(uint8(alignment))	// layout_gio.Alignment
	} else if ob.style == HBoxLayout || ob.style == VBoxLayout {
		ob.list_gio.Alignment = layout_gio.Alignment(uint8(alignment))	// layout_gio.Alignment
	}
}
func (ob *GoLayoutObj) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {
	ob.setSizePolicy(GetSizePolicy(horiz, vert))
}

func (ob *GoLayoutObj) SetSpacing(spacing GoLayoutSpacing) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flex_gio.Spacing = layout_gio.Spacing(uint8(spacing))	// layout_gio.Spacing
	}
}

func (ob *GoLayoutObj) addFlexControl(control GoObject) {
	if ob.style == HFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.FlexControl(control.sizePolicy().HFlex, control.sizePolicy().VFlex, 1, control.draw))
	} else if ob.style == VFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.FlexControl(control.sizePolicy().VFlex, control.sizePolicy().HFlex, 1, control.draw))
	}
}

func (ob *GoLayoutObj) addFlexedControl(control GoObject) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Flexed(1, control.draw))
	}
}

func (ob *GoLayoutObj) addRigidControl(control GoObject) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Rigid(control.draw))
	}
}

/*func (ob *GoLayoutObj) AddFlexLayout(direction GoLayoutDirection) (layout *GoLayoutObj) {
	if direction == Horizontal {
		return ob.addFlexedLayout(HFlexBoxLayout)
	} else {
		return ob.addFlexedLayout(VFlexBoxLayout)
	}
}*/

/*func (ob *GoLayoutObj) addFlexedLayout(style GoLayoutStyle) (layout *GoLayoutObj) {
	layout = GoFlexBoxLayout(ob, style)
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Flexed(1, layout.draw))
	}
	//ob.goObject.addControl(layout)
	return layout
}*/

/*func (ob *GoLayoutObj) addRigidLayout(layout *GoLayoutObj) {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = append(ob.flexControls, layout_gio.Rigid(layout.draw))
	}
	ob.goObject.addControl(layout)
}*/

func (ob *GoLayoutObj) draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions{Size: gtx.Constraints.Max,}
	if ob.visible {
		if ob.style == HBoxLayout || ob.style == VBoxLayout {
			//log.Println("BoxLayout style:", ob.style)
			dims = ob.goMargin.layout(gtx, func(gtx C) D {
				return ob.goBorder.layout(gtx, func(gtx C) D {
					return ob.goPadding.layout(gtx, func(gtx C) D {
						return ob.list_gio.Layout(gtx, len(ob.controls), func(gtx C, i int) D {
							return ob.controls[i].draw(gtx)
						})
					})
				})
			})
		} else if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
			//log.Println("FlexBoxLayout style:", ob.style)
			ob.repack()
			dims = ob.goMargin.layout(gtx, func(gtx C) D {
				return ob.goBorder.layout(gtx, func(gtx C) D {
					return ob.goPadding.layout(gtx, func(gtx C) D {
						return ob.flex_gio.Layout(gtx, ob.flexControls... ) 
					})
				})
			})
		}
	}
	return dims
}

func (ob *GoLayoutObj) objectType() (string) {
	return "GoLayoutObj"
}

func (ob *GoLayoutObj) repack() {
	if ob.style == HFlexBoxLayout || ob.style == VFlexBoxLayout {
		ob.flexControls = []layout_gio.FlexChild{}
		for i := 0; i < len(ob.controls); i++ {
			
			ob.addFlexControl(ob.controls[i])
			/*if ob.controls[i].sizePolicy().HFixed {
				ob.addRigidControl(ob.controls[i])
			} else {
				ob.addFlexedControl(ob.controls[i])
			}*/
			if ob.controls[i].objectType() == "GOLayoutObj" {
				ob.controls[i].(*GoLayoutObj).repack()
			}
		}
	}
}
