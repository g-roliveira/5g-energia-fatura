package app

import "time"

// timeOnly is just an alias to time.Time to keep the server_billing.go
// imports lean (the file was already getting long). The parser is ISO 8601
// date-only format (YYYY-MM-DD).
type timeOnly = time.Time

func parseISODate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}
