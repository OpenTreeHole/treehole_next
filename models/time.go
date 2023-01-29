package models

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

var locCST *time.Location

func init() {
	locCST, _ = time.LoadLocation("Asia/Shanghai")
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	ct.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		ct.Time, err = time.ParseInLocation(`2006-01-02T15:04:05`, s, locCST)
	}
	return err
}

func (ct *CustomTime) UnmarshalText(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	ct.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		ct.Time, err = time.ParseInLocation(`2006-01-02T15:04:05`, s, locCST)
	}
	return err
}
