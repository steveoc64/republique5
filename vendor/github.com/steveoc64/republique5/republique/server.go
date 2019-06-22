package republique

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/db"
	rp "github.com/steveoc64/republique5/proto"
	"strings"
	"sync"
)

// Server object stores the whole game state
type Server struct {
	sync.RWMutex
	log         *logrus.Logger
	version     string
	filename    string
	port        int
	web         int
	game        *rp.Game
	db          *db.DB
	stopWatch   int64
	tokenCache  map[string]*rp.Player
	mTokenCache sync.RWMutex
}

// NewServer returns a new republique server
func NewServer(log *logrus.Logger, version string, filename string, port int, web int) (*Server, error) {
	// load the DB
	if !strings.HasSuffix(filename, ".db") {
		filename = filename + ".db"
	}
	db, err := db.OpenDB(log, filename)
	if err != nil {
		return nil, err
	}
	gameData := &rp.Game{}
	err = db.Load("game", "state", gameData)
	if err != nil {
		return nil, err
	}
	gameData.InitGameState()

	return &Server{
		log:        log,
		version:    version,
		filename:   filename,
		port:       port,
		web:        web,
		game:       gameData,
		db:         db,
		tokenCache: NewTokenCache(gameData),
	}, nil
}

// Serve runs a game
func (s *Server) Serve() {
	s.log.WithFields(logrus.Fields{
		"version":  s.version,
		"port":     s.port,
		"web":      s.web,
		"filename": s.filename,
	}).Println("Starting Republique 5.0 Server")
	s.log.SetFormatter(&logrus.JSONFormatter{})

	// Setup REST endpoints, but only if we want web with it
	if s.web != 0 {
		go s.rpcProxy()
	}

	// Load GPPC endpoints
	s.grpcRun()
}

// Save saves the game state to the DB
func (s *Server) Save() {
	s.db.Save("game", "state", s.game)
}

// Close closes the DB - needed when you quit
func (s *Server) Close() {
	s.db.Close()
}
