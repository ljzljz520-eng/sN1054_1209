package store

import "example.com/familyitinerary/internal/model"

func (s *Store) SavePreferences(itineraryID string, preferences []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM preferences WHERE itinerary_id=?", itineraryID); err != nil {
		tx.Rollback()
		return err
	}
	for position, preference := range preferences {
		if _, err = tx.Exec("INSERT INTO preferences(itinerary_id,position,value) VALUES(?,?,?)", itineraryID, position, preference); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) LoadPreferences(itineraryID string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT value FROM preferences WHERE itinerary_id=? ORDER BY position", itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		items = append(items, value)
	}
	return items, rows.Err()
}

func (s *Store) SaveIntake(in model.Intake) error {
	if err := s.SaveFamily(in.Family); err != nil {
		return err
	}
	if err := s.SaveItinerary(in.Itinerary); err != nil {
		return err
	}
	for _, activity := range in.Activities {
		if err := s.SaveActivity(activity); err != nil {
			return err
		}
	}
	return s.SavePreferences(in.Itinerary.ID, in.Preferences)
}
