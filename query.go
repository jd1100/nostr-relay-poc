package main

import (
	//"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	//"strings"
	"github.com/timshannon/badgerhold/v4"
	"github.com/fiatjaf/go-nostr/event"
	"github.com/fiatjaf/go-nostr/filter"
)

func queryEvents(filter *filter.EventFilter) (events []event.Event, err error) {
	var id string
	var kind *int
	var authors []string
	var tagEvent string
	var tagProfile string
	var since uint32
	//var filterFields []interface{}
	var dbQuery = &badgerhold.Query{}
	//var eventsReturn []event.Event
	var filterFields = make(map[string]interface{})

	//dbQueryInterface := *badgerhold.Query{}
	//evtKindSetMetadataQuery := db.Find(&event.Event{}, badgerhold.Where("pubkey").Eq(evt.PubKey).And("kind").Eq(0))
	//var conditions []string
	//var params []interface{}

	if filter == nil {
		err = errors.New("filter cannot be null")
		return
	}

	if filter.ID != "" {
		//conditions = append(conditions, "id = ?")
		//params = append(params, filter.ID)
		id = filter.ID
		filterFields["id"] = id
		
		dbQuery = badgerhold.Where("ID").Eq(id)
	}

	if filter.Kind != nil && *filter.Kind != 0 {
		//conditions = append(conditions, "kind = ?")
		//params = append(params, filter.Kind)
		kind = filter.Kind
		//filterFields = append(filterFields, kind)
		filterFields["kind"] = kind
		dbQuery = dbQuery.And("Kind").Eq(kind)
	}

	if filter.Authors != nil {
		if len(filter.Authors) == 0 {
			// authors being [] means you won't get anything
			return
		} else {
			//inkeys := make([]string, 0, len(filter.Authors))
			for _, key := range filter.Authors {
				// to prevent sql attack here we will check if
				// these keys are valid 32byte hex
				parsed, err := hex.DecodeString(key)
				if err != nil || len(parsed) != 32 {
					continue
				}
				//inkeys = append(inkeys, fmt.Sprintf("'%x'", parsed))
				authors = append(authors, string(parsed))
				filterFields["authors"] = authors
				//filterFields = append(filterFields, authors)
			}
			dbQuery = dbQuery.And("Authors").In(authors)
			//conditions = append(conditions, `pubkey IN (`+strings.Join(inkeys, ",")+`)`)
		}
	}

	if filter.TagEvent != "" {
		//conditions = append(conditions, relatedEventsCondition)
		//params = append(params, filter.TagEvent)
		tagEvent = filter.TagEvent
		//filterFields = append(filterFields, tagEvent)
		filterFields["tagEvent"] = tagEvent
		dbQuery = dbQuery.And("Tags").Eq(tagEvent)
	}

	if filter.TagProfile != "" {
		//conditions = append(conditions, relatedEventsCondition)
		//params = append(params, filter.TagProfile)
		tagProfile = filter.TagProfile
		//filterFields = append(filterFields, tagProfile)
		filterFields["tagProfile"] = tagProfile
		dbQuery = dbQuery.And("Tags").Eq(tagProfile)
	}

	if filter.Since != 0 {
		//conditions = append(conditions, "created_at > ?")
		//params = append(params, filter.Since)
		since = filter.Since
		//filterFields = append(filterFields, since)
		filterFields["since"] = since
		dbQuery = dbQuery.And("CreatedAt").Gt(since)
	}

	for i, v := range filterFields {
		fmt.Println(i, v)
	}
	/*
	if len(conditions) == 0 {
		// fallback
		conditions = append(conditions, "true")
	}
	*/

	//dbQuery = "db.Find(&event.Event{}, " + dbQuery + ")"
	//fmt.Println(dbQuery)
	//dbQueryInterface = dbQuery
	/*
	query := db.Rebind("SELECT * FROM event WHERE " +
		strings.Join(conditions, " AND ") +
		" ORDER BY created_at LIMIT 100")
	*/
	err = db.Find(&events, dbQuery)
	//err = db.Select(&events, query, params...)
	
	if err != nil {
		log.Warn().Err(err).Interface("filter", filter).Msg("failed to fetch events")
		err = fmt.Errorf("failed to fetch events: %w", err)
	}
	
	return
}
