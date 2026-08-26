package store

import (
	"database/sql"
	"fmt"

	"example.com/familyitinerary/internal/model"
)

func (s *Store) SaveItinerary(i model.Itinerary) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("INSERT INTO itineraries(id,family_id,title,destination,start_date,end_date,status) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET family_id=excluded.family_id,title=excluded.title,destination=excluded.destination,start_date=excluded.start_date,end_date=excluded.end_date,status=excluded.status", i.ID, i.FamilyID, i.Title, i.Destination, i.StartDate, i.EndDate, model.NormalizeStatus(i.Status))
	return err
}

func (s *Store) GetItinerary(id string) (model.Itinerary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var item model.Itinerary
	if err := s.ensureOpen(); err != nil {
		return item, err
	}
	err := s.db.QueryRow("SELECT id,family_id,title,destination,start_date,end_date,status FROM itineraries WHERE id=?", id).Scan(&item.ID, &item.FamilyID, &item.Title, &item.Destination, &item.StartDate, &item.EndDate, &item.Status)
	if err == sql.ErrNoRows {
		return item, fmt.Errorf("itinerary %q not found", id)
	}
	return item, err
}

func (s *Store) ListItineraries(familyID string) ([]model.Itinerary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id,family_id,title,destination,start_date,end_date,status FROM itineraries WHERE family_id=? ORDER BY start_date,id", familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Itinerary, 0)
	for rows.Next() {
		var item model.Itinerary
		if err := rows.Scan(&item.ID, &item.FamilyID, &item.Title, &item.Destination, &item.StartDate, &item.EndDate, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateItineraryStatus(id, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE itineraries SET status=? WHERE id=?", model.NormalizeStatus(status), id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return fmt.Errorf("itinerary %q not found", id)
	}
	return err
}
