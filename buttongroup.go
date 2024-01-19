// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/buttongroup.go */

package utopia

import (
	_ "log"

	widget_int "github.com/utopiagio/utopia/internal/widget"

	"github.com/utopiagio/gio/layout"
)

func GoButtonGroup() (*GoButtonGroupObj) {
	enum := &widget_int.GioEnum{}
	return &GoButtonGroupObj{enum}
}

type GoButtonGroupObj struct {
	enum *widget_int.GioEnum
}

/*type GioEnum struct {
	Value    string
	hovered  string
	hovering bool

	focus   string
	focused bool

	keys []*enumKey
}*/

/*type enumKey struct {
	key   string
	click gesture.Click
	tag   struct{}
}*/

/*func (ob *GoButtonGroup) index(k string) *enumKey {
	return ob.enum.index(k)
}*/

func (ob *GoButtonGroupObj) Value() string {
	return ob.enum.Value
}

// Value has changed by user interaction.
func (ob *GoButtonGroupObj) Update(gtx layout.Context) bool {
	return ob.enum.Update(gtx)
}

// Hovered returns the key that is highlighted, or false if none are.
func (ob *GoButtonGroupObj) Hovered() (string, bool) {
	return ob.enum.Hovered()
}

// Focused reports the focused key, or false if no key is focused.
func (ob *GoButtonGroupObj) Focused() (string, bool) {
	return ob.enum.Focused()
}

// Layout adds the event handler for the key k.
func (ob *GoButtonGroupObj)Layout(gtx layout.Context, k string, content layout.Widget) layout.Dimensions {
	return ob.enum.Layout(gtx, k, content)
}