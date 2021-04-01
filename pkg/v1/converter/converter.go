package converter

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"reflect"
	"strconv"
)

func StringToFloat64(input string) float64 {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return -1
	}
	return f
}

// ComputeHMAC256 computes the HMAC256 hash of the message.
func GenerateHMAC256(k, message string) string {
	key := []byte(k)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Func to convert interface{} to int32
func InterfaceToInt32(data interface{}) int32 {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	switch typ.Kind() {
	case reflect.String:
		dataInt, _ := strconv.ParseInt(val.String(), 10, 64)
		return int32(dataInt)
	case reflect.Float64:
		return int32(data.(float64))
	case reflect.Float32:
		return int32(data.(float32))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return int32(val.Int())
	default:
		return int32(0)
	}
}

// InterfaceToInt64 Func to convert interface{} to int64
func InterfaceToInt64(data interface{}) int64 {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	switch typ.Kind() {
	case reflect.String:
		dataInt, _ := strconv.ParseInt(val.String(), 10, 64)
		return dataInt
	case reflect.Float64:
		return int64(data.(float64))
	case reflect.Float32:
		return int64(data.(float32))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return val.Int()
	default:
		return int64(0)
	}
}

// NewNullString Func to convert empty string to Null String SQL
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// Int32ToString Func to convert int32 to string
func Int32ToString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
