package store

// UI application store that things can listen in on
type Store struct {
	CommanderMap CommanderMap
}

func NewStore() *Store {
	return &Store{
		CommanderMap: newCommanderMap(),
	}
}
