package store

import (
	"example.com/familyitinerary/internal/model"
)

func (s *Store) SaveSuggestion(item model.DateSuggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("INSERT INTO date_suggestions(id,itinerary_id,suggested,reason,accepted) VALUES(?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET suggested=excluded.suggested,reason=excluded.reason,accepted=excluded.accepted", item.ID, item.Itinerary, item.Suggested, item.Reason, boolInt(item.Accepted))
	return err
}

func (s *Store) ListSuggestions(itineraryID string) ([]model.DateSuggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id,itinerary_id,suggested,reason,accepted FROM date_suggestions WHERE itinerary_id=? ORDER BY suggested,id", itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.DateSuggestion, 0)
	for rows.Next() {
		var item model.DateSuggestion
		var accepted int
		if err := rows.Scan(&item.ID, &item.Itinerary, &item.Suggested, &item.Reason, &accepted); err != nil {
			return nil, err
		}
		item.Accepted = intBool(accepted)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) AcceptSuggestion(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("UPDATE date_suggestions SET accepted=1 WHERE id=?", id)
	return err
}
