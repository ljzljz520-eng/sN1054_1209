package store

import (
	"database/sql"
	"fmt"

	"example.com/familyitinerary/internal/model"
)

func (s *Store) GetFamily(id string) (model.Family, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var family model.Family
	if err := s.ensureOpen(); err != nil {
		return family, err
	}
	err := s.db.QueryRow("SELECT id,name,children,notes FROM families WHERE id=?", id).Scan(&family.ID, &family.Name, &family.Children, &family.Notes)
	if err == sql.ErrNoRows {
		return family, fmt.Errorf("family %q not found", id)
	}
	return family, err
}

func (s *Store) ListFamilies() ([]model.Family, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id,name,children,notes FROM families ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Family, 0)
	for rows.Next() {
		var family model.Family
		if err := rows.Scan(&family.ID, &family.Name, &family.Children, &family.Notes); err != nil {
			return nil, err
		}
		items = append(items, family)
	}
	return items, rows.Err()
}
