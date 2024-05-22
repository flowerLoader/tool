package types

import "time"

const TIME_LAYOUT = "2006-01-02 15:04:05+00:00"

func FormatTime(t time.Time) string {
	return t.Format(TIME_LAYOUT)
}

func ParseTime(s string) (t time.Time, err error) {
	t, err = time.Parse(TIME_LAYOUT, s)
	return
}

func MustParseTime(s string) time.Time {
	t, err := ParseTime(s)
	if err != nil {
		panic(err)
	}
	return t
}
