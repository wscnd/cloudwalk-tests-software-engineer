package logparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogParser_detectMatches(t *testing.T) {
	logFilePath := "../../logs/qgames.log"
	file, err := os.Open(logFilePath)
	if err != nil {
		t.Fatalf("Failed to open logs/logs.log: %v", err)
	}
	defer file.Close()

	type fields struct {
		logfile *os.File
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Found 21 matches",
			fields: fields{
				logfile: file,
			},
			want: 21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := &LogParser{
				logfile: tt.fields.logfile,
			}

			matches := lp.detectMatches()

			require.Equal(t, tt.want, matches)
		})
	}
}
