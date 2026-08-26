package store

import (
	"fmt"

	"example.com/familyitinerary/internal/model"
)

func (s *Store) SaveMessage(m model.ChatMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	_, err := s.db.Exec("INSERT INTO messages(id,itinerary_id,sequence_no,sender,body,status,failure) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET sequence_no=excluded.sequence_no,sender=excluded.sender,body=excluded.body,status=excluded.status,failure=excluded.failure", m.ID, m.Itinerary, m.Sequence, m.Sender, m.Body, model.NormalizeStatus(m.Status), m.Failure)
	return err
}

func (s *Store) UpdateMessageStatus(id, status, failure string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureOpen(); err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE messages SET status=?,failure=? WHERE id=?", model.NormalizeStatus(status), failure, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return fmt.Errorf("message %q not found", id)
	}
	return err
}

func (s *Store) ListMessages(itineraryID string) ([]model.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	rows, err := s.db.Query("SELECT id,itinerary_id,sequence_no,sender,body,status,failure FROM messages WHERE itinerary_id=? ORDER BY sequence_no,id", itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.ChatMessage, 0)
	for rows.Next() {
		var item model.ChatMessage
		if err := rows.Scan(&item.ID, &item.Itinerary, &item.Sequence, &item.Sender, &item.Body, &item.Status, &item.Failure); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) NextMessageSequence(itineraryID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.ensureOpen(); err != nil {
		return 0, err
	}
	var sequence int
	if err := s.db.QueryRow("SELECT COALESCE(MAX(sequence_no),0)+1 FROM messages WHERE itinerary_id=?", itineraryID).Scan(&sequence); err != nil {
		return 0, err
	}
	return sequence, nil
}
