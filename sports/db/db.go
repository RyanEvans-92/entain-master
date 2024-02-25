package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME, level TEXT, sold_out INTEGER)`)
	if err == nil {
		_, err = statement.Exec()
	}

	// Anonymous function used to select a changing element from a string array.
	// Method is not random, but time will change quickly enough to give a variety
	// of results without additional imports. Good enough for creating a fake database.
	selectElement := func(arr []string) string {
		return arr[time.Now().Nanosecond()%len(arr)]
	}

	sportsList := []string{"Tennis", "Fencing", "Badminton", "Sportsketball", "Archery", "Caber Toss", "Football", "Soccer", "Competitive Crying", "Extreme Ironing", "Swimming", "Gymnastics", "Toe Wrestling", "Arguing"}
	levelList := []string{"Amateur", "Youth", "University", "Semi-Professional", "Professional", "International"}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO events(id, meeting_id, name, number, visible, advertised_start_time, level, sold_out) VALUES (?, ?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				selectElement(sportsList),
				faker.Number().Between(1, 12),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
				selectElement(levelList),
				faker.Number().Between(0, 1),
			)
		}
	}

	return err
}
