package store

import (
	"fmt"
	"strings"
)

type AuditEntry struct {
	Entity   string
	EntityID string
	Action   string
	Detail   string
}

func (s *Store) RecordAudit(entry AuditEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	if strings.TrimSpace(entry.Entity) == "" || strings.TrimSpace(entry.EntityID) == "" {
		return fmt.Errorf("audit entity identity is required")
	}
	if entry.Action == "" {
		return fmt.Errorf("audit action is required")
	}
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS audit_log (entity TEXT NOT NULL, entity_id TEXT NOT NULL, action TEXT NOT NULL, detail TEXT NOT NULL, position INTEGER PRIMARY KEY AUTOINCREMENT)")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO audit_log(entity,entity_id,action,detail) VALUES(?,?,?,?)", entry.Entity, entry.EntityID, entry.Action, entry.Detail)
	return err
}

func (s *Store) ListAudit(entityID string) ([]AuditEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT entity,entity_id,action,detail FROM audit_log WHERE entity_id=? ORDER BY position", entityID)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			return []AuditEntry{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	entries := make([]AuditEntry, 0)
	for rows.Next() {
		var entry AuditEntry
		if err := rows.Scan(&entry.Entity, &entry.EntityID, &entry.Action, &entry.Detail); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *Store) CountAudit(entityID string) (int, error) {
	entries, err := s.ListAudit(entityID)
	return len(entries), err
}
