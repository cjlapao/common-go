package restapi

import (
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
)

func JsonContentMiddlewareAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}

func CorrelationMiddlewareAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := execution_context.Get()
			logger := ctx.Services.Logger
			ctx.Refresh()
			logger.Info("Http request with correlation %v", ctx.CorrelationId)
			next.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddlewareAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			globalHttpListener.Logger.Info("[%v] %v from %v", r.Method, r.URL.Path, r.Host)
			next.ServeHTTP(w, r)
		})
	}
}
