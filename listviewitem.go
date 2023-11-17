package utopia

import (
	//"log"
	//"image"
	//"image/color"
	//"math"

	pointer_gio "github.com/utopiagio/gio/io/pointer"
	layout_gio "github.com/utopiagio/gio/layout"
	//op_gio "github.com/utopiagio/gio/op"
	//clip_gio "github.com/utopiagio/gio/op/clip"
	//paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"
	//widget_gio "github.com/utopiagio/gio/widget"
)

func GoListViewItem(parent GoObject, text string, menuId int, action func()) (hObj *GoListViewItemObj) {
	//var fontSize unit_gio.Sp = 14
	var theme *GoThemeObj = GoApp.Theme()

	object := GioObject{parent, parent.ParentWindow(), []GoObject{}, GetSizePolicy(FixedWidth, FixedHeight)}
	widget := GioWidget{
		GoBorder: GoBorder{BorderNone, Color_Black, 0, 0, 0},
		GoMargin: GoMargin{0,0,0,0},
		GoPadding: GoPadding{0,0,0,0},
		GoSize: GoSize{100, 30, 100, 30, 1000, 30},
		FocusPolicy: StrongFocus,
		Visible: true,
	}
	hListViewItem := &GoListViewItemObj{
		GioObject: object,
		GioWidget: widget,

		menuId: menuId,
		fontSize: theme.TextSize,
		text: text,
		color: theme.TextColor,
		background: Color_WhiteSmoke,
		cornerRadius: 0,
		inset: layout_gio.Inset{
			Top: 4, Bottom: 4,
			Left: 28, Right: 28,
		},
		shaper: theme.Shaper,
	}
	hListViewItem.SetSizePolicy(FixedWidth, FixedHeight)
	hListViewItem.SetOnPointerRelease(hListViewItem.Click)
	hListViewItem.SetOnPointerEnter(nil)
	hListViewItem.SetOnPointerLeave(nil)
	//hListViewItem.SetOnClick(action)
	parent.AddControl(hListViewItem)
	return hListViewItem
}

type GoListViewItemObj struct {
	GioObject
	GioWidget
	//theme *GoThemeObj
	font text_gio.Font
	fontSize unit_gio.Sp
	menuId int
	text string
	color GoColor
	background GoColor
	cornerRadius unit_gio.Dp
	inset layout_gio.Inset
	shaper *text_gio.Shaper
	onClick func()

	menuItems []*GoMenuItemObj
	itemLength int
	//textAlign text.Alignment
}

func (ob *GoListViewItemObj) Click(e pointer_gio.Event) {
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

func (ob *GoListViewItemObj) ObjectType() (string) {
	return "GoListViewItemObj"
}

func (ob *GoListViewItemObj) Widget() (*GioWidget) {
	return &ob.GioWidget
}