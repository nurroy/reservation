package header

import (
	"context"
	"github.com/golang/protobuf/proto"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net/http"
	"strconv"
)

type Message struct {
	Message_TYPE     string `json:"message_TYPE,omitempty"`
	MessageDesc      string `json:"messageDesc,omitempty"`
	MessageAddlnInfo string `json:"messageAddlnInfo,omitempty"`
	NameSpaceInfo    string `json:"nameSpaceInfo,omitempty"`
	MessageCode      string `json:"messageCode,omitempty"`
}

type Pagination struct {
	HasMoreRecords string `json:"hasMoreRecords,omitempty"`
	NumRecReturned string `json:"numRecReturned,omitempty"`
}

type Status struct {
	Message []Message `json:"message,omitempty"`
}

type ChannelCtx struct {
	Status     Status     `json:"status,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type ChannelCtx2 struct {
	Status Status `json:"status,omitempty"`
}

func HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := gwruntime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	Header(w)

	if vals := md.HeaderMD.Get("ChannelContext"); len(vals) > 0 {
		chanCtx := vals[0]
		w.Header().Set("ChannelContext", chanCtx)
		delete(w.Header(), "Grpc-Metadata-Channelcontext")
		delete(w.Header(), "Grpc-Metadata-Content-Type")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	}
	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
	}

	// delete the headers to not expose any grpc-metadata in http response
	delete(w.Header(), "Grpc-Metadata-Cache-Control")
	delete(w.Header(), "Grpc-Metadata-X-Powered-By")
	delete(w.Header(), "Grpc-Metadata-Access-Control-Allow-Headers")
	delete(w.Header(), "Grpc-Metadata-Server")
	delete(w.Header(), "Grpc-Metadata-Pragma")
	delete(w.Header(), "Grpc-Metadata-Access-Control-Expose-Headers")
	delete(w.Header(), "Grpc-Metadata-Access-Control-Allow-Origin")
	delete(w.Header(), "Grpc-Metadata-Access-Control-Allow-Methods")
	delete(w.Header(), "Grpc-Metadata-ChannelContext")
	delete(w.Header(), "Grpc-Metadata-Connection")

	return nil
}

// Header for http response
func Header(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("X-Powered-By", "Undertow/1")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Accept,Authorization,OCH_CATEGORY_ID,OCH_RANDOM_KEY,OCH_KEY_ID")
	w.Header().Set("Server", "PINUAT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Access-Control-Expose-Headers", "ChannelContext")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")
	w.Header().Set("Connection", "close")
}
