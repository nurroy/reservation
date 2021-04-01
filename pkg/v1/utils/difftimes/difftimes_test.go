package difftimes

import (
	"testing"
	"time"
)

func TestDiff(t *testing.T) {
	type args struct {
		aTimeString string
		bTimeString string
		formatTime  string
	}
	tests := []struct {
		name         string
		args         args
		wantDiffYear int
	}{
		{
			name: "sample success",
			args: args{
				formatTime:  "2006-01-02T15:04:05Z",
				aTimeString: "2019-02-01T00:00:00Z",
				bTimeString: "2023-02-01T00:00:00Z",
			},
			wantDiffYear: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aTime, _ := time.Parse(tt.args.formatTime, tt.args.aTimeString)
			bTime, _ := time.Parse(tt.args.formatTime, tt.args.bTimeString)

			year, _, _, _, _, _ := Diff(aTime, bTime)
			if year != tt.wantDiffYear {
				t.Errorf("Diff() get result = %v, wantDiffYear = %v", year, tt.wantDiffYear)
				return
			}
		})
	}
}

func TestGetAge(t *testing.T) {
	type args struct {
		birthday   string
		formatTime string
	}
	tests := []struct {
		name    string
		args    args
		wantAge int
	}{
		{
			name: "sample success",
			args: args{
				formatTime: "2006-01-02T15:04:05Z",
				birthday:   "1994-04-12T00:00:00Z",
			},
			wantAge: 26,
		},
		{
			name: "sample success -1 second before birthday",
			args: args{
				formatTime: "2006-01-02T15:04:05-0700",
				birthday:   "1994-06-15T01:00:01+0700",
			},
			wantAge: 25,
		},
		{
			name: "sample success +1 second after birthday",
			args: args{
				formatTime: "2006-01-02T15:04:05-0700",
				birthday:   "1994-06-15T00:59:59+0700",
			},
			wantAge: 26,
		},
		{
			name: "sample error",
			args: args{
				formatTime: "1994-04-12T00:00:00Z",
				birthday:   "2006-01-02T15:04:05Z",
			},
			wantAge: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			age := GetAge(tt.args.formatTime, tt.args.birthday, "testing")
			if age != tt.wantAge {
				t.Errorf("GetAge() error = %v, wantAge %v", age, tt.wantAge)
				return
			}
		})
	}
}
