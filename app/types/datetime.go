package types

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

const (
	DefaultDateTimeFormat = time.RFC3339
	JSONDateTimeFormat    = `"` + DefaultDateTimeFormat + `"`
)

var JST = time.FixedZone("Asia/Tokyo", 9*60*60)

var (
	dateTimeTZ = JST
)

type DateTime struct {
	time.Time
}

func NowDateTimeJST() *DateTime {
	t := time.Now()
	return NewDateTime(t)
}

func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t}
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	s, err := formatTime(dt.Time, JSONDateTimeFormat, dateTimeTZ)
	if err != nil {
		return nil, errors.Wrap(err, "error in types.DateTime #MarshalJSON")
	}
	return []byte(s), nil
}

func (dt *DateTime) UnmarshalJSON(v []byte) error {
	s := string(v)
	loc, err := ParseTime(s, JSONDateTimeFormat, dateTimeTZ)
	if err != nil {
		return errors.Wrap(err, "error in types.DateTime #UnmarshalJSON")
	}
	dt.Time = *loc
	return nil
}

func MarshalDateTime(dt DateTime) graphql.Marshaler {
	if dt.Time.IsZero() {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		m, _ := dt.MarshalJSON()
		_, _ = w.Write(m)
	})
}

func UnmarshalDateTime(v interface{}) (DateTime, error) {
	dt := DateTime{}
	if tmpStr, ok := v.(string); ok {
		tmpStr = fmt.Sprintf("\"%s\"", tmpStr)
		err := dt.UnmarshalJSON([]byte(tmpStr))
		return dt, err
	}
	return dt, fmt.Errorf("time should be RFC3339 formatted string")
}

func formatTime(t time.Time, f string, l *time.Location) (string, error) {
	loc := t.In(l)
	if y := loc.Year(); y < 0 || y >= 10000 {
		return "", fmt.Errorf("year outside of range")
	}
	return loc.Format(f), nil
}

func ParseTime(s string, f string, l *time.Location) (*time.Time, error) {
	loc, err := time.ParseInLocation(f, s, l)
	if err != nil {
		return nil, err
	}
	return &loc, err
}
