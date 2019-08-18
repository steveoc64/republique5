package store

import (
	"sync"

	"fyne.io/fyne"
	rp "github.com/steveoc64/republique5/proto"
)

// commander must implement a DataItem
type Commander struct {
	sync.Mutex
	Data    *rp.Command
	nextSub fyne.ListenerHandle
	subs    map[fyne.ListenerHandle]fyne.DataItemFunc
}

// String to implement DataItem
func (c *Commander) String() string {
	return c.Data.Name
}

// AddListener registers listener function
func (c *Commander) AddListener(fn fyne.DataItemFunc) fyne.ListenerHandle {
	c.Lock()
	defer c.Unlock()
	subID := c.nextSub
	if c.subs == nil {
		c.subs = make(map[fyne.ListenerHandle]fyne.DataItemFunc)
	}
	c.subs[subID] = fn
	c.nextSub++
	return subID
}

// DeleteListener removes a listener
func (c *Commander) DeleteListener(handle fyne.ListenerHandle) {
	c.Lock()
	defer c.Unlock()
	delete(c.subs, handle)
}

// Refresh will trigger all the listeners for this DataItem in a goroutine
func (c *Commander) Refresh() {
	c.Lock()
	defer c.Unlock()
	for _, fn := range c.subs {
		go fn(c)
	}
}
