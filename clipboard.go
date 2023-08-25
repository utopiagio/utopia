/* clipboard.go */

package utopia

import (
	"errors"
	"golang.design/x/clipboard"
)

func GoClipBoard() (hObj *GoClipBoardObj) {
	return &GoClipBoardObj{}
}

type GoClipBoardObj struct {
	available bool
}

func (clpbd *GoClipBoardObj) init() (e error) {
	err := clipboard.Init()
	if err != nil {
		clpbd.available = false
	    return err
	}
	clpbd.available = true
	return nil
}

// ReadData returns a byte data buffer from the clipboard
// or nil if the clipboard is empty or does not contain byte data.
func (clpbd *GoClipBoardObj) ReadData() (imdata []byte) {
	if clpbd.available {
		return clipboard.Read(clipboard.FmtImage)
	} else {
		return nil
	} 
}

// ReadText returns a string of text from the clipboard
// or nil if the clipboard is empty or does not contain text.
func (clpbd *GoClipBoardObj) ReadText() (text string) {
	if clpbd.available {
		return string(clipboard.Read(clipboard.FmtText))
	} else {
		return ""
	}
}

func (clpbd *GoClipBoardObj) WriteData(imdata []byte) {
	if clpbd.available {
		clipboard.Write(clipboard.FmtImage, imdata)
	}
}

// WriteText sends a string of text to the clipboard
// or nil if the clipboard is empty or does not contain text.
func (clpbd *GoClipBoardObj) WriteText(text string) (err error) {
	if clpbd.available {
		clipboard.Write(clipboard.FmtText, []byte(text))
		return nil
	}
	return errors.New("Clipboard not available.")
}