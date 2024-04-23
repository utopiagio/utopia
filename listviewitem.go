// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/listviewitem.go */

package utopia

import (
	//"log"
	"image"
	"image/draw"
	//"math"

	"github.com/utopiagio/utopia/internal/f32color"
	font_gio "github.com/utopiagio/gio/font"
	//pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
	widget_int "github.com/utopiagio/utopia/internal/widget"

	"github.com/utopiagio/utopia/metrics"

	"golang.org/x/exp/shiny/iconvg"
)

//const defaultIconColor = Color_Black
//const defaultIconSize = 24

func GoListViewItem(parent GoObject, data []byte, text string, listLevel int, listId int) (hObj *GoListViewItemObj) {
	// GoIcon returns a new Icon from IconVG data.
	var color GoColor
	var size int
	//var listView *GoListViewObj
	var theme *GoThemeObj = GoApp.Theme()
	color = defaultIconColor
	size = defaultIconSize
	//text = ""
	if data != nil {
		_, err := iconvg.DecodeMetadata(data)
		if err != nil {
			return nil
		}
	}
	
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(PreferredWidth, PreferredHeight)}
	tagCounter++
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{24, 24, 24, 24, 1000, 1000, 24, 24},
		FocusPolicy: StrongFocus,
		Visible: false,
		ForeColor: theme.ColorFg,
		tag: tagCounter,
	}
	
	hListViewItem := &GoListViewItemObj{
		GioObject: object,
		GioWidget: widget,
		color: theme.ColorFg,
		expanded: false,
		fontSize: theme.TextSize,
		icon: data,
		iconColor: color,
		iconSize: size,
		listView: nil,
		label: text,
		level: listLevel,
		id: listId,
		maxlines: 1,
		shaper: theme.Shaper,
	}
	hListViewItem.SetOnSetFocus(hListViewItem.SetHighlight)
	hListViewItem.SetOnClearFocus(hListViewItem.ClearHighlight)
	hListViewItem.SetOnPointerClick(hListViewItem.Clicked)
	hListViewItem.SetOnPointerDoubleClick(hListViewItem.DoubleClicked)
	hListViewItem.SetOnPointerEnter(hListViewItem.PointerEnter)
	hListViewItem.SetOnPointerLeave(hListViewItem.PointerLeave)
	switch parent.ObjectType() {
	case "GoListViewItemObj":
		hListViewItem.listView = parent.(*GoListViewItemObj).ListView()
		parent.AddControl(hListViewItem)
	case "GoListViewObj":
		hListViewItem.listView = parent.(*GoListViewObj)
	}
	return hListViewItem
}

type GoListViewItemObj struct {
	GioObject
	GioWidget
	//theme *GoThemeObj
	color GoColor					// foreground color
	expanded bool					// true if the tree node is expanded
	font font_gio.Font				// label font
	fontSize unit_gio.Sp			// label font size
	icon []byte						// icon svg data
	iconColor GoColor				// icon color
	iconSize int					// size of icon determines height of listViewItem
	id int							// position of listViewItem in parentItem
	//itemList []*GoListViewItemObj	// children of this listViewItem
	label string					// text displayed
	level int						// tree level 0 ...
	listView *GoListViewObj			// view to display all listViewItems
	maxlines int
	shaper *text_gio.Shaper
	
	//onClick func()
	
	// Cached values.
	op       paint_gio.ImageOp
	imgSize  int
	imgColor GoColor
}

func (ob *GoListViewItemObj) AddListItem(iconData []byte, labelText string) (listItem *GoListViewItemObj) {
	//log.Println("GoListViewItemObj::AddListItem()")
	listItem = GoListViewItem(ob, iconData, labelText, ob.level + 1, len(ob.Controls))
	listItem.SetMargin(20 * listItem.level, 0, 0, 0)
	listItem.SetIconSize(ob.iconSize)
	listItem.SetIconColor(ob.iconColor)
	if ob.IsExpanded() {
		listItem.Show()
	}
	ob.listView.AddControl(listItem)
	return listItem
}

func (ob *GoListViewItemObj) InsertListItem(iconData []byte, labelText string, idx int) (listItem *GoListViewItemObj) {
	//log.Println("GoListViewItemObj::InsertListItem()")
	listItem = GoListViewItem(ob, iconData, labelText, ob.level + 1, len(ob.Controls))
	listItem.SetMargin(20 * listItem.level, 0, 0, 0)
	listItem.SetIconSize(ob.iconSize)
	listItem.SetIconColor(ob.iconColor)
	if ob.IsExpanded() {
		listItem.Show()
	}
	ob.listView.InsertControl(listItem, idx)
	return listItem
}

func (ob *GoListViewItemObj) RemoveListItem(item GoObject, idx int) {
	//log.Println("GoListViewItemObj::RemoveListItem()")
	ob.RemoveControl(item)
	ob.listView.RemoveListItem(item)
}

func (ob *GoListViewItemObj) Clicked(e GoPointerEvent) {
	//log.Println("GoListViewItemObj.Clicked()-len(ob.Controls):", len(ob.Parent.Objects()))
	//log.Println(ob.Text(), "Clicked()")
	switch ob.Parent.ObjectType() {
		case "GoListViewItemObj":
			ob.Parent.(*GoListViewItemObj).ItemClicked([]int{ob.id})
		case "GoListViewObj":
			ob.Parent.(*GoListViewObj).ItemClicked([]int{ob.id})
	}	
}

func (ob *GoListViewItemObj) DoubleClicked(e GoPointerEvent) {
	//log.Println("GoListViewItemObj.DoubleClicked()-len(ob.Controls):", len(ob.Parent.Objects()))
	//log.Println("GoListViewItemObj.DoubleClicked()-Id:", ob.id)
	switch ob.Parent.ObjectType() {
	case "GoListViewItemObj":
		//log.Println("Parent GoListViewItemObj")
		ob.Parent.(*GoListViewItemObj).ItemDoubleClicked([]int{ob.id})
	case "GoListViewObj":
		//log.Println("Parent GoListViewObj")
		ob.Parent.(*GoListViewObj).ItemDoubleClicked([]int{ob.id})
	}
}

func (ob *GoListViewItemObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	gtx.Constraints = ob.SetConstraints(ob.Size(), gtx.Constraints)
	dims = layout_gio.Dimensions {Size: image.Point{X: 0, Y: 0,}}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C, ) D {
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

func (ob *GoListViewItemObj) ClearHighlight() {
	if ob.IsSelected() {
		ob.SetBackgroundColor(Color_LightGray)
	} else {
		ob.SetBackgroundColor(Color_Transparent)
	}
	//ob.SetBackgroundColor(ob.Highlight)
}



func (ob *GoListViewItemObj) Expand() {
	switch ob.Parent.ObjectType() {
		case "GoListViewItemObj":
			ob.Parent.(*GoListViewItemObj).ItemDoubleClicked([]int{ob.id})
		case "GoListViewObj":
			ob.Parent.(*GoListViewObj).ItemDoubleClicked([]int{ob.id})
	}
}


func (ob *GoListViewItemObj) Icon() ([]byte) {
	return ob.icon
}

func (ob *GoListViewItemObj) IconSize() (int) {
	return ob.iconSize
}

func (ob *GoListViewItemObj) Id() (int) {
	return ob.id
}

func (ob *GoListViewItemObj) IsExpanded() (expanded bool) {
	return ob.expanded
}

func (ob *GoListViewItemObj) Item(id int) (*GoListViewItemObj) {
	//return ob.itemList[id]
	return ob.Objects()[id].(*GoListViewItemObj)
}

func (ob *GoListViewItemObj) ItemCount() (count int) {
	//return len(ob.itemList)
	return len(ob.Objects())
}

func (ob *GoListViewItemObj) ItemClicked(nodeId []int) {
	switch ob.Parent.ObjectType() {
	case "GoListViewItemObj":
		nodeId = append([]int{ob.id}, nodeId...)
		ob.Parent.(*GoListViewItemObj).ItemClicked(nodeId)
	case "GoListViewObj":
		nodeId = append([]int{ob.id}, nodeId...)
		ob.Parent.(*GoListViewObj).ItemClicked(nodeId)
	}
}

func (ob *GoListViewItemObj) ItemDoubleClicked(nodeId []int) {
	//log.Println("GoListViewItemObj.ItemDoubleClicked()")
	
	switch ob.Parent.ObjectType() {
	case "GoListViewItemObj":
		nodeId = append([]int{ob.id}, nodeId...)
		//log.Println("Parent GoListViewItemObj.nodeId:", nodeId)
		ob.Parent.(*GoListViewItemObj).ItemDoubleClicked(nodeId)
	case "GoListViewObj":
		nodeId = append([]int{ob.id}, nodeId...)
		//log.Println("Parent GoListViewObj.nodeId:", nodeId)
		ob.Parent.(*GoListViewObj).ItemDoubleClicked(nodeId)
	}
}

func (ob *GoListViewItemObj) ListView() (*GoListViewObj) {
	return ob.listView
}

func (ob *GoListViewItemObj) MaxLines() (int) {
	return ob.maxlines
}

func (ob *GoListViewItemObj) ObjectType() (string) {
	return "GoListViewItemObj"
}

func (ob *GoListViewItemObj) PointerEnter(e GoPointerEvent) {
	if !ob.HasFocus() && !ob.IsSelected() {
		ob.SetBackgroundColor(Color_WhiteSmoke)
		ob.ParentWindow().Refresh()
	}
}

func (ob *GoListViewItemObj) PointerLeave(e GoPointerEvent) {
	if !ob.HasFocus() && !ob.IsSelected() {
		ob.SetBackgroundColor(Color_Transparent)
		ob.ParentWindow().Refresh()
	}
}

func (ob *GoListViewItemObj) SetExpanded(state bool) {
	//log.Println("GoListViewItemObj::SetExpanded(", state, ")")
	for _, lv := range ob.Controls {
		lvi := lv.(*GoListViewItemObj)
		if state == true {
			lvi.Show()
		} else {
			if lvi.IsExpanded() {
				lvi.SetExpanded(false)
			}
			lvi.Hide()
		}
	}
	ob.expanded = state
}

func (ob *GoListViewItemObj) SetIconColor(color GoColor) {
	ob.iconColor = color
}

func (ob *GoListViewItemObj) SetIconSize(size int) {
	ob.iconSize = size
}

func (ob *GoListViewItemObj) SetId(id int) {
	ob.id = id
}

func (ob *GoListViewItemObj) SetHighlight() {
	ob.SetBackgroundColor(Color_LightBlue)
	//ob.SetBackgroundColor(ob.Highlight)
}

func (ob *GoListViewItemObj) SetMaxLines(maxlines int) {
	ob.maxlines = maxlines
}

func (ob *GoListViewItemObj) Text() (string) {
	return ob.label
}

func (ob *GoListViewItemObj) Trigger() {
	ob.SetFocus()
	//ob.SetSelected(true)
	switch ob.Parent.ObjectType() {
		case "GoListViewItemObj":
			ob.Parent.(*GoListViewItemObj).ItemClicked([]int{ob.id})
		case "GoListViewObj":
			ob.Parent.(*GoListViewObj).ItemClicked([]int{ob.id})
	}
}

func (ob *GoListViewItemObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

// Layout displays the icon with its size set to the X minimum constraint.
func (ob *GoListViewItemObj) Layout(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	ob.ReceiveEvents(gtx, nil)
	textColorMacro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()
	//log.Println("GoListViewItemObj LabelText :", ob.label)
	if ob.icon != nil {
		dims = layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
			layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
				return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
					gtx.Constraints.Min = image.Point{X: ob.iconSize}
					return ob.layoutIcon(gtx)
				})
			}),
			layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
				return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
					paint_gio.ColorOp{Color: ob.ForeColor.NRGBA()}.Add(gtx.Ops)
					return widget_int.GioLabel{MaxLines: ob.maxlines}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.label, textColor)
				})
			}),
		)
	} else {
		dims = layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
			layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
				return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
					paint_gio.ColorOp{Color: ob.ForeColor.NRGBA()}.Add(gtx.Ops)
					return widget_int.GioLabel{MaxLines: ob.maxlines}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.label, textColor)
				})
			}),
		)
	}
	rect := image.Point{X: dims.Size.X, Y: dims.Size.Y}
	defer clip_gio.Rect{Max: rect}.Push(gtx.Ops).Pop()
	ob.SignalEvents(gtx)
	return dims
}

func (ob *GoListViewItemObj) layoutIcon(gtx layout_gio.Context) layout_gio.Dimensions {
	rect := image.Point{X: ob.iconSize, Y: ob.iconSize}
	defer clip_gio.Rect{Max: rect}.Push(gtx.Ops).Pop()

	icon := ob.image(ob.iconSize, ob.iconColor)
	icon.Add(gtx.Ops)
	paint_gio.PaintOp{}.Add(gtx.Ops)

	// add the events handler to receive widget pointer events
	//ob.SignalEvents(gtx)

	return layout_gio.Dimensions{
		Size: icon.Size(),
	}
}

func (ob *GoListViewItemObj) image(sz int, color GoColor) paint_gio.ImageOp {
	if sz == ob.imgSize && color == ob.imgColor {
		return ob.op
	}
	m, _ := iconvg.DecodeMetadata(ob.icon)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = f32color.NRGBAToLinearRGBA(color.NRGBA())
	iconvg.Decode(&ico, ob.icon, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ob.op = paint_gio.NewImageOp(img)
	ob.imgSize = sz
	ob.imgColor = color
	return ob.op
}