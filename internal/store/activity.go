package store

import (
	"example.com/familyitinerary/internal/model"
)

func (s *Store) SaveActivity(a model.Activity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("INSERT INTO activities(id,itinerary_id,day,title,category,duration,child_safe) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET itinerary_id=excluded.itinerary_id,day=excluded.day,title=excluded.title,category=excluded.category,duration=excluded.duration,child_safe=excluded.child_safe", a.ID, a.Itinerary, a.Day, a.Title, a.Category, a.Duration, boolInt(a.ChildSafe))
	return err
}

func (s *Store) ListActivities(itineraryID string) ([]model.Activity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id,itinerary_id,day,title,category,duration,child_safe FROM activities WHERE itinerary_id=? ORDER BY day,id", itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Activity, 0)
	for rows.Next() {
		var item model.Activity
		var childSafe int
		if err := rows.Scan(&item.ID, &item.Itinerary, &item.Day, &item.Title, &item.Category, &item.Duration, &childSafe); err != nil {
			return nil, err
		}
		item.ChildSafe = intBool(childSafe)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) DeleteActivity(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("DELETE FROM activities WHERE id=?", id)
	return err
}
