package execution_context

import (
	"os"
	"strings"

	"github.com/cjlapao/common-go/cache"
	"github.com/cjlapao/common-go/cache/jwt_token_cache"
	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/service_provider"
	"github.com/google/uuid"
)

var contextService *Context

// Context entity
type Context struct {
	Configuration *configuration.ConfigurationService
	Services      *service_provider.ServiceProvider
	Caches        *cache.CacheService
	TokenCache    *jwt_token_cache.JwtTokenCacheProvider
	Authorization *authorization_context.AuthorizationContext
	CorrelationId string
	Environment   string
	IsDevelopment bool
	Debug         bool
	Init          func() error
}

func New() (*Context, error) {
	if contextService != nil {
		contextService = nil
	}

	return InitNewContext(nil)
}

func InitNewContext(init func() error) (*Context, error) {
	contextService = &Context{
		IsDevelopment: false,
		Debug:         false,
		Init:          init,
	}

	contextService.Caches = cache.Get()
	contextService.TokenCache = jwt_token_cache.New()
	contextService.CorrelationId = uuid.NewString()
	contextService.Services = service_provider.Get()

	environment := os.Getenv("CJ_ENVIRONMENT")
	debug := os.Getenv("CJ_ENABLE_DEBUG")

	if !helper.IsNilOrEmpty(environment) {
		if strings.ToLower(environment) == "development" {
			contextService.IsDevelopment = true
			contextService.Environment = "Development"
		} else {
			contextService.IsDevelopment = false
			switch strings.ToLower(environment) {
			case "production":
				contextService.Environment = "Production"
			case "release":
				contextService.Environment = "Release"
			case "ci":
				contextService.IsDevelopment = true
				contextService.Environment = "CI"
			case "devprod":
				contextService.Environment = "DevProd"
			default:
				contextService.Environment = "Production"
			}
		}
	} else {
		contextService.IsDevelopment = false
		contextService.Environment = "Production"
	}

	if !helper.IsNilOrEmpty(debug) && strings.ToLower(debug) == "true" {
		contextService.Debug = true
	} else {
		contextService.Debug = false
	}

	contextService.Configuration = configuration.Get()

	if contextService.Init != nil {
		err := contextService.Init()
		if err != nil {
			contextService = nil
			return contextService, err
		}
	}

	// Authorization Context
	contextService.Authorization = authorization_context.New()

	return contextService, nil
}

func Get() *Context {
	if contextService != nil {
		return contextService
	}

	New()

	return contextService
}
