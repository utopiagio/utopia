// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia/history/history.go */

package history

import (
	"log"
	)

func GoFileHistory(path string) (history *GoHistoryObj) {
	hhistory := &GoHistoryObj{[]string{path}, 0, path}
	log.Println("GoHistoryObj Add(", path, ")")
	return hhistory
}

type GoHistoryObj struct {
	paths []string
	pos int
	currentPath string
}

func (ob *GoHistoryObj) Add(path string) {
	log.Println("GoHistoryObj Add(", path, ")")
	log.Println("ob.pos =", ob.pos)
	log.Println("len(ob.paths) =", len(ob.paths))
	if ob.pos == len(ob.paths) - 1 {
		ob.paths = append(ob.paths, path)
	} else {
		ob.paths = append(ob.paths[:ob.pos + 1], path)
	}
	ob.pos++
	log.Println("ob.pos =", ob.pos)
	log.Println("len(ob.paths) =", len(ob.paths))
}

func (ob *GoHistoryObj) Back() (path string) {
	log.Println("GoHistoryObj Back()")
	log.Println("len(ob.paths) =", len(ob.paths))
	if ob.pos > 0 {
		ob.pos--
	}
	if ob.pos > -1 && ob.pos <= len(ob.paths) {
		path = ob.paths[ob.pos]
	}
	log.Println("ob.pos =", ob.pos)
	log.Println("return path =", path)
	return path
}

func (ob *GoHistoryObj) CurrentPath() (path string) {
	return ob.paths[ob.pos]
}

func (ob *GoHistoryObj) Forward() (path string) {
	log.Println("GoHistoryObj Forward()")
	
	log.Println("len(ob.paths) =", len(ob.paths))
	if ob.pos < len(ob.paths) - 1 {
		ob.pos++
		path = ob.paths[ob.pos]
	} else {
		path = ob.paths[ob.pos]
	}
	log.Println("ob.pos =", ob.pos)
	log.Println("return path =", path)
	return path
}

func (ob *GoHistoryObj) Up() (path string) {
	
	paths := ob.paths[:ob.pos]
	ob.pos++
	path = paths[ob.pos]
	ob.paths = paths
	return path
}