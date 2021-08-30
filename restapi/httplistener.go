package restapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"net/http"

	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/helper"
	logger "github.com/cjlapao/common-go/log"
	"github.com/gorilla/mux"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type HttpListenerOptions struct {
	ApiPrefix               string
	HttpPort                string
	EnableTLS               bool
	TLSPort                 string
	TLSCertificate          string
	TLSPrivateKey           string
	UseAuthBackend          bool
	MongoDbConnectionString string
	DatabaseName            string
	EnableAuthentication    bool
}

// HttpListener HttpListener structure
type HttpListener struct {
	Router      *mux.Router
	Services    *executionctx.ServiceProvider
	Logger      *logger.Logger
	Options     *HttpListenerOptions
	Controllers []Controller
}

var globalHttpListener *HttpListener

// NewHttpListener  Creates a new controller
func NewHttpListener() *HttpListener {
	if globalHttpListener != nil {
		globalHttpListener = nil
	}

	listener := HttpListener{
		Services: executionctx.GetServiceProvider(),
		Router:   mux.NewRouter().StrictSlash(true),
	}

	listener.Logger = listener.Services.Logger

	listener.Controllers = make([]Controller, 0)

	listener.Options = listener.getDefaultConfiguration()

	globalHttpListener = &listener
	return globalHttpListener
}

func GetHttpListener() *HttpListener {
	if globalHttpListener != nil {
		return globalHttpListener
	}

	return NewHttpListener()
}

func (l *HttpListener) AddAuthentication() *HttpListener {
	if l.Options.UseAuthBackend {
		l.Logger.Info("Found MongoDB connection string, enabling MongoDb auth backend...")
	}

	l.AddController(l.Login, l.Options.ApiPrefix+"/login", "POST")
	l.AddController(l.Validate, l.Options.ApiPrefix+"/validate", "GET")
	l.Options.EnableAuthentication = true
	return l
}

func (l *HttpListener) AddHealthCheck() *HttpListener {
	l.AddController(l.Probe, l.Options.ApiPrefix+"/probe", "GET")
	l.Options.EnableAuthentication = true
	return l
}

func (l *HttpListener) AddLogger() *HttpListener {
	l.Router.Use(loggerMiddleware)

	return l
}

func (l *HttpListener) AddJsonContent() *HttpListener {
	l.Router.Use(jsonResponseMiddleware)

	return l
}

func (l *HttpListener) AddDefaultHomepage() *HttpListener {
	l.Router.Use(loggerMiddleware)

	return l
}

func (l *HttpListener) AddController(c Controller, path string, methods ...string) {
	l.Controllers = append(l.Controllers, c)
	if len(methods) > 0 {
		l.Router.HandleFunc(path, c).Methods(methods...)
	} else {
		l.Router.HandleFunc(path, c).Methods("GET")
	}
}

func (l *HttpListener) Start() {
	l.Logger.Notice("Starting %v Go Rest API v%v", l.Services.Version.Name, l.Services.Version.String())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	l.Router.HandleFunc(l.Options.ApiPrefix+"/", defaultHomepageController)

	// Creating and starting the http server
	srv := &http.Server{
		Addr:    ":" + l.Options.HttpPort,
		Handler: l.Router,
	}

	go func() {
		l.Logger.Info("Api listening on http://::" + l.Options.HttpPort + l.Options.ApiPrefix)
		l.Logger.Success("Finished Initiating http server")
		if err := srv.ListenAndServe(); err != nil {
			l.Logger.FatalError(err, "There was an error")
		}
	}()

	if l.Options.EnableTLS {
		cert, err := tls.X509KeyPair([]byte(l.Options.TLSCertificate), []byte(l.Options.TLSPrivateKey))
		if err == nil {
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			sslSrv := &http.Server{
				Addr:      ":" + l.Options.TLSPort,
				TLSConfig: tlsConfig,
				Handler:   l.Router,
			}

			go func() {
				l.Logger.Info("Api listening on https://::" + l.Options.TLSPort + l.Options.ApiPrefix)
				l.Logger.Success("Finished Initiating https server")
				if err := sslSrv.ListenAndServeTLS("", ""); err != nil {
					l.Logger.FatalError(err, "There was an error")
				}
			}()
		} else {
			l.Logger.Error("There was an error reading the certificates to enable HTTPS")
		}
	}

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	l.Logger.Info("Server shutdown")
	// http.ListenAndServeTLS(":10001", "./ssl/local-cluster.internal.crt", "./ssl/local-cluster.internal.key", router)
}

//region Private Methods
func (l *HttpListener) getDefaultConfiguration() *HttpListenerOptions {
	options := HttpListenerOptions{
		ApiPrefix:               l.Services.Context.Configuration.GetString("API_PREFIX"),
		HttpPort:                l.Services.Context.Configuration.GetString("HTTP_PORT"),
		EnableTLS:               l.Services.Context.Configuration.GetBool("ENABLE_TLS"),
		TLSPort:                 l.Services.Context.Configuration.GetString("TLS_PORT"),
		TLSCertificate:          l.Services.Context.Configuration.GetBase64("TLS_CERTIFICATE"),
		TLSPrivateKey:           l.Services.Context.Configuration.GetBase64("TLS_PRIVATE_KEY"),
		DatabaseName:            l.Services.Context.Configuration.GetString("MONGODB_DATABASENAME"),
		MongoDbConnectionString: l.Services.Context.Configuration.GetBase64("MONGODB_CONNECTION_STRING"),
	}

	if helper.IsNilOrEmpty(options.HttpPort) {
		options.HttpPort = "5000"
	}

	if helper.IsNilOrEmpty(options.TLSPort) {
		options.TLSPort = "5001"
	}

	if helper.IsNilOrEmpty(options.DatabaseName) {
		options.DatabaseName = "users"
	}

	if !helper.IsNilOrEmpty(options.MongoDbConnectionString) {
		options.UseAuthBackend = true
	}

	l.Options = &options

	return l.Options
}

func defaultHomepageController(w http.ResponseWriter, r *http.Request) {
	response := DefaultHomepage{
		CorrelationID: globalHttpListener.Services.Context.CorrelationId,
		Timestamp:     fmt.Sprint(time.Now().Format(time.RFC850)),
	}

	json.NewEncoder(w).Encode(response)
}

func jsonResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		globalHttpListener.Logger.Info("[%v] %v from %v", r.Method, r.URL.Path, r.Host)
		next.ServeHTTP(w, r)
	})
}

//endregion
