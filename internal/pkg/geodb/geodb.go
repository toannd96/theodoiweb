package geodb

import (
	geoip2 "github.com/oschwald/geoip2-golang"
)

func Open(path string) (*geoip2.Reader, error) {
	db, err := geoip2.Open(path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
