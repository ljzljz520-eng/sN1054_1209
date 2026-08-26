package store

import "database/sql"

const schema = `
CREATE TABLE IF NOT EXISTS families (id TEXT PRIMARY KEY, name TEXT NOT NULL, children INTEGER NOT NULL, notes TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS itineraries (id TEXT PRIMARY KEY, family_id TEXT NOT NULL, title TEXT NOT NULL, destination TEXT NOT NULL, start_date TEXT NOT NULL, end_date TEXT NOT NULL, status TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS activities (id TEXT PRIMARY KEY, itinerary_id TEXT NOT NULL, day INTEGER NOT NULL, title TEXT NOT NULL, category TEXT NOT NULL, duration INTEGER NOT NULL, child_safe INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS messages (id TEXT PRIMARY KEY, itinerary_id TEXT NOT NULL, sequence_no INTEGER NOT NULL, sender TEXT NOT NULL, body TEXT NOT NULL, status TEXT NOT NULL, failure TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS date_suggestions (id TEXT PRIMARY KEY, itinerary_id TEXT NOT NULL, suggested TEXT NOT NULL, reason TEXT NOT NULL, accepted INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS preferences (itinerary_id TEXT NOT NULL, position INTEGER NOT NULL, value TEXT NOT NULL, PRIMARY KEY (itinerary_id, position));
`

func initialize(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
