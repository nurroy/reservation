package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"belajar/reservation/pkg/v1/converter"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ConfigStore interface {
	GetKey(string) string
}

// Config is the struct to pass the Postgres configuration.
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

const POSTGRES string = "postgres"

var pgConn string    //global variable postgres connection


var reservationDB *sqlx.DB

var (
	errFilterFieldNotFound  = errors.New("filter where field not found")
	errNilData              = errors.New("Error data is nil")
	errStructIsRequired     = errors.New("Struct field is required for this function")
	errLogicalFilterInvalid = errors.New("logical filter in where clause invalid")
)

type ResInsertInfo struct {
	LastInsertId int64
	RowsAffected int64
}

type builderArgs struct {
	value    interface{}
	coalesce string
}


// New initializes the postgres writers
func New(config ConfigStore) error {
	// var err error
fmt.Println("POSTGRES")
	// set postgres och config from env
	pgConfig := PostgresConfig{
		Host:     config.GetKey("postgres.host"),
		Port:     config.GetKey("postgres.port"),
		User:     config.GetKey("postgres.user"),
		Password: config.GetKey("postgres.password"),
		Dbname:   config.GetKey("postgres.dbname"),
	}



	// Set postgres och config
	pgConn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.Dbname)

	dbOch, err := sql.Open(POSTGRES, pgConn)
	if err != nil {
		return err
	}
	fmt.Println("POSTGRES2")
	// Test postgres OCH connection
	sqlQueryOch := fmt.Sprintf(SELECT_MENU_LIMIT1)
	rowsOch, err := dbOch.Query(sqlQueryOch)
	defer dbOch.Close()
	if err != nil {
		return err
	}
	defer rowsOch.Close()

	// Initiate sqlx pinangDB Connection
	reservationDB, err = sqlx.Connect(POSTGRES, pgConn)
	if err != nil {
		return err
	}

	return nil
}

// QueryUpdateBuilder generate update query from map
//	sample:
// tableName -> tabelku
//	mapArgs -> map[string]interface{}{"cobafield": 1,"primary": "abc"}
// wheres -> []string{"primary"}
//
// result:
//	query -> "UPDATE tabelku SET cobafield = $1 WHERE primary = $2"
// args -> []interface{}{1, "abc"}
func QueryUpdateBuilder(tableName string, mapArgs map[string]interface{}, wheres []string, useDBTS bool) (query string, args []interface{}, err error) {
	query = "UPDATE " + tableName + " SET "

	whereMap := make(map[string]bool)
	for _, where := range wheres {
		whereClean := strings.Replace(where, "and", "", -1)
		whereClean = strings.Replace(where, "or", "", -1)
		whereClean = strings.Replace(where, " ", "", -1)
		whereMap[whereClean] = true
	}

	prefix := ""
	if useDBTS {
		prefix = "db_ts=db_ts+1, "
	}

	identifierNumberMap := make(map[string]int)
	identifier := 1

	// sorted map
	keys := make([]string, 0, len(mapArgs))
	for key := range mapArgs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// don't update filter value
		if _, ok := whereMap[key]; ok {
			continue
		}
		id, ok := identifierNumberMap[key]
		if !ok {
			identifierNumberMap[key] = identifier
			id = identifier
			identifier++
		}

		query += prefix + key + " = " + fmt.Sprintf("$%d", id)
		args = append(args, mapArgs[key])

		prefix = ", "
	}

	prefix = ""
	query += " WHERE "
	for _, where := range wheres {
		whereClean := strings.Replace(where, "and", "", -1)
		whereClean = strings.Replace(where, "or", "", -1)
		whereClean = strings.Replace(where, " ", "", -1)
		id, ok := identifierNumberMap[where]
		if !ok {
			identifierNumberMap[whereClean] = identifier
			id = identifier
			identifier++
		}

		if _, ok := mapArgs[whereClean]; !ok {
			return query, args, errFilterFieldNotFound
		}

		query += prefix + where + " = " + fmt.Sprintf("$%d", id)
		args = append(args, mapArgs[whereClean])
		prefix = " and "
	}

	return
}

// QuerySelectBuilder query builder for select statement, query result field will be sorted ascending
// please follow the sorting to scan the result to the struct after using the query
func QuerySelectBuilder(tableName string, structData interface{}, wheres []string) (query string, err error) {
	query = "SELECT "

	// get the field from struct
	mapArgs, err := convertFieldToMapArgs(structData)
	if err != nil {
		return "", err
	}

	// sorted map
	keys := make([]string, 0, len(mapArgs))
	coalesce := make(map[string]string, len(mapArgs))
	for key, args := range mapArgs {
		keys = append(keys, key)
		coalesce[key] = args.coalesce
	}
	sort.Strings(keys)

	// join the fields and combine with current query and table name
	prefix := ""
	for i := range keys {
		key := keys[i]
		if c, ok := coalesce[key]; ok && c != "" {
			key = c
		}
		query += prefix + key
		prefix = ", "
	}

	query += fmt.Sprintf(" FROM %s", tableName)
	if wheres == nil || len(wheres) == 0 {
		return query, nil
	}

	whereMap := make(map[string]bool)
	for _, where := range wheres {
		whereClean := strings.Replace(where, "and", "", -1)
		whereClean = strings.Replace(where, "or", "", -1)
		whereClean = strings.Replace(where, " ", "", -1)
		whereMap[whereClean] = true
	}

	identifierNumberMap := make(map[string]int)
	identifier := 1

	query += " WHERE "
	prefix = ""
	for _, where := range wheres {
		whereClean := strings.Replace(where, "and", "", -1)
		whereClean = strings.Replace(where, "or", "", -1)
		whereClean = strings.Replace(where, " ", "", -1)
		id, ok := identifierNumberMap[where]
		if !ok {
			identifierNumberMap[whereClean] = identifier
			id = identifier
			identifier++
		}

		if _, ok := mapArgs[whereClean]; !ok {
			return query, errFilterFieldNotFound
		}

		query += prefix + where + " = " + fmt.Sprintf("$%d", id)
		prefix = " and "
	}

	return
}

func convertFieldToMapArgs(data interface{}) (args map[string]builderArgs, err error) {
	if data == nil {
		return args, errNilData
	}
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		return args, errStructIsRequired
	}

	fields := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	values = reflect.Indirect(values)

	args = make(map[string]builderArgs)
	for i := 0; i < fields.NumField(); i++ {
		var resultCoalesce string
		var resultValue interface{}

		key := fields.Field(i).Tag.Get("db")
		if key == "" || key == "-" {
			continue
		}

		// checking if the field has COALESCE
		coalesceValue, coalesce := fields.Field(i).Tag.Lookup("coalesce")
		if coalesceValue == "" || coalesceValue == "-" {
			coalesceValue = ""
		}

		fieldValue := values.Field(i)
		// set based on REAL type
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			if coalesce {
				c, _ := strconv.ParseInt(coalesceValue, 10, 64)
				resultCoalesce = fmt.Sprintf("coalesce(%s, %v) as %s", key, c, key)
			}
			// convert to real int value
			int64Value := fieldValue.Int()
			if fieldValue.Kind() == reflect.Int {
				resultValue = int(int64Value)
			} else if fieldValue.Kind() == reflect.Int16 {
				resultValue = int16(int64Value)
			} else if fieldValue.Kind() == reflect.Int32 {
				resultValue = int32(int64Value)
			} else {
				resultValue = int64Value
			}
		case reflect.Float32, reflect.Float64:
			if coalesce {
				c := converter.StringToFloat64(coalesceValue)
				resultCoalesce = fmt.Sprintf("coalesce(%s, %v) as %s", key, c, key)
			}
			resultValue = fieldValue.Float()
			float64Value := fieldValue.Float()
			if fieldValue.Kind() == reflect.Float32 {
				resultValue = float32(float64Value)
			} else {
				resultValue = float64Value
			}
		case reflect.Bool:
			if coalesce {
				c, _ := strconv.ParseBool(coalesceValue)
				resultCoalesce = fmt.Sprintf("coalesce(%s, %v) as %s", key, c, key)
			}
			resultValue = fieldValue.Bool()
		case reflect.String:
			if coalesce {
				resultCoalesce = fmt.Sprintf("coalesce(%s, '%s') as %s", key, coalesceValue, key)
			}
			resultValue = fieldValue.String()
		}

		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			if coalesce {
				resultCoalesce = fmt.Sprintf("coalesce(%s, now()) as %s", key, key)
			}
			resultValue = fieldValue.String()
		}

		args[key] = builderArgs{
			value:    resultValue,
			coalesce: resultCoalesce,
		}
	}
	return
}
