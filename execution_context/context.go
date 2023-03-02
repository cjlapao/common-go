package execution_context

import (
	"os"
	"strings"

	cryptorand "github.com/cjlapao/common-go-cryptorand"
	"github.com/cjlapao/common-go/cache"
	"github.com/cjlapao/common-go/cache/jwt_token_cache"
	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/constants"
	"github.com/cjlapao/common-go/helper/reflect_helper"
	"github.com/cjlapao/common-go/service_provider"
)

var contextService *Context

// Context entity
type Context struct {
	Configuration *configuration.ConfigurationService
	Services      *service_provider.ServiceProvider
	Caches        *cache.CacheService
	TokenCache    *jwt_token_cache.JwtTokenCacheProvider
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

	contextService.Configuration = configuration.Get()
	contextService.Caches = cache.Get()
	contextService.TokenCache = jwt_token_cache.New()
	contextService.CorrelationId = cryptorand.GetRandomString(constants.ID_SIZE)
	contextService.Services = service_provider.Get()
	os.Setenv("CORRELATION_ID", contextService.CorrelationId)

	environment := contextService.Configuration.GetString(constants.ENVIRONMENT)
	debug := contextService.Configuration.GetBool(constants.DEBUG_ENVIRONMENT)

	if !reflect_helper.IsNilOrEmpty(environment) {
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

	if debug {
		contextService.Debug = true
	} else {
		contextService.Debug = false
	}

	if contextService.Init != nil {
		err := contextService.Init()
		if err != nil {
			contextService = nil
			return contextService, err
		}
	}

	return contextService, nil
}

func (c *Context) Refresh() *Context {
	c.CorrelationId = cryptorand.GetRandomString(constants.ID_SIZE)
	os.Setenv("CORRELATION_ID", c.CorrelationId)
	return c
}

func (c *Context) SetCorrelationId(correlationId string) *Context {
	c.CorrelationId = correlationId
	os.Setenv("CORRELATION_ID", c.CorrelationId)
	return c
}

func Get() *Context {
	if contextService != nil {
		return contextService
	}

	New()

	return contextService
}
