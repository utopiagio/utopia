// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/history.go */

package utopia

import (
	//"log"
	)

func GoFileHistory(path string) (history *GoHistoryObj) {
	hhistory := &GoHistoryObj{[]string{path}, 1}
	return hhistory
}

type GoHistoryObj struct {
	paths []string
	pos int
}

func (ob *GoHistoryObj) Back() (path string) {
	if ob.pos > 0 {
		ob.pos--
		path = ob.paths[ob.pos]
	}
	return path
}

func (ob *GoHistoryObj) Forward() (path string) {
	if ob.pos > len(ob.paths) {
		ob.pos++
		path = ob.paths[ob.pos]
	}
	return path
}

func (ob *GoHistoryObj) Up() (path string) {
	
	paths := ob.paths[:ob.pos]
	ob.pos++
	path = paths[ob.pos]
	ob.paths = paths
	return path
}