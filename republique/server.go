package republique

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/steveoc64/republique5/db"
	rp "github.com/steveoc64/republique5/proto"
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

	// Write the game state to the gamestate file
	s.writeRunFile()
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

func (s *Server) writeRunFile() error {
	filename := s.filename
	if strings.HasSuffix(filename, ".db") {
		filename = filename[:len(filename)-3]
	}
	savename := filepath.Join(os.Getenv("HOME"), ".republique", filename+".run")
	fp, err := os.Create(savename)
	if err != nil {
		return err
	}
	fmt.Fprintf(fp, "pid:%d\nturn:%d\nphase:%d\ntime:%s\nname:%s\ntable:%dx%d\nadmin:%s\n",
		os.Getpid(),
		s.game.TurnNumber,
		s.game.Phase,
		time.Now().String(),
		s.game.Name,
		s.game.TableX, s.game.TableY,
		s.game.AdminAccess,
	)
	for _, team := range s.game.Scenario.GetTeams() {
		for i, player := range team.GetPlayers() {
			fmt.Fprintf(fp, "%s-%d:%s:%s\n", team.Name, i+1, team.AccessCode, player.AccessCode)
		}
	}
	fp.Close()
	return nil
}
