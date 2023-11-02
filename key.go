/* key.go */

package utopia

import (
	//"errors"
	"log"
	//"time"
	//"os"
)

func GoKeyboard() (hObj *GoKeyboardObj) {
	return &GoKeyboardObj{}
}

type GoKeyboardObj struct {
	controlFocus *GioWidget
	focus bool
	focusControl *GioWidget
}

func (k *GoKeyboardObj) ClearFocus(w *GioWidget) bool {
	//log.Println("GoKeyboardObj::ClearFocus")
	k.focus = false
	k.controlFocus = nil
	return true
}

func (k *GoKeyboardObj) HasFocus() (focus bool) {
	return k.focus
}

func (k *GoKeyboardObj) GetFocus() (w *GioWidget) {
	return k.controlFocus
}

func (k *GoKeyboardObj) SetFocus(w *GioWidget) bool {
	if k.controlFocus == w {
		log.Println("GoKeyboardObj::ChangeFocus return true......")
		return true
	}
	if k.controlFocus != nil {
		if !k.controlFocus.ClearFocus() {
			log.Println("GoKeyboardObj::ClearFocus return false......")
			return false
		}
	}

	log.Println("GoKeyboardObj::ChangeFocus")
	
	if w == nil {
		k.controlFocus = w
		k.focus = false
	} else {
		if w.ChangeFocus(true) {
			log.Println("GoKeyboardObj::ChangeFocus return true......")
			k.controlFocus = w
			k.focus = true
		} else {
			log.Println("GoKeyboardObj::ChangeFocus return false......")
			return false
		}
	}
	return true
	
}

func (k *GoKeyboardObj) SetFocusControl(w *GioWidget) {
	k.focusControl = w
}

func (k *GoKeyboardObj) Update() (ok bool) {
	ok = true
	if k.focusControl != k.controlFocus {
		ok = k.SetFocus(k.focusControl)
	}
	return ok
}