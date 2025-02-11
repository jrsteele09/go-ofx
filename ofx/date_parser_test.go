package ofx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DateParsing(t *testing.T) {
	testDates := []string{
		"20231004123045.123+0200[PST]",
		"20231004",
		"20231004123045",
		"20231004123045.123",
		"20231004123045.123+0100[GMT]",
		"20230615120000.000[+1]",
		"20240503122717[-5:EST]",
	}

	for _, dateStr := range testDates {
		parsedDate, err := ParseOFXDateTime(dateStr)
		require.NoError(t, err)
		fmt.Println("Parsed date:", parsedDate)
	}

}
