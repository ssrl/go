// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"runtime"
	"testing"
)

var stop = make(chan bool, 1)

func perpetuumMobile() {
	select {
	case <-stop:
	default:
		go perpetuumMobile()
	}
}

func TestStopTheWorldDeadlock(t *testing.T) {
	if testing.Short() {
		t.Logf("skipping during short test")
		return
	}
	maxprocs := runtime.GOMAXPROCS(3)
	compl := make(chan bool, 2)
	go func() {
		for i := 0; i != 1000; i += 1 {
			runtime.GC()
		}
		compl <- true
	}()
	go func() {
		for i := 0; i != 1000; i += 1 {
			runtime.GOMAXPROCS(3)
		}
		compl <- true
	}()
	go perpetuumMobile()
	<-compl
	<-compl
	stop <- true
	runtime.GOMAXPROCS(maxprocs)
}
