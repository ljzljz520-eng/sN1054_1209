package store

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "modernc.org/sqlite"

	"example.com/familyitinerary/internal/model"
)

type Store struct {
	db   *sql.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	if err = initialize(db); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db, path: path}, nil
}

func OpenMemory() (*Store, error) {
	return Open("file:familyitinerary?mode=memory&cache=shared")
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) ensureOpen() error {
	if s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return nil
}

func (s *Store) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if _, err := s.db.Exec("DELETE FROM preferences; DELETE FROM date_suggestions; DELETE FROM messages; DELETE FROM activities; DELETE FROM itineraries; DELETE FROM families"); err != nil {
		return err
	}
	return nil
}

func RemoveDatabase(path string) error {
	if path == "" || path == ":memory:" {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func intBool(value int) bool { return value != 0 }

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) SaveFamily(f model.Family) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("INSERT INTO families(id,name,children,notes) VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,children=excluded.children,notes=excluded.notes", f.ID, f.Name, f.Children, f.Notes)
	return err
}
