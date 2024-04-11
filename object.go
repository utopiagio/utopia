// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/object.go */

package utopia

import (
	//"log"

	layout_gio "github.com/utopiagio/gio/layout"
	"github.com/utopiagio/utopia/metrics"
)

type GoSizePolicy struct {
	Horiz 	GoSizeType
	Vert 	GoSizeType
	HFlex 	bool
	VFlex 	bool
}

func GetSizePolicy(horiz GoSizeType, vert GoSizeType) (*GoSizePolicy) {
	var flexWidth bool = false
	var flexHeight bool = false
	if horiz == ExpandingWidth {flexWidth = true}
	if vert == ExpandingHeight {flexHeight = true}
	return &GoSizePolicy{horiz, vert, flexWidth, flexHeight}
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
	for k, _ := range ob.Controls {
	   	ob.Controls[k] = nil
	}
	ob.Controls = []GoObject{}
}

func (ob *GioObject) DeleteControl(control GoObject) {
	k := 0
	for _, v := range ob.Controls {
	    if v != control {
	        ob.Controls[k] = v
	        k++
	    } else {
	    	ob.Controls[k] = nil
	    }
	}
	ob.Controls = ob.Controls[:k] // set slice len to remaining elements
}


func (ob *GioObject) Draw(layout_gio.Context) (layout_gio.Dimensions) {
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

func (ob *GioObject) Objects() (controls []GoObject) {
	return ob.Controls
}

/*func (ob *GioObject) ObjectType() (string) {
	Implemented in each Utopia control.
}*/

func (ob *GioObject) ParentControl() (control GoObject) {
	return ob.Parent
}

func (ob *GioObject) ParentWindow() (window *GoWindowObj) {
	return ob.Window
}

/*func (ob *GioObject) Widget() *GioWidget {
	Implemented in each Utopia control.
}*/

func (ob *GioObject) RemoveControl(control GoObject) {
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

func (ob *GioObject) SetConstraints(size GoSize, cs layout_gio.Constraints) (layout_gio.Constraints) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	
	width := metrics.DpToPx(GoDpr, size.Width)
	height := metrics.DpToPx(GoDpr, size.Height)
	minWidth := metrics.DpToPx(GoDpr, size.MinWidth)
	minHeight := metrics.DpToPx(GoDpr, size.MinHeight)
	maxWidth := metrics.DpToPx(GoDpr, size.MaxWidth)
	maxHeight := metrics.DpToPx(GoDpr, size.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case FixedWidth:			// SizeHint is Fixed
		cs.Min.X = min(cs.Max.X, width)			// constrain ob.Width to cs.Max.X 
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = min(cs.Max.X, minWidth)		// constrain ob.MinWidth to cs.Max.X
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = minWidth						// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain ob.MaxWidth to cs.Max.X
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	case ExpandingWidth:		// SizeHint is Expanding
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain ob.MaxWidth to cs.Max.X
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case FixedHeight:			// SizeHint is Fixed 
		cs.Min.Y = min(cs.Max.Y, height)		// constrain to cs.Max.Y 
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = min(cs.Max.Y, minHeight)		// set to ob.MinHeight
		cs.Max.Y = cs.Min.Y						// set to ob.MinHeight
	case PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = minHeight					// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain ob.Height to cs.Max.Y
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	case ExpandingHeight:		// SizeHint is Expanding
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain ob.MaxHeight to cs.Max.Y
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}

	return cs
}

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