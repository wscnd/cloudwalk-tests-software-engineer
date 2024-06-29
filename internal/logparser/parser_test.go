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
		matches: &Match{
			Matches: make([][]string, 0, 21),
		},
	}

	err = lp.detectMatches()
	require.NoError(t, err)

	expectedMatches := 21
	actualMatches := len(lp.matches.Matches)
	require.Equal(t, expectedMatches, actualMatches)
}

	}{
		{
		},
	}


		})
	}
}
