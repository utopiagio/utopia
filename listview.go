// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/listview.go */

package utopia

import (
	"log"
	"image"
	//"image/color"
	"math"

	//pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	//unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"

	"github.com/utopiagio/utopia/metrics"
)

// ListStyle configures the presentation of a layout_gio.List with a scrollbar.
type GoListViewObj struct {
	GioObject
	GioWidget
	Scrollbar goScrollbar
	AnchorStrategy
	state *widget_gio.List
	columns int
	itemColor GoColor
	itemSize int
	//itemList []*GoListViewItemObj
	currentItem *GoListViewItemObj
	layout *GoLayoutObj

	onItemClicked func([]int)
	onItemDoubleClicked func([]int)
}

// List constructs a ListStyle using the provided theme and state.
func GoListView(parent GoObject) *GoListViewObj {
	theme := GoApp.Theme()
	lightFg := theme.ColorFg.NRGBA()
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	state := &widget_gio.List{}
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{4,4,4,4},
		GoSize: GoSize{100, 100, 100, 100, 1000, 1000, 100, 100},
		Visible: true,
		//target: nil,
	}
	scrollbar := goScrollbar{
		Scrollbar: &state.Scrollbar,
		Track: ScrollTrackStyle{
			MajorPadding: 2,
			MinorPadding: 2,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  8,
			MinorWidth:   6,
			CornerRadius: 3,
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
	hListView := &GoListViewObj {
		GioObject: object,
		GioWidget: widget,
		Scrollbar: scrollbar,
		//ScrollbarStyle: Scrollbar(theme, state.Scrollbar),
		AnchorStrategy: Occupy,
		state: 	state,
		columns: 1,
		itemSize: 24,
		itemColor: Color_Black,
		layout: nil,
	}
	hListView.layout = GoVFlexBoxLayout(hListView)
	parent.AddControl(hListView)
	return hListView
}

func (ob *GoListViewObj) AddListItem(iconData []byte, labelText string) (listItem *GoListViewItemObj) {
	log.Println("GoListViewObj::AddListItem(", labelText, ")")
	listItem = GoListViewItem(ob, iconData, labelText, 0, len(ob.Controls))
	listItem.SetIconSize(ob.itemSize)
	listItem.SetIconColor(ob.itemColor)
	ob.AddControl(listItem)
	return listItem
}

func (ob *GoListViewObj) CurrentSelection() (*GoListViewItemObj) {
	return ob.currentItem
}

func (ob *GoListViewObj) InsertListItem(iconData []byte, labelText string, idx int) (listItem *GoListViewItemObj) {
	log.Println("GoListViewObj::InsertListItem(", labelText, ")")
	listItem = GoListViewItem(ob, iconData, labelText, 0, len(ob.Controls))
	listItem.SetIconSize(ob.itemSize)
	listItem.SetIconColor(ob.itemColor)
	ob.InsertControl(listItem, idx)
	return listItem
}

func (ob *GoListViewObj) ClearList() {
	log.Println("GoListViewObj::ClearList()")
	ob.currentItem = nil
	ob.Clear()
}

func (ob *GoListViewObj) Item(nodeId []int) (*GoListViewItemObj) {
	var listViewItem  *GoListViewItemObj
	for level := 0; level < len(nodeId); level++ {
		listViewItem = listViewItem.Objects()[nodeId[level]].(*GoListViewItemObj)
	}
	return listViewItem
}

/*func (ob *GoListViewObj) Item(id int) (*GoListViewItemObj) {
	//return ob.itemList[id]
	return ob.Objects()[id].(*GoListViewItemObj)
}*/

func (ob *GoListViewObj) ItemClicked(nodeId []int) {
	ob.switchFocus()
	if ob.onItemClicked != nil {
		ob.onItemClicked(nodeId)
	}
}

func (ob *GoListViewObj) ItemDoubleClicked(nodeId []int) {
	ob.switchFocus()
	log.Println("GoListViewObj.ItemDoubleClicked()")
	//log.Println("Parent GoListViewObj.nodeId:", nodeId)
	if ob.onItemDoubleClicked != nil {
		ob.onItemDoubleClicked(nodeId)
	}
}

func (ob *GoListViewObj) RemoveListItem(item GoObject) {
	log.Println("GoListViewObj::RemoveListItem()")
	ob.RemoveControl(item)
}

func (ob *GoListViewObj) ObjectType() (string) {
	return "GoListViewObj"
}

func (ob *GoListViewObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

func (ob *GoListViewObj) SetLayoutMode(mode GoLayoutDirection) {
	ob.state.Axis = layout_gio.Axis(uint8(mode))
}

func (ob *GoListViewObj) SetOnItemClicked(f func([]int)) {
	ob.onItemClicked = f
}

func (ob *GoListViewObj) SetOnItemDoubleClicked(f func([]int)) {
	ob.onItemDoubleClicked = f
}

func (ob *GoListViewObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	log.Println("GoListViewObj::Draw()")
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
		log.Println("ExpandingWidth............")
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
					//return ob.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx, len(ob.Controls), func(gtx C, i int) D {
						return ob.Controls[i].Draw(gtx)
					})
				})
			})
		})
		ob.dims = dims
		ob.AbsWidth = metrics.PxToDp(GoDpr, dims.Size.X)
		ob.AbsHeight = metrics.PxToDp(GoDpr, dims.Size.Y)
	}
	return dims
}

/*func (ob *GoListViewObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx)
	width := gtx.Dp(unit_gio.Dp(ob.Width))
	height := gtx.Dp(unit_gio.Dp(ob.Height))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	dims := image.Point{X: width, Y: height}
	rr := gtx.Dp(ob.cornerRadius)
	defer clip_gio.UniformRRect(image.Rectangle{Max: dims}, rr).Push(gtx.Ops).Pop()
	// paint background
	background := ob.background.NRGBA()
	paint_gio.Fill(gtx.Ops, background)
	// add the events handler to receive widget pointer events
	ob.SignalEvents(gtx)
	ob.layout.Draw(gtx)
	return layout_gio.Dimensions{Size: dims}
}*/

// layout the list and its scrollbar.
func (ob *GoListViewObj) Layout(gtx layout_gio.Context, length int, w layout_gio.ListElement) layout_gio.Dimensions {
	originalConstraints := gtx.Constraints

	// Determine how much space the scrollbar occupies.
	barWidth := gtx.Dp(ob.Scrollbar.Width())

	if ob.AnchorStrategy == Occupy {

		// Reserve space for the scrollbar using the gtx constraints.
		max := ob.state.Axis.Convert(gtx.Constraints.Max)
		min := ob.state.Axis.Convert(gtx.Constraints.Min)
		max.Y -= barWidth
		if max.Y < 0 {
			max.Y = 0
		}
		min.Y -= barWidth
		if min.Y < 0 {
			min.Y = 0
		}
		gtx.Constraints.Max = ob.state.Axis.Convert(max)
		gtx.Constraints.Min = ob.state.Axis.Convert(min)
	}

	listDims := ob.state.List.Layout(gtx, length, w)
	gtx.Constraints = originalConstraints

	// Draw the scrollbar.
	anchoring := layout_gio.E
	if ob.state.Axis == layout_gio.Horizontal {
		anchoring = layout_gio.S
	}
	majorAxisSize := ob.state.Axis.Convert(listDims.Size).X
	start, end := fromListPosition(ob.state.Position, length, majorAxisSize)
	// layout.Direction respects the minimum, so ensure that the
	// scrollbar will be drawn on the correct edge even if the provided
	// layout.Context had a zero minimum constraint.
	gtx.Constraints.Min = listDims.Size
	if ob.AnchorStrategy == Occupy {
		min := ob.state.Axis.Convert(gtx.Constraints.Min)
		min.Y += barWidth
		gtx.Constraints.Min = ob.state.Axis.Convert(min)
	}
	anchoring.Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
		return ob.Scrollbar.Layout(gtx, ob.state.Axis, start, end)
	})

	if delta := ob.state.ScrollDistance(); delta != 0 {
		// Handle any changes to the list position as a result of user interaction
		// with the scrollbar.
		ob.state.List.Position.Offset += int(math.Round(float64(float32(ob.state.Position.Length) * delta)))

		// Ensure that the list pays attention to the Offset field when the scrollbar drag
		// is started while the bar is at the end of the list. Without this, the scrollbar
		// cannot be dragged away from the end.
		ob.state.List.Position.BeforeEnd = true
	}

	if ob.AnchorStrategy == Occupy {
		// Increase the width to account for the space occupied by the scrollbar.
		cross := ob.state.Axis.Convert(listDims.Size)
		cross.Y += barWidth
		listDims.Size = ob.state.Axis.Convert(cross)
	}
	//log.Println("listDims :", listDims)
	if ob.MinWidth > listDims.Size.X {
		listDims.Size.X = ob.MinWidth
	}
	return listDims
}

func (ob *GoListViewObj) switchFocus() {
	log.Println("GoListViewObj::switchFocus()")
	if ob.currentItem != nil {
		log.Println("ob.currentItem.SetSelected(false)")
		log.Println("ob.currentItem -", ob.currentItem.Text())
		ob.currentItem.SetSelected(false)
		ob.currentItem.ClearHighlight()	
	}	
	for _, item := range ob.Controls {
		if item.(*GoListViewItemObj).HasFocus() {
			ob.currentItem = item.(*GoListViewItemObj)
			ob.currentItem.SetSelected(true)
		}
	}
}