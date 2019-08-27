package store

// Store is a UI application store that things can listen in on
type Store struct {
	CommanderMap CommanderMap
}

// NewStore returns a clean new store
func NewStore() *Store {
	return &Store{
		CommanderMap: newCommanderMap(),
	}
}
