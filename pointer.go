// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/pointer.go */

package utopia

import (
	"image"
	"math"
	"time"


	"github.com/utopiagio/gio/f32"
	key_gio "github.com/utopiagio/gio/io/key"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
)


// Buttons is a set of mouse buttons
type Buttons uint8
const (
	// ButtonPrimary is the primary button, usually the left button for a
	// right-handed user.
	ButtonPrimary Buttons = 1 << iota
	// ButtonSecondary is the secondary button, usually the right button for a
	// right-handed user.
	ButtonSecondary
	// ButtonTertiary is the tertiary button, usually the middle button.
	ButtonTertiary
)

// ID is the id for the pointer to track press release 
type ID uint16

// Modifiers
type Modifiers uint32

const (
	// ModCtrl is the ctrl modifier key.
	ModCtrl Modifiers = 1 << iota
	// ModCommand is the command modifier key
	// found on Apple keyboards.
	ModCommand
	// ModShift is the shift modifier key.
	ModShift
	// ModAlt is the alt modifier key, or the option
	// key on Apple keyboards.
	ModAlt
	// ModSuper is the "logo" modifier key, often
	// represented by a Windows logo.
	ModSuper
)

// Priority of an Event.
type Priority uint8
const (
	// Shared priority is for handlers that
	// are part of a matching set larger than 1.
	Shared Priority = iota
	// Foremost priority is like Shared, but the
	// handler is the foremost of the matching set.
	Foremost
	// Grabbed is used for matching sets of size 1.
	Grabbed
)

// Source of an Event.
type Source uint8
const (
	// Mouse generated event.
	Mouse Source = iota
	// Touch generated event.
	Touch
)

// Type of an Event.
type Type uint
const (
	// A Cancel event is generated when the current gesture is
	// interrupted by other handlers or the system.
	Cancel Type = 1 << iota
	// Press of a pointer.
	Press
	// Release of a pointer.
	Release
	// Move of a pointer.
	Move
	// Drag of a pointer.
	Drag
	// Pointer enters an area watching for pointer input
	Enter
	// Pointer leaves an area watching for pointer input
	Leave
	// Scroll of a pointer.
	Scroll
)

type GoPointerEvent struct {
	Type   Type
	Source Source
	// PointerID is the id for the pointer and can be used
	// to track a particular pointer from Press to
	// Release or Cancel.
	PointerID ID
	// Priority is the priority of the receiving handler
	// for this event.
	Priority Priority
	// Time is when the event was received. The
	// timestamp is relative to an undefined base.
	Time time.Duration
	// Buttons are the set of pressed mouse buttons for this event.
	Buttons Buttons
	// Position is the coordinates of the event in the local coordinate
	// system of the receiving tag. The transformation from global window
	// coordinates to local coordinates is performed by the inverse of
	// the effective transformation of the tag.
	Position f32.Point
	// Scroll is the scroll amount, if any.
	Scroll f32.Point
	// Modifiers is the set of active modifiers when
	// the mouse button was pressed.
	Modifiers Modifiers
}



func GioPointerEvent(evt pointer_gio.Event) (ptrEvent GoPointerEvent) {
	pointerEvent := GoPointerEvent{
		Type: Type(evt.Kind),
		Source: Source(evt.Source),
		PointerID: ID(evt.PointerID),
		Priority: Priority(evt.Priority),
		Time: evt.Time,
		Buttons: Buttons(evt.Buttons),
		Position: evt.Position,
		Scroll: evt.Scroll,
		Modifiers: Modifiers(evt.Modifiers),
	}
	return pointerEvent
}

func (ev GoPointerEvent) ButtonState() (Buttons) {
	return ev.Buttons
}

func (evt GoPointerEvent) Gio() (e pointer_gio.Event) {
	e = pointer_gio.Event{
		Kind: pointer_gio.Kind(evt.Type),
		Source: pointer_gio.Source(evt.Source),
		PointerID: pointer_gio.ID(evt.PointerID),
		Priority: pointer_gio.Priority(evt.Priority),
		Time: evt.Time,
		Buttons: pointer_gio.Buttons(evt.Buttons),
		Position: evt.Position,
		Scroll: evt.Scroll,
		Modifiers: key_gio.Modifiers(evt.Modifiers),
	}
	return e
}

func (ev GoPointerEvent) KeyModifiers() (Modifiers) {
	return ev.Modifiers
}

func (ev GoPointerEvent) Pos() (pos image.Point) {
	 return image.Point{
		X: int(math.Round(float64(ev.Position.X))),
		Y: int(math.Round(float64(ev.Position.Y))),
	}
}

func (ev GoPointerEvent) X() (x int) {
	return int(math.Round(float64(ev.Position.X)))
}

func (ev GoPointerEvent) Y() (y int) {
	return int(math.Round(float64(ev.Position.Y)))
}