package unixtime2rfc

import (
	"errors"
	"fmt"
	"time"
)

var layoutMap = map[string]string{
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
	"DateTime":    time.DateTime,
	"DateOnly":    time.DateOnly,
	"TimeOnly":    time.TimeOnly,
}

// ErrInvalidLayout is returned when a non-empty, unsupported layout string is provided.
var ErrInvalidLayout = errors.New("invalid layout string")

// UnixSecToTime converts a Unix timestamp (seconds) to a UTC time.Time object.
func UnixSecToTime(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

// UnixMilliToTime converts a Unix timestamp (milliseconds) to a UTC time.Time object.
func UnixMilliToTime(msec int64) time.Time {
	return time.UnixMilli(msec).UTC()
}

// UnixMicroToTime converts a Unix timestamp (microseconds) to a UTC time.Time object.
func UnixMicroToTime(usec int64) time.Time {
	return time.UnixMicro(usec).UTC()
}

// FormatTime converts a time.Time object to a formatted string using the provided layout.
// It returns an error if a non-empty and invalid layout string is provided.
func FormatTime(timeVal time.Time, layout string) (string, error) {
	format := time.RFC3339 // Default

	var ok bool

	if layout != "" {
		format, ok = layoutMap[layout]
		if !ok {
			return "", fmt.Errorf("%w: %s", ErrInvalidLayout, layout)
		}
	}

	return timeVal.Format(format), nil
}
