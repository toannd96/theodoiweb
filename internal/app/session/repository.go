package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"analytics-api/configs"
	"analytics-api/internal/pkg/common"
	"analytics-api/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/sirupsen/logrus"
)

// Repository ...
type Repository interface {
	Insert(session models.Session, event models.Event) error
	GetTotalColumn(sessionID string) (int64, error)
	GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error
	GetSessionByID(sessionID string, session *models.Session) error
	GetAllSession(listID []string, session models.Session) ([]models.Session, error)
	GetAllSessionID() ([]string, error)
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetTotalColumn get total number column of session record
func (instance *repository) GetTotalColumn(sessionID string) (int64, error) {
	queryAPI := configs.Database.Client.QueryAPI(configs.Database.Organization)
	var numberColumn int64

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> group(columns: ["_measurement"])
	|> filter(fn: (r) => r["sessionID"] == "%s")
	|> group()
    |> count(column: "sessionID")`, configs.Database.Bucket, configs.Database.Measurement, sessionID)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, err
	}
	for result.Next() {
		numberColumn = result.Record().Values()["sessionID"].(int64)
	}

	return numberColumn, nil
}

// GetEventLimitBySessionID get limit event of session by session id
func (instance *repository) GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error {
	var event models.Event
	queryAPI := configs.Database.Client.QueryAPI(configs.Database.Organization)

	fmt.Printf("\n")
	logrus.Debug("limit ", limit, " offset ", offset)

	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -7d)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["sessionID"] == "%s")
		|> group()
		|> limit(n:%d, offset: %d)`, configs.Database.Bucket, configs.Database.Measurement, sessionID, limit, offset)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for result.Next() {
		values := result.Record().Values()

		timestampString := values["timestamp"].(string)
		timestamp, err := common.StringToInt64(timestampString)
		if err != nil {
			return err
		}
		event.Timestamp = timestamp

		typeString := values["type"].(string)
		typeEvent, err := common.StringToInt64(typeString)
		if err != nil {
			return err
		}
		event.Type = typeEvent

		json.Unmarshal([]byte(result.Record().Value().(string)), &event.Data)

		if len(events.Events) <= limit {
			logrus.Debug("event timestamp ", event.Timestamp)
			events.Events = append(events.Events, event)
		}
	}

	logrus.Debug("Number event ", len(events.Events))
	return nil
}

// GetSessionByID get session by session id
func (instance *repository) GetSessionByID(sessionID string, session *models.Session) error {
	var event models.Event
	data := event.Data

	queryAPI := configs.Database.Client.QueryAPI(configs.Database.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> filter(fn: (r) => r["sessionID"] == "%s")`, configs.Database.Bucket, configs.Database.Measurement, sessionID)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for result.Next() {
		values := result.Record().Values()

		session.SessionID = values["sessionID"].(string)
		session.OS = values["os"].(string)
		session.UserAgent = values["userAgent"].(string)
		session.Version = values["version"].(string)
		session.Browser = values["browser"].(string)
		session.SessionName = values["sessionName"].(string)
		session.UpdatedAt = values["updatedAt"].(string)

		timestampString := values["timestamp"].(string)
		timestamp, err := common.StringToInt64(timestampString)
		if err != nil {
			return err
		}
		event.SessionID = values["sessionID"].(string)
		event.Timestamp = timestamp

		typeString := values["type"].(string)
		typeEvent, err := common.StringToInt64(typeString)
		if err != nil {
			return err
		}
		event.Type = typeEvent

		json.Unmarshal([]byte(result.Record().Value().(string)), &data)
		event.Data = data

		session.Events = append(session.Events, event)
	}

	return nil
}

// GetAllSession get all session from list session id
func (instance *repository) GetAllSession(listID []string, session models.Session) ([]models.Session, error) {
	var sessions []models.Session

	queryAPI := configs.Database.Client.QueryAPI(configs.Database.Organization)

	for _, id := range listID {
		query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -7d)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["sessionID"] == "%s")`, configs.Database.Bucket, configs.Database.Measurement, id)

		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			return nil, err
		}
		for result.Next() {
			values := result.Record().Values()

			session.SessionID = values["sessionID"].(string)
			session.SessionName = values["sessionName"].(string)
			session.UpdatedAt = values["updatedAt"].(string)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetAllSessionID get all session id from all record
func (instance *repository) GetAllSessionID() ([]string, error) {
	var listID []string

	queryAPI := configs.Database.Client.QueryAPI(configs.Database.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> group(columns: ["_measurement"])`, configs.Database.Bucket, configs.Database.Measurement)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		listID = append(listID, result.Record().ValueByKey("sessionID").(string))
	}

	listID = common.RemoveDuplicateValues(listID)

	return listID, nil
}

// Insert insert session
func (instance *repository) Insert(session models.Session, event models.Event) error {
	mJson, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	jsonStr := string(mJson)

	p := influxdb2.NewPointWithMeasurement(configs.Database.Measurement).
		AddTag("sessionID", session.SessionID).
		AddTag("sessionName", session.SessionName).
		AddTag("browser", session.Browser).
		AddTag("os", session.OS).
		AddTag("userAgent", session.UserAgent).
		AddTag("version", session.Version).
		AddField("data", jsonStr).
		AddTag("type", strconv.FormatInt(event.Type, 10)).
		AddTag("timestamp", strconv.FormatInt(event.Timestamp, 10)).
		AddTag("updatedAt", session.UpdatedAt).
		SetTime(time.Now())

	writeAPI := configs.Database.Client.WriteAPIBlocking(configs.Database.Organization, configs.Database.Bucket)
	err = writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}
