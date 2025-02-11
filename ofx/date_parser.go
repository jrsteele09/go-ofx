package ofx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseOFXDateTime parses OFX date/time strings in various formats
func ParseOFXDateTime(dateTimeStr string) (time.Time, error) {
	layout := ""
	dateTimeStr = strings.TrimSpace(dateTimeStr)

	// Handle case: YYYYMMDD
	if len(dateTimeStr) == 8 {
		layout = "20060102"
		return time.Parse(layout, dateTimeStr)
	}

	// Handle case: YYYYMMDDHHMMSS
	if len(dateTimeStr) == 14 {
		layout = "20060102150405"
		return time.Parse(layout, dateTimeStr)
	}

	// Handle case: YYYYMMDDHHMMSS.XXX
	if len(dateTimeStr) > 15 && dateTimeStr[14] == '.' {
		layout = "20060102150405.999"
		if len(dateTimeStr) > 18 && (dateTimeStr[18] == '+' || dateTimeStr[18] == '-' || dateTimeStr[18] == '[') {
			layout += "-0700"
		}
		// Strip any timezone name or non-standard offset and convert if necessary
		dateTimeStr = convertCustomTimeZone(dateTimeStr)
		return time.Parse(layout, dateTimeStr)
	}

	// Handle case: YYYYMMDDHHMMSS.XXX [gmt offset[:tz name]]
	// Example: 20240503122717[-5:EST]
	var tzOffsetIndex int
	if tzOffsetIndex = strings.IndexAny(dateTimeStr, "+-"); tzOffsetIndex != -1 {
		layout = "20060102150405-0700"
		// Convert custom timezone formats like [-5:EST]
		dateTimeStr = convertCustomTimeZone(dateTimeStr)
		return time.Parse(layout, dateTimeStr)
	}

	// If no timezone info is present, assume GMT
	return time.Parse("20060102150405.999", dateTimeStr)
}

// convertCustomTimeZone converts custom timezone formats like [-5:EST] to -0500
func convertCustomTimeZone(dateTimeStr string) string {
	leftBracketIdx := strings.Index(dateTimeStr, "[")
	rightBracketIdx := strings.Index(dateTimeStr, "]")
	if leftBracketIdx == -1 || rightBracketIdx == -1 {
		return dateTimeStr
	}
	timeZoneStr := dateTimeStr[leftBracketIdx+1 : rightBracketIdx]
	colonIdx := strings.Index(timeZoneStr, ":")
	if colonIdx != -1 {
		timeZoneStr = timeZoneStr[:colonIdx]
	}

	timeZoneInt, err := strconv.Atoi(timeZoneStr)
	if err != nil {
		return dateTimeStr[:leftBracketIdx]
	}

	dateTimeStr = fmt.Sprintf("%s%+03d00", dateTimeStr[:leftBracketIdx], timeZoneInt)
	return dateTimeStr
}
