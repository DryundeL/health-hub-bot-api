package graphql

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalDate преобразует time.Time в строку для GraphQL
func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf(`"%s"`, t.Format("2006-01-02")))
	})
}

// UnmarshalDate преобразует строку из GraphQL в time.Time
func UnmarshalDate(v interface{}) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}
	str, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("date must be a string")
	}
	return time.Parse("2006-01-02", str)
}

// MarshalTime преобразует time.Time в строку для GraphQL
func MarshalTime(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf(`"%s"`, t.Format(time.RFC3339)))
	})
}

// UnmarshalTime преобразует строку из GraphQL в time.Time
func UnmarshalTime(v interface{}) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}
	str, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("time must be a string")
	}
	return time.Parse(time.RFC3339, str)
}

