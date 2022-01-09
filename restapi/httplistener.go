package restapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/reflect_helper"
	authControllers "github.com/cjlapao/common-go/identity/controllers"
	"github.com/cjlapao/common-go/identity/interfaces"
	"github.com/cjlapao/common-go/identity/middleware"
	logger "github.com/cjlapao/common-go/log"
	"github.com/gorilla/handlers"
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
	Router            *mux.Router
	Context           *execution_context.Context
	Logger            *logger.Logger
	Options           *HttpListenerOptions
	Controllers       []controllers.Controller
	DefaultAdapters   []controllers.Adapter
	Servers           []*http.Server
	shutdownRequest   chan bool
	shutdownRequested uint32
}

var globalHttpListener *HttpListener

// NewHttpListener  Creates a new controller
func NewHttpListener() *HttpListener {
	if globalHttpListener != nil {
		globalHttpListener = nil
		if len(globalHttpListener.Servers) > 0 {
			globalHttpListener.shutdownRequest <- true
		}
	}

	listener := HttpListener{
		Context: execution_context.Get(),
		Router:  mux.NewRouter().StrictSlash(true),
		Servers: make([]*http.Server, 0),
	}

	listener.shutdownRequest = make(chan bool)
	listener.Logger = listener.Context.Services.Logger

	listener.Controllers = make([]controllers.Controller, 0)
	listener.DefaultAdapters = make([]controllers.Adapter, 0)

	// Appending the correlationId renewal
	listener.DefaultAdapters = append(listener.DefaultAdapters, CorrelationMiddlewareAdapter())

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
	l.DefaultAdapters = append(l.DefaultAdapters, LoggerMiddlewareAdapter())
	return l
}

func (l *HttpListener) AddJsonContent() *HttpListener {
	l.DefaultAdapters = append(l.DefaultAdapters, JsonContentMiddlewareAdapter())
	return l
}

func (l *HttpListener) AddDefaultHomepage() *HttpListener {
	return l
}

func (l *HttpListener) WithDefaultAuthentication() *HttpListener {
	if l.Options.UseAuthBackend {
		l.Logger.Info("Found MongoDB connection string, enabling MongoDb auth backend...")
	}
	defaultAuthControllers := authControllers.NewDefaultAuthorizationControllers()

	l.AddController(defaultAuthControllers.Token(), "/auth/token", "POST")
	l.AddController(defaultAuthControllers.Token(), "/auth/{tenantId}/token", "POST")
	l.AddController(defaultAuthControllers.Introspection(), "/auth/token/introspect", "POST")
	l.AddController(defaultAuthControllers.Introspection(), "/auth/{tenantId}/token/introspect", "POST")
	l.AddController(defaultAuthControllers.Introspection(), "/auth/register", "POST")
	l.AddController(defaultAuthControllers.Introspection(), "/auth/{tenantId}/register", "POST")

	l.AddController(defaultAuthControllers.Configuration(), "/auth/.well-known/openid-configuration", "GET")
	l.AddController(defaultAuthControllers.Configuration(), "/auth/{tenantId}/.well-known/openid-configuration", "GET")
	l.AddController(defaultAuthControllers.Jwks(), "/auth/.well-known/openid-configuration/jwks", "GET")
	l.AddController(defaultAuthControllers.Jwks(), "/auth/{tenantId}/.well-known/openid-configuration/jwks", "GET")
	l.DefaultAdapters = append([]controllers.Adapter{middleware.EndAuthorizationMiddlewareAdapter()}, l.DefaultAdapters...)
	l.Options.EnableAuthentication = true
	return l
}

func (l *HttpListener) WithAuthentication(context interfaces.UserDatabaseAdapter) *HttpListener {
	ctx := execution_context.Get()
	if ctx.Authorization != nil {
		ctx.Authorization.ContextAdapter = context
		if l.Options.UseAuthBackend {
			l.Logger.Info("Found MongoDB connection string, enabling MongoDb auth backend...")
		}
		defaultAuthControllers := authControllers.NewAuthorizationControllers(context)

		l.AddController(defaultAuthControllers.Token(), "/auth/token", "POST")
		l.AddController(defaultAuthControllers.Token(), "/auth/{tenantId}/token", "POST")
		l.AddController(defaultAuthControllers.Introspection(), "/auth/token/introspect", "POST")
		l.AddController(defaultAuthControllers.Introspection(), "/auth/{tenantId}/token/introspect", "POST")
		l.AddAuthorizedControllerWithRoles(defaultAuthControllers.Register(), "/auth/register", []string{"_admin"}, []string{}, "POST")
		l.AddAuthorizedControllerWithRoles(defaultAuthControllers.Register(), "/auth/{tenantId}/register", []string{"_admin"}, []string{}, "POST")

		l.AddController(defaultAuthControllers.Configuration(), "/auth/.well-known/openid-configuration", "GET")
		l.AddController(defaultAuthControllers.Configuration(), "/auth/{tenantId}/.well-known/openid-configuration", "GET")
		l.AddController(defaultAuthControllers.Jwks(), "/auth/.well-known/openid-configuration/jwks", "GET")
		l.AddController(defaultAuthControllers.Jwks(), "/auth/{tenantId}/.well-known/openid-configuration/jwks", "GET")
		l.DefaultAdapters = append([]controllers.Adapter{middleware.EndAuthorizationMiddlewareAdapter()}, l.DefaultAdapters...)
		l.Options.EnableAuthentication = true
	} else {
		l.Logger.Error("No authorization context found, ignoring")
	}
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

	if l.Options.ApiPrefix != "" {
		path = l.Options.ApiPrefix + path
	}
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
	adapters = append(adapters, middleware.AuthorizationMiddlewareAdapter([]string{}, []string{}))

	if l.Options.ApiPrefix != "" {
		path = l.Options.ApiPrefix + path
	}

	subRouter.HandleFunc(path,
		controllers.Adapt(
			http.HandlerFunc(c),
			adapters...).ServeHTTP)
}

func (l *HttpListener) AddAuthorizedControllerWithRoles(c controllers.Controller, path string, roles []string, claims []string, methods ...string) {
	l.Controllers = append(l.Controllers, c)
	var subRouter *mux.Router
	if len(methods) > 0 {
		subRouter = l.Router.Methods(methods...).Subrouter()
	} else {
		subRouter = l.Router.Methods("GET").Subrouter()
	}
	adapters := make([]controllers.Adapter, 0)
	adapters = append(adapters, l.DefaultAdapters...)
	adapters = append(adapters, middleware.AuthorizationMiddlewareAdapter(roles, claims))

	if l.Options.ApiPrefix != "" {
		path = l.Options.ApiPrefix + path
	}

	subRouter.HandleFunc(path,
		controllers.Adapt(
			http.HandlerFunc(c),
			adapters...).ServeHTTP)
}

func (l *HttpListener) Start() {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "authorization", "Authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	l.Logger.Notice("Starting %v Go Rest API v%v", l.Context.Services.Version.Name, l.Context.Services.Version.String())

	done := make(chan bool)

	l.Router.HandleFunc(l.Options.ApiPrefix+"/", defaultHomepageController)
	l.Router.HandleFunc(l.Options.ApiPrefix+"/shutdown", globalHttpListener.ShutdownHandler)

	// Creating and starting the http server
	srv := &http.Server{
		Addr:    ":" + l.Options.HttpPort,
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(l.Router),
	}

	l.Servers = append(l.Servers, srv)

	go func() {
		l.Logger.Info("Api listening on http://::" + l.Options.HttpPort + l.Options.ApiPrefix)
		l.Logger.Success("Finished Initiating http server")
		if err := srv.ListenAndServe(); err != nil {
			if !strings.Contains(err.Error(), "http: Server closed") {
				l.Logger.Error("There was an error shutting down the http server: %v", err.Error())
			}
		}
		done <- true
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
					if !strings.Contains(err.Error(), "http: Server closed") {
						l.Logger.Error("There was an error shutting down the https server: %v", err.Error())
					}
				}
				done <- true
			}()
		} else {
			l.Logger.Error("There was an error reading the certificates to enable HTTPS")
		}
	}

	l.WaitAndShutdown()
	<-done

	l.Logger.Info("Server shut down successfully...")
}

func (l *HttpListener) WaitAndShutdown() {
	irqSign := make(chan os.Signal, 1)
	signal.Notify(irqSign, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-irqSign:
		l.Logger.Info("Server shutdown requested (signal: %v)", sig.String())
	case sig := <-l.shutdownRequest:
		l.Logger.Info("Server shutdown requested (/shutdown: %v)", fmt.Sprintf("%v", sig))
	}

	l.Logger.Info("Stopping the server...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Create shutdown context with 10 second timeout
	for _, s := range l.Servers {
		err := s.Shutdown(ctx)
		if err != nil {
			l.Logger.Error("Shutdown request error: %v", err.Error())
		}
	}
}

//region Private Methods
func (l *HttpListener) getDefaultConfiguration() *HttpListenerOptions {
	options := HttpListenerOptions{
		HttpPort:                l.Context.Configuration.GetString("HTTP_PORT"),
		EnableTLS:               l.Context.Configuration.GetBool("ENABLE_TLS"),
		TLSPort:                 l.Context.Configuration.GetString("TLS_PORT"),
		TLSCertificate:          l.Context.Configuration.GetBase64("TLS_CERTIFICATE"),
		TLSPrivateKey:           l.Context.Configuration.GetBase64("TLS_PRIVATE_KEY"),
		DatabaseName:            l.Context.Configuration.GetString("MONGODB_DATABASENAME"),
		MongoDbConnectionString: l.Context.Configuration.GetBase64("MONGODB_CONNECTION_STRING"),
	}

	if reflect_helper.IsNilOrEmpty(options.HttpPort) {
		options.HttpPort = "5000"
	}

	if reflect_helper.IsNilOrEmpty(options.TLSPort) {
		options.TLSPort = "5001"
	}

	if reflect_helper.IsNilOrEmpty(options.DatabaseName) {
		options.DatabaseName = "users"
	}

	apiPrefix := l.Context.Configuration.GetString("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/"
	}

	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = fmt.Sprintf("/%v", apiPrefix)
	}

	options.ApiPrefix = apiPrefix

	l.Options = &options

	return l.Options
}

func defaultHomepageController(w http.ResponseWriter, r *http.Request) {
	response := DefaultHomepage{
		CorrelationID: globalHttpListener.Context.CorrelationId,
		Timestamp:     fmt.Sprint(time.Now().Format(time.RFC850)),
	}

	json.NewEncoder(w).Encode(response)
}

func getDefaultBaseUrl() {

}

//endregion
