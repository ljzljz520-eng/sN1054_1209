package store

import (
	"fmt"

	"example.com/familyitinerary/internal/model"
)

func (s *Store) SaveIntakeAtomic(in model.Intake) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	rollback := func(cause error) error { _ = tx.Rollback(); return cause }
	if _, err = tx.Exec("INSERT INTO families(id,name,children,notes) VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,children=excluded.children,notes=excluded.notes", in.Family.ID, in.Family.Name, in.Family.Children, in.Family.Notes); err != nil {
		return rollback(err)
	}
	if _, err = tx.Exec("INSERT INTO itineraries(id,family_id,title,destination,start_date,end_date,status) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET family_id=excluded.family_id,title=excluded.title,destination=excluded.destination,start_date=excluded.start_date,end_date=excluded.end_date,status=excluded.status", in.Itinerary.ID, in.Itinerary.FamilyID, in.Itinerary.Title, in.Itinerary.Destination, in.Itinerary.StartDate, in.Itinerary.EndDate, model.NormalizeStatus(in.Itinerary.Status)); err != nil {
		return rollback(err)
	}
	if _, err = tx.Exec("DELETE FROM activities WHERE itinerary_id=?", in.Itinerary.ID); err != nil {
		return rollback(err)
	}
	for _, activity := range in.Activities {
		if _, err = tx.Exec("INSERT INTO activities(id,itinerary_id,day,title,category,duration,child_safe) VALUES(?,?,?,?,?,?,?)", activity.ID, activity.Itinerary, activity.Day, activity.Title, activity.Category, activity.Duration, boolInt(activity.ChildSafe)); err != nil {
			return rollback(err)
		}
	}
	if _, err = tx.Exec("DELETE FROM preferences WHERE itinerary_id=?", in.Itinerary.ID); err != nil {
		return rollback(err)
	}
	for position, preference := range in.Preferences {
		if _, err = tx.Exec("INSERT INTO preferences(itinerary_id,position,value) VALUES(?,?,?)", in.Itinerary.ID, position, preference); err != nil {
			return rollback(err)
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) VerifyItineraryConsistency(itineraryID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	var familyID string
	if err := s.db.QueryRow("SELECT family_id FROM itineraries WHERE id=?", itineraryID).Scan(&familyID); err != nil {
		return err
	}
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM families WHERE id=?", familyID).Scan(&count); err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("itinerary %q has no family", itineraryID)
	}
	return nil
}

func (s *Store) DeleteItinerary(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("DELETE FROM activities WHERE itinerary_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("DELETE FROM messages WHERE itinerary_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("DELETE FROM preferences WHERE itinerary_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("DELETE FROM date_suggestions WHERE itinerary_id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("DELETE FROM itineraries WHERE id=?", id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
