// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/metrics.go */

package metrics

/*
Package metrics implements device independent sizes.

Device independent pixel, or dip, is the unit for sizes independent of
the underlying display device.

Scaled pixels, or sp, is the unit for text sizes. An sp is like dp with
text scaling applied.

Finally, pixels, or px, is the unit for display dependent pixels. Their
size vary between platforms and displays.

To maintain a constant visual size across platforms and displays, always
use dps or sps to define user interfaces. Only use pixels for derived
values.
*/

import (
	"math"
)

// Dp converts v to pixels, rounded to the nearest integer value. PXPerDp
func DpToPx(dpr float32, v int) int {
	return int(math.Round(float64(nonZero(dpr)) * float64(v)))
}

// PxToDp converts v px to dp.		PXPerDp
func PxToDp(dpr float32, v int) int {
	return int(float64(v) / nonZero(dpr))
}

// Sp converts v to pixels, rounded to the nearest integer value. PxPerSp
func SpToPx(spr float32, v int) int {
	return int(math.Round(float64(nonZero(spr)) * float64(v)))
}

// PxToSp converts v px to sp.		PXPerSp
func PxToSp(spr float32, v int) int {
	return int(float64(v) / nonZero(spr))
}

// DpToSp converts v dp to sp.
func DpToSp(dpr float32, spr float32, v float32) int {
	return SpToPx(spr, int(float64(v) * nonZero(dpr) / nonZero(spr)))
}

// SpToDp converts v sp to dp.
func SpToDp(spr float32, dpr float32, v float32) int {
	return DpToPx(dpr, int(float64(v) * nonZero(spr) / nonZero(dpr)))
}

func nonZero(v float32) float64 {
	if v == 0. {
		return 1
	}
	return float64(v)
}