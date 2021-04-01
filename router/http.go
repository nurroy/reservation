package router

import (
	"belajar/reservation/pkg/v1/header"
	"belajar/reservation/pkg/v1/utils/errors"
	"context"
	"fmt"
	"github.com/kenshaw/envcfg"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	mnpb "belajar/reservation/proto/v1/menu"
)

// NewHTTPServer creates the http server serve mux.
func NewHTTPServer(config *envcfg.Envcfg, logger *logrus.Logger) error {
	gwruntime.HTTPError = errors.CustomHTTPError

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	addr := "0.0.0.0:" + config.GetKey("grpc.port")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return err
	}
	defer conn.Close()
	// Create new grpc-gateway
	rmux := gwruntime.NewServeMux(gwruntime.WithForwardResponseOption(header.HttpResponseModifier))

	// register gateway endpoints
	for _, f := range []func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error{
		// register grpc service handler
		//hlpb.RegisterHealthServiceHandler,
		// register grpc service handler
		mnpb.RegisterMenuHandler,
	} {
		if err = f(ctx, rmux, conn); err != nil {
			log.Fatal(err)
			return err
		}
	}

	// create http server mux
	mux := http.NewServeMux()
	mux.Handle("/", rmux)



	// run swagger server
	if config.GetKey("runtime.environment") == "development" {
		CreateSwagger(mux)
	}

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "X-CSRF-Token", "Authorization", "Timezone-Offset"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// running rest http server
	log.Println("[SERVER] REST HTTP server is ready")
	err = http.ListenAndServe("0.0.0.0:"+config.PortString(), handlers.CORS(headersOk, originsOk, methodsOk)(mux))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// CreateSwagger creates the swagger server serve mux.
func CreateSwagger(gwmux *http.ServeMux) {
	// register swagger service server
	gwmux.HandleFunc("/corp/rest/v1/banks/docs.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger/docs.json")
	})

	// load swagger-ui file
	fs := http.FileServer(http.Dir("swagger/swagger-ui"))
	gwmux.Handle("/corp/rest/v1/banks/docs/", http.StripPrefix("/corp/rest/v1/banks/docs", fs))
}

// logHandler is log middleware
func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Println(fmt.Sprintf("request:\n %q", x))
		rec := httptest.NewRecorder()
		fn(rec, r)
		log.Println("response: ")
		log.Println(fmt.Sprintf("Response Code: %v", rec.Code))
		log.Println(fmt.Sprintf("Headers: {%q}", rec.Header()))
		log.Println(fmt.Sprintf("Response: %q", rec.Body))

		// this copies the recorded response to the response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}
