package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	illegalCharFlag := false
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// Check if visible only has been supplied as true. If false,
	// or omitted, all races will be returned.
	if filter.VisibleOnly {
		clauses = append(clauses, "visible = ?")
		args = append(args, true)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	query += " ORDER BY "

	orderByStr := "advertised_start_time"

	// Function to check if an illegal character is entered by the user.
	strContainsIllegalChar := func(str string) bool {
		for i := 0; i < len(str); i++ {
			if str[i] == byte(';') {
				return true
			}
		}
		return false
	}

	/*	If the user supplied string can is a key that can be sorted by; if
		no sort_by field is supplied, order will default to by advertised_start_time.
		In the event a value that cannot be sorted by is supplied, the code already
		return an error code 2 for unknown strings, and code 3 for unsupported data
		types.
	*/
	if s := strings.ToLower(filter.SortBy); len(s) > 0 {
		query += s
	} else if strContainsIllegalChar(s) {
		illegalCharFlag = true
		query += orderByStr
	} else {
		query += orderByStr
	}

	/*	Since this is a binary option of ascending or descending order
		the code could be made more efficient by making sort_by a boolean
		rather than a string. However, it will make it easier for users if
		they can specify in natural language rather than 0 or 1. The default
		order, if not specified will be descending, making the most recent
		races appear first.
	*/

	if s := strings.ToUpper(filter.Order); s == "ASC" || s == "ASCENDING" {
		query += " ASC"
	} else if strContainsIllegalChar(s) {
		illegalCharFlag = true
		query += " DESC"
	} else {
		query += " DESC"
	}

	/* 	Not really possible, from my testing, to get extraneous SQL commands to
	execute on the server, but given the order_by part allows user input with errors
	handled higher up it might potentially be possible for a malicious user to sneak
	a semi-colon in and run a distinct SQL command, especially if the order of code
	is moved around throughout development. Probably not necessary, but worth noting
	something to do with security, especially since the API uses POST instead of GET.
	*/
	if illegalCharFlag {
		// Do something  useful here, like flag user.
		{
		}
	}

	return query, args
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		/* 	Compare advertised start time of race to current timestamp.
		If current time is after advertised start time, race is closed,
		otherwise it is still open. All races in the databaseappear to
		be circa 2021, so it is expected that they are all closed.
		*/
		if time.Now().After(time.Unix(ts.Seconds, int64(ts.Nanos))) {
			race.Status = "CLOSED"
		} else {
			race.Status = "OPEN"
		}

		races = append(races, &race)
	}

	return races, nil
}
