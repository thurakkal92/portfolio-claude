package content

import "time"

func formatDate(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case time.Time:
		return t.Format("2006-01-02")
	case string:
		return t
	}
	return ""
}
