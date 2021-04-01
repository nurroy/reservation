// Package errors is the custom error handler for BRI Open API.
package errors

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/metadata"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Success bool   `json:"success,omitempty"`
	Code    string `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

//type ErrorResponeErr struct {
//	Error ErrorResponse `json:"error,omitempty"`
//}

// Method is the method types.
type Method int
type Methode int
type Service int

const (
	FallbackError = `{"status": false, "code": "100", "msg": "internal error"}`
	InternalServerError = `Internal Server Error`

	// List of different Methods
	Authorization Method = iota
	Unauthenticated
	InvalidFormat
	InternalServer
	Unavailable
	Other
)

// CustomHTTPError is the custom http error handler for grpc gateway.
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, ierr error) {
	var desc string

	w.Header().Set("Content-Type", marshaler.ContentType())
	if s, ok := status.FromError(ierr); ok {
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
		desc = s.Message()
	} else {
		w.WriteHeader(runtime.HTTPStatusFromCode(codes.Unknown))
		desc = ierr.Error()
	}
	b := new(ErrorResponse)
	err := json.Unmarshal([]byte(desc), b)
	if err != nil {
		_, _ = w.Write([]byte(FallbackError))
	} else {
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			_, _ = w.Write([]byte(FallbackError))
		}
	}
}

// FormatError is the exposed function for generating errors.
func FormatError(m Method, c codes.Code, success bool, params ...string) error {
	// Default
	//errRes := &ErrorResponse{
	//	Success:       	success,
	//	RespCode:       params[0],
	//	RespDesc:       params[1],
	//}
	// Default change to code and msg
	errRes := &ErrorResponse{
		Success: success,
		Code:    params[0],
		Msg:     params[1],
	}

	// Custom error
	//errResErr := ErrorResponeErr{
	//	Error: ErrorResponse{
	//		Success: success,
	//		Code:    params[0],
	//		Msg:     params[1],
	//	},
	//}

	buf, err := json.Marshal(errRes)

	if err != nil {
		return status.Errorf(c, FallbackError)
	}

	return status.Errorf(c, string(buf))
}

// FormatErrorEncoded is the exposed function for generating errors.
func FormatErrorArray(w http.ResponseWriter, c codes.Code, params ...string) error {
	var buf []byte
	var err error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(c))

	if len(params) != 2 {
		_, _ = w.Write([]byte(InternalServerError))
		return status.Errorf(c, InternalServerError)
	}

	errRes := &ErrorResponse{
		Success: false,
		Code:    params[0],
		Msg:     params[1],
	}

	buf, err = json.Marshal(errRes)
	if err != nil {
		_, _ = w.Write([]byte(InternalServerError))
		return status.Errorf(c, InternalServerError)
	}

	_, _ = w.Write(buf)

	return status.Errorf(c, string(buf))
}

// FormatError is the exposed function for generating errors.
func FormatErrorHeader(ctx context.Context, md metadata.MD, m Method, c codes.Code, success bool, params ...string) error {
	// Default
	//errRes := &ErrorResponse{
	//	Success:       	success,
	//	RespCode:       params[0],
	//	RespDesc:       params[1],
	//}

	// Default change to code and msg
	errRes := &ErrorResponse{
		Success: success,
		Code:    params[0],
		Msg:     params[1],
	}

	// Custom error
	//errResErr := ErrorResponeErr{
	//	Error: ErrorResponse{
	//		Success: success,
	//		Code:    params[0],
	//		Msg:     params[1],
	//	},
	//}

	buf, err := json.Marshal(errRes)

	if err != nil {
		return status.Errorf(c, FallbackError)
	}

	return status.Errorf(c, string(buf))
}

// func to create error from string
func New(errs string) error {
	return errors.New(errs)
}

// WrapErr is the function wrapper to wrap the error
func WrapErr(err error, message string, args ...interface{}) error {
	return errors.Wrapf(err, message, args...)
}
