// Copyright 2021 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routemgr

import (
	"sync"

	"github.com/vishvananda/netlink"

	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/pkg/gateway/xnet"
)

// Linux is a one-way exporter of routes to Linux kernel.
type Linux struct {
	// Device is the device for exporting the routes.
	Device netlink.Link

	mtx sync.Mutex
	// exportedRoutes stores routes published by the local process.
	exportedRoutes RouteDB
	// externalRoutes stores routes received from quagga.
	closeChan chan struct{}
}

func (l *Linux) init() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.closeChan == nil {
		l.closeChan = make(chan struct{})
		go func() {
			defer log.HandlePanic()
			l.exportedRoutes.Run()
		}()
	}
}

func (l *Linux) NewPublisher() Publisher {
	return l.exportedRoutes.NewPublisher()
}

func (l *Linux) Close() {
	l.init()
	close(l.closeChan)
}

func (l *Linux) Run() {
	l.init()
	consumer := l.exportedRoutes.NewConsumer()
Top:
	for {
		select {
		case update := <-consumer.Updates():
			err := l.publishToLinux(update)
			if err != nil {
				log.Error("Error when publishing to Linux", "err", err)
			}
		case <-l.closeChan:
			// Closed by the user.
			break Top
		}
	}
	consumer.Close()
	l.exportedRoutes.Close()
}

func (l *Linux) publishToLinux(update RouteUpdate) error {
	if update.IsAdd {
		return xnet.AddRoute(0, l.Device, update.Prefix, update.Source)
	}
	return xnet.DeleteRoute(0, l.Device, update.Prefix, update.Source)
}

func (l *Linux) Diagnostics() Diagnostics {
	return l.exportedRoutes.Diagnostics()
}
