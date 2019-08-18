package store

import (
	"strconv"
	"sync"

	"fyne.io/fyne"
	rp "github.com/steveoc64/republique5/proto"
)

// CommanderOrdersMap map by commanderID of commanderOrder listeners
type CommanderMap struct {
	sync.RWMutex
	commanders map[int32]*Commander
	nextSub    fyne.ListenerHandle
	subs       map[fyne.ListenerHandle]fyne.DataMapFunc
}

// newCommanderMap - internal function to allocate a new map
func newCommanderMap() CommanderMap {
	return CommanderMap{
		commanders: make(map[int32]*Commander),
	}
}

// Reload resets the whole map from the commander array
// Note that it may get a whole new set of ptr data that overwrites the
// existing set, so listeners are preserved in this case
func (c *CommanderMap) Reload(data []*rp.Command) {
	c.Lock()
	defer c.Unlock()
	newCommanders := make(map[int32]*Commander)
	for _, command := range data {
		var subs map[fyne.ListenerHandle]fyne.DataItemFunc
		if oldCommander, ok := c.commanders[command.Id]; ok {
			subs = oldCommander.subs
		}
		newCommander := &Commander{
			Mutex:   sync.Mutex{},
			Data:    command,
			nextSub: 0,
			subs:    subs,
		}
		newCommanders[command.Id] = newCommander
	}
	c.commanders = newCommanders

	// fire all the commander subs
	for _, cmd := range c.commanders {
		cmd.Refresh()
	}
	// fire any map listeners
	for _, fn := range c.subs {
		go fn(c)
	}
}

func (c *CommanderMap) AddCommandListener(command *rp.Command, fn fyne.DataItemFunc) fyne.ListenerHandle {
	c.RLock()
	defer c.RUnlock()
	cmd, ok := c.commanders[command.Id]
	if !ok {
		cmd = &Commander{
			Mutex:   sync.Mutex{},
			Data:    command,
			nextSub: 0,
			subs:    nil,
		}
	}
	return cmd.AddListener(fn)
}

func (c *CommanderMap) AddListener(fn fyne.DataMapFunc) fyne.ListenerHandle {
	c.Lock()
	defer c.Unlock()
	subID := c.nextSub
	if c.subs == nil {
		c.subs = make(map[fyne.ListenerHandle]fyne.DataMapFunc)
	}
	c.subs[subID] = fn
	c.nextSub++
	return subID
}

func (c *CommanderMap) DeleteListener(handle fyne.ListenerHandle) {
	c.Lock()
	defer c.Unlock()
	delete(c.subs, handle)
}

func (c *CommanderMap) Get(s string) (fyne.DataItem, bool) {
	c.RLock()
	defer c.RUnlock()
	id, err := strconv.Atoi(s)
	if err != nil {
		return nil, false
	}
	cmd, ok := c.commanders[int32(id)]
	return cmd, ok
}

func (c *CommanderMap) Refresh(command *rp.Command) {
	c.RLock()
	defer c.RUnlock()
	if commander, ok := c.commanders[command.Id]; ok {
		commander.Refresh()
	}
}
