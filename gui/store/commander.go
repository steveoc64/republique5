package store

import (
	"sync"

	"fyne.io/fyne"
	rp "github.com/steveoc64/republique5/proto"
)

////////////////////////////////////////////////////////////////////////////////
// commander for tracking the orders per commander

// commander must implement a DataItem
type commander struct {
	sync.Mutex
	data    *rp.CommandGameState
	nextSub fyne.ListenerHandle
	subs    map[fyne.ListenerHandle]fyne.DataItemFunc
}

// String to implement DataItem
func (c commander) String() string {
	return c.data.GetOrders().String()
}

// AddListener registers listener function
func (c commander) AddListener(fn fyne.DataItemFunc) fyne.ListenerHandle {
	c.Lock()
	defer c.Unlock()
	c.subs[c.nextSub] = fn
	c.nextSub++
	return c.nextSub
}

// DeleteListener removes a listener
func (c commander) DeleteListener(handle fyne.ListenerHandle) {
	c.Lock()
	defer c.Unlock()
	delete(c.subs, handle)
}

// Refresh will trigger all the listeners for this DataItem in a goroutine
func (c commander) Refresh() {
	c.Lock()
	defer c.Unlock()
	for _, fn := range c.subs {
		go fn(c)
	}
}

////////////////////////////////////////////////////////////////////////////////
// CommanderMap for tracking all the commanders in a map

// CommanderOrdersMap map by commanderID of commanderOrder listeners
type CommanderMap struct {
	sync.RWMutex
	commanders map[int32]commander
}

// newCommanderMap - internal function to allocate a new map
func newCommanderMap() CommanderMap {
	return CommanderMap{
		commanders: make(map[int32]commander),
	}
}

// Load resets the whole map from the commander array
func (c *CommanderMap) Load(data []*rp.Command) {
	c.Lock()
	defer c.Unlock()
	// TODO - merge the data in, dont just overwrite it, otherwise the listeners all get trashed
	c.commanders = make(map[int32]commander)
	for _, command := range data {
		c.commanders[command.Id] = commander{
			Mutex:   sync.Mutex{},
			data:    command.GameState,
			nextSub: 0,
			subs:    make(map[fyne.ListenerHandle]fyne.DataItemFunc),
		}
	}
}

func (c *CommanderMap) AddListener(command *rp.Command, fn fyne.DataItemFunc) fyne.ListenerHandle {
	c.RLock()
	defer c.RUnlock()
	if commander, ok := c.commanders[command.Id]; ok {
		return commander.AddListener(fn)
	}
	return -1
}

func (c *CommanderMap) Refresh(command *rp.Command) {
	c.RLock()
	defer c.RUnlock()
	if commander, ok := c.commanders[command.Id]; ok {
		commander.Refresh()
	}
}
