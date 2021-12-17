package main

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	//"strings"

	"github.com/fiatjaf/go-nostr/event"
	"github.com/fiatjaf/go-nostr/filter"
)

func queryEvents(filter *filter.EventFilter) (events []event.Event, err error) {
	var id string
	var author string
	var kind uint8
	var authors []string
	var tagEvent string
	var tagProfile string
	var since uint32

	var eventsReturn []event.Event
	dbQueryInterface := *badgerhold.Query{}
	//evtKindSetMetadataQuery := db.Find(&event.Event{}, badgerhold.Where("pubkey").Eq(evt.PubKey).And("kind").Eq(0))
	var conditions []string
	var params []interface{}

	if filter == nil {
		err = errors.New("filter cannot be null")
		return
	}

	if filter.ID != "" {
		//conditions = append(conditions, "id = ?")
		//params = append(params, filter.ID)
		id = filter.ID
		dbQueryFilterID := badgerhold.Where("id").Eq(id)
		dbQuery := dbQueryFilterID
	}

	if filter.Author != "" {
		//conditions = append(conditions, "pubkey = ?")
		//params = append(params, filter.Author)
		author = filter.Author
		if dbQueryFilterID {
			// filter on ID and Author
			dbQueryFilterAuthorID := badgerhold.Where("id").Eq(id).And("pubkey").Eq(author)
			dbQuery := dbQueryFilterAuthorID 
		}
		else {
			// filter on just author 
			dbQueryFilterAuthor := badgerhold.Where("pubkey").Eq(author) 
		}
	}

	if filter.Kind != nil && *filter.Kind != 0 {
		conditions = append(conditions, "kind = ?")
		params = append(params, filter.Kind)
		kind = filter.Kind
		if dbQueryFilterAuthorID {
			// filter on ID, Author, and kind
			dbQueryFilterKindAuthorID := badgerhold.Where("id").Eq(id).And("pubkey").Eq(author).And("kind").Eq(kind)
			dbQuery := dbQueryFilterKindAuthorID 
		} 
		else if dbQueryFilterID {
			// filter on ID and Kind
			dbQueryFilterKindID := badgerhold.Where("id").Eq(id).And("kind").Eq(kind)
			dbQuery := dbQueryFilterKindID
		} else {
			// filter on just Kind
			dbQueryFilterKind := badgerhold.Where("kind").Eq(kind)
			dbQuery := dbQueryFilterKind
		}
	}

	if filter.Authors != nil {
		if len(filter.Authors) == 0 {
			// authors being [] means you won't get anything
			return
		} else {
			inkeys := make([]string, 0, len(filter.Authors))
			for _, key := range filter.Authors {
				// to prevent sql attack here we will check if
				// these keys are valid 32byte hex
				parsed, err := hex.DecodeString(key)
				if err != nil || len(parsed) != 32 {
					continue
				}
				//inkeys = append(inkeys, fmt.Sprintf("'%x'", parsed))
				authors = append(authors, parsed.(string))
			}
			if dbQueryFilterKindAuthorID {
				// filter on ID, Author, kind, and Authors
				dbQueryFilterAuthorsKindAuthorID := badgerhold.Where("id").Eq(id).And("pubkey").Eq(author).And("kind").Eq(kind).And("authors").In(authors)
				dbQuery := dbQueryFilterAuthorsKindAuthorID
			} 
			else if dbQueryFilterKindID {
				// filter on ID, Kind, and Authors
				dbQueryFilterAuthorsKindID := badgerhold.Where("id").Eq(id).And("kind").Eq(kind).And("authors").In(authors)
				dbQuery := dbQueryFilterAuthorsKindID
			} else if dbQueryFilterKind {
				// filter on Kind and Authors
				dbQueryFilterAuthorsKind := badgerhold.Where("kind").Eq(kind).And("authors").In(authors)
				dbQuery := dbQueryFilterAuthorsKind
			} 
			else {
				// filter on just Authors
				dbQueryFilterAuthors := badgerhold.Where("authors").In(authors)
			}
			//conditions = append(conditions, `pubkey IN (`+strings.Join(inkeys, ",")+`)`)
		}
	}

	if filter.TagEvent != "" {
		//conditions = append(conditions, relatedEventsCondition)
		//params = append(params, filter.TagEvent)
		tagEvent = filter.TagEvent
		
		if dbQueryFilterAuthorsKindAuthorID {
			// filter on ID, Author, Kind, Authors, and tagEvent
			dbQueryFilterAuthorsKindAuthorID := badgerhold.Where("id").Eq(id).And("pubkey").Eq(author).And("kind").Eq(kind).And("authors").In(authors).And("tags").Eq(tagEvent)
		}
		else if dbQueryFilterKindAuthorID {
			// filter on ID, Author, kind, and tagEvent
			dbQueryFilterTagEventKindAuthorID := badgerhold.Where("id").Eq(id).And("pubkey").Eq(author).And("kind").Eq(kind).And("tags").Eq(tagEvent)
			dbQuery := dbQueryFilterAuthorsKindAuthorID
		} 
		else if dbQueryFilterKindID {
			// filter on ID, Kind, and tagEvent
			dbQueryFilterTagEventKindID := badgerhold.Where("id").Eq(id).And("kind").Eq(kind).And("authors").In(authors)
			dbQuery := dbQueryFilterAuthorsKindID
		} else if dbQueryFilterKind {
			// filter on Kind and Authors
			dbQueryFilterAuthorsKind := badgerhold.Where("kind").Eq(kind).And("authors").In(authors)
			dbQuery := dbQueryFilterAuthorsKind
		} 
		else {
			// filter on just Authors
			dbQueryFilterAuthors := badgerhold.Where("authors").In(authors)
		}
	}

	if filter.TagProfile != "" {
		//conditions = append(conditions, relatedEventsCondition)
		//params = append(params, filter.TagProfile)
		tagProfile = filter.TagProfile
		if dbQuery == "badgerhold" {
			dbQuery = dbQuery + ".Where(\"tags\").Eq(tagProfile)"
		} else {
			dbQuery = dbQuery + ".And(\"tags\").Eq(tagProfile)"
		}
	}

	if filter.Since != 0 {
		//conditions = append(conditions, "created_at > ?")
		//params = append(params, filter.Since)
		since = filter.Since
		if dbQuery == "badgerhold" {
			dbQuery = dbQuery + ".Where(\"created_at\").Gt(since)"
		} else {
			dbQuery = dbQuery + ".And(\"created_at\").Gt(since)"
		}
	}

	/*
	if len(conditions) == 0 {
		// fallback
		conditions = append(conditions, "true")
	}
	*/

	fmt.Println(dbQuery)
	//dbQuery = "db.Find(&event.Event{}, " + dbQuery + ")"
	fmt.Println(dbQuery)
	dbQueryInterface = dbQuery
	/*
	query := db.Rebind("SELECT * FROM event WHERE " +
		strings.Join(conditions, " AND ") +
		" ORDER BY created_at LIMIT 100")
	*/
	events = db.Find(&event.Event{}, dbQuery )
	err = db.Select(&events, query, params...)
	if err != nil && err != sql.ErrNoRows {
		log.Warn().Err(err).Interface("filter", filter).Msg("failed to fetch events")
		err = fmt.Errorf("failed to fetch events: %w", err)
	}

	return
}
