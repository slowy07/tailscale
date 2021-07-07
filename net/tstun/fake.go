// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tstun

import (
	"io"
	"log"
	"os"

	"golang.zx2c4.com/wireguard/tun"
)

type fakeTUN struct {
	evchan    chan tun.Event
	closechan chan struct{}
}

// NewFake returns a tun.Device that does nothing.
func NewFake() tun.Device {
	return &fakeTUN{
		evchan:    make(chan tun.Event),
		closechan: make(chan struct{}),
	}
}

func (t *fakeTUN) File() *os.File {
	panic("fakeTUN.File() called, which makes no sense")
}

func (t *fakeTUN) Close() error {
	close(t.closechan)
	close(t.evchan)
	return nil
}

func (t *fakeTUN) Read(out []byte, offset int) (int, error) {
	log.Println("TSTUN : FAKE")
	<-t.closechan
	return 0, io.EOF
}

func (t *fakeTUN) Write(b []byte, n int) (int, error) {
	log.Println("FAKE : Write Called")
	select {
	case <-t.closechan:
		return 0, ErrClosed
	default:
	}
	return len(b), nil
}

func (t *fakeTUN) Flush() error           { return nil }
func (t *fakeTUN) MTU() (int, error)      { return 1500, nil }
func (t *fakeTUN) Name() (string, error)  { return "FakeTUN", nil }
func (t *fakeTUN) Events() chan tun.Event { return t.evchan }
