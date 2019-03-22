package redis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const DbKey = "drivers"
const FieldName = "payload"

// Ensure CabService implements driver_location.CabService.
var _ driver_location.CabService = &CabService{}

// CabService represents a service for managing cab requests.
type CabService struct {
	db     **redis.Client
	logger *log.Logger
	now    *int64

	interval    *int64
	intervalNum int64
}

// StoreLocation stores drivers:<id> as key, and location info along with time as a single string field value in Redis.
func (s *CabService) StoreLocation(id string, loc *driver_location.Location) error {
	// Set time.
	var timeNow time.Time
	if *s.interval == 0 {
		s.intervalNum = 0
	}
	if *s.now == 0 {
		timeNow = time.Now()
	} else {
		timeNow = time.Unix(*s.now+*s.interval*s.intervalNum, 0)
		s.intervalNum += 1
	}

	//set redis key
	stream := fmt.Sprintf("%s:%s", DbKey, id)

	//set redis ID
	//timeNow := time.now()
	timeNowId := strconv.FormatInt(timeNow.Unix(), 10)

	//set redis field value
	payload := fmt.Sprintf("{\"latitude\":%f,\"longitude\":%f,\"updated_at\":\"%s\"}",
		loc.Latitude,
		loc.Longitude,
		timeNow.Format(time.RFC3339))

	dbId, err := (*s.db).XAdd(&redis.XAddArgs{
		Stream: stream,
		ID:     timeNowId,
		Values: map[string]interface{}{FieldName: payload},
	}).Result()
	if err != nil {
		s.logger.Printf("ERROR: storing data failed %#v", err)
	} else {
		s.logger.Printf("data stored in database (key=%s, id=%s)\n", stream, dbId)
	}
	s.logger.Printf("Redis: xadd %s %s %s %s", stream, timeNowId, FieldName, payload)

	return err
}

// GetDriverLocations queries database for the history of driver id's records and returns them.
func (s *CabService) GetDriverLocations(id string, minutes float64) (*driver_location.Locations, error) {
	// Set time.
	var timeNow time.Time
	if *s.interval == 0 {
		s.intervalNum = 0
	}
	if *s.now == 0 {
		timeNow = time.Now()
	} else {
		timeNow = time.Unix(*s.now+*s.interval, 0)
	}

	rTimeWindow := strconv.FormatInt(timeNow.Unix(), 10)
	seconds := int64(minutes * 60)
	lTimeWindow := strconv.FormatInt(timeNow.Add(time.Duration(-seconds)*time.Second).Unix(), 10)

	startTime := time.Now()

	//query database for locations
	stream := fmt.Sprintf("%s:%s", DbKey, id)
	s.logger.Printf("Redis query: xrange %s %s %s", stream, lTimeWindow, rTimeWindow)
	locations, err := (*s.db).XRange(stream, lTimeWindow, rTimeWindow).Result()
	if err != nil {
		s.logger.Println("ERROR: fetching driver locations failed")
		return nil, err
	}
	locLen := len(locations)

	s.logger.Printf("database query processed in %s (#locations=%d)\n", time.Now().Sub(startTime), locLen)
	startTime = time.Now()

	//concatenate locations
	var loc strings.Builder
	if locLen > 0 {
		locSize := int(unsafe.Sizeof(loc))
		loc.Grow(len(locations[0].Values[FieldName].(string))*locLen + locLen + 2 + locSize)
	}
	_, err = fmt.Fprintf(&loc, "{\"locations\":[")
	if err != nil {
		s.logger.Println("ERROR: preparing driver locations failed")
		return nil, err
	}
	for i := 0; i < locLen; i++ {
		_, err = fmt.Fprintf(&loc, "%s", locations[i].Values[FieldName])
		if err != nil {
			s.logger.Println("ERROR: preparing driver locations failed")
			return nil, err
		}
		if i != locLen-1 {
			_, err = fmt.Fprintf(&loc, ",")
			if err != nil {
				s.logger.Println("ERROR: preparing driver locations failed")
				return nil, err
			}
		}
	}
	_, err = fmt.Fprintf(&loc, "]}")
	if err != nil {
		s.logger.Println("ERROR: preparing driver locations failed")
		return nil, err
	}

	s.logger.Printf("preparing results after query processed in %s\n", time.Now().Sub(startTime))

	// Decode response into JSON.
	var locStruct driver_location.Locations
	if err := json.NewDecoder(bytes.NewReader([]byte(loc.String()))).Decode(&locStruct); err != nil {
		return nil, err
	}
	//err = json.Unmarshal([]byte(loc.String()), &locStruct)
	//if err != nil {
	//	return nil, err
	//}

	return &locStruct, nil
}
