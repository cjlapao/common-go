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

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/identity"
	logger "github.com/cjlapao/common-go/log"
	"github.com/gorilla/mux"
)

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
	Router          *mux.Router
	Services        *executionctx.ServiceProvider
	Logger          *logger.Logger
	Options         *HttpListenerOptions
	Controllers     []controllers.Controller
	DefaultAdapters []controllers.Adapter
	Servers         []*http.Server
}

var globalHttpListener *HttpListener

// NewHttpListener  Creates a new controller
func NewHttpListener() *HttpListener {
	if globalHttpListener != nil {
		globalHttpListener = nil
		if len(globalHttpListener.Servers) > 0 {
			for _, srv := range globalHttpListener.Servers {
				srv.Close()
			}
		}
	}

	listener := HttpListener{
		Services: executionctx.GetServiceProvider(),
		Router:   mux.NewRouter().StrictSlash(true),
		Servers:  make([]*http.Server, 0),
	}

	listener.Logger = listener.Services.Logger

	listener.Controllers = make([]controllers.Controller, 0)
	listener.DefaultAdapters = make([]controllers.Adapter, 0)

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

func (l *HttpListener) AddHealthCheck() *HttpListener {
	l.AddController(l.Probe(), l.Options.ApiPrefix+"/probe", "GET")
	return l
}

func (l *HttpListener) AddLogger() *HttpListener {
	l.DefaultAdapters = append(l.DefaultAdapters, LoggerAdapter())
	return l
}

func (l *HttpListener) AddJsonContent() *HttpListener {
	l.DefaultAdapters = append(l.DefaultAdapters, JsonContentAdapter())
	return l
}

func (l *HttpListener) AddDefaultHomepage() *HttpListener {
	return l
}

func (l *HttpListener) WithDefaultAuthentication() *HttpListener {
	if l.Options.UseAuthBackend {
		l.Logger.Info("Found MongoDB connection string, enabling MongoDb auth backend...")
	}
	defaultAuthControllers := identity.NewDefaultAuthorizationControllers()

	l.AddController(defaultAuthControllers.Login(), l.Options.ApiPrefix+"/login", "POST")
	l.AddController(defaultAuthControllers.Validate(), l.Options.ApiPrefix+"/validate", "GET")
	l.DefaultAdapters = append([]controllers.Adapter{identity.EndAuthorizationAdapter()}, l.DefaultAdapters...)
	l.Options.EnableAuthentication = true
	return l
}

func (l *HttpListener) WithAuthentication(context identity.UserContext) *HttpListener {
	if l.Options.UseAuthBackend {
		l.Logger.Info("Found MongoDB connection string, enabling MongoDb auth backend...")
	}
	defaultAuthControllers := identity.NewAuthorizationControllers(context)

	l.AddController(defaultAuthControllers.Login(), l.Options.ApiPrefix+"/login", "POST")
	l.AddController(defaultAuthControllers.Validate(), l.Options.ApiPrefix+"/validate", "GET")
	l.DefaultAdapters = append([]controllers.Adapter{identity.EndAuthorizationAdapter()}, l.DefaultAdapters...)
	l.Options.EnableAuthentication = true
	return l
}

func (l *HttpListener) AddController(c controllers.Controller, path string, methods ...string) {
	l.Controllers = append(l.Controllers, c)
	var subRouter *mux.Router
	if len(methods) > 0 {
		subRouter = l.Router.Methods(methods...).Subrouter()
	} else {
		subRouter = l.Router.Methods("GET").Subrouter()
	}

	adapters := make([]controllers.Adapter, 0)
	adapters = append(adapters, l.DefaultAdapters...)

	subRouter.HandleFunc(path, controllers.Adapt(
		http.HandlerFunc(c),
		adapters...).ServeHTTP)
}

func (l *HttpListener) AddAuthorizedController(c controllers.Controller, path string, methods ...string) {
	l.Controllers = append(l.Controllers, c)
	var subRouter *mux.Router
	if len(methods) > 0 {
		subRouter = l.Router.Methods(methods...).Subrouter()
	} else {
		subRouter = l.Router.Methods("GET").Subrouter()
	}
	adapters := make([]controllers.Adapter, 0)
	adapters = append(adapters, l.DefaultAdapters...)
	adapters = append(adapters, identity.AuthorizationAdapter())
	subRouter.HandleFunc(path,
		controllers.Adapt(
			http.HandlerFunc(c),
			adapters...).ServeHTTP)
}

func (l *HttpListener) Start() {
	l.Logger.Notice("Starting %v Go Rest API v%v", l.Services.Version.Name, l.Services.Version.String())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	l.Router.HandleFunc(l.Options.ApiPrefix+"/", defaultHomepageController)

	// Creating and starting the http server
	srv := &http.Server{
		Addr:    ":" + l.Options.HttpPort,
		Handler: l.Router,
	}

	l.Servers = append(l.Servers, srv)

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

			l.Servers = append(l.Servers, sslSrv)

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

//endregion
