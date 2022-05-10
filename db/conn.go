package db

import (
	"analytics-api/configs"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// NewClient open new client to influxdb
func NewClient() influxdb2.Client {
	configs.Database.Client = influxdb2.NewClientWithOptions(
		configs.Database.URL,
		configs.Database.Token,
		influxdb2.DefaultOptions().SetHTTPRequestTimeout(
			configs.Database.RequestTimeout,
		))
	return configs.Database.Client
}
