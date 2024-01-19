// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/metrics_test.go */

package metrics

import (
	"testing"
)

var Dpr float32 = 1.25
var Spr float32 = 1.5

func TestDpToPx(t *testing.T) {
	for i := 0; i < 256; i++ {
		px := DpToPx(Dpr, i)
		dp := PxToDp(Dpr, px)
		if dp != i {
			t.Errorf("DpToPx(%d) = %d; PxToDp(%d) = %d; want %d", i, px, px, dp, i)
		}
	}
}

func TestSpToPx(t *testing.T) {

	for i := 0; i < 256; i++ {
		px := SpToPx(Spr, i)
		sp := PxToSp(Spr, px)
		if sp != i {
			t.Errorf("SpToPx(%d) = %d; PxToSp(%d) = %d; want %d", i, px, px, sp, i)
		}
	}
}

func TestDpToSp(t *testing.T) {
	for i := 0; i < 256; i++ {
		sp := DpToSp(Dpr, Spr, i)
		dp := SpToDp(Spr, Dpr, sp)
		if dp != i {
			t.Errorf("DpToSp(%d) = %f; SpToDp(%f) = %d; want %d", i, sp, sp, dp, i)
		}
	}
}