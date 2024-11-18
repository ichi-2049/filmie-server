package types

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

const (
	DefaultDateFormat = "2006-01-02"
	JSONDateFormat    = `"` + DefaultDateFormat + `"`
)

var (
	DateTZ = JST
)

type Date struct {
	time.Time
}

func NowDate() *Date {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, DateTZ)
	return &Date{today}
}

func NewDate(t time.Time) *Date {
	return &Date{t}
}

func (d Date) MarshalJSON() ([]byte, error) {
	s, err := formatTime(d.Time, JSONDateFormat, DateTZ)
	if err != nil {
		return nil, errors.Wrap(err, "error in types.date #MarshalJSON")
	}
	return []byte(s), nil
}

func (d *Date) UnmarshalJSON(v []byte) error {
	s := string(v)
	loc, err := ParseTime(s, JSONDateFormat, DateTZ)
	if err != nil {
		return errors.Wrap(err, "error in types.Date #UnmarshalJSON")
	}
	d.Time = *loc
	return nil
}

func MarshalDate(dt Date) graphql.Marshaler {
	if dt.Time.IsZero() {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		m, _ := dt.MarshalJSON()
		_, _ = w.Write(m)
	})
}

func UnmarshalDate(v interface{}) (Date, error) {
	d := Date{}
	if tmpStr, ok := v.(string); ok {
		tmpStr = fmt.Sprintf("\"%s\"", tmpStr)
		err := d.UnmarshalJSON([]byte(tmpStr))
		return d, err
	}
	return d, fmt.Errorf("time should be RFC3339 formatted string")
}

func ParseDate(layout string, str string, l *time.Location) (*Date, error) {
	loc, err := time.ParseInLocation(layout, str, l)
	if err != nil {
		return nil, err
	}
	return &Date{loc}, err
}
