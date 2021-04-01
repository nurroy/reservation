package postgres

import (
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func TestQueryUpdateBuilder(t *testing.T) {
	type args struct {
		tableName string
		mapArgs   map[string]interface{}
		wheres    []string
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantArgs  []interface{}
		wantErr   bool
	}{
		{
			name: "sample",
			args: args{
				tableName: "tabelku",
				mapArgs:   map[string]interface{}{"cobafield": 1, "primary": "abc"},
				wheres:    []string{"primary"},
			},
			wantQuery: "UPDATE tabelku SET cobafield = $1 WHERE primary = $2",
			wantArgs:  []interface{}{1, "abc"},
			wantErr:   false,
		},
		{
			name: "sample more than one where clause",
			args: args{
				tableName: "tabelku",
				mapArgs:   map[string]interface{}{"cobafield": 1, "primary": "abc", "primary2": "def"},
				wheres:    []string{"primary", "primary2"},
			},
			wantQuery: "UPDATE tabelku SET cobafield = $1 WHERE primary = $2 and primary2 = $3",
			wantArgs:  []interface{}{1, "abc", "def"},
			wantErr:   false,
		},
		{
			name: "invalid primary in where clause",
			args: args{
				tableName: "tabelku",
				mapArgs:   map[string]interface{}{"cobafield": 1, "primary": "abc"},
				wheres:    []string{"apa"},
			},
			wantQuery: "UPDATE tabelku SET cobafield = $1, primary = $2 WHERE ",
			wantArgs:  []interface{}{1, "abc"},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs, err := QueryUpdateBuilder(tt.args.tableName, tt.args.mapArgs, tt.args.wheres, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUpdateBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("QueryUpdateBuilder() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("QueryUpdateBuilder() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func Test_convertFieldToMapArgs(t *testing.T) {
	type TestingStruct struct {
		PrimaryKey   int64   `db:"primary_key"`
		SecondInt    int32   `db:"second_int" coalesce:"-"`
		ThirdString  string  `db:"third_string"`
		FourthBool   bool    `db:"fourth_bool" coalesce:"-"`
		FifthString  string  `db:"fifth_string" coalesce:"arkandia"`
		SixthFloat   float32 `db:"sixth_float"`
		SeventhFloat float64 `db:"seventh_float" coalesce:"1.2"`
		Eight        float32 `db:"-"`
	}

	var dataTest TestingStruct

	type args struct {
		data interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantArgs map[string]builderArgs
		wantErr  bool
	}{
		{
			name: "Test parsing nil",
			args: args{
				data: nil,
			},
			wantArgs: nil,
			wantErr:  true,
		},
		{
			name: "Test parsing non struct",
			args: args{
				data: "nil",
			},
			wantArgs: nil,
			wantErr:  true,
		},
		{
			name: "Test parsing empty struct",
			args: args{
				data: dataTest,
			},
			wantArgs: map[string]builderArgs{
				"primary_key":   {int64(0), ""},
				"second_int":    {int32(0), "coalesce(second_int, 0) as second_int"},
				"third_string":  {"", ""},
				"fourth_bool":   {false, "coalesce(fourth_bool, false) as fourth_bool"},
				"fifth_string":  {"", "coalesce(fifth_string, 'arkandia') as fifth_string"},
				"sixth_float":   {float32(0), ""},
				"seventh_float": {float64(0), "coalesce(seventh_float, 1.2) as seventh_float"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, err := convertFieldToMapArgs(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertFieldToMapArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("convertFieldToMapArgs() = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestQuerySelectBuilder(t *testing.T) {
	type TestingStruct struct {
		PrimaryKey  int64  `db:"primary_key"`
		SecondInt   int32  `db:"second_int" coalesce:"-"`
		ThirdString string `db:"third_string"`
	}

	type TestingPrimaryCoalesceStruct struct {
		PrimaryKey  int64  `db:"primary_key" coalesce:"-"`
		SecondInt   int32  `db:"second_int" coalesce:"-"`
		ThirdString string `db:"third_string"`
	}

	var dataTest TestingStruct
	var dataTest2 TestingPrimaryCoalesceStruct

	type args struct {
		tableName  string
		structData interface{}
		wheres     []string
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name: "testing valid struct with one where clause",
			args: args{
				tableName:  "coba",
				structData: dataTest,
				wheres:     []string{"primary_key"},
			},
			wantQuery: "SELECT primary_key, coalesce(second_int, 0) as second_int, third_string FROM coba WHERE primary_key = $1",
			wantErr:   false,
		},
		{
			name: "testing valid struct without where clause",
			args: args{
				tableName:  "coba",
				structData: dataTest,
				wheres:     []string{},
			},
			wantQuery: "SELECT primary_key, coalesce(second_int, 0) as second_int, third_string FROM coba",
			wantErr:   false,
		},
		{
			name: "testing valid struct with two where clause (and)",
			args: args{
				tableName:  "coba",
				structData: dataTest,
				wheres:     []string{"primary_key", "third_string"},
			},
			wantQuery: "SELECT primary_key, coalesce(second_int, 0) as second_int, third_string FROM coba WHERE primary_key = $1 and third_string = $2",
			wantErr:   false,
		},
		{
			name: "testing valid struct with one where clause - primary coalesce",
			args: args{
				tableName:  "coba",
				structData: dataTest2,
				wheres:     []string{"primary_key"},
			},
			wantQuery: "SELECT coalesce(primary_key, 0) as primary_key, coalesce(second_int, 0) as second_int, third_string FROM coba WHERE primary_key = $1",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, err := QuerySelectBuilder(tt.args.tableName, tt.args.structData, tt.args.wheres)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerySelectBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("QuerySelectBuilder() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}
