// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/object.go */

package utopia

import (
	"log"

	layout_gio "github.com/utopiagio/gio/layout"
)

func GetSizePolicy(horiz GoSizeType, vert GoSizeType) (*GoSizePolicy) {
	var fWidth bool = false
	var fHeight bool = false
	if horiz == ExpandingWidth {fWidth = true}
	if vert == ExpandingHeight {fHeight = true}
	return &GoSizePolicy{horiz, vert, fWidth, fHeight}
}

func ExpandingSizePolicy() (*GoSizePolicy) {
	return &GoSizePolicy{ExpandingWidth, ExpandingHeight, true, true}
}

func FixedSizePolicy() (*GoSizePolicy) {
	return &GoSizePolicy{FixedWidth, FixedHeight, false, false}
}

type GoSizeType int

const(
	FixedWidth GoSizeType		= 0x0000
	FixedHeight GoSizeType		= 0x0001
	MinimumWidth GoSizeType		= 0x0002
	MinimumHeight GoSizeType	= 0x0004
	MaximumWidth GoSizeType		= 0x0008
	MaximumHeight GoSizeType	= 0x0010
	PreferredWidth GoSizeType	= 0x0020
	PreferredHeight GoSizeType	= 0x0040
	ExpandingWidth GoSizeType	= 0x0080
	ExpandingHeight GoSizeType 	= 0x0100
)

type GoSizePolicy struct {
	Horiz 	GoSizeType
	Vert 	GoSizeType
	HFlex 	bool
	VFlex 	bool
}


type GoObject interface {
	AddControl(GoObject)
	Clear()
	Objects() ([]GoObject)
	DeleteControl(GoObject)
	Draw(layout_gio.Context) (layout_gio.Dimensions)
	InsertControl(GoObject, int)
	ObjectType() (string)
	ParentControl() (GoObject)
	ParentWindow() (*GoWindowObj)
	RemoveControl(GoObject)
	SizePolicy() (*GoSizePolicy)
	SetHorizSizePolicy(horiz GoSizeType)
	SetSizePolicy(horiz GoSizeType, vert GoSizeType)
	SetVertSizePolicy(vert GoSizeType)
	Widget() (*GioWidget)
}

type GioObject struct {
	Parent GoObject
	Window *GoWindowObj
	Controls []GoObject
	GoSizePolicy *GoSizePolicy
}

func (ob *GioObject) AddControl(control GoObject) {
	ob.Controls = append(ob.Controls, control)
}

func (ob *GioObject) Clear() {
	ob.Controls = []GoObject{}
}

func (ob *GioObject) DeleteControl(object GoObject) {
	k := 0
	for _, v := range ob.Controls {
	    if v != object {
	        ob.Controls[k] = v
	        k++
	    } else {
	    	ob.Controls[k] = nil
	    }
	}
	ob.Controls = ob.Controls[:k] // set slice len to remaining elements
}


func (ob *GioObject) Draw(layout_gio.Context) (layout_gio.Dimensions) {
	log.Println("GioObject.Draw()")
	return layout_gio.Dimensions{}
}

func (ob *GioObject) InsertControl(control GoObject, idx int) {
	if len(ob.Controls) < 1 || idx >= len(ob.Controls) {
		ob.Controls = append(ob.Controls, control)
	} else {
		ob.Controls = append(ob.Controls[:idx + 1], ob.Controls[idx:]...)
		ob.Controls[idx] = control
	}
}

func (ob *GioObject) Objects() []GoObject {
	return ob.Controls
}

/*func (ob *GioObject) ObjectType() (string) {
	log.Println("GioObject.ObjectType() -", ob)
	return ""
}*/

func (ob *GioObject) ParentControl() (GoObject) {
	return ob.Parent
}

func (ob *GioObject) ParentWindow() (*GoWindowObj) {
	return ob.Window
}

/*func (ob *GioObject) Widget() *GioWidget {
	return nil
}*/

func (ob *GioObject) RemoveControl(object GoObject) {
	k := 0
	for _, v := range ob.Controls {
	    if v != object {
	        ob.Controls[k] = v
	        k++
	    }
	}
	ob.Controls = ob.Controls[:k] // set slice len to remaining elements
}

/*func (ob *GioObject) RemoveIndex(idx int) {
	if idx >= 0 || idx < len(ob.Controls) {
		if len(ob.Controls) <= 1 || idx >= len(ob.Controls) {
			ob.Controls = ob.Controls[:idx]
		} else {
			ob.Controls = append(ob.Controls[:idx], ob.Controls[idx + 1:]...) // set slice len to remaining elements
		}
	}
}*/

func (ob *GioObject) SizePolicy() *GoSizePolicy {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	return ob.GoSizePolicy
}

func (ob *GioObject) SetHorizSizePolicy(horiz GoSizeType) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	ob.GoSizePolicy = GetSizePolicy(horiz, ob.GoSizePolicy.Vert)
}

func (ob *GioObject) SetSizePolicy(horiz GoSizeType, vert GoSizeType) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	ob.GoSizePolicy = GetSizePolicy(horiz, vert)
}

func (ob *GioObject) SetVertSizePolicy(vert GoSizeType) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	ob.GoSizePolicy = GetSizePolicy(ob.GoSizePolicy.Horiz, vert)
}