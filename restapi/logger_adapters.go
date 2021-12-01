package restapi

import (
	"net/http"

	"github.com/cjlapao/common-go/controllers"
)

func LoggerAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			globalHttpListener.Logger.Info("[%v] %v from %v", r.Method, r.URL.Path, r.Host)
			next.ServeHTTP(w, r)
		})
	}
}
