/* object.go */

package utopia

import (
	//"log"

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
	MinimuHeight GoSizeType		= 0x0004
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
	addControl(GoObject)
	objects() ([]GoObject)
	draw(layout_gio.Context) (layout_gio.Dimensions)
	objectType() (string)
	parentControl() (GoObject)
	parentWindow() (*GoWindowObj)
	removeControl(GoObject)
	sizePolicy() *GoSizePolicy
	setSizePolicy(*GoSizePolicy)
	update(layout_gio.Context)
	//repack()
	//wid() (*goWidget)
}

type goObject struct {
	parent GoObject
	window *GoWindowObj
	//controls  map[int]GoObject
	controls []GoObject
	goSizePolicy *GoSizePolicy
}

func (ob *goObject) addControl(control GoObject) {
	//ob.controls[id] = control
	ob.controls = append(ob.controls, control)
	//ob.index = append(object.index, id)
}

func (ob *goObject) objects() []GoObject {
	return ob.controls
}

func (ob *goObject) draw(layout_gio.Context) (layout_gio.Dimensions) {
	return layout_gio.Dimensions{}
}

func (ob *goObject) objectType() (string) {
	return ""
}

func (ob *goObject) parentControl() (GoObject) {
	return ob.parent
}

func (ob *goObject) parentWindow() (*GoWindowObj) {
	return ob.window
}

func (ob *goObject) removeControl(object GoObject) {
	k := 0
	for _, v := range ob.controls {
	    if v != object {
	        ob.controls[k] = v
	        k++
	    }
	}
	ob.controls = ob.controls[:k] // set slice len to remaining elements
}

func (ob *goObject) sizePolicy() *GoSizePolicy {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	return ob.goSizePolicy
}

func (ob *goObject) setSizePolicy(sizePolicy *GoSizePolicy) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	ob.goSizePolicy = sizePolicy
}
/*func (ob *goObject) repack() {
	
}*/

func (ob *goObject) update(gtx layout_gio.Context) {
	//if ob.objectType()
}