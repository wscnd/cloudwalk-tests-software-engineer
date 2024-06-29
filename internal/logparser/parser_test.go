package logparser

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	logFilePath = "../../logs/matches_start_end.log"
	file        *os.File
)

func setup() error {
	var err error
	file, err = os.Open(logFilePath)
	if err != nil {
		return err
	}
	return nil
}

func teardown() {
	if file != nil {
		file.Close()
	}
}

func TestLogParser_detectMatchesLength(t *testing.T) {
	err := setup()
	require.NoError(t, err)
	defer teardown()

	lp := &LogParser{
		logfile: file,

		matchesLog: make([][]string, 0, 21),
	}

	err = lp.detectMatches()
	require.NoError(t, err)

	expectedMatches := 21
	actualMatches := len(lp.matchesLog)
	require.Equal(t, expectedMatches, actualMatches)
}

func TestLogParser_detectMatchesStartAndEnd(t *testing.T) {
	err := setup()
	require.NoError(t, err)
	defer teardown()

	lp := &LogParser{
		logfile:    file,
		matchesLog: make([][]string, 0, 21),
	}

	err = lp.detectMatches()
	require.NoError(t, err)

	testCases := []struct {
		name             string
		expectedInitTime string
		expectedEndTime  string
	}{
		{
			name:             "Match 1",
			expectedInitTime: "15:00",
			expectedEndTime:  "20:37",
		},
		{
			name:             "Match 2",
			expectedInitTime: "20:38",
			expectedEndTime:  "26:09",
		},
		{
			name:             "Match 3",
			expectedInitTime: "0:25",
			expectedEndTime:  "1:47",
		},
		{
			name:             "Match 4",
			expectedInitTime: "1:47",
			expectedEndTime:  "12:13",
		},
		{
			name:             "Match 5",
			expectedInitTime: "12:14",
			expectedEndTime:  "54:21",
		},
		{
			name:             "Match 6",
			expectedInitTime: "0:07",
			expectedEndTime:  "3:32",
		},
		{
			name:             "Match 7",
			expectedInitTime: "3:32",
			expectedEndTime:  "11:22",
		},
		{
			name:             "Match 8",
			expectedInitTime: "11:23",
			expectedEndTime:  "16:35",
		},
		{
			name:             "Match 9",
			expectedInitTime: "16:36",
			expectedEndTime:  "21:52",
		},
		{
			name:             "Match 10",
			expectedInitTime: "0:00",
			expectedEndTime:  "3:47",
		},
		{
			name:             "Match 11",
			expectedInitTime: "0:00",
			expectedEndTime:  "2:33",
		},
		{
			name:             "Match 12",
			expectedInitTime: "2:33",
			expectedEndTime:  "10:28",
		},
		{
			name:             "Match 13",
			expectedInitTime: "10:28",
			expectedEndTime:  "11:03",
		},
		{
			name:             "Match 14",
			expectedInitTime: "11:04",
			expectedEndTime:  "16:53",
		},
		{
			name:             "Match 15",
			expectedInitTime: "16:53",
			expectedEndTime:  "981:27",
		},
		{
			name:             "Match 16",
			expectedInitTime: "981:27",
			expectedEndTime:  "981:39",
		},
		{
			name:             "Match 17",
			expectedInitTime: "0:00",
			expectedEndTime:  "1:53",
		},
		{
			name:             "Match 18",
			expectedInitTime: "0:00",
			expectedEndTime:  "0:35",
		},
		{
			name:             "Match 19",
			expectedInitTime: "0:35",
			expectedEndTime:  "6:10",
		},
		{
			name:             "Match 20",
			expectedInitTime: "6:10",
			expectedEndTime:  "6:34",
		},
		{
			name:             "Match 21",
			expectedInitTime: "6:34",
			expectedEndTime:  "14:11",
		},
	}

	for tcs := 0; tcs < len(testCases); tcs++ {
		tc := testCases[tcs]
		t.Run(tc.name, func(t *testing.T) {
			match := lp.matchesLog[tcs]

			firstTiming := strings.Fields(match[0])[0]
			lastTiming := strings.Fields(match[len(match)-1])[0]

			require.Equal(t, tc.expectedInitTime, firstTiming)
			require.Equal(t, tc.expectedEndTime, lastTiming)
		})
	}
}
