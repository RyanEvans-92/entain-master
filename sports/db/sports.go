package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"sports/proto/sports"
)

// EventsRepo provides repository access to events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the events repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy events.
		err = r.seed()
	})

	return err
}

func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
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
	// or omitted, all events will be returned.
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

func (m *eventsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.MeetingId, &event.Name, &event.Number, &event.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		/* 	Compare advertised start time of event to current timestamp.
		If current time is after advertised start time, event is closed,
		otherwise it is still open.
		*/
		if time.Now().After(time.Unix(ts.Seconds, int64(ts.Nanos))) {
			event.Status = "CLOSED"
		} else {
			event.Status = "OPEN"
		}

		events = append(events, &event)
	}

	return events, nil
}
