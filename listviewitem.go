// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/listviewitem.go */

package utopia

import (
	"log"
	"image"
	"image/draw"
	//"math"

	"github.com/utopiagio/utopia/internal/f32color"
	font_gio "github.com/utopiagio/gio/font"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
	widget_int "github.com/utopiagio/utopia/internal/widget"

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
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil
	}
	
	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{24, 24, 24, 24, 1000, 1000},
		FocusPolicy: StrongFocus,
		Visible: true,

		ForeColor: theme.ColorFg,
		
		//target: nil,
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
		//parentList: parent,
		shaper: theme.Shaper,
	}
	hListViewItem.SetSizePolicy(FixedWidth, FixedHeight)
	hListViewItem.SetOnSetFocus(hListViewItem.SetHighlight)
	hListViewItem.SetOnClearFocus(hListViewItem.ClearHighlight)
	hListViewItem.SetOnPointerClick(hListViewItem.Clicked)
	hListViewItem.SetOnPointerDoubleClick(hListViewItem.DoubleClicked)
	hListViewItem.SetOnPointerEnter(nil)
	hListViewItem.SetOnPointerLeave(nil)
	switch parent.ObjectType() {
	case "GoListViewItemObj":
		hListViewItem.listView = parent.(*GoListViewItemObj).ListView()
		parent.AddControl(hListViewItem)
	case "GoListViewObj":
		hListViewItem.listView = parent.(*GoListViewObj)
	}
	//hListViewItem.listView.AddControl(hListViewItem)
	//hListViewItem.listView.InsertControl(hListViewItem, listId)
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
	//parentList GoObject				// *GoListViewObj or *GoListViewItemObj
	
	
	
	shaper *text_gio.Shaper
	
	//onClick func()
	
	
	// Cached values.
	op       paint_gio.ImageOp
	imgSize  int
	imgColor GoColor
}

func (ob *GoListViewItemObj) AddListItem(iconData []byte, labelText string) (listItem *GoListViewItemObj) {
	log.Println("GoListViewItemObj::AddListItem()")
	listItem = GoListViewItem(ob, iconData, labelText, ob.level + 1, len(ob.Controls))
	listItem.SetMargin(20 * listItem.level, 0, 0, 0)
	listItem.SetIconSize(ob.iconSize)
	listItem.SetIconColor(ob.iconColor)
	ob.listView.AddControl(listItem)
	return listItem
}

func (ob *GoListViewItemObj) InsertListItem(iconData []byte, labelText string, idx int) (listItem *GoListViewItemObj) {
	log.Println("GoListViewItemObj::InsertListItem()")
	listItem = GoListViewItem(ob, iconData, labelText, ob.level + 1, len(ob.Controls))
	listItem.SetMargin(20 * listItem.level, 0, 0, 0)
	listItem.SetIconSize(ob.iconSize)
	listItem.SetIconColor(ob.iconColor)
	ob.listView.InsertControl(listItem, idx)
	return listItem
}

func (ob *GoListViewItemObj) RemoveListItem(item GoObject, idx int) {
	log.Println("GoListViewItemObj::RemoveListItem()")
	ob.RemoveControl(item)
	ob.listView.RemoveListItem(item)
}

func (ob *GoListViewItemObj) Clicked(e pointer_gio.Event) {
	log.Println("GoListViewItemObj.Clicked()-len(ob.Controls):", len(ob.Parent.Objects()))
	switch ob.Parent.ObjectType() {
	case "GoListViewItemObj":
		ob.Parent.(*GoListViewItemObj).ItemClicked([]int{ob.id})
	case "GoListViewObj":
		ob.Parent.(*GoListViewObj).ItemClicked([]int{ob.id})
	}
	
	/*if len(ob.menuItems) > 0 {
		popupMenu := ob.ParentWindow().AddPopupMenu()
		popupMenu.Clear()
		popupMenu.SetMargin(0, 25, 0, 0)
		//log.Println("modal.layout.SetMargin(ob.offset)=", ob.offset())
		offsetX, offsetY := ob.ItemOffset(ob.menuId)
		popupMenu.layout.SetMargin(offsetX, 25 + offsetY, 0, 0)
		for idx := 0; idx < len(ob.menuItems); idx++ {
			popupMenu.layout.AddControl(ob.menuItems[idx])
		}
	
		popupMenu.Show()
	} else {
		ob.ParentWindow().ClearPopupMenus()
		if ob.onClick != nil {
			ob.onClick()
		}
	}
	for idx := 0; idx < len(ob.ParentWindow().menuItems); idx++ {
		for idx := 0; idx < len(ob.menuItems); idx++ {
			ob.ParentWindow().MenuPopup(ob.id).layout.AddControl(ob.menuItems[idx])
	}*/
	
}

func (ob *GoListViewItemObj) DoubleClicked(e pointer_gio.Event) {
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
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	/*X = gtx.Constraints.Max.X
	Y = gtx.Constraints.Max.Y
	if ob.SizePolicy.Horiz == FixedWidth {X = ob.X}
	if ob.SizePolicy.Vert == FixedHeight {Y = ob.Y}*/
	gtx.Constraints.Min = image.Point{ob.MinWidth, 0}
	gtx.Constraints.Max = image.Point{ob.MaxWidth, 5000}
	if ob.Visible {
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx C, ) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.Layout(gtx)
				})
			})
		})
		ob.dims = dims
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
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

func (ob *GoListViewItemObj) SetHighlight() {
	ob.SetBackgroundColor(Color_LightBlue)
	//ob.SetBackgroundColor(ob.Highlight)
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

func (ob *GoListViewItemObj) ObjectType() (string) {
	return "GoListViewItemObj"
}

func (ob *GoListViewItemObj) SetExpanded(state bool) {
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

func (ob *GoListViewItemObj) Text() (string) {
	return ob.label
}

func (ob *GoListViewItemObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}

// Layout displays the icon with its size set to the X minimum constraint.
func (ob *GoListViewItemObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx)
	textColorMacro := op_gio.Record(gtx.Ops)
	paint_gio.ColorOp{Color: ob.color.NRGBA()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()
	dims := layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				gtx.Constraints.Min = image.Point{X: ob.iconSize}
				return ob.layoutIcon(gtx)
			})
		}),
		layout_gio.Rigid(func(gtx layout_gio.Context) layout_gio.Dimensions {
			return layout_gio.UniformInset(2).Layout(gtx, func(gtx layout_gio.Context) layout_gio.Dimensions {
				paint_gio.ColorOp{Color: ob.ForeColor.NRGBA()}.Add(gtx.Ops)
				return widget_int.GioLabel{}.Layout(gtx, ob.shaper, ob.font, ob.fontSize, ob.label, textColor)
			})
		}),
	)
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